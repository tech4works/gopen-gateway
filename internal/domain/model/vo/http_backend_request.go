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
	"fmt"

	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type HTTPBackendRequest struct {
	degradation Degradation
	host        string
	path        URLPath
	method      string
	header      Metadata
	query       Query
	body        *Payload
}

func NewHTTPBackendRequest(
	degradation Degradation,
	host,
	method string,
	path URLPath,
	header Metadata,
	query Query,
	body *Payload,
) *HTTPBackendRequest {
	return &HTTPBackendRequest{
		degradation: degradation,
		host:        host,
		path:        path,
		method:      method,
		header:      header,
		query:       query,
		body:        body,
	}
}

func (b HTTPBackendRequest) Degradation() Degradation {
	return b.degradation
}

func (b *HTTPBackendRequest) Degraded() bool {
	return b.Degradation().Any()
}

func (b *HTTPBackendRequest) HeaderDegraded() bool {
	return b.Degradation().Has(enum.DegradationKindMetadata)
}

func (b *HTTPBackendRequest) URLPathDegraded() bool {
	return b.Degradation().Has(enum.DegradationKindURLPath)
}

func (b *HTTPBackendRequest) QueryDegraded() bool {
	return b.Degradation().Has(enum.DegradationKindQuery)
}

func (b *HTTPBackendRequest) BodyDegraded() bool {
	return b.Degradation().Has(enum.DegradationKindPayload)
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

func (b *HTTPBackendRequest) URL() string {
	return fmt.Sprint(b.host, b.path.String())
}

func (b *HTTPBackendRequest) Params() Params {
	return b.path.Params()
}

func (b *HTTPBackendRequest) Method() string {
	return b.method
}

func (b *HTTPBackendRequest) Header() Metadata {
	return b.header
}

func (b *HTTPBackendRequest) Query() Query {
	return b.query
}

func (b *HTTPBackendRequest) Body() *Payload {
	return b.body
}

func (b *HTTPBackendRequest) HasBody() bool {
	return checker.NonNil(b.body)
}
