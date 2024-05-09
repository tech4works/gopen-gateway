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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

type securityCorsMiddleware struct {
	securityCors *vo.SecurityCors
}

type SecurityCors interface {
	Do(ctx *api.Context)
}

func NewSecurityCors(securityCors *vo.SecurityCors) SecurityCors {
	return securityCorsMiddleware{
		securityCors: securityCors,
	}
}

func (s securityCorsMiddleware) Do(ctx *api.Context) {
	// se a configuração não foi feita ja damos próximo
	if helper.IsNil(s.securityCors) {
		ctx.Next()
		return
	}

	// chamamos o objeto de valor para validar se o ip de origem é permitida a partir do objeto de valor fornecido
	if err := s.securityCors.AllowOrigins(ctx.Header().Get(consts.XForwardedFor)); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}
	// chamamos o objeto de valor para validar se o method é permitida a partir do objeto de valor fornecido
	if err := s.securityCors.AllowMethods(ctx.Method()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}
	// chamamos o domínio para validar se o headers fornecido estão permitidas a partir do objeto de valor fornecido
	if err := s.securityCors.AllowHeaders(ctx.Header()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}

	// se tudo ocorreu bem seguimos para o próximo
	ctx.Next()
}
