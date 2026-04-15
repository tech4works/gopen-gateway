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
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type client struct {
	engine *http.Client
}

func NewClient() app.HTTPClient {
	return client{
		engine: &http.Client{
			Transport: otelhttp.NewTransport(&http.Transport{}),
		},
	}
}

func (c client) MakeRequest(ctx context.Context, endpoint *vo.EndpointConfig, parent *vo.EndpointRequest, request *vo.HTTPBackendRequest) (
	*http.Response, error) {
	httpRequest, err := c.buildNetHTTPRequest(ctx, endpoint, parent, request)
	if checker.NonNil(err) {
		return nil, err
	}
	return c.engine.Do(httpRequest)
}

func (c client) buildNetHTTPRequest(ctx context.Context, endpoint *vo.EndpointConfig, parent *vo.EndpointRequest,
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

func (c client) buildNetHTTPRequestHeader(ctx context.Context, clientCfg *vo.RequestClientConfig, parent *vo.EndpointRequest,
	request *vo.HTTPBackendRequest) http.Header {
	httpHeader := http.Header(request.Header().Copy())

	// Prevent net/http from injecting "User-Agent: Go-http-client/x.x" automatically.
	// The gateway should be transparent — if User-Agent was deleted by a modifier it must not reappear.
	if _, exists := httpHeader["User-Agent"]; !exists {
		httpHeader["User-Agent"] = []string{""}
	}

	// Prevent net/http from injecting "Accept-Encoding: gzip" automatically.
	// The gateway should be transparent — if Accept-Encoding was deleted by a modifier it must not reappear.
	// If not explicitly configured, send only the encodings we support: gzip and deflate
	if _, exists := httpHeader["Accept-Encoding"]; !exists {
		httpHeader.Set("Accept-Encoding", "gzip, deflate")
	}

	// IP propagation (replaces unconditional X-Forwarded-For injection)
	if checker.NonNil(clientCfg) && clientCfg.IP().HasPropagateRequest() {
		httpHeader.Set(clientCfg.IP().Propagate().Request(), parent.ClientIP())
	}

	th := clientCfg.TransportHeadersRequest()

	// request-id: only inject if propagate.request is explicitly configured
	if checker.NonNil(clientCfg) && checker.NonNil(clientCfg.RequestID()) && clientCfg.RequestID().HasPropagateRequest() {
		httpHeader.Set(clientCfg.RequestID().Propagate().Request(), parent.ID())
	}

	// degradation group
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

	// timeout group
	if th.TimeoutEnabled() {
		if timeout, ok := ctx.Deadline(); ok {
			httpHeader.Set(app.XGopenTimeout, converter.ToString(time.Until(timeout).Milliseconds()))
		}
	}

	return httpHeader
}

func (c client) buildNetHTTPRequestBody(request *vo.HTTPBackendRequest) io.Reader {
	var body io.ReadCloser
	if request.HasBody() {
		body = io.NopCloser(request.Body().Buffer())
	}
	return body
}
