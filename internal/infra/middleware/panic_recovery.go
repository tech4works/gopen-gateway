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
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
	"runtime/debug"
)

// panicRecoveryMiddleware is a type that represents a middleware implementation for recovering from panics and handling
// the recovery process.
type panicRecoveryMiddleware struct {
}

// PanicRecovery represents a type that can recover from panics.
// It provides a method called Do that takes a *api.Context as a parameter and does the necessary recovery actions.
type PanicRecovery interface {
	// Do apply the necessary recovery actions for a given *api.Context parameter.
	Do(ctx *api.Context)
}

// NewPanicRecovery returns a PanicRecovery implementation.
// The returned PanicRecovery contains a Do method that recovers from panics and handles the recovery process.
// The Do method takes a *api.Context as a parameter and performs the necessary recovery actions.
// The recovery actions include logging the recovered panic and stack traceMiddleware, and writing an error response to the context.
// The Do method also calls ctx.Next() to proceed to the next request handling.
func NewPanicRecovery() PanicRecovery {
	return panicRecoveryMiddleware{}
}

// Do recovers from panics and handles the recovery process.
// It takes an *api.Context as a parameter and performs the necessary recovery actions.
// The recovery actions include logging the recovered panic and stack traceMiddleware,
// and writing an error response to the context.
// It also calls ctx.Next() to proceed to the next request handling.
func (p panicRecoveryMiddleware) Do(ctx *api.Context) {
	defer func() {
		if r := recover(); helper.IsNotNil(r) {
			logger.Errorf("%s:%s", r, string(debug.Stack()))
			err := errors.New("gateway panic error occurred! detail:", r)
			ctx.WriteError(http.StatusInternalServerError, err)
		}
	}()
	ctx.Next()
}
