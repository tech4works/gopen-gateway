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

func newLimiter(limiterDTO dto.Limiter) Limiter {
	return Limiter{
		maxHeaderSize:          NewBytes(limiterDTO.MaxHeaderSize),
		maxBodySize:            NewBytes(limiterDTO.MaxBodySize),
		maxMultipartMemorySize: NewBytes(limiterDTO.MaxMultipartMemorySize),
		rate:                   newRate(helper.IfNilReturns(limiterDTO.Rate, dto.Rate{})),
	}
}

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

func newSecurityCors(securityCorsDTO dto.SecurityCors) SecurityCors {
	return SecurityCors{
		allowOrigins: securityCorsDTO.AllowOrigins,
		allowMethods: securityCorsDTO.AllowMethods,
		allowHeaders: securityCorsDTO.AllowHeaders,
	}
}

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

func NewCacheFromEndpoint(duration time.Duration, strategyHeaders []string, allowCacheControl bool) Cache {
	return Cache{
		duration:          duration,
		strategyHeaders:   strategyHeaders,
		allowCacheControl: &allowCacheControl,
	}
}

func (g GOpen) Port() int {
	return g.port
}

func (g GOpen) HotReload() bool {
	return g.hotReload
}

func (g GOpen) Version() string {
	return g.version
}

func (g GOpen) Endpoints() []Endpoint {
	return g.endpoints
}

func (g GOpen) Timeout() time.Duration {
	if helper.IsGreaterThan(g.timeout, 0) {
		return g.timeout
	}
	return 30 * time.Second
}

func (g GOpen) LimiterMaxHeaderSize() Bytes {
	if helper.IsGreaterThan(g.limiter.maxHeaderSize, 0) {
		return g.limiter.maxHeaderSize
	}
	return NewBytes("1MB")
}

func (g GOpen) LimiterMaxBodySize() Bytes {
	if helper.IsGreaterThan(g.limiter.maxBodySize, 0) {
		return g.limiter.maxBodySize
	}
	return NewBytes("3MB")
}

func (g GOpen) LimiterMaxMultipartMemorySize() Bytes {
	if helper.IsGreaterThan(g.limiter.maxMultipartMemorySize, 0) {
		return g.limiter.maxMultipartMemorySize
	}
	return NewBytes("5MB")
}

func (g GOpen) LimiterRateCapacity() int {
	if helper.IsGreaterThan(g.limiter.rate.capacity, 0) {
		return g.limiter.rate.capacity
	}
	return 5
}

func (g GOpen) LimiterRateEvery() time.Duration {
	if helper.IsGreaterThan(g.limiter.rate.every, 0) {
		return g.limiter.rate.every
	}
	return time.Second
}

func (g GOpen) CacheDuration() time.Duration {
	return g.cache.duration
}

func (g GOpen) CacheStrategyHeaders() []string {
	return g.cache.strategyHeaders
}

func (g GOpen) AllowCacheControl() bool {
	return helper.IfNilReturns(g.cache.allowCacheControl, false)
}

func (g GOpen) SecurityCors() SecurityCors {
	return g.securityCors
}

func (g GOpen) CountMiddlewares() int {
	return len(g.middlewares)
}

func (g GOpen) CountEndpoints() int {
	return len(g.endpoints)
}

func (g GOpen) CountBackends() (count int) {
	count += g.CountMiddlewares()
	for _, endpointVO := range g.endpoints {
		count += endpointVO.CountBackends()
	}
	return count
}

func (g GOpen) CountModifiers() (count int) {
	for _, middlewareBackend := range g.middlewares {
		count += middlewareBackend.CountModifiers()
	}
	for _, endpointDTO := range g.endpoints {
		count += endpointDTO.CountModifiers()
	}
	return count
}

func (g GOpen) Middleware(key string) (Backend, bool) {
	return g.middlewares.Get(key)
}

func (g GOpen) Middlewares() Middlewares {
	return g.middlewares
}

func (e Endpoint) Path() string {
	return e.path
}

func (e Endpoint) Method() string {
	return e.method
}

func (e Endpoint) Equals(route gin.RouteInfo) (err error) {
	if helper.Equals(route.Path, e.path) && helper.Equals(route.Method, e.method) {
		err = errors.NewSkipCaller(2, "Error path:", e.path, "method:", e.method, "repeat route endpoint")
	}
	return err
}

func (e Endpoint) HasTimeout() bool {
	return helper.IsGreaterThan(e.timeout, 0)
}

func (e Endpoint) Timeout() time.Duration {
	return e.timeout
}

func (e Endpoint) HasLimiter() bool {
	return helper.IsNotEmpty(e.limiter)
}

func (e Endpoint) HasLimiterRateCapacity() bool {
	return helper.IsGreaterThan(e.limiter.rate.capacity, 0)
}

func (e Endpoint) LimiterRateCapacity() int {
	return e.limiter.rate.capacity
}

func (e Endpoint) HasLimiterRateEvery() bool {
	return helper.IsGreaterThan(e.limiter.rate.every, 0)
}

func (e Endpoint) LimiterRateEvery() time.Duration {
	return e.limiter.rate.every
}

func (e Endpoint) HasLimiterMaxHeaderSize() bool {
	return helper.IsGreaterThan(e.limiter.maxHeaderSize, 0)
}

func (e Endpoint) LimiterMaxHeaderSize() Bytes {
	return e.limiter.maxHeaderSize
}

func (e Endpoint) HasLimiterMaxBodySize() bool {
	return helper.IsGreaterThan(e.limiter.maxBodySize, 0)
}

func (e Endpoint) LimiterMaxBodySize() Bytes {
	return e.limiter.maxBodySize
}

func (e Endpoint) HasLimiterMaxMultipartFormSize() bool {
	return helper.IsGreaterThan(e.limiter.maxMultipartMemorySize, 0)
}

func (e Endpoint) LimiterMaxMultipartMemorySize() Bytes {
	return e.limiter.maxMultipartMemorySize
}

func (e Endpoint) HasCache() bool {
	return helper.IsNotEmpty(e.cache)
}

func (e Endpoint) HasCacheDuration() bool {
	return helper.IsGreaterThan(e.cache.duration, 0)
}

func (e Endpoint) CacheDuration() time.Duration {
	return e.cache.duration
}

func (e Endpoint) HasCacheStrategyHeaders() bool {
	return helper.IsNotEmpty(e.cache.strategyHeaders)
}

func (e Endpoint) CacheStrategyHeaders() []string {
	return e.cache.strategyHeaders
}

func (e Endpoint) HasAllowCacheControl() bool {
	return helper.IsNotNil(e.cache.allowCacheControl)
}

func (e Endpoint) AllowCacheControl() bool {
	return helper.IfNilReturns(e.cache.allowCacheControl, false)
}

func (e Endpoint) Beforeware() []string {
	return e.beforeware
}

func (e Endpoint) Backends() []Backend {
	return e.backends
}

func (e Endpoint) Afterware() []string {
	return e.afterware
}

func (e Endpoint) CountAllBackends() int {
	return len(e.beforeware) + len(e.backends) + len(e.afterware)
}

func (e Endpoint) CountBeforewares() int {
	return len(e.beforeware)
}

func (e Endpoint) CountAfterwares() int {
	return len(e.afterware)
}

func (e Endpoint) CountBackends() int {
	return len(e.backends)
}

func (e Endpoint) CountModifiers() (count int) {
	for _, backendDTO := range e.backends {
		count += backendDTO.CountModifiers()
	}
	return count
}

func (e Endpoint) Completed(responseHistorySize int) bool {
	return helper.Equals(responseHistorySize, e.CountAllBackends())
}

func (e Endpoint) AbortSequencial(responseVO Response) bool {
	if helper.IsEmpty(e.abortIfStatusCodes) {
		return helper.IsGreaterThanOrEqual(responseVO.statusCode, http.StatusBadRequest)
	}
	return helper.Contains(e.abortIfStatusCodes, responseVO.statusCode)
}

func (e Endpoint) ResponseEncode() enum.ResponseEncode {
	return e.responseEncode
}

func (e Endpoint) AggregateResponses() bool {
	return e.aggregateResponses
}

func (e Endpoint) AbortIfStatusCodes() []int {
	return e.abortIfStatusCodes
}

func (e Endpoint) Resume() string {
	return fmt.Sprintf("%s -> \"%s\" [beforeware: %v afterware: %v backends: %v modifiers: %v]", e.method, e.path,
		e.CountBeforewares(), e.CountAfterwares(), e.CountBackends(), e.CountModifiers())
}

func (c Cache) Duration() time.Duration {
	return c.duration
}

func (c Cache) Enabled() bool {
	return helper.IsGreaterThan(c.duration, 0)
}

func (c Cache) Disabled() bool {
	return !c.Enabled()
}

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

func (c Cache) CacheControlEnum(header Header) (cacheControl enum.CacheControl) {
	// caso esteja permitido o cache control obtemos do header
	if helper.IsNotNil(c.allowCacheControl) && *c.allowCacheControl {
		cacheControl = enum.CacheControl(header.Get("Cache-Control"))
	}
	return cacheControl
}

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

func (m Modifier) EqualsContext(context enum.ModifierContext) bool {
	return helper.IsEmpty(m.context) || helper.Equals(m.context, context)
}

func (m Modifier) NotEqualsContext(context enum.ModifierContext) bool {
	return !m.EqualsContext(context)
}

func (m Modifier) Context() enum.ModifierContext {
	return m.context
}

func (m Modifier) Scope() enum.ModifierScope {
	return m.scope
}

func (m Modifier) Action() enum.ModifierAction {
	return m.action
}

func (m Modifier) Global() bool {
	return m.global
}

func (m Modifier) Key() string {
	return m.key
}

func (m Modifier) Value() string {
	return m.value
}

func (m Modifier) Valid() bool {
	return helper.IsNotEmpty(m) && helper.IsNotEmpty(m.value)
}

func (s SecurityCors) AllowOriginsData() []string {
	return s.allowOrigins
}

func (s SecurityCors) AllowMethodsData() []string {
	return s.allowMethods
}

func (s SecurityCors) AllowHeadersData() []string {
	return s.allowHeaders
}

func (s SecurityCors) AllowOrigins(ip string) (err error) {
	// verificamos se na configuração security-cors.allow-origins ta vazia, tem * ou tá informado o ip da requisição
	if helper.IsNotEmpty(s.allowOrigins) && helper.NotContains(s.allowOrigins, "*") &&
		helper.NotContains(s.allowOrigins, ip) {
		err = errors.NewSkipCaller(2, "Origin not mapped on security-cors.allow-origins")
	}
	return err
}

func (s SecurityCors) AllowMethods(method string) (err error) {
	// verificamos se na configuração security-cors.allow-methods ta vazia, tem * ou tá informado com oq vem na requisição
	if helper.IsNotEmpty(s.allowMethods) && helper.NotContains(s.allowMethods, "*") &&
		helper.NotContains(s.allowMethods, method) {
		err = errors.NewSkipCaller(2, "method not mapped on security-cors.allow-methods")
	}
	return err
}

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
		return errors.NewSkipCaller(2,
			"Headers contains not mapped fields on security-cors.allow-headers:", headersFields)
	}

	// se tudo ocorreu bem retornamos nil
	return nil
}
