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
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
)

type HTTPBackendRequest struct {
	host   string
	path   URLPath
	method string
	header Header
	query  Query
	body   *Body
}

func NewHTTPBackendRequest(
	host,
	method string,
	path URLPath,
	header Header,
	query Query,
	body *Body,
) *HTTPBackendRequest {
	return &HTTPBackendRequest{
		host:   host,
		path:   path,
		method: method,
		header: header,
		query:  query,
		body:   body,
	}
}

func (b *HTTPBackendRequest) Path() URLPath {
	return b.path
}

func (b *HTTPBackendRequest) FullPath() (r string) {
	r = b.path.String()
	if !b.Query().IsEmpty() {
		r += "?" + b.Query().Encode()
	}
	return
}

func (b *HTTPBackendRequest) Url() string {
	return fmt.Sprint(b.host, b.path.String())
}

func (b *HTTPBackendRequest) Params() Params {
	return b.path.Params()
}

func (b *HTTPBackendRequest) Method() string {
	return b.method
}

func (b *HTTPBackendRequest) Header() *Header {
	return &b.header
}

func (b *HTTPBackendRequest) Query() Query {
	return b.query
}

func (b *HTTPBackendRequest) Body() *Body {
	return b.body
}

func (b *HTTPBackendRequest) HasBody() bool {
	return helper.IsNotNil(b.body)
}
