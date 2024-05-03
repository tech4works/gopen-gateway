/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"net/http"
	"strings"
	"time"
)

// Cache represents a cache configuration.
type Cache struct {
	// duration represents the duration of the cache in the Cache struct.
	duration Duration
	// strategyHeaders is a string slice that represents the list of httpRequest modifyHeaders used to generate a cache key.
	// The `StrategyKey` method in the `Cache` struct extracts values from these modifyHeaders and includes them in the cache key.
	// If no strategy values are found in the modifyHeaders, the cache key will be generated without them.
	strategyHeaders []string
	// AllowStatusCode checks if the given status code is allowed based on the onlyIfStatusCodes field in the Cache struct.
	// If the onlyIfStatusCodes field is empty or if the given status code is present in the onlyIfStatusCodes field, it returns true,
	// indicating that the status code is allowed. Otherwise, it returns false.
	onlyIfStatusCodes []int
	// onlyIfMethods is a field in the Cache struct that represents the list of httpRequest methods that are allowed for caching.
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

// newEndpointCache constructs a new EndpointCache object using the provided cacheVO and endpointCacheVO.
// If both cacheVO and endpointCacheVO are nil, it returns nil.
// If cacheVO is not nil, it retrieves values from cacheVO and assigns them to the corresponding fields in EndpointCache.
// If endpointCacheVO is not nil, it retrieves values from endpointCacheVO and assigns them to the corresponding fields in EndpointCache, overriding the values from cacheVO.
// It then returns the constructed EndpointCache object.
func newEndpointCache(cacheVO *Cache, endpointCacheVO *EndpointCacheJson) *EndpointCache {
	// se os dois cache VO estiver nil retornamos nil
	if helper.IsNil(cacheVO) && helper.IsNil(endpointCacheVO) {
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
	if helper.IsNotNil(cacheVO) {
		duration = cacheVO.Duration()
		strategyHeaders = cacheVO.StrategyHeaders()
		onlyIfStatusCodes = cacheVO.OnlyIfStatusCodes()
		onlyIfMethods = cacheVO.OnlyIfMethods()
		allowCacheControl = cacheVO.AllowCacheControl()
	}

	// caso seja informado no endpoint, damos prioridade
	if helper.IsNotNil(endpointCacheVO) {
		enabled = endpointCacheVO.Enabled
		ignoreQuery = endpointCacheVO.IgnoreQuery
		if endpointCacheVO.HasDuration() {
			duration = endpointCacheVO.Duration
		}
		if endpointCacheVO.HasStrategyHeaders() {
			strategyHeaders = endpointCacheVO.StrategyHeaders
		}
		if endpointCacheVO.HasAllowCacheControl() {
			allowCacheControl = endpointCacheVO.AllowCacheControl
		}
		if endpointCacheVO.HasOnlyIfStatusCodes() {
			onlyIfStatusCodes = endpointCacheVO.OnlyIfStatusCodes
		}
	}

	// construímos o objeto vo com os valores padrões ou informados no json
	return &EndpointCache{
		enabled:           enabled,
		duration:          duration,
		ignoreQuery:       ignoreQuery,
		strategyHeaders:   strategyHeaders,
		onlyIfStatusCodes: onlyIfStatusCodes,
		onlyIfMethods:     onlyIfMethods,
		allowCacheControl: allowCacheControl,
	}
}

// newCache constructs a new Cache object using the provided cacheJsonVO.
// If cacheJsonVO is nil, it returns nil.
// Otherwise, it assigns values from cacheJsonVO to the corresponding fields in Cache,
// and returns the constructed Cache object.
func newCache(cacheJsonVO *CacheJson) *Cache {
	if helper.IsNil(cacheJsonVO) {
		return nil
	}
	return &Cache{
		duration:          cacheJsonVO.Duration,
		strategyHeaders:   cacheJsonVO.StrategyHeaders,
		onlyIfStatusCodes: cacheJsonVO.OnlyIfStatusCodes,
		onlyIfMethods:     cacheJsonVO.OnlyIfMethods,
		allowCacheControl: cacheJsonVO.AllowCacheControl,
	}
}

// Duration returns the value of the duration field in the Cache struct.
// If the value is greater than zero, it returns the duration value.
// Otherwise, it returns a default value of 1 minute.
func (c Cache) Duration() Duration {
	if helper.IsGreaterThan(c.duration, 0) {
		return c.duration
	}
	return Duration(1 * time.Minute)
}

// StrategyHeaders returns the list of httpRequest headers used to generate a cache key.
func (c Cache) StrategyHeaders() []string {
	return c.strategyHeaders
}

// OnlyIfStatusCodes returns the list of status codes that are allowed for caching.
// If the onlyIfStatusCodes field is not empty, it returns the list of status codes specified in it.
// Otherwise, it returns an empty list.
func (c Cache) OnlyIfStatusCodes() []int {
	return c.onlyIfStatusCodes
}

// OnlyIfMethods returns the list of httpRequest methods that are allowed for caching.
// If the onlyIfMethods field is not empty, it returns the list of methods specified in it.
// Otherwise, it returns an empty list.
func (c Cache) OnlyIfMethods() []string {
	return c.onlyIfMethods
}

// AllowCacheControl checks if the caching is allowed or not.
// It uses the 'allowCacheControl' field in the 'Gopen' structure.
// In case of nil value, it defaults to 'false'.
func (c Cache) AllowCacheControl() *bool {
	return c.allowCacheControl
}

// Enabled returns the value of the enabled field in the EndpointCache struct.
// It returns the boolean value indicating whether the endpoint cache is enabled or not.
func (e EndpointCache) Enabled() bool {
	return e.enabled && e.HasDuration()
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
func (e EndpointCache) Duration() Duration {
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

// OnlyIfStatusCodes returns the list of status codes that the cache should be applied to.
// If the onlyIfStatusCodes field is empty, it means that the cache should be applied to all status codes.
// Otherwise, the cache is only applied to the specified status codes in the list.
func (e EndpointCache) OnlyIfStatusCodes() []int {
	return e.onlyIfStatusCodes
}

// OnlyIfMethods returns the list of httpRequest methods that are allowed for caching.
// If the onlyIfMethods field is empty or if the given method is present in the onlyIfMethods field, the method is allowed for caching.
// Otherwise, it is not allowed. This field is used by the AllowMethod method in the Cache struct.
func (e EndpointCache) OnlyIfMethods() []string {
	return e.onlyIfMethods
}

// CanRead checks if the cache is active and if the Cache-Control header in the httpRequest allows caching.
// It returns true if caching is allowed, false otherwise.
func (e EndpointCache) CanRead(requestVO *HttpRequest) bool {
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

// CanWrite checks if the cache is active and if the Cache-Control header in the httpResponse allows caching.
// It also checks if the httpRequest method and httpResponse status code are allowed for caching.
// It returns true if caching is allowed, false otherwise.
func (e EndpointCache) CanWrite(requestVO *HttpRequest, responseVO *HttpResponse) bool {
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

// StrategyKey generates a cache key based on the httpRequest information and strategy headers.
// If IgnoreQuery is enabled in the Cache struct, the query parameters will be ignored in the key generation.
// The generated key follows the pattern: "{HTTP Method}:{HttpRequest URL}:{Strategy Value 1}:{Strategy Value 2}:..."
// The Strategy Value is obtained from the httpRequest headers specified in the strategyHeaders field of the Cache struct.
// If no Strategy Value is found, the key will be generated without it.
// The final key is returned as a string.
func (e EndpointCache) StrategyKey(requestVO *HttpRequest) string {
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

// AllowMethod checks if the given method is allowed in the EndpointCache.
// If the onlyIfMethods field is nil, it allows all methods.
// If onlyIfMethods is empty and the method is GET, it allows the method.
// Otherwise, it checks if the method exists in the onlyIfMethods field and allows it.
// It returns true if the method is allowed, false otherwise.
func (e EndpointCache) AllowMethod(method string) bool {
	return helper.IsNil(e.onlyIfMethods) || (helper.IsEmpty(e.onlyIfMethods) && helper.Equals(method, http.MethodGet)) ||
		helper.Contains(e.onlyIfMethods, method)
}

// AllowStatusCode checks whether the provided status code is allowed based on the following conditions:
//  1. If the `onlyIfStatusCodes` field in the `EndpointCache` struct is nil, any status code is allowed.
//  2. If the `onlyIfStatusCodes` field in the `EndpointCache` struct is empty and the status code is between
//     200 and 299 (inclusive), the status code is allowed.
//  3. If the `onlyIfStatusCodes` field in the `EndpointCache` struct contains the status code, it is allowed.
//
// It returns true if the status code is allowed; otherwise, it returns false.
func (e EndpointCache) AllowStatusCode(statusCode int) bool {
	return helper.IsNil(e.onlyIfStatusCodes) || (helper.IsEmpty(e.onlyIfStatusCodes) &&
		helper.IsGreaterThanOrEqual(statusCode, 200) && helper.IsLessThanOrEqual(statusCode, 299)) ||
		helper.Contains(e.onlyIfStatusCodes, statusCode)
}
