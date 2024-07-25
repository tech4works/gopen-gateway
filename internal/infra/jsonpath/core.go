package jsonpath

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type core struct {
}

func New() domain.JSONPath {
	return core{}
}

func (c core) Parse(raw string) domain.JSONValue {
	return newValue(gjson.Parse(raw))
}

func (c core) ForEach(raw string, iterator func(key string, value domain.JSONValue) bool) {
	c.Parse(raw).ForEach(iterator)
}

func (c core) Add(raw, path, value string) (string, error) {
	if helper.IsEmpty(path) {
		return raw, nil
	}

	jsonValue := gjson.Get(raw, path)
	if jsonValue.Exists() && helper.IsNotEqualTo(jsonValue.Type, gjson.Null) {
		return sjson.SetRaw(raw, path, aggregateValue(jsonValue, value))
	}

	return sjson.SetRaw(raw, path, parseStringValueToRaw(value))
}

func (c core) AppendOnArray(raw, value string) (string, error) {
	return sjson.SetRaw(raw, "-1", value)
}

func (c core) Set(raw, path, value string) (string, error) {
	if helper.IsEmpty(path) {
		return raw, nil
	}
	return sjson.SetRaw(raw, path, parseStringValueToRaw(value))
}

func (c core) Replace(raw, path, value string) (string, error) {
	if helper.IsEmpty(path) {
		return raw, nil
	}

	jsonValue := gjson.Get(raw, path)
	if !jsonValue.Exists() {
		return raw, nil
	}

	return sjson.SetRaw(raw, path, parseStringValueToRaw(value))
}

func (c core) Delete(raw, path string) (string, error) {
	if helper.IsEmpty(path) {
		return raw, nil
	}
	return sjson.Delete(raw, path)
}

func (c core) Get(raw, path string) domain.JSONValue {
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
		if helper.Equals(v.Type, gjson.Null) || helper.IsEmpty(v.String()) {
			continue
		}
		if helper.IsNotEqualTo(i, 0) {
			newArrayJson += ","
		}
		newArrayJson += parseValueToRaw(v)
	}
	newArrayJson += "]"

	return newArrayJson
}

func parseStringValueToRaw(value string) string {
	parse := gjson.Parse(value)
	if helper.Equals(parse.Type, gjson.Null) {
		return "null"
	}
	return parse.Raw
}

func parseValueToRaw(value gjson.Result) string {
	if helper.Equals(value.Type, gjson.Null) {
		return "null"
	}
	return value.Raw
}
