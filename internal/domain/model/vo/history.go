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
	backends  []*Backend
	requests  []*HTTPBackendRequest
	responses []*HTTPBackendResponse
}

func NewEmptyHistory() *History {
	return &History{}
}

func NewHistory(backends []*Backend, requests []*HTTPBackendRequest, responses []*HTTPBackendResponse) *History {
	return &History{
		backends:  backends,
		requests:  requests,
		responses: responses,
	}
}

func (h *History) Add(backend *Backend, request *HTTPBackendRequest, response *HTTPBackendResponse) *History {
	return &History{
		backends:  append(h.backends, backend),
		requests:  append(h.requests, request),
		responses: append(h.responses, response),
	}
}

func (h *History) Get(i int) (*Backend, *HTTPBackendRequest, *HTTPBackendResponse) {
	return h.backends[i], h.requests[i], h.responses[i]
}

func (h *History) SingleResponse() bool {
	return checker.IsLessThanOrEqual(h.Size(), 1)
}

func (h *History) MultipleResponses() bool {
	return checker.IsGreaterThan(h.Size(), 1)
}

func (h *History) Size() int {
	return len(h.responses)
}

func (h *History) Last() *HTTPBackendResponse {
	return h.responses[len(h.responses)-1]
}

func (h *History) AllOK() bool {
	for _, httpBackendResponse := range h.responses {
		if !httpBackendResponse.OK() {
			return false
		}
	}
	return true
}

func (h *History) Map() (string, error) {
	var sliceOfMap []any
	for _, response := range h.responses {
		responseMap, err := response.Map()
		if checker.NonNil(err) {
			return "", err
		}
		sliceOfMap = append(sliceOfMap, responseMap)
	}
	return converter.ToStringWithErr(sliceOfMap)
}
