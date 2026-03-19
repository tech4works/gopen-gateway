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

	backends  []*vo.BackendConfig
	responses []*vo.BackendResponse
}

func NewHistoryWithSize(backendsSize int) *History {
	return &History{
		backends:  make([]*vo.BackendConfig, backendsSize),
		responses: make([]*vo.BackendResponse, backendsSize),
	}
}

func (h *History) Add(i int, backend *vo.BackendConfig, response *vo.BackendResponse) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.backends[i] = backend
	h.responses[i] = response
}

func (h *History) Get(i int) (*vo.BackendConfig, *vo.BackendResponse) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.backends[i], h.responses[i]
}

func (h *History) GetID(i int) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.backends[i].ID()
}

func (h *History) GetResponse(i int) *vo.BackendResponse {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.responses[i]
}

func (h *History) GetResponseLastest() *vo.BackendResponse {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for i := h.sizeUnlocked() - 1; checker.IsGreaterThanOrEqual(i, 0); i-- {
		resp := h.responses[i]
		if checker.NonNil(resp) {
			return resp
		}
	}
	return nil
}

func (h *History) GetLatest() (*vo.BackendConfig, *vo.BackendResponse) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for i := h.sizeUnlocked() - 1; checker.IsGreaterThanOrEqual(i, 0); i-- {
		resp := h.responses[i]
		if checker.NonNil(resp) {
			return h.backends[i], resp
		}
	}
	return nil, nil
}

func (h *History) DependenciesSatisfied(backend *vo.BackendConfig) bool {
	if !backend.HasDependencies() {
		return true
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, idx := range backend.Dependencies().Indexes() {
		resp := h.responses[idx]
		if checker.IsNil(resp) || !resp.Executed() {
			return false
		}
	}

	return true
}

func (h *History) IsSingleFinalResponse() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return checker.Equals(h.sizeForFinalResponseUnlocked(), 1)
}

func (h *History) IsMultipleFinalResponse() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return checker.IsGreaterThan(h.sizeForFinalResponseUnlocked(), 1)
}

func (h *History) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.sizeUnlocked()
}

func (h *History) Degradations() []vo.BackendDegradation {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var degradations []vo.BackendDegradation
	for i := 0; checker.IsLessThan(i, h.sizeForFinalResponseUnlocked()); i++ {
		resp := h.responses[i]
		if checker.NonNil(resp) && resp.ShouldInFinalResponse() && resp.Degraded() {
			degradations = append(degradations, vo.NewBackendDegradation(h.backends[i].ID(), resp.Degradation()))
		}
	}
	return degradations
}

func (h *History) AllOK() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for i := 0; checker.IsLessThan(i, h.sizeUnlocked()); i++ {
		resp := h.responses[i]
		if checker.IsNil(resp) || !resp.OK() {
			return false
		}
	}
	return true
}

func (h *History) AllExecuted() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for i := 0; checker.IsLessThan(i, h.sizeUnlocked()); i++ {
		resp := h.responses[i]
		if checker.IsNil(resp) || !resp.Executed() {
			return false
		}
	}
	return true
}

func (h *History) ResponsesMap() (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sliceOfMap := make([]any, h.sizeUnlocked())

	for i := 0; checker.IsLessThan(i, h.sizeUnlocked()); i++ {
		resp := h.responses[i]
		if checker.IsNil(resp) {
			continue
		}

		responseMap, err := resp.Map()
		if checker.NonNil(err) {
			return "", err
		}

		sliceOfMap[i] = responseMap
	}

	return converter.ToStringWithErr(sliceOfMap)
}

func (h *History) ResponsesMapByID() (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	byID := make(map[string]any, h.sizeUnlocked())

	for i := 0; checker.IsLessThan(i, h.sizeUnlocked()); i++ {
		backend := h.backends[i]
		if checker.IsNil(backend) {
			continue
		}

		resp := h.responses[i]
		if checker.IsNil(resp) {
			continue
		}

		responseMap, err := resp.Map()
		if checker.NonNil(err) {
			return "", err
		}

		byID[backend.ID()] = responseMap
	}

	return converter.ToStringWithErr(byID)
}

func (h *History) sizeUnlocked() int {
	return len(h.backends)
}

func (h *History) sizeForFinalResponseUnlocked() int {
	count := 0
	for i := 0; checker.IsLessThan(i, h.sizeUnlocked()); i++ {
		resp := h.responses[i]
		if checker.NonNil(resp) && resp.ShouldInFinalResponse() {
			count++
		}
	}
	return count
}
