/*
 * Copyright 2024 Gabriel Cataldo
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

package service

import (
	"context"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type endpointService struct {
	backendService Backend
}

type Endpoint interface {
	Execute(ctx context.Context, executeData *vo.ExecuteEndpoint) (*vo.HttpRequest, *vo.HttpResponse)
}

func NewEndpoint(backendService Backend) Endpoint {
	return endpointService{
		backendService: backendService,
	}
}

func (e endpointService) Execute(ctx context.Context, executeData *vo.ExecuteEndpoint) (*vo.HttpRequest, *vo.HttpResponse) {
	gopen := executeData.Gopen()
	endpoint := executeData.Endpoint()
	httpRequest := executeData.HttpRequest()
	httpResponse := vo.NewHttpResponse()

	httpRequest, httpResponse = e.processMiddlewares(ctx, enum.Beforewares, endpoint.Beforewares(), gopen, endpoint,
		httpRequest, httpResponse)
	if e.checkHttpResponse(httpResponse) {
		return httpRequest, httpResponse
	}

	httpRequest, httpResponse = e.processBackends(ctx, endpoint, httpRequest, httpResponse)
	if e.checkHttpResponse(httpResponse) {
		return httpRequest, httpResponse
	}

	httpRequest, httpResponse = e.processMiddlewares(ctx, enum.Afterwares, endpoint.Afterwares(), gopen, endpoint,
		httpRequest, httpResponse)
	if e.checkHttpResponse(httpResponse) {
		return httpRequest, httpResponse
	}

	return httpRequest, httpResponse
}

func (e endpointService) checkHttpResponse(httpResponse *vo.HttpResponse) bool {
	return httpResponse.Written() || httpResponse.Abort()
}

func (e endpointService) processMiddlewares(
	ctx context.Context,
	middlewareType enum.MiddlewareType,
	middlewareKeys []string,
	gopen *vo.Gopen,
	endpoint *vo.Endpoint,
	httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse,
) (*vo.HttpRequest, *vo.HttpResponse) {
	for _, middlewareKey := range middlewareKeys {
		middlewareBackend, ok := gopen.Middleware(middlewareKey)
		if !ok {
			logger.Warning(middlewareType, middlewareKey, "not configured on middlewares field!")
			continue
		}

		executeBackend := vo.NewExecuteBackend(endpoint, middlewareBackend, httpRequest, httpResponse)
		httpRequest, httpResponse = e.backendService.Execute(ctx, executeBackend)

		if e.checkHttpResponse(httpResponse) {
			break
		}
	}
	// retornamos os novos objetos de valor response e request
	return httpRequest, httpResponse
}

func (e endpointService) processBackends(
	ctx context.Context,
	endpoint *vo.Endpoint,
	httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse,
) (*vo.HttpRequest, *vo.HttpResponse) {
	// TODO: pensarmos em futuramente ter backends para chamadas concorrentes
	for _, backend := range endpoint.Backends() {
		executeBackend := vo.NewExecuteBackend(endpoint, &backend, httpRequest, httpResponse)
		httpRequest, httpResponse = e.backendService.Execute(ctx, executeBackend)
		if e.checkHttpResponse(httpResponse) {
			break
		}
	}
	return httpRequest, httpResponse
}
