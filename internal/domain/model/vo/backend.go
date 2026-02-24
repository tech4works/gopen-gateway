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
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type BackendPolymorphicResponse interface {
	OK() bool
	StatusCode() StatusCode
	Header() Header
	HasBody() bool
	Body() *Body
	Map() (map[string]any, error)
	Executed() bool
	Ignored() bool
	Cancelled() bool
}

type Backend struct {
	flow            enum.BackendFlow
	continueOnError bool
	onlyIf          []string
	ignoreIf        []string
	id              string
	dependencies    *BackendDependencies
	kind            enum.BackendKind
	http            *HTTP
	publisher       *Publisher
	response        *BackendResponse
}

type HTTP struct {
	hosts   []string
	path    string
	method  string
	request *BackendRequest
}

type BackendRequest struct {
	continueOnError bool
	concurrent      int
	async           bool
	header          *BackendRequestHeader
	urlPath         *BackendRequestURLPath
	query           *BackendRequestQuery
	body            *BackendRequestBody
}

type BackendResponse struct {
	continueOnError bool
	omit            bool
	header          *BackendResponseHeader
	body            *BackendResponseBody
}

type BackendResponseHeader struct {
	omit      bool
	mapper    *Mapper
	projector *Projector
	modifiers []Modifier
}

type BackendResponseBody struct {
	omit      bool
	group     string
	mapper    *Mapper
	projector *Projector
	modifiers []Modifier
	joins     []Join
}

type BackendRequestHeader struct {
	omit      bool
	mapper    *Mapper
	projector *Projector
	modifiers []Modifier
}

type BackendRequestURLPath struct {
	modifiers []Modifier
}

type BackendRequestQuery struct {
	omit      bool
	mapper    *Mapper
	projector *Projector
	modifiers []Modifier
}

type BackendDependencies struct {
	ids  []string
	idxs []int
}

type BackendRequestBody struct {
	omit            bool
	omitEmpty       bool
	contentType     enum.ContentType
	contentEncoding enum.ContentEncoding
	nomenclature    enum.Nomenclature
	mapper          *Mapper
	projector       *Projector
	modifiers       []Modifier
	joins           []Join
}

func NewBackendHTTP(
	flow enum.BackendFlow,
	onlyIf []string,
	ignoreIf []string,
	id string,
	dependencies *BackendDependencies,
	hosts []string,
	path string,
	method string,
	request *BackendRequest,
	response *BackendResponse,
) Backend {
	return Backend{
		flow:         flow,
		onlyIf:       onlyIf,
		ignoreIf:     ignoreIf,
		id:           id,
		dependencies: dependencies,
		kind:         enum.BackendKindHTTP,
		http: &HTTP{
			hosts:   hosts,
			path:    path,
			method:  method,
			request: request,
		},
		response: response,
	}
}

func NewBackendPublisher(
	flow enum.BackendFlow,
	onlyIf []string,
	ignoreIf []string,
	id string,
	dependencies *BackendDependencies,
	publisher Publisher,
	response *BackendResponse,
) Backend {
	return Backend{
		flow:         flow,
		onlyIf:       onlyIf,
		ignoreIf:     ignoreIf,
		id:           id,
		dependencies: dependencies,
		kind:         enum.BackendKindPublisher,
		publisher:    &publisher,
		response:     response,
	}
}

func NewBackendRequest(
	continueOnError bool,
	concurrent int,
	async bool,
	header *BackendRequestHeader,
	param *BackendRequestURLPath,
	query *BackendRequestQuery,
	body *BackendRequestBody,
) *BackendRequest {
	return &BackendRequest{
		continueOnError: continueOnError,
		concurrent:      concurrent,
		async:           async,
		header:          header,
		urlPath:         param,
		query:           query,
		body:            body,
	}
}

func NewBackendRequestHeader(
	omit bool,
	mapper *Mapper,
	projector *Projector,
	modifiers []Modifier,
) *BackendRequestHeader {
	return &BackendRequestHeader{
		omit:      omit,
		mapper:    mapper,
		projector: projector,
		modifiers: modifiers,
	}
}

func NewBackendRequestURLPath(modifiers []Modifier) *BackendRequestURLPath {
	return &BackendRequestURLPath{modifiers: modifiers}
}

func NewBackendRequestQuery(
	omit bool,
	mapper *Mapper,
	projector *Projector,
	modifiers []Modifier,
) *BackendRequestQuery {
	return &BackendRequestQuery{
		omit:      omit,
		mapper:    mapper,
		projector: projector,
		modifiers: modifiers,
	}
}

func NewBackendRequestBody(
	omit bool,
	omitEmpty bool,
	contentType enum.ContentType,
	contentEncoding enum.ContentEncoding,
	nomenclature enum.Nomenclature,
	mapper *Mapper,
	projector *Projector,
	modifiers []Modifier,
	joins []Join,
) *BackendRequestBody {
	return &BackendRequestBody{
		omit:            omit,
		omitEmpty:       omitEmpty,
		contentType:     contentType,
		contentEncoding: contentEncoding,
		nomenclature:    nomenclature,
		mapper:          mapper,
		projector:       projector,
		modifiers:       modifiers,
		joins:           joins,
	}
}

func NewBackendResponseForMiddleware(
	continueOnError,
	omit bool,
	header *BackendResponseHeader,
) *BackendResponse {
	return &BackendResponse{
		continueOnError: continueOnError,
		omit:            omit,
		header:          header,
		body: &BackendResponseBody{
			omit: true,
		},
	}
}

func NewBackendResponse(
	continueOnError,
	omit bool,
	header *BackendResponseHeader,
	body *BackendResponseBody,
) *BackendResponse {
	return &BackendResponse{
		continueOnError: continueOnError,
		omit:            omit,
		header:          header,
		body:            body,
	}
}

func NewBackendResponseHeader(
	omit bool,
	mapper *Mapper,
	projector *Projector,
	modifiers []Modifier,
) *BackendResponseHeader {
	return &BackendResponseHeader{
		omit:      omit,
		mapper:    mapper,
		projector: projector,
		modifiers: modifiers,
	}
}

func NewBackendResponseBody(
	omit bool,
	group string,
	mapper *Mapper,
	projector *Projector,
	modifiers []Modifier,
	joins []Join,
) *BackendResponseBody {
	return &BackendResponseBody{
		omit:      omit,
		group:     group,
		mapper:    mapper,
		projector: projector,
		modifiers: modifiers,
		joins:     joins,
	}
}

func NewBackendDependencies(ids []string, idxs []int) *BackendDependencies {
	return &BackendDependencies{
		ids:  ids,
		idxs: idxs,
	}
}

func (b *Backend) OnlyIf() []string {
	return b.onlyIf
}

func (b *Backend) IgnoreIf() []string {
	return b.ignoreIf
}

func (b *Backend) ID() string {
	return b.id
}

func (b *Backend) HasDependencies() bool {
	return checker.NonNil(b.dependencies)
}

func (b *Backend) Dependencies() *BackendDependencies {
	return b.dependencies
}

func (b *Backend) HTTP() *HTTP {
	return b.http
}

func (b *Backend) Publisher() *Publisher {
	return b.publisher
}

func (b *Backend) Kind() enum.BackendKind {
	return b.kind
}

func (b *Backend) IsHTTP() bool {
	return checker.Equals(b.kind, enum.BackendKindHTTP)
}

func (b *Backend) IsPublisher() bool {
	return checker.Equals(b.kind, enum.BackendKindPublisher)
}

func (b *Backend) Async() bool {
	if b.IsHTTP() {
		if b.HTTP().HasRequest() {
			return b.HTTP().Request().Async()
		}
		return false
	}

	if b.IsPublisher() {
		return b.publisher.Async()
	}

	return false
}

func (b *Backend) IsBeforeware() bool {
	return checker.Equals(b.flow, enum.BackendFlowBeforeware)
}

func (b *Backend) IsNormal() bool {
	return checker.Equals(b.flow, enum.BackendFlowNormal)
}

func (b *Backend) IsAfterware() bool {
	return checker.Equals(b.flow, enum.BackendFlowAfterware)
}

func (b *Backend) Flow() enum.BackendFlow {
	return b.flow
}

func (b *Backend) IsMiddleware() bool {
	return b.IsBeforeware() || b.IsAfterware()
}

func (b *Backend) HasResponse() bool {
	return checker.NonNil(b.response)
}

func (b *Backend) Response() *BackendResponse {
	return b.response
}

func (b *HTTP) Hosts() []string {
	return b.hosts
}

func (b *HTTP) Path() string {
	return b.path
}

func (b *HTTP) Method() string {
	return b.method
}

func (b *HTTP) HasRequest() bool {
	return checker.NonNil(b.request)
}

func (b *HTTP) Request() *BackendRequest {
	return b.request
}

func (b *HTTP) CountAllDataTransforms() (count int) {
	if b.HasRequest() {
		count += b.Request().CountAllDataTransforms()
	}
	return count
}

func (b *Backend) CountResponseDataTransforms() (count int) {
	if b.HasResponse() {
		count += b.Response().CountAllDataTransforms()
	}
	return count
}

func (b *Backend) CountAllDataTransforms() (count int) {
	switch b.Kind() {
	case enum.BackendKindHTTP:
		count += b.HTTP().CountAllDataTransforms()
	case enum.BackendKindPublisher:
		count += b.Publisher().CountAllDataTransforms()
	}
	count += b.CountResponseDataTransforms()
	return count
}

func (b *BackendRequest) IsConcurrent() bool {
	return checker.IsGreaterThanOrEqual(b.concurrent, 2)
}

func (b *BackendRequest) Concurrent() int {
	return b.concurrent
}

func (b *BackendRequest) Async() bool {
	return b.async
}

func (b *BackendRequest) CountAllDataTransforms() (count int) {
	if b.HasURLPath() {
		count += b.URLPath().CountDataTransforms()
	}
	if b.HasQuery() {
		count += b.Query().CountDataTransforms()
	}
	if b.HasHeader() {
		count += b.Header().CountDataTransforms()

	}
	if b.HasBody() {
		count += b.Body().CountDataTransforms()
	}
	return count
}

func (b BackendRequestHeader) Omit() bool {
	return b.omit
}

func (b BackendRequestQuery) Omit() bool {
	return b.omit
}

func (b BackendRequestBody) Omit() bool {
	return b.omit
}

func (b BackendRequestBody) HasContentType() bool {
	return b.contentType.IsEnumValid()
}

func (b BackendRequestBody) HasContentEncoding() bool {
	return b.contentEncoding.IsEnumValid()
}

func (b BackendRequestBody) ContentType() enum.ContentType {
	return b.contentType
}

func (b BackendRequestBody) ContentEncoding() enum.ContentEncoding {
	return b.contentEncoding
}

func (b BackendRequestBody) HasNomenclature() bool {
	return b.nomenclature.IsEnumValid()
}

func (b BackendRequestBody) Nomenclature() enum.Nomenclature {
	return b.nomenclature
}

func (b BackendRequestBody) OmitEmpty() bool {
	return b.omitEmpty
}

func (b BackendRequestHeader) Mapper() *Mapper {
	return b.mapper
}

func (b BackendRequestHeader) Projector() *Projector {
	return b.projector
}

func (b BackendRequestHeader) Modifiers() []Modifier {
	return b.modifiers
}

func (b BackendRequestURLPath) Modifiers() []Modifier {
	return b.modifiers
}

func (b BackendRequestQuery) Projector() *Projector {
	return b.projector
}

func (b BackendRequestQuery) Modifiers() []Modifier {
	return b.modifiers
}

func (b BackendRequestQuery) Mapper() *Mapper {
	return b.mapper
}

func (b BackendRequestBody) Projector() *Projector {
	return b.projector
}

func (b BackendRequestBody) Mapper() *Mapper {
	return b.mapper
}

func (b BackendRequestBody) Modifiers() []Modifier {
	return b.modifiers
}

func (b BackendRequestBody) Joins() []Join {
	return b.joins
}

func (b BackendRequestURLPath) CountDataTransforms() int {
	if b.HasModifiers() {
		return len(b.Modifiers())
	}
	return 0
}

func (b BackendRequestURLPath) HasModifiers() bool {
	return checker.IsNotEmpty(b.modifiers)
}

func (b BackendRequestQuery) CountDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	if b.HasMapper() {
		count += len(b.Mapper().Map().Keys())
	}
	if b.HasProjector() {
		count += len(b.Projector().Project().Keys())
	}
	if b.HasModifiers() {
		count += len(b.Modifiers())
	}
	return count
}

func (b BackendRequestQuery) HasMapper() bool {
	return checker.NonNil(b.mapper)
}

func (b BackendRequestQuery) HasProjector() bool {
	return checker.NonNil(b.projector)
}

func (b BackendRequestQuery) HasModifiers() bool {
	return checker.IsNotEmpty(b.modifiers)
}

func (b BackendRequestHeader) CountDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	if b.HasMapper() {
		count += len(b.Mapper().Map().Keys())
	}
	if b.HasProjector() {
		count += len(b.Projector().Project().Keys())
	}
	if b.HasModifiers() {
		count += len(b.Modifiers())
	}
	return count
}

func (b BackendRequestHeader) HasMapper() bool {
	return checker.NonNil(b.mapper)
}

func (b BackendRequestHeader) HasProjector() bool {
	return checker.NonNil(b.projector)
}

func (b BackendRequestHeader) HasModifiers() bool {
	return checker.IsNotEmpty(b.modifiers)
}

func (b BackendRequestBody) CountDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	if b.HasMapper() {
		count += len(b.Mapper().Map().Keys())
	}
	if b.HasProjector() {
		count += len(b.Projector().Project().Keys())
	}
	if b.HasModifiers() {
		count += len(b.Modifiers())
	}
	if b.HasJoins() {
		count += len(b.Joins())
	}
	return count
}

func (b BackendRequestBody) HasMapper() bool {
	return checker.NonNil(b.mapper)
}

func (b BackendRequestBody) HasProjector() bool {
	return checker.NonNil(b.projector)
}

func (b BackendRequestBody) HasModifiers() bool {
	return checker.IsNotEmpty(b.modifiers)
}

func (b BackendRequestBody) HasJoins() bool {
	return checker.IsNotEmpty(b.joins)
}

func (b BackendRequest) ContinueOnError() bool {
	return b.continueOnError
}

func (b BackendRequest) HasURLPath() bool {
	return checker.NonNil(b.urlPath)
}

func (b BackendRequest) URLPath() *BackendRequestURLPath {
	return b.urlPath
}

func (b BackendRequest) HasQuery() bool {
	return checker.NonNil(b.query)
}

func (b BackendRequest) Query() *BackendRequestQuery {
	return b.query
}

func (b BackendRequest) HasHeader() bool {
	return checker.NonNil(b.header)
}

func (b BackendRequest) Header() *BackendRequestHeader {
	return b.header
}

func (b BackendRequest) HasBody() bool {
	return checker.NonNil(b.body)
}

func (b BackendRequest) Body() *BackendRequestBody {
	return b.body
}

func (b BackendRequest) HasDataTransforms() bool {
	return checker.IsGreaterThan(b.CountAllDataTransforms(), 0)
}

func (b BackendResponse) ContinueOnError() bool {
	return b.continueOnError
}

func (b BackendResponse) Omit() bool {
	return b.omit
}

func (b BackendResponse) CountAllDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	if b.HasHeader() {
		count += b.Header().CountDataTransforms()
	}
	if b.HasBody() {
		count += b.Body().CountDataTransforms()
	}
	return count
}

func (b BackendResponse) HasHeader() bool {
	return checker.NonNil(b.header)
}

func (b BackendResponse) Header() *BackendResponseHeader {
	return b.header
}

func (b BackendResponse) HasBody() bool {
	return checker.NonNil(b.body)
}

func (b BackendResponse) Body() *BackendResponseBody {
	return b.body
}

func (b BackendResponseHeader) Omit() bool {
	return b.omit
}

func (b BackendResponseBody) Omit() bool {
	return b.omit
}

func (b BackendResponseHeader) Projector() *Projector {
	return b.projector
}

func (b BackendResponseHeader) Mapper() *Mapper {
	return b.mapper
}

func (b BackendResponseHeader) Modifiers() []Modifier {
	return b.modifiers
}

func (b BackendResponseBody) Projector() *Projector {
	return b.projector
}

func (b BackendResponseBody) Mapper() *Mapper {
	return b.mapper
}

func (b BackendResponseBody) Modifiers() []Modifier {
	return b.modifiers
}

func (b BackendResponseBody) Joins() []Join {
	return b.joins
}

func (b BackendResponseBody) HasMapper() bool {
	return checker.NonNil(b.mapper)
}

func (b BackendResponseBody) HasProjector() bool {
	return checker.NonNil(b.projector)
}

func (b BackendResponseBody) HasModifiers() bool {
	return checker.IsNotEmpty(b.modifiers)
}

func (b BackendResponseBody) HasJoins() bool {
	return checker.IsNotEmpty(b.joins)
}

func (b BackendResponseBody) Group() string {
	return b.group
}

func (b BackendResponseBody) HasGroup() bool {
	return checker.IsNotEmpty(b.group)
}

func (b BackendResponseHeader) HasMapper() bool {
	return checker.NonNil(b.mapper)
}

func (b BackendResponseHeader) HasProjector() bool {
	return checker.NonNil(b.projector)
}

func (b BackendResponseHeader) HasModifiers() bool {
	return checker.IsNotEmpty(b.modifiers)
}

func (b BackendResponseHeader) CountDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	if b.HasMapper() {
		count += len(b.Mapper().Map().Keys())
	}
	if b.HasProjector() {
		count += len(b.Projector().Project().Keys())
	}
	if b.HasModifiers() {
		count += len(b.Modifiers())
	}
	return count
}

func (b BackendResponseBody) CountDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	if b.HasMapper() {
		count += len(b.Mapper().Map().Keys())
	}
	if b.HasProjector() {
		count += len(b.Projector().Project().Keys())
	}
	if b.HasModifiers() {
		count += len(b.Modifiers())
	}
	if b.HasJoins() {
		count += len(b.Joins())
	}
	return count
}

func (d *BackendDependencies) IDs() []string {
	return d.ids
}

func (d *BackendDependencies) Indexes() []int {
	return d.idxs
}
