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
	"github.com/tech4works/checker"
)

type SecurityCors struct {
	allowOrigins []string
	allowMethods []string
	allowHeaders []string
}

func NewSecurityCors(allowsOrigins, allowMethods, allowHeaders []string) *SecurityCors {
	return &SecurityCors{
		allowOrigins: allowsOrigins,
		allowMethods: allowMethods,
		allowHeaders: allowHeaders,
	}
}

func (s SecurityCors) DisallowOrigin(origin string) bool {
	return checker.NotContains(s.allowOrigins, origin)
}

func (s SecurityCors) DisallowMethod(method string) bool {
	return checker.NotContains(s.allowMethods, method)
}

func (s SecurityCors) DisallowHeader(headerKey string) bool {
	return checker.NotContains(s.allowHeaders, headerKey)
}
