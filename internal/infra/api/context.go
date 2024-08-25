/*
 * Copyright 2024 Tech4Works
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
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"go.elastic.co/apm/v2"
	"golang.org/x/net/context"
	"io"
	"sync"
	"time"
)

type Context struct {
	startTime time.Time
	mutex     *sync.RWMutex
	engine    *gin.Context
	gopen     *vo.Gopen
	endpoint  *vo.Endpoint
	request   *vo.HTTPRequest
	response  *vo.HTTPResponse
}

func newContext(gin *gin.Context, gopen *vo.Gopen, endpoint *vo.Endpoint) app.Context {
	request := buildHTTPRequest(gin)
	return &Context{
		startTime: time.Now(),
		mutex:     &sync.RWMutex{},
		engine:    gin,
		gopen:     gopen,
		endpoint:  endpoint,
		request:   request,
	}
}

func buildHTTPRequest(gin *gin.Context) *vo.HTTPRequest {
	gin.Request.Header.Add(mapper.XForwardedFor, gin.ClientIP())
	header := vo.NewHeader(gin.Request.Header)

	query := vo.NewQuery(gin.Request.URL.Query())
	url := gin.Request.URL.Path
	if !query.IsEmpty() {
		url = fmt.Sprint(url, "?", query.Encode())
	}

	ginParams := map[string]string{}
	for _, param := range gin.Params {
		ginParams[param.Key] = param.Value
	}
	path := vo.NewURLPath(gin.FullPath(), ginParams)

	bodyBytes, _ := io.ReadAll(gin.Request.Body)
	body := vo.NewBody(gin.GetHeader(mapper.ContentType), gin.GetHeader(mapper.ContentEncoding), bytes.NewBuffer(bodyBytes))

	return vo.NewHTTPRequest(path, url, gin.Request.Method, header, query, body)
}

func (c *Context) Context() context.Context {
	return c.engine.Request.Context()
}

func (c *Context) Done() <-chan struct{} {
	return c.Context().Done()
}

func (c *Context) WithContext(ctx context.Context) {
	c.engine.Request = c.engine.Request.WithContext(ctx)
}

func (c *Context) Next() {
	c.engine.Next()
}

func (c *Context) Duration() time.Duration {
	return time.Now().Sub(c.startTime)
}

func (c *Context) TraceID() string {
	tx := apm.TransactionFromContext(c.Context())
	if checker.NonNil(tx) {
		return tx.TraceContext().Trace.String()
	}
	return "undefined"
}

func (c *Context) ClientIP() string {
	return c.Request().Header().GetFirst(mapper.XForwardedFor)
}

func (c *Context) Gopen() *vo.Gopen {
	return c.gopen
}

func (c *Context) Endpoint() *vo.Endpoint {
	return c.endpoint
}

func (c *Context) Request() *vo.HTTPRequest {
	return c.request
}

func (c *Context) Response() *vo.HTTPResponse {
	return c.response
}

func (c *Context) Write(response *vo.HTTPResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.engine.IsAborted() {
		return
	}

	var contentType vo.ContentType
	var rawBodyBytes []byte
	if response.HasBody() {
		contentType = response.Body().ContentType()
		rawBodyBytes = response.Body().RawBytes()
	}

	c.writeHeader(response.Header())
	if checker.IsNotEmpty(rawBodyBytes) {
		c.writeBody(response.StatusCode(), contentType.String(), rawBodyBytes)
	} else {
		c.writeStatusCode(response.StatusCode())
	}

	c.engine.Abort()
	c.response = response
}

func (c *Context) WriteError(code int, err error) {
	statusCode := vo.NewStatusCode(code)

	details := errors.Details(err)
	buffer := converter.ToBuffer(dto.ErrorBody{
		File:      details.File(),
		Line:      details.Line(),
		Endpoint:  c.endpoint.Path(),
		Message:   details.Message(),
		Timestamp: time.Now(),
	})
	body := vo.NewBodyJson(buffer)
	header := c.buildHeader(false, statusCode, body)

	c.Write(vo.NewHTTPResponse(statusCode, header, body))
}

func (c *Context) WriteCacheResponse(cacheResponse *vo.CacheResponse) {
	c.Write(vo.NewHTTPResponse(cacheResponse.StatusCode, c.buildCacheHeader(cacheResponse), cacheResponse.Body))
}

func (c *Context) WriteStatusCode(code int) {
	statusCode := vo.NewStatusCode(code)
	header := c.buildHeader(true, statusCode, nil)

	c.Write(vo.NewHTTPResponseStatusCode(statusCode, header))
}

func (c *Context) WriteString(code int, s string) {
	statusCode := vo.NewStatusCode(code)
	body := vo.NewBodyWithContentType(vo.NewContentTypeTextPlain(), converter.ToBuffer(s))
	header := c.buildHeader(true, statusCode, body)

	c.Write(vo.NewHTTPResponse(statusCode, header, body))
}

func (c *Context) WriteJson(code int, a any) {
	statusCode := vo.NewStatusCode(code)
	body := vo.NewBodyWithContentType(vo.NewContentTypeJson(), converter.ToBuffer(a))
	header := c.buildHeader(true, statusCode, body)

	c.Write(vo.NewHTTPResponse(statusCode, header, body))
}

func (c *Context) buildHeader(complete bool, statusCode vo.StatusCode, body *vo.Body) vo.Header {
	mapHeader := map[string][]string{
		mapper.XGopenCache:    {"false"},
		mapper.XGopenComplete: {converter.ToString(complete)},
		mapper.XGopenSuccess:  {converter.ToString(statusCode.OK())},
	}
	if checker.NonNil(body) {
		mapHeader[mapper.ContentType] = []string{body.ContentType().String()}
		mapHeader[mapper.ContentLength] = []string{body.SizeInString()}
	}
	return vo.NewHeader(mapHeader)
}

func (c *Context) buildCacheHeader(cacheResponse *vo.CacheResponse) vo.Header {
	copied := cacheResponse.Header.Copy()
	copied[mapper.XGopenCache] = []string{"true"}
	copied[mapper.XGopenCacheTTL] = []string{cacheResponse.TTL()}
	return vo.NewHeader(copied)
}

func (c *Context) writeStatusCode(statusCode vo.StatusCode) {
	if c.engine.IsAborted() {
		return
	}
	c.engine.Status(statusCode.Code())
}

func (c *Context) writeHeader(header vo.Header) {
	for _, key := range header.Keys() {
		c.engine.Header(key, header.Get(key))
	}
}

func (c *Context) writeBody(statusCode vo.StatusCode, contentType string, body []byte) {
	if c.engine.IsAborted() {
		return
	}
	c.engine.Data(statusCode.Code(), contentType, body)
}
