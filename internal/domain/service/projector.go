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
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type projectorService struct {
	jsonPath domain.JSONPath
}

type Projector interface {
	ProjectHeader(header vo.Header, projection *vo.Projection) vo.Header
	ProjectQuery(query vo.Query, projection *vo.Projection) vo.Query
	ProjectBody(body *vo.Body, projection *vo.Projection) (*vo.Body, []error)
}

func NewProjector(jsonPath domain.JSONPath) Projector {
	return projectorService{
		jsonPath: jsonPath,
	}
}

func (s projectorService) ProjectHeader(header vo.Header, projection *vo.Projection) vo.Header {
	if checker.IsNil(projection) || projection.IsEmpty() {
		return header
	} else if checker.Equals(projection.Type(), enum.ProjectionTypeRejection) {
		return s.projectRejectionHeader(header, projection)
	} else {
		return s.projectAdditionHeader(header, projection)
	}
}

func (s projectorService) ProjectQuery(query vo.Query, projection *vo.Projection) vo.Query {
	if checker.IsNil(projection) || projection.IsEmpty() {
		return query
	} else if checker.Equals(projection.Type(), enum.ProjectionValueRejection) {
		return s.projectRejectionQuery(query, projection)
	} else {
		return s.projectAdditionQuery(query, projection)
	}
}

func (s projectorService) ProjectBody(body *vo.Body, projection *vo.Projection) (*vo.Body, []error) {
	if checker.IsNil(projection) || projection.IsEmpty() || checker.IsNil(body) || body.ContentType().IsNotJSON() {
		return body, nil
	}

	bodyStr, err := body.String()
	if checker.NonNil(err) {
		return body, []error{err}
	}

	var projectedBody string
	var errs []error

	parsedJson := s.jsonPath.Parse(bodyStr)
	if parsedJson.IsArray() {
		projectedBody, errs = s.projectBodyJsonArray(parsedJson, projection)
	} else {
		projectedBody, errs = s.projectBodyJsonObject(parsedJson, projection)
	}
	if checker.IsNotEmpty(errs) {
		return body, errs
	}

	buffer, err := converter.ToBufferWithErr(projectedBody)
	if checker.NonNil(err) {
		return body, []error{err}
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}

func (s projectorService) projectRejectionHeader(header vo.Header, projection *vo.Projection) vo.Header {
	values := header.Copy()
	for _, key := range header.Keys() {
		if mapper.IsNotHeaderMandatoryKey(key) && projection.Exists(key) {
			delete(values, key)
		}
	}
	return vo.NewHeader(values)
}

func (s projectorService) projectAdditionHeader(header vo.Header, projection *vo.Projection) vo.Header {
	values := map[string][]string{}
	for _, key := range header.Keys() {
		if mapper.IsHeaderMandatoryKey(key) || projection.IsAddition(key) {
			values[key] = header.GetAll(key)
		}
	}
	return vo.NewHeader(values)
}

func (s projectorService) projectRejectionQuery(query vo.Query, projection *vo.Projection) vo.Query {
	values := query.Copy()
	for _, key := range query.Keys() {
		if projection.Exists(key) {
			delete(values, key)
		}
	}
	return vo.NewQuery(values)
}

func (s projectorService) projectAdditionQuery(query vo.Query, projection *vo.Projection) vo.Query {
	values := map[string][]string{}
	for _, key := range query.Keys() {
		if projection.IsAddition(key) {
			values[key] = query.GetAll(key)
		}
	}
	return vo.NewQuery(values)
}

func (s projectorService) projectBodyJsonObject(jsonObject domain.JSONValue, projection *vo.Projection) (string, []error) {
	if checker.Equals(projection.Type(), enum.ProjectionTypeRejection) {
		return s.projectRejectionBodyJsonObject(jsonObject, projection)
	}
	return s.projectionAdditionBodyJsonObject(jsonObject, projection)
}

func (s projectorService) projectionAdditionBodyJsonObject(jsonObject domain.JSONValue, projection *vo.Projection) (string,
	[]error) {
	var projectedJson = "{}"
	var errs []error

	for _, key := range projection.Keys() {
		if projection.IsRejection(key) {
			continue
		}

		jsonValue := jsonObject.Get(key)
		if !jsonValue.Exists() {
			continue
		}

		newProjectedJson, err := s.jsonPath.Set(projectedJson, key, jsonValue.Raw())
		if checker.NonNil(err) {
			errs = append(errs, err)
			continue
		}

		projectedJson = newProjectedJson
	}

	return projectedJson, errs
}

func (s projectorService) projectRejectionBodyJsonObject(jsonObject domain.JSONValue, projection *vo.Projection) (string,
	[]error) {
	var projectionJson = jsonObject.Raw()
	var errs []error

	for _, key := range projection.Keys() {
		newProjectionJson, err := s.jsonPath.Delete(projectionJson, key)
		if checker.NonNil(err) {
			errs = append(errs, err)
			continue
		}
		projectionJson = newProjectionJson
	}

	return projectionJson, errs
}

func (s projectorService) projectBodyJsonArray(jsonArray domain.JSONValue, projection *vo.Projection) (string, []error) {
	projectedArray, errs := s.projectBodyJsonArrayNormalKeys(jsonArray, projection)
	if checker.IsNotEmpty(errs) {
		return "", errs
	}

	projectedArray, errs = s.projectBodyJsonArrayNumericKeys(projectedArray, projection)
	if checker.IsNotEmpty(errs) {
		return "", errs
	}

	return projectedArray, errs
}

func (s projectorService) projectBodyJsonArrayNormalKeys(jsonArray domain.JSONValue, projection *vo.Projection) (string,
	[]error) {
	var projectedArray = "[]"
	var errs []error

	jsonArray.ForEach(func(key string, value domain.JSONValue) bool {
		var newProjectedArray string
		var err error
		if value.IsObject() {
			childObject, childErrs := s.projectBodyJsonObject(value, projection)
			if checker.IsNotEmpty(childErrs) {
				errs = append(errs, childErrs...)
				return true
			}
			newProjectedArray, err = s.jsonPath.AppendOnArray(projectedArray, childObject)
		} else if value.IsArray() {
			childArray, childErrs := s.projectBodyJsonArray(value, projection)
			if checker.IsNotEmpty(childErrs) {
				errs = append(errs, childErrs...)
				return true
			}
			newProjectedArray, err = s.jsonPath.AppendOnArray(projectedArray, childArray)
		} else {
			newProjectedArray, err = s.jsonPath.AppendOnArray(projectedArray, value.Raw())
		}

		if checker.NonNil(err) {
			errs = append(errs, err)
			return true
		}

		projectedArray = newProjectedArray
		return true
	})

	return projectedArray, errs
}

func (s projectorService) projectBodyJsonArrayNumericKeys(projectedArray string, projection *vo.Projection) (string, []error) {
	if projection.NotContainsNumericKey() {
		return projectedArray, nil
	} else if checker.Equals(projection.TypeNumeric(), enum.ProjectionTypeRejection) {
		return s.projectRejectionBodyJsonArray(projectedArray, projection)
	} else {
		return s.projectAdditionBodyJsonArray(projectedArray, projection)
	}
}

func (s projectorService) projectRejectionBodyJsonArray(projectedJson string, projection *vo.Projection) (string, []error) {
	var projectedArray = "[]"
	var errs []error

	s.jsonPath.ForEach(projectedJson, func(key string, value domain.JSONValue) bool {
		if checker.Contains(projection.Keys(), key) {
			return true
		}

		newProjectedArray, err := s.jsonPath.AppendOnArray(projectedArray, value.Raw())
		if checker.NonNil(err) {
			errs = append(errs, err)
			return true
		}

		projectedArray = newProjectedArray
		return true
	})

	return projectedArray, errs
}

func (s projectorService) projectAdditionBodyJsonArray(projectedJson string, projection *vo.Projection) (string, []error) {
	var projectedArray = "[]"
	var errs []error

	parsedProjectedJson := s.jsonPath.Parse(projectedJson)
	for _, key := range projection.Keys() {
		if checker.IsNotNumeric(key) || projection.IsRejection(key) {
			continue
		}

		jsonValue := parsedProjectedJson.Get(key)
		if !jsonValue.Exists() {
			continue
		}

		newProjectedArray, err := s.jsonPath.AppendOnArray(projectedArray, jsonValue.Raw())
		if checker.NonNil(err) {
			errs = append(errs, err)
			continue
		}

		projectedArray = newProjectedArray
	}

	return projectedArray, errs
}
