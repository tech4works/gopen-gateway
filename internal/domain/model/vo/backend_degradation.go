package vo

import "encoding/json"

type BackendDegradation struct {
	id          string
	degradation Degradation
}

func NewBackendDegradation(id string, degradation Degradation) BackendDegradation {
	return BackendDegradation{
		id:          id,
		degradation: degradation,
	}
}

func (b BackendDegradation) ID() string {
	return b.id
}

func (b BackendDegradation) Degradation() Degradation {
	return b.degradation
}

func (b BackendDegradation) Degraded() bool {
	return b.Degradation().Any()
}

func (b BackendDegradation) MarshalJSON() ([]byte, error) {
	type backendDegradationJSON struct {
		ID          string      `json:"id"`
		Degradation Degradation `json:"degradation"`
	}

	return json.Marshal(backendDegradationJSON{
		ID:          b.id,
		Degradation: b.degradation,
	})
}

func (b *BackendDegradation) UnmarshalJSON(data []byte) error {
	type backendDegradationJSON struct {
		ID          string      `json:"id"`
		Degradation Degradation `json:"degradation"`
	}

	var aux backendDegradationJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	b.id = aux.ID
	b.degradation = aux.Degradation
	return nil
}
