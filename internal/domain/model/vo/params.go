package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
)

type Params map[string]string

// NewParamsByPath filters the params from the parentParams map that contain the keys present in the path string.
// The filtered params are returned as a new Params map.
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

// Set assigns a value to a specific key in the Params map.
// It creates a shallow copy of the original Params map using the copy() method
// and assigns the new value to the specified key in the copied map.
// The copied Params map with the updated value is returned.
func (p Params) Set(key, value string) (r Params) {
	r = p.copy()
	r[key] = value
	return r
}

// Del removes a key-value pair from the Params map.
// It creates a shallow copy of the original Params map using the copy() method
// and deletes the specified key from the copied map.
// The copied Params map with the key-value pair removed is returned.
func (p Params) Del(key string) (r Params) {
	r = p.copy()
	delete(r, key)
	return r
}

// Get retrieves the value associated with the specified key from the Params map.
// It returns the value if found, otherwise it returns an empty string.
func (p Params) Get(key string) string {
	value, _ := p[key]
	return value
}

// copy creates a shallow copy of the Params map.
// It iterates over each key-value pair in the original Params map and assigns them to the new copy.
// The copied Params map is then returned.
func (p Params) copy() (r Params) {
	r = Params{}
	for key, value := range p {
		r[key] = value
	}
	return r
}
