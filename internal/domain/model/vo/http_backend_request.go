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

// NewHttpBackendRequest creates a new HttpBackendRequest using the specified backend, httpRequest, and httpResponse.
// The host of the HttpBackendRequest is set to the balanced host of the backend.
// The path of the HttpBackendRequest is set using the newBackendRequestPath function.
// The method of the HttpBackendRequest is set to the method of the backend.
// The header of the HttpBackendRequest is set using the newBackendRequestHeader function.
// The query of the HttpBackendRequest is set using the newBackendRequestQuery function.
// The body of the HttpBackendRequest is set using the newBackendRequestBody function.
// The newly created HttpBackendRequest is returned as a pointer.
func NewHttpBackendRequest(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) *HttpBackendRequest {
	body := newBackendRequestBody(backend, httpRequest, httpResponse)
	return &HttpBackendRequest{
		host:   backend.BalancedHost(),
		path:   newBackendRequestPath(backend, httpRequest, httpResponse),
		method: backend.Method(),
		header: newBackendRequestHeader(backend, body, httpRequest, httpResponse),
		query:  newBackendRequestQuery(backend, httpRequest, httpResponse),
		body:   body,
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

// newBackendRequestHeader creates a new header for the backend request using the specified backend, body, httpRequest,
// and httpResponse.
// The function checks if the backend request is nil. If it is, the header is set using the httpRequest's header.
// If the backend request has an omit header flag, an empty header is created.
// Otherwise, the httpRequest's header is used and mapped according to the backend request's header mapper function.
// The header is then projected based on the backend request's header projection.
// Finally, the header is modified using each modifier in the backend request's header modifiers slice.
// The modified header is written to the specified body and returned.
func newBackendRequestHeader(backend *Backend, body *Body, httpRequest *HttpRequest, httpResponse *HttpResponse) Header {
	backendRequest := backend.Request()

	var header Header
	if helper.IsNil(backendRequest) {
		header = httpRequest.Header()
	} else if backendRequest.OmitHeader() {
		header = httpRequest.Header().OnlyMandatoryKeys()
	} else {
		header = httpRequest.Header()
		header = header.Map(backendRequest.HeaderMapper())
		header = header.Projection(backendRequest.HeaderProjection())
		for _, modifier := range backendRequest.HeaderModifiers() {
			header = header.Modify(&modifier, httpRequest, httpResponse)
		}
	}

	return header.Write(body)
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

// newBackendRequestBody returns the body of the backend request based on the given backend, httpRequest, and httpResponse.
// If the backend request is nil, it returns the body of the httpRequest.
// If the httpRequest body is nil or the backend request specifies omitting the body, it returns nil.
// It applies body mapper, projection, and modifiers specified in the backend request.
// It also adjusts the content type and content encoding based on the backend request.
// Returns the modified body as a pointer.
// PARAMETERS:
// - backend: The backend configuration for the request.
// - httpRequest: The HTTP request object.
// - httpResponse: The HTTP response object.
// RETURNS:
//   - *Body: The modified body of the backend request.
//     nil if the backend request specifies omitting the body or if the httpRequest body is nil.
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

	contentType := body.ContentType()
	if backendRequest.HasContentType() {
		contentType = backendRequest.ContentType()
	}
	contentEncoding := body.ContentEncoding()
	if backendRequest.HasContentEncoding() {
		contentEncoding = backendRequest.ContentEncoding()
	}

	return body.ModifyContentType(contentType, contentEncoding)
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

func (h *HttpBackendRequest) BodyToReadCloser() io.ReadCloser {
	if helper.IsNil(h.Body()) {
		return nil
	}
	return io.NopCloser(h.Body().Buffer())
}

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
