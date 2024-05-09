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

package infra

import (
	berrors "errors"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"net/http"
	"net/url"
	"time"
)

// restTemplate is a struct that represents a template for making HTTP requests.
// It implements the interfaces.RestTemplate interface, which provides a method MakeRequest for sending
// an HTTP request and returning the corresponding HTTP response and an error, if any.
type restTemplate struct {
}

// NewRestTemplate returns a new instance of a restTemplate object.
// It implements the interfaces.RestTemplate interface.
func NewRestTemplate() interfaces.RestTemplate {
	return restTemplate{}
}

// MakeRequest sends an HTTP request to a backend and returns the HTTP response.
// It logs the backend HTTP request, handles HTTP client errors, and logs the backend
// HTTP response. If an error occurs during the request, it returns nil for the response
// and the error.
//
// Inputs:
//   - netHttpRequest: The HTTP request to be sent to the backend.
//
// Returns:
//   - *http.Response: The HTTP response obtained from the backend, if successful.
//   - error: The error occurred during the request, if any.
//
// Note: This method internally calls printNetHttpRequest to log the backend HTTP request.
// If an error occurs during the request, it calls treatHttpClientErr to handle the error,
// calls printNetHttpResponseError to log the error, and returns the error.
// Otherwise, it calls printNetHttpResponse to log the backend HTTP response and returns
// the response.
func (r restTemplate) MakeRequest(netHttpRequest *http.Request) (*http.Response, error) {
	startTime := time.Now()

	r.printNetHttpRequest(netHttpRequest)

	httpClient := http.Client{}
	netHttpResponse, err := httpClient.Do(netHttpRequest)

	latency := time.Since(startTime).String()

	err = r.treatHttpClientErr(err)
	if helper.IsNotNil(err) {
		r.printNetHttpResponseError(netHttpRequest, latency, err)
		return nil, err
	}

	r.printNetHttpResponse(netHttpRequest, netHttpResponse, latency)
	return netHttpResponse, nil
}

// printNetHttpRequest prints the information about an HTTP request being made to a backend.
// It includes the HTTP method and URL.
//
// Note: This method is used internally in MakeRequest to log the backend HTTP request.
func (r restTemplate) printNetHttpRequest(netHttpRequest *http.Request) {
	httpUrl := netHttpRequest.URL.String()
	httpMethod := netHttpRequest.Method

	msg := fmt.Sprintf("Backend HTTP request: %s --> %s", httpMethod, httpUrl)
	logger.Debug(msg)
}

// printNetHttpResponseError prints the information about an error that occurred during an HTTP response from the backend.
// It includes the HTTP method, URL, latency, and the error message.
//
// Note: This method is used internally in MakeRequest to log the backend HTTP response error.
func (r restTemplate) printNetHttpResponseError(netHttpRequest *http.Request, latency string, err error) {
	httpUrl := netHttpRequest.URL.String()
	httpMethod := netHttpRequest.Method
	logger.Errorf("Backend HTTP response: %s --> %s latency: %s err: %s", httpMethod, httpUrl, latency, err)
}

// printNetHttpResponse prints the information about the HTTP response obtained from the backend.
// It includes the HTTP method, URL, latency, and status code.
//
// Note: This method is used internally in MakeRequest to log the backend HTTP response.
func (r restTemplate) printNetHttpResponse(netHttpRequest *http.Request, netHttpResponse *http.Response, latency string) {
	httpUrl := netHttpRequest.URL.String()
	httpMethod := netHttpRequest.Method

	msg := fmt.Sprintf("Backend HTTP response: %s --> %s latency: %s statusCode: %v", httpMethod, httpUrl, latency,
		netHttpResponse.StatusCode)
	logger.Debug(msg)
}

// treatHttpClientErr handles and transforms HTTP client errors.
// If the error is nil, it returns nil.
// If the error is an *url.Error and has a timeout, it returns a new domainmapper.ErrGatewayTimeout error.
// For any other error, it returns a new domainmapper.ErrBadGateway error.
//
// Inputs:
//   - err: The HTTP client error to be treated.
//
// Returns:
//   - error: The transformed error, if any.
//
// Example:
//
//	err := r.treatHttpClientErr(http.ErrTimeout)
//
// Note: For more details, see mapper.NewErrGatewayTimeoutByErr and mapper.NewErrBadGateway declarations.
func (r restTemplate) treatHttpClientErr(err error) error {
	if helper.IsNil(err) {
		return nil
	}

	var urlErr *url.Error
	berrors.As(err, &urlErr)
	if urlErr.Timeout() {
		err = mapper.NewErrGatewayTimeoutByErr(err)
	} else {
		err = mapper.NewErrBadGateway(err)
	}
	return err
}
