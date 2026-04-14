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
	"github.com/tech4works/converter"
)

type EndpointCacheEntry struct {
	Degradation Degradation       `json:"degradation,omitempty"`
	Execution   EndpointExecution `json:"execution,omitempty"`
	Status      ResponseStatus    `json:"status"`
	Metadata    Metadata          `json:"metadata"`
	Payload     *Payload          `json:"payload,omitempty"`
	TTL         Duration          `json:"ttl"`
	CreatedAt   time.Time         `json:"createdAt"`
}

func NewEndpointCacheEntry(config *CacheConfig, response *EndpointResponse) *EndpointCacheEntry {
	return &EndpointCacheEntry{
		Degradation: response.Degradation(),
		Execution:   response.Execution(),
		Status:      response.Status(),
		Metadata:    response.Metadata(),
		Payload:     response.Payload(),
		TTL:         config.TTL(),
		CreatedAt:   time.Now(),
	}
}

func (e EndpointCacheEntry) Entry() (string, error) {
	return converter.ToCompactStringWithErr(e)
}

func (e EndpointCacheEntry) IsZero() bool {
	return checker.IsEmpty(e.Status.Value()) && checker.IsNil(e.Metadata.values) && checker.IsNil(e.Payload)
}

func (e EndpointCacheEntry) Response() *EndpointResponse {
	if e.IsZero() {
		return nil
	}
	return NewEndpointResponseWithAll(
		NewCacheInfo(true, e.TTL, NewDuration(e.TTL.Time()-time.Since(e.CreatedAt))),
		e.Degradation,
		e.Execution,
		e.Status,
		e.Metadata,
		e.Payload,
	)
}
