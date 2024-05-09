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

type limiterMiddleware struct {
}

type Limiter interface {
	Do(ctx *api.Context)
}

func NewLimiter() Limiter {
	return limiterMiddleware{}
}

func (l limiterMiddleware) Do(ctx *api.Context) {
	// instanciamos o limiter que esta no endpoint
	limiter := ctx.Endpoint().Limiter()

	// validamos com o objeto de valor se a requisição está dentro do limite de taxa permitido
	err := limiter.Rate().Allow(ctx.HttpRequest().Header().Get(consts.XForwardedFor))
	if helper.IsNotNil(err) {
		ctx.WriteError(http.StatusTooManyRequests, err)
		return
	}

	// validamos com o objeto de valor se a requisição está dentro do tamanho permitido
	err = limiter.Allow(ctx.HttpRequest())
	if errors.Contains(err, mapper.ErrPayloadTooLarge) {
		ctx.WriteError(http.StatusRequestEntityTooLarge, err)
		return
	} else if errors.Contains(err, mapper.ErrHeaderTooLarge) {
		ctx.WriteError(http.StatusRequestHeaderFieldsTooLarge, err)
		return
	}

	// se tudo ocorreu bem vamos para o próximo manipulador
	ctx.Next()
}
