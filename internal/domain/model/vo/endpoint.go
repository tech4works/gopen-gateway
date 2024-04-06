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

type Endpoint struct {
	path               string
	method             string
	timeout            time.Duration
	limiter            Limiter
	cache              EndpointCache
	responseEncode     enum.ResponseEncode
	aggregateResponses bool
	abortIfStatusCodes []int
	beforeware         []string
	afterware          []string
	backends           []Backend
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
		path:               endpointDTO.Path,
		method:             endpointDTO.Method,
		timeout:            timeout,
		limiter:            newLimiter(helper.IfNilReturns(endpointDTO.Limiter, dto.Limiter{})),
		cache:              newEndpointCache(helper.IfNilReturns(endpointDTO.Cache, dto.EndpointCache{})),
		responseEncode:     endpointDTO.ResponseEncode,
		aggregateResponses: endpointDTO.AggregateResponses,
		abortIfStatusCodes: endpointDTO.AbortIfStatusCodes,
		beforeware:         endpointDTO.Beforeware,
		afterware:          endpointDTO.Afterware,
		backends:           backends,
	}
}

// Path returns the path field of the Endpoint struct.
func (e Endpoint) Path() string {
	return e.path
}

// Method returns the value of the method field in the Endpoint struct.
func (e Endpoint) Method() string {
	return e.method
}

// Equals checks if the given route is equal to the Endpoint's path and method.
// If the route is equal, it returns an error indicating a repeat route endpoint.
func (e Endpoint) Equals(route gin.RouteInfo) (err error) {
	if helper.Equals(route.Path, e.path) && helper.Equals(route.Method, e.method) {
		err = errors.New("Error path:", e.path, "method:", e.method, "repeat route endpoint")
	}
	return err
}

// HasTimeout returns true if the timeout field in the Endpoint struct is greater than 0, false otherwise.
func (e Endpoint) HasTimeout() bool {
	return helper.IsGreaterThan(e.timeout, 0)
}

// Timeout returns the value of the timeout field in the Endpoint struct.
func (e Endpoint) Timeout() time.Duration {
	return e.timeout
}

// HasLimiter returns true if the Endpoint has a Limiter set, otherwise false.
func (e Endpoint) HasLimiter() bool {
	return helper.IsNotEmpty(e.limiter)
}

// HasLimiterRateCapacity returns true if the limiter rate capacity has been set and is greater than 0.
func (e Endpoint) HasLimiterRateCapacity() bool {
	return helper.IsGreaterThan(e.limiter.rate.capacity, 0)
}

// LimiterRateCapacity returns the value of the capacity field in the rate field of the Limiter struct.
// It is used to determine the rate limit capacity for a specific endpoint.
func (e Endpoint) LimiterRateCapacity() int {
	return e.limiter.rate.capacity
}

// HasLimiterRateEvery returns true if the endpoint has a limiter rate every greater than 0, false otherwise.
func (e Endpoint) HasLimiterRateEvery() bool {
	return helper.IsGreaterThan(e.limiter.rate.every, 0)
}

// LimiterRateEvery returns the value of the 'every' field in the 'rate' field of the Limiter struct in the Endpoint struct.
func (e Endpoint) LimiterRateEvery() time.Duration {
	return e.limiter.rate.every
}

// HasLimiterMaxHeaderSize returns true if the limiter's maxHeaderSize is greater than 0, false otherwise.
func (e Endpoint) HasLimiterMaxHeaderSize() bool {
	return helper.IsGreaterThan(e.limiter.maxHeaderSize, 0)
}

// LimiterMaxHeaderSize returns the value of the maxHeaderSize field in the Limiter struct.
func (e Endpoint) LimiterMaxHeaderSize() Bytes {
	return e.limiter.maxHeaderSize
}

// HasLimiterMaxBodySize returns true if the `maxBodySize` field in the `Limiter` struct of the `Endpoint` is greater
// than 0, otherwise it returns false.
func (e Endpoint) HasLimiterMaxBodySize() bool {
	return helper.IsGreaterThan(e.limiter.maxBodySize, 0)
}

// LimiterMaxBodySize returns the value of the maxBodySize field in the Limiter struct of the Endpoint struct.
func (e Endpoint) LimiterMaxBodySize() Bytes {
	return e.limiter.maxBodySize
}

// HasLimiterMaxMultipartFormSize returns true if the limiter's maxMultipartMemorySize value is greater than 0.
// Otherwise, it returns false.
func (e Endpoint) HasLimiterMaxMultipartFormSize() bool {
	return helper.IsGreaterThan(e.limiter.maxMultipartMemorySize, 0)
}

// LimiterMaxMultipartMemorySize returns the value of the maxMultipartMemorySize field in the Limiter struct of the
// Endpoint struct.
func (e Endpoint) LimiterMaxMultipartMemorySize() Bytes {
	return e.limiter.maxMultipartMemorySize
}

// HasCache returns a boolean value indicating whether the Endpoint has a cache.
func (e Endpoint) HasCache() bool {
	return e.cache.enabled
}

// HasCacheDuration returns true if the cache duration of the Endpoint is greater than 0, otherwise false.
func (e Endpoint) HasCacheDuration() bool {
	return helper.IsGreaterThan(e.cache.duration, 0)
}

// CacheDuration returns the cache duration of the Endpoint.
func (e Endpoint) CacheDuration() time.Duration {
	return e.cache.duration
}

// HasCacheStrategyHeaders checks if the Endpoint has any cache strategy headers defined.
// It returns true if there are cache strategy headers, and false otherwise.
func (e Endpoint) HasCacheStrategyHeaders() bool {
	return helper.IsNotNil(e.cache.strategyHeaders)
}

// CacheStrategyHeaders returns the strategyHeaders field in the Cache struct of the Endpoint.
// It contains the headers that define the caching strategy for the endpoint.
func (e Endpoint) CacheStrategyHeaders() []string {
	return e.cache.strategyHeaders
}

// OnlyIfStatusCodes returns the onlyIfStatusCodes field of the cache struct.
func (e Endpoint) OnlyIfStatusCodes() []int {
	return e.cache.onlyIfStatusCodes
}

// HasAllowCacheControl returns a boolean value indicating whether the `allowCacheControl` field in the Cache struct
// of the Endpoint is not nil.
func (e Endpoint) HasAllowCacheControl() bool {
	return helper.IsNotNil(e.cache.allowCacheControl)
}

// AllowCacheControl returns the value of the allowCacheControl field in the Cache struct of the Endpoint.
// If the allowCacheControl field is nil, it returns false.
// This method is used to determine whether cache control is allowed for the endpoint.
func (e Endpoint) AllowCacheControl() bool {
	return helper.IfNilReturns(e.cache.allowCacheControl, false)
}

// CacheIgnoreQuery returns the value of the ignoreQuery field in the Cache struct,
// which determines whether to ignore the query parameters when caching a response.
func (e Endpoint) CacheIgnoreQuery() bool {
	return e.cache.ignoreQuery
}

// Beforeware returns the slice of strings representing the beforeware keys configured for the Endpoint.Beforeware
// middlewares are executed before the main backends.
func (e Endpoint) Beforeware() []string {
	return e.beforeware
}

// Backends returns the slice of backends in the Endpoint struct.
func (e Endpoint) Backends() []Backend {
	return e.backends
}

// Afterware returns the slice of strings representing the afterware keys configured for the Endpoint.Afterware
// middlewares are executed after the main backends.
func (e Endpoint) Afterware() []string {
	return e.afterware
}

// CountAllBackends calculates the total number of beforeware, backends, and afterware in the Endpoint struct.
// It returns the sum of the lengths of these slices.
func (e Endpoint) CountAllBackends() int {
	return len(e.beforeware) + len(e.backends) + len(e.afterware)
}

// CountBeforewares returns the number of beforewares in the Endpoint struct.
func (e Endpoint) CountBeforewares() int {
	return len(e.beforeware)
}

// CountAfterwares returns the number of afterwares in the Endpoint struct.
func (e Endpoint) CountAfterwares() int {
	return len(e.afterware)
}

// CountBackends returns the number of backends in the Endpoint struct.
func (e Endpoint) CountBackends() int {
	return len(e.backends)
}

// CountModifiers counts the total number of modifiers in an Endpoint by summing the count of modifiers in each
// Backend associated with it.
func (e Endpoint) CountModifiers() (count int) {
	for _, backendDTO := range e.backends {
		count += backendDTO.CountModifiers()
	}
	return count
}

// Completed checks if the response history size is equal to the count of all backends in the Endpoint struct.
func (e Endpoint) Completed(responseHistorySize int) bool {
	return helper.Equals(responseHistorySize, e.CountAllBackends())
}

// AbortSequencial checks if the response should be aborted based on the status codes defined in the abortIfStatusCodes
// field of the Endpoint struct. It returns true if the response status code matches any of the abortIfStatusCodes,
// otherwise it returns false. If the abortIfStatusCodes field is empty, it returns true if the response status code is
// greater than or equal to http.StatusBadRequest.
func (e Endpoint) AbortSequencial(responseVO Response) bool {
	if helper.IsEmpty(e.abortIfStatusCodes) {
		return helper.IsGreaterThanOrEqual(responseVO.statusCode, http.StatusBadRequest)
	}
	return helper.Contains(e.abortIfStatusCodes, responseVO.statusCode)
}

// ResponseEncode returns the value of the responseEncode field in the Endpoint struct.
func (e Endpoint) ResponseEncode() enum.ResponseEncode {
	return e.responseEncode
}

// AggregateResponses returns the value of the aggregateResponses field in the Endpoint struct.
func (e Endpoint) AggregateResponses() bool {
	return e.aggregateResponses
}

// AbortIfStatusCodes returns the value of the abortIfStatusCodes field in the Endpoint struct.
func (e Endpoint) AbortIfStatusCodes() []int {
	return e.abortIfStatusCodes
}

// Resume returns a string representation of the Endpoint, including information about the method,
// path, the number of beforewares, afterwares, backends, and modifiers.
// The format of the string is as follows:
// "{method} -> \"{path}\" [beforeware: {CountBeforewares} afterware: {CountAfterwares} backends: {CountBackends} modifiers: {CountModifiers}]"
func (e Endpoint) Resume() string {
	return fmt.Sprintf("%s -> \"%s\" [beforeware: %v afterware: %v backends: %v modifiers: %v]", e.method, e.path,
		e.CountBeforewares(), e.CountAfterwares(), e.CountBackends(), e.CountModifiers())
}
