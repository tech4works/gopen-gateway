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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

type timeoutMiddleware struct {
}

// Timeout represents an interface for implementing timeoutMiddleware functionality. It has a single method Do
// that takes a timeoutMiddleware duration and returns a HandlerFunc.
type Timeout interface {
	// Do is a method of the Timeout interface that takes a timeoutMiddleware duration and returns a HandlerFunc.
	// It is used to implement timeoutMiddleware functionality for HTTP route handlers.
	Do(timeoutDuration vo.Duration) api.HandlerFunc
}

// NewTimeout returns a new instance of the `timeoutMiddleware` type that implements the `Timeout` interface.
func NewTimeout() Timeout {
	return timeoutMiddleware{}
}

// Do execute a timeoutMiddleware handler for an HTTP request.
// It initializes the context with the timeoutMiddleware provided in the gateway configuration.
// The timeoutMiddleware context is set in the current request to propagate it to other handlers.
// It creates finishChan and starts a goroutine to call the next handler in the request.
// If the execution finishes on time, it signals the finishChan channel.
// The main goroutine waits for either the finishChan channel or the context to be notified.
// If the timeoutMiddleware is reached, it writes a gateway timeoutMiddleware error to the response.
// If the execution finishes before the timeoutMiddleware, it breaks from the select block.
func (t timeoutMiddleware) Do(timeoutDuration vo.Duration) api.HandlerFunc {
	return func(ctx *api.Context) {
		// inicializamos o context com timeoutMiddleware fornecido na config do gateway
		timeoutContext, cancel := context.WithTimeout(ctx.Context(), timeoutDuration.Time())
		defer cancel()

		// setamos esse context na request atual para propagar para os outros manipuladores
		ctx.RequestWithContext(timeoutContext)

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
		case <-ctx.Context().Done():
			err := errors.New("gateway timeoutMiddleware:", timeoutDuration.String())
			ctx.WriteError(http.StatusGatewayTimeout, err)
		}
	}
}
