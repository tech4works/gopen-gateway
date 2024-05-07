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
	configVO "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/enum"
	"github.com/gin-gonic/gin"
	"io"
)

type HttpRequest struct {
	url     string
	path    UrlPath
	method  string
	header  Header
	query   Query
	body    *Body
	history []*HttpBackendRequest
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

	params := map[string]string{}
	for _, param := range gin.Params {
		params[param.Key] = param.Value
	}

	// montamos o VO de requisição
	return &HttpRequest{
		path:   NewUrlPath(gin.FullPath(), params),
		url:    url,
		method: gin.Request.Method,
		header: NewHeader(gin.Request.Header),
		query:  query,
		body:   NewBody(contentType, contentEncoding, bodyBuffer),
	}
}

func (h *HttpRequest) AddHeader(key, value string) *HttpRequest {
	newHeader := h.Header().Add(key, value)
	return &HttpRequest{
		url:     h.url,
		path:    h.path,
		method:  h.method,
		header:  newHeader,
		query:   h.query,
		body:    h.body,
		history: h.history,
	}
}

func (h *HttpRequest) SetHeader(key, value string) *HttpRequest {
	newHeader := h.Header().Set(key, value)
	return &HttpRequest{
		url:     h.url,
		path:    h.path,
		method:  h.method,
		header:  newHeader,
		query:   h.query,
		body:    h.body,
		history: h.history,
	}
}

func (h *HttpRequest) Modify(backendRequest *configVO.BackendRequest, httpResponse *HttpResponse) *HttpRequest {
	// se for nil, retornamos a request atual
	if helper.IsNil(backendRequest) {
		return h
	}
	// chamamos o modify do header
	header := h.modifyHeader(backendRequest, httpResponse)
	// chamamos o modify do path
	path := h.modifyPath(backendRequest, httpResponse)
	// chamamos o modify do query
	query := h.modifyQuery(backendRequest, httpResponse)
	// chamamos o modify do body
	body := h.modifyBody(backendRequest, httpResponse)

	// montamos o novo HttpRequest com os valores possivelmente modificados
	return &HttpRequest{
		url:     h.url,
		path:    path,
		method:  h.method,
		header:  header,
		query:   query,
		body:    body,
		history: h.history,
	}
}

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

func (h *HttpRequest) Path() UrlPath {
	return h.path
}

func (h *HttpRequest) Params() Params {
	return h.Path().Params()
}

// Method returns the HTTP method of the HttpRequest.
// It retrieves the value of the `method` field from the HttpRequest struct.
func (h *HttpRequest) Method() string {
	return h.method
}

// Header returns the httpRequest header of a HttpRequest.
// It returns an instance of Header.
func (h *HttpRequest) Header() Header {
	return h.header
}

// Query returns the query parameter map of the HttpRequest.
func (h *HttpRequest) Query() Query {
	return h.query
}

// Body returns the body of the httpRequest.
func (h *HttpRequest) Body() *Body {
	return h.body
}

func (h *HttpRequest) History() []*HttpBackendRequest {
	return h.history
}

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

func (h *HttpRequest) modifyHeader(backendRequest *configVO.BackendRequest, httpResponse *HttpResponse) Header {
	// instanciamos o header com o valor atual
	header := h.Header()
	// iteramos os modificadores de header
	for _, modifier := range backendRequest.HeaderModifiers() {
		// se for propagate executamos
		if modifier.Propagate() {
			header = header.Modify(NewModify(&modifier, h, httpResponse))
		}
	}
	// retornamos o header possivelmente alterado
	return header
}

func (h *HttpRequest) modifyPath(backendRequest *configVO.BackendRequest, httpResponse *HttpResponse) UrlPath {
	// instanciamos o path com o valor atual
	path := h.Path()
	// iteramos os modificadores de header
	for _, modifier := range backendRequest.ParamModifiers() {
		// se for propagate executamos
		if modifier.Propagate() {
			path = path.Modify(NewModify(&modifier, h, httpResponse))
		}
	}
	// retornamos o path possivelmente alterado
	return path
}

func (h *HttpRequest) modifyQuery(backendRequest *configVO.BackendRequest, httpResponse *HttpResponse) Query {
	// inicializamos o query com o valor atual
	query := h.Query()
	// iteramos os modificadores de header
	for _, modifier := range backendRequest.QueryModifiers() {
		// se for propagate executamos
		if modifier.Propagate() {
			query = query.Modify(NewModify(&modifier, h, httpResponse))
		}
	}
	// retornamos o query possivelmente alterado
	return query
}

func (h *HttpRequest) modifyBody(backendRequest *configVO.BackendRequest, httpResponse *HttpResponse) *Body {
	// inicializamos o query com o valor atual
	body := h.Body()
	// se o mesmo for nil retornamos nil
	if helper.IsNil(body) {
		return nil
	}
	// iteramos os modificadores de header
	for _, modifier := range backendRequest.BodyModifiers() {
		// se for propagate executamos
		if modifier.Propagate() {
			body = body.Modify(NewModify(&modifier, h, httpResponse))
		}
	}
	// retornamos o query possivelmente alterado
	return body
}

func (h *HttpRequest) CacheControl() enum.CacheControl {
	return enum.CacheControl(h.Header().Get("Cache-Control"))
}
