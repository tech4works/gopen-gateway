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

type Cache struct {
	duration          time.Duration
	ignoreQuery       bool
	strategyHeaders   []string
	onlyIfStatusCodes []int
	onlyIfMethods     []string
	allowCacheControl *bool
}

type EndpointCache struct {
	enabled           bool
	ignoreQuery       bool
	duration          time.Duration
	strategyHeaders   []string
	onlyIfStatusCodes []int
	allowCacheControl *bool
}

// NewCacheFromEndpoint creates a new instance of Cache based on the provided duration, strategyHeaders, allowCacheControl and enabled.
// It initializes the fields of Cache with the given values.
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

	return Cache{
		duration:          duration,
		ignoreQuery:       endpointVO.CacheIgnoreQuery(),
		strategyHeaders:   strategyHeaders,
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

// CanRead checks if it is possible to read from the cache based on the Cache-Control header
// and the HTTP method of the request.
// If the cache is disabled, it returns false.
// It retrieves the Cache-Control enum from the request header.
// It returns false if the Cache-Control header contains "no-cache" or the HTTP method is not in the onlyIfMethods field;
// otherwise, it returns true.
func (c Cache) CanRead(requestVO Request) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := c.CacheControlEnum(requestVO.Header())

	// verificamos se no Cache-Control enviado veio como "no-cache" e se o método da requisição contains no campo
	// de permissão
	return helper.IsNotEqualTo(enum.CacheControlNoCache, cacheControl) &&
		helper.Contains(c.onlyIfMethods, requestVO.Method())
}

func (c Cache) CanWrite(requestVO Request, responseVO Response) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := c.CacheControlEnum(responseVO.Header())

	// verificamos se no Cache-Control enviado veio como "no-store" e se o método da requisição contains no campo
	// de permissão, também verificamos o código de
	return helper.IsNotEqualTo(enum.CacheControlNoStore, cacheControl) &&
		helper.Contains(c.onlyIfMethods, requestVO.Method()) &&
		helper.Contains(c.onlyIfStatusCodes, responseVO.StatusCode())
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

// StrategyKey generates a key for caching based on the HTTP method, URL, and header values.
// The key is initially constructed with the request method and URL.
// Then, the method iterates through the strategyHeaders field of the Cache object
// to collect the corresponding values from the header parameter.
// If the values are found, they are separated with a colon delimiter.
// If the strategyKey is not empty, it is appended to the key string.
// The final key is returned.
func (c Cache) StrategyKey(requestVO Request) string {
	// inicializamos a url da requisição completa
	url := requestVO.Url()
	// caso o cache queira ignorar as queries, ele ignora
	if c.IgnoreQuery() {
		url = requestVO.Uri()
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
