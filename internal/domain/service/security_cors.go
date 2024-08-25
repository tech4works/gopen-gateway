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

package service

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"strings"
)

type securityCorsService struct {
}

type SecurityCors interface {
	ValidateOrigin(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error
	ValidateMethod(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error
	ValidateHeaders(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error
}

func NewSecurityCors() SecurityCors {
	return securityCorsService{}
}

func (s securityCorsService) ValidateOrigin(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error {
	if securityCors.DisallowOrigin(request.ClientIP()) {
		return errors.New("Origin not mapped on security-cors.allow-origins")
	}
	return nil
}

func (s securityCorsService) ValidateMethod(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error {
	if securityCors.DisallowMethod(request.Method()) {
		return errors.New("Method not mapped on security-cors.allow-methods")
	}
	return nil
}

func (s securityCorsService) ValidateHeaders(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error {
	var headersNotAllowed []string
	for _, key := range request.Header().Keys() {
		if checker.NotEquals(key, mapper.XForwardedFor) && securityCors.DisallowHeader(key) {
			headersNotAllowed = append(headersNotAllowed, key)
		}
	}

	if checker.IsNotEmpty(headersNotAllowed) {
		keys := strings.Join(headersNotAllowed, ", ")
		return errors.Newf("Headers contain not mapped fields on security-cors.allow-headers: %s", keys)
	}

	return nil
}
