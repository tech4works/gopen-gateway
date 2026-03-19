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

	"github.com/tech4works/converter"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type BackendCacheEntry struct {
	Kind        enum.BackendKind    `json:"kind"`
	Outcome     enum.BackendOutcome `json:"outcome"`
	Degradation Degradation         `json:"degradation"`
	Status      ResponseStatus      `json:"status"`
	Metadata    Metadata            `json:"metadata"`
	Payload     *Payload            `json:"payload,omitempty"`
	TTL         Duration            `json:"ttl"`
	CreatedAt   time.Time           `json:"createdAt"`
}

func NewBackendCacheEntry(config *CacheConfig, response *BackendResponse) *BackendCacheEntry {
	return &BackendCacheEntry{
		Kind:        response.Kind(),
		Outcome:     response.Outcome(),
		Degradation: response.Degradation(),
		Status:      response.Status(),
		Metadata:    response.Metadata(),
		Payload:     response.Payload(),
		TTL:         config.TTL(),
		CreatedAt:   time.Now(),
	}
}

func (b BackendCacheEntry) Entry() (string, error) {
	return converter.ToCompactStringWithErr(b)
}

func (b BackendCacheEntry) Response(duration time.Duration) *BackendResponse {
	return NewBackendResponseWithAll(
		b.Kind,
		NewCacheInfo(true, b.TTL, NewDuration(b.TTL.Time()-time.Since(b.CreatedAt))),
		b.Outcome,
		b.Degradation,
		duration,
		b.Status,
		b.Metadata,
		b.Payload,
	)
}
