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
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type aggregator struct {
	jsonPath domain.JSONPath
}

type Aggregator interface {
	AggregateMetadata(base, value vo.Metadata, ignoreKeys []string) vo.Metadata
	AggregatePayloadToKey(key string, payload *vo.Payload) (*vo.Payload, error)
	AggregatePayloadsIntoSlice(history *aggregate.History) (*vo.Payload, []error)
	AggregatePayloads(history *aggregate.History) (*vo.Payload, []error)
}

func NewAggregator(jsonPath domain.JSONPath) Aggregator {
	return aggregator{
		jsonPath: jsonPath,
	}
}

func (a aggregator) AggregateMetadata(base, value vo.Metadata, ignoreKeys []string) vo.Metadata {
	aggregated := base.Copy()
	for _, key := range value.Keys() {
		if checker.NotContains(ignoreKeys, key) {
			aggregated[key] = append(aggregated[key], value.GetAll(key)...)
		}
	}
	return vo.NewMetadata(aggregated)
}

func (a aggregator) AggregatePayloadToKey(key string, payload *vo.Payload) (*vo.Payload, error) {
	if checker.IsEmpty(key) ||
		checker.IsNil(payload) ||
		(payload.ContentType().IsNotJSON() && payload.ContentType().IsNotPlainText()) {
		return payload, nil
	}

	raw, err := payload.Raw()
	if checker.NonNil(err) {
		return payload, errors.Inheritf(err, "aggregator failed: op=raw key=%s", key)
	}

	jsonValue, err := a.jsonPath.Set("{}", key, raw)
	if checker.NonNil(err) {
		return payload, errors.Inheritf(err, "aggregator failed: op=set key=%s", key)
	}

	buffer, err := converter.ToBufferWithErr(jsonValue)
	if checker.NonNil(err) {
		return payload, errors.Inheritf(err, "aggregator failed: op=buffer key=%s", key)
	}

	return vo.NewPayloadJSON(buffer), nil
}

func (a aggregator) AggregatePayloadsIntoSlice(history *aggregate.History) (*vo.Payload, []error) {
	result := "[]"

	var errs []error
	for i := 0; checker.IsLessThan(i, history.Size()); i++ {
		backendResponse := history.GetResponse(i)
		if !backendResponse.ShouldInFinalResponse() {
			continue
		}

		var err error
		var raw string

		backendID := history.GetID(i)
		defaultJSON := a.buildPayloadDefaultForSlice(backendResponse)

		if !backendResponse.HasBody() {
			result, err = a.jsonPath.AppendOnArray(result, defaultJSON)
			if checker.NonNil(err) {
				errs = append(errs, errors.Inheritf(err, "aggregator failed: op=append-default idx=%d backend=%s",
					i, backendID))
			}
			continue
		}

		raw, err = backendResponse.Payload().Raw()
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "aggregator failed: op=raw idx=%d backend=%s", i, backendID))

			result, err = a.jsonPath.AppendOnArray(result, defaultJSON)
			if checker.NonNil(err) {
				errs = append(errs, errors.Inheritf(err, "aggregator failed: op=append-default idx=%d backend=%s",
					i, backendID))
			}
			continue
		}

		newJSONStr, mergeErrs := a.merge(defaultJSON, backendID, raw)
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

		result, err = a.jsonPath.AppendOnArray(result, newJSONStr)
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

	return a.buildPayloadJSON(result, errs)
}

func (a aggregator) AggregatePayloads(history *aggregate.History) (*vo.Payload, []error) {
	result := "{}"

	var errs []error
	for i := 0; checker.IsLessThan(i, history.Size()); i++ {
		backendResponse := history.GetResponse(i)
		if !backendResponse.ShouldInFinalResponse() {
			continue
		}

		var err error
		var raw string

		backendID := history.GetID(i)

		raw, err = backendResponse.Payload().Raw()
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "aggregator failed: op=raw idx=%d backend=%s", i, backendID))
			continue
		}

		newJSONStr, mergeErrs := a.merge(result, backendID, raw)
		if checker.IsNotEmpty(mergeErrs) {
			for _, me := range mergeErrs {
				errs = append(errs, errors.Inheritf(me, "aggregator failed: op=merge idx=%d backend=%s", i, backendID))
			}
			continue
		}

		result = newJSONStr
	}

	return a.buildPayloadJSON(result, errs)
}

func (a aggregator) buildPayloadDefaultForSlice(backendResponse *vo.BackendResponse) string {
	jsonStr := "{}"
	jsonStr, _ = a.jsonPath.Set(jsonStr, "ok", converter.ToString(backendResponse.OK()))
	jsonStr, _ = a.jsonPath.Set(jsonStr, "code", backendResponse.Status().String())

	return jsonStr
}

func (a aggregator) merge(jsonStr, key, raw string) (string, []error) {
	if checker.IsNotJSON(raw) || checker.IsSlice(raw) {
		return a.mergeJSONInKey(jsonStr, key, raw)
	}
	return a.mergeJSON(jsonStr, raw)
}

func (a aggregator) mergeJSONInKey(jsonStr, key, raw string) (string, []error) {
	newJSONStr, err := a.jsonPath.Set(jsonStr, key, raw)
	if checker.NonNil(err) {
		return jsonStr, errors.InheritAsSlicef(err, "aggregator failed: op=set backendKey=%s", key)
	}
	return newJSONStr, nil
}

func (a aggregator) mergeJSON(jsonStr, raw string) (string, []error) {
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

func (a aggregator) buildPayloadJSON(result string, errs []error) (*vo.Payload, []error) {
	buffer, err := converter.ToBufferWithErr(result)
	if checker.NonNil(err) {
		return nil, append(errs, errors.Inherit(err, "aggregator failed: op=buffer"))
	}

	return vo.NewPayloadJSON(buffer), errs
}
