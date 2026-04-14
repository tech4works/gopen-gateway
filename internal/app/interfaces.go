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
	"net/http"
	"time"

	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/app/model/publisher"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
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
	Handle(gopen *vo.GopenConfig, endpoint *vo.EndpointConfig, handles ...HandlerFunc)
}

type Context interface {
	Context() context.Context
	WithContext(ctx context.Context)
	Done() <-chan struct{}
	Next()
	Abort()
	IsAborted() bool
	Duration() time.Duration
	Gopen() *vo.GopenConfig
	Endpoint() *vo.EndpointConfig
	Request() *vo.EndpointRequest
	Response() *vo.EndpointResponse
	Write(response *vo.EndpointResponse)
	WriteError(status enum.ResponseStatus, err error)
	WriteString(status enum.ResponseStatus, s string)
	WriteJSON(status enum.ResponseStatus, a any)
	WriteStatus(status enum.ResponseStatus)
	WriteMetadata(metadata vo.Metadata)
}

type HTTPClient interface {
	MakeRequest(ctx context.Context, endpoint *vo.EndpointConfig, parent *vo.EndpointRequest, request *vo.HTTPBackendRequest) (*http.Response, error)
}

type PublisherClient interface {
	Publish(ctx context.Context, parent *vo.EndpointRequest, request *vo.PublisherBackendRequest) (*publisher.Response,
		error)
}

type HTTPLog interface {
	PrintRequest(ctx Context)
	PrintResponse(ctx Context)
}

type MiddlewareLog interface {
	PrintInfof(ctx Context, format string, msg ...any)
	PrintInfo(ctx Context, msg ...any)
	PrintWarnf(ctx Context, format string, msg ...any)
	PrintWarn(ctx Context, msg ...any)
	PrintErrorf(ctx Context, format string, msg ...any)
	PrintError(ctx Context, msg ...any)
}

type EndpointLog interface {
	PrintInfof(executeData dto.ExecuteEndpoint, format string, msg ...any)
	PrintInfo(executeData dto.ExecuteEndpoint, msg ...any)
	PrintWarnf(executeData dto.ExecuteEndpoint, format string, msg ...any)
	PrintWarn(executeData dto.ExecuteEndpoint, msg ...any)
	PrintErrorf(executeData dto.ExecuteEndpoint, format string, msg ...any)
	PrintError(executeData dto.ExecuteEndpoint, msg ...any)
}

type BackendLog interface {
	PrintHTTPRequest(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, request *vo.HTTPBackendRequest)
	PrintPublisherRequest(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, request *vo.PublisherBackendRequest)
	PrintResponse(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, response *vo.BackendResponse)
	PrintInfof(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, format string, msg ...any)
	PrintInfo(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, msg ...any)
	PrintWarnf(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, format string, msg ...any)
	PrintWarn(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, msg ...any)
	PrintErrorf(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, format string, msg ...any)
	PrintError(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, msg ...any)
}
