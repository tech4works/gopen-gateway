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

// timeoutMiddleware is a type that implements the Timeout interface. It is responsible for
// executing the middleware logic for handling timeouts. It gets the timeout duration from the
// endpoint in the given context and creates a new context with that timeout duration. Then it
// launches a goroutine to execute the next middleware and waits for it to finish using a
// finishChan channel. If the execution of the next middleware takes more time than the timeout
// duration, the context's Done channel gets closed and the timeout error is written to the response.
type timeoutMiddleware struct {
}

// Timeout is an interface that defines a method for handling timeouts with a given context.
type Timeout interface {
	// Do is a method that performs a specific action using the provided context.
	// Parameters:
	//   - ctx: the context object used to perform the action
	Do(ctx *api.Context)
}

// NewTimeout is a function that returns a new instance of Timeout, implemented by timeoutMiddleware.
func NewTimeout() Timeout {
	return timeoutMiddleware{}
}

// Do executes the middleware logic for handling timeouts. It gets the timeout duration from the
// endpoint in the given context and creates a new context with that timeout duration. Then it
// launches a goroutine to execute the next middleware and waits for it to finish using a
// finishChan channel. If the execution of the next middleware takes more time than the timeout
// duration, the context's Done channel gets closed and the timeout error is written to the response.
func (t timeoutMiddleware) Do(ctx *api.Context) {
	timeout := ctx.Endpoint().Timeout()

	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), timeout.Time())
	defer cancel()

	ctx.WithContext(timeoutCtx)

	finishChan := make(chan interface{}, 1)
	go func() {
		ctx.Next()
		finishChan <- struct{}{}
	}()
	select {
	case <-finishChan:
	case <-ctx.Done():
		ctx.WriteError(http.StatusGatewayTimeout, errors.New("gateway timeout:", timeout.String()))
	}
}
