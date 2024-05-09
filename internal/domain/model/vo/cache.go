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

// newCache returns a new Cache object with values constructed from the provided CacheJson and EndpointCacheJson objects.
// If both cacheJson and endpointCacheJson are nil, it returns nil.
// The function checks if cacheJson is not nil and sets the specified values from cacheJson.
// The function then checks if endpointCacheJson is not nil and updates the values from endpointCacheJson with priority.
// Finally, it constructs a new Cache object with the obtained values and returns it.
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

// Enabled returns the value of the enabled field in the Cache struct.
func (c Cache) Enabled() bool {
	return c.enabled
}

// Disabled returns the inverse value of Enabled(). It indicates if caching is disabled for the endpoint.
func (c Cache) Disabled() bool {
	return !c.Enabled()
}

// IgnoreQuery returns the value of the ignoreQuery field in the Cache struct.
func (c Cache) IgnoreQuery() bool {
	return c.ignoreQuery
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

// StrategyHeaders returns the value of the strategyHeaders field in the Cache struct.
func (c Cache) StrategyHeaders() []string {
	return c.strategyHeaders
}

// CacheControl returns the Cache-Control header value from the provided HttpRequest header,
// if the cache is allowed for the endpoint based on the Cache configuration. Otherwise, it returns an empty string.
//
// Parameters:
// - httpRequest: the HttpRequest object used to retrieve the Cache-Control header value.
//
// Returns:
// - enum.CacheControl: the Cache-Control header value, represented as an enum value of type enum.CacheControl.
// If caching is not allowed for the endpoint or if the Cache-Control header is not present in the HttpRequest header,
// it returns an empty enum value.
func (c Cache) CacheControl(httpRequest *HttpRequest) enum.CacheControl {
	var control enum.CacheControl
	if c.AllowCacheControl() {
		control = enum.CacheControl(httpRequest.Header().Get("Cache-Control"))
	}
	return control
}

// CanRead checks if the cache is enabled and allows reading based on the cache configuration and the HTTP request.
// It returns true if caching is enabled, the cache control value is not "no-cache", and the HTTP method is allowed;
// otherwise, it returns false.
func (c Cache) CanRead(httpRequest *HttpRequest) bool {
	if c.Disabled() {
		return false
	}
	control := c.CacheControl(httpRequest)
	return helper.IsNotEqualTo(enum.CacheControlNoCache, control) && c.AllowMethod(httpRequest.Method())
}

// CantRead returns the inverse value of CanRead method by passing the httpRequest parameter.
//
// Parameters:
// - httpRequest: the HTTP request used to check if caching is allowed.
//
// Returns:
// - bool: true if caching is not allowed based on the CanRead method, false otherwise.
func (c Cache) CantRead(httpRequest *HttpRequest) bool {
	return !c.CanRead(httpRequest)
}

// CanWrite checks if caching is enabled and allows writing based on the cache configuration and the HTTP request
// and response.
// It returns true if caching is enabled, the cache control value is not "no-store", the HTTP method is allowed, and
// the response status code is allowed;
// otherwise, it returns false.
func (c Cache) CanWrite(httpRequest *HttpRequest, httpResponse *HttpResponse) bool {
	if c.Disabled() {
		return false
	}
	control := c.CacheControl(httpRequest)
	return helper.IsNotEqualTo(enum.CacheControlNoStore, control) && c.AllowMethod(httpRequest.Method()) &&
		c.AllowStatusCode(httpResponse.StatusCode())
}

// CantWrite returns the opposite value of the CanWrite method by passing the httpRequest and httpResponse parameters.
func (c Cache) CantWrite(httpRequest *HttpRequest, httpResponse *HttpResponse) bool {
	return !c.CanWrite(httpRequest, httpResponse)
}

// StrategyKey generates a key for caching based on the provided HTTP request.
// The key is generated by concatenating the HTTP method and URL of the request.
// If the cache is configured to ignore query parameters, only the path of the URL is used.
// Additionally, any header values specified in the `strategyHeaders` field of the Cache struct
// are appended to the key in the format "key1:value1:key2:value2:...".
// The generated key is returned as a string.
func (c Cache) StrategyKey(httpRequest *HttpRequest) string {
	url := httpRequest.Url()
	if c.IgnoreQuery() {
		url = httpRequest.Path().String()
	}
	strategyKey := fmt.Sprintf("%s:%s", httpRequest.Method(), url)

	var strategyHeaderValues []string
	for _, strategyHeaderKey := range c.strategyHeaders {
		valueByStrategyKey := httpRequest.Header().Get(strategyHeaderKey)
		if helper.IsNotEmpty(valueByStrategyKey) {
			strategyHeaderValues = append(strategyHeaderValues, valueByStrategyKey)
		}
	}
	if helper.IsNotEmpty(strategyHeaderValues) {
		strategyKey = fmt.Sprintf("%s:%s", strategyKey, strings.Join(strategyHeaderValues, ":"))
	}

	return strategyKey
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

// AllowCacheControl returns the value of the allowCacheControl field in the Cache struct.
func (c Cache) AllowCacheControl() bool {
	return helper.IfNilReturns(c.allowCacheControl, false)
}
