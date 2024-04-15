package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"net/http"
	"strings"
	"time"
)

// Cache represents a cache configuration.
type Cache struct {
	// duration represents the duration of the cache in the Cache struct.
	duration time.Duration
	// strategyHeaders is a string slice that represents the list of request modifyHeaders used to generate a cache key.
	// The `StrategyKey` method in the `Cache` struct extracts values from these modifyHeaders and includes them in the cache key.
	// If no strategy values are found in the modifyHeaders, the cache key will be generated without them.
	strategyHeaders []string
	// AllowStatusCode checks if the given status code is allowed based on the onlyIfStatusCodes field in the Cache struct.
	// If the onlyIfStatusCodes field is empty or if the given status code is present in the onlyIfStatusCodes field, it returns true,
	// indicating that the status code is allowed. Otherwise, it returns false.
	onlyIfStatusCodes []int
	// onlyIfMethods is a field in the Cache struct that represents the list of request methods that are allowed for caching.
	// If the onlyIfMethods field is empty or if the given method is present in the onlyIfMethods field, the method is allowed for caching.
	// Otherwise, it is not allowed. This field is used by the AllowMethod method in the Cache struct.
	onlyIfMethods []string
	// allowCacheControl represents a pointer to a boolean indicating whether the cache should
	// honor the Cache-Control header
	allowCacheControl *bool
}

// EndpointCache represents the cache configuration for an endpoint.
type EndpointCache struct {
	// enabled represents a boolean indicating whether caching is enabled for an endpoint.
	enabled bool
	// ignoreQuery represents a boolean indicating whether to ignore query parameters when caching.
	ignoreQuery bool
	// duration represents the duration configuration for caching an endpoint response.
	duration time.Duration
	// strategyHeaders represents a slice of strings for strategy modifyHeaders
	strategyHeaders []string
	// onlyIfStatusCodes represents the status codes that the cache should be applied to.
	onlyIfStatusCodes []int
	// onlyIfMethods is a field in the Cache struct that represents the list of request methods that are allowed for caching.
	// If the onlyIfMethods field is empty or if the given method is present in the onlyIfMethods field, the method is allowed for caching.
	// Otherwise, it is not allowed. This field is used by the AllowMethod method in the Cache struct.
	onlyIfMethods []string
	// allowCacheControl represents a boolean value indicating whether the cache control header is allowed for the endpoint cache.
	allowCacheControl *bool
}

// newEndpointCache creates a new instance of Cache based on the provided Gopen and Endpoint.
// It initializes the fields of Cache based on values from Gopen and Endpoint and sets default values for empty fields.
func newEndpointCache(gopenVO Gopen, endpointVO Endpoint) EndpointCache {
	// se o endpoint não tem cache retornamos vazio
	if !endpointVO.HasCache() {
		return EndpointCache{}
	}

	// instanciamos o gopen cache
	cacheVO := gopenVO.Cache()
	// instanciamos o endpoint cache
	endpointCacheVO := endpointVO.Cache()

	// obtemos o valor do pai
	duration := cacheVO.Duration()
	// caso seja informado no endpoint, damos prioridade
	if endpointCacheVO.HasDuration() {
		duration = endpointCacheVO.Duration()
	}
	// obtemos o valor do pai
	strategyHeaders := cacheVO.StrategyHeaders()
	// caso seja informado no endpoint, damos prioridade
	if endpointCacheVO.HasStrategyHeaders() {
		strategyHeaders = endpointCacheVO.StrategyHeaders()
	}
	// obtemos o valor do pai
	allowCacheControl := cacheVO.AllowCacheControl()
	// caso seja informado no endpoint, damos prioridade
	if endpointCacheVO.HasAllowCacheControl() {
		allowCacheControl = *endpointCacheVO.AllowCacheControl()
	}
	// obtemos o valor do pai
	onlyIfStatusCodes := cacheVO.OnlyIfStatusCodes()
	// caso seja informado no endpoint, damos prioridade
	if endpointCacheVO.HasOnlyIfStatusCodes() {
		onlyIfStatusCodes = endpointCacheVO.OnlyIfStatusCodes()
	}

	// construímos o objeto vo com os valores padrões ou informados no json
	return EndpointCache{
		duration:          duration,
		ignoreQuery:       endpointCacheVO.IgnoreQuery(),
		strategyHeaders:   strategyHeaders,
		onlyIfStatusCodes: onlyIfStatusCodes,
		onlyIfMethods:     cacheVO.OnlyIfMethods(),
		allowCacheControl: &allowCacheControl,
	}
}

// newCacheFromDTO creates a new instance of Cache based on the provided cacheDTO.
// It initializes the fields of Cache based on values from cacheDTO and sets default values for empty fields.
func newCacheFromDTO(cacheDTO dto.Cache) Cache {
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
		onlyIfStatusCodes: cacheDTO.OnlyIfStatusCodes,
		onlyIfMethods:     cacheDTO.OnlyIfMethods,
		allowCacheControl: cacheDTO.AllowCacheControl,
	}
}

// newEndpointCacheFromDTO creates a new instance of EndpointCache based on the provided EndpointCacheDTO.
// It initializes the fields of EndpointCache based on values from EndpointCacheDTO and sets default values for empty fields.
func newEndpointCacheFromDTO(endpointCacheDTO dto.EndpointCache) EndpointCache {
	var duration time.Duration
	var err error
	if helper.IsNotEmpty(endpointCacheDTO.Duration) {
		duration, err = time.ParseDuration(endpointCacheDTO.Duration)
		if helper.IsNotNil(err) {
			logger.Warning("Parse duration endpoint.cache.duration err:", err)
		}
	}
	return EndpointCache{
		enabled:           endpointCacheDTO.Enabled,
		ignoreQuery:       endpointCacheDTO.IgnoreQuery,
		duration:          duration,
		strategyHeaders:   endpointCacheDTO.StrategyHeaders,
		onlyIfStatusCodes: endpointCacheDTO.OnlyIfStatusCodes,
		allowCacheControl: endpointCacheDTO.AllowCacheControl,
	}
}

// Duration returns the value of the duration field in the Cache struct.
// If the value is greater than zero, it returns the duration value.
// Otherwise, it returns a default value of 1 minute.
func (c Cache) Duration() time.Duration {
	if helper.IsGreaterThan(c.duration, 0) {
		return c.duration
	}
	return 1 * time.Minute
}

// StrategyHeaders returns the list of request headers used to generate a cache key.
func (c Cache) StrategyHeaders() []string {
	return c.strategyHeaders
}

// OnlyIfStatusCodes returns the list of status codes that are allowed for caching.
// If the onlyIfStatusCodes field is not empty, it returns the list of status codes specified in it.
// Otherwise, it returns a default list containing commonly used status codes for successful responses.
func (c Cache) OnlyIfStatusCodes() []int {
	if helper.IsNotEmpty(c.onlyIfStatusCodes) {
		return c.onlyIfStatusCodes
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

// OnlyIfMethods returns the list of request methods that are allowed for caching.
// If the onlyIfMethods field is not empty, it returns the list of methods specified in it.
// Otherwise, it returns a default list containing only the GET method.
func (c Cache) OnlyIfMethods() []string {
	if helper.IsNotEmpty(c.onlyIfMethods) {
		return c.onlyIfMethods
	}
	return []string{
		http.MethodGet,
	}
}

// AllowCacheControl checks if the caching is allowed or not.
// It uses the 'allowCacheControl' field in the 'Gopen' structure.
// In case of nil value, it defaults to 'false'.
func (c Cache) AllowCacheControl() bool {
	return helper.IfNilReturns(c.allowCacheControl, false)
}

// Enabled returns true if the cache is enabled, false otherwise.
// The cache is considered enabled if the duration field in the Cache struct is greater than 0.
func (e EndpointCache) Enabled() bool {
	return helper.IsGreaterThan(e.duration, 0)
}

// Disabled returns the opposite of the Enabled method. It indicates if the cache is disabled or not.
func (e EndpointCache) Disabled() bool {
	return !e.Enabled()
}

// IgnoreQuery returns the value of the ignoreQuery field in the EndpointCache struct.
func (e EndpointCache) IgnoreQuery() bool {
	return e.ignoreQuery
}

// HasDuration returns true if the cache duration of the Endpoint is greater than 0, otherwise false.
func (e EndpointCache) HasDuration() bool {
	return helper.IsGreaterThan(e.duration, 0)
}

// Duration returns the value of the duration field in the Cache struct.
func (e EndpointCache) Duration() time.Duration {
	return e.duration
}

// DurationStr returns the value of the duration field in the EndpointCache struct as a string.
// If the value is empty, it returns an empty string.
// Otherwise, it returns the duration value as a string.
func (e EndpointCache) DurationStr() string {
	if helper.IsEmpty(e.duration) {
		return ""
	}
	return e.duration.String()
}

// HasStrategyHeaders returns a boolean value indicating whether the `strategyHeaders` field in the EndpointCache.
// struct is not nil.
func (e EndpointCache) HasStrategyHeaders() bool {
	return helper.IsNotNil(e.strategyHeaders)
}

// StrategyHeaders returns the list of request headers used to generate a cache key.
func (e EndpointCache) StrategyHeaders() []string {
	return e.strategyHeaders
}

// HasAllowCacheControl returns a boolean value indicating whether the `allowCacheControl` field in the EndpointCache struct
// is not nil. If the field is not nil, it means that the cache control header is allowed for the endpoint cache, and the
// function returns true. Otherwise, it returns false.
func (e EndpointCache) HasAllowCacheControl() bool {
	return helper.IsNotNil(e.allowCacheControl)
}

// AllowCacheControl returns the value of the allowCacheControl field in the EndpointCache struct.
func (e EndpointCache) AllowCacheControl() *bool {
	return e.allowCacheControl
}

// HasOnlyIfStatusCodes returns a boolean value indicating whether the `onlyIfStatusCodes` field in the EndpointCache struct
// is not nil. If the field is not nil, it means that the cache should only be applied to the specified status codes, and the
// function returns true. Otherwise, it returns false.
func (e EndpointCache) HasOnlyIfStatusCodes() bool {
	return helper.IsNotNil(e.onlyIfStatusCodes)
}

// OnlyIfStatusCodes returns the list of status codes that the cache should be applied to.
// If the onlyIfStatusCodes field is empty, it means that the cache should be applied to all status codes.
// Otherwise, the cache is only applied to the specified status codes in the list.
func (e EndpointCache) OnlyIfStatusCodes() []int {
	return e.onlyIfStatusCodes
}

// OnlyIfMethods returns the list of request methods that are allowed for caching.
// If the onlyIfMethods field is empty or if the given method is present in the onlyIfMethods field, the method is allowed for caching.
// Otherwise, it is not allowed. This field is used by the AllowMethod method in the Cache struct.
func (e EndpointCache) OnlyIfMethods() []string {
	return e.onlyIfMethods
}

// CanRead checks if the cache is active and if the Cache-Control header in the request allows caching.
// It returns true if caching is allowed, false otherwise.
func (e EndpointCache) CanRead(requestVO Request) bool {
	// verificamos se ta ativo
	if e.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := e.CacheControlEnum(requestVO.Header())

	// verificamos se no Cache-Control enviado veio como "no-cache" e se o método da requisição contains no campo
	// de permissão, ou esse campo esteja vazio
	return helper.IsNotEqualTo(enum.CacheControlNoCache, cacheControl) && e.AllowMethod(requestVO.Method())
}

// CanWrite checks if the cache is active and if the Cache-Control header in the response allows caching.
// It also checks if the request method and response status code are allowed for caching.
// It returns true if caching is allowed, false otherwise.
func (e EndpointCache) CanWrite(requestVO Request, responseVO Response) bool {
	// verificamos se ta ativo
	if e.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := e.CacheControlEnum(responseVO.Header())

	// verificamos se no Cache-Control enviado veio como "no-store" e se o método da requisição contains no campo
	// de permissão, também verificamos o código de
	return helper.IsNotEqualTo(enum.CacheControlNoStore, cacheControl) && e.AllowMethod(requestVO.Method()) &&
		e.AllowStatusCode(responseVO.StatusCode())
}

// CacheControlEnum takes a Header and returns the CacheControl enum value.
// If caching is allowed and the Cache-Control header is present, it tries to parse the value and convert it
func (e EndpointCache) CacheControlEnum(header Header) (cacheControl enum.CacheControl) {
	// caso esteja permitido o cache control obtemos do header
	if helper.IsNotNil(e.allowCacheControl) && *e.allowCacheControl {
		cacheControl = enum.CacheControl(header.Get("Cache-Control"))
	}
	// retornamos a enum do cache control vazia ou não, dependendo da configuração
	return cacheControl
}

// StrategyKey generates a cache key based on the request information and strategy headers.
// If IgnoreQuery is enabled in the Cache struct, the query parameters will be ignored in the key generation.
// The generated key follows the pattern: "{HTTP Method}:{Request URL}:{Strategy Value 1}:{Strategy Value 2}:..."
// The Strategy Value is obtained from the request headers specified in the strategyHeaders field of the Cache struct.
// If no Strategy Value is found, the key will be generated without it.
// The final key is returned as a string.
func (e EndpointCache) StrategyKey(requestVO Request) string {
	// inicializamos a url da requisição completa
	url := requestVO.Url()
	// caso o cache queira ignorar as queries, ele ignora
	if e.IgnoreQuery() {
		url = requestVO.Path()
	}

	// construímos a chave inicialmente com os valores de requisição
	key := fmt.Sprintf("%s:%s", requestVO.Method(), url)

	var strategyValues []string
	// iteramos as chaves para obter os valores
	for _, strategyKey := range e.strategyHeaders {
		valueByStrategyKey := requestVO.Header().Get(strategyKey)
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

	// retornamos a key construída
	return key
}

// AllowMethod checks if the given method is allowed based on the onlyIfMethods field in the Cache struct.
// If the onlyIfMethods field is empty or if the given method is present in the onlyIfMethods field, it returns true,
// indicating that the method is allowed. Otherwise, it returns false.
func (e EndpointCache) AllowMethod(method string) bool {
	return helper.IsEmpty(e.onlyIfMethods) || helper.Contains(e.onlyIfMethods, method)
}

// AllowStatusCode checks if the given status code is allowed based on the onlyIfStatusCodes field in the Cache struct.
// If the onlyIfStatusCodes field is empty or if the given status code is present in the onlyIfStatusCodes field, it returns true,
// indicating that the status code is allowed. Otherwise, it returns false.
func (e EndpointCache) AllowStatusCode(statusCode int) bool {
	return helper.IsEmpty(e.onlyIfStatusCodes) || helper.Contains(e.onlyIfStatusCodes, statusCode)
}
