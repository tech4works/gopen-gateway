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
	configEnum "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/enum"
	"net/http"
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
	limiter *EndpointLimiter
	// cache represents the `cache` configuration for an endpoint.
	// The default value is EndpointCache empty with enabled false.
	cache *EndpointCache
	// abortIfStatusCodes represents a slice of integers representing the HTTP status codes
	// for which the API endpoint should abort. It is a field in the Endpoint struct.
	abortIfStatusCodes *[]int
	// todo:
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

type EndpointResponse struct {
	// aggregate represents a boolean indicating whether the API endpoint should aggregate responses
	// from multiple backends.
	aggregate bool
	// encode represents the encoding format for the API endpoint httpResponse. The ResponseEncode
	// field is an enum.Encode value, which can have one of the following values:
	// - enum.EncodeText: for encoding the httpResponse as plain text.
	// - enum.EncodeJson: for encoding the httpResponse as JSON.
	// - enum.EncodeXml: for encoding the httpResponse as XML.
	// The default value is empty. If not provided, the httpResponse will be encoded by type, if the string is json it
	// returns json, otherwise it responds to plain text
	encode       configEnum.Encode
	nomenclature configEnum.Nomenclature
	omitEmpty    bool
}

func NewEndpointStatic(gopenVO *Gopen, path, method string) Endpoint {
	return Endpoint{
		path:    path,
		method:  method,
		timeout: Duration(10 * time.Second),
		limiter: newEndpointLimiterStatic(gopenVO.Limiter()),
	}
}

// todo:
func newEndpoint(gopen *Gopen, endpointJson *EndpointJson) Endpoint {
	// por padrão obtemos o timeout configurado na raiz, caso não informado um valor padrão é retornado
	timeoutDuration := gopen.Timeout()
	// se o timeout foi informado no endpoint damos prioridade a ele
	if helper.IsGreaterThan(endpointJson.Timeout, 0) {
		timeoutDuration = endpointJson.Timeout
	}
	// construímos o limiter com os valores de configuração global
	endpointLimiter := newEndpointLimiter(gopen.Limiter(), endpointJson.Limiter)

	// construímos o endpoint cache com os valores de configuração global
	endpointCache := newEndpointCache(gopen.Cache(), endpointJson.Cache)

	// fazemos o parse dos backends
	var backends []Backend
	for _, backendJsonVO := range endpointJson.Backends {
		backends = append(backends, newBackend(&backendJsonVO))
	}

	// construímos o endpoint VO ja com os valores padrões
	return Endpoint{
		comment:            endpointJson.Comment,
		path:               endpointJson.Path,
		method:             endpointJson.Method,
		timeout:            timeoutDuration,
		limiter:            endpointLimiter,
		cache:              endpointCache,
		abortIfStatusCodes: endpointJson.AbortIfStatusCodes,
		response:           newEndpointResponse(endpointJson.Response),
		beforewares:        endpointJson.Beforewares,
		afterwares:         endpointJson.Afterwares,
		backends:           backends,
	}
}

// todo:
func newEndpointResponse(endpointResponseJson *EndpointResponseJson) *EndpointResponse {
	if helper.IsNil(endpointResponseJson) {
		return nil
	}
	return &EndpointResponse{
		aggregate:    endpointResponseJson.Aggregate,
		encode:       endpointResponseJson.Encode,
		nomenclature: endpointResponseJson.Nomenclature,
		omitEmpty:    endpointResponseJson.OmitEmpty,
	}
}

// Comment returns the comment field of the Endpoint struct.
func (e *Endpoint) Comment() string {
	return e.comment
}

// Path returns the path field of the Endpoint struct.
func (e *Endpoint) Path() string {
	return e.path
}

// Method returns the value of the method field in the Endpoint struct.
func (e *Endpoint) Method() string {
	return e.method
}

// Timeout returns the value of the timeout field in the Endpoint struct.
func (e *Endpoint) Timeout() Duration {
	return e.timeout
}

// Limiter returns the limiter field of the Endpoint struct.
func (e *Endpoint) Limiter() *EndpointLimiter {
	return e.limiter
}

// Cache returns the cache field of the Endpoint struct.
func (e *Endpoint) Cache() *EndpointCache {
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

// CountAllBackends calculates the total number of beforeware, backends, and afterware in the Endpoint struct.
// It returns the sum of the lengths of these slices.
func (e *Endpoint) CountAllBackends() int {
	return e.CountBeforewares() + e.CountBackends() + e.CountAfterwares()
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

func (e *Endpoint) CountBackendsNonOmit() int {
	count := 0
	for _, backend := range e.Backends() {
		if helper.IsNil(backend.Response()) || !backend.Response().Omit() {
			count++
		}
	}
	return count
}

func (e *Endpoint) CountAllDataTransforms() (count int) {
	if helper.IsNotNil(e.Response()) {
		count += e.Response().CountAllDataTransforms()
	}
	for _, backend := range e.backends {
		count += backend.CountAllDataTransforms()
	}
	return count
}

func (e *Endpoint) Completed(responseHistorySize int) bool {
	return helper.Equals(responseHistorySize, e.CountBackendsNonOmit())
}

// Abort checks if the given statusCode is present in the abortIfStatusCodes
// slice of the Endpoint struct. If the abortIfStatusCodes slice is nil, it returns
// true if the statusCode is greater than or equal to http.StatusBadRequest, otherwise false.
// Otherwise, it returns true if the given statusCode is present in the abortIfStatusCodes
// slice, otherwise false.
func (e *Endpoint) Abort(statusCode int) bool {
	if helper.IsNil(e.abortIfStatusCodes) {
		return helper.IsGreaterThanOrEqual(statusCode, http.StatusBadRequest)
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

func (e *Endpoint) Resume() string {
	return fmt.Sprintf("%s --> \"%s\" (beforeware:%v, afterware:%v, backends:%v, transformations:%v)",
		e.method, e.path, e.CountBeforewares(), e.CountAfterwares(), e.CountBackends(), e.CountAllDataTransforms())
}

func (r EndpointResponse) HasEncode() bool {
	return r.encode.IsEnumValid()
}

func (r EndpointResponse) Encode() configEnum.Encode {
	return r.encode
}

func (r EndpointResponse) Aggregate() bool {
	return r.aggregate
}

func (r EndpointResponse) OmitEmpty() bool {
	return r.omitEmpty
}

func (r EndpointResponse) HasNomenclature() bool {
	return r.nomenclature.IsEnumValid()
}

func (r EndpointResponse) Nomenclature() configEnum.Nomenclature {
	return r.nomenclature
}

func (r EndpointResponse) CountAllDataTransforms() (count int) {
	if r.Aggregate() {
		count++
	}
	if r.OmitEmpty() {
		count++
	}
	if r.HasEncode() {
		count++
	}
	if r.HasNomenclature() {
		count++
	}
	return count
}
