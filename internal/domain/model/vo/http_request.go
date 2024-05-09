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
	"bytes"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"time"
)

// HttpRequest represents an HTTP request.
type HttpRequest struct {
	// url represents the URL of the HttpRequest.
	url string
	// path represents the URL path of the HttpRequest.
	path UrlPath
	// method represents the HTTP method of the HttpRequest. It is a string field on the HttpRequest struct.
	method string
	// Header represents the HTTP headers of an HttpRequest.
	header Header
	// query represents the query parameter map of the HttpRequest struct.
	// It is an instance of the Query type.
	query Query
	// Body represents the body of an HTTP request.
	// It is a pointer to the Body struct.
	body *Body
	// history is a field representing a collection of HttpBackendRequest objects.
	// It stores the history of backend requests made by an HttpRequest.
	// The field is a slice of pointers to HttpBackendRequest objects.
	// Each HttpBackendRequest object in the history represents a single backend request.
	history []*HttpBackendRequest
}

// NewHttpRequest creates a new instance of HttpRequest.
// It initializes the HttpRequest with data from the gin.Context parameter.
// The function performs various operations to populate the fields of the HttpRequest struct,
// such as updating headers, constructing the URL path, processing and parsing the query parameters,
// reading the request body, and setting the HTTP method.
// The function returns the created HttpRequest object.
func NewHttpRequest(gin *gin.Context) *HttpRequest {
	header := NewHeader(gin.Request.Header)
	header = header.Add(consts.XForwardedFor, gin.ClientIP())
	if helper.IsEmpty(header.Get(consts.XTraceId)) {
		u := uuid.New().String()
		unixNano := time.Now().UnixNano()
		header = header.Set(consts.XTraceId, fmt.Sprintf("%s%d", u[:8], unixNano)[:16])
	}

	query := NewQuery(gin.Request.URL.Query())
	url := gin.Request.URL.Path
	if helper.IsNotEmpty(query) {
		url += "?" + query.Encode()
	}

	bodyBytes, _ := io.ReadAll(gin.Request.Body)
	gin.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	params := map[string]string{}
	for _, param := range gin.Params {
		params[param.Key] = param.Value
	}

	contentType := gin.GetHeader("Content-Type")
	contentEncoding := gin.GetHeader("Content-Encoding")
	return &HttpRequest{
		path:   NewUrlPath(gin.FullPath(), params),
		url:    url,
		method: gin.Request.Method,
		header: header,
		query:  query,
		body:   NewBody(contentType, contentEncoding, bytes.NewBuffer(bodyBytes)),
	}
}

// AddOnHeader adds a new header key-value pair to the existing HttpRequest and returns a new instance of HttpRequest
// with the updated header.
// The new header is obtained by calling the Add() method of the existing HttpRequest's header.
// The new HttpRequest instance is created by copying the existing HttpRequest's properties and updating the header
// field with the new header.
// The other fields of the new HttpRequest remain the same as the existing HttpRequest.
// This method does not modify the existing HttpRequest. It creates a new instance with the updated header.
// The method takes two string arguments: key and value, representing the key-value pair to be added to the header.
// The method returns a pointer to the new HttpRequest instance.
func (h *HttpRequest) AddOnHeader(key, value string) *HttpRequest {
	return &HttpRequest{
		url:     h.url,
		path:    h.path,
		method:  h.method,
		header:  h.Header().Add(key, value),
		query:   h.query,
		body:    h.body,
		history: h.history,
	}
}

// SetOnHeader sets the value of the specified key in the header of the existing HttpRequest and returns a new instance
// of HttpRequest with the updated header.
// The new header is obtained by calling the Set() method of the existing HttpRequest's header with the specified key
// and value.
// The new HttpRequest instance is created by copying the existing HttpRequest's properties and updating the header
// field with the new header.
// The other fields of the new HttpRequest remain the same as the existing HttpRequest.
// This method does not modify the existing HttpRequest. It creates a new instance with the updated header.
// The method takes two string arguments: key and value, representing the key-value pair to be set in the header.
// The method returns a pointer to the new HttpRequest instance.
func (h *HttpRequest) SetOnHeader(key, value string) *HttpRequest {
	return &HttpRequest{
		url:     h.url,
		path:    h.path,
		method:  h.method,
		header:  h.Header().Set(key, value),
		query:   h.query,
		body:    h.body,
		history: h.history,
	}
}

// Modify creates a new instance of HttpRequest with modified properties based on the provided BackendRequest
// and HttpResponse.
// If the backendRequest argument is nil, it returns the existing HttpRequest without any modifications.
// The new HttpRequest instance is created by copying the existing HttpRequest's properties and updating the following
// fields:
// - path: It is modified by calling the modifyPath function, passing the backendRequest and httpResponse arguments.
// - header: It is modified by calling the modifyHeader function, passing the backendRequest and httpResponse arguments.
// - query: It is modified by calling the modifyQuery function, passing the backendRequest and httpResponse arguments.
// - body: It is modified by calling the modifyBody function, passing the backendRequest and httpResponse arguments.
// The other fields of the new HttpRequest remain the same as the existing HttpRequest.
// This method does not modify the existing HttpRequest. It creates a new instance with the modified properties.
// The method takes two pointer arguments: backendRequest and httpResponse, representing the BackendRequest and
// HttpResponse objects respectively.
// The method returns a pointer to the new HttpRequest instance.
func (h *HttpRequest) Modify(backendRequest *BackendRequest, httpResponse *HttpResponse) *HttpRequest {
	if helper.IsNil(backendRequest) {
		return h
	}
	return &HttpRequest{
		url:     h.url,
		path:    h.modifyPath(backendRequest, httpResponse),
		method:  h.method,
		header:  h.modifyHeader(backendRequest, httpResponse),
		query:   h.modifyQuery(backendRequest, httpResponse),
		body:    h.modifyBody(backendRequest, httpResponse),
		history: h.history,
	}
}

// Append appends the given HttpBackendRequest to the history slice of the HttpRequest and
// returns a new instance of HttpRequest with the updated history.
// The new HttpRequest instance is created by copying the existing HttpRequest's properties and
// appending the httpBackendRequest to the history slice.
// The other fields of the new HttpRequest remain the same as the existing HttpRequest.
// This method does not modify the existing HttpRequest. It creates a new instance with the updated history.
// The method takes one argument: httpBackendRequest, which is a pointer to the HttpBackendRequest to be appended
// to the history.
// The method returns a pointer to the new HttpRequest instance.
func (h *HttpRequest) Append(httpBackendRequest *HttpBackendRequest) *HttpRequest {
	return &HttpRequest{
		path:    h.path,
		url:     h.url,
		method:  h.method,
		header:  h.header,
		query:   h.query,
		body:    h.body,
		history: append(h.history, httpBackendRequest),
	}
}

// Url returns the URL of the HttpRequest.
// It retrieves the value of the `url` field from the HttpRequest struct.
func (h *HttpRequest) Url() string {
	return h.url
}

// Path returns the UrlPath of the existing HttpRequest.
func (h *HttpRequest) Path() UrlPath {
	return h.path
}

// Method returns the HTTP method of the HttpRequest.
// It retrieves the value of the `method` field from the HttpRequest struct.
func (h *HttpRequest) Method() string {
	return h.method
}

// Header returns the HTTP header of a HttpRequest.
// It returns an instance of Header.
func (h *HttpRequest) Header() Header {
	return h.header
}

// Params returns the parameters of the HttpRequest's path as a Params instance.
func (h *HttpRequest) Params() Params {
	return h.Path().Params()
}

// Query returns the query parameter map of the HttpRequest.
func (h *HttpRequest) Query() Query {
	return h.query
}

// Body returns the body of the httpRequest.
func (h *HttpRequest) Body() *Body {
	return h.body
}

// History returns the slice of HttpBackendRequest objects representing the history of HttpBackendRequest objects
func (h *HttpRequest) History() []*HttpBackendRequest {
	return h.history
}

// HeaderSize calculates and returns the size of the HTTP header in bytes.
// The size is determined by iterating over the key-value pairs in the h.Header() map,
// where each key represents a header field and each value is a slice of header values.
// The size is calculated by adding the length of each key, including the ': ' separator,
// and then adding the length of each value, separated by ', ' if there are multiple values.
// Finally, the size is increased by 2 bytes for each line terminator '\r\n', and the total size is returned.
// The method does not modify the HttpRequest instance.
func (h *HttpRequest) HeaderSize() int {
	size := 0
	for name, values := range h.Header() {
		size += len(name) + 2
		for _, value := range values {
			size += len(value)
			size += 2
		}
		size -= 2
		size += 2
	}
	size += 2
	return size
}

// Map creates a string representation of the HttpRequest object.
//
// The method iterates over the history of the HttpRequest and
// maps each HttpBackendRequest to a string using the Map method.
// The mapped strings are then stored in the 'history' slice.
//
// The 'body' variable is used to store the body of the HttpRequest
// if it is not nil. The body is obtained by calling the Interface
// method of the Body instance.
//
// The method returns a string representation of the HttpRequest.
// The representation is obtained by converting a map containing
// the following key-value pairs to a string using the
// SimpleConvertToString helper method:
//   - "header": the Header object of the HttpRequest
//   - "params": the Params object of the HttpRequest
//   - "query": the Query object of the HttpRequest
//   - "body": the 'body' variable
//   - "history": the 'history' slice
func (h *HttpRequest) Map() string {
	var history []any
	for _, httpBackendRequest := range h.History() {
		history = append(history, httpBackendRequest.Map())
	}
	var body any
	if helper.IsNotNil(h.body) {
		body = h.Body().Interface()
	}
	return helper.SimpleConvertToString(map[string]any{
		"header":  h.Header(),
		"params":  h.Params(),
		"query":   h.Query(),
		"body":    body,
		"history": history,
	})
}

// modifyHeader modifies the header of the current HttpRequest based on the provided BackendRequest
// and HttpResponse.
// It retrieves the existing header of the current HttpRequest using the Header() method.
// It then iterates through the list of header modifiers obtained from the BackendRequest and applies
// the ones that need to be propagated.
// For each modifier that needs to be propagated, it calls the Modify() method of the header, passing
// the current modifier, current HttpRequest, and HttpResponse.
// The modified header is assigned back to the 'header' variable.
// Finally, it returns the modified header.
// The method takes two arguments: backendRequest of type *BackendRequest and httpResponse of type
// *HttpResponse.
// The method returns a Header instance.
func (h *HttpRequest) modifyHeader(backendRequest *BackendRequest, httpResponse *HttpResponse) Header {
	header := h.Header()
	for _, modifier := range backendRequest.HeaderModifiers() {
		if modifier.Propagate() {
			header = header.Modify(&modifier, h, httpResponse)
		}
	}
	return header
}

// modifyPath takes in a BackendRequest and a HttpResponse and returns a copy of the current HttpRequest's path.
// The method iterates through the ParamModifiers of the BackendRequest and, for each modifier with Propagate() set to
// true, it calls the Modify() method of the current HttpRequest's path, passing the modifier, the current HttpRequest,
// and the HttpResponse as arguments.
// The returned modified path is then returned by the method.
// The method does not modify the current HttpRequest or any of its fields.
// The method has three parameters:
// - backendRequest: a pointer to the BackendRequest that contains the ParamModifiers to iterate through
// - httpResponse: a pointer to the HttpResponse that may be used during the Modify() method calls
// The method returns a copy of the current HttpRequest's path, possibly modified by the ParamModifiers.
func (h *HttpRequest) modifyPath(backendRequest *BackendRequest, httpResponse *HttpResponse) UrlPath {
	path := h.Path()
	for _, modifier := range backendRequest.ParamModifiers() {
		if modifier.Propagate() {
			path = path.Modify(&modifier, h, httpResponse)
		}
	}
	return path
}

// modifyQuery modifies the Query object of the HttpRequest based on the QueryModifiers of the BackendRequest
// It takes two parameters:
// - backendRequest: a pointer to a BackendRequest object that contains the QueryModifiers
// - httpResponse: a pointer to a HttpResponse object used in the Modify function of the QueryModifiers
// The method iterates over each QueryModifier in the BackendRequest's QueryModifiers.
// If the modifier's Propagate method returns true, the query object is modified by calling the Modify method
// passing the modifier, current HttpRequest, and httpResponse as arguments.
// The resulting modified query object is returned by the method.
// The original HttpRequest's Query object is obtained by calling the Query method of the HttpRequest.
// This method does not modify the existing HttpRequest, it returns the modified query as a separate Query object.
// The returned Query object represents the updated state of the HttpRequest's query.
// The method returns the modified Query object.
func (h *HttpRequest) modifyQuery(backendRequest *BackendRequest, httpResponse *HttpResponse) Query {
	query := h.Query()
	for _, modifier := range backendRequest.QueryModifiers() {
		if modifier.Propagate() {
			query = query.Modify(&modifier, h, httpResponse)
		}
	}
	return query
}

// modifyBody extracts the body from the existing HttpRequest and initializes it to the variable 'body'.
// If the body is nil, it returns nil.
// It iterates through the body modifiers obtained from the BackendRequest and checks for the 'Propagate' modifier.
// If the modifier is set to true, it modifies the body by calling the 'Modify' method on it, passing the modifier,
// the existing HttpRequest, and the HttpResponse.
// Finally, it returns the modified or unmodified body.
// This method does not modify the existing HttpRequest.
// The method takes two pointer arguments: backendRequest of type *BackendRequest and httpResponse of type *HttpResponse.
// The method returns a pointer to the body of the existing HttpRequest.
func (h *HttpRequest) modifyBody(backendRequest *BackendRequest, httpResponse *HttpResponse) *Body {
	body := h.Body()
	if helper.IsNil(body) {
		return nil
	}

	for _, modifier := range backendRequest.BodyModifiers() {
		if modifier.Propagate() {
			body = body.Modify(&modifier, h, httpResponse)
		}
	}

	return body
}
