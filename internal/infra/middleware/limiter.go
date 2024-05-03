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
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

type limiter struct {
}

// Limiter is an interface that defines a method for handling rate limiting and size limiting
// for API endpoints. The Do method takes a rateLimiterProvider and a sizeLimiterProvider
// as arguments and returns a HandlerFunc. The returned HandlerFunc can be used as an HTTP
// route handler.
type Limiter interface {
	// Do handles rate limiting and size limiting for API endpoints.
	// It takes a rateLimiterProvider and a sizeLimiterProvider as arguments
	// and returns a HandlerFunc that can be used as an HTTP route handler.
	Do(rateLimiterProvider infra.RateLimiterProvider, sizeLimiterProvider infra.SizeLimiterProvider) api.HandlerFunc
}

// NewLimiter creates a new instance of Limiter.
// It returns a Limiter value which implements the Limiter interface.
func NewLimiter() Limiter {
	return limiter{}
}

// Do execute a handler that implements the api.HandlerFunc interface, providing rate limiting and size limiting functionality.
// It takes a RateLimiterProvider instance and a SizeLimiterProvider instance as input parameters.
// The rateLimiterProvider is used to check whether the request is allowed based on the rate limit.
// The sizeLimiterProvider is used to check whether the request size is within the allowed limit.
// If the request is not allowed, it returns an error and writes an error response to the request.
// If the request is allowed, it calls the Next() method of the HttpRequest object to execute the next handler in the chain.
func (l limiter) Do(rateLimiterProvider infra.RateLimiterProvider, sizeLimiterProvider infra.SizeLimiterProvider,
) api.HandlerFunc {
	return func(ctx *api.Context) {
		// aqui ja verificamos se a chave hoje sendo ela o IP está permitida
		err := rateLimiterProvider.Allow(ctx.HeaderValue(consts.XForwardedFor))
		if helper.IsNotNil(err) {
			ctx.WriteError(http.StatusTooManyRequests, err)
			return
		}

		// verificamos o tamanho da requisição, e tratamos o erro logo em seguida
		err = sizeLimiterProvider.Allow(ctx.Http())
		if errors.Contains(err, domainmapper.ErrHeaderTooLarge) {
			ctx.WriteError(http.StatusRequestHeaderFieldsTooLarge, err)
			return
		} else if helper.IsNotNil(err) {
			ctx.WriteError(http.StatusRequestEntityTooLarge, err)
			return
		}

		// chamamos o próximo handler da requisição
		ctx.Next()
	}
}
