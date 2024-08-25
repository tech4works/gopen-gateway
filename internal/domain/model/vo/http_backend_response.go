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
)

type HTTPBackendResponse struct {
	statusCode StatusCode
	header     Header
	body       *Body
}

func NewHTTPBackendResponse(statusCode StatusCode, header Header, body *Body) *HTTPBackendResponse {
	return &HTTPBackendResponse{
		statusCode: statusCode,
		header:     header,
		body:       body,
	}
}

func (h *HTTPBackendResponse) OK() bool {
	return h.statusCode.OK()
}

func (h *HTTPBackendResponse) StatusCode() StatusCode {
	return h.statusCode
}

func (h *HTTPBackendResponse) Header() Header {
	return h.header
}

func (h *HTTPBackendResponse) Body() *Body {
	return h.body
}

func (h *HTTPBackendResponse) Map() (any, error) {
	var body any
	if checker.NonNil(h.body) {
		bodyMap, err := h.body.Map()
		if checker.NonNil(err) {
			return nil, err
		}
		body = bodyMap
	}
	return map[string]any{
		"statusCode": h.statusCode,
		"header":     h.header,
		"body":       body,
	}, nil
}

func (h *HTTPBackendResponse) HasBody() bool {
	return checker.NonNil(h.body)
}
