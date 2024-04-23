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

// securityCors implements the SecurityCors interface.
// It represents the configuration options for Cross-Origin Resource Sharing (CORS) settings in Gopen.
type securityCors struct {
	securityCorsVO *vo.SecurityCors
}

// SecurityCors is an interface that defines the behavior of handling Cross-Origin Resource Sharing (CORS) settings in Gopen.
// Implementing types must define the Do method, which takes a *api.Context as an argument.
// The Do method is responsible for handling CORS-related tasks, such as validating and processing CORS requests.
type SecurityCors interface {
	// Do handles Cross-Origin Resource Sharing (CORS) related tasks.
	// It takes a *api.Context as an argument and is responsible for validating and processing CORS requests.
	Do(ctx *api.Context)
}

// NewSecurityCors is a function that creates a new instance of SecurityCors with the given securityCorsVO configuration.
// It returns the new SecurityCors object.
func NewSecurityCors(securityCorsVO *vo.SecurityCors) SecurityCors {
	return securityCors{
		securityCorsVO: securityCorsVO,
	}
}

// Do is a method that handles the Cross-Origin Resource Sharing (CORS) settings in Gopen.
// It allows or denies access to resources on a web page from another domain.
// It checks if the configuration is set. If not, it proceeds to the next middleware.
// It validates if the origin IP is allowed based on the provided SecurityCors configuration.
// It validates if the HTTP method is allowed based on the provided SecurityCors configuration.
// It validates if the headers are allowed based on the provided SecurityCors configuration.
// If any validation fails, it returns a forbidden error.
// If all validations pass, it proceeds to the next middleware.
// The method takes a Context as the input parameter.
// The method does not return anything.
func (c securityCors) Do(ctx *api.Context) {
	// se a configuração não foi feita ja damos próximo
	if helper.IsNil(c.securityCorsVO) {
		ctx.Next()
		return
	}

	// chamamos o objeto de valor para validar se o ip de origem é permitida a partir do objeto de valor fornecido
	if err := c.securityCorsVO.AllowOrigins(ctx.HeaderValue(consts.XForwardedFor)); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}
	// chamamos o objeto de valor para validar se o method é permitida a partir do objeto de valor fornecido
	if err := c.securityCorsVO.AllowMethods(ctx.Method()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}
	// chamamos o domínio para validar se o headers fornecido estão permitidas a partir do objeto de valor fornecido
	if err := c.securityCorsVO.AllowHeaders(ctx.Header()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}

	// se tudo ocorreu bem seguimos para o próximo
	ctx.Next()
}
