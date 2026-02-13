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
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type Projector struct {
	onlyIf   []string
	ignoreIf []string
	project  Project
}

type Project struct {
	keys   []string
	values map[string]enum.ProjectValue
}

func NewProjector(onlyIf, ignoreIf []string, project Project) *Projector {
	return &Projector{
		onlyIf:   onlyIf,
		ignoreIf: ignoreIf,
		project:  project,
	}
}

func (p Projector) OnlyIf() []string {
	return p.onlyIf
}

func (p Projector) IgnoreIf() []string {
	return p.ignoreIf
}

func (p Projector) Project() *Project {
	return &p.project
}

func (p *Project) IsEmpty() bool {
	return checker.IsEmpty(p.Keys())
}

func (p *Project) Exists(key string) bool {
	return checker.Contains(p.keys, key)
}

func (p *Project) ContainsNumericKey() bool {
	for _, key := range p.keys {
		if checker.IsNumeric(key) {
			return true
		}
	}
	return false
}

func (p *Project) NotContainsNumericKey() bool {
	return !p.ContainsNumericKey()
}

func (p *Project) Keys() []string {
	return p.keys
}

func (p *Project) Get(key string) enum.ProjectValue {
	return p.values[key]
}

func (p *Project) Kind() enum.ProjectKind {
	addition := p.allAddition()
	rejection := p.allRejection()
	if addition && !rejection {
		return enum.ProjectKindAddition
	} else if rejection && !addition {
		return enum.ProjectKindRejection
	}
	return enum.ProjectKindAll
}

func (p *Project) NumericKind() enum.ProjectKind {
	addition := p.allNumericAddition()
	rejection := p.allNumericRejection()
	if addition && !rejection {
		return enum.ProjectKindAddition
	} else if rejection && !addition {
		return enum.ProjectKindRejection
	}
	return enum.ProjectKindAll
}

func (p *Project) IsAddition(key string) bool {
	return p.Exists(key) && checker.Equals(p.Get(key), enum.ProjectValueAddition)
}

func (p *Project) IsRejection(key string) bool {
	return p.Exists(key) && checker.Equals(p.Get(key), enum.ProjectValueRejection)
}

func (p *Project) UnmarshalJSON(data []byte) error {
	if checker.IsEmpty(data) || checker.Equals(strings.TrimSpace(string(data)), "{}") {
		return nil
	}

	iter := jsoniter.ParseString(jsoniter.ConfigFastest, string(data))

	p.keys = []string{}
	p.values = map[string]enum.ProjectValue{}
	for field := iter.ReadObject(); checker.IsNotEmpty(field); field = iter.ReadObject() {
		p.keys = append(p.keys, field)
		p.values[field] = enum.ProjectValue(iter.ReadInt())
	}

	return iter.Error
}

func (p *Project) MarshalJSON() ([]byte, error) {
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

func (p *Project) allAddition() bool {
	for _, key := range p.Keys() {
		if checker.NotEquals(p.Get(key), enum.ProjectValueAddition) {
			return false
		}
	}
	return true
}

func (p *Project) allNumericAddition() bool {
	for _, key := range p.Keys() {
		if checker.IsNumeric(key) && checker.NotEquals(p.Get(key), enum.ProjectValueAddition) {
			return false
		}
	}
	return true
}

func (p *Project) allRejection() bool {
	for _, key := range p.Keys() {
		if checker.NotEquals(p.Get(key), enum.ProjectValueRejection) {
			return false
		}
	}
	return true
}

func (p *Project) allNumericRejection() bool {
	for _, key := range p.Keys() {
		if checker.IsNumeric(key) && checker.NotEquals(p.Get(key), enum.ProjectValueRejection) {
			return false
		}
	}
	return true
}
