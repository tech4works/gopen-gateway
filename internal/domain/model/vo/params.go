package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
)

type Params map[string]string

func NewParams(params map[string]string) Params {
	return params
}

func NewParamsByPath(path string, parentParams Params) Params {
	r := Params{}

	// filtramos os params que contem no path passado
	for key, value := range parentParams {
		paramKey := fmt.Sprint(":", key)
		if helper.ContainsIgnoreCase(path, paramKey) {
			r[key] = value
		}
	}
	return r
}

func (p Params) Set(key, value string) (r Params) {
	r = p.copy()
	r[key] = value
	return r
}

func (p Params) Del(key string) (r Params) {
	r = p.copy()
	delete(r, key)
	return r
}

func (p Params) Get(key string) string {
	value, _ := p[key]
	return value
}

func (p Params) copy() (r Params) {
	r = Params{}
	for key, value := range p {
		r[key] = value
	}
	return r
}
