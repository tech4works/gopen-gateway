/*
 * Copyright 2024 Tech4Works
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

package interceptor

import (
	"runtime/debug"

	"github.com/tech4works/checker"
	"github.com/tech4works/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type panicRecoveryMiddleware struct {
	log app.MiddlewareLog
}

type PanicRecovery interface {
	Do(ctx app.Context)
}

func NewPanicRecovery(log app.MiddlewareLog) PanicRecovery {
	return panicRecoveryMiddleware{
		log: log,
	}
}

func (p panicRecoveryMiddleware) Do(ctx app.Context) {
	defer func() {
		if r := recover(); checker.NonNil(r) {
			p.log.PrintErrorf(ctx, "%s:%s", r, string(debug.Stack()))

			err := errors.New("Gateway panic error occurred! detail:", r)

			span := trace.SpanFromContext(ctx.Context())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			ctx.WriteError(enum.ResponseStatusInternalError, err)
		}
	}()
	ctx.Next()
}
