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

import "github.com/GabrielHCataldo/go-helper/helper"

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

func (c Cache) AllowCacheControlNonNil() bool {
	return helper.IfNilReturns(c.allowCacheControl, false)
}

func (c Cache) HasOnlyIfMethods() bool {
	return helper.IsNotNil(c.onlyIfMethods)
}

func (c Cache) HasAnyOnlyIfMethods() bool {
	return helper.IsNotEmpty(c.onlyIfMethods)
}

func (c Cache) HasOnlyIfStatusCodes() bool {
	return helper.IsNotNil(c.onlyIfStatusCodes)
}

func (c Cache) HasAnyOnlyIfStatusCodes() bool {
	return helper.IsNotEmpty(c.onlyIfStatusCodes)
}
