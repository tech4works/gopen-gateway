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

package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
)

type HTTPRequest struct {
	url    string
	path   URLPath
	method string
	header Header
	query  Query
	body   *Body
}

func NewHTTPRequest(path URLPath, url, method string, header Header, query Query, body *Body) *HTTPRequest {
	return &HTTPRequest{
		path:   path,
		url:    url,
		method: method,
		header: header,
		query:  query,
		body:   body,
	}
}

func (h *HTTPRequest) Url() string {
	return h.url
}

func (h *HTTPRequest) Path() URLPath {
	return h.path
}

func (h *HTTPRequest) Method() string {
	return h.method
}

func (h *HTTPRequest) Header() Header {
	return h.header
}

func (h *HTTPRequest) Params() Params {
	return h.Path().Params()
}

func (h *HTTPRequest) Query() Query {
	return h.query
}

func (h *HTTPRequest) Body() *Body {
	return h.body
}

func (h *HTTPRequest) Map() (string, error) {
	var body any
	if checker.NonNil(h.Body()) {
		bodyMap, err := h.Body().Map()
		if checker.NonNil(err) {
			return "", err
		}
		body = bodyMap
	}
	return converter.ToStringWithErr(map[string]any{
		"header": h.Header().Map(),
		"params": h.Params().Map(),
		"query":  h.Query().Map(),
		"body":   body,
	})
}

func (h *HTTPRequest) ClientIP() string {
	return h.Header().GetFirst(mapper.XForwardedFor)
}

func (h *HTTPRequest) HasBody() bool {
	return checker.NonNil(h.body)
}
