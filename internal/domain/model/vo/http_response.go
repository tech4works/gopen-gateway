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

package vo

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"net/http"
)

// HttpResponse represents the gateway HTTP httpResponse.
type HttpResponse struct {
	// statusCode stores the integer HTTP status code of the HttpResponse object.
	statusCode StatusCode
	// header represents the header of the HttpResponse object.
	header Header
	// Body represents the body of the gateway HTTP httpResponse.
	body *Body
	// abort represents the property "abort" of the HttpResponse object. If `abort` is true,
	// it indicates that the httpResponse should be aborted. Returns a boolean value representing the `abort` property.
	abort bool
	// written indicates whether the HttpResponse has been written.
	written bool
	// history represents the history of backend responses in the HttpResponse object.
	history httpResponseHistory
}

// NewHttpResponse creates a new HttpResponse object with the default status code set to http.StatusNoContent.
func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		statusCode: http.StatusNoContent,
	}
}

// NewHttpResponseAborted constructs a new HttpResponse object with the data from the aborted backend.
// It creates a new ResponseHeader object using the completed status of the endpoint and the Ok status of the backend response.
// The header is then aggregated with the backend response header. The constructed HttpResponse object has the backend's
// status code, header, body, and abort flag set.
func NewHttpResponseAborted(endpoint *Endpoint, httpBackendResponse *HttpBackendResponse) *HttpResponse {
	header := NewResponseHeader(endpoint.Completed(1), httpBackendResponse.Ok())
	header = header.Aggregate(httpBackendResponse.Header())
	return &HttpResponse{
		statusCode: httpBackendResponse.StatusCode(),
		header:     header,
		body:       httpBackendResponse.Body(),
		abort:      true,
	}
}

// NewHttpResponseByStatusCode creates a new HttpResponse object with the provided status code.
// The header is set to a new ResponseHeader object with the "Ok" field set based on the provided status code.
// The "Ok" field is set to true if the status code is within the valid range of 200 to 299, otherwise false.
func NewHttpResponseByStatusCode(statusCode StatusCode) *HttpResponse {
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewResponseHeader(true, statusCode.Ok()),
	}
}

// NewHttpResponseByString creates a new HttpResponse object with the specified StatusCode and body string.
// The StatusCode parameter determines the HTTP status code of the response.
// The body parameter is used to set the body of the response as a string.
// This function returns a pointer to the created HttpResponse object.
func NewHttpResponseByString(statusCode StatusCode, body string) *HttpResponse {
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewResponseHeader(true, statusCode.Ok()),
		body:       NewBodyByString(body),
	}
}

// NewHttpResponseByJson creates a new HttpResponse object with the given status code and body.
// It determines if the status code is within the valid range (200-299) using helper.IsGreaterThanOrEqual
// or helper.IsLessThanOrEqual. It then initializes the header with NewResponseHeader, where the 'complete'
// parameter is set to true if the statusCode is within range, and the 'success' parameter is set to the
// value of 'ok'. The body is initialized with NewBodyByJson, using the provided 'body' parameter.
// The created HttpResponse object is returned.
func NewHttpResponseByJson(statusCode StatusCode, body any) *HttpResponse {
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewResponseHeader(true, statusCode.Ok()),
		body:       NewBodyByJson(body),
	}
}

// NewHttpResponseByCache creates a new HttpResponse object with the status code, header, body,
// and written attribute set based on the given CacheResponse object.
// The XGopenCache and XGopenCacheTTL headers are added to the header of the HttpResponse object
// with appropriate values.
// The body is created by calling the NewBodyByCache function with the CacheResponse object's body.
func NewHttpResponseByCache(cacheResponse *CacheResponse) *HttpResponse {
	header := cacheResponse.Header
	header = header.Set(consts.XGopenCache, helper.SimpleConvertToString(true))
	header = header.Set(consts.XGopenCacheTTL, cacheResponse.TTL())
	return &HttpResponse{
		statusCode: cacheResponse.StatusCode,
		header:     header,
		body:       NewBodyByCache(cacheResponse.Body),
		written:    true,
	}
}

// NewHttpResponseByErr creates a new HttpResponse object with the given status code,
// a failed header, a body constructed from the given path and error, and sets abort to true.
func NewHttpResponseByErr(path string, statusCode StatusCode, err error) *HttpResponse {
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewHeaderFailed(),
		body:       NewBodyByError(path, err),
		abort:      true,
	}
}

// Append appends the provided HttpBackendResponse object to the history list of the HttpResponse object.
// If the provided HttpBackendResponse object is nil, it is ignored and the original HttpResponse object is returned.
// Otherwise, it creates a new list of history by appending the HttpBackendResponse object to the existing history list.
// A new HttpResponse object is returned with the same status code, header, body, and the updated history list.
// Returns the constructed HttpResponse object with the appended HttpBackendResponse object in the history list.
func (r *HttpResponse) Append(httpBackendResponse *HttpBackendResponse) *HttpResponse {
	if helper.IsNil(httpBackendResponse) {
		return r
	}

	history := r.history
	history = append(history, httpBackendResponse)

	return &HttpResponse{
		statusCode: r.StatusCode(),
		header:     r.Header(),
		body:       r.Body(),
		history:    history,
	}
}

// Error sets the status code, header, and body of the HttpResponse object based on the provided path and error.
// If the error contains ErrBadGateway, the status code is set to http.StatusBadGateway.
// If the error contains ErrGatewayTimeout, the status code is set to http.StatusGatewayTimeout.
// Otherwise, the status code is set to http.StatusInternalServerError.
// The header is set to a new Header object with failed status values for consts.XGopenCache, consts.XGopenComplete,
// and consts.XGopenSuccess.
// The body is set to a new Body object constructed from the path and error using the NewBodyByError function.
// The abort property of the HttpResponse object is set to true.
// Returns the constructed HttpResponse object with the updated status code, header, body, and abort property.
func (r *HttpResponse) Error(path string, err error) *HttpResponse {
	var statusCode StatusCode
	if errors.Contains(err, mapper.ErrBadGateway) {
		statusCode = http.StatusBadGateway
	} else if errors.Contains(err, mapper.ErrGatewayTimeout) {
		statusCode = http.StatusGatewayTimeout
	} else {
		statusCode = http.StatusInternalServerError
	}

	return &HttpResponse{
		statusCode: statusCode,
		header:     NewHeaderFailed(),
		body:       NewBodyByError(path, err),
		abort:      true,
	}
}

// Abort returns the value of the `abort` property of the HttpResponse object.
// If `abort` is true, it indicates that the httpResponse should be aborted.
// Returns a boolean value representing the `abort` property.
func (r *HttpResponse) Abort() bool {
	return r.abort
}

// Written returns a boolean value indicating whether the HttpResponse has been written.
// Returns true if the HttpResponse has been written, false otherwise.
func (r *HttpResponse) Written() bool {
	return r.written
}

// StatusCode returns the status code of the HttpResponse object.
func (r *HttpResponse) StatusCode() StatusCode {
	return r.statusCode
}

// Header returns the header of the HttpResponse object.
func (r *HttpResponse) Header() Header {
	return r.header
}

func (r *HttpResponse) ContentType() enum.ContentType {
	if helper.IsNotNil(r.Body()) {
		return r.Body().ContentType()
	}
	return ""
}

// Body returns the body of the HttpResponse object.
func (r *HttpResponse) Body() *Body {
	return r.body
}

// BodyBytes returns the body of the HttpResponse object as a slice of bytes.
// If the body is nil, it returns nil.
// Otherwise, it returns the result of calling the Bytes method on the body object.
func (r *HttpResponse) BodyBytes() []byte {
	if helper.IsNil(r.body) {
		return nil
	}
	return r.Body().Bytes()
}

// HasHistory returns a boolean value indicating whether the HttpResponse object has a history.
// It checks if the history list stored in the HttpResponse object is not empty.
// Returns true if the history list is not empty, false otherwise.
func (r *HttpResponse) HasHistory() bool {
	return helper.IsNotEmpty(r.History())
}

// History returns the history list of backend responses stored in the HttpResponse object.
func (r *HttpResponse) History() httpResponseHistory {
	return r.history
}

// Write writes the HttpResponse object based on the provided Endpoint, HttpRequest, and HttpResponse objects.
// If the HttpResponse object has already been written, it returns the same HttpResponse object without any modifications.
// Otherwise, it creates new variables to store the current status code, header, and body of the HttpResponse object.
// If the HttpResponse object has a history, it calls the writeByHistory method to update the status code, header, and body based on the history.
// Then, it calls the writeByEndpointConfig method to write the response based on the endpoint configuration.
// Returns the modified HttpResponse object after writing the response.
func (r *HttpResponse) Write(endpoint *Endpoint, httpRequest *HttpRequest, httpResponse *HttpResponse) *HttpResponse {
	if r.Written() {
		return r
	}

	statusCode := r.StatusCode()
	header := r.Header()
	body := r.Body()

	if r.HasHistory() {
		statusCode, header, body = r.writeByHistory(endpoint, httpRequest, httpResponse)
	}

	return r.writeByEndpointConfig(endpoint, statusCode, header, body)
}

// Map returns a string representation of the history list stored in the HttpResponse object.
func (r *HttpResponse) Map() string {
	return r.History().Map()
}

// writeByHistory updates the status code, header, and body of the HttpResponse object
// based on the filtered history list, the endpoint, http request, and http response objects.
// It filters the history list based on the provided http request and http response objects.
// It obtains the status code from the filtered history, creates the header based on the completed and success values,
// and aggregates the headers from the filtered history.
// It checks if it needs to aggregate the body based on the endpoint configuration.
// It obtains the body from the filtered history.
// Returns the updated status code, header, and body.
func (r *HttpResponse) writeByHistory(endpoint *Endpoint, httpRequest *HttpRequest, httpResponse *HttpResponse) (
	statusCode StatusCode, header Header, body *Body) {
	endpointResponse := endpoint.Response()

	filteredHistory := r.History().Filter(httpRequest, httpResponse)

	statusCode = filteredHistory.StatusCode()
	header = NewResponseHeader(endpoint.Completed(filteredHistory.Size()), filteredHistory.Success())
	header = header.Aggregate(filteredHistory.Header())

	aggregate := false
	if helper.IsNotNil(endpointResponse) {
		aggregate = endpointResponse.Aggregate()
	}
	body = filteredHistory.Body(aggregate)

	return statusCode, header, body
}

// writeByEndpointConfig updates the body of the HttpResponse object based on the provided Endpoint and Body objects.
// It first retrieves the response object from the endpoint.
// If either the response object or the body object is nil, it returns the original body object.
// If the omitEmpty flag in the response object is true, it removes any empty values from the body object.
// If the nomenclature flag in the response object is true, it converts the body object to the specified case.
// It then determines the content type based on the response object's encode properties and the body object's content
// type.
// If the content type is different from the body object's content type, it converts the body object to the specified
// content type using byte conversion.
// The updated body object is returned.
func (r *HttpResponse) writeByEndpointConfig(endpoint *Endpoint, statusCode StatusCode, header Header, body *Body,
) *HttpResponse {
	return &HttpResponse{
		statusCode: statusCode,
		header:     header,
		body:       r.writeBodyByEndpointConfig(endpoint, body),
		written:    true,
	}
}

// writeBodyByEndpointConfig updates the body of the HttpResponse object based on the provided Endpoint and Body objects.
// It first retrieves the response object from the endpoint.
// If either the response object or the body object is nil, it returns the original body object.
// If the omitEmpty flag in the response object is true, it removes any empty values from the body object.
// If the nomenclature flag in the response object is true, it converts the body object to the specified case.
// It then determines the content type based on the response object's encode properties and the body object's content
// type.
// If the content type is different from the body object's content type, it converts the body object to the specified
// content type using byte conversion.
// The updated body object is returned.
func (r *HttpResponse) writeBodyByEndpointConfig(endpoint *Endpoint, body *Body) *Body {
	endpointResponse := endpoint.Response()
	if helper.IsNil(endpointResponse) || helper.IsNil(body) {
		return body
	}

	if endpointResponse.OmitEmpty() {
		body = body.OmitEmpty()
	}
	if endpointResponse.HasNomenclature() {
		body = body.ToCase(endpointResponse.Nomenclature())
	}

	var contentType enum.ContentType
	if helper.IsNotNil(endpoint.Response()) && endpoint.Response().HasEncode() {
		contentType = endpoint.Response().Encode().ContentType()
	} else if helper.IsNotNil(body) {
		contentType = body.ContentType()
	}

	if helper.IsNotEqualTo(contentType, body.ContentType()) {
		bs := body.BytesByContentType(contentType)
		body = NewBodyByContentType(contentType.String(), helper.SimpleConvertToBuffer(bs))
	}

	return body
}
