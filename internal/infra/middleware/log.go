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
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"time"
)

// logMiddleware is a struct that represents a logging component. It uses an implementation of the LogProvider interface
// to perform logging operations.
type logMiddleware struct {
	logProvider infra.LogProvider
}

// Log is an interface that defines a logging operation. Implementation of this interface
// should provide a method Do() that takes a context object of type *api.Context as an argument.
type Log interface {
	// Do perform a logging operation using the provided context object.
	// The context object is of type *api.Context and is passed as an argument to the method.
	Do(ctx *api.Context)
}

// NewLog creates a new instance of the Log interface using the provided LogProvider.
func NewLog(logProvider infra.LogProvider) Log {
	return logMiddleware{
		logProvider: logProvider,
	}
}

// Do is a method that performs logging for a request.
// It keeps track of the request start time, initializes the logger options with traceMiddleware ID and XForwardedFor,
// prints the start logMiddleware, calls the next request handler, and prints the finish logMiddleware.
// It takes a *api.Context as a parameter.
func (l logMiddleware) Do(ctx *api.Context) {
	// mantemos o tempo que a requisição começou
	startTime := time.Now()

	// inicializamos a logger options global, com o traceId e XForwardedFor
	l.logProvider.InitializeLoggerOptions(ctx)

	// imprimimos o logMiddleware de start
	logger.Info("Start!", l.logProvider.BuildStartRequestMessage(ctx))

	// chamamos o próximo handler da requisição
	ctx.Next()

	// imprimimos o logMiddleware de finish
	logger.Info("Finish!", l.logProvider.BuildFinishRequestMessage(ctx.HttpResponse(), startTime))
}
