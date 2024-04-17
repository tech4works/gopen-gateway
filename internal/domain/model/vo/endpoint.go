package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Endpoint represents the configuration for an API endpoint in the Gopen application.
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
	timeout time.Duration
	// limiter represents the configuration for rate limiting in the Gopen application.
	// The default value is nil. If not provided, the `limiter` will be Gopen.limiter.
	limiter *EndpointLimiter
	// cache represents the `cache` configuration for an endpoint.
	// The default value is EndpointCache empty with enabled false.
	cache *EndpointCache
	// responseEncode represents the encoding format for the API endpoint response. The ResponseEncode
	// field is an enum.ResponseEncode value, which can have one of the following values:
	// - enum.ResponseEncodeText: for encoding the response as plain text.
	// - enum.ResponseEncodeJson: for encoding the response as JSON.
	// - enum.ResponseEncodeXml: for encoding the response as XML.
	// The default value is empty. If not provided, the response will be encoded by type, if the string is json it
	// returns json, otherwise it responds to plain text
	responseEncode enum.ResponseEncode
	// aggregateResponses represents a boolean indicating whether the API endpoint should aggregate responses
	// from multiple backends.
	aggregateResponses bool
	// abortIfStatusCodes represents a slice of integers representing the HTTP status codes
	// for which the API endpoint should abort. It is a field in the Endpoint struct.
	abortIfStatusCodes *[]int
	// beforeware represents a slice of strings containing the names of the beforeware middlewares that should be
	// applied before processing the API endpoint.
	beforeware []string
	// afterware represents the configuration for the afterware middlewares to apply after processing the API endpoint.
	// It is a slice of strings representing the names of the afterware middlewares to apply.
	// The names specify the behavior and settings of each afterware middleware.
	// If not provided, the default value is an empty slice.
	// The afterware middleware is executed after processing the API endpoint, allowing for modification or
	// transformation of the response or performing any additional actions.
	// Afterware can be used for logging, error handling, response modification, etc.
	afterware []string
	// Backends represents the backend configurations for an API endpoint in the Gopen application.
	// It is a slice of Backend structs.
	backends []Backend
}

// newEndpoint creates a new instance of Endpoint based on the provided endpointDTO.
// It initializes the fields of Endpoint based on values from endpointDTO and sets default values for empty fields.
// The function returns the created Endpoint.
func newEndpoint(endpointDTO dto.Endpoint) Endpoint {
	var backends []Backend
	for _, backendDTO := range endpointDTO.Backends {
		backends = append(backends, newBackend(backendDTO))
	}

	var timeout time.Duration
	var err error
	if helper.IsNotEmpty(endpointDTO.Timeout) {
		timeout, err = time.ParseDuration(endpointDTO.Timeout)
		if helper.IsNotNil(err) {
			logger.Warning("Parse duration endpoint.timeout err:", err)
		}
	}

	return Endpoint{
		comment:            endpointDTO.Comment,
		path:               endpointDTO.Path,
		method:             endpointDTO.Method,
		timeout:            timeout,
		limiter:            newEndpointLimiterFromDTO(endpointDTO.Limiter),
		cache:              newEndpointCacheFromDTO(endpointDTO.Cache),
		responseEncode:     endpointDTO.ResponseEncode,
		aggregateResponses: endpointDTO.AggregateResponses,
		abortIfStatusCodes: endpointDTO.AbortIfStatusCodes,
		beforeware:         endpointDTO.Beforeware,
		afterware:          endpointDTO.Afterware,
		backends:           backends,
	}
}

// fillDefaultValues sets default values for an Endpoint object based on a given Gopen object.
// The timeout value is obtained from the Gopen object by default, unless a timeout value is specified in the Endpoint,
// in which case, that value takes priority. The limiter and cache values are constructed using the global configuration
// from the Gopen object and the Endpoint object. The method returns a new Endpoint object with the default values set.
func (e *Endpoint) fillDefaultValues(gopenVO *Gopen) Endpoint {
	// por padrão obtemos o timeout configurado na raiz, caso não informado um valor padrão é retornado
	timeoutDuration := gopenVO.Timeout()
	// se o timeout foi informado no endpoint damos prioridade a ele
	if e.HasTimeout() {
		timeoutDuration = e.Timeout()
	}
	// construímos o limiter com os valores de configuração global
	endpointLimiterVO := newEndpointLimiter(gopenVO.Limiter(), e.Limiter())

	// construímos o endpoint cache com os valores de configuração global
	endpointCacheVO := newEndpointCache(gopenVO.Cache(), e.Cache())

	// construímos o VO com os valores padrões construídos a partir do Gopen e o próprio endpoint
	return Endpoint{
		path:               e.path,
		method:             e.method,
		timeout:            timeoutDuration,
		limiter:            endpointLimiterVO,
		cache:              endpointCacheVO,
		responseEncode:     e.responseEncode,
		aggregateResponses: e.aggregateResponses,
		abortIfStatusCodes: e.abortIfStatusCodes,
		beforeware:         e.beforeware,
		afterware:          e.afterware,
		backends:           e.backends,
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

// Equals checks if the given route is equal to the Endpoint's path and method.
// If the route is equal, it returns an error indicating a repeat route endpoint.
func (e *Endpoint) Equals(route gin.RouteInfo) (err error) {
	if helper.Equals(route.Path, e.path) && helper.Equals(route.Method, e.method) {
		err = errors.New("Error path:", e.path, "method:", e.method, "repeat route endpoint")
	}
	return err
}

// HasTimeout returns true if the timeout field in the Endpoint struct is greater than 0, false otherwise.
func (e *Endpoint) HasTimeout() bool {
	return helper.IsGreaterThan(e.timeout, 0)
}

// Timeout returns the value of the timeout field in the Endpoint struct.
func (e *Endpoint) Timeout() time.Duration {
	return e.timeout
}

// TimeoutStr returns the string representation of the timeout value in the Endpoint struct.
// If the Endpoint has a timeout value greater than 0, it returns the string representation of the timeout value.
// Otherwise, an empty string is returned.
func (e *Endpoint) TimeoutStr() string {
	if e.HasTimeout() {
		return e.timeout.String()
	}
	return ""
}

// Limiter returns the limiter field of the Endpoint struct.
func (e *Endpoint) Limiter() *EndpointLimiter {
	return e.limiter
}

// HasLimiter returns true if the Endpoint has a Limiter set, otherwise false.
func (e *Endpoint) HasLimiter() bool {
	return helper.IsNotEmpty(e.limiter)
}

// Cache returns the cache field of the Endpoint struct.
func (e *Endpoint) Cache() *EndpointCache {
	return e.cache
}

// HasCache returns a boolean value indicating whether the Endpoint has a cache.
func (e *Endpoint) HasCache() bool {
	return helper.IsNotNil(e.cache) && e.cache.enabled
}

// Beforeware returns the slice of strings representing the beforeware keys configured for the Endpoint.Beforeware
// middlewares are executed before the main backends.
func (e *Endpoint) Beforeware() []string {
	return e.beforeware
}

// Backends returns the slice of backends in the Endpoint struct.
func (e *Endpoint) Backends() []Backend {
	return e.backends
}

// Afterware returns the slice of strings representing the afterware keys configured for the Endpoint.Afterware
// middlewares are executed after the main backends.
func (e *Endpoint) Afterware() []string {
	return e.afterware
}

// CountAllBackends calculates the total number of beforeware, backends, and afterware in the Endpoint struct.
// It returns the sum of the lengths of these slices.
func (e *Endpoint) CountAllBackends() int {
	return e.CountBeforewares() + e.CountBackends() + e.CountAfterwares()
}

// CountBeforewares returns the number of beforewares in the Endpoint struct.
func (e *Endpoint) CountBeforewares() int {
	if helper.IsNil(e.Beforeware()) {
		return 0
	}
	return len(e.Beforeware())
}

// CountAfterwares returns the number of afterwares in the Endpoint struct.
func (e *Endpoint) CountAfterwares() int {
	if helper.IsNil(e.Afterware()) {
		return 0
	}
	return len(e.Afterware())
}

// CountBackends returns the number of backends in the Endpoint struct.
func (e *Endpoint) CountBackends() int {
	if helper.IsNil(e.Backends()) {
		return 0
	}
	return len(e.Backends())
}

// CountModifiers counts the total number of modifiers in an Endpoint by summing the count of modifiers in each
// Backend associated with it.
func (e *Endpoint) CountModifiers() (count int) {
	for _, backendDTO := range e.backends {
		count += backendDTO.CountModifiers()
	}
	return count
}

// Completed checks if the response history size is equal to the count of all backends in the Endpoint struct.
func (e *Endpoint) Completed(responseHistorySize int) bool {
	return helper.Equals(responseHistorySize, e.CountAllBackends())
}

// AbortSequencial checks if the given statusCode is present in the abortIfStatusCodes
// slice of the Endpoint struct. If the abortIfStatusCodes slice is nil, it returns
// true if the statusCode is greater than or equal to http.StatusBadRequest, otherwise false.
// Otherwise, it returns true if the given statusCode is present in the abortIfStatusCodes
// slice, otherwise false.
func (e *Endpoint) AbortSequencial(statusCode int) bool {
	if helper.IsNil(e.abortIfStatusCodes) {
		return helper.IsGreaterThanOrEqual(statusCode, http.StatusBadRequest)
	}
	return helper.Contains(e.abortIfStatusCodes, statusCode)
}

// ResponseEncode returns the value of the responseEncode field in the Endpoint struct.
func (e *Endpoint) ResponseEncode() enum.ResponseEncode {
	return e.responseEncode
}

// AggregateResponses returns the value of the aggregateResponses field in the Endpoint struct.
func (e *Endpoint) AggregateResponses() bool {
	return e.aggregateResponses
}

// AbortIfStatusCodes returns the value of the abortIfStatusCodes field in the Endpoint struct.
func (e *Endpoint) AbortIfStatusCodes() *[]int {
	return e.abortIfStatusCodes
}

// Resume returns a string representation of the Endpoint, including information about the method,
// path, the number of beforewares, afterwares, backends, and modifiers.
// The format of the string is as follows:
// "{method} -> \"{path}\" [beforeware: {CountBeforewares} afterware: {CountAfterwares} backends: {CountBackends} modifiers: {CountModifiers}]"
func (e *Endpoint) Resume() string {
	return fmt.Sprintf("%s --> \"%s\" [beforeware: %v afterware: %v backends: %v modifiers: %v]", e.method, e.path,
		e.CountBeforewares(), e.CountAfterwares(), e.CountBackends(), e.CountModifiers())
}
