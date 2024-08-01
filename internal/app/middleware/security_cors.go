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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"net/http"
)

type securityCorsMiddleware struct {
	service service.SecurityCors
}

type SecurityCors interface {
	Do(ctx app.Context)
}

func NewSecurityCors(service service.SecurityCors) SecurityCors {
	return securityCorsMiddleware{
		service: service,
	}
}

func (s securityCorsMiddleware) Do(ctx app.Context) {
	if !ctx.Gopen().HasSecurityCors() {
		ctx.Next()
	} else if err := s.service.ValidateOrigin(ctx.Gopen().SecurityCors(), ctx.Request()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
	} else if err = s.service.ValidateMethod(ctx.Gopen().SecurityCors(), ctx.Request()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
	} else if err = s.service.ValidateHeaders(ctx.Gopen().SecurityCors(), ctx.Request()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
	} else {
		ctx.Next()
	}

}
