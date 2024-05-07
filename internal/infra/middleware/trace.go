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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

// traceMiddleware represents a type responsible for performing tracing logic for a request.
// It contains a traceProvider that implements the infra.TraceProvider interface.
type traceMiddleware struct {
	traceProvider infra.TraceProvider
}

// Trace represents an interface for performing tracing logic for a request.
// The Do method accepts a ctx parameter of type *api.Context to perform the tracing operation.
type Trace interface {
	// Do perform the tracing operation on the given context.
	Do(ctx *api.Context)
}

// NewTrace creates a new Trace instance.
func NewTrace(traceProvider infra.TraceProvider) Trace {
	return traceMiddleware{
		traceProvider: traceProvider,
	}
}

// Do perform the tracing logic for the request.
// It adds the X-Forwarded-For header to the request with the remote address,
// and sets the X-TraceId header if it is not already specified.
// Then it proceeds to the next function in the request.
func (t traceMiddleware) Do(ctx *api.Context) {
	// adicionamos na requisição o X-Forwarded-For
	ctx.AddHeader(consts.XForwardedFor, ctx.RemoteAddr())
	// caso não tenha traceMiddleware id informado, setamos
	if helper.IsEmpty(ctx.HeaderValue(consts.XTraceId)) {
		ctx.SetHeader(consts.XTraceId, t.traceProvider.GenerateTraceId())
	}
	// seguimos para a próxima func da requisição
	ctx.Next()
}
