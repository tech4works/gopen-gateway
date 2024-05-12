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
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
)

// Mapper represents a mapping structure that maps keys to values.
type Mapper struct {
	// keys is a field of the Mapper struct that represents a slice of strings. It stores the keys used for mapping
	// values in the Mapper.
	keys []string
	// values is a field of the Mapper struct that represents a mapping of keys to values.
	// The keys are of type string and the values are also of type string.
	// It stores the values used for mapping in the Mapper struct.
	values map[string]string
}

// IsEmpty checks if the Mapper object is empty.
// It returns true if the Mapper object is empty, false otherwise.
// A Mapper object is considered empty if it has no keys or values.
// This method calls the helper.IsEmpty function with the Mapper object as its argument.
func (m *Mapper) IsEmpty() bool {
	return helper.IsEmpty(m)
}

func (m *Mapper) IsNotEmpty() bool {
	return !m.IsEmpty()
}

// Exists checks if the specified key exists in the Mapper object.
// It returns true if the key exists, false otherwise.
func (m *Mapper) Exists(key string) bool {
	return helper.Contains(m.keys, key)
}

// Keys returns the slice of strings that represents the keys used for mapping values in the Mapper.
// It simply returns the value of the keys field of the Mapper struct.
func (m *Mapper) Keys() []string {
	return m.keys
}

// Get retrieves the value associated with the specified key from the Mapper object.
// It returns the value as a string. If the key does not exist, it returns an empty string.
// It uses the key to access the values map in the Mapper object.
// Example usage:
// value := myMapper.Get("key")
func (m *Mapper) Get(key string) string {
	return m.values[key]
}

// UnmarshalJSON is a method of the Mapper struct that is used to unmarshal JSON data into the Mapper object.
// It takes a byte slice 'data' as input, which represents the JSON data to be unmarshalled.
// The method first checks if the JSON data is empty or if it represents an empty JSON object.
// If it is empty or represents an empty JSON object, the method returns nil, indicating success.
// Otherwise, it uses the jsoniter library to parse the JSON data into a jsoniter.Iterator object.
// It then initializes the keys and values fields of the Mapper object.
// The method iterates over each key-value pair in the JSON object and appends the keys to the keys field,
// and assigns the corresponding values to the values field in the Mapper object.
// Finally, it returns any error that occurred during parsing and unmarshalling of the JSON data.
func (m *Mapper) UnmarshalJSON(data []byte) error {
	if helper.IsEmpty(data) || helper.Equals(strings.TrimSpace(string(data)), "{}") {
		return nil
	}

	iter := jsoniter.ParseString(jsoniter.ConfigFastest, string(data))

	m.keys = []string{}
	m.values = map[string]string{}

	for field := iter.ReadObject(); helper.IsNotEmpty(field); field = iter.ReadObject() {
		m.keys = append(m.keys, field)
		m.values[field] = iter.ReadString()
	}

	return iter.Error
}

// MarshalJSON converts the Mapper object to its JSON representation.
// If the Mapper object is empty, it returns the JSON string "null".
// Otherwise, it constructs a JSON object string by iterating over each key-value pair in the Mapper.
// It quotes both the key and value using strconv.Quote function and constructs the "{key:value}" string.
// The method then joins all the key-value strings in the data slice using strings.Join with "," as the separator.
// Finally, it returns the byte representation of the JSON object and nil error.
func (m *Mapper) MarshalJSON() ([]byte, error) {
	if m.IsEmpty() {
		return []byte("null"), nil
	}

	var data []string
	for _, key := range m.Keys() {
		value := m.values[key]
		data = append(data, fmt.Sprintf("%s:%s", strconv.Quote(key), strconv.Quote(value)))
	}

	return []byte(fmt.Sprintf("{%s}", strings.Join(data, ","))), nil
}
