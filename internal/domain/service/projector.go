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
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type projectorService struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Projector interface {
	ProjectHeader(projector *vo.Projector, header vo.Header, request *vo.HTTPRequest, history *vo.History) (vo.Header, error)
	ProjectQuery(projector *vo.Projector, query vo.Query, request *vo.HTTPRequest, history *vo.History) (vo.Query, error)
	ProjectBody(projector *vo.Projector, body *vo.Body, request *vo.HTTPRequest, history *vo.History) (*vo.Body, []error)
}

func NewProjector(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Projector {
	return projectorService{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (s projectorService) ProjectHeader(projector *vo.Projector, header vo.Header, request *vo.HTTPRequest,
	history *vo.History) (vo.Header, error) {
	if checker.IsNil(projector) || projector.Project().IsEmpty() {
		return header, nil
	}

	shouldRun, err := s.evalProjectorGuards("header", projector, request, history)
	if checker.NonNil(err) {
		return header, err
	} else if !shouldRun {
		return header, nil
	}

	if projector.Project().IsAllRejection() {
		return s.projectRejectionHeader(projector.Project(), header)
	} else {
		return s.projectAdditionHeader(projector.Project(), header)
	}
}

func (s projectorService) ProjectQuery(projector *vo.Projector, query vo.Query, request *vo.HTTPRequest,
	history *vo.History) (vo.Query, error) {
	if checker.IsNil(projector) || projector.Project().IsEmpty() {
		return query, nil
	}

	shouldRun, err := s.evalProjectorGuards("query", projector, request, history)
	if checker.NonNil(err) {
		return query, err
	} else if !shouldRun {
		return query, nil
	}

	if projector.Project().IsAllRejection() {
		return s.projectRejectionQuery(projector.Project(), query)
	} else {
		return s.projectAdditionQuery(projector.Project(), query)
	}
}

func (s projectorService) ProjectBody(projector *vo.Projector, body *vo.Body, request *vo.HTTPRequest,
	history *vo.History) (*vo.Body, []error) {
	if checker.IsNil(projector) || projector.Project().IsEmpty() || checker.IsNil(body) || body.ContentType().IsNotJSON() {
		return body, nil
	}

	shouldRun, err := s.evalProjectorGuards("body", projector, request, history)
	if checker.NonNil(err) {
		return body, converter.ToSlice(err)
	} else if !shouldRun {
		return body, nil
	}

	bodyStr, err := body.String()
	if checker.NonNil(err) {
		return body, errors.InheritAsSlice(err, "body projector: failed to stringify body")
	}

	var projectedBody string
	var errs []error

	parsedJson := s.jsonPath.Parse(bodyStr)
	if parsedJson.IsArray() {
		projectedBody, errs = s.projectBodyJsonArray(projector.Project(), parsedJson)
	} else {
		projectedBody, errs = s.projectBodyJsonObject(projector.Project(), parsedJson)
	}

	buffer, err := converter.ToBufferWithErr(projectedBody)
	if checker.NonNil(err) {
		return body, append(errs, errors.Inherit(err, "body projector: failed to build buffer from projected body"))
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), errs
}

func (s projectorService) evalProjectorGuards(kind string, projector *vo.Projector, request *vo.HTTPRequest,
	history *vo.History) (bool, error) {
	shouldRun, _, errs := s.dynamicValueService.EvalGuards(projector.OnlyIf(), projector.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInheritf(errs, ", ", "failed to evaluate guard for %s projector", kind)
	}
	return shouldRun, nil
}

func (s projectorService) projectRejectionHeader(project *vo.Project, header vo.Header) (vo.Header, error) {
	values := header.Copy()
	for _, key := range header.Keys() {
		if mapper.IsNotHeaderMandatoryKey(key) && project.Exists(key) {
			delete(values, key)
		}
	}
	return vo.NewHeader(values), nil
}

func (s projectorService) projectAdditionHeader(project *vo.Project, header vo.Header) (vo.Header, error) {
	values := map[string][]string{}
	for _, key := range header.Keys() {
		if mapper.IsHeaderMandatoryKey(key) || project.IsAddition(key) {
			values[key] = header.GetAll(key)
		}
	}
	return vo.NewHeader(values), nil
}

func (s projectorService) projectRejectionQuery(project *vo.Project, query vo.Query) (vo.Query, error) {
	values := query.Copy()
	for _, key := range query.Keys() {
		if project.Exists(key) {
			delete(values, key)
		}
	}
	return vo.NewQuery(values), nil
}

func (s projectorService) projectAdditionQuery(project *vo.Project, query vo.Query) (vo.Query, error) {
	values := map[string][]string{}
	for _, key := range query.Keys() {
		if project.IsAddition(key) {
			values[key] = query.GetAll(key)
		}
	}
	return vo.NewQuery(values), nil
}

func (s projectorService) projectBodyJsonObject(project *vo.Project, jsonObject domain.JSONValue) (string, []error) {
	if project.IsAllRejection() {
		return s.projectRejectionBodyJsonObject(project, jsonObject)
	}
	return s.projectionAdditionBodyJsonObject(project, jsonObject)
}

func (s projectorService) projectionAdditionBodyJsonObject(project *vo.Project, jsonObject domain.JSONValue) (string,
	[]error) {
	var projectedJson = "{}"
	var errs []error

	for _, key := range project.Keys() {
		if project.IsRejection(key) {
			continue
		}

		jsonValue := jsonObject.Get(key)
		if !jsonValue.Exists() {
			continue
		}

		newProjectedJson, err := s.jsonPath.Set(projectedJson, key, jsonValue.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "body projector: op=set mode=addition path=%s", key))
			continue
		}

		projectedJson = newProjectedJson
	}

	return projectedJson, errs
}

func (s projectorService) projectRejectionBodyJsonObject(project *vo.Project, jsonObject domain.JSONValue) (string,
	[]error) {
	var projectionJson = jsonObject.Raw()
	var errs []error

	for _, key := range project.Keys() {
		newProjectionJson, err := s.jsonPath.Delete(projectionJson, key)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inherit(err, "body projector: op=delete mode=rejection path=%s", key))
			continue
		}

		projectionJson = newProjectionJson
	}

	return projectionJson, errs
}

func (s projectorService) projectBodyJsonArray(project *vo.Project, jsonArray domain.JSONValue) (string, []error) {
	projectedArray, errs := s.projectBodyJsonArrayNormalKeys(project, jsonArray)
	if checker.IsNotEmpty(errs) {
		return projectedArray, errs
	}

	projectedArray, errs = s.projectBodyJsonArrayNumericKeys(project, projectedArray)
	if checker.IsNotEmpty(errs) {
		return projectedArray, errs
	}

	return projectedArray, errs
}

func (s projectorService) projectBodyJsonArrayNormalKeys(project *vo.Project, jsonArray domain.JSONValue) (string,
	[]error) {
	var projectedArray = "[]"
	var errs []error

	jsonArray.ForEach(func(key string, value domain.JSONValue) bool {
		var newProjectedArray string
		var err error
		if value.IsObject() {
			childObject, childErrs := s.projectBodyJsonObject(project, value)
			if checker.IsNotEmpty(childErrs) {
				for _, ce := range childErrs {
					errs = append(errs, errors.Inheritf(ce, "body projector array normal keys with child object: array idx=%s", key))
				}
				return true
			}
			newProjectedArray, err = s.jsonPath.AppendOnArray(projectedArray, childObject)
		} else if value.IsArray() {
			childArray, childErrs := s.projectBodyJsonArray(project, value)
			if checker.IsNotEmpty(childErrs) {
				for _, ce := range childErrs {
					errs = append(errs, errors.Inheritf(ce, "body projector array normal keys with child array: array idx=%s", key))
				}
				return true
			}
			newProjectedArray, err = s.jsonPath.AppendOnArray(projectedArray, childArray)
		} else {
			newProjectedArray, err = s.jsonPath.AppendOnArray(projectedArray, value.Raw())
		}

		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "body projector: op=append array idx=%s", key))
			return true
		}

		projectedArray = newProjectedArray
		return true
	})

	return projectedArray, errs
}

func (s projectorService) projectBodyJsonArrayNumericKeys(project *vo.Project, projectedArray string) (string, []error) {
	if project.NotContainsNumericKey() {
		return projectedArray, nil
	} else if project.IsAllRejection() {
		return s.projectRejectionBodyJsonArray(project, projectedArray)
	} else {
		return s.projectAdditionBodyJsonArray(project, projectedArray)
	}
}

func (s projectorService) projectRejectionBodyJsonArray(project *vo.Project, projectedJson string) (string, []error) {
	var projectedArray = "[]"
	var errs []error

	s.jsonPath.ForEach(projectedJson, func(key string, value domain.JSONValue) bool {
		if checker.Contains(project.Keys(), key) {
			return true
		}

		newProjectedArray, err := s.jsonPath.AppendOnArray(projectedArray, value.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "body projector: op=append mode=rejection numericKey=%s", key))
			return true
		}

		projectedArray = newProjectedArray
		return true
	})

	return projectedArray, errs
}

func (s projectorService) projectAdditionBodyJsonArray(project *vo.Project, projectedJson string) (string, []error) {
	var projectedArray = "[]"
	var errs []error

	parsedProjectedJson := s.jsonPath.Parse(projectedJson)
	for _, key := range project.Keys() {
		if checker.IsNotNumeric(key) || project.IsRejection(key) {
			continue
		}

		jsonValue := parsedProjectedJson.Get(key)
		if !jsonValue.Exists() {
			continue
		}

		newProjectedArray, err := s.jsonPath.AppendOnArray(projectedArray, jsonValue.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "body projector: op=append mode=addition numericKey=%s", key))
			continue
		}

		projectedArray = newProjectedArray
	}

	return projectedArray, errs
}
