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
// It takes a parameter restTemplate of type interfaces.RestTemplate, which represents a template for making HTTP
// requests. The function returns a backendService object that implements the Backend interface.
func NewBackend(restTemplate interfaces.RestTemplate) Backend {
	return backendService{
		restTemplate: restTemplate,
	}
}

// Execute is a method that executes a backend request and handles the corresponding response.
// It takes two parameters: ctx of type context.Context, which represents the context of the request,
// and executeData of type *vo.ExecuteBackend, which represents the execution context for a backend in the application.
// The method retrieves the endpoint, backend, httpRequest, and httpResponse from the executeData object.
// It modifies the httpRequest by calling the Modify() method of the executeData's backend and passing in the request
// and response.
// It creates a new httpBackendRequest object using the backend, modified httpRequest, and httpResponse.
// The httpRequest is then appended with httpBackendRequest using the Append() method.
// The method makes an HTTP request by calling the makeNetHttpRequest() method with the context and httpBackendRequest,
// capturing the returned netHttpResponse and error.
// If an error occurs, it returns the httpRequest and an error response created by the httpResponse.Error() method.
// The netHttpResponse is closed using the closeNetHttpResponse() method.
// It creates a new httpBackendResponse using the backend, netHttpResponse, httpRequest, and httpResponse.
// If httpBackendResponse is not nil and endpoint.Abort() returns true with the httpBackendResponse's status code,
// it returns the httpRequest and a new response of type HttpResponseAborted using the NewHttpResponseAborted() method.
// The httpResponse is then appended with the httpBackendResponse.
// Finally, it returns the httpRequest and httpResponse.
func (b backendService) Execute(ctx context.Context, executeData *vo.ExecuteBackend) (*vo.HttpRequest, *vo.HttpResponse) {
	endpoint := executeData.Endpoint()
	backend := executeData.Backend()
	httpResponse := executeData.HttpResponse()
	httpRequest := executeData.HttpRequest()

	httpRequest = httpRequest.Modify(executeData.Backend().Request(), httpResponse)

	httpBackendRequest := vo.NewHttpBackendRequest(executeData.Backend(), httpRequest, executeData.HttpResponse())
	httpRequest = httpRequest.Append(httpBackendRequest)

	netHttpResponse, err := b.makeNetHttpRequest(ctx, httpBackendRequest)
	if helper.IsNotNil(err) {
		return httpRequest, vo.NewHttpResponseByErr(executeData.Endpoint().Path(), err)
	}
	defer b.closeNetHttpResponse(netHttpResponse)

	httpBackendResponse := vo.NewHttpBackendResponse(backend, netHttpResponse, httpRequest, httpResponse)
	if helper.IsNotNil(httpBackendResponse) && endpoint.Abort(httpBackendResponse.StatusCode()) {
		return httpRequest, vo.NewHttpResponseAborted(endpoint, httpBackendResponse)
	}
	httpResponse = httpResponse.Append(httpBackendResponse)

	return httpRequest, httpResponse
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
