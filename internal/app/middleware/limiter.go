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

package middleware

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"net/http"
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
	err := l.service.AllowRate(ctx.Request(), ctx.Endpoint().Limiter().Rate())
	if helper.IsNotNil(err) {
		ctx.WriteError(http.StatusTooManyRequests, err)
		return
	}

	err = l.service.AllowSize(ctx.Request(), ctx.Endpoint().Limiter())
	if errors.Contains(err, mapper.ErrPayloadTooLarge) {
		ctx.WriteError(http.StatusRequestEntityTooLarge, err)
	} else if errors.Contains(err, mapper.ErrHeaderTooLarge) {
		ctx.WriteError(http.StatusRequestHeaderFieldsTooLarge, err)
	} else {
		ctx.Next()
	}
}
