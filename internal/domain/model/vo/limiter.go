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
	"github.com/tech4works/checker"
	"time"
)

type Limiter struct {
	maxHeaderSize          Bytes
	maxBodySize            Bytes
	maxMultipartMemorySize Bytes
	rate                   Rate
}

type Rate struct {
	capacity int
	every    Duration
}

func NewLimiter(maxHeaderSize, maxBodySize, maxMultipartForm Bytes, rate Rate) Limiter {
	return Limiter{
		maxHeaderSize:          maxHeaderSize,
		maxBodySize:            maxBodySize,
		maxMultipartMemorySize: maxMultipartForm,
		rate:                   rate,
	}
}

func NewLimiterDefault() Limiter {
	return Limiter{}
}

func NewRate(every Duration, capacity int) Rate {
	return Rate{
		capacity: capacity,
		every:    every,
	}
}

func (l Limiter) MaxHeaderSize() Bytes {
	if checker.IsGreaterThan(l.maxHeaderSize, 0) {
		return l.maxHeaderSize
	}
	return NewBytes("1MB")
}

func (l Limiter) MaxBodySize() Bytes {
	if checker.IsGreaterThan(l.maxBodySize, 0) {
		return l.maxBodySize
	}
	return NewBytes("3MB")
}

func (l Limiter) MaxMultipartMemorySize() Bytes {
	if checker.IsGreaterThan(l.maxMultipartMemorySize, 0) {
		return l.maxMultipartMemorySize
	}
	return NewBytes("5MB")
}

func (l Limiter) Rate() Rate {
	return l.rate
}

func (r Rate) IsEmpty() bool {
	return checker.IsLessThanOrEqual(r.Capacity(), 0) && checker.IsLessThanOrEqual(r.Every(), 0)
}

func (r Rate) Capacity() int {
	return r.capacity
}

func (r Rate) Every() Duration {
	return r.every
}

func (r Rate) EveryTime() time.Duration {
	return r.every.Time()
}
