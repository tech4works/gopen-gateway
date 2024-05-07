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
	configEnum "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/enum"
	configVO "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/vo"
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

func (s endpointService) Execute(ctx context.Context, executeData *vo.ExecuteEndpoint) (*vo.HttpRequest, *vo.HttpResponse) {
	// instanciamos o objeto gopen
	gopen := executeData.Gopen()
	// instanciamos o objeto de valor do endpoint
	endpoint := executeData.Endpoint()
	// instanciamos o objeto de valor da requisição
	httpRequest := executeData.HttpRequest()
	// inicializamos o objeto de valor de resposta do endpoint
	httpResponse := vo.NewHttpResponse()

	// iteramos o beforeware, chaves configuradas para middlewares antes das requisições principais
	httpRequest, httpResponse = s.processMiddlewares(ctx, configEnum.Beforewares, endpoint.Beforewares(), gopen,
		endpoint, httpRequest, httpResponse)
	// verificamos a resposta já foi escrita ou abortada
	if httpResponse.Written() || httpResponse.Abort() {
		return httpRequest, httpResponse
	}

	// iteramos os backends principais para executa-las
	httpRequest, httpResponse = s.processBackends(ctx, endpoint, httpRequest, httpResponse)
	// verificamos a resposta já foi escrita ou abortada
	if httpResponse.Written() || httpResponse.Abort() {
		return httpRequest, httpResponse
	}

	// iteramos o afterware, chaves configuradas para middlewares depois das requisições principais
	httpRequest, httpResponse = s.processMiddlewares(ctx, configEnum.Afterwares, endpoint.Afterwares(), gopen,
		endpoint, httpRequest, httpResponse)
	// verificamos a resposta já foi escrita ou abortada
	if httpResponse.Written() || httpResponse.Abort() {
		return httpRequest, httpResponse
	}

	// retornamos o objeto de valor de resposta final
	return httpRequest, httpResponse
}

func (s endpointService) processMiddlewares(
	ctx context.Context,
	middlewareType configEnum.MiddlewareType,
	middlewareKeys []string,
	gopen *configVO.Gopen,
	endpoint *configVO.Endpoint,
	httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse,
) (*vo.HttpRequest, *vo.HttpResponse) {
	// iteramos as chaves de middlewares
	for _, middlewareKey := range middlewareKeys {
		// verificamos se essa chave foram configuradas no campo middlewares
		middlewareBackend, ok := gopen.Middleware(middlewareKey)
		if !ok {
			logger.Warning(middlewareType, middlewareKey, "not configured on middlewares field!")
			continue
		}
		// instanciamos o objeto de valor de execução do backend
		executeBackend := vo.NewExecuteBackend(endpoint, middlewareBackend, httpRequest, httpResponse)
		// processamos o backend do middleware
		httpRequest, httpResponse = s.backendService.Execute(ctx, executeBackend)
		// verificamos a resposta já foi escrita ou abortada
		if httpResponse.Written() || httpResponse.Abort() {
			break
		}
	}
	// retornamos os novos objetos de valor response e request
	return httpRequest, httpResponse
}

func (s endpointService) processBackends(
	ctx context.Context,
	endpoint *configVO.Endpoint,
	httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse,
) (*vo.HttpRequest, *vo.HttpResponse) {
	// TODO: pensarmos em futuramente ter backends para chamadas concorrentes
	// iteramos os backends fornecidos
	for _, backend := range endpoint.Backends() {
		// instanciamos o vo de execução
		executeBackend := vo.NewExecuteBackend(endpoint, &backend, httpRequest, httpResponse)
		// processamos o backend principal iterado
		httpRequest, httpResponse = s.backendService.Execute(ctx, executeBackend)
		// verificamos a resposta já foi escrita ou abortada
		if httpResponse.Written() || httpResponse.Abort() {
			break
		}
	}
	return httpRequest, httpResponse
}
