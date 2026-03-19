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
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	vo "github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type projector struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Projector interface {
	ProjectMetadata(
		config *vo.ProjectorConfig,
		metadata vo.Metadata,
		ignoreKeys []string,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (vo.Metadata, error)
	ProjectQuery(config *vo.ProjectorConfig, query vo.Query, request *vo.EndpointRequest, history *aggregate.History,
	) (vo.Query, error)
	ProjectPayload(config *vo.ProjectorConfig, payload *vo.Payload, request *vo.EndpointRequest, history *aggregate.History,
	) (*vo.Payload, []error)
}

func NewProjector(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Projector {
	return projector{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (s projector) ProjectMetadata(
	config *vo.ProjectorConfig,
	metadata vo.Metadata,
	ignoreKeys []string,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Metadata, error) {
	if checker.IsNil(config) || config.Project().IsEmpty() {
		return metadata, nil
	}

	shouldRun, err := s.evalProjectorGuards(config, "metadata", request, history)
	if checker.NonNil(err) {
		return metadata, err
	} else if !shouldRun {
		return metadata, nil
	}

	if config.Project().IsAllRejection() {
		return s.projectRejectionMetadata(config.Project(), metadata, ignoreKeys)
	} else {
		return s.projectAdditionMetadata(config.Project(), metadata, ignoreKeys)
	}
}

func (s projector) ProjectQuery(
	config *vo.ProjectorConfig,
	query vo.Query,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Query, error) {
	if checker.IsNil(config) || config.Project().IsEmpty() {
		return query, nil
	}

	shouldRun, err := s.evalProjectorGuards(config, "query", request, history)
	if checker.NonNil(err) {
		return query, err
	} else if !shouldRun {
		return query, nil
	}

	if config.Project().IsAllRejection() {
		return s.projectRejectionQuery(config.Project(), query)
	} else {
		return s.projectAdditionQuery(config.Project(), query)
	}
}

func (s projector) ProjectPayload(
	config *vo.ProjectorConfig,
	payload *vo.Payload,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.Payload, []error) {
	if checker.IsNil(config) || config.Project().IsEmpty() || checker.IsNil(payload) || payload.ContentType().IsNotJSON() {
		return payload, nil
	}

	shouldRun, err := s.evalProjectorGuards(config, "payload", request, history)
	if checker.NonNil(err) {
		return payload, converter.ToSlice(err)
	} else if !shouldRun {
		return payload, nil
	}

	str, err := payload.String()
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlice(err, "projector failed: failed to stringify payload")
	}

	var projectedPayload string
	var errs []error

	parsedJSON := s.jsonPath.Parse(str)
	if parsedJSON.IsArray() {
		projectedPayload, errs = s.projectPayloadJSONArray(config.Project(), parsedJSON)
	} else {
		projectedPayload, errs = s.projectPayloadJSONObject(config.Project(), parsedJSON)
	}

	buffer, err := converter.ToBufferWithErr(projectedPayload)
	if checker.NonNil(err) {
		return payload, append(errs, errors.Inherit(err, "projector failed: failed to build buffer from projected payload"))
	}

	return vo.NewPayloadWithContentType(payload.ContentType(), buffer), errs
}

func (s projector) evalProjectorGuards(
	config *vo.ProjectorConfig,
	kind string,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (bool, error) {
	shouldRun, _, errs := s.dynamicValueService.EvalGuards(config.OnlyIf(), config.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInheritf(errs, ", ", "failed to evaluate guard for %s projector", kind)
	}
	return shouldRun, nil
}

func (s projector) projectRejectionMetadata(config *vo.ProjectConfig, metadata vo.Metadata, ignoreKeys []string) (
	vo.Metadata, error) {
	values := metadata.Copy()
	for _, key := range metadata.Keys() {
		if checker.NotContains(ignoreKeys, key) && config.Exists(key) {
			delete(values, key)
		}
	}
	return vo.NewMetadata(values), nil
}

func (s projector) projectAdditionMetadata(config *vo.ProjectConfig, metadata vo.Metadata, ignoreKeys []string) (
	vo.Metadata, error) {
	values := map[string][]string{}
	for _, key := range metadata.Keys() {
		if checker.Contains(ignoreKeys, key) || config.IsAddition(key) {
			values[key] = metadata.GetAll(key)
		}
	}
	return vo.NewMetadata(values), nil
}

func (s projector) projectRejectionQuery(config *vo.ProjectConfig, query vo.Query) (vo.Query, error) {
	values := query.Copy()
	for _, key := range query.Keys() {
		if config.Exists(key) {
			delete(values, key)
		}
	}
	return vo.NewQuery(values), nil
}

func (s projector) projectAdditionQuery(config *vo.ProjectConfig, query vo.Query) (vo.Query, error) {
	values := map[string][]string{}
	for _, key := range query.Keys() {
		if config.IsAddition(key) {
			values[key] = query.GetAll(key)
		}
	}
	return vo.NewQuery(values), nil
}

func (s projector) projectPayloadJSONObject(config *vo.ProjectConfig, jsonObject domain.JSONValue) (string, []error) {
	if config.IsAllRejection() {
		return s.projectRejectionPayloadJSONObject(config, jsonObject)
	}
	return s.projectionAdditionPayloadJSONObject(config, jsonObject)
}

func (s projector) projectionAdditionPayloadJSONObject(config *vo.ProjectConfig, jsonObject domain.JSONValue) (string,
	[]error) {
	var projectedJSON = "{}"
	var errs []error

	for _, key := range config.Keys() {
		if config.IsRejection(key) {
			continue
		}

		jsonValue := jsonObject.Get(key)
		if !jsonValue.Exists() {
			continue
		}

		newProjectedJSON, err := s.jsonPath.Set(projectedJSON, key, jsonValue.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "projector failed: op=set mode=addition path=%s", key))
			continue
		}

		projectedJSON = newProjectedJSON
	}

	return projectedJSON, errs
}

func (s projector) projectRejectionPayloadJSONObject(config *vo.ProjectConfig, jsonObject domain.JSONValue) (string,
	[]error) {
	var projectionJSON = jsonObject.Raw()
	var errs []error

	for _, key := range config.Keys() {
		newProjectionJSON, err := s.jsonPath.Delete(projectionJSON, key)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inherit(err, "projector failed: op=delete mode=rejection path=%s", key))
			continue
		}

		projectionJSON = newProjectionJSON
	}

	return projectionJSON, errs
}

func (s projector) projectPayloadJSONArray(config *vo.ProjectConfig, jsonArray domain.JSONValue) (string, []error) {
	projectedArray, errs := s.projectPayloadJSONArrayNormalKeys(config, jsonArray)
	if checker.IsNotEmpty(errs) {
		return projectedArray, errs
	}

	projectedArray, errs = s.projectPayloadJSONArrayNumericKeys(config, projectedArray)
	if checker.IsNotEmpty(errs) {
		return projectedArray, errs
	}

	return projectedArray, errs
}

func (s projector) projectPayloadJSONArrayNormalKeys(config *vo.ProjectConfig, jsonArray domain.JSONValue) (string,
	[]error) {
	var projectedArray = "[]"
	var errs []error

	jsonArray.ForEach(func(key string, value domain.JSONValue) bool {
		var newProjectedArray string
		var err error
		if value.IsObject() {
			childObject, childErrs := s.projectPayloadJSONObject(config, value)
			if checker.IsNotEmpty(childErrs) {
				for _, ce := range childErrs {
					errs = append(errs, errors.Inheritf(ce, "projector failed: op=project-payload-json-object idx=%s",
						key))
				}
				return true
			}
			newProjectedArray, err = s.jsonPath.AppendOnArray(projectedArray, childObject)
		} else if value.IsArray() {
			childArray, childErrs := s.projectPayloadJSONArray(config, value)
			if checker.IsNotEmpty(childErrs) {
				for _, ce := range childErrs {
					errs = append(errs, errors.Inheritf(ce, "projector failed: op=project-payload-json-array idx=%s",
						key))
				}
				return true
			}
			newProjectedArray, err = s.jsonPath.AppendOnArray(projectedArray, childArray)
		} else {
			newProjectedArray, err = s.jsonPath.AppendOnArray(projectedArray, value.Raw())
		}

		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "projector failed: op=append idx=%s", key))
			return true
		}

		projectedArray = newProjectedArray
		return true
	})

	return projectedArray, errs
}

func (s projector) projectPayloadJSONArrayNumericKeys(config *vo.ProjectConfig, projectedArray string) (string, []error) {
	if config.NotContainsNumericKey() {
		return projectedArray, nil
	} else if config.IsAllRejection() {
		return s.projectRejectionPayloadJSONArray(config, projectedArray)
	} else {
		return s.projectAdditionPayloadJSONArray(config, projectedArray)
	}
}

func (s projector) projectRejectionPayloadJSONArray(config *vo.ProjectConfig, projectedJSON string) (string, []error) {
	var projectedArray = "[]"
	var errs []error

	s.jsonPath.ForEach(projectedJSON, func(key string, value domain.JSONValue) bool {
		if checker.Contains(config.Keys(), key) {
			return true
		}

		newProjectedArray, err := s.jsonPath.AppendOnArray(projectedArray, value.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "projector failed: op=append mode=rejection numericKey=%s", key))
			return true
		}

		projectedArray = newProjectedArray
		return true
	})

	return projectedArray, errs
}

func (s projector) projectAdditionPayloadJSONArray(config *vo.ProjectConfig, projectedJSON string) (string, []error) {
	var projectedArray = "[]"
	var errs []error

	parsedProjectedJSON := s.jsonPath.Parse(projectedJSON)
	for _, key := range config.Keys() {
		if checker.IsNotNumeric(key) || config.IsRejection(key) {
			continue
		}

		jsonValue := parsedProjectedJSON.Get(key)
		if !jsonValue.Exists() {
			continue
		}

		newProjectedArray, err := s.jsonPath.AppendOnArray(projectedArray, jsonValue.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "projector failed: op=append mode=addition key=%s", key))
			continue
		}

		projectedArray = newProjectedArray
	}

	return projectedArray, errs
}
