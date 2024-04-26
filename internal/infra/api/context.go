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

// Context is a struct that represents the context of the current request.
// It contains various fields including mutex to synchronize access to the context,
// framework to handle the request, gopen configuration, endpoint information, and request and response data.
type Context struct {
	// mutex is a pointer to a sync.RWMutex structure which provides mutual
	// exclusion locking using read-write locks.
	mutex *sync.RWMutex
	// framework represents a context object for the Gin framework. It contains information about
	// the current HTTP request and response.
	framework *gin.Context
	// gopen represents a variable of type vo.Gopen. It is used to access and manipulate data using the desired
	// application settings
	gopen *vo.Gopen
	// endpoint represents the configuration of the endpoint that is receiving the current request, widely used to take
	// execution guidelines and response customization
	endpoint *vo.Endpoint
	// request represents a data structure for current `request`.
	request *vo.Request
	// response is a structure that represents the HTTP response, written by the context.
	response *vo.Response
}

// Context returns the context of the Context. It delegates the call to the underlying framework's Context.Context() method.
func (c *Context) Context() context.Context {
	return c.framework.Request.Context()
}

// Gopen returns the Gopen object associated with the Context. It retrieves the Gopen value from the Context object.
func (c *Context) Gopen() *vo.Gopen {
	return c.gopen
}

// Endpoint returns the endpoint associated with the request.
// It retrieves the endpoint value from the `endpoint` field of the Context struct.
func (c *Context) Endpoint() *vo.Endpoint {
	return c.endpoint
}

// Request returns the request object of the Context.
// It returns the `request` field of the Context struct.
func (c *Context) Request() *vo.Request {
	return c.request
}

// Response returns the response of the Context. It returns the response object stored
// in the Context struct.
func (c *Context) Response() *vo.Response {
	return c.response
}

// Http returns the underlying HTTP request object of the Context.
// It delegates the call to the underlying framework's Request property.
func (c *Context) Http() *http.Request {
	return c.framework.Request
}

// RequestWithContext sets the context of the Context to the provided context.
// It updates the underlying framework's Context.Context() method to use the new context.
func (c *Context) RequestWithContext(ctx context.Context) {
	c.framework.Request = c.framework.Request.WithContext(ctx)
}

// Header returns the `vo.Header` of the `Request`. It creates a new `vo.Header` using the underlying `http.Header`
// from the `Request`.
func (c *Context) Header() vo.Header {
	return c.request.Header()
}

// HeaderValue returns the value of the specified header key. It delegates the call to the underlying Context's
// Header().Get method.
func (c *Context) HeaderValue(key string) string {
	return c.Header().Get(key)
}

// AddHeader adds a new header to the HTTP request.
// It takes a key and value as parameters and adds them to the request's headers.
// Example usage:
//
//	req.AddHeader("Content-Type", "application/json")
//	req.AddHeader("Authorization", "Bearer token123")
//
// The method first creates a new header using the provided key and value.
// It then adds the header to the request using the Header() method of the context.
// Finally, it sets the updated request header using the SetHeader() method of the request.
func (c *Context) AddHeader(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	header := c.Header().Add(key, value)
	c.request = c.Request().SetHeader(header)
}

// SetHeader sets the value of the specified header key for the Context object.
// It delegates the call to the underlying framework's Request.Header.Set() method.
// Example usage: req.SetHeader("X-Forwarded-For", req.RemoteAddr()) and req.SetHeader("X-TraceId", t.traceProvider.GenerateTraceId())
// The SetHeader method takes a key and value as parameters, set the key value pair in the Context object's header.
// It uses the underlying framework's Request.Header.Set() method to update the header value.
// It returns nothing.
func (c *Context) SetHeader(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	header := c.Header().Set(key, value)
	c.request = c.Request().SetHeader(header)
}

// RemoteAddr returns the client's remote network address in the format "IP:port". It delegates the call to the
// underlying framework's ClientIP() method.
func (c *Context) RemoteAddr() string {
	return c.framework.ClientIP()
}

// Method returns the HTTP method of the Context.
// It delegates the call to the underlying framework's Request.Method() method.
func (c *Context) Method() string {
	return c.Request().Method()
}

// Url returns the URL of the request.
// It delegates the call to the underlying framework's Request.Url() method.
func (c *Context) Url() string {
	return c.Request().Url()
}

// Uri returns the URI of the Context. It delegates the call to the
// underlying Request's Uri() method.
func (c *Context) Uri() string {
	return c.Request().Path()
}

// Body returns the body of the Context. It delegates the call to the
// underlying framework's Request.Body() method.
func (c *Context) Body() *vo.Body {
	return c.Request().Body()
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
// It delegates the call to the underlying Request's Params() method.
func (c *Context) Params() vo.Params {
	return c.Request().Params()
}

// Query returns the query object associated with the current context. It retrieves the query object
// by delegating the call to the underlying framework's Request().Query() method.
func (c *Context) Query() vo.Query {
	return c.Request().Query()
}

// Next calls the underlying framework's Next method to proceed to the next handler in the request chain.
func (c *Context) Next() {
	c.framework.Next()
}

// Write writes the response to the client.
// It first checks if the request has already been aborted, in which case it does nothing.
// Then, it writes the response headers.
// It retrieves the status code and body from the responseVO.
// If the body is not empty, it writes the body along with the status code.
// Otherwise, it only writes the status code.
func (c *Context) Write(responseVO *vo.Response) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// se ja tiver abortado não fazemos nada
	if c.framework.IsAborted() {
		return
	}

	// escrevemos os headers de resposta
	c.writeHeader(responseVO.Header())

	// instanciamos os valores a serem utilizados
	statusCode := responseVO.StatusCode()
	contentType := responseVO.ContentType()
	bodyBytes := responseVO.BytesBody()

	// verificamos se tem valor o body
	if helper.IsNotEmpty(bodyBytes) {
		c.writeBody(statusCode, contentType.String(), bodyBytes)
	} else {
		c.writeStatusCode(statusCode)
	}

	// abortamos a requisição
	c.framework.Abort()

	// setamos a resposta VO escrita
	c.response = responseVO
}

// WriteCacheResponse writes the cache response to the client's response.
// It creates a new response using the cache response and writes it.
func (c *Context) WriteCacheResponse(cacheResponse *vo.CacheResponse) {
	// preparamos a resposta
	responseVO := vo.NewResponseByCache(c.endpoint, cacheResponse)
	// escrevemos a resposta
	c.Write(responseVO)
}

// WriteError writes an error response to the client.
// It creates a new Response object with the provided code and error, and delegates the writing of the response to the
// Write method.
func (c *Context) WriteError(code int, err error) {
	// preparamos a resposta
	responseVO := vo.NewResponseByErr(c.Endpoint(), code, err)
	// escrevemos a resposta
	c.Write(responseVO)
}

// writeHeader sets the headers in the HTTP response as received in the `header` argument, excluding certain headers.
// Headers to be ignored: "Content-Length", "Content-Type", "Date".
// The method delegates the actual header setting to the underlying framework's `Header()` method.
// Example usage can be found in the `Write()` method.
func (c *Context) writeHeader(header vo.Header) {
	for key := range header {
		if helper.EqualsIgnoreCase(key, "Content-Length") || helper.EqualsIgnoreCase(key, "Content-Type") ||
			helper.EqualsIgnoreCase(key, "Date") {
			continue
		}
		c.framework.Header(key, header.Get(key))
	}
}

// writeBody writes the response body based on the configured response encoding of the endpoint.
// If the framework is aborted, the method returns early without writing the body.
// If the response encoding is set to ResponseEncodeText, the body is written as a string using the given code.
// If the response encoding is set to ResponseEncodeJson, the body is written as JSON using the given code.
// If the response encoding is set to ResponseEncodeXml, the body is written as XML using the given code.
// If the response encoding is set to ResponseEncodeYaml, the body is written as YAML using the given code.
// If none of the above cases match and the body is of JSON type, it is written as JSON using the given code.
// If none of the above cases match and the body is not of JSON type, it is written as a string using the given code.
func (c *Context) writeBody(code int, contentType string, body []byte) {
	if c.framework.IsAborted() {
		return
	}
	c.framework.Data(code, contentType, body)
}

// writeStatusCode writes the HTTP status code to the response.
// If the request is already aborted, it does nothing.
// It sets the status code in the underlying framework using the given code.
// Parameter:
//   - code: the HTTP status code to be set in the response.
func (c *Context) writeStatusCode(code int) {
	if c.framework.IsAborted() {
		return
	}
	c.framework.Status(code)
}
