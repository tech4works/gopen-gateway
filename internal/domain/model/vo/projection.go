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
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"strconv"
	"strings"
)

type Projection struct {
	keys   []string
	values map[string]enum.ProjectionValue
}

func (p *Projection) IsEmpty() bool {
	return helper.IsEmpty(p.Keys())
}

func (p *Projection) Exists(key string) bool {
	return helper.Contains(p.keys, key)
}

func (p *Projection) ContainsNumericKey() bool {
	for _, key := range p.keys {
		if helper.IsNumeric(key) {
			return true
		}
	}
	return false
}

func (p *Projection) NotContainsNumericKey() bool {
	return !p.ContainsNumericKey()
}

func (p *Projection) Keys() []string {
	return p.keys
}

func (p *Projection) Get(key string) enum.ProjectionValue {
	return p.values[key]
}

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

func (p *Projection) IsAddition(key string) bool {
	return p.Exists(key) && helper.Equals(p.Get(key), enum.ProjectionValueAddition)
}

func (p *Projection) IsRejection(key string) bool {
	return p.Exists(key) && helper.Equals(p.Get(key), enum.ProjectionValueRejection)
}

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

func (p *Projection) allAddition() bool {
	for _, key := range p.Keys() {
		if helper.IsNotEqualTo(p.Get(key), enum.ProjectionValueAddition) {
			return false
		}
	}
	return true
}

func (p *Projection) allNumericAddition() bool {
	for _, key := range p.Keys() {
		if helper.IsNumeric(key) && helper.IsNotEqualTo(p.Get(key), enum.ProjectionValueAddition) {
			return false
		}
	}
	return true
}

func (p *Projection) allRejection() bool {
	for _, key := range p.Keys() {
		if helper.IsNotEqualTo(p.Get(key), enum.ProjectionValueRejection) {
			return false
		}
	}
	return true
}

func (p *Projection) allNumericRejection() bool {
	for _, key := range p.Keys() {
		if helper.IsNumeric(key) && helper.IsNotEqualTo(p.Get(key), enum.ProjectionValueRejection) {
			return false
		}
	}
	return true
}
