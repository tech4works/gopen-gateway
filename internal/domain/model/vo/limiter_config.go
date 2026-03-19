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

import (
	"time"

	"github.com/tech4works/checker"
)

type LimiterConfig struct {
	size *LimiterSizeConfig
	rate *LimiterRateConfig
}

type LimiterSizeConfig struct {
	maxMetadata Bytes
	maxPayload  Bytes
}

type LimiterRateConfig struct {
	capacity int
	every    Duration
}

func NewLimiterConfig(size *LimiterSizeConfig, rate *LimiterRateConfig) *LimiterConfig {
	return &LimiterConfig{
		size: size,
		rate: rate,
	}
}

func NewLimiterConfigDefault() *LimiterConfig {
	return &LimiterConfig{
		size: &LimiterSizeConfig{},
	}
}

func NewLimiterSizeConfig(maxMetadata, maxPayload Bytes) *LimiterSizeConfig {
	return &LimiterSizeConfig{
		maxMetadata: maxMetadata,
		maxPayload:  maxPayload,
	}
}

func NewLimiterRateConfig(every Duration, capacity int) *LimiterRateConfig {
	return &LimiterRateConfig{
		capacity: capacity,
		every:    every,
	}
}

func (l *LimiterConfig) Size() *LimiterSizeConfig {
	return l.size
}

func (l *LimiterConfig) Rate() *LimiterRateConfig {
	return l.rate
}

func (l *LimiterSizeConfig) MaxMetadata() Bytes {
	if checker.IsGreaterThan(l.maxMetadata, 0) {
		return l.maxMetadata
	}
	return NewBytes("1MB")
}

func (l *LimiterSizeConfig) MaxPayload() Bytes {
	if checker.IsGreaterThan(l.maxPayload, 0) {
		return l.maxPayload
	}
	return NewBytes("3MB")
}

func (r *LimiterRateConfig) Capacity() int {
	return r.capacity
}

func (r *LimiterRateConfig) Every() Duration {
	return r.every
}

func (r *LimiterRateConfig) EveryTime() time.Duration {
	return r.every.Time()
}
