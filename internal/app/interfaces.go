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

package app

import (
	"context"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"net/http"
	"time"
)

type Boot interface {
	Init() *dto.Gopen
	Start(gopen *dto.Gopen)
	Stop()
}

type BootLog interface {
	PrintLogo()
	PrintTitle(title string)
	PrintInfo(msg ...any)
	PrintInfof(format string, msg ...any)
	PrintWarn(msg ...any)
	PrintWarnf(format string, msg ...any)
	PrintError(msg ...any)
	SkipLine()
}

type HandlerFunc func(ctx Context)

type Router interface {
	Engine() http.Handler
	Handle(gopen *vo.Gopen, endpoint *vo.Endpoint, handles ...HandlerFunc)
}

type Context interface {
	Context() context.Context
	WithContext(ctx context.Context)
	Done() <-chan struct{}
	Next()
	Duration() time.Duration
	TraceID() string
	ClientIP() string
	Gopen() *vo.Gopen
	Endpoint() *vo.Endpoint
	Request() *vo.HTTPRequest
	Response() *vo.HTTPResponse
	Write(response *vo.HTTPResponse)
	WriteCacheResponse(cacheResponse *vo.CacheResponse)
	WriteError(code int, err error)
	WriteString(code int, s string)
	WriteJson(code int, a any)
	WriteStatusCode(code int)
}

type HTTPClient interface {
	MakeRequest(ctx context.Context, request *vo.HTTPBackendRequest) (*http.Response, error)
}

type HTTPLog interface {
	PrintRequest(ctx Context)
	PrintResponse(ctx Context)
}

type EndpointLog interface {
	PrintInfof(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID, format string, msg ...any)
	PrintInfo(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID string, msg ...any)
	PrintWarnf(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID, format string, msg ...any)
	PrintWarn(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID string, msg ...any)
	PrintErrorf(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID, format string, msg ...any)
	PrintError(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID string, msg ...any)
}

type BackendLog interface {
	PrintRequest(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest)
	PrintResponse(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest, response *vo.HTTPBackendResponse, duration time.Duration)
	PrintInfof(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest, format string, msg ...any)
	PrintInfo(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest, msg ...any)
	PrintWarnf(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest, format string, msg ...any)
	PrintWarn(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest, msg ...any)
	PrintErrorf(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest, format string, msg ...any)
	PrintError(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest, msg ...any)
}
