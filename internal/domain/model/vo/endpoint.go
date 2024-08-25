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
	"time"
)

type Endpoint struct {
	path               string
	method             string
	timeout            Duration
	limiter            Limiter
	cache              *Cache
	abortIfStatusCodes *[]int
	response           *EndpointResponse
	backends           []Backend
}

type EndpointResponse struct {
	aggregate       bool
	contentType     enum.ContentType
	contentEncoding enum.ContentEncoding
	nomenclature    enum.Nomenclature
	omitEmpty       bool
}

func NewEndpoint(
	path,
	method string,
	timeout Duration,
	limiter Limiter,
	cache *Cache,
	abortIfStatusCodes *[]int,
	response *EndpointResponse,
	backends []Backend,
) Endpoint {
	return Endpoint{
		path:               path,
		method:             method,
		timeout:            timeout,
		limiter:            limiter,
		cache:              cache,
		abortIfStatusCodes: abortIfStatusCodes,
		response:           response,
		backends:           backends,
	}
}

func NewEndpointStatic(path, method string) Endpoint {
	return Endpoint{
		path:    path,
		method:  method,
		timeout: NewDuration(10 * time.Second),
		limiter: NewLimiterDefault(),
	}
}

func NewEndpointResponse(
	aggregate bool,
	contentType enum.ContentType,
	contentEncoding enum.ContentEncoding,
	nomenclature enum.Nomenclature,
	omitEmpty bool,
) *EndpointResponse {
	return &EndpointResponse{
		aggregate:       aggregate,
		contentType:     contentType,
		contentEncoding: contentEncoding,
		nomenclature:    nomenclature,
		omitEmpty:       omitEmpty,
	}
}

func (e *Endpoint) Path() string {
	return e.path
}

func (e *Endpoint) Method() string {
	return e.method
}

func (e *Endpoint) Timeout() Duration {
	if checker.IsGreaterThan(e.timeout, 0) {
		return e.timeout
	}
	return NewDuration(30 * time.Second)
}

func (e *Endpoint) Limiter() Limiter {
	return e.limiter
}

func (e *Endpoint) Cache() *Cache {
	return e.cache
}

func (e *Endpoint) Backends() []Backend {
	return e.backends
}

func (e *Endpoint) CountBeforewares() (count int) {
	for _, backend := range e.backends {
		if backend.IsBeforeware() {
			count++
		}
	}
	return count
}

func (e *Endpoint) CountAfterwares() (count int) {
	for _, backend := range e.backends {
		if backend.IsAfterware() {
			count++
		}
	}
	return count
}

func (e *Endpoint) CountBackends() (count int) {
	for _, backend := range e.backends {
		if backend.IsNormal() {
			count++
		}
	}
	return count
}

func (e *Endpoint) CountBackendsNonOmit() int {
	count := 0
	for _, backend := range e.Backends() {
		if checker.IsNil(backend.Response()) || !backend.Response().Omit() {
			count++
		}
	}
	return count
}

func (e *Endpoint) CountAllDataTransforms() (count int) {
	if checker.NonNil(e.Response()) {
		count += e.Response().CountAllDataTransforms()
	}
	for _, backend := range e.backends {
		count += backend.CountAllDataTransforms()
	}
	return count
}

func (e *Endpoint) HasAbortStatusCodes() bool {
	return checker.NonNil(e.abortIfStatusCodes)
}

func (e *Endpoint) Response() *EndpointResponse {
	return e.response
}

func (e *Endpoint) AbortIfStatusCodes() *[]int {
	return e.abortIfStatusCodes
}

func (e *Endpoint) Resume() string {
	return fmt.Sprintf("%s --> \"%s\" (beforeware:%v, afterware:%v, backends:%v, transformations:%v)",
		e.method, e.path, e.CountBeforewares(), e.CountAfterwares(), e.CountBackends(), e.CountAllDataTransforms())
}

func (e *Endpoint) NoCache() bool {
	return checker.IsNil(e.Cache()) || e.Cache().Disabled()
}

func (e *Endpoint) HasResponse() bool {
	return checker.NonNil(e.response)
}

func (e EndpointResponse) HasContentType() bool {
	return e.contentType.IsEnumValid()
}

func (e EndpointResponse) HasContentEncoding() bool {
	return e.contentEncoding.IsEnumValid()
}

func (e EndpointResponse) ContentType() enum.ContentType {
	return e.contentType
}

func (e EndpointResponse) ContentEncoding() enum.ContentEncoding {
	return e.contentEncoding
}

func (e EndpointResponse) Aggregate() bool {
	return e.aggregate
}

func (e EndpointResponse) OmitEmpty() bool {
	return e.omitEmpty
}

func (e EndpointResponse) HasNomenclature() bool {
	return e.nomenclature.IsEnumValid()
}

func (e EndpointResponse) Nomenclature() enum.Nomenclature {
	return e.nomenclature
}

func (e EndpointResponse) CountAllDataTransforms() (count int) {
	if e.Aggregate() {
		count++
	}
	if e.OmitEmpty() {
		count++
	}
	if e.HasContentType() {
		count++
	}
	if e.HasNomenclature() {
		count++
	}
	return count
}
