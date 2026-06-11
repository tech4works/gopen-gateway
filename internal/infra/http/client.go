/*
 * Copyright 2024 Tech4Works
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type client struct {
	engine   *http.Client
	breakers sync.Map
	cfg      *vo.ClientConfig
	log      app.BootLog
}

// NewClient creates an HTTP client with optimized transport, circuit breaker per host,
// and retry for GET/HEAD on transient errors. All settings are derived from the
// server.client JSON config block (populated via env vars).
func NewClient(gopen *dto.Gopen, log app.BootLog) app.HTTPClient {
	cfg := buildClientConfig(gopen)

	transport := &http.Transport{
		MaxIdleConns:        cfg.MaxIdleConns(),
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost(),
		MaxConnsPerHost:     0,
		IdleConnTimeout:     cfg.IdleConnTimeout(),
		DisableKeepAlives:   false,
		ForceAttemptHTTP2:   true,
	}

	log.PrintInfof("HTTP client config: timeout=%s max-idle-conns=%d max-idle-conns-per-host=%d idle-conn-timeout=%s",
		cfg.Timeout(), cfg.MaxIdleConns(), cfg.MaxIdleConnsPerHost(), cfg.IdleConnTimeout())
	log.PrintInfof("HTTP client circuit breaker: failure-threshold=%d success-threshold=%d open-timeout=%s half-open-max-requests=%d",
		cfg.CBFailureThreshold(), cfg.CBSuccessThreshold(), cfg.CBOpenTimeout(), cfg.CBHalfOpenMaxReqs())
	log.PrintInfof("HTTP client retry: max-retries=%d backoff=%s",
		cfg.RetryMaxRetries(), cfg.RetryBackoff())

	return &client{
		engine: &http.Client{
			Transport: otelhttp.NewTransport(transport),
			Timeout:   cfg.Timeout(),
		},
		cfg: cfg,
		log: log,
	}
}

func buildClientConfig(gopen *dto.Gopen) *vo.ClientConfig {
	var timeout, idleConnTimeout vo.Duration
	var maxIdleConns, maxIdleConnsPerHost int
	var cb vo.CircuitBreakerConfig
	var retry vo.RetryConfig

	if checker.NonNil(gopen) && checker.NonNil(gopen.Server) && checker.NonNil(gopen.Server.Client) {
		c := gopen.Server.Client
		if checker.NonNil(c.Timeout) {
			timeout = *c.Timeout
		}
		if checker.NonNil(c.MaxIdleConns) {
			maxIdleConns = *c.MaxIdleConns
		}
		if checker.NonNil(c.MaxIdleConnsPerHost) {
			maxIdleConnsPerHost = *c.MaxIdleConnsPerHost
		}
		if checker.NonNil(c.IdleConnTimeout) {
			idleConnTimeout = *c.IdleConnTimeout
		}
		if checker.NonNil(c.CircuitBreaker) {
			var ft, st, hom int
			var ot vo.Duration
			if checker.NonNil(c.CircuitBreaker.FailureThreshold) {
				ft = *c.CircuitBreaker.FailureThreshold
			}
			if checker.NonNil(c.CircuitBreaker.SuccessThreshold) {
				st = *c.CircuitBreaker.SuccessThreshold
			}
			if checker.NonNil(c.CircuitBreaker.OpenTimeout) {
				ot = *c.CircuitBreaker.OpenTimeout
			}
			if checker.NonNil(c.CircuitBreaker.HalfOpenMaxReqs) {
				hom = *c.CircuitBreaker.HalfOpenMaxReqs
			}
			cb = vo.NewCircuitBreakerConfig(ft, st, ot, hom)
		}
		if checker.NonNil(c.Retry) {
			var mr int
			var b vo.Duration
			if checker.NonNil(c.Retry.MaxRetries) {
				mr = *c.Retry.MaxRetries
			}
			if checker.NonNil(c.Retry.Backoff) {
				b = *c.Retry.Backoff
			}
			retry = vo.NewRetryConfig(mr, b)
		}
	}

	return vo.NewClientConfig(timeout, maxIdleConns, maxIdleConnsPerHost, idleConnTimeout, cb, retry)
}

func (c *client) MakeRequest(ctx context.Context, endpoint *vo.EndpointConfig, parent *vo.EndpointRequest,
	request *vo.HTTPBackendRequest) (*http.Response, error) {
	httpRequest, err := c.buildNetHTTPRequest(ctx, endpoint, parent, request)
	if checker.NonNil(err) {
		return nil, err
	}

	host := httpRequest.URL.Host
	cb := c.getOrCreateBreaker(host)

	// Circuit breaker check
	if !cb.Allow() {
		return nil, errors.Newf("Circuit breaker open for host=%s", host)
	}

	var resp *http.Response
	maxAttempts := 1
	if checker.Equals(httpRequest.Method, http.MethodGet) || checker.Equals(httpRequest.Method, http.MethodHead) {
		maxAttempts += c.cfg.RetryMaxRetries()
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if checker.IsGreaterThan(attempt, 0) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.cfg.RetryBackoff()):
			}
		}

		resp, err = c.engine.Do(httpRequest)

		if checker.IsLessThan(attempt, maxAttempts-1) && c.isTransientError(err, resp) {
			// Drain body before retry
			if checker.NonNil(resp) {
				_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 512*1024))
				resp.Body.Close()
				resp = nil
			}
			continue
		}
		break
	}

	// Record in circuit breaker
	if c.isTransientError(err, resp) {
		cb.RecordFailure()
	} else {
		cb.RecordSuccess()
	}

	if checker.NonNil(err) {
		return nil, err
	}

	return resp, nil
}

func (c *client) getOrCreateBreaker(host string) *circuitBreaker {
	if v, ok := c.breakers.Load(host); ok {
		return v.(*circuitBreaker)
	}
	cb := newCircuitBreaker(
		c.cfg.CBFailureThreshold(),
		c.cfg.CBSuccessThreshold(),
		c.cfg.CBOpenTimeout(),
		c.cfg.CBHalfOpenMaxReqs(),
	)
	actual, _ := c.breakers.LoadOrStore(host, cb)
	return actual.(*circuitBreaker)
}

// isTransientError returns true for infrastructure failures that should trip the CB
// and be retried (connection errors, 502, 503, 504).
func (c *client) isTransientError(err error, resp *http.Response) bool {
	if checker.NonNil(err) {
		return true
	}
	if checker.NonNil(resp) {
		code := resp.StatusCode
		return checker.Equals(code, http.StatusBadGateway) ||
			checker.Equals(code, http.StatusServiceUnavailable) ||
			checker.Equals(code, http.StatusGatewayTimeout)
	}
	return false
}

func (c *client) buildNetHTTPRequest(ctx context.Context, endpoint *vo.EndpointConfig, parent *vo.EndpointRequest,
	request *vo.HTTPBackendRequest) (*http.Request, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, request.Method(), request.URL(), c.buildNetHTTPRequestBody(request))
	if checker.NonNil(err) {
		return nil, err
	}

	var clientCfg *vo.RequestClientConfig
	if checker.NonNil(endpoint) {
		clientCfg = endpoint.RequestClient()
	}
	httpRequest.Header = c.buildNetHTTPRequestHeader(ctx, clientCfg, parent, request)
	httpRequest.URL.RawQuery = request.Query().Encode()

	return httpRequest, nil
}

func (c *client) buildNetHTTPRequestHeader(ctx context.Context, clientCfg *vo.RequestClientConfig, parent *vo.EndpointRequest,
	request *vo.HTTPBackendRequest) http.Header {
	httpHeader := http.Header(request.Header().Copy())

	if _, exists := httpHeader["User-Agent"]; !exists {
		httpHeader["User-Agent"] = []string{""}
	}

	if _, exists := httpHeader["Accept-Encoding"]; !exists {
		httpHeader.Set("Accept-Encoding", "gzip, deflate")
	}

	if checker.NonNil(clientCfg) && clientCfg.IP().HasPropagateRequest() {
		httpHeader.Set(clientCfg.IP().Propagate().Request(), parent.ClientIP())
	}

	th := clientCfg.TransportHeadersRequest()

	if checker.NonNil(clientCfg) && checker.NonNil(clientCfg.RequestID()) && clientCfg.RequestID().HasPropagateRequest() {
		httpHeader.Set(clientCfg.RequestID().Propagate().Request(), parent.ID())
	}

	if th.DegradationEnabled() {
		httpHeader.Set(app.XGopenDegraded, converter.ToString(request.Degraded()))
		httpHeader.Set(app.XGopenURLPathDegraded, converter.ToString(request.URLPathDegraded()))
		httpHeader.Set(app.XGopenHeaderDegraded, converter.ToString(request.HeaderDegraded()))
		httpHeader.Set(app.XGopenQueryDegraded, converter.ToString(request.QueryDegraded()))
		httpHeader.Set(app.XGopenBodyDegraded, converter.ToString(request.BodyDegraded()))
	}

	if request.HasBody() {
		httpHeader.Set(app.ContentType, request.Body().ContentType().String())
		httpHeader.Set(app.ContentLength, request.Body().SizeInString())

		if request.Body().HasContentEncoding() {
			httpHeader.Set(app.ContentEncoding, request.Body().ContentEncoding().String())
		}
	}

	if th.TimeoutEnabled() {
		if timeout, ok := ctx.Deadline(); ok {
			httpHeader.Set(app.XGopenTimeout, converter.ToString(time.Until(timeout).Milliseconds()))
		}
	}

	return httpHeader
}

func (c *client) buildNetHTTPRequestBody(request *vo.HTTPBackendRequest) io.Reader {
	var body io.ReadCloser
	if request.HasBody() {
		body = io.NopCloser(request.Body().Buffer())
	}
	return body
}
