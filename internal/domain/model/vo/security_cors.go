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

package vo

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"strings"
)

// SecurityCors represents the configuration options for Cross-Origin Resource Sharing (CORS) settings in Gopen.
type SecurityCors struct {
	// allowOrigins is a field in the SecurityCors struct that represents a list of allowed origins for Cross-Origin
	// Resource Sharing (CORS).
	allowOrigins []string
	// allowMethods is a field in the SecurityCors struct that represents a list of allowed HTTP methods for Cross-Origin
	// Resource Sharing (CORS).
	allowMethods []string
	// allowHeaders is a field in the SecurityCors struct that represents a list of allowed HTTP headers for Cross-Origin
	// Resource Sharing (CORS).
	allowHeaders []string
}

func newSecurityCors(securityCorsDTO *dto.SecurityCors) *SecurityCors {
	if helper.IsNil(securityCorsDTO) {
		return nil
	}
	return &SecurityCors{
		allowOrigins: securityCorsDTO.AllowOrigins,
		allowMethods: securityCorsDTO.AllowMethods,
		allowHeaders: securityCorsDTO.AllowHeaders,
	}
}

// AllowOriginsData returns the allowOrigins field in the SecurityCors struct.
func (s SecurityCors) AllowOriginsData() []string {
	return s.allowOrigins
}

// AllowMethodsData returns the allowMethods field in the SecurityCors struct.
func (s SecurityCors) AllowMethodsData() []string {
	return s.allowMethods
}

// AllowHeadersData returns the allowHeaders field in the SecurityCors struct.
func (s SecurityCors) AllowHeadersData() []string {
	return s.allowHeaders
}

// AllowOrigins checks if the given IP is allowed based on the allowOrigins field in the SecurityCors struct.
// If the allowOrigins field is not empty and does not contain the given IP, it returns an error "Origin not mapped
// on security-cors.allow-origins". Otherwise, it returns nil.
func (s SecurityCors) AllowOrigins(ip string) (err error) {
	// verificamos se na configuração security-cors.allow-origins ta vazia, ou tá informado o ip da requisição
	if helper.IsNotEmpty(s.allowOrigins) && helper.NotContains(s.allowOrigins, ip) {
		err = errors.New("Origin not mapped on security-cors.allow-origins")
	}
	return err
}

// AllowMethods checks if the given HTTP method is allowed based on the allowMethods field in the SecurityCors struct.
// If the allowMethods field is not empty and does not contain the given method, it returns an error "Method not mapped
// on security-cors.allow-methods". Otherwise, it returns nil.
func (s SecurityCors) AllowMethods(method string) (err error) {
	// verificamos se na configuração security-cors.allow-methods ta vazia ou tá informado com oq vem na requisição
	if helper.IsNotEmpty(s.allowMethods) && helper.NotContains(s.allowMethods, method) {
		err = errors.New("Method not mapped on security-cors.allow-methods")
	}
	return err
}

// AllowHeaders checks if the requested HTTP headers are allowed based on the allowHeaders field in the SecurityCors struct.
// If the allowHeaders field is empty, it returns nil.
// Otherwise, it iterates over the headers of the request and adds any headers that are not mapped in the allowHeaders list,
// except for the X-Forwarded-For and X-Trace-Id headers.
// If there are headers that are not allowed, it returns an error "Headers contains not mapped fields on security-cors.allow-headers"
// along with the list of headers that are not allowed.
// If all headers are allowed, it returns nil.
func (s SecurityCors) AllowHeaders(header Header) (err error) {
	// verificamos se na configuração security-cors.allow-headers ta vazia
	if helper.IsEmpty(s.allowHeaders) {
		return nil
	}
	// inicializamos os headers não permitidos
	var headersNotAllowed []string
	// iteramos o header da requisição para verificar os headers que contain
	for key := range header {
		// caso o campo do header não esteja mapeado na lista security-cors.allow-headers e nao seja X-Forwarded-For
		// e nem X-Trace-Id adicionamos na lista
		if helper.NotContains(s.allowHeaders, key) && helper.IsNotEqualToIgnoreCase(key, consts.XForwardedFor) &&
			helper.IsNotEqualToIgnoreCase(key, consts.XTraceId) {
			headersNotAllowed = append(headersNotAllowed, key)
		}
	}
	// caso a lista não esteja vazia, quer dizer que tem headers não permitidos
	if helper.IsNotEmpty(headersNotAllowed) {
		headersFields := strings.Join(headersNotAllowed, ", ")
		return errors.New("Headers contains not mapped fields on security-cors.allow-headers:", headersFields)
	}

	// se tudo ocorreu bem retornamos nil
	return nil
}
