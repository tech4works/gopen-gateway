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
	"github.com/GabrielHCataldo/go-helper/helper"
)

type SecurityCors struct {
	allowOrigins []string
	allowMethods []string
	allowHeaders []string
}

func newSecurityCors(securityCorsJson *SecurityCorsJson) *SecurityCors {
	if helper.IsNil(securityCorsJson) {
		return nil
	}
	return &SecurityCors{
		allowOrigins: securityCorsJson.AllowOrigins,
		allowMethods: securityCorsJson.AllowMethods,
		allowHeaders: securityCorsJson.AllowHeaders,
	}
}
func (s SecurityCors) IsValidOrigin(origin string) bool {
	return helper.IsEmpty(s.allowOrigins) || helper.Contains(s.allowOrigins, origin)
}

func (s SecurityCors) IsValidMethod(method string) bool {
	return helper.IsEmpty(s.allowMethods) || helper.Contains(s.allowMethods, method)
}

func (s SecurityCors) IsValidHeader(header string) bool {
	return helper.IsEmpty(s.allowHeaders) || helper.Contains(s.allowHeaders, header)
}
