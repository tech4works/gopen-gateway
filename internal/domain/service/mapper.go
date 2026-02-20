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

package service

import (
	"strings"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	domainMapper "github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type mapperService struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Mapper interface {
	MapHeader(mapper *vo.Mapper, header vo.Header, request *vo.HTTPRequest, history *vo.History) (vo.Header, error)
	MapQuery(mapper *vo.Mapper, query vo.Query, request *vo.HTTPRequest, history *vo.History) (vo.Query, error)
	MapBody(mapper *vo.Mapper, body *vo.Body, request *vo.HTTPRequest, history *vo.History) (*vo.Body, []error)
}

func NewMapper(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Mapper {
	return mapperService{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (m mapperService) MapHeader(mapper *vo.Mapper, header vo.Header, request *vo.HTTPRequest, history *vo.History,
) (vo.Header, error) {
	if checker.IsNil(mapper) || mapper.Map().IsEmpty() {
		return header, nil
	}

	shouldRun, err := m.evalMapperGuards("header", mapper, request, history)
	if checker.NonNil(err) {
		return header, err
	} else if !shouldRun {
		return header, nil
	}

	mappedHeader := map[string][]string{}
	for _, key := range header.Keys() {
		if domainMapper.IsNotHeaderMandatoryKey(key) && mapper.Map().Exists(key) {
			mappedHeader[mapper.Map().Get(key)] = header.GetAll(key)
		} else {
			mappedHeader[key] = header.GetAll(key)
		}
	}

	return vo.NewHeader(mappedHeader), nil
}

func (m mapperService) MapQuery(mapper *vo.Mapper, query vo.Query, request *vo.HTTPRequest, history *vo.History) (
	vo.Query, error) {
	if checker.IsNil(mapper) || mapper.Map().IsEmpty() {
		return query, nil
	}

	shouldRun, err := m.evalMapperGuards("query", mapper, request, history)
	if checker.NonNil(err) {
		return query, err
	} else if !shouldRun {
		return query, nil
	}

	mappedQuery := map[string][]string{}
	for _, key := range query.Keys() {
		if mapper.Map().Exists(key) {
			mappedQuery[mapper.Map().Get(key)] = query.GetAll(key)
		} else {
			mappedQuery[key] = query.GetAll(key)
		}
	}
	return vo.NewQuery(mappedQuery), nil
}

func (m mapperService) MapBody(mapper *vo.Mapper, body *vo.Body, request *vo.HTTPRequest, history *vo.History,
) (*vo.Body, []error) {
	if checker.IsNil(mapper) || mapper.Map().IsEmpty() || checker.IsNil(body) {
		return body, nil
	}

	shouldRun, err := m.evalMapperGuards("body", mapper, request, history)
	if checker.NonNil(err) {
		return body, []error{err}
	} else if !shouldRun {
		return body, nil
	}

	if body.ContentType().IsText() {
		return m.mapBodyText(mapper.Map(), body)
	} else if body.ContentType().IsJSON() {
		return m.mapBodyJson(mapper.Map(), body)
	} else {
		return body, nil
	}
}

func (m mapperService) evalMapperGuards(kind string, mapper *vo.Mapper, request *vo.HTTPRequest, history *vo.History) (
	bool, error) {
	shouldRun, _, errs := m.dynamicValueService.EvalGuards(mapper.OnlyIf(), mapper.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInheritf(errs, ", ", "failed to evaluate guard for %s mapper", kind)
	}
	return shouldRun, nil
}

func (m mapperService) mapBodyText(mMap *vo.Map, body *vo.Body) (*vo.Body, []error) {
	mappedBody, err := body.String()
	if checker.NonNil(err) {
		return body, []error{errors.Inherit(err, "body mapper text: failed to stringify body")}
	}

	for _, key := range mMap.Keys() {
		newKey := mMap.Get(key)
		if checker.NotEquals(key, newKey) {
			mappedBody = strings.ReplaceAll(mappedBody, key, newKey)
		}
	}

	buffer, err := converter.ToBufferWithErr(mappedBody)
	if checker.NonNil(err) {
		return body, []error{errors.Inherit(err, "body mapper text: failed to build buffer from mapped body")}
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}

func (m mapperService) mapBodyJson(mMap *vo.Map, body *vo.Body) (*vo.Body, []error) {
	bodyStr, err := body.String()
	if checker.NonNil(err) {
		return body, []error{errors.Inherit(err, "body mapper json: failed to stringify body")}
	}

	var mappedBodyStr string
	var errs []error

	parsedJson := m.jsonPath.Parse(bodyStr)
	if parsedJson.IsArray() {
		mappedBodyStr, errs = m.mapBodyJsonArray(mMap, parsedJson)
	} else {
		mappedBodyStr, errs = m.mapBodyJsonObject(mMap, parsedJson)
	}

	buffer, err := converter.ToBufferWithErr(mappedBodyStr)
	if checker.NonNil(err) {
		return body, append(errs, errors.Inherit(err, "body mapper json: failed to build buffer from mapped body"))
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}

func (m mapperService) mapBodyJsonArray(mMap *vo.Map, jsonArray domain.JSONValue) (string, []error) {
	var mappedArray = "[]"
	var errs []error

	jsonArray.ForEach(func(key string, value domain.JSONValue) bool {
		var newMappedArray string
		var err error
		if value.IsObject() {
			childObject, childErrs := m.mapBodyJsonObject(mMap, value)
			if checker.IsNotEmpty(childErrs) {
				for _, ce := range childErrs {
					errs = append(errs, errors.Inheritf(ce, "body mapper json array: child object idx=%s", key))
				}
				return true
			}
			newMappedArray, err = m.jsonPath.AppendOnArray(mappedArray, childObject)
		} else if value.IsArray() {
			childArray, childErrs := m.mapBodyJsonArray(mMap, value)
			if checker.IsNotEmpty(childErrs) {
				for _, ce := range childErrs {
					errs = append(errs, errors.Inheritf(ce, "body mapper json array: child array idx=%s", key))
				}
				return true
			}
			newMappedArray, err = m.jsonPath.AppendOnArray(mappedArray, childArray)
		} else {
			newMappedArray, err = m.jsonPath.AppendOnArray(mappedArray, value.Raw())
		}

		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "body mapper json array: op=append idx=%s", key))
			return true
		}

		mappedArray = newMappedArray
		return true
	})

	return mappedArray, errs
}

func (m mapperService) mapBodyJsonObject(mMap *vo.Map, jsonObject domain.JSONValue) (string, []error) {
	var mappedJson = jsonObject.Raw()
	var errs []error

	for _, key := range mMap.Keys() {
		newKey := mMap.Get(key)
		if checker.Equals(key, newKey) {
			continue
		}

		jsonValue := jsonObject.Get(key)
		if jsonValue.NotExists() {
			continue
		}

		newMappedJson, err := m.jsonPath.Set(mappedJson, newKey, jsonValue.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "body mapper json object: op=set from=%s to=%s", key, newKey))
			continue
		}

		newMappedJson, err = m.jsonPath.Delete(newMappedJson, key)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inherit(err, "body mapper json object: op=delete from=%s to=%s", key, newKey))
			continue
		}

		mappedJson = newMappedJson
	}

	return mappedJson, errs
}
