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

// Backend represents the backend configuration for an application or service.
type Backend struct {
	// hosts represents a slice of strings that contains the hosts configuration for a Backend instance.
	hosts []string
	// path represents the path configuration for a Backend instance.
	// It is a string that specifies the path for the backend request.
	path string
	// Method represents the method configuration for a Backend instance.
	method string
	// request represents the request configuration for a Backend instance.
	request *BackendRequest
	// response is a pointer to a BackendResponse object.
	// BackendResponse represents the response configuration for a Backend instance.
	response *BackendResponse
}

// BackendRequest represents the configuration for a backend request.
type BackendRequest struct {
	// omitHeader represents a boolean flag indicating whether the header should be omitted in the BackendRequest config.
	omitHeader bool
	// omitQuery represents a boolean flag indicating whether the query should be omitted in the BackendRequest config.
	omitQuery bool
	// omitBody represents a boolean flag indicating whether the body should be omitted in the BackendRequest config.
	omitBody bool
	// contentType represents the encoding format of the backend request body.
	contentType enum.ContentType
	// contentEncoding represents the encoding format for the backend request.
	contentEncoding enum.ContentEncoding
	// nomenclature is an enumeration type representing the case format for text values.
	nomenclature enum.Nomenclature
	// omitEmpty is a boolean value that indicates whether the JSON field in the body should be omitted or not if empty.
	omitEmpty bool
	// headerMapper is a field of type *Mapper in the BackendRequest struct.
	// It represents the mapper for the header fields in the backend request config.
	// The mapper is used to map keys to values for the header fields in the request.
	headerMapper *Mapper
	// queryMapper is a field of type *Mapper in the BackendRequest struct. It represents
	// the mapper for the query fields in the backend request config. The mapper is used to
	// map keys to values for the query fields in the request.
	queryMapper *Mapper
	// bodyMapper represents the mapper for the body fields in the BackendRequest config.
	// The mapper is used to map keys to values for the body fields in the request.
	bodyMapper *Mapper
	// headerProjection is a field of type *Projection in the BackendRequest struct.
	// It represents the mapper for the header fields in the backend request config.
	// The projection is used to map header keys to projection values for the request.
	headerProjection *Projection
	// queryProjection is a pointer to a Projection struct, which represents the projection configuration for a
	// backend request's query parameters.
	queryProjection *Projection
	// bodyProjection is a pointer to the Projection struct. It represents the projection configuration for the
	// request body data. The `Projection` struct is used to define which fields should be included or excluded in
	// the projection of the request body data. The `bodyProjection` field can be accessed using the BodyProjection()
	// method of the BackendRequest struct.
	bodyProjection *Projection
	// headerModifiers is a slice of Modifier that represents modifications that can be applied to
	// the headers of a BackendRequest in the Gopen application.
	headerModifiers []Modifier
	// paramModifiers is a slice of Modifier structs that represents the modifications
	// that can be applied to the parameters of a BackendRequest in the Gopen application.
	paramModifiers []Modifier
	// queryModifiers is a slice of Modifier that represents the modifications to be applied to the query of a
	// BackendRequest.
	queryModifiers []Modifier
	// bodyModifiers is a slice of `Modifier` structs representing the modifications that can be applied to the body
	// of a backend request.
	bodyModifiers []Modifier
}

// BackendResponse represents the response configuration for a Backend instance.
type BackendResponse struct {
	// apply represents the scope of applying a BackendResponse.
	// It is used to indicate whether the response config should be applied early or late.
	apply enum.BackendResponseApply
	// omit is a boolean field that represents whether a response should be omitted or not.
	// If true, the response will be omitted.
	omit bool
	// omitHeader is a boolean field that represents whether a response header should be omitted or not.
	// If true, the response header will be omitted.
	omitHeader bool
	// omitBody represents a boolean field that indicates whether a response body should be omitted or not.
	// If true, the response body will be omitted.
	omitBody bool
	// group is a field of type string that represents a grouping identifier for a BackendResponse instance.
	group string
	// headerMapper is a field of type *Mapper in the BackendResponse struct.
	// It represents a mapper for mapping response headers in a BackendResponse instance.
	headerMapper *Mapper
	// bodyMapper is a field of type *Mapper in the BackendResponse struct.
	// It represents a mapper for mapping response bodies in a BackendResponse instance.
	bodyMapper *Mapper
	// headerProjection is a field in the BackendResponse struct. It represents a projection
	// for mapping response headers in a BackendResponse instance.
	headerProjection *Projection
	// bodyProjection represents the projection configuration for the response body of a BackendResponse instance.
	bodyProjection *Projection
	// headerModifiers represents an array of Modifier objects that define modifications to be applied to the response
	// headers in a BackendResponse.
	headerModifiers []Modifier
	// bodyModifiers is a field in the `BackendResponse` struct that represents a list of modifiers to be applied to
	// the response body.
	bodyModifiers []Modifier
}

// newBackend creates a new Backend instance based on the provided backendJson.
// It takes the fields from backendJson and assigns them to the corresponding fields in the Backend struct.
// The function returns the created Backend instance.
func newBackend(backendJson *BackendJson) Backend {
	return Backend{
		hosts:    backendJson.Hosts,
		path:     backendJson.Path,
		method:   backendJson.Method,
		request:  newBackendRequest(backendJson.Request),
		response: newBackendResponse(backendJson.Response),
	}
}

// newMiddlewareBackend creates a new Backend instance based on the provided backend.
// It takes the fields from backend and assigns them to the corresponding fields in the Backend struct.
// The response field is set to the result of calling newBackendResponseForMiddleware().
// The function returns a pointer to the created Backend instance.
func newMiddlewareBackend(backend *Backend) *Backend {
	return &Backend{
		hosts:    backend.hosts,
		path:     backend.path,
		method:   backend.method,
		request:  backend.request,
		response: newBackendResponseForMiddleware(),
	}
}

// newBackendResponseForMiddleware creates a new BackendResponse instance with default values
// for middleware purposes. It sets the apply field to enum.BackendResponseApplyLate and
// the omit field to true. The function returns a pointer to the created BackendResponse instance.
func newBackendResponseForMiddleware() *BackendResponse {
	return &BackendResponse{
		apply: enum.BackendResponseApplyLate,
		omit:  true,
	}
}

// newBackendRequest creates a new BackendRequest instance based on the provided backendRequestJson.
// It takes the fields from backendRequestJson and assigns them to the corresponding fields in the BackendRequest struct.
// The function returns the created BackendRequest instance.
// It also initializes the headerModifiers, paramModifiers, queryModifiers, and bodyModifiers by iterating over the
// corresponding fields in the backendRequestJson.
func newBackendRequest(backendRequestJson *BackendRequestJson) *BackendRequest {
	if helper.IsNil(backendRequestJson) {
		return nil
	}

	var headerModifiers []Modifier
	for _, modifierJson := range backendRequestJson.HeaderModifiers {
		headerModifiers = append(headerModifiers, newModifier(modifierJson))
	}
	var paramModifiers []Modifier
	for _, modifierJson := range backendRequestJson.ParamModifiers {
		paramModifiers = append(paramModifiers, newModifier(modifierJson))
	}
	var queryModifiers []Modifier
	for _, modifierJson := range backendRequestJson.QueryModifiers {
		queryModifiers = append(queryModifiers, newModifier(modifierJson))
	}
	var bodyModifiers []Modifier
	for _, modifierJson := range backendRequestJson.BodyModifiers {
		bodyModifiers = append(bodyModifiers, newModifier(modifierJson))
	}

	return &BackendRequest{
		omitHeader:       backendRequestJson.OmitHeader,
		omitQuery:        backendRequestJson.OmitQuery,
		omitBody:         backendRequestJson.OmitBody,
		contentType:      backendRequestJson.ContentType,
		contentEncoding:  backendRequestJson.ContentEncoding,
		nomenclature:     backendRequestJson.Nomenclature,
		omitEmpty:        backendRequestJson.OmitEmpty,
		headerMapper:     backendRequestJson.HeaderMapper,
		queryMapper:      backendRequestJson.QueryMapper,
		bodyMapper:       backendRequestJson.BodyMapper,
		headerProjection: backendRequestJson.HeaderProjection,
		queryProjection:  backendRequestJson.QueryProjection,
		bodyProjection:   backendRequestJson.BodyProjection,
		headerModifiers:  headerModifiers,
		paramModifiers:   paramModifiers,
		queryModifiers:   queryModifiers,
		bodyModifiers:    bodyModifiers,
	}
}

// newBackendResponse creates a new BackendResponse instance based on the provided backendResponseJson.
// It checks if the backendResponseJson is nil, if it is, it returns a nil BackendResponse.
// It then creates an empty array for headerModifiers and bodyModifiers, and populates them by iterating over the
// corresponding JSON arrays.
// Finally, it returns a pointer to the created BackendResponse instance.
func newBackendResponse(backendResponseJson *BackendResponseJson) *BackendResponse {
	if helper.IsNil(backendResponseJson) {
		return nil
	}

	var headerModifiers []Modifier
	for _, headerModifierJson := range backendResponseJson.HeaderModifiers {
		headerModifiers = append(headerModifiers, newModifier(headerModifierJson))
	}
	var bodyModifiers []Modifier
	for _, bodyModifierJson := range backendResponseJson.BodyModifiers {
		bodyModifiers = append(bodyModifiers, newModifier(bodyModifierJson))
	}

	return &BackendResponse{
		apply:            backendResponseJson.Apply,
		omit:             backendResponseJson.Omit,
		omitHeader:       backendResponseJson.OmitHeader,
		omitBody:         backendResponseJson.OmitBody,
		group:            backendResponseJson.Group,
		headerMapper:     backendResponseJson.HeaderMapper,
		bodyMapper:       backendResponseJson.BodyMapper,
		headerProjection: backendResponseJson.HeaderProjection,
		bodyProjection:   backendResponseJson.BodyProjection,
		headerModifiers:  headerModifiers,
		bodyModifiers:    bodyModifiers,
	}
}

// BalancedHost returns a randomly selected host from the backend's hosts slice.
// If there is only one host, it will always be returned.
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

// Request returns the BackendRequest object associated with the Backend instance.
func (b *Backend) Request() *BackendRequest {
	return b.request
}

// Response returns the BackendResponse object associated with the Backend instance.
func (b *Backend) Response() *BackendResponse {
	return b.response
}

// CountAllDataTransforms returns the total number of data transforms performed by the Backend instance.
// It includes the number of data transforms performed by the request and the response objects associated with the Backend.
// The count is obtained by summing the counts of data transforms from the request and the response.
func (b *Backend) CountAllDataTransforms() (count int) {
	if helper.IsNotNil(b.Request()) {
		count += b.Request().CountAllDataTransforms()
	}
	if helper.IsNotNil(b.Response()) {
		count += b.Response().CountAllDataTransforms()
	}
	return count
}

// OmitHeader returns a boolean value indicating whether the header should be omitted or not.
func (b BackendRequest) OmitHeader() bool {
	return b.omitHeader
}

// OmitQuery returns a boolean value indicating whether the query should be omitted or not.
func (b BackendRequest) OmitQuery() bool {
	return b.omitQuery
}

// OmitBody returns a boolean value indicating whether the body should be omitted or not.
func (b BackendRequest) OmitBody() bool {
	return b.omitBody
}

func (b BackendRequest) HasContentType() bool {
	return b.contentType.IsEnumValid()
}

func (b BackendRequest) HasContentEncoding() bool {
	return b.contentEncoding.IsEnumValid()
}

func (b BackendRequest) ContentType() ContentType {
	switch b.contentType {
	default:
		return NewContentTypeTextPlain()
	case enum.ContentTypeJson:
		return NewContentTypeJson()
	case enum.ContentTypeXml:
		return NewContentTypeXml()
	}
}

func (b BackendRequest) ContentEncoding() ContentEncoding {
	switch b.contentEncoding {
	default:
		return ""
	case enum.ContentEncodingGzip:
		return NewContentEncodingGzip()
	case enum.ContentEncodingDeflate:
		return NewContentEncodingDeflate()
	}
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

// HeaderMapper returns the header mapper configuration of the BackendRequest.
func (b BackendRequest) HeaderMapper() *Mapper {
	return b.headerMapper
}

// HeaderProjection returns the header projection configuration of the BackendRequest.
func (b BackendRequest) HeaderProjection() *Projection {
	return b.headerProjection
}

// HeaderModifiers returns the slice of header modifiers specified in the BackendRequest.
func (b BackendRequest) HeaderModifiers() []Modifier {
	return b.headerModifiers
}

// ParamModifiers returns the slice of param modifiers specified in the BackendRequest.
func (b BackendRequest) ParamModifiers() []Modifier {
	return b.paramModifiers
}

// QueryProjection returns the query projection configuration of the BackendRequest.
func (b BackendRequest) QueryProjection() *Projection {
	return b.queryProjection
}

// QueryMapper returns the query mapper configuration of the BackendRequest.
func (b BackendRequest) QueryMapper() *Mapper {
	return b.queryMapper

}

// QueryModifiers returns the slice of query modifiers specified in the BackendRequest.
func (b BackendRequest) QueryModifiers() []Modifier {
	return b.queryModifiers
}

// BodyProjection returns the body projection configuration of the BackendRequest.
func (b BackendRequest) BodyProjection() *Projection {
	return b.bodyProjection
}

// BodyMapper returns the body mapper configuration of the BackendRequest.
func (b BackendRequest) BodyMapper() *Mapper {
	return b.bodyMapper
}

// BodyModifiers returns the slice of body modifiers specified in the BackendRequest.
func (b BackendRequest) BodyModifiers() []Modifier {
	return b.bodyModifiers
}

// CountAllDataTransforms returns the total number of data transformations applied to the BackendRequest.
// It calculates the count by summing the counts from CountParamDataTransforms(),
// CountHeaderDataTransforms(), CountQueryDataTransforms(), and CountBodyDataTransforms(),
// and returns the final count.
func (b BackendRequest) CountAllDataTransforms() (count int) {
	count += b.CountParamDataTransforms()
	count += b.CountHeaderDataTransforms()
	count += b.CountQueryDataTransforms()
	count += b.CountBodyDataTransforms()
	return count
}

// CountQueryDataTransforms calculates the number of data transformations applied to the query
// of the BackendRequest. If the query is omitted, it returns 1. Otherwise, it adds the number
// of keys in the query mapper, the number of keys in the query projection, and the total number
// of query modifiers. The calculated count is then returned.
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

// CountHeaderDataTransforms calculates the number of data transformations applied
// to the header of the BackendRequest. If the header is omitted, it returns 1.
// Otherwise, it adds the number of keys in the header mapper, the number of keys
// in the header projection, and the total number of header modifiers.
// The calculated count is then returned.
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

// CountBodyDataTransforms calculates the number of data transformations applied
// to the body of the BackendRequest. If the body is omitted, it returns 1.
// Otherwise, it adds the number of keys in the body mapper, the number of keys
// in the body projection, and the total number of body modifiers.
// The calculated count is then returned.
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

// CountParamDataTransforms returns the number of data transformations applied to the param
// of the BackendRequest. If the param is nil, it returns 0. Otherwise, it returns the number
// of param modifiers in the BackendRequest.
func (b BackendRequest) CountParamDataTransforms() int {
	if helper.IsNotNil(b.ParamModifiers()) {
		return len(b.ParamModifiers())
	}
	return 0
}

// Apply returns the enum.BackendResponseApply value of the BackendResponse instance.
// If the apply value is a valid enum, it is returned. Otherwise, enum.BackendResponseApplyEarly is returned.
func (b BackendResponse) Apply() enum.BackendResponseApply {
	if b.apply.IsEnumValid() {
		return b.apply
	}
	return enum.BackendResponseApplyEarly
}

// Omit returns a boolean value indicating whether the backend response should be omitted.
func (b BackendResponse) Omit() bool {
	return b.omit
}

// OmitHeader returns a boolean value indicating whether the backend response's header should be omitted.
func (b BackendResponse) OmitHeader() bool {
	return b.omitHeader
}

// OmitBody returns a boolean value indicating whether the backend response's body should be omitted.
func (b BackendResponse) OmitBody() bool {
	return b.omitBody
}

// HeaderProjection returns the Projection object that represents the header projection configuration
// for the BackendResponse instance.
func (b BackendResponse) HeaderProjection() *Projection {
	return b.headerProjection
}

// HeaderMapper returns the Mapper object that represents the header mapping configuration
// for the BackendResponse instance.
func (b BackendResponse) HeaderMapper() *Mapper {
	return b.headerMapper
}

// HeaderModifiers returns the slice of Modifier objects associated with the BackendResponse instance.
func (b BackendResponse) HeaderModifiers() []Modifier {
	return b.headerModifiers
}

// BodyProjection returns the Projection object that represents the body projection configuration
// for the BackendResponse instance.
func (b BackendResponse) BodyProjection() *Projection {
	return b.bodyProjection
}

// BodyMapper returns the Mapper object that represents the body mapping configuration for the BackendResponse instance.
// It maps the current body based on the configuration provided.
func (b BackendResponse) BodyMapper() *Mapper {
	return b.bodyMapper
}

// BodyModifiers returns the slice of Modifier objects associated with the BackendResponse instance.
func (b BackendResponse) BodyModifiers() []Modifier {
	return b.bodyModifiers
}

// HasGroup returns a boolean value indicating whether the BackendResponse instance has a group value.
// It checks if the group field is not empty.
func (b BackendResponse) HasGroup() bool {
	return helper.IsNotEmpty(b.group)
}

// Group returns the group value of the BackendResponse instance.
func (b BackendResponse) Group() string {
	return b.group
}

// CountAllDataTransforms returns the total number of data transforms applied to the BackendResponse instance.
// It includes both header and body data transforms.
// If the response is omitted, it returns 1.
func (b BackendResponse) CountAllDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	count += b.CountHeaderDataTransforms()
	count += b.CountBodyDataTransforms()
	return count
}

// CountHeaderDataTransforms returns the number of data transforms applied to the header of the BackendResponse instance.
// If the header is omitted, it returns 1. Otherwise, it counts the number of keys in the header mapper,
// the number of keys in the header projection, and the number of header modifiers.
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

// CountBodyDataTransforms calculates the number of data transforms applied to the body of the BackendResponse instance.
// If the body is omitted, it returns 1. Otherwise, it counts the number of keys in the body mapper,
// the number of keys in the body projection, and the number of body modifiers.
// The result is returned as an integer representing the total count.
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
