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

// newSecurityCors creates a new SecurityCors object based on the given SecurityCorsJson object.
// If the securityCorsJson is nil, it returns nil.
// Otherwise, it initializes a new SecurityCors object with the allowOrigins, allowMethods, and allowHeaders fields
// populated with the corresponding values from the SecurityCorsJson object.
// The newly created SecurityCors object is returned as a pointer.
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

// AllowOrigins checks if the given IP is allowed based on the allowOrigins field in the SecurityCors struct.
// If the allowOrigins field is not empty and does not contain the given IP, it returns an error "Origin not mapped
// on security-cors.allow-origins". Otherwise, it returns nil.
func (s SecurityCors) AllowOrigins(ip string) (err error) {
	if helper.IsNotEmpty(s.allowOrigins) && helper.NotContains(s.allowOrigins, ip) {
		err = errors.New("Origin not mapped on security-cors.allow-origins")
	}
	return err
}

// AllowMethods checks if the given HTTP method is allowed based on the allowMethods field in the SecurityCors struct.
// If the allowMethods field is not empty and does not contain the given method, it returns an error "Method not mapped
// on security-cors.allow-methods". Otherwise, it returns nil.
func (s SecurityCors) AllowMethods(method string) (err error) {
	if helper.IsNotEmpty(s.allowMethods) && helper.NotContains(s.allowMethods, method) {
		err = errors.New("Method not mapped on security-cors.allow-methods")
	}
	return err
}

// AllowHeaders checks if the given HTTP headers are allowed based on the allowHeaders field in the SecurityCors struct.
// If the allowHeaders field is empty, it returns nil.
// Otherwise, it iterates through the keys in the header map and checks if the key is not equal to consts.XForwardedFor,
// consts.XTraceId, and if it is not contained in the allowHeaders list.
// If any headers are found that are not allowed, it appends them to the headersNotAllowed slice.
// Finally, if headersNotAllowed is not empty, it joins the values into a comma-separated string and returns an error
// "Headers contains not mapped fields on security-cors.allow-headers: <comma-separated headers>".
// Otherwise, it returns nil.
func (s SecurityCors) AllowHeaders(header Header) (err error) {
	if helper.IsEmpty(s.allowHeaders) {
		return nil
	}

	var headersNotAllowed []string
	for key := range header {
		if helper.IsNotEqualTo(key, consts.XForwardedFor) && helper.IsNotEqualToIgnoreCase(key, consts.XTraceId) &&
			helper.NotContains(s.allowHeaders, key) {
			headersNotAllowed = append(headersNotAllowed, key)
		}
	}
	if helper.IsNotEmpty(headersNotAllowed) {
		headersFields := strings.Join(headersNotAllowed, ", ")
		err = errors.New("Headers contains not mapped fields on security-cors.allow-headers:", headersFields)
	}

	return err
}
