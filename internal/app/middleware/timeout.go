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

package middleware

import (
	"context"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app"
	"net/http"
)

type timeoutMiddleware struct {
}

type Timeout interface {
	Do(ctx app.Context)
}

func NewTimeout() Timeout {
	return timeoutMiddleware{}
}

func (t timeoutMiddleware) Do(ctx app.Context) {
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
