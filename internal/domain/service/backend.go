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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

// backendService is a type that represents a backend service in the application.
// It is responsible for executing backend requests and handling the corresponding responses.
// The type contains a restTemplate object of type interfaces.RestTemplate, which is used to send HTTP requests.
type backendService struct {
	// restTemplate represents an interface for making HTTP requests.
	restTemplate interfaces.RestTemplate
}

// Backend is an interface representing a backend service in the application.
type Backend interface {
	// Execute is a method of the Backend interface that is responsible for executing backend operations.
	//
	// It takes a context.Context object as the first parameter, which provides the execution context and allows
	// cancellation or timeouts for the operation.
	//
	// The second parameter, executeData, is an instance of the vo.ExecuteBackend struct, which holds the execution
	// context for the backend operation. It contains pointers to the Endpoint, Backend, HttpRequest, and HttpResponse
	// structs, representing the specific API endpoint configuration, backend configuration, incoming HTTP request,
	// and outgoing HTTP response, respectively.
	//
	// The method returns two pointers: a vo.HttpRequest pointer and a vo.HttpResponse pointer. The HttpRequest pointer
	// represents the modified incoming HTTP request, after any necessary modifications or processing by the backend.
	// The HttpResponse pointer represents the outgoing HTTP response, ready to be sent back to the client.
	//
	// Note: The method should be implemented by types that implement the Backend interface.
	Execute(ctx context.Context, executeData *vo.ExecuteBackend) (*vo.HttpRequest, *vo.HttpResponse)
}

// NewBackend is a function that creates and returns a new instance of the Backend interface.
// It takes a parameter restTemplate of type interfaces.RestTemplate, which represents a template for making HTTP
// requests. The function returns a backendService object that implements the Backend interface.
func NewBackend(restTemplate interfaces.RestTemplate) Backend {
	return backendService{
		restTemplate: restTemplate,
	}
}

func (b backendService) Execute(ctx context.Context, executeData *vo.ExecuteBackend) (*vo.HttpRequest, *vo.HttpResponse) {
	endpoint := executeData.Endpoint()
	backend := executeData.Backend()
	httpResponse := executeData.HttpResponse()
	httpRequest := executeData.HttpRequest()

	httpRequest = httpRequest.Modify(executeData.Backend().Request(), httpResponse)
	httpBackendRequest := vo.NewHttpBackendRequest(executeData.Backend(), httpRequest, executeData.HttpResponse())
	httpRequest = httpRequest.Append(httpBackendRequest)

	httpBackendResponse, err := b.restTemplate.MakeRequest(ctx, backend, httpBackendRequest)
	if helper.IsNotNil(err) {
		return httpRequest, vo.NewHttpResponseByErr(executeData.Endpoint().Path(), err)
	}

	httpBackendResponse = httpBackendResponse.ApplyConfig(enum.BackendResponseApplyEarly, httpRequest, httpResponse)
	if helper.IsNotNil(httpBackendResponse) && endpoint.Abort(httpBackendResponse.StatusCode()) {
		return httpRequest, vo.NewHttpResponseAborted(endpoint, httpBackendResponse)
	}
	httpResponse = httpResponse.Append(httpBackendResponse)

	return httpRequest, httpResponse
}
