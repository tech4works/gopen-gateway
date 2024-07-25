package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
)

type History struct {
	backends  []*Backend
	responses []*HTTPBackendResponse
}

func NewEmptyHistory() *History {
	return &History{}
}

func NewHistory(backends []*Backend, responses []*HTTPBackendResponse) *History {
	return &History{
		backends:  backends,
		responses: responses,
	}
}

func (h *History) Add(backend *Backend, response *HTTPBackendResponse) *History {
	return &History{
		backends:  append(h.backends, backend),
		responses: append(h.responses, response),
	}
}

func (h *History) Get(i int) (*Backend, *HTTPBackendResponse) {
	return h.backends[i], h.responses[i]
}

func (h *History) SingleResponse() bool {
	return helper.IsLessThanOrEqual(h.Size(), 1)
}

func (h *History) MultipleResponses() bool {
	return helper.IsGreaterThan(h.Size(), 1)
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
		if helper.IsNotNil(err) {
			return "", err
		}
		sliceOfMap = append(sliceOfMap, responseMap)
	}
	return helper.ConvertToString(sliceOfMap)
}
