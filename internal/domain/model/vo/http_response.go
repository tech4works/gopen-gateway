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

type HTTPResponse struct {
	statusCode StatusCode
	header     Header
	body       *Body
}

func NewHTTPResponse(statusCode StatusCode, header Header, body *Body) *HTTPResponse {
	return &HTTPResponse{
		statusCode: statusCode,
		header:     header,
		body:       body,
	}
}

func NewHTTPResponseStatusCode(statusCode StatusCode, header Header) *HTTPResponse {
	return &HTTPResponse{
		statusCode: statusCode,
		header:     header,
	}
}

func (h *HTTPResponse) StatusCode() StatusCode {
	return h.statusCode
}

func (h *HTTPResponse) Header() Header {
	return h.header
}

func (h *HTTPResponse) Body() *Body {
	return h.body
}

func (h *HTTPResponse) HasBody() bool {
	return checker.NonNil(h.body)
}
