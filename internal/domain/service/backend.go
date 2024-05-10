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
// It takes a parameter restTemplate of type interfaces.RestTemplate, which represents a template for making HTTP requests.
// The function returns a backendService object that implements the Backend interface.
func NewBackend(restTemplate interfaces.RestTemplate) Backend {
	return backendService{
		restTemplate: restTemplate,
	}
}

// Execute is a method that executes a backend request and handles the corresponding response.
// It takes two parameters: ctx of type context.Context, which represents the context of the request,
// and executeData of type *vo.ExecuteBackend, which represents the execution context for a backend request.
// The method retrieves the endpoint and backend from the executeData object.
// It then calls the buildRequest() method to construct the httpRequest and httpResponse objects.
// The method calls the executeRequest() method to send the httpRequest to the backend service.
// If an error occurs during the execution of the request, the method returns an HTTP error response
// using the httpRequest and httpResponse objects.
// The method calls the closeNetHttpResponse() method to close the netHttpResponse object.
// Finally, the method calls the buildResponse() method to construct the final httpResponse object
// based on the netHttpResponse object and the httpRequest and httpResponse objects.
// The method returns the httpRequest and httpResponse objects.
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

// buildRequest is a method that constructs a *vo.HttpRequest and *vo.HttpResponse objects
// using the given *vo.ExecuteBackend object.
//
// It takes one parameter: executeData of type *vo.ExecuteBackend, which represents the execution context
// for a backend request. This object contains an endpoint, backend, httpRequest, and httpResponse.
//
// The method creates an *vo.HttpRequest object by modifying the backend's request with the incoming request
// and response objects. This is done by calling the Modify() method of the executeData.HttpRequest object.
//
// It then creates an *vo.HttpBackendRequest object using the given executeData.Backend(), httpRequest, and
// executeData.HttpResponse() objects. The creation is performed by calling the NewHttpBackendRequest() function
// of the vo package.
//
// Next, it appends the httpBackendRequest to the httpRequest using the Append() method of the httpRequest object.
//
// Finally, the method returns the httpRequest and executeData.HttpResponse objects.
func (b backendService) buildRequest(executeData *vo.ExecuteBackend) (*vo.HttpRequest, *vo.HttpResponse) {
	httpRequest := executeData.HttpRequest().Modify(executeData.Backend().Request(), executeData.HttpResponse())

	httpBackendRequest := vo.NewHttpBackendRequest(executeData.Backend(), httpRequest, executeData.HttpResponse())
	httpRequest = httpRequest.Append(httpBackendRequest)

	return httpRequest, executeData.HttpResponse()
}

// executeRequest is a method that sends an HTTP request to the backend service.
// It takes four parameters: ctx of type context.Context, which represents the context of the request,
// backend of type *vo.Backend, which represents the backend configuration,
// httpRequest of type *vo.HttpRequest, which represents the HTTP request object,
// and httpResponse of type *vo.HttpResponse, which represents the HTTP response object.
// The method creates an *vo.HttpBackendRequest object using the given parameters. This object represents the
// backend HTTP request. The creation is performed by calling the NewHttpBackendRequest() function of the vo package.
// The method then calls the makeNetHttpRequest() method to send the HTTP request and returns the
// *http.Response and error returned by the makeNetHttpRequest() method.
func (b backendService) executeRequest(
	ctx context.Context,
	backend *vo.Backend,
	httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse,
) (*http.Response, error) {
	httpBackendRequest := vo.NewHttpBackendRequest(backend, httpRequest, httpResponse)
	return b.makeNetHttpRequest(ctx, httpBackendRequest)
}

// buildResponse is a method that constructs a *vo.HttpResponse object based on the given parameters.
// It takes five parameters: endpoint of type *vo.Endpoint, backend of type *vo.Backend, netHttpResponse of type
// *http.Response, httpRequest of type *vo.HttpRequest, and httpResponse of type *vo.HttpResponse.
//
// The method creates an *vo.HttpBackendResponse object using the given parameters. This object represents the
// backend HTTP response. The creation is performed by calling the NewHttpBackendResponse() function of the vo package.
//
// If the httpBackendResponse object is not nil and the endpoint's Abort() method returns true for the response's
// status code, the method creates and returns a new *vo.HttpResponseAborted object. This object represents the
// aborted HTTP response. The creation is performed by calling the NewHttpResponseAborted() function of the vo package.
//
// If the httpBackendResponse object is nil or the endpoint's Abort() method returns false for the response's
// status code, the method appends the httpBackendResponse object to the given httpResponse object and returns it.
// The appending is performed by calling the Append() method of the httpResponse object.
//
// The method returns the *vo.HttpResponse object.
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

// makeNetHttpRequest is a method that sends an HTTP request to the backend service.
// It takes two parameters: ctx of type context.Context, which represents the context of the request,
// and httpBackendRequest of type *vo.HttpBackendRequest, which represents the backend HTTP request object.
// The method calls the NetHttp() method of the httpBackendRequest object to create an *http.Request instance.
// If an error occurs during the creation of the *http.Request instance, the method returns nil and the error.
// Otherwise, it calls the MakeRequest() method of the restTemplate object to send the HTTP request.
// The method returns the *http.Response and error returned by the MakeRequest() method.
func (b backendService) makeNetHttpRequest(ctx context.Context, httpBackendRequest *vo.HttpBackendRequest) (
	*http.Response, error) {
	netHttpRequest, err := httpBackendRequest.NetHttp(ctx)
	if helper.IsNotNil(err) {
		return nil, err
	}
	return b.restTemplate.MakeRequest(netHttpRequest)
}

// closeNetHttpResponse is a method that closes the HTTP response body.
// It takes a parameter netHttpResponse of type *http.Response, which represents the HTTP response object.
// The method calls the Body.Close() function to close the response body.
// If an error occurs while closing the response body, it logs a warning message using the logger.WarningSkipCaller
// function.
func (b backendService) closeNetHttpResponse(netHttpResponse *http.Response) {
	err := netHttpResponse.Body.Close()
	if helper.IsNotNil(err) {
		logger.WarningSkipCaller(2, "Error close http response:", err)
	}
}
