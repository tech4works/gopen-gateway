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

package api

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"sync"
)

// Context is a struct that represents the context of the current httpRequest.
// It contains various fields including mutex to synchronize access to the context,
// framework to handle the httpRequest, gopen configuration, endpoint information, and httpRequest and httpResponse data.
type Context struct {
	// mutex is a pointer to a sync.RWMutex structure which provides mutual
	// exclusion locking using read-write locks.
	mutex *sync.RWMutex
	// framework represents a context object for the Gin framework. It contains information about
	// the current HTTP httpRequest and httpResponse.
	framework *gin.Context
	// gopen represents a variable of type vo.Gopen. It is used to access and manipulate data using the desired
	// application settings
	gopen *vo.Gopen
	// endpoint represents the configuration of the endpoint that is receiving the current httpRequest, widely used to take
	// execution guidelines and httpResponse customization
	endpoint *vo.Endpoint
	// httpRequest represents a data structure for current `httpRequest`.
	httpRequest *vo.HttpRequest
	// httpResponse is a structure that represents the HTTP httpResponse, written by the context.
	httpResponse *vo.HttpResponse
}

// Context returns the context of the Context. It delegates the call to the underlying framework's Context.Context() method.
func (c *Context) Context() context.Context {
	return c.framework.Request.Context()
}

// Gopen returns the Gopen object associated with the Context. It retrieves the Gopen value from the Context object.
func (c *Context) Gopen() *vo.Gopen {
	return c.gopen
}

// Endpoint returns the endpoint associated with the httpRequest.
// It retrieves the endpoint value from the `endpoint` field of the Context struct.
func (c *Context) Endpoint() *vo.Endpoint {
	return c.endpoint
}

// HttpRequest returns the httpRequest object of the Context.
// It returns the `httpRequest` field of the Context struct.
func (c *Context) HttpRequest() *vo.HttpRequest {
	return c.httpRequest
}

// HttpResponse returns the httpResponse of the Context. It returns the httpResponse object stored
// in the Context struct.
func (c *Context) HttpResponse() *vo.HttpResponse {
	return c.httpResponse
}

// Http returns the underlying HTTP httpRequest object of the Context.
// It delegates the call to the underlying framework's HttpRequest property.
func (c *Context) Http() *http.Request {
	return c.framework.Request
}

// RequestWithContext sets the context of the Context to the provided context.
// It updates the underlying framework's Context.Context() method to use the new context.
func (c *Context) RequestWithContext(ctx context.Context) {
	c.framework.Request = c.framework.Request.WithContext(ctx)
}

// Header returns the `vo.Header` of the `HttpRequest`. It creates a new `vo.Header` using the underlying `http.Header`
// from the `HttpRequest`.
func (c *Context) Header() vo.Header {
	return c.httpRequest.Header()
}

// HeaderValue returns the value of the specified header key. It delegates the call to the underlying Context's
// Header().Get method.
func (c *Context) HeaderValue(key string) string {
	return c.Header().Get(key)
}

// AddHeader adds a new header to the HTTP httpRequest.
// It takes a key and value as parameters and adds them to the httpRequest's headers.
// Example usage:
//
//	req.AddHeader("Content-Type", "application/json")
//	req.AddHeader("Authorization", "Bearer token123")
//
// The method first creates a new header using the provided key and value.
// It then adds the header to the httpRequest using the Header() method of the context.
// Finally, it sets the updated httpRequest header using the SetHeader() method of the httpRequest.
func (c *Context) AddHeader(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	header := c.Header().Add(key, value)
	c.httpRequest = c.HttpRequest().SetHeader(header)
}

// SetHeader sets the value of the specified header key for the Context object.
// It delegates the call to the underlying framework's HttpRequest.Header.Set() method.
// Example usage: req.SetHeader("X-Forwarded-For", req.RemoteAddr()) and req.SetHeader("X-TraceId", t.traceProvider.GenerateTraceId())
// The SetHeader method takes a key and value as parameters, set the key value pair in the Context object's header.
// It uses the underlying framework's HttpRequest.Header.Set() method to update the header value.
// It returns nothing.
func (c *Context) SetHeader(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	header := c.Header().Set(key, value)
	c.httpRequest = c.HttpRequest().SetHeader(header)
}

// RemoteAddr returns the client's remote network address in the format "IP:port". It delegates the call to the
// underlying framework's ClientIP() method.
func (c *Context) RemoteAddr() string {
	return c.framework.ClientIP()
}

// Method returns the HTTP method of the Context.
// It delegates the call to the underlying framework's HttpRequest.Method() method.
func (c *Context) Method() string {
	return c.HttpRequest().Method()
}

// Url returns the URL of the httpRequest.
// It delegates the call to the underlying framework's HttpRequest.Url() method.
func (c *Context) Url() string {
	return c.HttpRequest().Url()
}

// Uri returns the URI of the Context. It delegates the call to the
// underlying HttpRequest's Uri() method.
func (c *Context) Uri() string {
	return c.HttpRequest().Path()
}

// Body returns the body of the Context. It delegates the call to the
// underlying framework's HttpRequest.Body() method.
func (c *Context) Body() *vo.Body {
	return c.HttpRequest().Body()
}

// BodyString returns the string representation of the body of the Context.
// It retrieves the body from the Context and converts it to a string using the String() method.
// If the body is nil, it returns an empty string.
func (c *Context) BodyString() string {
	body := c.Body()
	if helper.IsNotNil(body) {
		return body.String()
	}
	return ""
}

// Params returns the params of the Context.
// It delegates the call to the underlying HttpRequest's Params() method.
func (c *Context) Params() vo.Params {
	return c.HttpRequest().Params()
}

// Query returns the query object associated with the current context. It retrieves the query object
// by delegating the call to the underlying framework's HttpRequest().Query() method.
func (c *Context) Query() vo.Query {
	return c.HttpRequest().Query()
}

// Next calls the underlying framework's Next method to proceed to the next handler in the httpRequest chain.
func (c *Context) Next() {
	c.framework.Next()
}

func (c *Context) Write(httpResponseVO *vo.HttpResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// se ja tiver abortado não fazemos nada
	if c.framework.IsAborted() {
		return
	}

	// escrevemos a resposta a partir da resposta
	httpResponseWrittenVO := httpResponseVO.Write(c.Endpoint())

	// escrevemos os headers de resposta
	c.writeHeader(httpResponseWrittenVO.Header())

	// instanciamos os valores a serem utilizados
	statusCode := httpResponseWrittenVO.StatusCode()
	contentType := httpResponseWrittenVO.ContentType()
	bodyBytes := httpResponseWrittenVO.BytesBody()

	// verificamos se tem valor o body
	if helper.IsNotEmpty(bodyBytes) {
		c.writeBody(statusCode, contentType.String(), bodyBytes)
	} else {
		c.writeStatusCode(statusCode)
	}

	// abortamos a requisição
	c.framework.Abort()

	// setamos a resposta VO escrita
	c.httpResponse = httpResponseWrittenVO
}

func (c *Context) WriteStatusCode(code int) {
	// preparamos a resposta com o status code recebido
	httpResponseVO := vo.NewHttpResponseByStatusCode(code)
	// escrevemos a resposta
	c.Write(httpResponseVO)
}

func (c *Context) WriteString(code int, body string) {
	// preparamos a resposta com a string
	httpResponseVO := vo.NewHttpResponseByString(code, body)
	// escrevemos a resposta
	c.Write(httpResponseVO)
}

func (c *Context) WriteJson(code int, body any) {
	// preparamos a resposta com o body any
	httpResponseVO := vo.NewHttpResponseByJson(code, body)
	// escrevemos a resposta
	c.Write(httpResponseVO)
}

func (c *Context) WriteCacheResponse(cacheResponse *vo.CacheResponse) {
	// preparamos a resposta a partir da resposta obtida do cache
	httpResponseVO := vo.NewHttpResponseByCache(cacheResponse)
	// escrevemos a resposta
	c.Write(httpResponseVO)
}

func (c *Context) WriteError(code int, err error) {
	// preparamos a resposta a partir do status code e error
	httpResponseVO := vo.NewHttpResponseByErr(c.Endpoint().Path(), code, err)
	// escrevemos a resposta
	c.Write(httpResponseVO)
}

// writeStatusCode writes the HTTP status code to the httpResponse.
// If the httpRequest is already aborted, it does nothing.
// It sets the status code in the underlying framework using the given code.
// Parameter:
//   - code: the HTTP status code to be set in the httpResponse.
func (c *Context) writeStatusCode(code int) {
	if c.framework.IsAborted() {
		return
	}
	c.framework.Status(code)
}

// writeHeader writes the given header to the underlying framework's httpResponse. It skips certain headers including
// "Content-Length", "Content-Type", "Content-Encoding" (if it contains "gzip"), and "Date". The method delegates the
// call to the underlying framework's Header method for each non-skipped header.
func (c *Context) writeHeader(header vo.Header) {
	for key := range header {
		headerValue := header.Get(key)
		if helper.EqualsIgnoreCase(key, "Content-Length") || helper.EqualsIgnoreCase(key, "Content-Type") ||
			(helper.EqualsIgnoreCase(key, "Content-Encoding") && helper.ContainsIgnoreCase(headerValue, "gzip")) ||
			helper.EqualsIgnoreCase(key, "Date") {
			continue
		}
		c.framework.Header(key, header.Get(key))
	}
}

// writeBody writes the httpResponse body based on the configured httpResponse encoding of the endpoint.
// If the framework is aborted, the method returns early without writing the body.
// If the httpResponse encoding is set to ResponseEncodeText, the body is written as a string using the given code.
// If the httpResponse encoding is set to ResponseEncodeJson, the body is written as JSON using the given code.
// If the httpResponse encoding is set to ResponseEncodeXml, the body is written as XML using the given code.
// If the httpResponse encoding is set to ResponseEncodeYaml, the body is written as YAML using the given code.
// If none of the above cases match and the body is of JSON type, it is written as JSON using the given code.
// If none of the above cases match and the body is not of JSON type, it is written as a string using the given code.
func (c *Context) writeBody(code int, contentType string, body []byte) {
	if c.framework.IsAborted() {
		return
	}
	c.framework.Data(code, contentType, body)
}
