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
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
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
	return converter.ToString(p.values)
}

func (p Params) Exists(key string) bool {
	_, ok := p.values[key]
	return ok
}
