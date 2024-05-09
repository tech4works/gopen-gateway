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

type Context struct {
	startTime    time.Time
	mutex        *sync.RWMutex
	framework    *gin.Context
	gopen        *vo.Gopen
	endpoint     *vo.Endpoint
	httpRequest  *vo.HttpRequest
	httpResponse *vo.HttpResponse
}

func (c *Context) Latency() time.Duration {
	return time.Now().Sub(c.startTime)
}

// Context returns the context of the Context. It delegates the call to the underlying framework's Context.Context() method.
func (c *Context) Context() context.Context {
	return c.framework.Request.Context()
}

func (c *Context) Done() <-chan struct{} {
	return c.Context().Done()
}

func (c *Context) Gopen() *vo.Gopen {
	return c.gopen
}

func (c *Context) Endpoint() *vo.Endpoint {
	return c.endpoint
}

func (c *Context) HttpRequest() *vo.HttpRequest {
	return c.httpRequest
}

func (c *Context) HttpResponse() *vo.HttpResponse {
	return c.httpResponse
}

func (c *Context) Http() *http.Request {
	return c.framework.Request
}

func (c *Context) RequestWithContext(ctx context.Context) {
	c.framework.Request = c.framework.Request.WithContext(ctx)
}

func (c *Context) Header() vo.Header {
	return c.httpRequest.Header()
}

func (c *Context) AddOnHeader(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.httpRequest = c.HttpRequest().SetHeader(key, value)
}

func (c *Context) SetOnHeader(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.httpRequest = c.HttpRequest().SetHeader(key, value)
}

func (c *Context) RemoteAddr() string {
	return c.framework.ClientIP()
}

func (c *Context) Method() string {
	return c.HttpRequest().Method()
}

func (c *Context) Url() string {
	return c.HttpRequest().Url()
}

func (c *Context) Body() *vo.Body {
	return c.HttpRequest().Body()
}

func (c *Context) BodyString() string {
	body := c.Body()
	if helper.IsNotNil(body) {
		return body.String()
	}
	return ""
}

func (c *Context) Params() vo.Params {
	return c.HttpRequest().Params()
}

func (c *Context) Query() vo.Query {
	return c.HttpRequest().Query()
}

func (c *Context) Next() {
	c.framework.Next()
}

func (c *Context) Write(httpRequest *vo.HttpRequest, httpResponse *vo.HttpResponse) {
	// inserimos o http request atualizado no contexto
	c.httpRequest = httpRequest
	// escrevemos o http response para o cliente final
	c.WriteHttpResponse(httpResponse)
}

func (c *Context) WriteHttpResponse(httpResponse *vo.HttpResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// se ja tiver abortado não fazemos nada
	if c.framework.IsAborted() {
		return
	}

	// escrevemos a resposta a partir da resposta
	httpResponseWritten := httpResponse.Write(c.Endpoint(), c.HttpRequest(), httpResponse)

	// escrevemos os headers de resposta
	c.writeHeader(httpResponseWritten.Header())

	// instanciamos os valores a serem utilizados
	statusCode := httpResponseWritten.StatusCode()
	contentType := httpResponseWritten.ContentType()
	bodyBytes := httpResponseWritten.BytesBody()

	// verificamos se tem valor o body
	if helper.IsNotEmpty(bodyBytes) {
		c.writeBody(statusCode, contentType.String(), bodyBytes)
	} else {
		c.writeStatusCode(statusCode)
	}

	// abortamos a requisição
	c.framework.Abort()

	// setamos a resposta VO escrita
	c.httpResponse = httpResponseWritten
}

func (c *Context) WriteStatusCode(code int) {
	// preparamos a resposta com o status code recebido
	httpResponseVO := vo.NewHttpResponseByStatusCode(code)
	// escrevemos a resposta
	c.WriteHttpResponse(httpResponseVO)
}

func (c *Context) WriteString(code int, body string) {
	// preparamos a resposta com a string
	httpResponseVO := vo.NewHttpResponseByString(code, body)
	// escrevemos a resposta
	c.WriteHttpResponse(httpResponseVO)
}

func (c *Context) WriteJson(code int, body any) {
	// preparamos a resposta com o body any
	httpResponseVO := vo.NewHttpResponseByJson(code, body)
	// escrevemos a resposta
	c.WriteHttpResponse(httpResponseVO)
}

func (c *Context) WriteCacheResponse(cacheResponse *vo.CacheResponse) {
	// preparamos a resposta a partir da resposta obtida do cache
	httpResponseVO := vo.NewHttpResponseByCache(cacheResponse)
	// escrevemos a resposta
	c.WriteHttpResponse(httpResponseVO)
}

func (c *Context) WriteError(code int, err error) {
	// preparamos a resposta a partir do status code e error
	httpResponseVO := vo.NewHttpResponseByErr(c.Endpoint().Path(), code, err)
	// escrevemos a resposta
	c.WriteHttpResponse(httpResponseVO)
}

func (c *Context) writeStatusCode(code int) {
	if c.framework.IsAborted() {
		return
	}
	c.framework.Status(code)
}

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

func (c *Context) writeBody(code int, contentType string, body []byte) {
	if c.framework.IsAborted() {
		return
	}
	c.framework.Data(code, contentType, body)
}
