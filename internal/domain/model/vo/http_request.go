/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a Copy of the License at
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
	"github.com/gin-gonic/gin"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"io"
)

type HTTPRequest struct {
	url    string
	path   URLPath
	method string
	header Header
	query  Query
	body   *Body
}

func NewHTTPRequest(gin *gin.Context) *HTTPRequest {
	gin.Request.Header.Add(mapper.XForwardedFor, gin.ClientIP())
	header := NewHeader(gin.Request.Header)

	query := NewQuery(gin.Request.URL.Query())
	url := gin.Request.URL.Path
	if helper.IsNotEmpty(query) {
		url = fmt.Sprint(url, "?", query.Encode())
	}

	ginParams := map[string]string{}
	for _, param := range gin.Params {
		ginParams[param.Key] = param.Value
	}
	path := NewURLPath(gin.FullPath(), ginParams)

	bodyBytes, _ := io.ReadAll(gin.Request.Body)
	body := NewBody(gin.GetHeader(mapper.ContentType), gin.GetHeader(mapper.ContentEncoding), bytes.NewBuffer(bodyBytes))

	return &HTTPRequest{
		path:   path,
		url:    url,
		method: gin.Request.Method,
		header: header,
		query:  query,
		body:   body,
	}
}

func (r *HTTPRequest) Url() string {
	return r.url
}

func (r *HTTPRequest) Path() URLPath {
	return r.path
}

func (r *HTTPRequest) Method() string {
	return r.method
}

func (r *HTTPRequest) Header() *Header {
	return &r.header
}

func (r *HTTPRequest) Params() Params {
	return r.Path().Params()
}

func (r *HTTPRequest) Query() Query {
	return r.query
}

func (r *HTTPRequest) Body() *Body {
	return r.body
}

func (r *HTTPRequest) Map() (string, error) {
	var body any
	if helper.IsNotNil(r.Body()) {
		bodyMap, err := r.Body().Map()
		if helper.IsNotNil(err) {
			return "", err
		}
		body = bodyMap
	}
	return helper.ConvertToString(map[string]any{
		"header": r.Header().Map(),
		"params": r.Params().Map(),
		"query":  r.Query().Map(),
		"body":   body,
	})
}

func (r *HTTPRequest) ClientIP() string {
	return r.Header().GetFirst(mapper.XForwardedFor)
}
