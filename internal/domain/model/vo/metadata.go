/*
 * Copyright 2024 Tech4Works
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
	"strings"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
)

type Metadata struct {
	values map[string][]string
}

func NewEmptyMetadata() Metadata {
	return Metadata{values: map[string][]string{}}
}

func NewMetadata(values map[string][]string) Metadata {
	cleanValues := map[string][]string{}
	for k, v := range values {
		if checker.IsEmpty(v) {
			continue
		}

		nk := normalizeKey(k)

		if existing, ok := cleanValues[nk]; ok {
			cleanValues[nk] = append(existing, v...)
		} else {
			cleanValues[nk] = v
		}
	}
	return Metadata{values: cleanValues}
}

func (m Metadata) String() string {
	return converter.ToCompactString(m.values)
}

func (m Metadata) GetAll(key string) []string {
	return m.values[normalizeKey(key)]
}

func (m Metadata) Get(key string) string {
	valuesByKey := m.values[normalizeKey(key)]
	if checker.IsNotEmpty(valuesByKey) {
		return strings.Join(valuesByKey, ", ")
	}
	return ""
}

func (m Metadata) GetFirst(key string) string {
	valuesByKey := m.values[normalizeKey(key)]
	if checker.IsNotEmpty(valuesByKey) {
		return valuesByKey[0]
	}
	return ""
}

func (m Metadata) Exists(key string) bool {
	_, ok := m.values[normalizeKey(key)]
	return ok
}

func (m Metadata) NotExists(key string) bool {
	return !m.Exists(key)
}

func (m Metadata) Copy() map[string][]string {
	copiedValues := map[string][]string{}
	for key, value := range m.values {
		copiedValues[key] = value
	}
	return copiedValues
}

func (m Metadata) Map() any {
	return m.Copy()
}

func (m Metadata) Size() int {
	size := 0
	for key, values := range m.values {
		size += len(key) + 2
		for _, value := range values {
			size += len(value)
			size += 2
		}
		size -= 2
		size += 2
	}
	size += 2
	return size
}

func (m Metadata) SizeStr() string {
	bs := NewBytesByInt(m.Size())
	return converter.ToCompactString(bs.String())
}

func (m Metadata) Keys() (ss []string) {
	for key := range m.values {
		ss = append(ss, key)
	}
	return ss
}

func (m Metadata) MarshalJSON() ([]byte, error) {
	return converter.ToBytesWithErr(m.Map())
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	var values map[string][]string

	err := converter.ToDestWithErr(data, &values)
	if checker.NonNil(err) {
		return err
	}

	m.values = values
	return nil
}

func normalizeKey(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
}
