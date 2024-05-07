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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/vo"
	"net/http"
)

type backendService struct {
	restTemplate interfaces.RestTemplate
}

type Backend interface {
	Execute(ctx context.Context, executeData *vo.ExecuteBackend) (*vo.HttpRequest, *vo.HttpResponse)
}

// NewBackend initializes and returns a new Backend instance.
//
// This function serves as a factory, accepting implementations of Modifier and interfaces.RestTemplate
// as arguments, creating and returning an instance of the backend type that satisfies the Backend interface.
//
// Parameters:
// modifierService: Provides the service for modifying backend information. Must conform to the Modifier interface.
// restTemplate: Provides the functionality for conducting RESTful operations. Must conform to the RestTemplate
// interface from interfaces package.
//
// Returns:
// A Backend instance with modifierService and restTemplate composed in.
func NewBackend(restTemplate interfaces.RestTemplate) Backend {
	return backendService{
		restTemplate: restTemplate,
	}
}

func (b backendService) Execute(ctx context.Context, executeData *vo.ExecuteBackend) (*vo.HttpRequest, *vo.HttpResponse) {
	// instanciamos o endpoint vo
	endpoint := executeData.Endpoint()
	// instanciamos o backend vo
	backend := executeData.Backend()
	// instanciamos o objeto de valor de request
	httpRequest := executeData.HttpRequest()
	// instanciamos o objeto de valor de response
	httpResponse := executeData.HttpResponse()

	// caso tenha modifiers de request, com propagate true e escopo REQUEST ele ira modificar, se não, retorna ele mesmo
	httpRequest = httpRequest.Modify(backend.Request(), httpResponse)

	// construímos o backend request
	httpBackendRequest := vo.NewHttpBackendRequest(backend, httpRequest, httpResponse)

	// damos o append no backend request no histórico de requisições
	httpRequest = httpRequest.Append(httpBackendRequest)

	// realizamos com base no http backend request
	netHttpResponse, err := b.makeNetHttpRequest(ctx, httpBackendRequest)
	if helper.IsNotNil(err) {
		return httpRequest, httpResponse.Error(executeData.Endpoint().Path(), err)
	}
	// chamamos para fechar o body assim que possível
	defer b.closeNetHttpResponse(netHttpResponse)

	// construímos o backend response a partir do netHttpResponse e configurações no EARLY
	httpBackendResponse := vo.NewHttpBackendResponse(backend, netHttpResponse, httpRequest, httpResponse)

	// se resposta é para abortar, e retornamos a mesma
	if helper.IsNotNil(httpBackendResponse) && endpoint.Abort(httpBackendResponse.StatusCode()) {
		return httpRequest, vo.NewHttpResponseAborted(endpoint, httpBackendResponse)
	}

	// adicionamos o novo backend request no objeto de valor de resposta como histórico
	httpResponse = httpResponse.Append(httpBackendResponse)

	// se tudo ocorrer bem retornamos a requisição e o response resultante
	return httpRequest, httpResponse
}

func (b backendService) makeNetHttpRequest(ctx context.Context, httpBackendRequest *vo.HttpBackendRequest) (
	*http.Response, error) {
	// montamos o http request com o context
	netHttpRequest, err := httpBackendRequest.NetHttp(ctx)
	// caso ocorra um erro na montagem, retornamos o mesmo
	if helper.IsNotNil(err) {
		return nil, err
	}
	// chamamos a interface de infra para chamar a conexão http e tratar a resposta
	return b.restTemplate.MakeRequest(netHttpRequest)
}

// closeNetHttpResponse closes the HTTP response body.
// If there is an error while closing the body, a warning message will be logged.
func (b backendService) closeNetHttpResponse(netHttpResponse *http.Response) {
	err := netHttpResponse.Body.Close()
	if helper.IsNotNil(err) {
		logger.WarningSkipCaller(2, "Error close http response:", err)
	}
}
