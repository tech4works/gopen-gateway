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

// Backend is a type that represents a backend server configuration.
type Backend struct {
	// hosts is an array of host addresses.
	hosts []string
	// path is a string that represents the path of the backend server configuration.
	path string
	// method is the HTTP method to be used for requests to the backend.
	method   string
	request  *BackendRequest
	response *BackendResponse
}

type BackendRequest struct {
	omitHeader       bool
	omitQuery        bool
	omitBody         bool
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
	apply              enum.BackendResponseApply
	omit               bool
	omitHeader         bool
	omitBody           bool
	group              string
	headerMapper       *Mapper
	bodyMapper         *Mapper
	headerProjection   *Projection
	bodyProjection     *Projection
	statusCodeModifier int
	headerModifiers    []Modifier
	bodyModifiers      []Modifier
}

// newBackend creates a new instance of Backend based on the provided BackendJson.
// It assigns the values from the BackendJson fields to the corresponding fields in Backend.
// Returns the newly created Backend instance.
func newBackend(backendJson *BackendJson) Backend {
	return Backend{
		hosts:    backendJson.Hosts,
		path:     backendJson.Path,
		method:   backendJson.Method,
		request:  newBackendRequest(backendJson.Request),
		response: newBackendResponse(backendJson.Response),
	}
}

// newMiddlewareBackend creates a new Backend instance based on the provided backendVO and backendExtraConfigVO.
// It takes the fields from backendVO and assigns them to the corresponding fields in the Backend struct.
// It assigns the backendExtraConfigVO parameter to the extraConfig field of the Backend struct.
// The function returns the created Backend instance.
func newMiddlewareBackend(backend *Backend) *Backend {
	return &Backend{
		hosts:    backend.hosts,
		path:     backend.path,
		method:   backend.method,
		request:  backend.request,
		response: newBackendResponseForMiddleware(),
	}
}

func newBackendResponseForMiddleware() *BackendResponse {
	return &BackendResponse{
		apply: enum.BackendResponseApplyLate,
		omit:  true,
	}
}

func newBackendRequest(backendRequestJson *BackendRequestJson) *BackendRequest {
	if helper.IsNil(backendRequestJson) {
		return nil
	}
	return &BackendRequest{
		omitHeader:       backendRequestJson.OmitHeader,
		omitQuery:        backendRequestJson.OmitQuery,
		omitBody:         backendRequestJson.OmitBody,
		headerMapper:     backendRequestJson.HeaderMapper,
		queryMapper:      backendRequestJson.QueryMapper,
		bodyMapper:       backendRequestJson.BodyMapper,
		headerProjection: backendRequestJson.HeaderProjection,
		queryProjection:  backendRequestJson.QueryProjection,
		bodyProjection:   backendRequestJson.BodyProjection,
	}
}

func newBackendResponse(backendResponseJson *BackendResponseJson) *BackendResponse {
	if helper.IsNil(backendResponseJson) {
		return nil
	}

	var headerModifiers []Modifier
	for _, headerModifierJsonVO := range backendResponseJson.HeaderModifiers {
		headerModifiers = append(headerModifiers, newModifier(headerModifierJsonVO))
	}
	var bodyModifiers []Modifier
	for _, bodyModifierJsonVO := range backendResponseJson.BodyModifiers {
		bodyModifiers = append(bodyModifiers, newModifier(bodyModifierJsonVO))
	}

	return &BackendResponse{
		apply:              backendResponseJson.Apply,
		omit:               backendResponseJson.Omit,
		omitHeader:         backendResponseJson.OmitHeader,
		omitBody:           backendResponseJson.OmitBody,
		group:              backendResponseJson.Group,
		headerMapper:       backendResponseJson.HeaderMapper,
		bodyMapper:         backendResponseJson.BodyMapper,
		headerProjection:   backendResponseJson.HeaderProjection,
		bodyProjection:     backendResponseJson.BodyProjection,
		statusCodeModifier: backendResponseJson.StatusCodeModifier,
		headerModifiers:    headerModifiers,
		bodyModifiers:      bodyModifiers,
	}
}

// BalancedHost returns a balanced host from the Backend instance. If there is only one host, it is returned directly.
// Otherwise, a random host is selected from the available ones.
func (b *Backend) BalancedHost() string {
	// todo: aqui obtemos o host (correto é criar um domínio chamado balancer aonde ele vai retornar o host
	//  disponível pegando como base, se ele esta de pé ou não, e sua config de porcentagem)
	if helper.EqualsLen(b.hosts, 1) {
		return b.hosts[0]
	}
	return b.hosts[helper.RandomNumber(0, len(b.hosts)-1)]
}

// Path returns the path of the Backend instance.
func (b *Backend) Path() string {
	return b.path
}

// Method returns the method of the Backend instance.
func (b *Backend) Method() string {
	return b.method
}

func (b *Backend) Request() *BackendRequest {
	return b.request
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
	// contamos as modificações de param
	count += b.CountParamDataTransforms()
	// contamos as modificações de header
	count += b.CountHeaderDataTransforms()
	// contamos as modificações de query
	count += b.CountQueryDataTransforms()
	// contamos as modificações de body
	return b.CountBodyDataTransforms()
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

func (b BackendResponse) Apply() enum.BackendResponseApply {
	if b.apply.IsEnumValid() {
		return b.apply
	}
	return enum.BackendResponseApplyEarly
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
	// se ele quer omitir a resposta retornamos apenas 1
	if b.Omit() {
		return 1
	}
	// contamos as modificações de header
	count += b.CountHeaderDataTransforms()
	// contamos as modificações de body
	return b.CountBodyDataTransforms()
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
