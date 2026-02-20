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
	"github.com/tech4works/converter"
)

type History struct {
	backends           []*Backend
	httpRequests       []*HTTPBackendRequest
	httpResponses      []*HTTPBackendResponse
	publisherRequests  []*PublisherBackendRequest
	publisherResponses []*PublisherBackendResponse
}

func NewHistoryWithSize(backendsSize int) *History {
	return &History{
		backends:           make([]*Backend, backendsSize),
		httpRequests:       make([]*HTTPBackendRequest, backendsSize),
		httpResponses:      make([]*HTTPBackendResponse, backendsSize),
		publisherRequests:  make([]*PublisherBackendRequest, backendsSize),
		publisherResponses: make([]*PublisherBackendResponse, backendsSize),
	}
}

func (h *History) AddBackend(
	i int,
	backend *Backend,
	httpRequest *HTTPBackendRequest,
	httpResponse *HTTPBackendResponse,
	publisherRequest *PublisherBackendRequest,
	publisherResponse *PublisherBackendResponse,
) *History {
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

	return &History{
		backends:           nb,
		httpRequests:       nhr,
		httpResponses:      nhs,
		publisherRequests:  npr,
		publisherResponses: nps,
	}
}

func (h *History) GetHTTPBackend(i int) (*Backend, *HTTPBackendRequest, *HTTPBackendResponse) {
	return h.backends[i], h.httpRequests[i], h.httpResponses[i]
}

func (h *History) GetPublisherBackend(i int) (*Backend, *PublisherBackendRequest, *PublisherBackendResponse) {
	return h.backends[i], h.publisherRequests[i], h.publisherResponses[i]
}

func (h *History) GetBackend(i int) (*Backend, *HTTPBackendRequest, *HTTPBackendResponse, *PublisherBackendRequest, *PublisherBackendResponse) {
	return h.backends[i], h.httpRequests[i], h.httpResponses[i], h.publisherRequests[i], h.publisherResponses[i]
}

func (h *History) GetBackendID(i int) string {
	return h.backends[i].ID()
}

func (h *History) GetBackendResponse(i int) BackendPolymorphicResponse {
	if checker.NonNil(h.httpResponses[i]) {
		return h.httpResponses[i]
	} else if checker.NonNil(h.publisherResponses[i]) {
		return h.publisherResponses[i]
	}
	return nil
}

func (h *History) IsSingleResponse() bool {
	return checker.Equals(h.NormalSize(), 1)
}

func (h *History) IsMultipleResponses() bool {
	return checker.IsGreaterThan(h.NormalSize(), 1)
}

func (h *History) BackendResponseLastest() BackendPolymorphicResponse {
	for i := h.Size() - 1; checker.IsGreaterThanOrEqual(i, 0); i-- {
		if checker.NonNil(h.httpResponses[i]) {
			return h.httpResponses[i]
		} else if checker.NonNil(h.publisherResponses[i]) {
			return h.publisherResponses[i]
		}
	}
	return nil
}

func (h *History) Size() int {
	return len(h.backends)
}

func (h *History) NormalSize() int {
	count := 0
	for i := 0; checker.IsLessThan(i, h.Size()); i++ {
		backend := h.backends[i]
		if checker.IsNil(backend) || !backend.IsNormal() {
			continue
		}
		count++
	}
	return count
}

func (h *History) AllOK() bool {
	for i := 0; checker.IsLessThan(i, h.Size()); i++ {
		backendResponse := h.GetBackendResponse(i)
		if checker.IsNil(backendResponse) {
			continue
		} else if !backendResponse.OK() {
			return false
		}
	}
	return true
}

func (h *History) AllBackendsExecuted() bool {
	for i := 0; checker.IsLessThan(i, h.Size()); i++ {
		backendResponse := h.GetBackendResponse(i)
		if checker.IsNil(backendResponse) {
			return false
		}
	}
	return true
}

func (h *History) ResponsesMap() (string, error) {
	sliceOfMap := make([]any, h.Size())

	for i := 0; checker.IsLessThan(i, h.Size()); i++ {
		backendResponse := h.GetBackendResponse(i)
		if checker.IsNil(backendResponse) {
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
