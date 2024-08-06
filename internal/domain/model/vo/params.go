package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tech4works/checker"
)

type Params struct {
	values map[string]string
}

func NewParams(values map[string]string) Params {
	cleanValues := map[string]string{}
	for key, value := range values {
		if checker.IsNotEmpty(values) {
			cleanValues[key] = value
		}
	}
	return Params{values: cleanValues}
}

func (p Params) Get(key string) string {
	return p.values[key]
}

func (p Params) Keys() (ss []string) {
	for key := range p.values {
		ss = append(ss, key)
	}
	return ss
}

func (p Params) Map() any {
	return p.Copy()
}

func (p Params) Copy() map[string]string {
	copiedMap := map[string]string{}
	for key, value := range p.values {
		copiedMap[key] = value
	}
	return copiedMap
}

func (p Params) Length() int {
	return len(p.values)
}

func (p Params) IsEmpty() bool {
	return checker.Equals(p.Length(), 0)
}

func (p Params) String() string {
	return helper.SimpleConvertToString(p.values)
}

func (p Params) Exists(key string) bool {
	_, ok := p.values[key]
	return ok
}
