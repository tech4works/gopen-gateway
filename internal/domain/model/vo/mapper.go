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
	"fmt"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/tech4works/checker"
)

type Mapper struct {
	onlyIf   []string
	ignoreIf []string
	mMap     Map
}

type Map struct {
	keys   []string
	values map[string]string
}

func NewMapper(onlyIf, ignoreIf []string, mMap Map) *Mapper {
	return &Mapper{
		onlyIf:   onlyIf,
		ignoreIf: ignoreIf,
		mMap:     mMap,
	}
}

func (m Mapper) OnlyIf() []string {
	return m.onlyIf
}

func (m Mapper) IgnoreIf() []string {
	return m.ignoreIf
}

func (m Mapper) Map() *Map {
	return &m.mMap
}

func (m *Map) IsEmpty() bool {
	return checker.IsEmpty(m)
}

func (m *Map) IsNotEmpty() bool {
	return !m.IsEmpty()
}

func (m *Map) Exists(key string) bool {
	return checker.Contains(m.keys, key)
}

func (m *Map) Keys() []string {
	return m.keys
}

func (m *Map) Get(key string) string {
	return m.values[key]
}

func (m *Map) UnmarshalJSON(data []byte) error {
	if checker.IsEmpty(data) || checker.Equals(strings.TrimSpace(string(data)), "{}") {
		return nil
	}

	iter := jsoniter.ParseString(jsoniter.ConfigFastest, string(data))

	m.keys = []string{}
	m.values = map[string]string{}

	for field := iter.ReadObject(); checker.IsNotEmpty(field); field = iter.ReadObject() {
		m.keys = append(m.keys, field)
		m.values[field] = iter.ReadString()
	}

	return iter.Error
}

func (m *Map) MarshalJSON() ([]byte, error) {
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
