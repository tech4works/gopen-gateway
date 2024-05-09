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
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

type timeoutMiddleware struct {
}

type Timeout interface {
	Do(ctx *api.Context)
}

// NewTimeout returns a new instance of the `timeoutMiddleware` type that implements the `Timeout` interface.
func NewTimeout() Timeout {
	return timeoutMiddleware{}
}

func (t timeoutMiddleware) Do(ctx *api.Context) {
	// instanciamos a configuração do endpoint de timeout para aplicar
	timeout := ctx.Endpoint().Timeout()

	// inicializamos o context com timeoutMiddleware fornecido na config do gateway
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), timeout.Time())
	defer cancel()

	// setamos esse context na request atual para propagar para os outros manipuladores
	ctx.RequestWithContext(timeoutCtx)

	// criamos os canais de alerta
	finishChan := make(chan interface{}, 1)

	go func() {
		// chamamos o próximo handler na requisição
		ctx.Next()
		// se finalizou a tempo, chamamos o channel para seguir normalmente
		finishChan <- struct{}{}
	}()

	// seguramos o goroutine principal aguardando os canais ou o context serem notificados
	select {
	case <-finishChan:
	case <-ctx.Done():
		ctx.WriteError(http.StatusGatewayTimeout, errors.New("gateway timeout:", timeout.String()))
	}
}
