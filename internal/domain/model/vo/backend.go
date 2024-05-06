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
	path UrlPath
	// method is the HTTP method to be used for requests to the backend.
	method   string
	request  *BackendRequest
	response *BackendResponse
	// modifiers is an instance of BackendModifiers containing modifiers for the backend request and response.
	modifiers *BackendModifiers
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
}

type BackendResponse struct {
	apply            enum.BackendResponseApply
	omit             bool
	omitHeader       bool
	omitBody         bool
	group            string
	headerMapper     *Mapper
	bodyMapper       *Mapper
	headerProjection *Projection
	bodyProjection   *Projection
}

// BackendModifiers is a type that represents the set of modifiers for a backend configuration.
// It contains fields for the status code, header, params, query, and body modifiers.
type BackendModifiers struct {
	// statusCode represents the status code modifier for a BackendModifiers instance.
	// It is an integer value that specifies the desired status code for a httpResponse.
	statusCode int
	// header represents an array of Modifier instances that modify the modifyHeaders of a httpRequest or httpResponse
	// from Endpoint or only current backend.
	header []Modifier
	// param is a field in the BackendModifiers struct.
	// It represents an array of Modifier instances that modify the parameters of a httpRequest from Endpoint or only current
	// backend.
	param []Modifier
	// `query` is a field in the `BackendModifiers` struct. It represents an array of `Modifier` instances that
	// modify the query parameters of a httpRequest from `Endpoint` or only the current backend.
	query []Modifier
	// body represents an array of Modifier instances that modify the body of a httpRequest or httpResponse
	// from Endpoint or only current backend.
	body []Modifier
}

// newBackend creates a new instance of Backend based on the provided BackendJson.
// It assigns the values from the BackendJson fields to the corresponding fields in Backend.
// Returns the newly created Backend instance.
func newBackend(backendJsonVO *BackendJson) Backend {
	return Backend{
		hosts:     backendJsonVO.Hosts,
		path:      backendJsonVO.Path,
		method:    backendJsonVO.Method,
		request:   newBackendRequest(backendJsonVO.Request),
		response:  newBackendResponse(backendJsonVO.Response),
		modifiers: newBackendModifier(backendJsonVO.Modifiers),
	}
}

// newMiddlewareBackend creates a new Backend instance based on the provided backendVO and backendExtraConfigVO.
// It takes the fields from backendVO and assigns them to the corresponding fields in the Backend struct.
// It assigns the backendExtraConfigVO parameter to the extraConfig field of the Backend struct.
// The function returns the created Backend instance.
func newMiddlewareBackend(backendVO *Backend) *Backend {
	return &Backend{
		hosts:     backendVO.hosts,
		path:      backendVO.path,
		method:    backendVO.method,
		request:   backendVO.request,
		response:  newBackendResponseForMiddleware(),
		modifiers: backendVO.modifiers,
	}
}

func newBackendResponseForMiddleware() *BackendResponse {
	return &BackendResponse{
		omit: true,
	}
}

func newBackendRequest(backendRequestJsonVO *BackendRequestJson) *BackendRequest {
	if helper.IsNil(backendRequestJsonVO) {
		return nil
	}
	return &BackendRequest{
		omitHeader:       backendRequestJsonVO.OmitHeader,
		omitQuery:        backendRequestJsonVO.OmitQuery,
		omitBody:         backendRequestJsonVO.OmitBody,
		headerMapper:     backendRequestJsonVO.HeaderMapper,
		queryMapper:      backendRequestJsonVO.QueryMapper,
		bodyMapper:       backendRequestJsonVO.BodyMapper,
		headerProjection: backendRequestJsonVO.HeaderProjection,
		queryProjection:  backendRequestJsonVO.QueryProjection,
		bodyProjection:   backendRequestJsonVO.BodyProjection,
	}
}

func newBackendResponse(backendResponseJsonVO *BackendResponseJson) *BackendResponse {
	if helper.IsNil(backendResponseJsonVO) {
		return nil
	}
	return &BackendResponse{
		apply:            backendResponseJsonVO.Apply,
		omit:             backendResponseJsonVO.Omit,
		omitHeader:       backendResponseJsonVO.OmitHeader,
		omitBody:         backendResponseJsonVO.OmitBody,
		group:            backendResponseJsonVO.Group,
		headerMapper:     backendResponseJsonVO.HeaderMapper,
		bodyMapper:       backendResponseJsonVO.BodyMapper,
		headerProjection: backendResponseJsonVO.HeaderProjection,
		bodyProjection:   backendResponseJsonVO.BodyProjection,
	}
}

// newBackendModifier creates a new instance of BackendModifiers based on the provided BackendModifiersJson.
// If the provided BackendModifiersJson is nil, it returns nil.
// It creates an array of Modifier instances for header, param, query, and body by iterating through the respective slices
// in the BackendModifiersJson and creating a new Modifier instance for each item.
// It constructs the BackendModifiers object and returns it.
func newBackendModifier(backendModifierJsonVO *BackendModifiersJson) *BackendModifiers {
	if helper.IsNil(backendModifierJsonVO) {
		return nil
	}

	var header []Modifier
	for _, modifierDTO := range backendModifierJsonVO.Header {
		header = append(header, *newModifier(&modifierDTO))
	}
	var params []Modifier
	for _, modifierDTO := range backendModifierJsonVO.Param {
		params = append(params, *newModifier(&modifierDTO))
	}
	var query []Modifier
	for _, modifierDTO := range backendModifierJsonVO.Query {
		query = append(query, *newModifier(&modifierDTO))
	}
	var body []Modifier
	for _, modifierDTO := range backendModifierJsonVO.Body {
		body = append(body, *newModifier(&modifierDTO))
	}

	return &BackendModifiers{
		statusCode: backendModifierJsonVO.StatusCode,
		header:     header,
		param:      params,
		query:      query,
		body:       body,
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
func (b *Backend) Path() UrlPath {
	return b.path
}

// Method returns the method of the Backend instance.
func (b *Backend) Method() string {
	return b.method
}

// CountModifiers returns the number of modifiers present in the Backend instance.
// If the modifiers field is not nil, it counts all the modifiers using the CountAll() method of BackendModifiers.
// Otherwise, it returns 0.
func (b *Backend) CountModifiers() int {
	if helper.IsNotNil(b.modifiers) {
		return b.modifiers.CountAll()
	}
	return 0
}

func (b *Backend) Request() *BackendRequest {
	return b.request
}

func (b *Backend) Response() *BackendResponse {
	return b.response
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

func (b BackendRequest) HeaderProjection() *Projection {
	return b.headerProjection
}

func (b BackendRequest) HeaderMapper() *Mapper {
	return b.headerMapper
}

func (b BackendRequest) QueryProjection() *Projection {
	return b.queryProjection
}

func (b BackendRequest) QueryMapper() *Mapper {
	return b.queryMapper

}

func (b BackendRequest) BodyProjection() *Projection {
	return b.bodyProjection
}

func (b BackendRequest) BodyMapper() *Mapper {
	return b.bodyMapper
}

func (r BackendResponse) Apply() enum.BackendResponseApply {
	if r.apply.IsEnumValid() {
		return r.apply
	}
	return enum.BackendResponseApplyEarly
}

func (r BackendResponse) Omit() bool {
	return r.omit
}

func (r BackendResponse) OmitHeader() bool {
	return r.omitHeader
}

func (r BackendResponse) OmitBody() bool {
	return r.omitBody
}

func (r BackendResponse) HeaderProjection() *Projection {
	return r.headerProjection
}

func (r BackendResponse) HeaderMapper() *Mapper {
	return r.headerMapper
}

func (r BackendResponse) BodyProjection() *Projection {
	return r.bodyProjection
}

func (r BackendResponse) BodyMapper() *Mapper {
	return r.bodyMapper
}

func (r BackendResponse) HasGroup() bool {
	return helper.IsNotEmpty(r.group)
}

func (r BackendResponse) Group() string {
	return r.group
}

// CountAll returns the total count of modifiers for a BackendModifiers instance.
// It counts the number of valid `statusCode` and the length of `header`, `params`, `query`, and `body` slices,
// and adds them up to get the total count.
func (b *BackendModifiers) CountAll() (count int) {
	if helper.IsNotEmpty(b.statusCode) {
		count++
	}
	count += len(b.header) + len(b.param) + len(b.query) + len(b.body)
	return count
}
