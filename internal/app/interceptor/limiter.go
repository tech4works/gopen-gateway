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

package interceptor

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type limiterMiddleware struct {
	service service.Limiter
}

type Limiter interface {
	Do(ctx app.Context)
}

func NewLimiter(service service.Limiter) Limiter {
	return limiterMiddleware{
		service: service,
	}
}

func (l limiterMiddleware) Do(ctx app.Context) {
	if !ctx.Endpoint().HasLimiter() {
		ctx.Next()
		return
	}

	err := l.service.AllowRate(ctx.Endpoint().Limiter().Rate(), ctx.Request())
	if errors.Is(err, domain.ErrLimiterTooManyRequests) {
		ctx.WriteError(enum.ResponseStatusResourceExhausted, err)
		return
	} else if checker.NonNil(err) {
		ctx.WriteError(enum.ResponseStatusInternalError, err)
		return
	}

	err = l.service.AllowSize(ctx.Endpoint().Limiter().Size(), ctx.Request())
	if errors.Is(err, domain.ErrLimiterMetadataTooLarge) {
		ctx.WriteError(enum.ResponseStatusMetadataTooLarge, err)
	} else if errors.Is(err, domain.ErrLimiterPayloadTooLarge) {
		ctx.WriteError(enum.ResponseStatusPayloadTooLarge, err)
	} else if checker.NonNil(err) {
		ctx.WriteError(enum.ResponseStatusInternalError, err)
	} else {
		ctx.Next()
	}
}
