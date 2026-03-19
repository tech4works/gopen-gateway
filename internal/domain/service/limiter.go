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
	"io"
	"net/http"
	"sync"

	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	timerate "golang.org/x/time/rate"
)

type limiter struct {
	keys  map[string]*timerate.Limiter
	mutex *sync.RWMutex
}

type Limiter interface {
	AllowSize(config *vo.LimiterSizeConfig, request *vo.EndpointRequest) error
	AllowRate(config *vo.LimiterRateConfig, request *vo.EndpointRequest) error
}

func NewLimiter() Limiter {
	return &limiter{
		keys:  map[string]*timerate.Limiter{},
		mutex: &sync.RWMutex{},
	}
}

func (s *limiter) AllowSize(config *vo.LimiterSizeConfig, request *vo.EndpointRequest) error {
	if checker.IsNil(config) {
		return nil
	}

	maxMetadataSize := config.MaxMetadata()
	if checker.IsGreaterThan(request.Metadata().Size(), maxMetadataSize) {
		return domain.NewErrLimiterMetadataTooLarge(maxMetadataSize.String())
	}

	if !request.HasPayload() {
		return nil
	}

	maxPayloadSize := config.MaxPayload()

	reader := http.MaxBytesReader(nil, io.NopCloser(request.Payload().Buffer()), int64(maxPayloadSize))
	defer reader.Close()

	_, err := io.ReadAll(reader)
	if checker.NonNil(err) {
		return domain.NewErrLimiterPayloadTooLarge(maxPayloadSize.String())
	}

	return nil
}

func (s *limiter) AllowRate(config *vo.LimiterRateConfig, request *vo.EndpointRequest) error {
	if checker.IsNil(config) {
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientIP := request.Metadata().GetFirst(app.XForwardedFor)

	rateLimiter, exists := s.keys[clientIP]
	if !exists {
		rateLimiter = timerate.NewLimiter(timerate.Every(config.EveryTime()), config.Capacity())
		s.keys[clientIP] = rateLimiter
	}

	if !rateLimiter.Allow() {
		return domain.NewErrLimiterTooManyRequests(config.Capacity(), config.EveryTime())
	}

	return nil
}
