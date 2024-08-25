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

package jsonpath

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tidwall/gjson"
)

type value struct {
	result gjson.Result
}

func newValue(result gjson.Result) domain.JSONValue {
	return value{
		result: result,
	}
}

func (v value) Get(path string) domain.JSONValue {
	return newValue(v.result.Get(path))
}

func (v value) ForEach(iterator func(key string, value domain.JSONValue) bool) {
	v.result.ForEach(func(key, value gjson.Result) bool {
		return iterator(key.String(), newValue(value))
	})
}

func (v value) Exists() bool {
	return v.result.Exists()
}

func (v value) NotExists() bool {
	return !v.Exists()
}

func (v value) IsObject() bool {
	return v.result.IsObject()
}

func (v value) IsArray() bool {
	return v.result.IsArray()
}

func (v value) Raw() string {
	if checker.Equals(v.result.Type, gjson.Null) {
		return "null"
	}
	return v.result.Raw
}

func (v value) String() string {
	return v.result.String()
}

func (v value) Interface() any {
	return v.result.Value()
}
