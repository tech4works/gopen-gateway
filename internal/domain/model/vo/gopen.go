package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"time"
)

type Gopen struct {
	env          string
	version      string
	hotReload    bool
	port         int
	timeout      time.Duration
	limiter      Limiter
	cache        Cache
	securityCors SecurityCors
	middlewares  Middlewares
	endpoints    []Endpoint
}

// NewGOpen creates a new instance of Gopen based on the provided environment and gopenDTO.
// It initializes the fields of Gopen based on values from gopenDTO and sets default values for empty fields.
func NewGOpen(env string, gopenDTO dto.Gopen) Gopen {
	var endpoints []Endpoint
	for _, endpointDTO := range gopenDTO.Endpoints {
		endpoints = append(endpoints, newEndpoint(endpointDTO))
	}

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

// NewCacheFromEndpoint creates a new instance of Cache based on the provided duration, strategyHeaders, and allowCacheControl.
// It initializes the fields of Cache with the given values.
func NewCacheFromEndpoint(duration time.Duration, strategyHeaders []string, allowCacheControl bool) Cache {
	return Cache{
		duration:          duration,
		strategyHeaders:   strategyHeaders,
		allowCacheControl: &allowCacheControl,
	}
}

// newLimiter creates a new instance of Limiter based on the provided limiterDTO.
// It initializes the fields of Limiter based on values from limiterDTO and sets default values for empty fields.
func newLimiter(limiterDTO dto.Limiter) Limiter {
	return Limiter{
		maxHeaderSize:          NewBytes(limiterDTO.MaxHeaderSize),
		maxBodySize:            NewBytes(limiterDTO.MaxBodySize),
		maxMultipartMemorySize: NewBytes(limiterDTO.MaxMultipartMemorySize),
		rate:                   newRate(helper.IfNilReturns(limiterDTO.Rate, dto.Rate{})),
	}
}

// newRate creates a new instance of Rate based on the provided rateDTO.
// It initializes the fields of Rate based on values from rateDTO and sets default values for empty fields.
func newRate(rateDTO dto.Rate) Rate {
	var every time.Duration
	var err error
	if helper.IsNotEmpty(rateDTO.Every) {
		every, err = time.ParseDuration(rateDTO.Every)
		if helper.IsNotNil(err) {
			logger.Warning("Parse duration limiter.rate.every err:", err)
		}
	}

	return Rate{
		capacity: rateDTO.Capacity,
		every:    every,
	}
}

// newCache creates a new instance of Cache based on the provided cacheDTO.
// It initializes the fields of Cache based on values from cacheDTO and sets default values for empty fields.
func newCache(cacheDTO dto.Cache) Cache {
	var duration time.Duration
	var err error
	if helper.IsNotEmpty(cacheDTO.Duration) {
		duration, err = time.ParseDuration(cacheDTO.Duration)
		if helper.IsNotNil(err) {
			logger.Warning("Parse duration cache.duration err:", err)
		}
	}

	return Cache{
		duration:          duration,
		strategyHeaders:   cacheDTO.StrategyHeaders,
		allowCacheControl: cacheDTO.AllowCacheControl,
	}
}

// newSecurityCors creates a new instance of SecurityCors based on the provided securityCorsDTO.
// It sets the allowOrigins, allowMethods, and allowHeaders fields of SecurityCors based on the values from securityCorsDTO.
func newSecurityCors(securityCorsDTO dto.SecurityCors) SecurityCors {
	return SecurityCors{
		allowOrigins: securityCorsDTO.AllowOrigins,
		allowMethods: securityCorsDTO.AllowMethods,
		allowHeaders: securityCorsDTO.AllowHeaders,
	}
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
		cache:              newCache(helper.IfNilReturns(endpointDTO.Cache, dto.Cache{})),
		responseEncode:     endpointDTO.ResponseEncode,
		aggregateResponses: endpointDTO.AggregateResponses,
		abortIfStatusCodes: endpointDTO.AbortIfStatusCodes,
		beforeware:         endpointDTO.Beforeware,
		afterware:          endpointDTO.Afterware,
		backends:           backends,
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

// AllowCacheControl checks if the caching is allowed or not.
// It uses the 'allowCacheControl' field in the 'Gopen' structure.
// In case of nil value, it defaults to 'false'.
func (g Gopen) AllowCacheControl() bool {
	return helper.IfNilReturns(g.cache.allowCacheControl, false)
}

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
