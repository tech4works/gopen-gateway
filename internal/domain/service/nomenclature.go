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

type nomenclatureService struct {
	jsonPath     domain.JSONPath
	nomenclature domain.Nomenclature
}

type Nomenclature interface {
	ToCase(nomenclature enum.Nomenclature, body *vo.Body) (*vo.Body, []error)
}

func NewNomenclature(jsonPath domain.JSONPath, nomenclature domain.Nomenclature) Nomenclature {
	return nomenclatureService{
		jsonPath:     jsonPath,
		nomenclature: nomenclature,
	}
}

func (n nomenclatureService) ToCase(nm enum.Nomenclature, body *vo.Body) (*vo.Body, []error) {
	if !nm.IsEnumValid() {
		return body, nil
	} else if checker.IsNil(body) || body.ContentType().IsNotJSON() {
		return body, nil
	}

	raw, err := body.Raw()
	if checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "nomenclature failed: op=raw case=%s", nm)
	}

	jsonStr, errs := n.convertKeysToCase(raw, nm)

	buffer, err := converter.ToBufferWithErr(jsonStr)
	if checker.NonNil(err) {
		return body, append(errs, errors.Inheritf(err, "nomenclature failed: op=buffer case=%s", nm))
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), errs
}

func (n nomenclatureService) convertKeysToCase(raw string, nm enum.Nomenclature) (string, []error) {
	root := n.jsonPath.Parse(raw)

	result := "{}"
	if root.IsArray() {
		result = "[]"
	}

	var allErrs []error
	root.ForEach(func(key string, value domain.JSONValue) bool {
		var err error

		newKey := n.nomenclature.Parse(nm, key)

		valueResolved, errs := n.resolveValueIfDeepJSON(value, nm)
		if checker.IsNotEmpty(errs) {
			allErrs = append(allErrs, errs...)
		}

		result, err = n.jsonPath.Set(result, newKey, valueResolved)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "nomenclature failed: op=set case=%s key=%s newKey=%s type=%s",
				nm, key, newKey, value.Type()))
		}

		return true
	})
	return result, allErrs
}

func (n nomenclatureService) resolveValueIfDeepJSON(value domain.JSONValue, nm enum.Nomenclature) (string, []error) {
	raw := value.Raw()

	if value.IsObject() || value.IsArray() {
		return n.convertKeysToCase(raw, nm)
	}

	return raw, nil
}
