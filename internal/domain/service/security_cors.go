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
	"strings"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type securityCors struct {
	dynamicValueService DynamicValue
}

type SecurityCors interface {
	Validate(config *vo.SecurityCorsConfig, request *vo.EndpointRequest) error
	BuildResponseMetadataByConfig(config *vo.SecurityCorsConfig, request *vo.EndpointRequest) (vo.Metadata, error)
}

func NewSecurityCors(dynamicValueService DynamicValue) SecurityCors {
	return securityCors{
		dynamicValueService: dynamicValueService,
	}
}

func (s securityCors) Validate(config *vo.SecurityCorsConfig, request *vo.EndpointRequest) error {
	if !request.IsCORS() {
		return nil
	}

	err := s.evalSecurityCorsGuards(config, request)
	if checker.NonNil(err) {
		return err
	}

	origin := request.Metadata().Get("Origin")
	if config.DisallowOrigin(origin) {
		return errors.New("security-cors failed: origin not mapped on security-cors.allow-origins")
	}

	method := request.Operation()
	if request.IsPreflight() {
		method = request.Metadata().Get("Access-Control-Request-Method")
	}

	if config.DisallowMethod(method) {
		return errors.New("security-cors failed: method not mapped on security-cors.allow-methods")
	}

	if request.IsPreflight() {
		var notAllowed []string
		for _, h := range request.Metadata().GetAll("Access-Control-Request-Headers") {
			if config.DisallowHeader(h) {
				notAllowed = append(notAllowed, h)
			}
		}
		if checker.IsNotEmpty(notAllowed) {
			return errors.Newf("security-cors failed: request-headers contain not mapped fields on security-cors.allow-headers=%s",
				strings.Join(notAllowed, ", "))
		}
	}

	return nil
}

func (s securityCors) BuildResponseMetadataByConfig(config *vo.SecurityCorsConfig, request *vo.EndpointRequest) (
	vo.Metadata, error) {
	if !request.IsCORS() {
		return vo.NewEmptyMetadata(), nil
	}

	err := s.evalSecurityCorsGuards(config, request)
	if checker.NonNil(err) && errors.IsNot(err, domain.ErrEvalGuards) {
		return vo.NewEmptyMetadata(), err
	}

	mapMetadata := map[string][]string{
		"Access-Control-Allow-Origin": request.Metadata().GetAll("Origin"),
		"Vary":                        {"Origin"},
	}

	if config.AllowCredentials() {
		mapMetadata["Access-Control-Allow-Credentials"] = converter.ToSlice("true")
	}

	if request.IsPreflight() {
		mapMetadata["Access-Control-Allow-Methods"] = request.Metadata().GetAll("Access-Control-Request-Method")

		requestedHeaders := request.Metadata().GetAll("Access-Control-Request-Headers")
		if checker.IsNotEmpty(requestedHeaders) {
			mapMetadata["Access-Control-Allow-Headers"] = requestedHeaders
		}
	}

	return vo.NewMetadata(mapMetadata), nil
}

func (s securityCors) evalSecurityCorsGuards(config *vo.SecurityCorsConfig, request *vo.EndpointRequest) error {
	errs := s.dynamicValueService.EvalGuardsWithErr(config.OnlyIf(), config.IgnoreIf(), request, nil)
	if errors.Only(errs, domain.ErrEvalGuards) {
		return errs[0]
	} else if checker.IsNotEmpty(errs) {
		return errors.JoinInheritf(errs, ", ", "failed to evaluate guard for security-cors")
	} else {
		return nil
	}
}
