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

package service

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	timerate "golang.org/x/time/rate"
	"io"
	"net/http"
	"sync"
)

type limiterService struct {
	keys  map[string]*timerate.Limiter
	mutex *sync.RWMutex
}

type Limiter interface {
	AllowRate(request *vo.HTTPRequest, rate vo.Rate) error
	AllowSize(request *vo.HTTPRequest, limiter vo.Limiter) error
}

func NewLimiter() Limiter {
	return &limiterService{
		keys:  map[string]*timerate.Limiter{},
		mutex: &sync.RWMutex{},
	}
}

func (s *limiterService) AllowRate(request *vo.HTTPRequest, rate vo.Rate) (err error) {
	if rate.IsEmpty() {
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientIP := request.Header().GetFirst(mapper.XForwardedFor)

	rateLimiter, exists := s.keys[clientIP]
	if !exists {
		rateLimiter = timerate.NewLimiter(timerate.Every(rate.EveryTime()), rate.Capacity())
		s.keys[clientIP] = rateLimiter
	}

	if !rateLimiter.Allow() {
		err = mapper.NewErrTooManyRequests(rate.Capacity(), rate.EveryTime())
	}

	return err
}

func (s *limiterService) AllowSize(request *vo.HTTPRequest, limiter vo.Limiter) error {
	maxHeaderSize := limiter.MaxHeaderSize()
	if checker.IsGreaterThan(request.Header().Size(), maxHeaderSize) {
		return mapper.NewErrHeaderTooLarge(maxHeaderSize.String())
	}

	maxBodySize := limiter.MaxBodySize()
	if checker.ContainsIgnoreCase(request.Header().Get(mapper.ContentType), "multipart/form-data") {
		maxBodySize = limiter.MaxMultipartMemorySize()
	}

	if !request.HasBody() {
		return nil
	}

	bodyBuffer := request.Body().Buffer()
	readCloser := http.MaxBytesReader(nil, io.NopCloser(bodyBuffer), int64(maxBodySize))

	_, err := io.ReadAll(readCloser)
	if checker.NonNil(err) {
		return mapper.NewErrPayloadTooLarge(maxBodySize.String())
	}

	return nil
}
