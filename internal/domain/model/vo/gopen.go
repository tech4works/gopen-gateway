package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type GOpen struct {
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

type Limiter struct {
	maxHeaderSize          Bytes
	maxBodySize            Bytes
	maxMultipartMemorySize Bytes
	rate                   Rate
}

type Rate struct {
	capacity int
	every    time.Duration
}

type Cache struct {
	duration          time.Duration
	strategyHeaders   []string
	allowCacheControl *bool
}

type SecurityCors struct {
	allowOrigins []string
	allowMethods []string
	allowHeaders []string
}

type Endpoint struct {
	path               string
	method             string
	timeout            time.Duration
	limiter            Limiter
	cache              Cache
	responseEncode     enum.ResponseEncode
	aggregateResponses bool
	abortIfStatusCodes []int
	beforeware         []string
	afterware          []string
	backends           []Backend
}

// NewGOpen creates a new instance of GOpen based on the provided environment and gopenDTO.
// It initializes the fields of GOpen based on values from gopenDTO and sets default values for empty fields.
func NewGOpen(env string, gopenDTO dto.GOpen) GOpen {
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

	return GOpen{
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

// Port returns the value of the port field in the GOpen struct.
func (g GOpen) Port() int {
	return g.port
}

// HotReload returns the value of the hotReload field in the GOpen struct.
func (g GOpen) HotReload() bool {
	return g.hotReload
}

// Version returns the value of the version field in the GOpen struct.
func (g GOpen) Version() string {
	return g.version
}

// Endpoints returns the value of the endpoints field in the GOpen struct.
func (g GOpen) Endpoints() []Endpoint {
	return g.endpoints
}

// Timeout returns the value of the timeout field in the GOpen struct. If the timeout is greater than 0,
// it returns the timeout value. Otherwise, it returns a default timeout of 30 seconds
func (g GOpen) Timeout() time.Duration {
	if helper.IsGreaterThan(g.timeout, 0) {
		return g.timeout
	}
	return 30 * time.Second
}

// LimiterMaxHeaderSize returns the value of maxHeaderSize field in the Limiter struct.
// If the value is greater than 0, it returns the value.
// Otherwise, it returns a newly created Bytes instance with the value "1MB".
func (g GOpen) LimiterMaxHeaderSize() Bytes {
	if helper.IsGreaterThan(g.limiter.maxHeaderSize, 0) {
		return g.limiter.maxHeaderSize
	}
	return NewBytes("1MB")
}

// LimiterMaxBodySize returns the maximum body size limit for the GOpen struct.
// It checks if the Limiter's maxBodySize field is greater than 0 and returns it.
// If it is not greater than 0, it returns a new Bytes object initialized with the value "3MB".
func (g GOpen) LimiterMaxBodySize() Bytes {
	if helper.IsGreaterThan(g.limiter.maxBodySize, 0) {
		return g.limiter.maxBodySize
	}
	return NewBytes("3MB")
}

// LimiterMaxMultipartMemorySize returns the value of the maxMultipartMemorySize field in the Limiter struct of the GOpen object.
// If the value is greater than zero, it returns the value of that field.
// Otherwise, it returns a new Bytes object initialized with a value of "5MB".
func (g GOpen) LimiterMaxMultipartMemorySize() Bytes {
	if helper.IsGreaterThan(g.limiter.maxMultipartMemorySize, 0) {
		return g.limiter.maxMultipartMemorySize
	}
	return NewBytes("5MB")
}

// LimiterRateCapacity returns the value of the capacity field in the LimiterRateCapacity method.
// If the capacity is greater than 0, it returns the capacity field value; otherwise, it returns 5.
// This method is used to retrieve the rate capacity for a limiter.
func (g GOpen) LimiterRateCapacity() int {
	if helper.IsGreaterThan(g.limiter.rate.capacity, 0) {
		return g.limiter.rate.capacity
	}
	return 5
}

// LimiterRateEvery returns the value of the rate.every field in the Limiter struct.
// If the value is greater than 0, it returns the rate.every value.
// Otherwise, it returns a default value of 1 second.
func (g GOpen) LimiterRateEvery() time.Duration {
	if helper.IsGreaterThan(g.limiter.rate.every, 0) {
		return g.limiter.rate.every
	}
	return time.Second
}

// CacheDuration returns the duration of the cache field in the GOpen struct.
func (g GOpen) CacheDuration() time.Duration {
	return g.cache.duration
}

// CacheStrategyHeaders returns the strategyHeaders field in the Cache struct, which represents the headers
// used in cache strategy.
func (g GOpen) CacheStrategyHeaders() []string {
	return g.cache.strategyHeaders
}

// AllowCacheControl checks if the caching is allowed or not.
// It uses the 'allowCacheControl' field in the 'GOpen' structure.
// In case of nil value, it defaults to 'false'.
func (g GOpen) AllowCacheControl() bool {
	return helper.IfNilReturns(g.cache.allowCacheControl, false)
}

func (g GOpen) SecurityCors() SecurityCors {
	return g.securityCors
}

// CountMiddlewares returns the number of middlewares in the GOpen instance.
func (g GOpen) CountMiddlewares() int {
	return len(g.middlewares)
}

// CountEndpoints returns the number of endpoints in the GOpen struct.
func (g GOpen) CountEndpoints() int {
	return len(g.endpoints)
}

// CountBackends returns the total number of backends present in the `GOpen` struct and its nested `Endpoint` structs.
// It calculates the count by summing the number of middlewares in `GOpen` and recursively iterating through each `Endpoint`
// to count their backends.
// Returns an integer indicating the total count of backends.
func (g GOpen) CountBackends() (count int) {
	count += g.CountMiddlewares()
	for _, endpointVO := range g.endpoints {
		count += endpointVO.CountBackends()
	}
	return count
}

// CountModifiers counts the total number of modifiers in the GOpen struct.
// It iterates through all the middleware backends and endpoint VOs,
// and calls the CountModifiers method on each of them to calculate the count.
// The count is incremented for each modifier found and the final count is returned.
func (g GOpen) CountModifiers() (count int) {
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
func (g GOpen) Middleware(key string) (Backend, bool) {
	return g.middlewares.Get(key)
}

// Middlewares returns the value of the middlewares field in the GOpen struct.
func (g GOpen) Middlewares() Middlewares {
	return g.middlewares
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
	return helper.IsNotEmpty(e.cache)
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
	return helper.IsNotEmpty(e.cache.strategyHeaders)
}

// CacheStrategyHeaders returns the strategyHeaders field in the Cache struct of the Endpoint.
// It contains the headers that define the caching strategy for the endpoint.
func (e Endpoint) CacheStrategyHeaders() []string {
	return e.cache.strategyHeaders
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

// Duration returns the value of the duration field in the Cache struct.
func (c Cache) Duration() time.Duration {
	return c.duration
}

// Enabled returns true if the cache is enabled, false otherwise.
// The cache is considered enabled if the duration field in the Cache struct is greater than 0.
func (c Cache) Enabled() bool {
	return helper.IsGreaterThan(c.duration, 0)
}

// Disabled returns the opposite of the Enabled method. It indicates if the cache is disabled or not.
func (c Cache) Disabled() bool {
	return !c.Enabled()
}

// CanRead checks if the cache is active and if the Cache-Control header value is not "no-cache" and the HTTP method is "GET".
// If the cache is disabled, it returns false.
// It returns true if the cache is active and the conditions are met, otherwise false.
func (c Cache) CanRead(httpMethod string, header Header) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := c.CacheControlEnum(header)

	// verificamos se no Cache-Control enviado veio como "no-cache" e se o método da requisição é GET
	return helper.IsNotEqualTo(enum.CacheControlNoCache, cacheControl) && helper.Equals(httpMethod, http.MethodGet)
}

// CanWrite checks if it is possible to write to the cache based on the HTTP method and header.
// If the cache is disabled, it returns false.
// It retrieves the cache control enum from the request header.
// It returns false if the Cache-Control header contains "no-store" and the HTTP method is GET, otherwise it returns true.
func (c Cache) CanWrite(httpMethod string, header Header) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := c.CacheControlEnum(header)

	// verificamos se no Cache-Control enviado veio como "no-store" e se o método da requisição é GET
	return helper.IsNotEqualTo(enum.CacheControlNoStore, cacheControl) && helper.Equals(httpMethod, http.MethodGet)
}

// CacheControlEnum takes a Header and returns the CacheControl enum value.
// If caching is allowed and the Cache-Control header is present, it tries to parse the value and convert it
func (c Cache) CacheControlEnum(header Header) (cacheControl enum.CacheControl) {
	// caso esteja permitido o cache control obtemos do header
	if helper.IsNotNil(c.allowCacheControl) && *c.allowCacheControl {
		cacheControl = enum.CacheControl(header.Get("Cache-Control"))
	}
	return cacheControl
}

// StrategyKey generates a key for caching based on the HTTP method, URL, and header values.
// The key is initially constructed with the request method and URL.
// Then, the method iterates through the strategyHeaders field of the Cache object
// to collect the corresponding values from the header parameter.
// If the values are found, they are separated with a colon delimiter.
// If the strategyKey is not empty, it is appended to the key string.
// The final key is returned.
func (c Cache) StrategyKey(httpMethod string, httpUrl string, header Header) string {
	// construímos a chave inicialmente com os valores de requisição
	key := fmt.Sprintf("%s:%s", httpMethod, httpUrl)

	var strategyValues []string
	// iteramos as chaves para obter os valores
	for _, strategyKey := range c.strategyHeaders {
		valueByStrategyKey := header.Get(strategyKey)
		if helper.IsNotEmpty(valueByStrategyKey) {
			strategyValues = append(strategyValues, valueByStrategyKey)
		}
	}
	// caso tenha encontrado valores, separamos os mesmos
	strategyKey := strings.Join(strategyValues, ":")

	// caso o valor não esteja vazio retornamos o key padrão com a estratégia imposto no objeto de valor
	if helper.IsNotEmpty(strategyKey) {
		key = fmt.Sprintf("%s:%s", key, strategyKey)
	}
	return key
}

// EqualsContext checks if the context of the Modifier is equal to the given enum.ModifierContext.
// It returns true if the context is empty or if it is equal to the given context, otherwise it returns false.
func (m Modifier) EqualsContext(context enum.ModifierContext) bool {
	return helper.IsEmpty(m.context) || helper.Equals(m.context, context)
}

// NotEqualsContext returns `true` if the `Modifier` context is not equal to the specified `enum.ModifierContext`, otherwise `false`.
// It uses the `EqualsContext` method to check for equality.
func (m Modifier) NotEqualsContext(context enum.ModifierContext) bool {
	return !m.EqualsContext(context)
}

// Context returns the value of the context field in the Modifier struct.
func (m Modifier) Context() enum.ModifierContext {
	return m.context
}

// Scope returns the value of the scope field in the Modifier struct.
func (m Modifier) Scope() enum.ModifierScope {
	return m.scope
}

// Action returns the value of the action field in the Modifier struct.
func (m Modifier) Action() enum.ModifierAction {
	return m.action
}

// Global returns the value of the global field in the Modifier struct.
func (m Modifier) Global() bool {
	return m.global
}

// Key returns the value of the key field in the Modifier struct.
func (m Modifier) Key() string {
	return m.key
}

// Value returns the value of the value field in the Modifier struct.
func (m Modifier) Value() string {
	return m.value
}

// Valid checks if a Modifier is valid.
// A Modifier is considered valid if both the Modifier and its value are not empty.
func (m Modifier) Valid() bool {
	return helper.IsNotEmpty(m) && helper.IsNotEmpty(m.value)
}

// AllowOriginsData returns the allowOrigins field in the SecurityCors struct.
func (s SecurityCors) AllowOriginsData() []string {
	return s.allowOrigins
}

// AllowMethodsData returns the allowMethods field in the SecurityCors struct.
func (s SecurityCors) AllowMethodsData() []string {
	return s.allowMethods
}

// AllowHeadersData returns the allowHeaders field in the SecurityCors struct.
func (s SecurityCors) AllowHeadersData() []string {
	return s.allowHeaders
}

// AllowOrigins checks if the given IP is allowed by the security-cors.allow-origins configuration
// It returns an error if the configuration is not empty, does not contain "*", and does not contain the IP
func (s SecurityCors) AllowOrigins(ip string) (err error) {
	// verificamos se na configuração security-cors.allow-origins ta vazia, tem * ou tá informado o ip da requisição
	if helper.IsNotEmpty(s.allowOrigins) && helper.NotContains(s.allowOrigins, "*") &&
		helper.NotContains(s.allowOrigins, ip) {
		err = errors.New("Origin not mapped on security-cors.allow-origins")
	}
	return err
}

// AllowMethods verifies if the specified method is allowed based on the configuration in SecurityCors.
// The method checks if the allow-methods configuration is not empty, does not contain "*", and does not have the specified method.
// If the method is not allowed, it returns an error with the message "method not mapped on security-cors.allow-methods".
// It returns the error, if any, otherwise it returns nil.
func (s SecurityCors) AllowMethods(method string) (err error) {
	// verificamos se na configuração security-cors.allow-methods ta vazia, tem * ou tá informado com oq vem na requisição
	if helper.IsNotEmpty(s.allowMethods) && helper.NotContains(s.allowMethods, "*") &&
		helper.NotContains(s.allowMethods, method) {
		err = errors.New("Method not mapped on security-cors.allow-methods")
	}
	return err
}

// AllowHeaders checks if the headers in the provided Header map are allowed according to the configuration in the SecurityCors struct.
// It returns an error if any of the headers are not allowed, based on the SecurityCors.allowHeaders field.
// If the SecurityCors.allowHeaders field is empty or contains "*", all headers are considered allowed.
// Headers listed in SecurityCors.allowHeaders, as well as the X-Forwarded-For and X-Trace-Id headers, are always allowed.
// If there are headers that are not allowed, the function returns an error with a message indicating which headers are not allowed.
// If there are no headers that are not allowed, it returns nil.
func (s SecurityCors) AllowHeaders(header Header) (err error) {
	// verificamos se na configuração security-cors.allow-headers ta vazia, tem * para retornar ok
	if helper.IsEmpty(s.allowHeaders) || helper.Contains(s.allowHeaders, "*") {
		return nil
	}
	// inicializamos os headers não permitidos
	var headersNotAllowed []string
	// iteramos o header da requisição para verificar os headers que contain
	for key := range header {
		// caso o campo do header não esteja mapeado na lista security-cors.allow-headers e nao seja X-Forwarded-For
		// e nem X-Trace-Id adicionamos na lista
		if helper.NotContains(s.allowHeaders, key) && helper.IsNotEqualToIgnoreCase(key, consts.XForwardedFor) &&
			helper.IsNotEqualToIgnoreCase(key, consts.XTraceId) {
			headersNotAllowed = append(headersNotAllowed, key)
		}
	}
	// caso a lista não esteja vazia, quer dizer que tem headers não permitidos
	if helper.IsNotEmpty(headersNotAllowed) {
		headersFields := strings.Join(headersNotAllowed, ", ")
		return errors.New("Headers contains not mapped fields on security-cors.allow-headers:", headersFields)
	}

	// se tudo ocorreu bem retornamos nil
	return nil
}
