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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

// limiterMiddleware is a type that represents a middleware that handles requests by checking if they are within the
// allowed rate limit and size limit. It implements the Limiter interface.
type limiterMiddleware struct {
}

// Limiter is an interface that defines the behavior of a limiter.
// Implementations of this interface must have a Do method that takes a context as input.
type Limiter interface {
	// Do execute the logic of the limiter.
	//
	// It takes a context as input to allow cancellation or passing of values between participating entities.
	Do(ctx *api.Context)
}

// NewLimiter creates a new instance of a Limiter and returns it. The Limiter is implemented by the limiterMiddleware
// struct.
func NewLimiter() Limiter {
	return limiterMiddleware{}
}

// Do handle the request by checking if it is within the allowed rate limit and size limit.
// It first gets the limiter from the endpoint, then checks if the request is within the rate limit.
// If the rate limit is exceeded, it writes an error response with status code 429 (Too Many Requests).
// Next, it checks if the request is within the size limit.
// If the payload size limit is exceeded, it writes an error response with status code 413 (Request Entity Too Large).
// If the header size limit is exceeded, it writes an error response with status code 431 (Request Header Fields Too Large).
// If all checks pass, it proceeds to the next handler.
func (l limiterMiddleware) Do(ctx *api.Context) {
	limiter := ctx.Endpoint().Limiter()

	err := limiter.Rate().Allow(ctx.HttpRequest().Header().Get(consts.XForwardedFor))
	if helper.IsNotNil(err) {
		ctx.WriteError(http.StatusTooManyRequests, err)
		return
	}

	err = limiter.Allow(ctx.HttpRequest())
	if errors.Contains(err, mapper.ErrPayloadTooLarge) {
		ctx.WriteError(http.StatusRequestEntityTooLarge, err)
		return
	} else if errors.Contains(err, mapper.ErrHeaderTooLarge) {
		ctx.WriteError(http.StatusRequestHeaderFieldsTooLarge, err)
		return
	}

	ctx.Next()
}
