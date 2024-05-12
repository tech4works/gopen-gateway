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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"time"
)

// Endpoint represent the configuration for an API endpoint in the Gopen application.
type Endpoint struct {
	// comment is a string field representing the comment associated with an API endpoint.
	comment string
	// path is a string representing the path of the API endpoint. It is a field in the Endpoint struct.
	path string
	// method represents the HTTP method of an API endpoint.
	method string
	// timeout represents the timeout duration for the API endpoint.
	// It is a string value specified in the JSON configuration.
	// The default value is empty. If not provided, the timeout will be Gopen.timeout.
	timeout Duration
	// limiter represents the configuration for rate limiting in the Gopen application.
	// The default value is nil. If not provided, the `limiter` will be Gopen.limiter.
	limiter *Limiter
	// cache represents the `cache` configuration for an endpoint.
	// The default value is Cache empty with enabled false.
	cache *Cache
	// abortIfStatusCodes represents a slice of integers representing the HTTP status codes
	// for which the API endpoint should abort. It is a field in the Endpoint struct.
	abortIfStatusCodes *[]int
	// response represents the configuration for the response of an API endpoint in the Gopen application.
	response *EndpointResponse
	// beforewares represents a slice of strings containing the names of the beforewares middlewares that should be
	// applied before processing the API endpoint.
	beforewares []string
	// afterwares represents the configuration for the afterwares middlewares to apply after processing the API endpoint.
	// It is a slice of strings representing the names of the afterwares middlewares to apply.
	// The names specify the behavior and settings of each afterwares middleware.
	// If not provided, the default value is an empty slice.
	// The afterwares middleware is executed after processing the API endpoint, allowing for modification or
	// transformation of the response or performing any additional actions.
	// Afterwares can be used for logging, error handling, response modification, etc.
	afterwares []string
	// Backends represents the backend configurations for an API endpoint in the Gopen application.
	// It is a slice of Backend structs.
	backends []Backend
}

// EndpointResponse represents the configuration for the response of an API endpoint in the Gopen application.
type EndpointResponse struct {
	// aggregate represents a boolean indicating whether the API endpoint should aggregate responses
	// from multiple backends.
	aggregate bool
	// contentType represents the encoding format for the API endpoint httpResponse. The ResponseEncode
	// field is an enum.ContentType value, which can have one of the following values:
	// - enum.ContentTypePlainText: for encoding the httpResponse as plain text.
	// - enum.ContentTypeJson: for encoding the httpResponse as JSON.
	// - enum.ContentTypeXml: for encoding the httpResponse as XML.
	// The default value is empty. If not provided, the httpResponse will be encoded by type, if the string is json it
	// returns json, otherwise it responds to plain text.
	contentType     enum.ContentType
	contentEncoding enum.ContentEncoding
	// nomenclature represents the case format for json text fields.
	nomenclature enum.Nomenclature
	// omitEmpty represents a boolean indicating whether empty fields should be omitted in the API endpoint response.
	omitEmpty bool
}

// NewEndpointStatic creates a new instance of the Endpoint struct with the provided path and method.
// It sets the timeout to 10 seconds and creates a default limiter using newLimiterDefault().
//
// Parameters:
// - path: The path of the endpoint.
// - method: The HTTP method of the endpoint.
//
// Returns:
// - The created instance of the Endpoint struct.
func NewEndpointStatic(path, method string) Endpoint {
	return Endpoint{
		path:    path,
		method:  method,
		timeout: NewDuration(10 * time.Second),
		limiter: newLimiterDefault(),
	}
}

// newEndpoint creates a new instance of the Endpoint struct with the provided GopenJson and EndpointJson.
// It initializes the fields of the Endpoint struct using values from the GopenJson and EndpointJson objects.
// The timeout value is obtained from the GopenJson, but if it is provided in the EndpointJson, it takes priority.
// The limiter and cache fields are initialized using the global and local configuration values from GopenJson and EndpointJson.
// The backends slice is populated by iterating over the backendJson objects in the EndpointJson and creating a new backend for each.
// The created instance of the Endpoint struct is returned.
func newEndpoint(gopenJson *GopenJson, endpointJson *EndpointJson) Endpoint {
	timeout := gopenJson.Timeout
	if helper.IsGreaterThan(endpointJson.Timeout, 0) {
		timeout = endpointJson.Timeout
	}
	limiter := newLimiter(gopenJson.Limiter, endpointJson.Limiter)
	cache := newCache(gopenJson.Cache, endpointJson.Cache)

	var backends []Backend
	for _, backendJson := range endpointJson.Backends {
		backends = append(backends, newBackend(&backendJson))
	}

	return Endpoint{
		comment:            endpointJson.Comment,
		path:               endpointJson.Path,
		method:             endpointJson.Method,
		timeout:            timeout,
		limiter:            limiter,
		cache:              cache,
		abortIfStatusCodes: endpointJson.AbortIfStatusCodes,
		response:           newEndpointResponse(endpointJson.Response),
		beforewares:        endpointJson.Beforewares,
		afterwares:         endpointJson.Afterwares,
		backends:           backends,
	}
}

// newEndpointResponse creates a new instance of the EndpointResponse struct using the values from the provided
// EndpointResponseJson object.
//
// Parameters:
//
// - endpointResponseJson: A pointer to an EndpointResponseJson object containing the configuration for the endpoint
// response.
//
// Returns:
//
// - A pointer to the created instance of the EndpointResponse struct.
func newEndpointResponse(endpointResponseJson *EndpointResponseJson) *EndpointResponse {
	if helper.IsNil(endpointResponseJson) {
		return nil
	}
	return &EndpointResponse{
		aggregate:    endpointResponseJson.Aggregate,
		contentType:  endpointResponseJson.ContentType,
		nomenclature: endpointResponseJson.Nomenclature,
		omitEmpty:    endpointResponseJson.OmitEmpty,
	}
}

// Path returns the path field of the Endpoint struct.
func (e *Endpoint) Path() string {
	return e.path
}

// Method returns the value of the method field in the Endpoint struct.
func (e *Endpoint) Method() string {
	return e.method
}

// Timeout returns the timeout duration configured for the Endpoint.
// It checks if the configured timeout is greater than 0. If it is, it returns
// the configured timeout. Otherwise, it returns a default timeout of 30 seconds.
func (e *Endpoint) Timeout() Duration {
	if helper.IsGreaterThan(e.timeout, 0) {
		return e.timeout
	}
	return NewDuration(30 * time.Second)
}

// Limiter returns the limiter field of the Endpoint struct.
func (e *Endpoint) Limiter() *Limiter {
	return e.limiter
}

// Cache returns the cache field of the Endpoint struct.
func (e *Endpoint) Cache() *Cache {
	return e.cache
}

// Beforewares returns the slice of strings representing the beforeware keys configured for the Endpoint.Beforewares
// middlewares are executed before the main backends.
func (e *Endpoint) Beforewares() []string {
	return e.beforewares
}

// Afterwares returns the slice of strings representing the afterware keys configured for the Endpoint.Afterwares
// middlewares are executed after the main backends.
func (e *Endpoint) Afterwares() []string {
	return e.afterwares
}

// Backends returns the slice of backends in the Endpoint struct.
func (e *Endpoint) Backends() []Backend {
	return e.backends
}

// CountBeforewares returns the number of beforewares in the Endpoint struct.
func (e *Endpoint) CountBeforewares() int {
	if helper.IsNil(e.Beforewares()) {
		return 0
	}
	return len(e.Beforewares())
}

// CountAfterwares returns the number of afterwares in the Endpoint struct.
func (e *Endpoint) CountAfterwares() int {
	if helper.IsNil(e.Afterwares()) {
		return 0
	}
	return len(e.Afterwares())
}

// CountBackends returns the number of backends in the Endpoint struct.
func (e *Endpoint) CountBackends() int {
	return len(e.Backends())
}

// CountBackendsNonOmit returns the number of backends in the Endpoint struct that have a non-nil response and are not omitted.
func (e *Endpoint) CountBackendsNonOmit() int {
	count := 0
	for _, backend := range e.Backends() {
		if helper.IsNil(backend.Response()) || !backend.Response().Omit() {
			count++
		}
	}
	return count
}

// CountAllDataTransforms recursively counts the number of data transforms in the Endpoint struct
// by adding the count of transforms in the Endpoint's response and the count of transforms in each backend.
// It returns the total count of transforms.
func (e *Endpoint) CountAllDataTransforms() (count int) {
	if helper.IsNotNil(e.Response()) {
		count += e.Response().CountAllDataTransforms()
	}
	for _, backend := range e.backends {
		count += backend.CountAllDataTransforms()
	}
	return count
}

// Completed checks if the given responseHistorySize is equal to the count of non-omitted backends in the endpoint.
func (e *Endpoint) Completed(responseHistorySize int) bool {
	return helper.Equals(responseHistorySize, e.CountBackendsNonOmit())
}

// Abort checks if the given StatusCode should cause the endpoint to be aborted.
// It returns true if the StatusCode is present in the abortIfStatusCodes slice.
// If the abortIfStatusCodes slice is nil, it returns whether the StatusCode is a Failed status.
func (e *Endpoint) Abort(statusCode StatusCode) bool {
	if helper.IsNil(e.abortIfStatusCodes) {
		return statusCode.Failed()
	}
	return helper.Contains(e.abortIfStatusCodes, statusCode)
}

// Response returns the response field of the Endpoint struct.
func (e *Endpoint) Response() *EndpointResponse {
	return e.response
}

// AbortIfStatusCodes returns the value of the abortIfStatusCodes field in the Endpoint struct.
func (e *Endpoint) AbortIfStatusCodes() *[]int {
	return e.abortIfStatusCodes
}

// Resume returns a string representation of the Endpoint's method, path, beforeware count, afterware count,
// backend count, and data transformation count. It formats the string as follows:
// "{Method} --> "{Path}" (beforeware:{BeforewareCount}, afterware:{AfterwareCount},
// backends:{BackendCount}, transformations:{DataTransformationCount})"
func (e *Endpoint) Resume() string {
	return fmt.Sprintf("%s --> \"%s\" (beforeware:%v, afterware:%v, backends:%v, transformations:%v)",
		e.method, e.path, e.CountBeforewares(), e.CountAfterwares(), e.CountBackends(), e.CountAllDataTransforms())
}

// NoCache returns true if the cache field of the Endpoint struct is nil or disabled.
func (e *Endpoint) NoCache() bool {
	return helper.IsNil(e.Cache()) || e.Cache().Disabled()
}

// HasContentType checks if the encode field of EndpointResponse is a valid enumeration value.
// It returns true if the encode is either ContentTypePlainText, EncodeJson, or ContentTypeXml, otherwise it returns false.
func (e EndpointResponse) HasContentType() bool {
	return e.contentType.IsEnumValid()
}

func (e EndpointResponse) HasContentEncoding() bool {
	return e.contentEncoding.IsEnumValid()
}

// ContentType returns the ContentType based on the encoding format specified in the EndpointResponse.
//
// The function performs a switch on the encode field of the EndpointResponse struct and returns
// different ContentType objects based on the value of encode:
// - If encode is enum.ContentTypeJson, it returns NewContentTypeJson().
// - If encode is enum.ContentTypeXml, it returns NewContentTypeXml().
// - For any other value of encode, it returns NewContentTypeTextPlain() as the default.
//
// The returned ContentType represents the encoding format for the HttpResponse.
// Example usage:
//
//	endpointResponse := EndpointResponse{encode: enum.EncodeJson}
//	contentType := endpointResponse.ContentType() // Returns NewContentTypeJson()
func (e EndpointResponse) ContentType() ContentType {
	switch e.contentType {
	default:
		return NewContentTypeTextPlain()
	case enum.ContentTypeJson:
		return NewContentTypeJson()
	case enum.ContentTypeXml:
		return NewContentTypeXml()
	}
}

func (e EndpointResponse) ContentEncoding() ContentEncoding {
	switch e.contentEncoding {
	default:
		return ""
	case enum.ContentEncodingGzip:
		return NewContentEncodingGzip()
	case enum.ContentEncodingDeflate:
		return NewContentEncodingDeflate()
	}
}

// Aggregate returns the value of the aggregate field in the EndpointResponse struct.
// The aggregate field represents a boolean indicating whether the API endpoint should aggregate responses
// from multiple backends.
func (e EndpointResponse) Aggregate() bool {
	return e.aggregate
}

// OmitEmpty returns the value of the omitEmpty field in the EndpointResponse struct.
// The omitEmpty field represents a boolean indicating whether empty fields should be omitted in the API endpoint
// response.
func (e EndpointResponse) OmitEmpty() bool {
	return e.omitEmpty
}

// HasNomenclature checks if the nomenclature field of EndpointResponse is a valid enumeration value.
// It returns true if the nomenclature is either NomenclatureCamel, NomenclatureLowerCamel, NomenclatureSnake,
// or NomenclatureKebab, otherwise it returns false.
func (e EndpointResponse) HasNomenclature() bool {
	return e.nomenclature.IsEnumValid()
}

// Nomenclature returns the value of the nomenclature field in the EndpointResponse struct.
// The nomenclature field represents the case format for json text fields.
func (e EndpointResponse) Nomenclature() enum.Nomenclature {
	return e.nomenclature
}

// CountAllDataTransforms returns the number of data transforms applied in the EndpointResponse struct.
// It increments the count by 1 for each data transform:
//   - Aggregate: if the API endpoint should aggregate responses from multiple backends.
//   - OmitEmpty: if empty fields should be omitted in the API endpoint response.
//   - HasContentType: if the encode field is a valid enumeration value (ContentTypePlainText, EncodeJson, ContentTypeXml).
//   - HasNomenclature: if the nomenclature field is a valid enumeration value (NomenclatureCamel, NomenclatureLowerCamel,
//     NomenclatureSnake, NomenclatureKebab).
//
// The count is returned as the result.
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
