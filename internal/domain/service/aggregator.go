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

type aggregatorService struct {
	jsonPath domain.JSONPath
}

type Aggregator interface {
	AggregateHeaders(base, header vo.Header) vo.Header
	AggregateBodyToKey(key string, body *vo.Body) (*vo.Body, error)
	AggregateBodiesIntoSlice(history *vo.History) (*vo.Body, []error)
	AggregateBodies(history *vo.History) (*vo.Body, []error)
}

func NewAggregator(jsonPath domain.JSONPath) Aggregator {
	return aggregatorService{
		jsonPath: jsonPath,
	}
}

func (a aggregatorService) AggregateHeaders(base, value vo.Header) vo.Header {
	aggregated := base.Copy()
	for _, key := range value.Keys() {
		if mapper.IsNotHeaderMandatoryKey(key) {
			aggregated[key] = append(aggregated[key], value.GetAll(key)...)
		}
	}
	return vo.NewHeader(aggregated)
}

func (a aggregatorService) AggregateBodyToKey(key string, body *vo.Body) (*vo.Body, error) {
	if checker.IsEmpty(key) || body.ContentType().IsNotJSON() || checker.IsNil(body) {
		return body, nil
	}

	raw, err := body.Raw()
	if checker.NonNil(err) {
		return body, errors.Inheritf(err, "aggregator failed: op=raw key=%s", key)
	}

	jsonValue, err := a.jsonPath.Set("{}", key, raw)
	if checker.NonNil(err) {
		return body, errors.Inheritf(err, "aggregator failed: op=set key=%s", key)
	}

	buffer, err := converter.ToBufferWithErr(jsonValue)
	if checker.NonNil(err) {
		return body, errors.Inheritf(err, "aggregator failed: op=buffer key=%s", key)
	}

	return vo.NewBodyWithContentType(vo.NewContentTypeJson(), buffer), nil
}

func (a aggregatorService) AggregateBodiesIntoSlice(history *vo.History) (*vo.Body, []error) {
	result := "[]"

	var errs []error
	for i := 0; checker.IsLessThan(i, history.Size()); i++ {
		backendResponse := history.GetBackendResponse(i)
		if checker.IsNil(backendResponse) {
			continue
		}

		var err error
		var raw string

		backendID := history.GetBackendID(i)
		defaultJSON := a.buildBodyDefaultForSlice(backendResponse)

		if !backendResponse.HasBody() {
			result, err = a.jsonPath.AppendOnArray(result, defaultJSON)
			if checker.NonNil(err) {
				errs = append(errs, errors.Inheritf(
					err, "aggregator failed: op=append-default idx=%d backend=%s", i, backendID))
			}
			continue
		}

		raw, err = backendResponse.Body().Raw()
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err,
				"aggregator failed: op=raw idx=%d backend=%s", i, backendID))

			result, err = a.jsonPath.AppendOnArray(result, defaultJSON)
			if checker.NonNil(err) {
				errs = append(errs, errors.Inheritf(err,
					"aggregator failed: op=append-default idx=%d backend=%s", i, backendID))
			}
			continue
		}

		newJsonStr, mergeErrs := a.merge(defaultJSON, backendID, raw)
		if checker.IsNotEmpty(mergeErrs) {
			for _, me := range mergeErrs {
				errs = append(errs, errors.Inheritf(me, "aggregator failed: op=merge idx=%d backend=%s", i, backendID))
			}

			result, err = a.jsonPath.AppendOnArray(result, defaultJSON)
			if checker.NonNil(err) {
				errs = append(errs, errors.Inheritf(err, "aggregator failed: op=append-default idx=%d backend=%s",
					i, backendID))
			}
			continue
		}

		result, err = a.jsonPath.AppendOnArray(result, newJsonStr)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "aggregator failed: op=append-merged idx=%d backend=%s",
				i, backendID))

			result, err = a.jsonPath.AppendOnArray(result, defaultJSON)

			if checker.NonNil(err) {
				errs = append(errs, errors.Inheritf(err, "aggregator failed: op=append-default idx=%d backend=%s",
					i, backendID))
			}
		}
	}

	return a.buildBodyJson(result, errs)
}

func (a aggregatorService) AggregateBodies(history *vo.History) (*vo.Body, []error) {
	result := "{}"

	var errs []error
	for i := 0; checker.IsLessThan(i, history.Size()); i++ {
		backendResponse := history.GetBackendResponse(i)
		if checker.IsNil(backendResponse) || !backendResponse.HasBody() {
			continue
		}
		var err error
		var raw string

		backendID := history.GetBackendID(i)

		raw, err = backendResponse.Body().Raw()
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "aggregator failed: op=raw idx=%d backend=%s", i, backendID))
			continue
		}

		newJsonStr, mergeErrs := a.merge(result, backendID, raw)
		if checker.IsNotEmpty(mergeErrs) {
			for _, me := range mergeErrs {
				errs = append(errs, errors.Inheritf(me, "aggregator failed: op=merge idx=%d backend=%s", i, backendID))
			}
			continue
		}

		result = newJsonStr
	}

	return a.buildBodyJson(result, errs)
}

func (a aggregatorService) buildBodyDefaultForSlice(backendResponse vo.BackendPolymorphicResponse) string {
	jsonStr := "{}"
	jsonStr, _ = a.jsonPath.Set(jsonStr, "ok", converter.ToString(backendResponse.OK()))
	jsonStr, _ = a.jsonPath.Set(jsonStr, "code", backendResponse.StatusCode().String())

	return jsonStr
}

func (a aggregatorService) merge(jsonStr, key, raw string) (string, []error) {
	if checker.IsNotJSON(raw) || checker.IsSlice(raw) {
		return a.mergeJSONInKey(jsonStr, key, raw)
	}
	return a.mergeJSON(jsonStr, raw)
}

func (a aggregatorService) mergeJSONInKey(jsonStr, key, raw string) (string, []error) {
	newJsonStr, err := a.jsonPath.Set(jsonStr, key, raw)
	if checker.NonNil(err) {
		return jsonStr, errors.InheritAsSlicef(err, "aggregator failed: op=set backendKey=%s", key)
	}
	return newJsonStr, nil
}

func (a aggregatorService) mergeJSON(jsonStr, raw string) (string, []error) {
	var result string
	var errs []error

	result = jsonStr
	a.jsonPath.Parse(raw).ForEach(func(key string, value domain.JSONValue) bool {
		newResult, err := a.jsonPath.Add(result, key, value.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "aggregator failed: op=add path=%s", key))
			return true
		}
		result = newResult
		return true
	})

	return result, errs
}

func (a aggregatorService) buildBodyJson(result string, errs []error) (*vo.Body, []error) {
	buffer, err := converter.ToBufferWithErr(result)
	if checker.NonNil(err) {
		return nil, append(errs, errors.Inherit(err, "aggregator failed: op=buffer"))
	}

	return vo.NewBodyWithContentType(vo.NewContentTypeJson(), buffer), errs
}
