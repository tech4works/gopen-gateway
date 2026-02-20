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
	"strconv"

	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
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
	var newRaw string
	var err error

	jsonValue := gjson.Get(raw, path)
	if jsonValue.Exists() && checker.NotEquals(jsonValue.Type, gjson.Null) {
		newRaw, err = sjson.SetRaw(raw, path, aggregateValue(jsonValue, value))
	} else {
		newRaw, err = sjson.SetRaw(raw, path, parseStringValueToRaw(value))
	}

	return treatModifierResult("add", newRaw, path, value, newRaw, err)
}

func (p provider) AppendOnArray(raw, value string) (string, error) {
	path := "-1"

	newRaw, err := sjson.SetRaw(raw, path, value)

	return treatModifierResult("append-on-array", newRaw, path, value, newRaw, err)
}

func (p provider) Set(raw, path, value string) (string, error) {
	newRaw, err := sjson.SetRaw(raw, path, parseStringValueToRaw(value))

	return treatModifierResult("set", raw, path, value, newRaw, err)
}

func (p provider) Replace(raw, path, value string) (string, error) {
	jsonValue := gjson.Get(raw, path)
	if !jsonValue.Exists() {
		return raw, nil
	}

	newRaw, err := sjson.SetRaw(raw, path, parseStringValueToRaw(value))

	return treatModifierResult("replace", raw, path, value, newRaw, err)
}

func (p provider) Delete(raw, path string) (string, error) {
	if checker.IsEmpty(path) {
		return raw, nil
	}

	newRaw, err := sjson.Delete(raw, path)

	return treatModifierResult("delete", raw, path, "", newRaw, err)
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
		if checker.IsEmpty(value) || checker.Equals(value, "null") {
			return "null"
		} else {
			parse = gjson.Parse(strconv.Quote(value))
		}
	}
	return parse.Raw
}

func parseValueToRaw(value gjson.Result) string {
	if checker.Equals(value.Type, gjson.Null) {
		return "null"
	}
	return value.Raw
}

func treatModifierResult(op, raw, path, value, newRaw string, err error) (string, error) {
	if checker.NonNil(err) {
		return raw, err
	} else if checker.Equals(raw, newRaw) {
		return raw, mapper.NewErrJSONNotModified(op, path, value)
	}
	return newRaw, nil
}
