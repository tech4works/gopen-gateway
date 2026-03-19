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
	"regexp"
	"sort"
	"strings"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type mapper struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Mapper interface {
	MapMetadata(
		config *vo.MapperConfig,
		metadata vo.Metadata,
		ignoreKeys []string,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (vo.Metadata, error)
	MapQuery(
		config *vo.MapperConfig,
		query vo.Query,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (vo.Query, error)
	MapPayload(
		config *vo.MapperConfig,
		payload *vo.Payload,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (*vo.Payload, []error)
}

func NewMapper(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Mapper {
	return mapper{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (m mapper) MapMetadata(
	config *vo.MapperConfig,
	metadata vo.Metadata,
	ignoreKeys []string,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Metadata, error) {
	if checker.IsNil(config) || config.Map().IsEmpty() {
		return metadata, nil
	}

	shouldRun, err := m.evalMapperGuards(config, "metadata", request, history)
	if checker.NonNil(err) {
		return metadata, err
	} else if !shouldRun {
		return metadata, nil
	}

	if config.ShouldDropUnmapped() {
		return m.mapMetadataDropUnmapped(config.Map(), metadata, ignoreKeys), nil
	} else {
		return m.mapMetadataKeepAll(config.Map(), metadata, ignoreKeys), nil
	}
}

func (m mapper) MapQuery(
	config *vo.MapperConfig,
	query vo.Query,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Query, error) {
	if checker.IsNil(config) || config.Map().IsEmpty() {
		return query, nil
	}

	shouldRun, err := m.evalMapperGuards(config, "query", request, history)
	if checker.NonNil(err) {
		return query, err
	} else if !shouldRun {
		return query, nil
	}

	if config.ShouldDropUnmapped() {
		return m.mapQueryDropUnmapped(config.Map(), query), nil
	} else {
		return m.mapQueryKeepAll(config.Map(), query), nil
	}
}

func (m mapper) MapPayload(
	config *vo.MapperConfig,
	payload *vo.Payload,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.Payload, []error) {
	if checker.IsNil(config) || config.Map().IsEmpty() || checker.IsNil(payload) {
		return payload, nil
	}

	shouldRun, err := m.evalMapperGuards(config, "payload", request, history)
	if checker.NonNil(err) {
		return payload, converter.ToSlice(err)
	} else if !shouldRun {
		return payload, nil
	}

	if payload.ContentType().IsPlainText() {
		return m.mapPayloadPlainText(config, payload)
	} else if payload.ContentType().IsJSON() {
		return m.mapPayloadJSON(config, payload)
	}

	return payload, nil
}

func (m mapper) mapMetadataKeepAll(config *vo.MapConfig, header vo.Metadata, ignoreKeys []string) vo.Metadata {
	out := map[string][]string{}

	for _, key := range header.Keys() {
		if checker.NotContains(ignoreKeys, key) && config.Exists(key) {
			out[config.Get(key)] = header.GetAll(key)
			continue
		}
		out[key] = header.GetAll(key)
	}

	return vo.NewMetadata(out)
}

func (m mapper) mapMetadataDropUnmapped(config *vo.MapConfig, header vo.Metadata, ignoreKeys []string) vo.Metadata {
	out := map[string][]string{}

	for _, key := range header.Keys() {
		if checker.Contains(ignoreKeys, key) {
			out[key] = header.GetAll(key)
		} else if config.Exists(key) {
			out[config.Get(key)] = header.GetAll(key)
		}
	}

	return vo.NewMetadata(out)
}

func (m mapper) mapQueryKeepAll(config *vo.MapConfig, query vo.Query) vo.Query {
	out := map[string][]string{}

	for _, key := range query.Keys() {
		if config.Exists(key) {
			out[config.Get(key)] = query.GetAll(key)
			continue
		}
		out[key] = query.GetAll(key)
	}

	return vo.NewQuery(out)
}

func (m mapper) mapQueryDropUnmapped(config *vo.MapConfig, query vo.Query) vo.Query {
	out := map[string][]string{}

	for _, key := range query.Keys() {
		if config.Exists(key) {
			out[config.Get(key)] = query.GetAll(key)
		}
	}

	return vo.NewQuery(out)
}

func (m mapper) mapPayloadPlainText(config *vo.MapperConfig, payload *vo.Payload) (*vo.Payload, []error) {
	raw, err := payload.String()
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlice(err, "mapper failed: kind=text op=stringify-payload")
	}

	re, err := m.buildMapperRegex(config.Map())
	if checker.IsNil(re) {
		return payload, errors.NewAsSlice("mapper failed: kind=text op=build-regex: regex compilation failed")
	}

	var out string
	if config.ShouldDropUnmapped() {
		out = m.dropUnmappedPlainText(config.Map(), re, raw)
	} else {
		out = m.keepUnmappedPlainText(config.Map(), re, raw)
	}

	return m.newPayloadWithString(payload, out, nil)
}

func (m mapper) mapPayloadJSON(config *vo.MapperConfig, payload *vo.Payload) (*vo.Payload, []error) {
	raw, err := payload.String()
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlice(err, "mapper failed: kind=json op=stringify-payload")
	}

	parsed := m.jsonPath.Parse(raw)

	var out string
	var errs []error

	if parsed.IsArray() {
		out, errs = m.mapJSONArray(config, parsed)
	} else {
		out, errs = m.mapJSONObject(config, parsed)
	}

	return m.newPayloadWithString(payload, out, errs)
}

func (m mapper) mapJSONArray(config *vo.MapperConfig, jsonArray domain.JSONValue) (string, []error) {
	var out = "[]"
	var errs []error
	jsonArray.ForEach(func(idx string, value domain.JSONValue) bool {
		next, itemErrs := m.mapJSONArrayItem(config, idx, value, out)
		if checker.IsNotEmpty(itemErrs) {
			errs = append(errs, itemErrs...)
			return true
		}
		out = next
		return true
	})
	return out, errs
}

func (m mapper) mapJSONArrayItem(
	config *vo.MapperConfig,
	idx string,
	value domain.JSONValue,
	current string,
) (string, []error) {
	var (
		next string
		err  error
	)

	if value.IsObject() {
		obj, childErrs := m.mapJSONObject(config, value)
		if checker.IsNotEmpty(childErrs) {
			return current, m.inheritIdxErrs(childErrs, "map-json-object", idx)
		}
		next, err = m.jsonPath.AppendOnArray(current, obj)
	} else if value.IsArray() {
		arr, childErrs := m.mapJSONArray(config, value)
		if checker.IsNotEmpty(childErrs) {
			return current, m.inheritIdxErrs(childErrs, "map-json-array", idx)
		}
		next, err = m.jsonPath.AppendOnArray(current, arr)
	} else {
		next, err = m.jsonPath.AppendOnArray(current, value.Raw())
	}
	if checker.NonNil(err) {
		return current, errors.InheritAsSlice(err, "mapper failed: kind=json op=append-on-array idx=%s", idx)
	}

	return next, nil
}

func (m mapper) inheritIdxErrs(errs []error, op string, idx string) []error {
	out := make([]error, 0, len(errs))
	for _, e := range errs {
		out = append(out, errors.Inheritf(e, "mapper failed: kind=json op=%s idx=%s", op, idx))
	}
	return out
}

func (m mapper) mapJSONObject(config *vo.MapperConfig, jsonObject domain.JSONValue) (string, []error) {
	if config.ShouldDropUnmapped() {
		return m.mapJSONObjectDropUnmapped(config.Map(), jsonObject)
	} else {
		return m.mapJSONObjectKeepAll(config.Map(), jsonObject)
	}
}

func (m mapper) mapJSONObjectKeepAll(config *vo.MapConfig, jsonObject domain.JSONValue) (string, []error) {
	mapped := jsonObject.Raw()
	var errs []error

	for _, key := range config.Keys() {
		newKey := config.Get(key)
		if checker.Equals(key, newKey) {
			continue
		}

		val := jsonObject.Get(key)
		if val.NotExists() {
			continue
		}

		out, err := m.jsonPath.Set(mapped, newKey, val.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "mapper failed: kind=json op=set from=%s to=%s", key, newKey))
			continue
		}

		out, err = m.jsonPath.Delete(out, key)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "mapper failed: kind=json op=delete from=%s to=%s", key, newKey))
			continue
		}

		mapped = out
	}

	return mapped, errs
}

func (m mapper) mapJSONObjectDropUnmapped(config *vo.MapConfig, jsonObject domain.JSONValue) (string, []error) {
	mapped := "{}"
	var errs []error

	for _, key := range config.Keys() {
		newKey := config.Get(key)

		val := jsonObject.Get(key)
		if val.NotExists() {
			continue
		}

		out, err := m.jsonPath.Set(mapped, newKey, val.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "mapper failed: kind=json op=set policy=drop-unmapped from=%s to=%s",
				key, newKey))
			continue
		}

		mapped = out
	}

	return mapped, errs
}

func (m mapper) buildMapperRegex(config *vo.MapConfig) (*regexp.Regexp, error) {
	keys := config.Keys()

	sort.Slice(keys, func(i, j int) bool {
		return checker.IsGreaterThan(len(keys[i]), len(keys[j]))
	})

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, regexp.QuoteMeta(k))
	}

	return regexp.Compile(strings.Join(parts, "|"))
}

func (m mapper) keepUnmappedPlainText(config *vo.MapConfig, re *regexp.Regexp, raw string) string {
	return re.ReplaceAllStringFunc(raw, func(match string) string {
		if config.Exists(match) {
			return config.Get(match)
		}
		return match
	})
}

func (m mapper) dropUnmappedPlainText(config *vo.MapConfig, re *regexp.Regexp, raw string) string {
	var b strings.Builder
	for _, match := range re.FindAllString(raw, -1) {
		if config.Exists(match) {
			b.WriteString(config.Get(match))
		}
	}
	return b.String()
}

func (m mapper) evalMapperGuards(config *vo.MapperConfig, kind string, request *vo.EndpointRequest,
	history *aggregate.History) (bool, error) {
	shouldRun, _, errs := m.dynamicValueService.EvalGuards(config.OnlyIf(), config.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInheritf(errs, ", ", "failed to evaluate guard for %s mapper", kind)
	}
	return shouldRun, nil
}

func (m mapper) newPayloadWithString(payload *vo.Payload, str string, errs []error) (*vo.Payload, []error) {
	buffer, err := converter.ToBufferWithErr(str)
	if checker.NonNil(err) {
		return payload, append(errs, errors.Inherit(err, "mapper failed: op=buffer"))
	}
	return vo.NewPayloadWithContentType(payload.ContentType(), buffer), errs
}
