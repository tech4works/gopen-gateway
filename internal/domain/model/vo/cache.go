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

type Cache struct {
	enabled           bool
	ignoreQuery       bool
	duration          Duration
	strategyHeaders   []string
	onlyIfStatusCodes []int
	onlyIfMethods     []string
	allowCacheControl *bool
}

func NewCache(
	enabled,
	ignoreQuery bool,
	duration Duration,
	strategyHeaders []string,
	onlyIfStatusCodes []int,
	onlyIfMethods []string,
	allowCacheControl *bool,
) *Cache {
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

func (c Cache) Enabled() bool {
	return c.enabled
}

func (c Cache) Disabled() bool {
	return !c.Enabled()
}

func (c Cache) IgnoreQuery() bool {
	return c.ignoreQuery
}

func (c Cache) Duration() Duration {
	return c.duration
}

func (c Cache) OnlyIfStatusCodes() []int {
	return c.onlyIfStatusCodes
}

func (c Cache) OnlyIfMethods() []string {
	return c.onlyIfMethods
}

func (c Cache) StrategyHeaders() []string {
	return c.strategyHeaders
}

// todo: realocar daq pra baixo tudo para o servico de dominio

func (c Cache) CanRead(request *HTTPRequest) bool {
	if c.Disabled() {
		return false
	}
	control := c.GetCacheControl(request)
	return helper.IsNotEqualTo(enum.CacheControlNoCache, control) && c.AllowMethod(request.Method())
}

func (c Cache) CantRead(request *HTTPRequest) bool {
	return !c.CanRead(request)
}

func (c Cache) CanWrite(request *HTTPRequest, response *HTTPResponse) bool {
	if c.Disabled() {
		return false
	}
	control := c.GetCacheControl(request)
	return helper.IsNotEqualTo(enum.CacheControlNoStore, control) && c.AllowMethod(request.Method()) &&
		c.AllowStatusCode(response.StatusCode())
}

func (c Cache) CantWrite(request *HTTPRequest, response *HTTPResponse) bool {
	return !c.CanWrite(request, response)
}

func (c Cache) StrategyKey(request *HTTPRequest) string {
	url := request.Url()
	if c.IgnoreQuery() {
		url = request.Path().String()
	}
	strategyKey := fmt.Sprintf("%s:%s", request.Method(), url)

	var strategyHeaderValues []string
	for _, strategyHeaderKey := range c.strategyHeaders {
		valueByStrategyKey := request.Header().Get(strategyHeaderKey)
		if helper.IsNotEmpty(valueByStrategyKey) {
			strategyHeaderValues = append(strategyHeaderValues, valueByStrategyKey)
		}
	}
	if helper.IsNotEmpty(strategyHeaderValues) {
		strategyKey = fmt.Sprintf("%s:%s", strategyKey, strings.Join(strategyHeaderValues, ":"))
	}

	return strategyKey
}

func (c Cache) AllowMethod(method string) bool {
	return helper.IsNil(c.onlyIfMethods) || (helper.IsEmpty(c.onlyIfMethods) && helper.Equals(method, http.MethodGet)) ||
		helper.Contains(c.onlyIfMethods, method)
}

func (c Cache) AllowStatusCode(statusCode StatusCode) bool {
	return helper.IsNil(c.onlyIfStatusCodes) || (helper.IsEmpty(c.onlyIfStatusCodes) && statusCode.OK()) ||
		helper.Contains(c.onlyIfStatusCodes, statusCode)
}

func (c Cache) AllowCacheControl() bool {
	return helper.IfNilReturns(c.allowCacheControl, false)
}

func (c Cache) GetCacheControl(request *HTTPRequest) enum.CacheControl {
	var cacheControl enum.CacheControl
	if c.AllowCacheControl() {
		cacheControl = enum.CacheControl(request.Header().GetFirst("Cache-Control"))
	}
	return cacheControl
}
