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

package domain

import (
	"fmt"
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

const (
	codeErrDynamicValueNotFound            = "DYNAMIC_VALUE_NOT_FOUND"
	codeErrModifierActionNotImplemented    = "MODIFIER_ACTION_NOT_IMPLEMENTED"
	codeErrModifierIncompatibleContentType = "MODIFIER_INCOMPATIBLE_CONTENT_TYPE"
	codeErrLimiterMetadataTooLarge         = "METADATA_TOO_LARGE"
	codeErrLimiterPayloadTooLarge          = "PAYLOAD_TOO_LARGE"
	codeErrLimiterTooManyRequests          = "TOO_MANY_REQUESTS"
	codeErrCacheNotFound                   = "CACHE_NOT_FOUND"
	codeErrEvalGuards                      = "EVAL_GUARDS"
	codeErrJSONPathNotModified             = "JSON_PATH_NOT_MODIFIED"
)

const (
	msgErrDynamicValueNotFound            = "dynamic-value failed: value not found by syntax=%s"
	msgErrModifierActionNotImplemented    = "modifier failed: op=modify kind=%s action=%s not implemented"
	msgErrModifierIncompatibleContentType = "modifier failed: op=%s incompatible payload content-type=%s to modify"
	msgErrLimiterMetadataTooLarge         = "limiter failed: header too large error permitted=%s"
	msgErrLimiterPayloadTooLarge          = "limiter failed: payload too large error permitted=%s"
	msgErrLimiterTooManyRequests          = "limiter failed: too many requests error permitted=%s every=%s"
	msgErrCacheNotFound                   = "cache failed: not found by key=%s"
	msgErrEvalGuards                      = "eval guards: op=eval-guards reason=%s should-run=false"
	msgErrJSONPathNotModified             = "jsonpath failed: op=%s not modified %s"
)

var (
	ErrDynamicValueNotFound            = errors.TargetWithCode(codeErrDynamicValueNotFound)
	ErrModifierActionNotImplemented    = errors.TargetWithCode(codeErrModifierActionNotImplemented)
	ErrModifierIncompatibleContentType = errors.TargetWithCode(codeErrModifierIncompatibleContentType)
	ErrCacheNotFound                   = errors.TargetWithCode(codeErrCacheNotFound)
	ErrLimiterMetadataTooLarge         = errors.TargetWithCode(codeErrLimiterMetadataTooLarge)
	ErrLimiterPayloadTooLarge          = errors.TargetWithCode(codeErrLimiterPayloadTooLarge)
	ErrLimiterTooManyRequests          = errors.TargetWithCode(codeErrLimiterTooManyRequests)
	ErrEvalGuards                      = errors.TargetWithCode(codeErrEvalGuards)
	ErrJSONNotModified                 = errors.TargetWithCode(codeErrJSONPathNotModified)
)

func NewErrDynamicValueNotFound(syntax string) error {
	return errors.NewWithSkipCallerAndCodef(
		2,
		codeErrDynamicValueNotFound,
		msgErrDynamicValueNotFound,
		syntax,
	)
}

func NewErrModifierActionNotImplemented(kind string, action enum.ModifierAction) error {
	return errors.NewWithSkipCallerAndCodef(
		2,
		codeErrModifierActionNotImplemented,
		msgErrModifierActionNotImplemented,
		kind,
		action,
	)
}

func NewErrModifierIncompatibleContentType(op, contentType string) error {
	return errors.NewWithSkipCallerAndCodef(
		2,
		codeErrModifierIncompatibleContentType,
		msgErrModifierIncompatibleContentType,
		op,
		contentType,
	)
}

func NewErrLimiterPayloadTooLarge(limit string) error {
	return errors.NewWithSkipCallerAndCodef(2, codeErrLimiterPayloadTooLarge, msgErrLimiterPayloadTooLarge, limit)
}

func NewErrLimiterMetadataTooLarge(limit string) error {
	return errors.NewWithSkipCallerAndCodef(2, codeErrLimiterMetadataTooLarge, msgErrLimiterMetadataTooLarge, limit)
}

func NewErrLimiterTooManyRequests(capacity int, every time.Duration) error {
	return errors.NewWithSkipCaller(
		2,
		codeErrLimiterTooManyRequests,
		msgErrLimiterTooManyRequests,
		capacity,
		every.String(),
	)
}

func NewErrCacheNotFound(key string) error {
	return errors.NewWithSkipCallerAndCodef(2, codeErrCacheNotFound, msgErrCacheNotFound, key)
}

func NewErrEvalGuards(reason string) error {
	return errors.NewWithSkipCallerAndCodef(2, codeErrEvalGuards, msgErrEvalGuards, reason)
}

func NewErrJSONNotModified(op, path, value string) error {
	msg := fmt.Sprintf("path: %s", path)
	if checker.IsNotEmpty(value) {
		msg += fmt.Sprintf(" value: %s", value)
	}
	return errors.NewWithSkipCallerAndCodef(2, codeErrJSONPathNotModified, msgErrJSONPathNotModified, op, msg)
}
