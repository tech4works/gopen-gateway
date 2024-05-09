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
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"io"
	"net/http"
	"net/url"
)

// HttpBackendRequest represents a type that holds information about a backend HTTP request.
// It contains fields for the host, path, method, header, query, and body of the request.
type HttpBackendRequest struct {
	// host represents the host of the backend httpRequest.
	// It is a string field that contains the host information used for constructing the URL.
	host string
	// path is a string field that represents the path of the backend httpRequest.
	// It contains the path information used for constructing the URL.
	path UrlPath
	// method is a string field that represents the HTTP method to be used for the backend httpRequest.
	// It contains information about the desired HTTP method, such as GET, POST, PUT, DELETE, etc.
	// The value of the method field can be accessed using the Method() `method`.
	method string
	// header represents the header fields of a backend httpRequest.
	header Header
	// query represents the query fields of a backend httpRequest.
	query Query
	// body represents the body of a backend httpRequest.
	body *Body
}

// NewHttpBackendRequest creates a new HttpBackendRequest with the specified backend, httpRequest, and httpResponse.
// The host of the HttpBackendRequest is set to the balanced host from the backend.
// The path is set using the newBackendRequestPath function.
// The method is set to the method used by the backend.
// The header is set using the newBackendRequestHeader function.
// The query is set using the newBackendRequestQuery function.
// The body is set using the newBackendRequestBody function.
// The returned HttpBackendRequest contains the updated fields based on the function parameters.
func NewHttpBackendRequest(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) *HttpBackendRequest {
	return &HttpBackendRequest{
		host:   backend.BalancedHost(),
		path:   newBackendRequestPath(backend, httpRequest, httpResponse),
		method: backend.Method(),
		header: newBackendRequestHeader(backend, httpRequest, httpResponse),
		query:  newBackendRequestQuery(backend, httpRequest, httpResponse),
		body:   newBackendRequestBody(backend, httpRequest, httpResponse),
	}
}

// newBackendRequestPath creates a new UrlPath for the HttpBackendRequest using the specified backend, httpRequest, and httpResponse.
// The UrlPath is initialized with the backend's path and the parameters from the httpRequest.
// If the backend has a non-nil request object, each parameter modifier from the request is applied to the UrlPath.
// The final modified UrlPath is returned as the result.
func newBackendRequestPath(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) UrlPath {
	backendRequest := backend.Request()

	path := NewUrlPath(backend.Path(), httpRequest.Params())
	if helper.IsNotNil(backendRequest) {
		for _, modifier := range backendRequest.ParamModifiers() {
			path = path.Modify(&modifier, httpRequest, httpResponse)
		}
	}
	return path
}

// newBackendRequestHeader initializes the backendRequest variable to configure the header.
// If the backendRequest is nil, the function returns the header from the httpRequest.
// If the backendRequest wants to omit the header, an empty header is returned.
// The header is instantiated from the existing information in the httpRequest.
// The header is mapped according to the backendRequest's headerMapper configuration.
// The header is projected according to the backendRequest's headerProjection configuration.
// The header is modified based on the backendRequest's headerModifiers.
// The function returns the modified header.
func newBackendRequestHeader(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) Header {
	backendRequest := backend.Request()

	if helper.IsNil(backendRequest) {
		return httpRequest.Header()
	} else if backendRequest.OmitHeader() {
		return NewEmptyHeader()
	}

	header := httpRequest.Header()
	header = header.Map(backendRequest.HeaderMapper())
	header = header.Projection(backendRequest.HeaderProjection())
	for _, modifier := range backendRequest.HeaderModifiers() {
		header = header.Modify(&modifier, httpRequest, httpResponse)
	}

	return header
}

// newBackendRequestQuery creates a new Query object based on the specified backend, httpRequest, and httpResponse.
// It initializes the backendRequest to be used for configuring the query.
// If the backendRequest is nil or wants to omit the query, it returns the query from the httpRequest.
// It maps and projects the query based on the configuration from the backendRequest.
// It runs the query modifications based on the backendRequest's queryModifiers.
// The final query is returned.
func newBackendRequestQuery(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) Query {
	backendRequest := backend.Request()

	if helper.IsNil(backendRequest) {
		return httpRequest.Query()
	} else if backendRequest.OmitQuery() {
		return NewEmptyQuery()
	}

	query := httpRequest.Query()
	query = query.Map(backendRequest.QueryMapper())
	query = query.Projection(backendRequest.QueryProjection())
	for _, modifier := range backendRequest.QueryModifiers() {
		query = query.Modify(&modifier, httpRequest, httpResponse)
	}

	return query
}

// newBackendRequestBody creates a new Body object for the HttpBackendRequest based on the provided backend,
// httpRequest, and httpResponse.
// If the backend's request is nil, it returns the body from the httpRequest.
// If the httpRequest's body is nil or the backend's request has OmitBody set to true, it returns nil.
// Otherwise, it performs the following operations:
//   - Maps the body based on the backend's request.BodyMapper()
//   - Projects the body based on the backend's request.BodyProjection()
//   - Modifies the body using the backend's request.BodyModifiers()
//
// The resulting body is then returned.
func newBackendRequestBody(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) *Body {
	backendRequest := backend.Request()

	if helper.IsNil(backendRequest) {
		return httpRequest.Body()
	} else if helper.IsNil(httpRequest.Body()) || backendRequest.OmitBody() {
		return nil
	}

	body := httpRequest.Body()
	body = body.Map(backendRequest.BodyMapper())
	body = body.Projection(backendRequest.BodyProjection())
	for _, modifier := range backendRequest.BodyModifiers() {
		body = body.Modify(&modifier, httpRequest, httpResponse)
	}

	return body
}

// Path returns the URL path of the HttpBackendRequest instance.
func (h *HttpBackendRequest) Path() UrlPath {
	return h.path
}

// Url returns the full URL string representation of the HttpBackendRequest instance.
// It concatenates the host and the URL path using the `fmt.Sprint` function.
// The resulting URL string represents the complete URL of the backend HTTP request.
func (h *HttpBackendRequest) Url() string {
	return fmt.Sprint(h.host, h.Path().String())
}

// Params returns the `Params` field of the `HttpBackendRequest` instance.
// `Params` represents a map of string key-value pairs that store additional parameters for the URL path.
// The `Params` method allows you to access and manipulate these parameters.
// The returned `Params` field is of type `map[string]string`.
func (h *HttpBackendRequest) Params() Params {
	return h.Path().Params()
}

// Method returns the method of the HttpBackendRequest instance.
func (h *HttpBackendRequest) Method() string {
	return h.method
}

// Header returns the Header of the HttpBackendRequest instance.
func (h *HttpBackendRequest) Header() Header {
	return h.header
}

// Query returns the query of the HttpBackendRequest.
func (h *HttpBackendRequest) Query() Query {
	return h.query
}

// RawQuery encodes the query parameters into a string representation.
// It returns the encoded query string that can be appended to the URL.
func (h *HttpBackendRequest) RawQuery() string {
	return url.Values(h.query).Encode()
}

// Body returns the `Body` field of the `httpBackendRequest` instance.
// The `Body` field represents the httpRequest body.
// The `Body` method allows you to access the httpRequest body for further manipulation or inspection.
func (h *HttpBackendRequest) Body() *Body {
	return h.body
}

// BodyToReadCloser returns the body to send as an `io.ReadCloser` interface.
// If `omitRequestBody` is set to `true` or `body` is `nil`, it returns `nil`.
//
// It converts the body to bytes using the desired encoding (XML, JSON, TEXT/PLAIN) based on `Content-Type` config.
//
// If there is an error during the conversion, it returns `nil`.
//
// Finally, it returns the `io.ReadCloser` interface with the bytes of the body.
func (h *HttpBackendRequest) BodyToReadCloser() io.ReadCloser {
	if helper.IsNil(h.body) {
		return nil
	}
	// todo: aqui podemos futuramente colocar encode de httpRequest customizado
	return io.NopCloser(h.body.Buffer())
}

// NetHttp creates a *http.Request instance using the information from the HttpBackendRequest instance.
// It sets the context, method, URL, body, header, and query of the created request.
// It returns the created *http.Request instance and an error, if any occurred.
func (h *HttpBackendRequest) NetHttp(ctx context.Context) (*http.Request, error) {
	netHttpRequest, err := http.NewRequestWithContext(ctx, h.Method(), h.Url(), h.BodyToReadCloser())
	if helper.IsNotNil(err) {
		return nil, err
	}

	netHttpRequest.Header = h.Header().Http()
	netHttpRequest.URL.RawQuery = h.RawQuery()
	return netHttpRequest, nil
}

// Map returns a map[string]any containing the header, params, query, and body
// fields of the HttpBackendRequest instance. If the body is not nil, it will
// be converted to an interface{} and included in the map.
func (h *HttpBackendRequest) Map() any {
	var body any
	if helper.IsNotNil(body) {
		body = h.Body().Interface()
	}
	return map[string]any{
		"header": h.Header(),
		"params": h.Params(),
		"query":  h.Query(),
		"body":   body,
	}
}
