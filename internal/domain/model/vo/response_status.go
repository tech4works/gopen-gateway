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
	"encoding/json"
	"fmt"

	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type ResponseStatus struct {
	value       enum.ResponseStatus
	raw         any
	description string
}

func NewResponseStatus(value enum.ResponseStatus, raw any, description string) ResponseStatus {
	return ResponseStatus{
		value:       value,
		raw:         raw,
		description: description,
	}
}

func NewResponseStatusByValue(value enum.ResponseStatus) ResponseStatus {
	return NewResponseStatus(value, nil, "")
}

func (s ResponseStatus) Value() enum.ResponseStatus {
	return s.value
}

func (s ResponseStatus) HasRaw() bool {
	return checker.NonNil(s.raw)
}

func (s ResponseStatus) Raw() any {
	return s.raw
}

func (s ResponseStatus) Description() string {
	return s.description
}

func (s ResponseStatus) OK() bool {
	return checker.Equals(s.value, enum.ResponseStatusOK)
}

func (s ResponseStatus) Failed() bool {
	return !s.OK()
}

func (s ResponseStatus) ClientError() bool {
	switch s.value {
	case enum.ResponseStatusInvalidArgument,
		enum.ResponseStatusNotFound,
		enum.ResponseStatusAlreadyExists,
		enum.ResponseStatusPermissionDenied,
		enum.ResponseStatusUnauthenticated,
		enum.ResponseStatusFailedPrecondition,
		enum.ResponseStatusOutOfRange,
		enum.ResponseStatusConflict:
		return true
	default:
		return false
	}
}

func (s ResponseStatus) ServerError() bool {
	switch s.value {
	case enum.ResponseStatusInternalError,
		enum.ResponseStatusUnavailable,
		enum.ResponseStatusDeadlineExceeded,
		enum.ResponseStatusDataLoss,
		enum.ResponseStatusResourceExhausted,
		enum.ResponseStatusUnimplemented:
		return true
	default:
		return false
	}
}

func (s ResponseStatus) Retryable() bool {
	switch s.value {
	case enum.ResponseStatusUnavailable,
		enum.ResponseStatusDeadlineExceeded,
		enum.ResponseStatusMetadataTooLarge,
		enum.ResponseStatusPayloadTooLarge,
		enum.ResponseStatusResourceExhausted:
		return true
	default:
		return false
	}
}

func (s ResponseStatus) String() string {
	return fmt.Sprintf("status=%s raw=%v description=%s", s.value, s.raw, s.description)
}

func (s ResponseStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status      string `json:"status"`
		Raw         any    `json:"raw"`
		Description string `json:"description,omitempty"`
	}{
		Status:      string(s.value),
		Raw:         s.raw,
		Description: s.description,
	})
}

func (s *ResponseStatus) UnmarshalJSON(data []byte) error {
	var aux struct {
		Status      string `json:"status"`
		Raw         any    `json:"raw"`
		Description string `json:"description,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); checker.NonNil(err) {
		return err
	}

	s.value = enum.ResponseStatus(aux.Status)
	s.raw = aux.Raw
	s.description = aux.Description

	return nil
}

func (s ResponseStatus) Map() any {
	return map[string]any{
		"value":       s.value,
		"raw":         s.raw,
		"description": s.description,
	}
}
