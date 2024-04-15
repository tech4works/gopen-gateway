package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"net/http"
	"time"
)

// Gopen is a struct that represents the configuration for a Gopen server.
type Gopen struct {
	// env is a string field that represents the environment in which the Gopen server is running.
	env string
	// version is a string field that represents the version of the Gopen server configured in the configuration json.
	version string
	// port represents the port number on which the Gopen application will listen for incoming requests.
	// It is an integer value and can be specified in the Gopen configuration JSON file.
	port int
	// hotReload represents a boolean flag indicating whether hot-reloading is enabled or not.
	// It is a field in the Gopen struct and is specified in the Gopen configuration JSON file.
	// It is used to control whether the Gopen application will automatically reload the configuration file
	// and apply the changes and restart the server.
	// If the value is true, hot-reloading is enabled. If the value is false, hot-reloading is disabled.
	// By default, hot-reloading is disabled, so if the field is not specified in the JSON file, it will be set to false.
	hotReload bool
	// timeout represents the timeout duration for a request or operation.
	// It is specified in string format and can be parsed into a time.Duration value.
	// The default value is empty. If not provided, the timeout will be 30s.
	timeout time.Duration
	// limiter represents the configuration for rate limiting.
	// It specifies the maximum header size, maximum body size, maximum multipart memory size, and the rate of allowed requests.
	limiter Limiter
	// cache is a struct representing the `cache` configuration in the Gopen struct. It contains the following fields:
	// - Duration: a string representing the duration of the `cache` in a format compatible with Go's time.ParseDuration
	// function. It defaults to an empty string. If not provided, the duration will be 30s.
	// - StrategyHeaders: a slice of strings representing the modifyHeaders used to determine the `cache` strategy. It defaults
	// to an empty slice.
	// - OnlyIfStatusCodes: A slice of integers representing the HTTP status codes for which the `cache` should be used.
	// Default is an empty slice. If not provided, the default value is 2xx success HTTP status codes
	// - OnlyIfMethods: a slice of strings representing the HTTP methods for which the `cache` should be used. The default
	// is an empty slice. If not provided by default, we only consider the http GET method.
	//- AllowCacheControl: a pointer to a boolean indicating whether the `cache` should honor the Cache-Control header.
	// It defaults to empty.
	cache Cache
	// securityCors represents the configuration options for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	securityCors SecurityCors
	// middlewares is a map that represents the middleware configuration in Gopen.
	// The keys of the map are the names of the middlewares, and the values are
	// Backend objects that define the properties of each middleware.
	// The Backend struct contains fields like name, hosts, path, method, forwardHeaders,
	// forwardQueries, modifiers, and extraConfig, which specify the behavior
	// and settings of the middleware.
	middlewares Middlewares
	// endpoints is a field in the Gopen struct that represents a slice of Endpoint objects.
	// Each Endpoint object defines a specific API endpoint with its corresponding settings such as path, method,
	// timeout, limiter, cache, etc.
	endpoints []Endpoint
}

// NewGOpen creates a new instance of Gopen based on the provided environment and gopenDTO.
// It initializes the fields of Gopen based on values from gopenDTO and sets default values for empty fields.
func NewGOpen(env string, gopenDTO dto.Gopen) Gopen {
	// damos o parse dos endpoints
	var endpoints []Endpoint
	for _, endpointDTO := range gopenDTO.Endpoints {
		endpoints = append(endpoints, newEndpoint(endpointDTO))
	}

	// damos o parse do timeout
	var timeout time.Duration
	var err error
	if helper.IsNotEmpty(gopenDTO.Timeout) {
		timeout, err = time.ParseDuration(gopenDTO.Timeout)
		if helper.IsNotNil(err) {
			logger.Warning("Parse duration timeout err:", err)
		}
	}

	return Gopen{
		env:          env,
		version:      gopenDTO.Version,
		port:         gopenDTO.Port,
		timeout:      timeout,
		limiter:      newLimiter(helper.IfNilReturns(gopenDTO.Limiter, dto.Limiter{})),
		cache:        newCache(helper.IfNilReturns(gopenDTO.Cache, dto.Cache{})),
		securityCors: newSecurityCors(helper.IfNilReturns(gopenDTO.SecurityCors, dto.SecurityCors{})),
		middlewares:  newMiddlewares(gopenDTO.Middlewares),
		endpoints:    endpoints,
	}
}

// Port returns the value of the port field in the Gopen struct.
func (g Gopen) Port() int {
	return g.port
}

// HotReload returns the value of the hotReload field in the Gopen struct.
func (g Gopen) HotReload() bool {
	return g.hotReload
}

// Version returns the value of the version field in the Gopen struct.
func (g Gopen) Version() string {
	return g.version
}

// Endpoints returns the value of the endpoints field in the Gopen struct.
func (g Gopen) Endpoints() []Endpoint {
	return g.endpoints
}

// Timeout returns the value of the timeout field in the Gopen struct. If the timeout is greater than 0,
// it returns the timeout value. Otherwise, it returns a default timeout of 30 seconds
func (g Gopen) Timeout() time.Duration {
	if helper.IsGreaterThan(g.timeout, 0) {
		return g.timeout
	}
	return 30 * time.Second
}

// LimiterMaxHeaderSize returns the value of maxHeaderSize field in the Limiter struct.
// If the value is greater than 0, it returns the value.
// Otherwise, it returns a newly created Bytes instance with the value "1MB".
func (g Gopen) LimiterMaxHeaderSize() Bytes {
	if helper.IsGreaterThan(g.limiter.maxHeaderSize, 0) {
		return g.limiter.maxHeaderSize
	}
	return NewBytes("1MB")
}

// LimiterMaxBodySize returns the maximum body size limit for the Gopen struct.
// It checks if the Limiter's maxBodySize field is greater than 0 and returns it.
// If it is not greater than 0, it returns a new Bytes object initialized with the value "3MB".
func (g Gopen) LimiterMaxBodySize() Bytes {
	if helper.IsGreaterThan(g.limiter.maxBodySize, 0) {
		return g.limiter.maxBodySize
	}
	return NewBytes("3MB")
}

// LimiterMaxMultipartMemorySize returns the value of the maxMultipartMemorySize field in the Limiter struct of the Gopen object.
// If the value is greater than zero, it returns the value of that field.
// Otherwise, it returns a new Bytes object initialized with a value of "5MB".
func (g Gopen) LimiterMaxMultipartMemorySize() Bytes {
	if helper.IsGreaterThan(g.limiter.maxMultipartMemorySize, 0) {
		return g.limiter.maxMultipartMemorySize
	}
	return NewBytes("5MB")
}

// LimiterRateCapacity returns the value of the capacity field in the LimiterRateCapacity method.
// If the capacity is greater than 0, it returns the capacity field value; otherwise, it returns 5.
// This method is used to retrieve the rate capacity for a limiter.
func (g Gopen) LimiterRateCapacity() int {
	if helper.IsGreaterThan(g.limiter.rate.capacity, 0) {
		return g.limiter.rate.capacity
	}
	return 5
}

// LimiterRateEvery returns the value of the rate.every field in the Limiter struct.
// If the value is greater than 0, it returns the rate.every value.
// Otherwise, it returns a default value of 1 second.
func (g Gopen) LimiterRateEvery() time.Duration {
	if helper.IsGreaterThan(g.limiter.rate.every, 0) {
		return g.limiter.rate.every
	}
	return time.Second
}

// CacheDuration returns the duration of the cache field in the Gopen struct.
func (g Gopen) CacheDuration() time.Duration {
	return g.cache.duration
}

// CacheStrategyHeaders returns the strategyHeaders field in the Cache struct, which represents the headers
// used in cache strategy.
func (g Gopen) CacheStrategyHeaders() []string {
	return g.cache.strategyHeaders
}

// CacheOnlyIfStatusCodes returns the array of status codes that are used for conditional caching.
// If the 'onlyIfStatusCodes' field in the Gopen cache is not empty, it returns that array.
// Otherwise, it returns a default array of commonly used status codes for conditional caching.
func (g Gopen) CacheOnlyIfStatusCodes() []int {
	if helper.IsNotEmpty(g.cache.onlyIfStatusCodes) {
		return g.cache.onlyIfStatusCodes
	}
	return []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent,
		http.StatusMultiStatus,
		http.StatusAlreadyReported,
		http.StatusIMUsed,
	}
}

// CacheOnlyIfMethods returns the array of HTTP methods used for conditional caching.
// If the 'onlyIfMethods' field in the Gopen cache is not empty, it returns that array.
// Otherwise, it returns a default array containing only the HTTP GET method.
func (g Gopen) CacheOnlyIfMethods() []string {
	if helper.IsNotEmpty(g.cache.onlyIfMethods) {
		return g.cache.onlyIfMethods
	}
	return []string{
		http.MethodGet,
	}
}

// AllowCacheControl checks if the caching is allowed or not.
// It uses the 'allowCacheControl' field in the 'Gopen' structure.
// In case of nil value, it defaults to 'false'.
func (g Gopen) AllowCacheControl() bool {
	return helper.IfNilReturns(g.cache.allowCacheControl, false)
}

// SecurityCors returns the value of the securityCors field in the Gopen struct.
func (g Gopen) SecurityCors() SecurityCors {
	return g.securityCors
}

// CountMiddlewares returns the number of middlewares in the Gopen instance.
func (g Gopen) CountMiddlewares() int {
	return len(g.middlewares)
}

// CountEndpoints returns the number of endpoints in the Gopen struct.
func (g Gopen) CountEndpoints() int {
	return len(g.endpoints)
}

// CountBackends returns the total number of backends present in the `Gopen` struct and its nested `Endpoint` structs.
// It calculates the count by summing the number of middlewares in `Gopen` and recursively iterating through each `Endpoint`
// to count their backends.
// Returns an integer indicating the total count of backends.
func (g Gopen) CountBackends() (count int) {
	count += g.CountMiddlewares()
	for _, endpointVO := range g.endpoints {
		count += endpointVO.CountBackends()
	}
	return count
}

// CountModifiers counts the total number of modifiers in the Gopen struct.
// It iterates through all the middleware backends and endpoint VOs,
// and calls the CountModifiers method on each of them to calculate the count.
// The count is incremented for each modifier found and the final count is returned.
func (g Gopen) CountModifiers() (count int) {
	for _, middlewareBackend := range g.middlewares {
		count += middlewareBackend.CountModifiers()
	}
	for _, endpointDTO := range g.endpoints {
		count += endpointDTO.CountModifiers()
	}
	return count
}

// Middleware retrieves a backend from the middlewares map based on the given key and returns it with a boolean
// indicating whether it exists or not. The returned backend is wrapped in a new middleware backend with the omitResponse
// field set to true.
func (g Gopen) Middleware(key string) (Backend, bool) {
	return g.middlewares.Get(key)
}

// Middlewares returns the value of the middlewares field in the Gopen struct.
func (g Gopen) Middlewares() Middlewares {
	return g.middlewares
}
