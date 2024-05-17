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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

// logMiddleware is a struct that represents a middleware for logging HTTP requests and responses.
// It contains an instance of HttpLoggerProvider interface to handle HTTP logging.
type logMiddleware struct {
	logger interfaces.LoggerProvider
}

// Log is an interface that represents a logging operation.
// The Do method is responsible for performing the logging operation using the provided context object.
// The context object is of type *api.Context and is passed as an argument to the method.
type Log interface {
	// Do perform a logging operation using the provided context object.
	// The context object is of type *api.Context and is passed as an argument to the method.
	Do(ctx *api.Context)
}

// NewLog is a function that creates a new instance of logMiddleware struct and
// returns it as a value of Log interface.
func NewLog(loggerProvider interfaces.LoggerProvider) Log {
	return logMiddleware{
		logger: loggerProvider,
	}
}

func (l logMiddleware) Do(ctx *api.Context) {
	ctx.Next()

	l.logger.PrintEndpointResponseInfo(ctx)
	l.logger.PrintHttpResponseInfo(ctx)
}
