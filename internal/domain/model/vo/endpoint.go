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
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type Endpoint struct {
	path                   string
	method                 string
	timeout                Duration
	limiter                Limiter
	cache                  *Cache
	abortIfHTTPStatusCodes *[]int
	parallelism            bool
	response               *EndpointResponse
	backends               []Backend
}

type EndpointResponse struct {
	continueOnError bool
	header          *EndpointResponseHeader
	body            *EndpointResponseBody
}

type EndpointResponseHeader struct {
	mapper    *Mapper
	projector *Projector
}

type EndpointResponseBody struct {
	aggregate       bool
	omitEmpty       bool
	contentType     enum.ContentType
	contentEncoding enum.ContentEncoding
	nomenclature    enum.Nomenclature
	mapper          *Mapper
	projector       *Projector
}

func NewEndpoint(
	path,
	method string,
	timeout Duration,
	limiter Limiter,
	cache *Cache,
	abortIfHTTPStatusCodes *[]int,
	parallelism bool,
	response *EndpointResponse,
	backends []Backend,
) Endpoint {
	return Endpoint{
		path:                   path,
		method:                 method,
		timeout:                timeout,
		limiter:                limiter,
		cache:                  cache,
		abortIfHTTPStatusCodes: abortIfHTTPStatusCodes,
		parallelism:            parallelism,
		response:               response,
		backends:               backends,
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
	continueOnError bool,
	header *EndpointResponseHeader,
	body *EndpointResponseBody,
) *EndpointResponse {
	return &EndpointResponse{
		continueOnError: continueOnError,
		header:          header,
		body:            body,
	}
}

func NewEndpointResponseHeader(mapper *Mapper, projector *Projector) *EndpointResponseHeader {
	return &EndpointResponseHeader{
		mapper:    mapper,
		projector: projector,
	}
}

func NewEndpointResponseBody(
	aggregate bool,
	omitEmpty bool,
	contentType enum.ContentType,
	contentEncoding enum.ContentEncoding,
	nomenclature enum.Nomenclature,
	mapper *Mapper,
	projector *Projector,
) *EndpointResponseBody {
	return &EndpointResponseBody{
		aggregate:       aggregate,
		omitEmpty:       omitEmpty,
		contentType:     contentType,
		contentEncoding: contentEncoding,
		nomenclature:    nomenclature,
		mapper:          mapper,
		projector:       projector,
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

func (e *Endpoint) CountAllDataTransforms() (count int) {
	if e.HasResponse() && e.Response().HasBody() {
		count += e.Response().Body().CountAllDataTransforms()
	}
	for _, backend := range e.backends {
		count += backend.CountAllDataTransforms()
	}
	return count
}

func (e *Endpoint) HasAbortIfHTTPStatusCodes() bool {
	return checker.NonNil(e.abortIfHTTPStatusCodes)
}

func (e *Endpoint) Response() *EndpointResponse {
	return e.response
}

func (e *Endpoint) Parallelism() bool {
	return e.parallelism
}

func (e *Endpoint) AbortIfHTTPStatusCodes() *[]int {
	return e.abortIfHTTPStatusCodes
}

func (e *Endpoint) Resume() string {
	return fmt.Sprintf("%s --> \"%s\" (beforeware:%v, backends:%v, afterware:%v, transformations:%v)",
		e.method, e.path, e.CountBeforewares(), e.CountAfterwares(), e.CountBackends(), e.CountAllDataTransforms())
}

func (e *Endpoint) NoCache() bool {
	return checker.IsNil(e.Cache()) || e.Cache().Disabled()
}

func (e *Endpoint) HasResponse() bool {
	return checker.NonNil(e.response)
}

func (e EndpointResponse) ContinueOnError() bool {
	return e.continueOnError
}

func (e EndpointResponse) HasBody() bool {
	return checker.NonNil(e.body)
}

func (e EndpointResponse) Body() *EndpointResponseBody {
	return e.body
}

func (e EndpointResponseBody) HasContentType() bool {
	return e.contentType.IsEnumValid()
}

func (e EndpointResponseBody) HasContentEncoding() bool {
	return e.contentEncoding.IsEnumValid()
}

func (e EndpointResponseBody) ContentType() enum.ContentType {
	return e.contentType
}

func (e EndpointResponseBody) ContentEncoding() enum.ContentEncoding {
	return e.contentEncoding
}

func (e EndpointResponseBody) Aggregate() bool {
	return e.aggregate
}

func (e EndpointResponseBody) OmitEmpty() bool {
	return e.omitEmpty
}

func (e EndpointResponseHeader) Mapper() *Mapper {
	return e.mapper
}

func (e EndpointResponseHeader) Projector() *Projector {
	return e.projector
}

func (e EndpointResponseBody) Mapper() *Mapper {
	return e.mapper
}

func (e EndpointResponseBody) Projector() *Projector {
	return e.projector
}

func (e EndpointResponseBody) HasNomenclature() bool {
	return e.nomenclature.IsEnumValid()
}

func (e EndpointResponseBody) Nomenclature() enum.Nomenclature {
	return e.nomenclature
}

func (e EndpointResponseBody) CountAllDataTransforms() (count int) {
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
