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
	domainMapper "github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type mapperService struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Mapper interface {
	MapHeader(mapper *vo.Mapper, header vo.Header, request *vo.HTTPRequest, history *aggregate.History) (vo.Header, error)
	MapQuery(mapper *vo.Mapper, query vo.Query, request *vo.HTTPRequest, history *aggregate.History) (vo.Query, error)
	MapBody(mapper *vo.Mapper, body *vo.Body, request *vo.HTTPRequest, history *aggregate.History) (*vo.Body, []error)
}

func NewMapper(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Mapper {
	return mapperService{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (m mapperService) MapHeader(
	mapper *vo.Mapper,
	header vo.Header,
	request *vo.HTTPRequest,
	history *aggregate.History,
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

	if mapper.ShouldDropUnmapped() {
		return m.mapHeaderDropUnmapped(mapper.Map(), header), nil
	} else {
		return m.mapHeaderKeepAll(mapper.Map(), header), nil
	}
}

func (m mapperService) MapQuery(
	mapper *vo.Mapper,
	query vo.Query,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (vo.Query, error) {
	if checker.IsNil(mapper) || mapper.Map().IsEmpty() {
		return query, nil
	}

	shouldRun, err := m.evalMapperGuards("query", mapper, request, history)
	if checker.NonNil(err) {
		return query, err
	} else if !shouldRun {
		return query, nil
	}

	if mapper.ShouldDropUnmapped() {
		return m.mapQueryDropUnmapped(mapper.Map(), query), nil
	} else {
		return m.mapQueryKeepAll(mapper.Map(), query), nil
	}
}

func (m mapperService) MapBody(
	mapper *vo.Mapper,
	body *vo.Body,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (*vo.Body, []error) {
	if checker.IsNil(mapper) || mapper.Map().IsEmpty() || checker.IsNil(body) {
		return body, nil
	}

	shouldRun, err := m.evalMapperGuards("body", mapper, request, history)
	if checker.NonNil(err) {
		return body, converter.ToSlice(err)
	} else if !shouldRun {
		return body, nil
	}

	if body.ContentType().IsText() {
		return m.mapBodyText(mapper, body)
	} else if body.ContentType().IsJSON() {
		return m.mapBodyJSON(mapper, body)
	}

	return body, nil
}

func (m mapperService) mapHeaderKeepAll(mMap *vo.Map, header vo.Header) vo.Header {
	out := map[string][]string{}

	for _, key := range header.Keys() {
		if domainMapper.IsNotHeaderMandatoryKey(key) && mMap.Exists(key) {
			out[mMap.Get(key)] = header.GetAll(key)
			continue
		}
		out[key] = header.GetAll(key)
	}

	return vo.NewHeader(out)
}

func (m mapperService) mapHeaderDropUnmapped(mMap *vo.Map, header vo.Header) vo.Header {
	out := map[string][]string{}

	for _, key := range header.Keys() {
		if domainMapper.IsHeaderMandatoryKey(key) {
			out[key] = header.GetAll(key)
		} else if mMap.Exists(key) {
			out[mMap.Get(key)] = header.GetAll(key)
		}
	}

	return vo.NewHeader(out)
}

func (m mapperService) mapQueryKeepAll(mMap *vo.Map, query vo.Query) vo.Query {
	out := map[string][]string{}

	for _, key := range query.Keys() {
		if mMap.Exists(key) {
			out[mMap.Get(key)] = query.GetAll(key)
			continue
		}
		out[key] = query.GetAll(key)
	}

	return vo.NewQuery(out)
}

func (m mapperService) mapQueryDropUnmapped(mMap *vo.Map, query vo.Query) vo.Query {
	out := map[string][]string{}

	for _, key := range query.Keys() {
		if mMap.Exists(key) {
			out[mMap.Get(key)] = query.GetAll(key)
		}
	}

	return vo.NewQuery(out)
}

func (m mapperService) mapBodyText(mapper *vo.Mapper, body *vo.Body) (*vo.Body, []error) {
	raw, err := body.String()
	if checker.NonNil(err) {
		return body, errors.InheritAsSlice(err, "mapper failed: kind=text op=stringify-body")
	}

	re, err := m.buildMapperRegex(mapper.Map())
	if checker.IsNil(re) {
		return body, errors.NewAsSlice("mapper failed: kind=text op=build-regex: regex compilation failed")
	}

	var out string
	if mapper.ShouldDropUnmapped() {
		out = m.dropUnmappedText(mapper.Map(), re, raw)
	} else {
		out = m.keepUnmappedText(mapper.Map(), re, raw)
	}

	return m.newBodyWithString(body, out, nil)
}

func (m mapperService) mapBodyJSON(mapper *vo.Mapper, body *vo.Body) (*vo.Body, []error) {
	raw, err := body.String()
	if checker.NonNil(err) {
		return body, errors.InheritAsSlice(err, "mapper failed: kind=json op=stringify-body")
	}

	parsed := m.jsonPath.Parse(raw)

	var out string
	var errs []error

	if parsed.IsArray() {
		out, errs = m.mapJSONArray(mapper, parsed)
	} else {
		out, errs = m.mapJSONObject(mapper, parsed)
	}

	return m.newBodyWithString(body, out, errs)
}

func (m mapperService) mapJSONArray(mapper *vo.Mapper, jsonArray domain.JSONValue) (string, []error) {
	var out = "[]"
	var errs []error
	jsonArray.ForEach(func(idx string, value domain.JSONValue) bool {
		next, itemErrs := m.mapJSONArrayItem(mapper, idx, value, out)
		if checker.IsNotEmpty(itemErrs) {
			errs = append(errs, itemErrs...)
			return true
		}
		out = next
		return true
	})
	return out, errs
}

func (m mapperService) mapJSONArrayItem(
	mapper *vo.Mapper,
	idx string,
	value domain.JSONValue,
	current string,
) (string, []error) {
	var (
		next string
		err  error
	)

	if value.IsObject() {
		obj, childErrs := m.mapJSONObject(mapper, value)
		if checker.IsNotEmpty(childErrs) {
			return current, m.inheritIdxErrs(childErrs, "map-json-object", idx)
		}
		next, err = m.jsonPath.AppendOnArray(current, obj)
	} else if value.IsArray() {
		arr, childErrs := m.mapJSONArray(mapper, value)
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

func (m mapperService) inheritIdxErrs(errs []error, op string, idx string) []error {
	out := make([]error, 0, len(errs))
	for _, e := range errs {
		out = append(out, errors.Inheritf(e, "mapper failed: kind=json op=%s idx=%s", op, idx))
	}
	return out
}

func (m mapperService) mapJSONObject(mapper *vo.Mapper, jsonObject domain.JSONValue) (string, []error) {
	if mapper.ShouldDropUnmapped() {
		return m.mapJSONObjectDropUnmapped(mapper.Map(), jsonObject)
	} else {
		return m.mapJSONObjectKeepAll(mapper.Map(), jsonObject)
	}
}

func (m mapperService) mapJSONObjectKeepAll(mMap *vo.Map, jsonObject domain.JSONValue) (string, []error) {
	mapped := jsonObject.Raw()
	var errs []error

	for _, key := range mMap.Keys() {
		newKey := mMap.Get(key)
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

func (m mapperService) mapJSONObjectDropUnmapped(mMap *vo.Map, jsonObject domain.JSONValue) (string, []error) {
	mapped := "{}"
	var errs []error

	for _, key := range mMap.Keys() {
		newKey := mMap.Get(key)

		val := jsonObject.Get(key)
		if val.NotExists() {
			continue
		}

		out, err := m.jsonPath.Set(mapped, newKey, val.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(
				err,
				"mapper failed: kind=json op=set policy=drop-unmapped from=%s to=%s",
				key,
				newKey,
			))
			continue
		}

		mapped = out
	}

	return mapped, errs
}

func (m mapperService) buildMapperRegex(mMap *vo.Map) (*regexp.Regexp, error) {
	keys := mMap.Keys()

	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, regexp.QuoteMeta(k))
	}

	return regexp.Compile(strings.Join(parts, "|"))
}

func (m mapperService) keepUnmappedText(mMap *vo.Map, re *regexp.Regexp, raw string) string {
	return re.ReplaceAllStringFunc(raw, func(match string) string {
		if mMap.Exists(match) {
			return mMap.Get(match)
		}
		return match
	})
}

func (m mapperService) dropUnmappedText(mMap *vo.Map, re *regexp.Regexp, raw string) string {
	var b strings.Builder
	for _, match := range re.FindAllString(raw, -1) {
		if mMap.Exists(match) {
			b.WriteString(mMap.Get(match))
		}
	}
	return b.String()
}

func (m mapperService) evalMapperGuards(kind string, mapper *vo.Mapper, request *vo.HTTPRequest, history *aggregate.History) (
	bool, error) {
	shouldRun, _, errs := m.dynamicValueService.EvalGuards(mapper.OnlyIf(), mapper.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInheritf(errs, ", ", "failed to evaluate guard for %s mapper", kind)
	}
	return shouldRun, nil
}

func (m mapperService) newBodyWithString(body *vo.Body, str string, errs []error) (*vo.Body, []error) {
	buffer, err := converter.ToBufferWithErr(str)
	if checker.NonNil(err) {
		return body, append(errs, errors.Inherit(err, "mapper failed: op=buffer"))
	}
	return vo.NewBodyWithContentType(body.ContentType(), buffer), errs
}
