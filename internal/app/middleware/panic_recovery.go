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
	"github.com/tech4works/checker"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app"
	"go.elastic.co/apm/v2"
	"net/http"
	"runtime/debug"
)

type panicRecoveryMiddleware struct {
	log app.EndpointLog
}

type PanicRecovery interface {
	Do(ctx app.Context)
}

func NewPanicRecovery(log app.EndpointLog) PanicRecovery {
	return panicRecoveryMiddleware{
		log: log,
	}
}

func (p panicRecoveryMiddleware) Do(ctx app.Context) {
	defer func() {
		if r := recover(); checker.NonNil(r) {
			p.printErrorf(ctx, "%s:%s", r, string(debug.Stack()))

			err := errors.New("Gateway panic error occurred! detail:", r)

			tx := apm.TransactionFromContext(ctx.Context())
			if checker.NonNil(tx) {
				apmErr := apm.DefaultTracer().NewError(err)
				apmErr.SetTransaction(tx)
				apmErr.Send()
			}

			ctx.WriteError(http.StatusInternalServerError, err)
		}
	}()
	ctx.Next()
}

func (p panicRecoveryMiddleware) printErrorf(ctx app.Context, format string, msg ...any) {
	p.log.PrintErrorf(ctx.Endpoint(), ctx.Request(), ctx.ClientIP(), ctx.TraceID(), format, msg...)
}
