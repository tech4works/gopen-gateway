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

package vo

import (
	"github.com/tech4works/checker"
)

type SecurityCorsConfig struct {
	onlyIf           []string
	ignoreIf         []string
	allowOrigins     []string
	allowMethods     []string
	allowHeaders     []string
	allowCredentials bool
}

func NewSecurityCorsConfig(onlyIf, ignoreIf, allowsOrigins, allowMethods, allowHeaders []string, allowCredentials bool,
) *SecurityCorsConfig {
	return &SecurityCorsConfig{
		onlyIf:           onlyIf,
		ignoreIf:         ignoreIf,
		allowOrigins:     allowsOrigins,
		allowMethods:     allowMethods,
		allowHeaders:     allowHeaders,
		allowCredentials: allowCredentials,
	}
}

func (s SecurityCorsConfig) OnlyIf() []string {
	return s.onlyIf
}

func (s SecurityCorsConfig) IgnoreIf() []string {
	return s.ignoreIf
}

func (s SecurityCorsConfig) AllowMethods() []string {
	return s.allowMethods
}

func (s SecurityCorsConfig) AllowHeaders() []string {
	return s.allowHeaders
}

func (s SecurityCorsConfig) AllowCredentials() bool {
	return s.allowCredentials
}

func (s SecurityCorsConfig) DisallowOrigin(origin string) bool {
	return checker.IsNotEmpty(s.allowOrigins) && checker.NotContains(s.allowOrigins, origin)
}

func (s SecurityCorsConfig) DisallowMethod(method string) bool {
	return checker.IsNotEmpty(s.allowMethods) && checker.NotContains(s.allowMethods, method)
}

func (s SecurityCorsConfig) DisallowHeader(headerKey string) bool {
	return checker.IsNotEmpty(s.allowHeaders) && checker.NotContains(s.allowHeaders, headerKey)
}
