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
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type nomenclature struct {
	jsonPath domain.JSONPath
	provider domain.Nomenclature
}

type Nomenclature interface {
	ToCase(config enum.Nomenclature, payload *vo.Payload) (*vo.Payload, []error)
}

func NewNomenclature(jsonPath domain.JSONPath, provider domain.Nomenclature) Nomenclature {
	return nomenclature{
		jsonPath: jsonPath,
		provider: provider,
	}
}

func (n nomenclature) ToCase(config enum.Nomenclature, payload *vo.Payload) (*vo.Payload, []error) {
	if !config.IsEnumValid() || checker.IsNil(payload) || payload.ContentType().IsNotJSON() {
		return payload, nil
	}

	raw, err := payload.Raw()
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlicef(err, "nomenclature failed: op=raw case=%s", config)
	}

	jsonStr, errs := n.convertKeysToCase(config, raw)

	buffer, err := converter.ToBufferWithErr(jsonStr)
	if checker.NonNil(err) {
		return payload, append(errs, errors.Inheritf(err, "nomenclature failed: op=buffer case=%s", config))
	}

	return vo.NewPayloadWithContentType(payload.ContentType(), buffer), errs
}

func (n nomenclature) convertKeysToCase(config enum.Nomenclature, raw string) (string, []error) {
	root := n.jsonPath.Parse(raw)

	result := "{}"
	if root.IsArray() {
		result = "[]"
	}

	var allErrs []error
	root.ForEach(func(key string, value domain.JSONValue) bool {
		var err error

		newKey := n.provider.Parse(config, key)

		valueResolved, errs := n.resolveValueIfDeepJSON(config, value)
		if checker.IsNotEmpty(errs) {
			allErrs = append(allErrs, errs...)
		}

		result, err = n.jsonPath.Set(result, newKey, valueResolved)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "nomenclature failed: op=set case=%s key=%s newKey=%s type=%s",
				config, key, newKey, value.Type()))
		}

		return true
	})
	return result, allErrs
}

func (n nomenclature) resolveValueIfDeepJSON(config enum.Nomenclature, value domain.JSONValue) (string, []error) {
	raw := value.Raw()

	if value.IsObject() || value.IsArray() {
		return n.convertKeysToCase(config, raw)
	}

	return raw, nil
}
