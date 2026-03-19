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
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type omitterService struct {
	jsonPath domain.JSONPath
}

type Omitter interface {
	OmitEmptyValuesFromPayload(payload *vo.Payload) (*vo.Payload, []error)
}

func NewOmitter(jsonPath domain.JSONPath) Omitter {
	return omitterService{jsonPath: jsonPath}
}

func (o omitterService) OmitEmptyValuesFromPayload(payload *vo.Payload) (*vo.Payload, []error) {
	if checker.IsNil(payload) {
		return nil, nil
	}

	switch {
	case payload.ContentType().IsNotJSON():
		return o.omitEmptyValuesFromPayloadPlainText(payload)
	case payload.ContentType().IsJSON():
		return o.omitEmptyValuesFromPayloadJSON(payload)
	default:
		return payload, nil
	}
}

func (o omitterService) omitEmptyValuesFromPayloadPlainText(payload *vo.Payload) (*vo.Payload, []error) {
	str, err := payload.String()
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlice(err, "omitter failed: op=payload-string")
	}

	compact := converter.ToCompactString(str)

	buffer, err := converter.ToBufferWithErr(compact)
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlicef(err, "omitter failed: op=buffer-text")
	}

	return vo.NewPayloadWithContentType(payload.ContentType(), buffer), nil
}

func (o omitterService) omitEmptyValuesFromPayloadJSON(payload *vo.Payload) (*vo.Payload, []error) {
	raw, err := payload.Raw()
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlice(err, "omitter failed: op=raw")
	}

	newRaw, errs := o.removeAllEmptyFields(raw)
	if checker.IsNotEmpty(errs) {
		for i, e := range errs {
			errs[i] = errors.Inheritf(e, "omitter failed: op=remove-empty-fields")
		}
	}

	buffer, err := converter.ToBufferWithErr(newRaw)
	if checker.NonNil(err) {
		return payload, append(errs, errors.Inherit(err, "omitter failed: op=buffer-json"))
	}

	return vo.NewPayloadWithContentType(payload.ContentType(), buffer), errs
}

func (o omitterService) removeAllEmptyFields(jsonStr string) (string, []error) {
	root := o.jsonPath.Parse(jsonStr)
	result := jsonStr

	var errs []error

	root.ForEach(func(key string, value domain.JSONValue) bool {
		if value.IsObject() || value.IsArray() {
			childClean, childErrs := o.removeAllEmptyFields(value.Raw())
			if checker.IsNotEmpty(childErrs) {
				errs = append(errs, errors.JoinInheritf(childErrs, ", ", "omitter failed: op=child key=%s", key))
			}
			value = o.jsonPath.Parse(childClean)
		}

		var (
			next string
			err  error
		)

		if checker.IsEmpty(value.Interface()) {
			next, err = o.jsonPath.Delete(result, key)
			if checker.NonNil(err) {
				errs = append(errs, errors.Inheritf(err, "omitter failed: op=delete key=%s", key))
				return true
			}
			result = next
			return true
		}

		next, err = o.jsonPath.Set(result, key, value.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "omitter failed: op=set key=%s", key))
			return true
		}

		result = next
		return true
	})

	return result, errs
}
