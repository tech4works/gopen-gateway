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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
)

// Projection represents a data structure that stores keys and values of a projection.
type Projection struct {
	// keys represents a slice of strings that store the keys of a Projection.
	keys []string

	// values map[string]enum.ProjectionValue
	values map[string]enum.ProjectionValue
}

// IsEmpty checks if the Projection is empty by checking if the keys are empty. Returns true if the Projection is empty,
// false otherwise.
func (p *Projection) IsEmpty() bool {
	return helper.IsEmpty(p.Keys())
}

// Exists checks if the key exists in the projection's keys. Returns true if the key exists, false otherwise.
func (p *Projection) Exists(key string) bool {
	return helper.Contains(p.keys, key)
}

// ContainsNumericKey checks if the Projection contains any numeric key. Returns true if a numeric key is found,
// false otherwise.
func (p *Projection) ContainsNumericKey() bool {
	for _, key := range p.keys {
		if helper.IsNumeric(key) {
			return true
		}
	}
	return false
}

// NotContainsNumericKey checks if the Projection does not contain any numeric key.
// Returns true if no numeric key is found, false otherwise.
func (p *Projection) NotContainsNumericKey() bool {
	return !p.ContainsNumericKey()
}

// Keys returns the keys of the Projection, which represent the data structure that stores keys and values
// of a projection.
func (p *Projection) Keys() []string {
	return p.keys
}

// Get returns the value of a projection for the given key.
// The key is used to fetch the corresponding value from the `values` map in the `Projection` struct.
// If the key does not exist in the `values` map, it returns the zero value of the `ProjectionValue` type.
func (p *Projection) Get(key string) enum.ProjectionValue {
	return p.values[key]
}

// Type returns the ProjectionType of the Projection based on the values of its keys.
// If all values are ProjectionTypeAddition and no values are ProjectionTypeRejection, it returns ProjectionTypeAddition.
// If all values are ProjectionTypeRejection and no values are ProjectionTypeAddition, it returns ProjectionTypeRejection.
// Otherwise, it returns ProjectionTypeAll.
func (p *Projection) Type() enum.ProjectionType {
	addition := p.allAddition()
	rejection := p.allRejection()
	if addition && !rejection {
		return enum.ProjectionTypeAddition
	} else if rejection && !addition {
		return enum.ProjectionTypeRejection
	}
	return enum.ProjectionTypeAll
}

// TypeNumeric returns the ProjectionType of the Projection based on the values of its numeric keys.
// If all the numeric keys have ProjectionTypeAddition and no numeric keys have ProjectionTypeRejection,
// it returns ProjectionTypeAddition.
// If all the numeric keys have ProjectionTypeRejection and no numeric keys have ProjectionTypeAddition,
// it returns ProjectionTypeRejection.
// Otherwise, it returns ProjectionTypeAll.
func (p *Projection) TypeNumeric() enum.ProjectionType {
	addition := p.allNumericAddition()
	rejection := p.allNumericRejection()
	if addition && !rejection {
		return enum.ProjectionTypeAddition
	} else if rejection && !addition {
		return enum.ProjectionTypeRejection
	}
	return enum.ProjectionTypeAll
}

// IsAddition checks if the key exists in the projection's keys and if the value for that key is enum.ProjectionValueAddition.
// Returns true if the key exists and the value is enum.ProjectionValueAddition, false otherwise.
func (p *Projection) IsAddition(key string) bool {
	return p.Exists(key) && helper.Equals(p.Get(key), enum.ProjectionValueAddition)
}

// IsRejection checks if the key exists in the projection's keys and if the value for that key is enum.ProjectionValueRejection.
// Returns true if the key exists and the value is enum.ProjectionValueRejection, false otherwise.
func (p *Projection) IsRejection(key string) bool {
	return p.Exists(key) && helper.Equals(p.Get(key), enum.ProjectionValueRejection)
}

// UnmarshalJSON unmarshals the JSON data into the Projection struct.
// It populates the keys and values of the Projection based on the JSON data.
// If the JSON data is empty or represents an empty object, it returns nil.
// Otherwise, it parses the JSON data using jsoniter and iterates through the fields.
// For each field, it appends the key to the keys slice and assigns the corresponding
// ProjectionValue to the values map. Finally, it returns any error that occurred during parsing.
func (p *Projection) UnmarshalJSON(data []byte) error {
	if helper.IsEmpty(data) || helper.Equals(strings.TrimSpace(string(data)), "{}") {
		return nil
	}

	iter := jsoniter.ParseString(jsoniter.ConfigFastest, string(data))

	p.keys = []string{}
	p.values = map[string]enum.ProjectionValue{}
	for field := iter.ReadObject(); helper.IsNotEmpty(field); field = iter.ReadObject() {
		p.keys = append(p.keys, field)
		p.values[field] = enum.ProjectionValue(iter.ReadInt())
	}

	return iter.Error
}

// MarshalJSON converts a Projection to its JSON representation.
// If the Projection is empty, it returns a JSON null value.
// Otherwise, it iterates through the keys and values of the Projection,
// formats them as key-value pairs in a JSON object, and returns the JSON byte array.
// The keys are quoted using strconv.Quote to ensure valid JSON syntax.
// The values are appended as string representations.
// Any error during the process will be returned.
func (p *Projection) MarshalJSON() ([]byte, error) {
	if p.IsEmpty() {
		return []byte("null"), nil
	}

	var data []string
	for _, key := range p.Keys() {
		value := p.values[key]
		data = append(data, fmt.Sprintf("%s:%v", strconv.Quote(key), value))
	}

	return []byte(fmt.Sprintf("{%s}", strings.Join(data, ","))), nil
}

// allAddition checks if all values in the Projection are enum.ProjectionValueAddition.
// It iterates through each key in the Projection, and if any value is not equal to enum.ProjectionValueAddition,
// it returns false. If all values are equal to enum.ProjectionValueAddition, it returns true.
func (p *Projection) allAddition() bool {
	for _, key := range p.Keys() {
		if helper.IsNotEqualTo(p.Get(key), enum.ProjectionValueAddition) {
			return false
		}
	}
	return true
}

// allNumericAddition checks if all numeric keys in the Projection have a value equal to enum.ProjectionValueAddition.
// It iterates through each numeric key in the Projection, and if any value is not equal to enum.ProjectionValueAddition,
// it returns false. If all values for numeric keys are equal to enum.ProjectionValueAddition, it returns true.
func (p *Projection) allNumericAddition() bool {
	for _, key := range p.Keys() {
		if helper.IsNumeric(key) && helper.IsNotEqualTo(p.Get(key), enum.ProjectionValueAddition) {
			return false
		}
	}
	return true
}

// allRejection checks if all values in the Projection are enum.ProjectionValueRejection.
// It iterates through each key in the Projection, and if any value is not equal to enum.ProjectionValueRejection,
// it returns false. If all values are equal to enum.ProjectionValueRejection, it returns true.
func (p *Projection) allRejection() bool {
	for _, key := range p.Keys() {
		if helper.IsNotEqualTo(p.Get(key), enum.ProjectionValueRejection) {
			return false
		}
	}
	return true
}

// allNumericRejection checks if all numeric keys in the Projection have a value equal to enum.ProjectionValueRejection.
// It iterates through each numeric key in the Projection, and if any value is not equal to enum.ProjectionValueRejection,
// it returns false. If all values for numeric keys are equal to enum.ProjectionValueRejection, it returns true.
func (p *Projection) allNumericRejection() bool {
	for _, key := range p.Keys() {
		if helper.IsNumeric(key) && helper.IsNotEqualTo(p.Get(key), enum.ProjectionValueRejection) {
			return false
		}
	}
	return true
}
