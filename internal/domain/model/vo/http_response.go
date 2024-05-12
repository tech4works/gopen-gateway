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

// NewHttpResponseByStatusCodeAndErr creates a new HttpResponse object with the provided status code and error.
// It sets the statusCode field of the HttpResponse object to the provided status code.
// It creates a new HeaderFailed object and assigns it to the header field of the HttpResponse object.
// It generates a new Body object based on the provided path and error using the NewBodyByError function.
// It sets the abort field of the HttpResponse object to true.
// Returns the newly created HttpResponse object.
func NewHttpResponseByStatusCodeAndErr(path string, statusCode StatusCode, err error) *HttpResponse {
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewHeaderFailed(),
		body:       NewBodyByError(path, err),
		abort:      true,
	}
}

// NewHttpResponseByErr creates a new HttpResponse object with the status code determined by the provided error.
// If the error is mapper.ErrBadGateway, the status code will be http.StatusBadGateway.
// If the error is mapper.ErrGatewayTimeout, the status code will be http.StatusGatewayTimeout.
// Otherwise, the default status code will be http.StatusInternalServerError.
// The header will be created using the NewHeaderFailed function.
// The body will be created using the NewBodyByError function with the provided path and error.
// The abort property will be set to true.
// Returns the newly created HttpResponse object.
func NewHttpResponseByErr(path string, err error) *HttpResponse {
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

// Append appends the provided HttpBackendResponse object to the history list of the HttpResponse object.
// If the provided HttpBackendResponse object is nil, it is ignored and the original HttpResponse object is returned.
// Otherwise, it creates a new list of history by appending the HttpBackendResponse object to the existing history list.
// A new HttpResponse object is returned with the same status code, header, body, and the updated history list.
// Returns the constructed HttpResponse object with the appended HttpBackendResponse object in the history list.
func (h *HttpResponse) Append(httpBackendResponse *HttpBackendResponse) *HttpResponse {
	if helper.IsNil(httpBackendResponse) {
		return h
	}

	history := h.history
	history = append(history, httpBackendResponse)

	return &HttpResponse{
		statusCode: h.StatusCode(),
		header:     h.Header(),
		body:       h.Body(),
		history:    history,
	}
}

// Abort returns the value of the `abort` property of the HttpResponse object.
// If `abort` is true, it indicates that the httpResponse should be aborted.
// Returns a boolean value representing the `abort` property.
func (h *HttpResponse) Abort() bool {
	return h.abort
}

// Written returns a boolean value indicating whether the HttpResponse has been written.
// Returns true if the HttpResponse has been written, false otherwise.
func (h *HttpResponse) Written() bool {
	return h.written
}

// StatusCode returns the status code of the HttpResponse object.
func (h *HttpResponse) StatusCode() StatusCode {
	return h.statusCode
}

// Header returns the header of the HttpResponse object.
func (h *HttpResponse) Header() Header {
	return h.header
}

func (h *HttpResponse) ContentType() ContentType {
	if helper.IsNotNil(h.Body()) {
		return h.Body().ContentType()
	}
	return ""
}

func (h *HttpResponse) HasContentEncoding() bool {
	return helper.IsNotNil(h.Body()) && h.Body().HasContentEncoding()
}

func (h *HttpResponse) ContentEncoding() ContentEncoding {
	if helper.IsNotNil(h.Body()) {
		return h.Body().ContentEncoding()
	}
	return ""
}

// Body returns the body of the HttpResponse object.
func (h *HttpResponse) Body() *Body {
	return h.body
}

func (h *HttpResponse) RawBodyBytes() []byte {
	if helper.IsNil(h.body) {
		return nil
	}
	return h.Body().RawBytes()
}

// HasHistory returns a boolean value indicating whether the HttpResponse object has a history.
// It checks if the history list stored in the HttpResponse object is not empty.
// Returns true if the history list is not empty, false otherwise.
func (h *HttpResponse) HasHistory() bool {
	return helper.IsNotEmpty(h.History())
}

// History returns the history list of backend responses stored in the HttpResponse object.
func (h *HttpResponse) History() httpResponseHistory {
	return h.history
}

// Write writes the HttpResponse object based on the provided Endpoint, HttpRequest, and HttpResponse objects.
// If the HttpResponse object has already been written, it returns the same HttpResponse object without any modifications.
// Otherwise, it creates new variables to store the current status code, header, and body of the HttpResponse object.
// If the HttpResponse object has a history, it calls the writeByHistory method to update the status code, header, and body based on the history.
// Then, it calls the writeByEndpointConfig method to write the response based on the endpoint configuration.
// Returns the modified HttpResponse object after writing the response.
func (h *HttpResponse) Write(endpoint *Endpoint, httpRequest *HttpRequest, httpResponse *HttpResponse) *HttpResponse {
	if h.Written() {
		return h
	}

	statusCode := h.StatusCode()
	header := h.Header()
	body := h.Body()

	if h.HasHistory() {
		statusCode, header, body = h.writeByHistory(endpoint, httpRequest, httpResponse)
	}

	return h.writeByEndpointConfig(endpoint, statusCode, header, body)
}

// Map returns a string representation of the history list stored in the HttpResponse object.
func (h *HttpResponse) Map() string {
	return h.History().Map()
}

// writeByHistory updates the status code, header, and body of the HttpResponse object
// based on the filtered history list, the endpoint, http request, and http response objects.
// It filters the history list based on the provided http request and http response objects.
// It obtains the status code from the filtered history, creates the header based on the completed and success values,
// and aggregates the headers from the filtered history.
// It checks if it needs to aggregate the body based on the endpoint configuration.
// It obtains the body from the filtered history.
// Returns the updated status code, header, and body.
func (h *HttpResponse) writeByHistory(endpoint *Endpoint, httpRequest *HttpRequest, httpResponse *HttpResponse) (
	statusCode StatusCode, header Header, body *Body) {
	endpointResponse := endpoint.Response()

	filteredHistory := h.History().Filter(httpRequest, httpResponse)

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

// writeByEndpointConfig updates the HttpResponse object based on the provided Endpoint object,
// StatusCode, Header, and Body parameters. It handles updating the Body and Header based on the Endpoint's configuration.
// It returns a new HttpResponse object with the updated StatusCode, Header, Body, and the written flag set to true.
func (h *HttpResponse) writeByEndpointConfig(endpoint *Endpoint, statusCode StatusCode, header Header, body *Body,
) *HttpResponse {
	body = h.writeBodyByEndpointConfig(endpoint, body)
	return &HttpResponse{
		statusCode: statusCode,
		header:     header.Write(body),
		body:       body,
		written:    true,
	}
}

func (h *HttpResponse) writeBodyByEndpointConfig(endpoint *Endpoint, body *Body) *Body {
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

	contentType := body.ContentType()
	if endpointResponse.HasContentType() {
		contentType = endpointResponse.ContentType()
	}
	contentEncoding := body.ContentEncoding()
	if endpointResponse.HasContentEncoding() {
		contentEncoding = endpointResponse.ContentEncoding()
	}

	return body.ModifyContentType(contentType, contentEncoding)
}
