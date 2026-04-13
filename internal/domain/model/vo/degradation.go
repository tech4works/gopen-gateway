package vo

import (
	"encoding/json"

	"github.com/tech4works/checker"

	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type Degradation struct {
	kinds []enum.DegradationKind
}

func NewEmptyDegradation() Degradation {
	return NewDegradation()
}

func NewDegradation(kinds ...enum.DegradationKind) Degradation {
	return Degradation{
		kinds: kinds,
	}
}

func (d Degradation) Has(kind enum.DegradationKind) bool {
	return checker.IsNotEmpty(d.kinds) && checker.Contains(d.kinds, kind)
}

func (d Degradation) Any() bool {
	return checker.IsNotEmpty(d.kinds)
}

func (d Degradation) IsNone() bool {
	return checker.IsEmpty(d.kinds)
}

func (d Degradation) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.kinds)
}

func (d *Degradation) UnmarshalJSON(data []byte) error {
	var kinds []enum.DegradationKind

	if err := json.Unmarshal(data, &kinds); checker.NonNil(err) {
		return err
	}

	d.kinds = kinds
	return nil
}
