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
	"github.com/tech4works/gopen-gateway/internal/app"
	"net/http"
)

type panicRecoveryMiddleware struct {
}

type PanicRecovery interface {
	Do(ctx app.Context)
}

func NewPanicRecovery() PanicRecovery {
	return panicRecoveryMiddleware{}
}

func (p panicRecoveryMiddleware) Do(ctx app.Context) {
	defer func() {
		if r := recover(); helper.IsNotNil(r) {
			//todo
			// p.logger.PrintEndpointErrorf(ctx, "%s:%s", r, string(debug.Stack()))
			ctx.WriteError(http.StatusInternalServerError, errors.New("gateway panic error occurred! detail:", r))
		}
	}()
	ctx.Next()
}
