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

package aggregate

import (
	"sync"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type History struct {
	mu sync.RWMutex

	backends           []*vo.Backend
	httpRequests       []*vo.HTTPBackendRequest
	httpResponses      []*vo.HTTPBackendResponse
	publisherRequests  []*vo.PublisherBackendRequest
	publisherResponses []*vo.PublisherBackendResponse
}

func NewHistoryWithSize(backendsSize int) *History {
	return &History{
		backends:           make([]*vo.Backend, backendsSize),
		httpRequests:       make([]*vo.HTTPBackendRequest, backendsSize),
		httpResponses:      make([]*vo.HTTPBackendResponse, backendsSize),
		publisherRequests:  make([]*vo.PublisherBackendRequest, backendsSize),
		publisherResponses: make([]*vo.PublisherBackendResponse, backendsSize),
	}
}

func (h *History) AddBackend(
	i int,
	backend *vo.Backend,
	httpRequest *vo.HTTPBackendRequest,
	httpResponse *vo.HTTPBackendResponse,
	publisherRequest *vo.PublisherBackendRequest,
	publisherResponse *vo.PublisherBackendResponse,
) {
	h.mu.Lock()
	defer h.mu.Unlock()

	nb := h.backends
	nhr := h.httpRequests
	nhs := h.httpResponses
	npr := h.publisherRequests
	nps := h.publisherResponses

	nb[i] = backend
	nhr[i] = httpRequest
	nhs[i] = httpResponse
	npr[i] = publisherRequest
	nps[i] = publisherResponse
}

func (h *History) GetHTTPBackend(i int) (*vo.Backend, *vo.HTTPBackendRequest, *vo.HTTPBackendResponse) {
	return h.backends[i], h.httpRequests[i], h.httpResponses[i]
}

func (h *History) GetPublisherBackend(i int) (*vo.Backend, *vo.PublisherBackendRequest, *vo.PublisherBackendResponse) {
	return h.backends[i], h.publisherRequests[i], h.publisherResponses[i]
}

func (h *History) GetBackend(i int) (
	*vo.Backend,
	*vo.HTTPBackendRequest,
	*vo.HTTPBackendResponse,
	*vo.PublisherBackendRequest,
	*vo.PublisherBackendResponse,
) {
	return h.backends[i], h.httpRequests[i], h.httpResponses[i], h.publisherRequests[i], h.publisherResponses[i]
}

func (h *History) GetBackendID(i int) string {
	return h.backends[i].ID()
}

func (h *History) GetBackendResponse(i int) vo.BackendPolymorphicResponse {
	if h.backends[i].IsHTTP() {
		return h.httpResponses[i]
	} else if h.backends[i].IsPublisher() {
		return h.publisherResponses[i]
	}
	return nil
}

func (h *History) IsSingleResponse() bool {
	return checker.Equals(h.SizeOfExecuted(), 1)
}

func (h *History) IsMultipleResponses() bool {
	return checker.IsGreaterThan(h.SizeOfExecuted(), 1)
}

func (h *History) BackendResponseLastest() vo.BackendPolymorphicResponse {
	for i := h.Size() - 1; checker.IsGreaterThanOrEqual(i, 0); i-- {
		backendResponse := h.GetBackendResponse(i)
		if backendResponse.Executed() {
			return backendResponse
		}
	}
	return nil
}

func (h *History) Size() int {
	return len(h.backends)
}

func (h *History) SizeOfExecuted() int {
	count := 0
	for i := 0; checker.IsLessThan(i, h.Size()); i++ {
		backendResponse := h.GetBackendResponse(i)
		if backendResponse.Executed() {
			count++
		}
	}
	return count
}

func (h *History) AllOK() bool {
	for i := 0; checker.IsLessThan(i, h.Size()); i++ {
		backendResponse := h.GetBackendResponse(i)
		if backendResponse.Executed() && !backendResponse.OK() {
			return false
		}
	}
	return true
}

func (h *History) AllBackendsExecuted() bool {
	for i := 0; checker.IsLessThan(i, h.Size()); i++ {
		backendResponse := h.GetBackendResponse(i)
		if !backendResponse.Executed() {
			return false
		}
	}
	return true
}

func (h *History) ResponsesMap() (string, error) {
	sliceOfMap := make([]any, h.Size())

	for i := 0; checker.IsLessThan(i, h.Size()); i++ {
		backendResponse := h.GetBackendResponse(i)
		if !backendResponse.Executed() {
			continue
		}

		responseMap, err := backendResponse.Map()
		if checker.NonNil(err) {
			return "", err
		}

		sliceOfMap[i] = responseMap
	}

	return converter.ToStringWithErr(sliceOfMap)
}

func (h *History) ResponsesMapByID() (string, error) {
	byID := make(map[string]any, h.Size())

	for i := 0; checker.IsLessThan(i, h.Size()); i++ {
		backend := h.backends[i]
		if checker.IsNil(backend) || checker.IsEmpty(backend.ID()) {
			continue
		}

		backendResponse := h.GetBackendResponse(i)
		if checker.IsNil(backendResponse) {
			continue
		}

		responseMap, err := backendResponse.Map()
		if checker.NonNil(err) {
			return "", err
		}

		byID[backend.ID()] = responseMap
	}

	return converter.ToStringWithErr(byID)
}
