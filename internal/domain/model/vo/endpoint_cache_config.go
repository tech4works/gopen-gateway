/*
 * Copyright 2024 Tech4Works
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

import "github.com/tech4works/gopen-gateway/internal/domain/model/enum"

type CacheConfig struct {
	kind  enum.CacheKind
	read  CacheDecisionConfig
	write CacheDecisionConfig
	key   string
	ttl   Duration
}

type CacheDecisionConfig struct {
	onlyIf   []string
	ignoreIf []string
}

func NewCacheConfig(
	kind enum.CacheKind,
	read,
	write CacheDecisionConfig,
	key string,
	ttl Duration,
) *CacheConfig {
	return &CacheConfig{
		kind:  kind,
		read:  read,
		write: write,
		key:   key,
		ttl:   ttl,
	}
}

func NewCacheDecisionConfig(onlyIf, ignoreIf []string) CacheDecisionConfig {
	return CacheDecisionConfig{
		onlyIf:   onlyIf,
		ignoreIf: ignoreIf,
	}
}

func (c CacheConfig) Kind() enum.CacheKind {
	return c.kind
}

func (c CacheConfig) Read() CacheDecisionConfig {
	return c.read
}

func (c CacheConfig) Write() CacheDecisionConfig {
	return c.write
}

func (c CacheConfig) Key() string {
	return c.key
}

func (c CacheConfig) TTL() Duration {
	return c.ttl
}

func (c CacheDecisionConfig) OnlyIf() []string {
	return c.onlyIf
}

func (c CacheDecisionConfig) IgnoreIf() []string {
	return c.ignoreIf
}
