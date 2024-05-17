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

// securityCorsMiddleware represents a middleware for handling Cross-Origin Resource Sharing (CORS) settings in Gopen.
// It checks if the request IP, HTTP method, and request headers are allowed based on the configuration provided in the
// SecurityCors struct.
// If any of these checks fail, it returns a Forbidden error response. Otherwise, it passes the request to the next
// middleware.
type securityCorsMiddleware struct {
	securityCors *vo.SecurityCors
}

// SecurityCors represents an interface for handling Cross-Origin Resource Sharing (CORS) settings in a system.
// Its methods Do take in a context object and should be implemented to perform the necessary CORS checks for a request.
// If the request fails any of the checks, it should handle the response accordingly, e.g., by returning a Forbidden error.
// Otherwise, it should pass the request to the next middleware.
type SecurityCors interface {
	// Do perform the necessary CORS checks for a request.
	// If the request fails any of the checks, it should handle the response accordingly.
	// Otherwise, it should pass the request to the next middleware.
	Do(ctx *api.Context)
}

// NewSecurityCors creates a new instance of the SecurityCors interface with the provided SecurityCors configuration.
// The returned instance is of type securityCorsMiddleware, which is a middleware for handling Cross-Origin Resource
// Sharing (CORS) settings.
func NewSecurityCors(securityCors *vo.SecurityCors) SecurityCors {
	return securityCorsMiddleware{
		securityCors: securityCors,
	}
}

// Do is a method of the securityCorsMiddleware type used for handling Cross-Origin Resource Sharing (CORS) settings.
// It checks if the request IP, HTTP method, and request headers are allowed based on the configuration provided in the
// SecurityCors struct.
// If any of these checks fail, it returns a Forbidden error response. Otherwise, it passes the request to the
// next middleware.
func (s securityCorsMiddleware) Do(ctx *api.Context) {
	if helper.IsNil(s.securityCors) {
		ctx.Next()
		return
	}

	if err := s.securityCors.AllowOrigins(ctx.Header().Get(consts.XForwardedFor)); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	} else if err = s.securityCors.AllowMethods(ctx.Method()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	} else if err = s.securityCors.AllowHeaders(ctx.Header()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}

	ctx.Next()
}
