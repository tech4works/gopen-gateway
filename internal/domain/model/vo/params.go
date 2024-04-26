/*
 * Copyright 2024 Gabriel Cataldo
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
	"fmt"
	"github.com/gin-gonic/gin"
)

// Params is a type alias for a map[string]string. It represents a collection of key-value pairs,
// where the keys are strings and the values are strings.
type Params map[string]string

// NewParams creates a new Params map from the given gin.Params slice.
// Each element in the slice is assigned as a key-value pair in the Params map.
// The resulting Params map is returned.
func NewParams(params gin.Params) Params {
	result := Params{}
	for _, param := range params {
		result[param.Key] = param.Value
	}
	return result
}

// NewParamsByPath filters the params from the parentParams map that contain the keys present in the path string.
// The filtered params are returned as a new Params map.
func NewParamsByPath(path UrlPath, parentParams Params) Params {
	r := Params{}

	// filtramos os params que contem no path passado
	for key, value := range parentParams {
		paramKey := fmt.Sprint(":", key)
		if path.ContainsParam(paramKey) {
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

// Replace replaces the value of a specific key in the Params map with a new value.
// It checks if the key exists in the Params map using the NotExists method.
// If the key does not exist, it returns the original Params map as is.
// If the key exists, it calls the Set method to create a new Params map with the updated value,
// and returns the new Params map.
// The original Params map is not modified.
func (p Params) Replace(key string, value string) Params {
	if p.NotExists(key) {
		return p
	}
	return p.Set(key, value)
}

// Rename renames a key in the Params map.
// It first checks if the oldKey exists in the Params map using the NotExists method.
// If the oldKey does not exist, it returns the original Params map as is.
// If the oldKey exists, it creates a shallow copy of the original Params map using the copy() method,
// assigns the value of the oldKey to the newKey in the copied map, and deletes the oldKey from the copied map.
// The copied Params map with the newKey and without the oldKey is returned.
func (p Params) Rename(oldKey, newKey string) (r Params) {
	if p.NotExists(oldKey) {
		return p
	}
	r = p.copy()
	r[newKey] = r[oldKey]
	delete(r, oldKey)
	return r
}

// Delete removes a key-value pair from the Params map.
// It creates a shallow copy of the original Params map using the copy() method
// and deletes the specified key from the copied map.
// The copied Params map with the key-value pair removed is returned.
func (p Params) Delete(key string) (r Params) {
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

// Exists checks if a key exists in the Params map.
// It returns true if the key exists, otherwise it returns false.
func (p Params) Exists(key string) bool {
	_, ok := p[key]
	return ok
}

// NotExists checks if a key exists in the Params map.
// It returns true if the key does not exist, otherwise it returns false.
func (p Params) NotExists(key string) bool {
	return !p.Exists(key)
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
