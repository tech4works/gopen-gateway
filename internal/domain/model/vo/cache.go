package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"net/http"
	"strings"
)

// Cache represents the cache configuration for an endpoint.
type Cache struct {
	// enabled represents a boolean indicating whether caching is enabled for an endpoint.
	enabled bool
	// ignoreQuery represents a boolean indicating whether to ignore query parameters when caching.
	ignoreQuery bool
	// duration represents the duration configuration for caching an endpoint httpResponse.
	duration Duration
	// strategyHeaders represents a slice of strings for strategy modifyHeaders
	strategyHeaders []string
	// onlyIfStatusCodes represents the status codes that the cache should be applied to.
	onlyIfStatusCodes []int
	// onlyIfMethods is a field in the Cache struct that represents the list of httpRequest methods that are allowed for caching.
	// If the onlyIfMethods field is empty or if the given method is present in the onlyIfMethods field, the method is allowed for caching.
	// Otherwise, it is not allowed. This field is used by the AllowMethod method in the Cache struct.
	onlyIfMethods []string
	// allowCacheControl represents a boolean value indicating whether the cache control header is allowed for the endpoint cache.
	allowCacheControl *bool
}

func newCache(cacheJson *CacheJson, endpointCacheJson *EndpointCacheJson) *Cache {
	// se os dois cache VO estiver nil retornamos nil
	if helper.IsNil(cacheJson) && helper.IsNil(endpointCacheJson) {
		return nil
	}

	// obtemos o valor do pai
	var enabled bool
	var ignoreQuery bool
	var duration Duration
	var strategyHeaders []string
	var onlyIfStatusCodes []int
	var onlyIfMethods []string
	var allowCacheControl *bool

	// caso seja informado na raiz
	if helper.IsNotNil(cacheJson) {
		duration = cacheJson.Duration
		strategyHeaders = cacheJson.StrategyHeaders
		onlyIfStatusCodes = cacheJson.OnlyIfStatusCodes
		onlyIfMethods = cacheJson.OnlyIfMethods
		allowCacheControl = cacheJson.AllowCacheControl
	}

	// caso seja informado no endpoint, damos prioridade
	if helper.IsNotNil(endpointCacheJson) {
		enabled = endpointCacheJson.Enabled
		ignoreQuery = endpointCacheJson.IgnoreQuery
		if endpointCacheJson.HasDuration() {
			duration = endpointCacheJson.Duration
		}
		if endpointCacheJson.HasStrategyHeaders() {
			strategyHeaders = endpointCacheJson.StrategyHeaders
		}
		if endpointCacheJson.HasAllowCacheControl() {
			allowCacheControl = endpointCacheJson.AllowCacheControl
		}
		if endpointCacheJson.HasOnlyIfStatusCodes() {
			onlyIfStatusCodes = endpointCacheJson.OnlyIfStatusCodes
		}
	}

	// construímos o objeto de valor com os valores padrões ou informados no json
	return &Cache{
		enabled:           enabled,
		duration:          duration,
		ignoreQuery:       ignoreQuery,
		strategyHeaders:   strategyHeaders,
		onlyIfStatusCodes: onlyIfStatusCodes,
		onlyIfMethods:     onlyIfMethods,
		allowCacheControl: allowCacheControl,
	}
}

// Duration returns the value of the duration field in the Cache struct.
func (c Cache) Duration() Duration {
	return c.duration
}

// OnlyIfStatusCodes returns the list of status codes that the cache should be applied to.
// If the onlyIfStatusCodes field is empty, it means that the cache should be applied to all status codes.
// Otherwise, the cache is only applied to the specified status codes in the list.
func (c Cache) OnlyIfStatusCodes() []int {
	return c.onlyIfStatusCodes
}

// OnlyIfMethods returns the list of httpRequest methods that are allowed for caching.
// If the onlyIfMethods field is empty or if the given method is present in the onlyIfMethods field, the method is allowed for caching.
// Otherwise, it is not allowed. This field is used by the AllowMethod method in the Cache struct.
func (c Cache) OnlyIfMethods() []string {
	return c.onlyIfMethods
}

func (c Cache) Enabled() bool {
	return c.enabled
}

func (c Cache) Disabled() bool {
	return !c.Enabled()
}

func (c Cache) IgnoreQuery() bool {
	return c.ignoreQuery
}

func (c Cache) StrategyHeaders() []string {
	return c.strategyHeaders
}

func (c Cache) CacheControl(httpRequest *HttpRequest) enum.CacheControl {
	var control enum.CacheControl
	if c.AllowCacheControl() {
		control = httpRequest.CacheControl()
	}
	return control
}

func (c Cache) CanRead(httpRequest *HttpRequest) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}
	// obtemos o cache-control
	control := c.CacheControl(httpRequest)
	// verificamos se no Cache-Control enviado veio como "no-cache" e se o método da requisição contains no campo
	// de permissão, ou esse campo esteja vazio
	return helper.IsNotEqualTo(enum.CacheControlNoCache, control) && c.AllowMethod(httpRequest.Method())
}

func (c Cache) CantRead(httpRequest *HttpRequest) bool {
	return !c.CanRead(httpRequest)
}

func (c Cache) CanWrite(httpRequest *HttpRequest, httpResponse *HttpResponse) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}
	// obtemos o cache-control
	control := c.CacheControl(httpRequest)
	// verificamos se no Cache-Control enviado veio como "no-store" e se o método da requisição contains no campo
	// de permissão, também verificamos o código de
	return helper.IsNotEqualTo(enum.CacheControlNoStore, control) && c.AllowMethod(httpRequest.Method()) &&
		c.AllowStatusCode(httpResponse.StatusCode())
}

func (c Cache) CantWrite(httpRequest *HttpRequest, httpResponse *HttpResponse) bool {
	return !c.CanWrite(httpRequest, httpResponse)
}

func (c Cache) StrategyKey(httpRequest *HttpRequest) string {
	// inicializamos a url da requisição completa
	url := httpRequest.Url()
	// caso o cache queira ignorar as queries, ele ignora
	if c.IgnoreQuery() {
		url = httpRequest.Path().String()
	}

	// construímos a chave inicialmente com os valores de requisição
	key := fmt.Sprintf("%s:%s", httpRequest.Method(), url)

	var strategyValues []string
	// iteramos as chaves para obter os valores
	for _, strategyKey := range c.strategyHeaders {
		valueByStrategyKey := httpRequest.Header().Get(strategyKey)
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

// AllowMethod checks if the given method is allowed in the Cache.
// If the onlyIfMethods field is nil, it allows all methods.
// If onlyIfMethods is empty and the method is GET, it allows the method.
// Otherwise, it checks if the method exists in the onlyIfMethods field and allows it.
// It returns true if the method is allowed, false otherwise.
func (c Cache) AllowMethod(method string) bool {
	return helper.IsNil(c.onlyIfMethods) || (helper.IsEmpty(c.onlyIfMethods) && helper.Equals(method, http.MethodGet)) ||
		helper.Contains(c.onlyIfMethods, method)
}

// AllowStatusCode checks whether the provided status code is allowed based on the following conditions:
//  1. If the `onlyIfStatusCodes` field in the `EndpointCache` struct is nil, any status code is allowed.
//  2. If the `onlyIfStatusCodes` field in the `EndpointCache` struct is empty and the status code is between
//     200 and 299 (inclusive), the status code is allowed.
//  3. If the `onlyIfStatusCodes` field in the `EndpointCache` struct contains the status code, it is allowed.
//
// It returns true if the status code is allowed; otherwise, it returns false.
func (c Cache) AllowStatusCode(statusCode int) bool {
	return helper.IsNil(c.onlyIfStatusCodes) || (helper.IsEmpty(c.onlyIfStatusCodes) &&
		helper.IsGreaterThanOrEqual(statusCode, 200) && helper.IsLessThanOrEqual(statusCode, 299)) ||
		helper.Contains(c.onlyIfStatusCodes, statusCode)
}

func (c Cache) AllowCacheControl() bool {
	return helper.IfNilReturns(c.allowCacheControl, false)
}
