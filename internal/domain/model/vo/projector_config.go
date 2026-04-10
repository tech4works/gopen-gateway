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

type ProjectorConfig struct {
	onlyIf   []string
	ignoreIf []string
	project  ProjectConfig
}

type ProjectConfig struct {
	keys   []string
	values map[string]enum.ProjectValue
}

func NewProjectorConfig(onlyIf, ignoreIf []string, project ProjectConfig) *ProjectorConfig {
	return &ProjectorConfig{
		onlyIf:   onlyIf,
		ignoreIf: ignoreIf,
		project:  project,
	}
}

func (p ProjectorConfig) OnlyIf() []string {
	return p.onlyIf
}

func (p ProjectorConfig) IgnoreIf() []string {
	return p.ignoreIf
}

func (p ProjectorConfig) Project() *ProjectConfig {
	return &p.project
}

func (p *ProjectConfig) IsEmpty() bool {
	return checker.IsEmpty(p.Keys())
}

func (p *ProjectConfig) Exists(key string) bool {
	return checker.IsNotEmpty(p.keys) && checker.Contains(p.keys, key)
}

func (p *ProjectConfig) ContainsNumericKey() bool {
	for _, key := range p.keys {
		if checker.IsNumeric(key) {
			return true
		}
	}
	return false
}

func (p *ProjectConfig) NotContainsNumericKey() bool {
	return !p.ContainsNumericKey()
}

func (p *ProjectConfig) Keys() []string {
	return p.keys
}

func (p *ProjectConfig) Get(key string) enum.ProjectValue {
	return p.values[key]
}

func (p *ProjectConfig) Kind() enum.ProjectKind {
	addition := p.allAddition()
	rejection := p.allRejection()
	if addition && !rejection {
		return enum.ProjectKindAddition
	} else if rejection && !addition {
		return enum.ProjectKindRejection
	}
	return enum.ProjectKindAll
}

func (p *ProjectConfig) NumericKind() enum.ProjectKind {
	addition := p.allNumericAddition()
	rejection := p.allNumericRejection()
	if addition && !rejection {
		return enum.ProjectKindAddition
	} else if rejection && !addition {
		return enum.ProjectKindRejection
	}
	return enum.ProjectKindAll
}

func (p *ProjectConfig) IsAddition(key string) bool {
	return p.Exists(key) && checker.Equals(p.Get(key), enum.ProjectValueAddition)
}

func (p *ProjectConfig) IsRejection(key string) bool {
	return p.Exists(key) && checker.Equals(p.Get(key), enum.ProjectValueRejection)
}

func (p *ProjectConfig) IsAllRejection() bool {
	return checker.Equals(p.Kind(), enum.ProjectKindRejection)
}

func (p *ProjectConfig) UnmarshalJSON(data []byte) error {
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

func (p *ProjectConfig) MarshalJSON() ([]byte, error) {
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

func (p *ProjectConfig) allAddition() bool {
	for _, key := range p.Keys() {
		if checker.NotEquals(p.Get(key), enum.ProjectValueAddition) {
			return false
		}
	}
	return true
}

func (p *ProjectConfig) allNumericAddition() bool {
	for _, key := range p.Keys() {
		if checker.IsNumeric(key) && checker.NotEquals(p.Get(key), enum.ProjectValueAddition) {
			return false
		}
	}
	return true
}

func (p *ProjectConfig) allRejection() bool {
	for _, key := range p.Keys() {
		if checker.NotEquals(p.Get(key), enum.ProjectValueRejection) {
			return false
		}
	}
	return true
}

func (p *ProjectConfig) allNumericRejection() bool {
	for _, key := range p.Keys() {
		if checker.IsNumeric(key) && checker.NotEquals(p.Get(key), enum.ProjectValueRejection) {
			return false
		}
	}
	return true
}
