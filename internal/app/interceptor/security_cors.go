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

package interceptor

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type securityCorsMiddleware struct {
	service service.SecurityCors
}

type SecurityCors interface {
	Do(ctx app.Context)
}

func NewSecurityCors(service service.SecurityCors) SecurityCors {
	return securityCorsMiddleware{
		service: service,
	}
}

func (s securityCorsMiddleware) Do(ctx app.Context) {
	if !ctx.Endpoint().HasSecurityCode() {
		ctx.Next()
		return
	}

	err := s.service.Validate(ctx.Endpoint().SecurityCors(), ctx.Request())
	if checker.NonNil(err) && errors.IsNot(err, domain.ErrEvalGuards) {
		ctx.WriteError(enum.ResponseStatusPermissionDenied, err)
		return
	}

	metadata, err := s.service.BuildResponseMetadataByConfig(ctx.Endpoint().SecurityCors(), ctx.Request())
	if checker.NonNil(err) {
		ctx.WriteError(enum.ResponseStatusInternalError, err)
		return
	}

	ctx.WriteMetadata(metadata)
	ctx.Next()
}
