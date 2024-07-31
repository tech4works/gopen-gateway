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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"io"
	"net/http"
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

func NewRate(every Duration, capacity int) Rate {
	return Rate{
		capacity: capacity,
		every:    every,
	}
}

func NewRateDefault() Rate {
	return Rate{
		capacity: 0,
		every:    0,
	}
}

func (l Limiter) MaxHeaderSize() Bytes {
	if helper.IsGreaterThan(l.maxHeaderSize, 0) {
		return l.maxHeaderSize
	}
	return NewBytes("1MB")
}

func (l Limiter) MaxBodySize() Bytes {
	if helper.IsGreaterThan(l.maxBodySize, 0) {
		return l.maxBodySize
	}
	return NewBytes("3MB")
}

func (l Limiter) MaxMultipartMemorySize() Bytes {
	if helper.IsGreaterThan(l.maxMultipartMemorySize, 0) {
		return l.maxMultipartMemorySize
	}
	return NewBytes("5MB")
}

func (l Limiter) Rate() Rate {
	return l.rate
}

func (l Limiter) Allow(request *HTTPRequest) (err error) {
	err = l.allowHeader(request)
	if helper.IsNil(err) && helper.IsNotNil(request.Body()) {
		err = l.allowBody(request)
	}
	return err
}

func (r Rate) HasData() bool {
	return helper.IsGreaterThan(r.Capacity(), 0) && helper.IsGreaterThan(r.Every(), 0)
}

func (r Rate) NoData() bool {
	return !r.HasData()
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

func (l Limiter) allowHeader(request *HTTPRequest) (err error) {
	maxSizeAllowed := l.MaxHeaderSize()
	if helper.IsGreaterThan(request.Header().Size(), maxSizeAllowed) {
		err = mapper.NewErrHeaderTooLarge(maxSizeAllowed.String())
	}
	return err
}

func (l Limiter) allowBody(request *HTTPRequest) (err error) {
	maxSizeAllowed := l.MaxBodySize()
	if helper.ContainsIgnoreCase(request.Header().Get("Content-Type"), "multipart/form-data") {
		maxSizeAllowed = l.MaxMultipartMemorySize()
	}

	bodyBuffer := request.Body().Buffer()
	readCloser := http.MaxBytesReader(nil, io.NopCloser(bodyBuffer), int64(maxSizeAllowed))

	_, err = io.ReadAll(readCloser)
	if helper.IsNotNil(err) {
		err = mapper.NewErrPayloadTooLarge(maxSizeAllowed.String())
	}

	return err
}
