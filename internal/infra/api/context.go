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

package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type Context struct {
	startTime time.Time
	mutex     *sync.RWMutex
	engine    *engine
	gopen     *vo.GopenConfig
	endpoint  *vo.EndpointConfig
	request   *vo.EndpointRequest
	response  *vo.EndpointResponse
}

type engine struct {
	http *gin.Context
}

func newHTTPContext(gin *gin.Context, gopen *vo.GopenConfig, endpoint *vo.EndpointConfig) app.Context {
	request := buildHTTPRequest(gin, endpoint)
	return &Context{
		startTime: time.Now(),
		mutex:     &sync.RWMutex{},
		engine:    &engine{http: gin},
		gopen:     gopen,
		endpoint:  endpoint,
		request:   request,
	}
}

func resolveClientIP(gin *gin.Context, cfg *vo.RequestClientConfig) string {
	if checker.NonNil(cfg) && cfg.IP().HasHeaders() {
		for _, header := range cfg.IP().Headers() {
			if val := gin.GetHeader(header); checker.IsNotEmpty(val) {
				if isIPTrusted(gin.ClientIP(), cfg.IP().TrustedProxies()) {
					return val
				}
			}
		}
	}
	return gin.ClientIP()
}

func isIPTrusted(clientIP string, trustedProxies []string) bool {
	if checker.IsEmpty(trustedProxies) {
		return true
	}
	for _, proxy := range trustedProxies {
		if clientIP == proxy {
			return true
		}
	}
	return false
}

func resolveRequestID(gin *gin.Context, cfg *vo.RequestClientConfig) string {
	if checker.NonNil(cfg) && checker.NonNil(cfg.RequestID()) && cfg.RequestID().HasHeaders() {
		return cfg.RequestID().ResolveRequestID(gin.Request.Header)
	}
	return uuid.New().String() // current behavior when request-id not configured
}

func buildHTTPRequest(gin *gin.Context, endpoint *vo.EndpointConfig) *vo.EndpointRequest {
	requestID := resolveRequestID(gin, endpoint.RequestClient())

	var traceID string
	span := trace.SpanFromContext(gin.Request.Context())
	if span.SpanContext().IsValid() {
		traceID = span.SpanContext().TraceID().String()
	}

	clientIP := resolveClientIP(gin, endpoint.RequestClient())

	header := vo.NewMetadata(gin.Request.Header)

	query := vo.NewQuery(gin.Request.URL.Query())
	url := gin.Request.URL.Path
	if !query.IsEmpty() {
		url = fmt.Sprint(url, "?", query.Encode())
	}

	ginParams := map[string]string{}
	for _, param := range gin.Params {
		ginParams[param.Key] = param.Value
	}
	path := vo.NewURLPath(gin.FullPath(), ginParams)

	bodyBytes, err := io.ReadAll(gin.Request.Body)
	if checker.NonNil(err) {
		panic(err)
	}

	body := vo.NewPayload(gin.GetHeader(app.ContentType), gin.GetHeader(app.ContentEncoding), bytes.NewBuffer(bodyBytes))

	return vo.NewHTTPEndpointRequest(requestID, traceID, clientIP, url, path, query, gin.Request.Method, header, body)
}

func (c *Context) Context() context.Context {
	return c.engine.http.Request.Context()
}

func (c *Context) Done() <-chan struct{} {
	return c.Context().Done()
}

func (c *Context) WithContext(ctx context.Context) {
	c.engine.http.Request = c.engine.http.Request.WithContext(ctx)
}

func (c *Context) Next() {
	c.engine.http.Next()
}

func (c *Context) Abort() {
	c.engine.http.Abort()
}

func (c *Context) IsAborted() bool {
	return c.engine.http.IsAborted()
}

func (c *Context) Duration() time.Duration {
	return time.Now().Sub(c.startTime)
}

func (c *Context) Gopen() *vo.GopenConfig {
	return c.gopen
}

func (c *Context) Endpoint() *vo.EndpointConfig {
	return c.endpoint
}

func (c *Context) Request() *vo.EndpointRequest {
	return c.request
}

func (c *Context) Response() *vo.EndpointResponse {
	return c.response
}

func (c *Context) Write(response *vo.EndpointResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.writeHTTPResponse(response)
}

func (c *Context) WriteError(status enum.ResponseStatus, err error) {
	wrapped := errors.Wrap(err)
	payload := vo.NewPayloadJSON(converter.ToBuffer(dto.ErrorPayload{
		File:      wrapped.File(),
		Line:      wrapped.Line(),
		Endpoint:  c.endpoint.Path(),
		Message:   wrapped.Message(),
		Stack:     wrapped.Stack(),
		Timestamp: time.Now(),
	}))

	c.Write(vo.NewEndpointResponseWithOnlyStatusAndPayload(vo.NewResponseStatusByValue(status), payload))
}

func (c *Context) WriteStatus(status enum.ResponseStatus) {
	c.Write(vo.NewEndpointResponseWithOnlyStatus(vo.NewResponseStatusByValue(status)))
}

func (c *Context) WriteString(status enum.ResponseStatus, s string) {
	payload := vo.NewPayloadWithContentType(vo.NewContentTypeTextPlain(), converter.ToBuffer(s))

	c.Write(vo.NewEndpointResponseWithOnlyStatusAndPayload(vo.NewResponseStatusByValue(status), payload))
}

func (c *Context) WriteJSON(status enum.ResponseStatus, a any) {
	payload := vo.NewPayloadWithContentType(vo.NewContentTypeJSON(), converter.ToBuffer(a))

	c.Write(vo.NewEndpointResponseWithOnlyStatusAndPayload(vo.NewResponseStatusByValue(status), payload))
}

func (c *Context) WriteMetadata(metadata vo.Metadata) {
	for _, key := range metadata.Keys() {
		c.engine.http.Header(key, metadata.Get(key))
	}
}

func (c *Context) writeHTTPResponse(response *vo.EndpointResponse) {
	if c.IsAborted() {
		return
	}

	var statusCode int
	if response.Status().HasRaw() && converter.CouldBeInt(response.Status().Raw()) {
		statusCode = converter.ToInt(response.Status().Raw())
	} else {
		statusCode = c.parseResponseStatusToHTTPStatusCode(response.Status().Value())
	}

	var contentType vo.ContentType
	var rawBodyBytes []byte
	if response.HasPayload() {
		contentType = response.Payload().ContentType()
		rawBodyBytes = response.Payload().RawBytes()
	}

	for _, key := range response.Metadata().Keys() {
		c.engine.http.Header(key, response.Metadata().Get(key))
	}

	c.decorateHTTPTransportHeaders(response, statusCode)

	if checker.IsNotEmpty(rawBodyBytes) {
		c.engine.http.Data(statusCode, contentType.String(), rawBodyBytes)
	} else {
		c.engine.http.Status(statusCode)
	}

	c.Abort()
	c.response = response
}

func (c *Context) parseResponseStatusToHTTPStatusCode(responseStatus enum.ResponseStatus) int {
	switch responseStatus {
	case enum.ResponseStatusOK:
		return http.StatusOK
	case enum.ResponseStatusCancelled:
		return http.StatusRequestTimeout
	case enum.ResponseStatusInvalidArgument:
		return http.StatusBadRequest
	case enum.ResponseStatusDeadlineExceeded:
		return http.StatusGatewayTimeout
	case enum.ResponseStatusNotFound:
		return http.StatusNotFound
	case enum.ResponseStatusAlreadyExists:
		return http.StatusConflict
	case enum.ResponseStatusPermissionDenied:
		return http.StatusForbidden
	case enum.ResponseStatusUnauthenticated:
		return http.StatusUnauthorized
	case enum.ResponseStatusResourceExhausted:
		return http.StatusTooManyRequests
	case enum.ResponseStatusFailedPrecondition:
		return http.StatusPreconditionFailed
	case enum.ResponseStatusAborted:
		return http.StatusConflict
	case enum.ResponseStatusOutOfRange:
		return http.StatusBadRequest
	case enum.ResponseStatusUnimplemented:
		return http.StatusNotImplemented
	case enum.ResponseStatusInternalError:
		return http.StatusInternalServerError
	case enum.ResponseStatusUnavailable:
		return http.StatusServiceUnavailable
	case enum.ResponseStatusDataLoss:
		return http.StatusInternalServerError
	case enum.ResponseStatusConflict:
		return http.StatusConflict
	case enum.ResponseStatusBadGateway:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

func (c *Context) parseResponseStatusToGRPCCode(responseStatus enum.ResponseStatus) codes.Code {
	switch responseStatus {
	case enum.ResponseStatusOK:
		return codes.OK
	case enum.ResponseStatusCancelled:
		return codes.Canceled
	case enum.ResponseStatusInvalidArgument:
		return codes.InvalidArgument
	case enum.ResponseStatusDeadlineExceeded:
		return codes.DeadlineExceeded
	case enum.ResponseStatusNotFound:
		return codes.NotFound
	case enum.ResponseStatusAlreadyExists:
		return codes.AlreadyExists
	case enum.ResponseStatusPermissionDenied:
		return codes.PermissionDenied
	case enum.ResponseStatusUnauthenticated:
		return codes.Unauthenticated
	case enum.ResponseStatusResourceExhausted:
		return codes.ResourceExhausted
	case enum.ResponseStatusFailedPrecondition:
		return codes.FailedPrecondition
	case enum.ResponseStatusAborted:
		return codes.Aborted
	case enum.ResponseStatusOutOfRange:
		return codes.OutOfRange
	case enum.ResponseStatusUnimplemented:
		return codes.Unimplemented
	case enum.ResponseStatusInternalError:
		return codes.Internal
	case enum.ResponseStatusUnavailable:
		return codes.Unavailable
	case enum.ResponseStatusDataLoss:
		return codes.DataLoss
	case enum.ResponseStatusConflict:
		return codes.Aborted
	case enum.ResponseStatusBadGateway:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}

func (c *Context) decorateHTTPTransportHeaders(response *vo.EndpointResponse, statusCode int) {
	th := c.endpoint.RequestClient().TransportHeadersResponse()

	// request-id: only inject if propagate.response is explicitly configured
	requestClient := c.endpoint.RequestClient()
	if checker.NonNil(requestClient) && checker.NonNil(requestClient.RequestID()) && requestClient.RequestID().HasPropagateResponse() {
		c.engine.http.Header(requestClient.RequestID().Propagate().Response(), c.request.ID())
	}

	// IP propagation via propagate.response (independent of transport-headers groups)
	if c.endpoint.RequestClient().IP().HasPropagateResponse() {
		c.engine.http.Header(c.endpoint.RequestClient().IP().Propagate().Response(), c.request.ClientIP())
	}

	// cache group
	if th.CacheEnabled() {
		// X-Gopen-Cache: true if endpoint cache OR all backends from cache
		cacheHit := response.ComesFromCache()
		if !cacheHit && th.BackendCacheEnabled() {
			cacheHit = response.AllBackendsFromCache()
		}
		c.engine.http.Header(app.XGopenCache, converter.ToString(cacheHit))

		// X-Gopen-Cache-TTL: endpoint TTL or newest backend TTL
		if response.ComesFromCache() {
			c.engine.http.Header(app.XGopenCacheTTL, converter.ToString(response.Cache().RemainingTTL().Time().Milliseconds()))
		} else if th.BackendCacheEnabled() {
			if newestTTL := response.NewestBackendCacheTTLMillis(); newestTTL >= 0 {
				c.engine.http.Header(app.XGopenCacheTTL, converter.ToString(newestTTL))
			}
		}

		// X-Gopen-Backend-Cache: list of cached backend IDs
		if th.BackendCacheEnabled() {
			cachedIDs := response.BackendsCachedIDs()
			if checker.IsNotEmpty(cachedIDs) {
				c.engine.http.Header(app.XGopenBackendCache, strings.Join(cachedIDs, ", "))
			}
		}
	}

	// execution-status group
	if th.ExecutionStatusEnabled() {
		c.engine.http.Header(app.XGopenSuccess, converter.ToString(response.Execution().AllOK()))
		c.engine.http.Header(app.XGopenComplete, converter.ToString(response.Execution().AllExecuted()))
	}

	// degradation group
	if th.DegradationEnabled() {
		c.engine.http.Header(app.XGopenDegraded, converter.ToString(response.Degradation().Any()))

		degradation := response.Degradation()
		c.engine.http.Header(app.XGopenHeaderDegraded, converter.ToString(degradation.Has(enum.DegradationKindMetadata)))
		c.engine.http.Header(app.XGopenMetadataDegraded, converter.ToString(degradation.Has(enum.DegradationKindMetadata)))
		c.engine.http.Header(app.XGopenQueryDegraded, converter.ToString(degradation.Has(enum.DegradationKindQuery)))
		c.engine.http.Header(app.XGopenURLPathDegraded, converter.ToString(degradation.Has(enum.DegradationKindURLPath)))
		c.engine.http.Header(app.XGopenBodyDegraded, converter.ToString(degradation.Has(enum.DegradationKindPayload)))
		c.engine.http.Header(app.XGopenPayloadDegraded, converter.ToString(degradation.Has(enum.DegradationKindPayload)))
		c.engine.http.Header(app.XGopenDeduplicationIDDegraded, converter.ToString(degradation.Has(enum.DegradationKindDeduplicationID)))
		c.engine.http.Header(app.XGopenGroupIDDegraded, converter.ToString(degradation.Has(enum.DegradationKindGroupID)))
		c.engine.http.Header(app.XGopenAttributeDegraded, converter.ToString(degradation.Has(enum.DegradationKindAttributes)))

		backendsDegraded := response.Execution().Degradations()
		if checker.IsNotEmpty(backendsDegraded) {
			ids := make([]string, len(backendsDegraded))
			for i, backendDegraded := range backendsDegraded {
				ids[i] = backendDegraded.ID()
			}
			c.engine.http.Header(app.XGopenDegradedBackendCount, converter.ToString(len(ids)))
			c.engine.http.Header(app.XGopenDegradedBackends, strings.Join(ids, ", "))
		}
	}

	// Content headers (always injected, not controlled by transport-headers)
	if response.HasPayload() {
		c.engine.http.Header(app.ContentType, response.Payload().ContentType().String())
		c.engine.http.Header(app.ContentLength, response.Payload().SizeInString())
		if response.Payload().HasContentEncoding() {
			c.engine.http.Header(app.ContentEncoding, response.Payload().ContentEncoding().String())
		}
	}
}
