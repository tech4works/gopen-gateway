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
	"time"
)

// Context represents the context of a request being handled by the server.
// It contains various objects related to the request and response.
type Context struct {
	// startTime represents the start time of the request processing in the Context struct.
	startTime time.Time
	// mutex is a pointer to a sync.RWMutex, which implements the sync.Locker interface.
	// It is used for controlling concurrent access to shared resources in the Context struct.
	mutex *sync.RWMutex
	// framework represents the context of a request and response in the Gin framework.
	// It contains methods for handling the request and response objects, as well as
	// accessing various properties and performing operations related to the request and response.
	framework *gin.Context
	// Gopen represents the configuration for a Gopen server.
	// It contains the port number for listening to incoming requests, the CORS settings, middleware configuration,
	// and a slice of endpoints.
	// The Gopen struct has methods to retrieve the port, CORS settings, middlewares, and endpoints.
	// It also has methods to count the number of endpoints, middlewares, backends, and data transforms.
	gopen *vo.Gopen
	// endpoint represents the endpoint of a request being handled by the server.
	// It contains information about the path, method, and other details of the request.
	// The endpoint is stored in the Context struct and can be accessed using the Endpoint() method.
	// It is used for handling the request and response objects.
	// The Endpoint struct has methods to retrieve the path, method, and other properties of the endpoint.
	// It also has methods to set the header, body, and status code of the response.
	// The endpoint is an important part of the request processing flow in the server.
	endpoint *vo.Endpoint
	// httpRequest is a pointer to the vo.HttpRequest object.
	// It represents the HTTP request received by the server.
	// It contains information about the request method, URL, headers, body, and other related details.
	httpRequest *vo.HttpRequest
	// httpResponse is a value object representing the HTTP response.
	// It contains the status code, content type, and response body.
	// It is used to write the HTTP response to the client.
	httpResponse *vo.HttpResponse
}

// Latency returns the duration between the current time and the start time of the Context.
func (c *Context) Latency() time.Duration {
	return time.Now().Sub(c.startTime)
}

// Context returns the context of the Context. It delegates the call to the underlying framework's Context.Context() method.
func (c *Context) Context() context.Context {
	return c.framework.Request.Context()
}

// Done returns a channel `<-chan struct{}` from the underlying framework's `Context().Done()` method.
func (c *Context) Done() <-chan struct{} {
	return c.Context().Done()
}

// Gopen returns the Gopen object associated with the Context. It represents the configuration
// for a Gopen server. It contains various fields such as port, securityCors, middlewares, and
// endpoints, which define the behavior and settings of the Gopen server.
func (c *Context) Gopen() *vo.Gopen {
	return c.gopen
}

// Endpoint returns the endpoint associated with the Context. It represents the configuration
// for a specific API endpoint. It contains fields such as method, path, parameters, and
// response types, which define the behavior and settings of the endpoint.
func (c *Context) Endpoint() *vo.Endpoint {
	return c.endpoint
}

// HttpRequest returns the HttpRequest object associated with the Context.
// It represents the HTTP request received by the server, containing information
// such as the request method, URL, headers, and body.
func (c *Context) HttpRequest() *vo.HttpRequest {
	return c.httpRequest
}

// HttpResponse returns the HttpResponse object associated with the Context.
// It represents the HTTP response that will be sent back to the client,
// containing information such as the response status code, headers, and body.
func (c *Context) HttpResponse() *vo.HttpResponse {
	return c.httpResponse
}

// Http returns the *http.Request object associated with the Context.
// It represents the HTTP request received by the server, containing information
// such as the request method, URL, headers, and body.
func (c *Context) Http() *http.Request {
	return c.framework.Request
}

// WithContext sets the context of the Context by delegating the call to the underlying framework's
// Request.WithContext method. This allows the propagation of the context to other handlers.
func (c *Context) WithContext(ctx context.Context) {
	c.framework.Request = c.framework.Request.WithContext(ctx)
}

// Header returns the header object associated with the Context.
// It represents the header of the HTTP request received by the server,
// containing information such as the request method, URL, headers, and body.
func (c *Context) Header() vo.Header {
	return c.httpRequest.Header()
}

// AddOnHeader adds a header to the HttpRequest object associated with the Context.
// It acquires a mutex lock to ensure thread safety, sets the header key-value pair,
// and updates the HttpRequest object. It is used to modify the headers of the HTTP request.
// The changes made by this method will affect the subsequent processing of the request.
func (c *Context) AddOnHeader(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.httpRequest = c.HttpRequest().AddOnHeader(key, value)
}

// SetOnHeader sets the value of the specified header key in the HttpRequest object associated with the Context.
// It acquires a mutex lock to ensure thread safety, sets the header key-value pair,
// and updates the HttpRequest object. It is used to modify a header value in the HTTP request.
// The changes made by this method will affect the subsequent processing of the request.
func (c *Context) SetOnHeader(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.httpRequest = c.HttpRequest().SetOnHeader(key, value)
}

// RemoteAddr returns the remote network address of the client making the request.
func (c *Context) RemoteAddr() string {
	return c.framework.ClientIP()
}

// Method returns the HTTP method of the HttpRequest object associated with the Context.
func (c *Context) Method() string {
	return c.HttpRequest().Method()
}

// Url returns the URL from the HttpRequest object associated with the Context.
func (c *Context) Url() string {
	return c.HttpRequest().Url()
}

// Body returns the Body object associated with the Context.
// It represents the body of the HTTP request received by the server,
// containing information such as the request body content.
func (c *Context) Body() *vo.Body {
	return c.HttpRequest().Body()
}

// BodyString returns the body content of the HTTP request as a string.
// It retrieves the body object associated with the context. If the body object is not nil,
// it returns the string representation of the body. Otherwise, it returns an empty string.
func (c *Context) BodyString() string {
	body := c.Body()
	if helper.IsNotNil(body) {
		return body.String()
	}
	return ""
}

// Params returns the Params object associated with the Context.
// It represents the parameters extracted from the URL of the HTTP request,
// typically used in routing to match specific endpoints or retrieve query parameters.
func (c *Context) Params() vo.Params {
	return c.HttpRequest().Params()
}

// Query returns the Query object associated with the Context.
// It represents the query parameters extracted from the URL of the HTTP request.
func (c *Context) Query() vo.Query {
	return c.HttpRequest().Query()
}

// Next calls the Next method of the underlying framework's Context object.
func (c *Context) Next() {
	c.framework.Next()
}

// Write updates the current HTTP request in the context and writes the HTTP response to the client.
// It takes an updated httpRequest of type *vo.HttpRequest and an httpResponse of type *vo.HttpResponse as input parameters.
// It first updates the httpRequest in the context by assigning it to c.httpRequest.
// Then, it calls the WriteHttpResponse method to write the httpResponse to the client.
// Note: This method assumes that the httpRequest and httpResponse are already populated with the necessary data.
// If not, the method may not function as expected.
func (c *Context) Write(httpRequest *vo.HttpRequest, httpResponse *vo.HttpResponse) {
	c.httpRequest = httpRequest
	c.WriteHttpResponse(httpResponse)
}

// WriteHttpResponse writes the provided http response into the context's response writer.
// It acquires a lock to ensure exclusive access to the context, checks if the request has been aborted, and if so, returns.
// It then writes the response headers and obtains the status code, content type, and body bytes.
// If the body is not empty, it writes the body along with the status code and content type.
// Otherwise, it writes only the status code.
// After writing the response, it aborts the request and sets the written http response in the context.
func (c *Context) WriteHttpResponse(httpResponse *vo.HttpResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.framework.IsAborted() {
		return
	}

	httpResponseWritten := httpResponse.Write(c.Endpoint(), c.HttpRequest(), httpResponse)
	statusCode := httpResponseWritten.StatusCode()
	contentType := httpResponseWritten.ContentType()
	bodyBytes := httpResponseWritten.BodyBytes()

	c.writeHeader(httpResponseWritten.Header())
	if helper.IsNotEmpty(bodyBytes) {
		c.writeBody(statusCode, contentType.String(), bodyBytes)
	} else {
		c.writeStatusCode(statusCode)
	}
	c.framework.Abort()

	c.httpResponse = httpResponseWritten
}

// WriteStatusCode writes the given status code to the HTTP response.
// It creates a new HttpResponse object based on the status code,
// and then calls WriteHttpResponse to write the response to the client.
func (c *Context) WriteStatusCode(code int) {
	httpResponse := vo.NewHttpResponseByStatusCode(code)
	c.WriteHttpResponse(httpResponse)
}

// WriteString takes an HTTP status code and a string body as input and creates
// an HTTP response using the given code and body. It then writes this response
// to the underlying HTTP writer by calling WriteHttpResponse.
func (c *Context) WriteString(code int, body string) {
	httpResponse := vo.NewHttpResponseByString(code, body)
	c.WriteHttpResponse(httpResponse)
}

// WriteJson writes a JSON response based on the given status code and body.
// It creates a new HttpResponse using the provided code and body, then calls
// WriteHttpResponse with the created HttpResponse object.
func (c *Context) WriteJson(code int, body any) {
	httpResponse := vo.NewHttpResponseByJson(code, body)
	c.WriteHttpResponse(httpResponse)
}

// WriteCacheResponse writes the cache response to the context's HTTP response.
// It creates a new HTTP response using the cache response and calls WriteHttpResponse to write the response.
func (c *Context) WriteCacheResponse(cacheResponse *vo.CacheResponse) {
	httpResponse := vo.NewHttpResponseByCache(cacheResponse)
	c.WriteHttpResponse(httpResponse)
}

// WriteError writes an HTTP response with the given status code and error message.
// It creates a new HTTP response using the given error and endpoint path, and
// then calls WriteHttpResponse to send the response back to the client.
func (c *Context) WriteError(code int, err error) {
	httpResponse := vo.NewHttpResponseByErr(c.Endpoint().Path(), code, err)
	c.WriteHttpResponse(httpResponse)
}

// writeStatusCode writes the specified status code to the framework's response.
// If the framework is already aborted, the function returns immediately without modifying the response.
func (c *Context) writeStatusCode(code int) {
	if c.framework.IsAborted() {
		return
	}
	c.framework.Status(code)
}

// writeHeader writes the specified header to the response. It skips certain headers such as
// "Content-Length", "Content-Type", "Content-Encoding" (if it contains "gzip"), and "Date".
// Other headers are passed to the underlying framework to be written to the response.
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

// writeBody writes the given body with the specified status code and content type to the response.
// It first checks if the framework's IsAborted method returns true. If true, it returns immediately
// without writing the response. Otherwise, it calls the framework's Data method to set the response
// code, content type, and body.
func (c *Context) writeBody(code int, contentType string, body []byte) {
	if c.framework.IsAborted() {
		return
	}
	c.framework.Data(code, contentType, body)
}
