package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"strings"
	"time"
)

// Cache represents a cache configuration.
type Cache struct {
	// duration represents the duration of the cache in the Cache struct.
	duration time.Duration
	// ignoreQuery is a boolean field in the Cache struct. When enabled, it indicates that the
	// query parameters should be ignored when generating a cache key in the StrategyKey method.
	ignoreQuery bool
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
	// allowCacheControl represents a boolean value indicating whether the cache control header is allowed for the endpoint cache.
	allowCacheControl *bool
}

// NewCacheFromEndpoint creates a new instance of Cache based on the provided Gopen and Endpoint.
// It initializes the fields of Cache based on values from Gopen and Endpoint and sets default values for empty fields.
func NewCacheFromEndpoint(gopenVO Gopen, endpointVO Endpoint) Cache {
	// se o endpoint não tem cache retornamos vazio
	if !endpointVO.HasCache() {
		return Cache{}
	}

	// obtemos o valor do pai
	duration := gopenVO.CacheDuration()
	// caso seja informado no endpoint, damos prioridade
	if endpointVO.HasCacheDuration() {
		duration = endpointVO.CacheDuration()
	}
	// obtemos o valor do pai
	strategyHeaders := gopenVO.CacheStrategyHeaders()
	// caso seja informado no endpoint, damos prioridade
	if endpointVO.HasCacheStrategyHeaders() {
		strategyHeaders = endpointVO.CacheStrategyHeaders()
	}
	// obtemos o valor do pai
	allowCacheControl := gopenVO.AllowCacheControl()
	// caso seja informado no endpoint, damos prioridade
	if endpointVO.HasAllowCacheControl() {
		allowCacheControl = endpointVO.AllowCacheControl()
	}

	// obtemos o valor do pai
	onlyIfStatusCodes := gopenVO.CacheOnlyIfStatusCodes()
	// caso seja informado no endpoint, damos prioridade
	if endpointVO.HasCacheOnlyIfStatusCodes() {
		onlyIfStatusCodes = endpointVO.CacheOnlyIfStatusCodes()
	}

	// construímos o objeto vo com os valores padrões ou informados no json
	return Cache{
		duration:          duration,
		ignoreQuery:       endpointVO.CacheIgnoreQuery(),
		strategyHeaders:   strategyHeaders,
		onlyIfStatusCodes: onlyIfStatusCodes,
		onlyIfMethods:     gopenVO.CacheOnlyIfMethods(),
		allowCacheControl: &allowCacheControl,
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
		ignoreQuery:       false,
		strategyHeaders:   cacheDTO.StrategyHeaders,
		onlyIfStatusCodes: cacheDTO.OnlyIfStatusCodes,
		onlyIfMethods:     cacheDTO.OnlyIfMethods,
		allowCacheControl: cacheDTO.AllowCacheControl,
	}
}

// newEndpointCache creates a new instance of EndpointCache based on the provided EndpointCacheDTO.
// It initializes the fields of EndpointCache based on values from EndpointCacheDTO and sets default values for empty fields.
func newEndpointCache(endpointCacheDTO dto.EndpointCache) EndpointCache {
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

// IgnoreQuery returns the value of the ignoreQuery field in the Cache struct.
func (c Cache) IgnoreQuery() bool {
	return c.ignoreQuery
}

// CanRead checks if the cache is active and if the Cache-Control header in the request allows caching.
// It returns true if caching is allowed, false otherwise.
func (c Cache) CanRead(requestVO Request) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := c.CacheControlEnum(requestVO.Header())

	// verificamos se no Cache-Control enviado veio como "no-cache" e se o método da requisição contains no campo
	// de permissão, ou esse campo esteja vazio
	return helper.IsNotEqualTo(enum.CacheControlNoCache, cacheControl) && c.AllowMethod(requestVO.Method())
}

// CanWrite checks if the cache is active and if the Cache-Control header in the response allows caching.
// It also checks if the request method and response status code are allowed for caching.
// It returns true if caching is allowed, false otherwise.
func (c Cache) CanWrite(requestVO Request, responseVO Response) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := c.CacheControlEnum(responseVO.Header())

	// verificamos se no Cache-Control enviado veio como "no-store" e se o método da requisição contains no campo
	// de permissão, também verificamos o código de
	return helper.IsNotEqualTo(enum.CacheControlNoStore, cacheControl) && c.AllowMethod(requestVO.Method()) &&
		c.AllowStatusCode(responseVO.StatusCode())
}

// CacheControlEnum takes a Header and returns the CacheControl enum value.
// If caching is allowed and the Cache-Control header is present, it tries to parse the value and convert it
func (c Cache) CacheControlEnum(header Header) (cacheControl enum.CacheControl) {
	// caso esteja permitido o cache control obtemos do header
	if helper.IsNotNil(c.allowCacheControl) && *c.allowCacheControl {
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
func (c Cache) StrategyKey(requestVO Request) string {
	// inicializamos a url da requisição completa
	url := requestVO.Url()
	// caso o cache queira ignorar as queries, ele ignora
	if c.IgnoreQuery() {
		url = requestVO.Path()
	}

	// construímos a chave inicialmente com os valores de requisição
	key := fmt.Sprintf("%s:%s", requestVO.Method(), url)

	var strategyValues []string
	// iteramos as chaves para obter os valores
	for _, strategyKey := range c.strategyHeaders {
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
func (c Cache) AllowMethod(method string) bool {
	return helper.IsEmpty(c.onlyIfMethods) || helper.Contains(c.onlyIfMethods, method)
}

// AllowStatusCode checks if the given status code is allowed based on the onlyIfStatusCodes field in the Cache struct.
// If the onlyIfStatusCodes field is empty or if the given status code is present in the onlyIfStatusCodes field, it returns true,
// indicating that the status code is allowed. Otherwise, it returns false.
func (c Cache) AllowStatusCode(statusCode int) bool {
	return helper.IsEmpty(c.onlyIfStatusCodes) || helper.Contains(c.onlyIfStatusCodes, statusCode)
}
