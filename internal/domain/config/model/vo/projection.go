package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	configEnum "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/enum"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
)

type Projection struct {
	keys   []string
	values map[string]configEnum.ProjectionValue
}

func (p *Projection) IsEmpty() bool {
	return helper.IsEmpty(p.Keys())
}

func (p *Projection) Exists(key string) bool {
	return helper.Contains(p.keys, key)
}

func (p *Projection) NotExists(key string) bool {
	return p.Exists(key)
}

func (p *Projection) ContainsNumericKey() bool {
	for _, key := range p.keys {
		if helper.IsNumeric(key) {
			return true
		}
	}
	return false
}

func (p *Projection) NotContainsNumericKey() bool {
	return !p.ContainsNumericKey()
}

func (p *Projection) Keys() []string {
	return p.keys
}

func (p *Projection) Get(key string) configEnum.ProjectionValue {
	return p.values[key]
}

func (p *Projection) Type() configEnum.ProjectionType {
	// verificamos se todos os valores são de adição 1
	addition := p.allAddition()
	rejection := p.allRejection()
	if addition && !rejection {
		return configEnum.ProjectionTypeAddition
	} else if rejection && !addition {
		return configEnum.ProjectionTypeRejection
	}
	return configEnum.ProjectionTypeAll
}

func (p *Projection) TypeNumeric() configEnum.ProjectionType {
	// verificamos se todos os valores são de adição 1
	addition := p.allNumericAddition()
	rejection := p.allNumericRejection()
	if addition && !rejection {
		return configEnum.ProjectionTypeAddition
	} else if rejection && !addition {
		return configEnum.ProjectionTypeRejection
	}
	return configEnum.ProjectionTypeAll
}

func (p *Projection) IsAddition(key string) bool {
	return p.Exists(key) && helper.Equals(p.Get(key), configEnum.ProjectionValueAddition)
}

func (p *Projection) IsRejection(key string) bool {
	return p.Exists(key) && helper.Equals(p.Get(key), configEnum.ProjectionValueRejection)
}

func (p *Projection) UnmarshalJSON(data []byte) error {
	if helper.IsEmpty(data) || helper.Equals(strings.TrimSpace(string(data)), "{}") {
		return nil
	}

	iter := jsoniter.ParseString(jsoniter.ConfigFastest, string(data))

	p.keys = []string{}
	p.values = map[string]configEnum.ProjectionValue{}

	for field := iter.ReadObject(); helper.IsNotEmpty(field); field = iter.ReadObject() {
		p.keys = append(p.keys, field)
		p.values[field] = configEnum.ProjectionValue(iter.ReadInt())
	}

	return iter.Error
}

func (p *Projection) MarshalJSON() ([]byte, error) {
	if p.IsEmpty() {
		return []byte("null"), nil
	}
	var data []string
	for _, key := range p.Keys() {
		value := p.values[key]
		data = append(data, fmt.Sprintf("%s:%v", strconv.Quote(key), value))
	}
	obj := fmt.Sprintf("{%s}", strings.Join(data, ","))
	return []byte(obj), nil
}

func (p *Projection) allAddition() bool {
	// passamos campo por campo e se tiver uma rejeição retornamos false
	for _, key := range p.Keys() {
		if helper.IsNotEqualTo(p.Get(key), configEnum.ProjectionValueAddition) {
			return false
		}
	}
	return true
}

func (p *Projection) allNumericAddition() bool {
	// passamos campo por campo e se tiver uma rejeição retornamos false
	for _, key := range p.Keys() {
		if helper.IsNumeric(key) && helper.IsNotEqualTo(p.Get(key), configEnum.ProjectionValueAddition) {
			return false
		}
	}
	return true
}

func (p *Projection) allRejection() bool {
	// passamos campo por campo e se tiver uma adição retornamos false
	for _, key := range p.Keys() {
		if helper.IsNotEqualTo(p.Get(key), configEnum.ProjectionValueRejection) {
			return false
		}
	}
	return true
}

func (p *Projection) allNumericRejection() bool {
	// passamos campo por campo e se tiver uma adição retornamos false
	for _, key := range p.Keys() {
		if helper.IsNumeric(key) && helper.IsNotEqualTo(p.Get(key), configEnum.ProjectionValueRejection) {
			return false
		}
	}
	return true
}
