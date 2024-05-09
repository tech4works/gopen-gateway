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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"net/http"
)

type backendService struct {
	restTemplate interfaces.RestTemplate
}

type Backend interface {
	Execute(ctx context.Context, executeData *vo.ExecuteBackend) (*vo.HttpRequest, *vo.HttpResponse)
}

func NewBackend(restTemplate interfaces.RestTemplate) Backend {
	return backendService{
		restTemplate: restTemplate,
	}
}

func (b backendService) Execute(ctx context.Context, executeData *vo.ExecuteBackend) (*vo.HttpRequest, *vo.HttpResponse) {
	endpoint := executeData.Endpoint()
	backend := executeData.Backend()
	httpRequest, httpResponse := b.buildRequest(executeData)

	netHttpResponse, err := b.executeRequest(ctx, backend, httpRequest, httpResponse)
	if helper.IsNotNil(err) {
		return httpRequest, httpResponse.Error(executeData.Endpoint().Path(), err)
	}
	defer b.closeNetHttpResponse(netHttpResponse)

	httpResponse = b.buildResponse(endpoint, backend, netHttpResponse, httpRequest, httpResponse)

	return httpRequest, httpResponse
}

func (b backendService) buildRequest(executeData *vo.ExecuteBackend) (*vo.HttpRequest, *vo.HttpResponse) {
	httpRequest := executeData.HttpRequest().Modify(executeData.Backend().Request(), executeData.HttpResponse())
	httpBackendRequest := vo.NewHttpBackendRequest(executeData.Backend(), httpRequest, executeData.HttpResponse())
	httpRequest = httpRequest.Append(httpBackendRequest)

	return httpRequest, executeData.HttpResponse()
}

func (b backendService) executeRequest(
	ctx context.Context,
	backend *vo.Backend,
	httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse,
) (*http.Response, error) {
	httpBackendRequest := vo.NewHttpBackendRequest(backend, httpRequest, httpResponse)
	return b.makeNetHttpRequest(ctx, httpBackendRequest)
}

func (b backendService) buildResponse(
	endpoint *vo.Endpoint,
	backend *vo.Backend,
	netHttpResponse *http.Response,
	httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse,
) *vo.HttpResponse {
	httpBackendResponse := vo.NewHttpBackendResponse(backend, netHttpResponse, httpRequest, httpResponse)
	if helper.IsNotNil(httpBackendResponse) && endpoint.Abort(httpBackendResponse.StatusCode()) {
		return vo.NewHttpResponseAborted(endpoint, httpBackendResponse)
	}
	return httpResponse.Append(httpBackendResponse)
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

func (b backendService) closeNetHttpResponse(netHttpResponse *http.Response) {
	err := netHttpResponse.Body.Close()
	if helper.IsNotNil(err) {
		logger.WarningSkipCaller(2, "Error close http response:", err)
	}
}
