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
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type BackendResponse struct {
	kind enum.BackendKind

	cache       CacheInfo
	outcome     enum.BackendOutcome
	degradation Degradation
	duration    time.Duration

	status   ResponseStatus
	metadata Metadata
	payload  *Payload
}

func NewBackendResponse(
	kind enum.BackendKind,
	outcome enum.BackendOutcome,
	duration time.Duration,
	status ResponseStatus,
	metadata Metadata,
	payload *Payload,
) *BackendResponse {
	return NewBackendResponseWithAll(kind, NewEmptyCacheInfo(), outcome, NewEmptyDegradation(), duration, status, metadata, payload)
}

func NewBackendResponseWithDegradation(
	kind enum.BackendKind,
	outcome enum.BackendOutcome,
	degradation Degradation,
	duration time.Duration,
	status ResponseStatus,
	metadata Metadata,
	payload *Payload,
) *BackendResponse {
	return NewBackendResponseWithAll(kind, NewEmptyCacheInfo(), outcome, degradation, duration, status, metadata, payload)
}

func NewBackendResponseWithAll(
	kind enum.BackendKind,
	cache CacheInfo,
	outcome enum.BackendOutcome,
	degradation Degradation,
	duration time.Duration,
	status ResponseStatus,
	metadata Metadata,
	payload *Payload,
) *BackendResponse {
	return &BackendResponse{
		kind:        kind,
		cache:       cache,
		outcome:     outcome,
		degradation: degradation,
		duration:    duration,
		status:      status,
		metadata:    metadata,
		payload:     payload,
	}
}

func (b *BackendResponse) Kind() enum.BackendKind {
	return b.kind
}

func (b *BackendResponse) Cache() CacheInfo {
	return b.cache
}

func (b *BackendResponse) ComesFromCache() bool {
	return b.cache.Hit()
}

func (b *BackendResponse) Outcome() enum.BackendOutcome {
	return b.outcome
}

func (b *BackendResponse) Duration() time.Duration {
	return b.duration
}

func (b *BackendResponse) OK() bool {
	return b.status.OK()
}

func (b *BackendResponse) Ignored() bool {
	return checker.Equals(b.outcome, enum.BackendOutcomeIgnored)
}

func (b *BackendResponse) Cancelled() bool {
	return checker.Equals(b.outcome, enum.BackendOutcomeCancelled)
}

func (b *BackendResponse) Error() bool {
	return checker.Equals(b.outcome, enum.BackendOutcomeError)
}

func (b *BackendResponse) Executed() bool {
	return checker.Equals(b.outcome, enum.BackendOutcomeExecuted)
}

func (b *BackendResponse) Degradation() Degradation {
	return b.degradation
}

func (b *BackendResponse) Degraded() bool {
	return b.Degradation().Any()
}

func (b *BackendResponse) MetadataDegraded() bool {
	return b.Degradation().Has(enum.DegradationKindMetadata)
}

func (b *BackendResponse) PayloadDegraded() bool {
	return b.Degradation().Has(enum.DegradationKindPayload)
}

func (b *BackendResponse) ShouldIgnoreFinalResponseBuild() bool {
	return !b.ShouldInFinalResponse()
}

func (b *BackendResponse) ShouldInFinalResponse() bool {
	return b.Executed() || b.Error()
}

func (b *BackendResponse) Status() ResponseStatus {
	return b.status
}

func (b *BackendResponse) Metadata() Metadata {
	return b.metadata
}

func (b *BackendResponse) HasBody() bool {
	return checker.NonNil(b.payload) && checker.IsGreaterThan(b.payload.Size(), 0)
}

func (b *BackendResponse) Payload() *Payload {
	return b.payload
}

func (b *BackendResponse) Map() (map[string]any, error) {
	var payload any
	if checker.NonNil(b.payload) {
		payloadMap, err := b.payload.Map()
		if checker.NonNil(err) {
			return nil, err
		}
		payload = payloadMap
	}

	baseMap := map[string]any{
		"cache":     b.Cache().Map(),
		"outcome":   b.Outcome(),
		"ok":        b.OK(),
		"ignored":   b.Ignored(),
		"cancelled": b.Cancelled(),
		"error":     b.Error(),
		"executed":  b.Executed(),
		"status":    b.status.Map(),
	}

	switch b.kind {
	case enum.BackendKindHTTP, enum.BackendKindPublisher:
		baseMap["header"] = b.Metadata().Map()
		baseMap["body"] = payload
	default:
		baseMap["metadata"] = b.Metadata().Map()
		baseMap["payload"] = payload
	}

	return baseMap, nil
}
