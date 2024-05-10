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

// endpointService is a struct type that represents an implementation of the Endpoint interface.
//
// It contains a backendService field of type Backend, which represents the backend service used by the endpoint.
// The backendService field is responsible for executing backend operations.
//
// Note: The endpointService type should be used in conjunction with the NewEndpoint function to create new instances.
// It should also implement the Execute method of the Endpoint interface.
type endpointService struct {
	// backendService represents a backend service in the application.
	backendService Backend
}

// Endpoint is an interface that represents an API endpoint in the Gopen application.
// It contains the Execute method, which accepts a context, ExecuteEndpoint data, and returns an HttpRequest and HttpResponse.
// The Execute method is responsible for executing the endpoint logic and generating the appropriate HTTP request and response.
type Endpoint interface {
	// Execute is a method that executes the logic of an API endpoint in the Gopen application.
	//
	// Parameters:
	// - ctx: the context of the execution.
	// - executeData: the ExecuteEndpoint data representing the execution context for the endpoint.
	//
	// Returns:
	// - httpRequest: the generated HttpRequest based on the execution result.
	// - httpResponse: the generated HttpResponse based on the execution result.
	Execute(ctx context.Context, executeData *vo.ExecuteEndpoint) (*vo.HttpRequest, *vo.HttpResponse)
}

// NewEndpoint initializes a new instance of the endpointService struct, which implements the Endpoint interface.
//
// Parameters:
// - backendService: an instance of the Backend interface representing the backend service used by the endpoint.
//
// Returns:
// - endpointService: an instance of the endpointService struct that implements the Endpoint interface.
func NewEndpoint(backendService Backend) Endpoint {
	return endpointService{
		backendService: backendService,
	}
}

// Execute executes the API endpoint by processing middlewares, backends, and afterwares.
// It returns the updated vo.HttpRequest and vo.HttpResponse objects.
//
// Parameters:
// - ctx: the current context.Context
// - executeData: a pointer to the vo.ExecuteEndpoint struct representing the execution context
//
// Returns:
// - *vo.HttpRequest: a pointer to the updated vo.HttpRequest struct
// - *vo.HttpResponse: a pointer to the updated vo.HttpResponse struct
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

// checkHttpResponse checks if the httpResponse has been written or if it needs to be aborted.
// It returns true if the httpResponse has been written or needs to be aborted, false otherwise.
//
// Parameters:
// - httpResponse: a pointer to the vo.HttpResponse struct representing the HTTP response
//
// Returns:
// - bool: true if the httpResponse has been written or needs to be aborted, false otherwise.
func (e endpointService) checkHttpResponse(httpResponse *vo.HttpResponse) bool {
	return httpResponse.Written() || httpResponse.Abort()
}

// processMiddlewares executes the middleware operations defined in the given endpoint.
//
// Parameters:
// - ctx: the current context.Context
// - middlewareType: the type of the middleware, which is the enum.MiddlewareType
// - middlewareKeys: a slice of strings representing the keys of the middlewares to be executed
// - gopen: a pointer to the vo.Gopen struct representing the Gopen server configuration
// - endpoint: a pointer to the vo.Endpoint struct representing the API endpoint
// - httpRequest: a pointer to the vo.HttpRequest struct representing the HTTP request
// - httpResponse: a pointer to the vo.HttpResponse struct representing the HTTP response
//
// Returns:
// - httpRequest: a pointer to the updated vo.HttpRequest struct
// - httpResponse: a pointer to the updated vo.HttpResponse struct
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
	return httpRequest, httpResponse
}

// processBackends executes the backend operations defined in the given endpoint.
//
// Parameters:
// - ctx: the current context.Context
// - endpoint: a pointer to the vo.Endpoint struct representing the API endpoint
// - httpRequest: a pointer to the vo.HttpRequest struct representing the HTTP request
// - httpResponse: a pointer to the vo.HttpResponse struct representing the HTTP response
//
// Returns:
// - httpRequest: a pointer to the updated vo.HttpRequest struct
// - httpResponse: a pointer to the updated vo.HttpResponse struct
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
