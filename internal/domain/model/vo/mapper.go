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

type Mapper struct {
	keys   []string
	values map[string]string
}

func (m *Mapper) IsEmpty() bool {
	return helper.IsEmpty(m)
}

func (m *Mapper) IsNotEmpty() bool {
	return !m.IsEmpty()
}

func (m *Mapper) Exists(key string) bool {
	return helper.Contains(m.keys, key)
}

func (m *Mapper) Keys() []string {
	return m.keys
}

func (m *Mapper) Get(key string) string {
	return m.values[key]
}

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
