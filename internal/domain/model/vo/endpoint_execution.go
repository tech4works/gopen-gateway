package vo

import (
	"encoding/json"

	"github.com/tech4works/checker"
)

type EndpointExecution struct {
	allExecuted  bool
	allOK        bool
	degradations []BackendDegradation
}

func NewEmptyEndpointExecution() EndpointExecution {
	return EndpointExecution{}
}

func NewEndpointExecution(allExecuted, allOK bool, degradations []BackendDegradation) EndpointExecution {
	return EndpointExecution{
		allExecuted:  allExecuted,
		allOK:        allOK,
		degradations: degradations,
	}
}

func (e EndpointExecution) AllExecuted() bool {
	return e.allExecuted
}

func (e EndpointExecution) AllOK() bool {
	return e.allOK
}

func (e EndpointExecution) Degradations() []BackendDegradation {
	return e.degradations
}

func (e EndpointExecution) MarshalJSON() ([]byte, error) {
	type Alias EndpointExecution
	return json.Marshal(&struct {
		AllExecuted  bool                 `json:"allExecuted"`
		AllOK        bool                 `json:"allOK"`
		Degradations []BackendDegradation `json:"degradations,omitempty"`
	}{
		AllExecuted:  e.allExecuted,
		AllOK:        e.allOK,
		Degradations: e.degradations,
	})
}

func (e *EndpointExecution) UnmarshalJSON(data []byte) error {
	aux := &struct {
		AllExecuted  bool                 `json:"allExecuted"`
		AllOK        bool                 `json:"allOK"`
		Degradations []BackendDegradation `json:"degradations,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); checker.NonNil(err) {
		return err
	}

	e.allExecuted = aux.AllExecuted
	e.allOK = aux.AllOK
	e.degradations = aux.Degradations

	return nil
}
