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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/gin-gonic/gin"
	"io"
)

// HttpRequest represents an HTTP httpRequest and contains information such as the httpRequest path, URL, method, header,
// parameters, query parameters, httpRequest body, and httpRequest history.
type HttpRequest struct {
	// path represents the URI of the HttpRequest.
	// It is a string field in the HttpRequest struct.
	path string
	// url represents the URL of the HttpRequest.
	url string
	// method represents the HTTP method of the HttpRequest.
	// It is a string field in the HttpRequest struct.
	method string
	// header represents the httpRequest header of a HttpRequest.
	header Header
	// params represents the parameters associated with the httpRequest.
	params Params
	// query represents the query parameter map of the HttpRequest.
	query Query
	// Body represents the body of an HTTP httpRequest.
	// It is a field in the HttpRequest struct.
	body *Body
	// history represents the history of backend requests made by the HttpRequest object.
	// It is a slice of httpBackendRequest objects.
	history []*httpBackendRequest
}

func NewHttpRequest(gin *gin.Context) *HttpRequest {
	// instanciamos o query VO para obter funções de montagem da url por ele
	query := NewQuery(gin.Request.URL.Query())

	// preparamos a url ordenando as chaves de busca
	url := gin.Request.URL.Path
	if helper.IsNotEmpty(gin.Request.URL.RawQuery) {
		url += "?" + query.Encode()
	}

	// obtemos os bytes da requisição
	bodyBytes, _ := io.ReadAll(gin.Request.Body)
	bodyBuffer := bytes.NewBuffer(bodyBytes)
	gin.Request.Body = io.NopCloser(bodyBuffer)

	// instanciamos os valores necessários do header para montar o body
	contentType := gin.GetHeader("Content-Type")
	contentEncoding := gin.GetHeader("Content-Encoding")

	// montamos o VO de requisição
	return &HttpRequest{
		path:   gin.Request.URL.Path,
		url:    url,
		method: gin.Request.Method,
		header: NewHeader(gin.Request.Header),
		params: NewParams(gin.Params),
		query:  query,
		body:   NewBody(contentType, contentEncoding, bodyBuffer),
	}
}

// SetHeader takes a Header object as an argument and returns a new HttpRequest with the provided header.
// The other fields of the new HttpRequest remain unchanged.
//
// Parameters:
// header - The Header object to be set in the new HttpRequest.
//
// Returns:
// HttpRequest - A new HttpRequest instance with the updated header and the original values for the other fields.
func (r *HttpRequest) SetHeader(header Header) *HttpRequest {
	return &HttpRequest{
		path:    r.path,
		url:     r.url,
		method:  r.method,
		header:  header,
		params:  r.params,
		query:   r.query,
		body:    r.body,
		history: r.history,
	}
}

// ModifyHeader creates a new HttpRequest from an existing one with modifications to the httpRequest header.
// Also modifies the httpRequest history by adding a new backendRequestVO at the end.
// 'header' which is an instance of Header will replace the existing httpRequest header.
// 'backendRequestVO' is an instance of a backend httpRequest object to be changed in the history array.
// The function returns a new instance of HttpRequest with the modified header and history.
func (r *HttpRequest) ModifyHeader(header Header, httpBackendRequestVO *httpBackendRequest) *HttpRequest {
	history := r.history
	history[len(history)-1] = httpBackendRequestVO

	return &HttpRequest{
		path:    r.path,
		url:     r.url,
		method:  r.method,
		header:  header,
		params:  r.params,
		query:   r.query,
		body:    r.body,
		history: history,
	}
}

// ModifyParams takes a Params object and a backend httpRequest object of type httpBackendRequest.
// It modifies the history of the httpRequest by replacing the last element with the provided httpBackendRequest object,
// and updates the params of the original httpRequest.
// It returns a new HttpRequest with the updated history and params.
func (r *HttpRequest) ModifyParams(params Params, httpBackendRequestVO *httpBackendRequest) *HttpRequest {
	history := r.history
	history[len(history)-1] = httpBackendRequestVO

	return &HttpRequest{
		path:    r.path,
		url:     r.url,
		method:  r.method,
		header:  r.header,
		params:  params,
		query:   r.query,
		body:    r.body,
		history: history,
	}
}

// ModifyQuery takes a Query and a httpBackendRequest as arguments and modifies the HttpRequest's history field.
// It returns a new HttpRequest object with updated values.
//
// It allocates a new slice, then copies the values of the original HttpRequest's fields into the new one,
// and modifies its history field by updating the last element with the provided httpBackendRequest.
// The query field is also updated with the provided query argument.
//
// The original HttpRequest is not modified. All modifications are performed on the new HttpRequest object.
//
// Parameters:
// query - A Query object that is used to update the new HttpRequest's query field.
// backendRequestVO - A httpBackendRequest object that is used to update the last element of the history field in the new HttpRequest.
//
// Returns:
// HttpRequest - a new HttpRequest object with updated query and history fields.
func (r *HttpRequest) ModifyQuery(query Query, httpBackendRequestVO *httpBackendRequest) *HttpRequest {
	history := r.history
	history[len(history)-1] = httpBackendRequestVO

	return &HttpRequest{
		path:    r.path,
		url:     r.url,
		method:  r.method,
		header:  r.header,
		params:  r.params,
		query:   query,
		body:    r.body,
		history: history,
	}
}

// ModifyBody takes a Body and a backendRequestVO as parameters,
// replaces the last item in the history slice of the HttpRequest receiver with backendRequestVO, and
// returns a new HttpRequest with the modified history, preserving all other fields.
//
// Parameters:
//
//	body - The new body for the HttpRequest.
//	backendRequestVO - The backend httpRequest value object to be added in the httpRequest history.
//
// Returns:
// A new HttpRequest instance with the updated history.
func (r *HttpRequest) ModifyBody(body *Body, backendRequestVO *httpBackendRequest) *HttpRequest {
	history := r.history
	history[len(history)-1] = backendRequestVO
	return &HttpRequest{
		path:    r.path,
		url:     r.url,
		method:  r.method,
		header:  r.header,
		params:  r.params,
		query:   r.query,
		body:    body,
		history: history,
	}
}

// Append is a method for the HttpRequest type. It adds a httpBackendRequest to the history of the HttpRequest
// and returns the updated HttpRequest.
//
// The method receives a httpBackendRequest as an argument.
// It constructs a new HttpRequest with the same parameters as the original httpRequest (url, method, header, params, query, body),
// but with the httpBackendRequest added to the history.
//
// Parameters:
//
//	httpBackendRequest : The httpRequest to append to the history of the HttpRequest.
//
// Returns:
//
//	HttpRequest - A new HttpRequest with the httpBackendRequest added to its history.
func (r *HttpRequest) Append(backendRequest *httpBackendRequest) *HttpRequest {
	return &HttpRequest{
		path:    r.path,
		url:     r.url,
		method:  r.method,
		header:  r.header,
		params:  r.params,
		query:   r.query,
		body:    r.body,
		history: append(r.history, backendRequest),
	}
}

// LastHttpBackendRequest returns the last httpBackendRequest object in the httpRequest's history array.
func (r *HttpRequest) LastHttpBackendRequest() *httpBackendRequest {
	return r.history[len(r.history)-1]
}

// Url returns the URL of the HttpRequest.
// It retrieves the value of the `url` field from the HttpRequest struct.
func (r *HttpRequest) Url() string {
	return r.url
}

// Path returns the URI of the HttpRequest.
// It retrieves the value of the `path` field from the HttpRequest struct.
func (r *HttpRequest) Path() string {
	return r.path
}

// Method returns the HTTP method of the HttpRequest.
// It retrieves the value of the `method` field from the HttpRequest struct.
func (r *HttpRequest) Method() string {
	return r.method
}

// Header returns the httpRequest header of a HttpRequest.
// It returns an instance of Header.
func (r *HttpRequest) Header() Header {
	return r.header
}

// Params returns the parameters associated with the httpRequest.
// The returned value is of type Params.
func (r *HttpRequest) Params() Params {
	return r.params
}

// Query returns the query parameter map of the HttpRequest.
func (r *HttpRequest) Query() Query {
	return r.query
}

// Body returns the body of the httpRequest.
func (r *HttpRequest) Body() *Body {
	return r.body
}

// Json takes no arguments and returns a string representation of the HttpRequest object.
// It iterates over the HttpRequest's history field and calls the Json method on each element,
// appending the result to the evalHistory slice.
// It constructs a map with the following keys:
//   - "header": The header field of the HttpRequest.
//   - "params": The params field of the HttpRequest.
//   - "query": The query field of the HttpRequest.
//   - "body": The Interface method of the body field of the HttpRequest.
//   - "history": The evalHistory slice.
//
// It then calls the SimpleConvertToString method of the helper package, passing the map as an argument,
// and returns the resulting string.
func (r *HttpRequest) Json() string {
	var evalHistory []any
	for _, backendRequestVO := range r.history {
		evalHistory = append(evalHistory, backendRequestVO.Eval())
	}
	var evalBody any
	if helper.IsNotNil(r.body) {
		evalBody = r.body.Interface()
	}

	mapEval := map[string]any{
		"header":  r.header,
		"params":  r.params,
		"query":   r.query,
		"body":    evalBody,
		"history": evalHistory,
	}
	return helper.SimpleConvertToString(mapEval)
}
