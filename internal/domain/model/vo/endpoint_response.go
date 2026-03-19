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
	"github.com/tech4works/checker"
)

type EndpointResponse struct {
	cache       CacheInfo
	degradation Degradation

	execution EndpointExecution

	status   ResponseStatus
	metadata Metadata
	payload  *Payload
}

func NewEndpointResponse(status ResponseStatus, metadata Metadata, payload *Payload) *EndpointResponse {
	return NewEndpointResponseWithAll(
		NewEmptyCacheInfo(),
		NewEmptyDegradation(),
		NewEmptyEndpointExecution(),
		status,
		metadata,
		payload,
	)
}

func NewEndpointResponseWithOnlyStatus(status ResponseStatus) *EndpointResponse {
	return NewEndpointResponse(status, NewEmptyMetadata(), nil)
}

func NewEndpointResponseWithOnlyStatusAndPayload(status ResponseStatus, payload *Payload) *EndpointResponse {
	return NewEndpointResponse(status, NewEmptyMetadata(), payload)
}

func NewEndpointResponseWithExecution(
	degradation Degradation,
	execution EndpointExecution,
	status ResponseStatus,
	metadata Metadata,
	payload *Payload,
) *EndpointResponse {
	return NewEndpointResponseWithAll(NewEmptyCacheInfo(), degradation, execution, status, metadata, payload)
}

func NewEndpointResponseWithAll(
	cache CacheInfo,
	degradation Degradation,
	execution EndpointExecution,
	status ResponseStatus,
	metadata Metadata,
	payload *Payload,
) *EndpointResponse {
	return &EndpointResponse{
		cache:       cache,
		degradation: degradation,
		execution:   execution,
		status:      status,
		metadata:    metadata,
		payload:     payload,
	}
}

func (e *EndpointResponse) Cache() CacheInfo {
	return e.cache
}

func (e *EndpointResponse) ComesFromCache() bool {
	return e.cache.Hit()
}

func (e *EndpointResponse) Degradation() Degradation {
	return e.degradation
}

func (e *EndpointResponse) Execution() EndpointExecution {
	return e.execution
}

func (e *EndpointResponse) Status() ResponseStatus {
	return e.status
}

func (e *EndpointResponse) Metadata() Metadata {
	return e.metadata
}

func (e *EndpointResponse) HasPayload() bool {
	return checker.NonNil(e.payload)
}

func (e *EndpointResponse) Payload() *Payload {
	return e.payload
}
