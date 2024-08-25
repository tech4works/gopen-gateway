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
	"github.com/tidwall/sjson"
)

type provider struct {
}

func New() domain.JSONPath {
	return provider{}
}

func (p provider) Parse(raw string) domain.JSONValue {
	return newValue(gjson.Parse(raw))
}

func (p provider) ForEach(raw string, iterator func(key string, value domain.JSONValue) bool) {
	p.Parse(raw).ForEach(iterator)
}

func (p provider) Add(raw, path, value string) (string, error) {
	if checker.IsEmpty(path) {
		return raw, nil
	}

	jsonValue := gjson.Get(raw, path)
	if jsonValue.Exists() && checker.NotEquals(jsonValue.Type, gjson.Null) {
		return sjson.SetRaw(raw, path, aggregateValue(jsonValue, value))
	}

	return sjson.SetRaw(raw, path, parseStringValueToRaw(value))
}

func (p provider) AppendOnArray(raw, value string) (string, error) {
	return sjson.SetRaw(raw, "-1", value)
}

func (p provider) Set(raw, path, value string) (string, error) {
	if checker.IsEmpty(path) {
		return raw, nil
	}
	return sjson.SetRaw(raw, path, parseStringValueToRaw(value))
}

func (p provider) Replace(raw, path, value string) (string, error) {
	if checker.IsEmpty(path) {
		return raw, nil
	}

	jsonValue := gjson.Get(raw, path)
	if !jsonValue.Exists() {
		return raw, nil
	}

	return sjson.SetRaw(raw, path, parseStringValueToRaw(value))
}

func (p provider) Delete(raw, path string) (string, error) {
	if checker.IsEmpty(path) {
		return raw, nil
	}
	return sjson.Delete(raw, path)
}

func (p provider) Get(raw, path string) domain.JSONValue {
	return newValue(gjson.Get(raw, path))
}

func aggregateValue(jsonValue gjson.Result, newValue string) string {
	var newArray []gjson.Result

	if jsonValue.IsArray() {
		newArray = jsonValue.Array()
	} else {
		newArray = []gjson.Result{jsonValue}
	}

	newParsedValue := gjson.Parse(newValue)
	if newParsedValue.IsArray() {
		newArray = append(newArray, newParsedValue.Array()...)
	} else {
		newArray = append(newArray, newParsedValue)
	}

	newArrayJson := "["
	for i, v := range newArray {
		if checker.Equals(v.Type, gjson.Null) || checker.IsEmpty(v.String()) {
			continue
		}
		if checker.NotEquals(i, 0) {
			newArrayJson += ","
		}
		newArrayJson += parseValueToRaw(v)
	}
	newArrayJson += "]"

	return newArrayJson
}

func parseStringValueToRaw(value string) string {
	parse := gjson.Parse(value)
	if checker.Equals(parse.Type, gjson.Null) {
		return "null"
	}
	return parse.Raw
}

func parseValueToRaw(value gjson.Result) string {
	if checker.Equals(value.Type, gjson.Null) {
		return "null"
	}
	return value.Raw
}
