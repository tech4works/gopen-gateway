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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

type Backend struct {
	hosts    []string
	path     string
	method   string
	request  *BackendRequest
	response *BackendResponse
}

type BackendRequest struct {
	omitHeader       bool
	omitQuery        bool
	omitBody         bool
	contentType      enum.ContentType
	contentEncoding  enum.ContentEncoding
	nomenclature     enum.Nomenclature
	omitEmpty        bool
	headerMapper     *Mapper
	queryMapper      *Mapper
	bodyMapper       *Mapper
	headerProjection *Projection
	queryProjection  *Projection
	bodyProjection   *Projection
	headerModifiers  []Modifier
	paramModifiers   []Modifier
	queryModifiers   []Modifier
	bodyModifiers    []Modifier
}

type BackendResponse struct {
	omit             bool
	omitHeader       bool
	omitBody         bool
	group            string
	headerMapper     *Mapper
	bodyMapper       *Mapper
	headerProjection *Projection
	bodyProjection   *Projection
	headerModifiers  []Modifier
	bodyModifiers    []Modifier
}

func NewBackend(
	hosts []string,
	path,
	method string,
	request *BackendRequest,
	response *BackendResponse,
) Backend {
	return Backend{
		hosts:    hosts,
		path:     path,
		method:   method,
		request:  request,
		response: response,
	}
}

func NewBackendRequest(
	omitHeader,
	omitQuery,
	omitBody bool,
	contentType enum.ContentType,
	contentEncoding enum.ContentEncoding,
	nomenclature enum.Nomenclature,
	omitEmpty bool,
	headerMapper,
	queryMapper,
	bodyMapper *Mapper,
	headerProjection,
	queryProjection,
	bodyProjection *Projection,
	headerModifiers,
	paramModifiers,
	queryModifiers,
	bodyModifiers []Modifier,
) *BackendRequest {
	return &BackendRequest{
		omitHeader:       omitHeader,
		omitQuery:        omitQuery,
		omitBody:         omitBody,
		contentType:      contentType,
		contentEncoding:  contentEncoding,
		nomenclature:     nomenclature,
		omitEmpty:        omitEmpty,
		headerMapper:     headerMapper,
		queryMapper:      queryMapper,
		bodyMapper:       bodyMapper,
		headerProjection: headerProjection,
		queryProjection:  queryProjection,
		bodyProjection:   bodyProjection,
		headerModifiers:  headerModifiers,
		paramModifiers:   paramModifiers,
		queryModifiers:   queryModifiers,
		bodyModifiers:    bodyModifiers,
	}
}

func NewBackendResponse(
	omit,
	omitHeader,
	omitBody bool,
	group string,
	headerMapper,
	bodyMapper *Mapper,
	headerProjection,
	bodyProjection *Projection,
	headerModifiers,
	bodyModifiers []Modifier,
) *BackendResponse {
	return &BackendResponse{
		omit:             omit,
		omitHeader:       omitHeader,
		omitBody:         omitBody,
		group:            group,
		headerMapper:     headerMapper,
		bodyMapper:       bodyMapper,
		headerProjection: headerProjection,
		bodyProjection:   bodyProjection,
		headerModifiers:  headerModifiers,
		bodyModifiers:    bodyModifiers,
	}
}

func NewBackendResponseForMiddleware() *BackendResponse {
	return &BackendResponse{
		omit: true,
	}
}

func (b *Backend) Hosts() []string {
	return b.hosts
}

func (b *Backend) Path() string {
	return b.path
}

func (b *Backend) Method() string {
	return b.method
}

func (b *Backend) HasRequest() bool {
	return helper.IsNotNil(b.request)
}

func (b *Backend) Request() *BackendRequest {
	return b.request
}

func (b *Backend) HasResponse() bool {
	return helper.IsNotNil(b.response)
}

func (b *Backend) Response() *BackendResponse {
	return b.response
}

func (b *Backend) CountAllDataTransforms() (count int) {
	if helper.IsNotNil(b.Request()) {
		count += b.Request().CountAllDataTransforms()
	}
	if helper.IsNotNil(b.Response()) {
		count += b.Response().CountAllDataTransforms()
	}
	return count
}

func (b BackendRequest) OmitHeader() bool {
	return b.omitHeader
}

func (b BackendRequest) OmitQuery() bool {
	return b.omitQuery
}

func (b BackendRequest) OmitBody() bool {
	return b.omitBody
}

func (b BackendRequest) HasContentType() bool {
	return b.contentType.IsEnumValid()
}

func (b BackendRequest) HasContentEncoding() bool {
	return b.contentEncoding.IsEnumValid()
}

func (b BackendRequest) ContentType() enum.ContentType {
	return b.contentType
}

func (b BackendRequest) ContentEncoding() enum.ContentEncoding {
	return b.contentEncoding
}

func (b BackendRequest) HasNomenclature() bool {
	return b.nomenclature.IsEnumValid()
}

func (b BackendRequest) Nomenclature() enum.Nomenclature {
	return b.nomenclature
}

func (b BackendRequest) OmitEmpty() bool {
	return b.omitEmpty
}

func (b BackendRequest) HeaderMapper() *Mapper {
	return b.headerMapper
}

func (b BackendRequest) HeaderProjection() *Projection {
	return b.headerProjection
}

func (b BackendRequest) HeaderModifiers() []Modifier {
	return b.headerModifiers
}

func (b BackendRequest) ParamModifiers() []Modifier {
	return b.paramModifiers
}

func (b BackendRequest) QueryProjection() *Projection {
	return b.queryProjection
}

func (b BackendRequest) QueryMapper() *Mapper {
	return b.queryMapper

}

func (b BackendRequest) QueryModifiers() []Modifier {
	return b.queryModifiers
}

func (b BackendRequest) BodyProjection() *Projection {
	return b.bodyProjection
}

func (b BackendRequest) BodyMapper() *Mapper {
	return b.bodyMapper
}

func (b BackendRequest) BodyModifiers() []Modifier {
	return b.bodyModifiers
}

func (b BackendRequest) CountAllDataTransforms() (count int) {
	count += b.CountParamDataTransforms()
	count += b.CountHeaderDataTransforms()
	count += b.CountQueryDataTransforms()
	count += b.CountBodyDataTransforms()
	return count
}

func (b BackendRequest) CountQueryDataTransforms() (count int) {
	if b.OmitQuery() {
		return 1
	}
	if helper.IsNotNil(b.QueryMapper()) {
		count += len(b.QueryMapper().Keys())
	}
	if helper.IsNotNil(b.QueryProjection()) {
		count += len(b.QueryProjection().Keys())
	}
	if helper.IsNotNil(b.QueryModifiers()) {
		count += len(b.QueryModifiers())
	}
	return count
}

func (b BackendRequest) CountHeaderDataTransforms() (count int) {
	if b.OmitHeader() {
		return 1
	}
	if helper.IsNotNil(b.HeaderMapper()) {
		count += len(b.HeaderMapper().Keys())
	}
	if helper.IsNotNil(b.HeaderProjection()) {
		count += len(b.HeaderProjection().Keys())
	}
	if helper.IsNotNil(b.HeaderModifiers()) {
		count += len(b.HeaderModifiers())
	}
	return count
}

func (b BackendRequest) CountBodyDataTransforms() (count int) {
	if b.OmitBody() {
		return 1
	}
	if helper.IsNotNil(b.BodyMapper()) {
		count += len(b.BodyMapper().Keys())
	}
	if helper.IsNotNil(b.BodyProjection()) {
		count += len(b.BodyProjection().Keys())
	}
	if helper.IsNotNil(b.BodyModifiers()) {
		count += len(b.BodyModifiers())
	}
	return count
}

func (b BackendRequest) CountParamDataTransforms() int {
	if helper.IsNotNil(b.ParamModifiers()) {
		return len(b.ParamModifiers())
	}
	return 0
}

func (b BackendResponse) Omit() bool {
	return b.omit
}

func (b BackendResponse) OmitHeader() bool {
	return b.omitHeader
}

func (b BackendResponse) OmitBody() bool {
	return b.omitBody
}

func (b BackendResponse) HeaderProjection() *Projection {
	return b.headerProjection
}

func (b BackendResponse) HeaderMapper() *Mapper {
	return b.headerMapper
}

func (b BackendResponse) HeaderModifiers() []Modifier {
	return b.headerModifiers
}

func (b BackendResponse) BodyProjection() *Projection {
	return b.bodyProjection
}

func (b BackendResponse) BodyMapper() *Mapper {
	return b.bodyMapper
}

func (b BackendResponse) BodyModifiers() []Modifier {
	return b.bodyModifiers
}

func (b BackendResponse) HasGroup() bool {
	return helper.IsNotEmpty(b.group)
}

func (b BackendResponse) Group() string {
	return b.group
}

func (b BackendResponse) CountAllDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	count += b.CountHeaderDataTransforms()
	count += b.CountBodyDataTransforms()
	return count
}

func (b BackendResponse) CountHeaderDataTransforms() (count int) {
	if b.OmitHeader() {
		return 1
	}
	if helper.IsNotNil(b.HeaderMapper()) {
		count += len(b.HeaderMapper().Keys())
	}
	if helper.IsNotNil(b.HeaderProjection()) {
		count += len(b.HeaderProjection().Keys())
	}
	if helper.IsNotNil(b.HeaderModifiers()) {
		count += len(b.HeaderModifiers())
	}
	return count
}

func (b BackendResponse) CountBodyDataTransforms() (count int) {
	if b.OmitBody() {
		return 1
	}
	if helper.IsNotNil(b.BodyMapper()) {
		count += len(b.BodyMapper().Keys())
	}
	if helper.IsNotNil(b.BodyProjection()) {
		count += len(b.BodyProjection().Keys())
	}
	if helper.IsNotNil(b.BodyModifiers()) {
		count += len(b.BodyModifiers())
	}
	return count
}
