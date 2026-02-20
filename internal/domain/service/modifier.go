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
	"fmt"
	"strings"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type modifierService struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Modifier interface {
	ExecuteURLPathModifiers(modifiers []vo.Modifier, urlPath vo.URLPath, request *vo.HTTPRequest, history *vo.History) (vo.URLPath, []error)
	ExecuteHeaderModifiers(modifiers []vo.Modifier, header vo.Header, request *vo.HTTPRequest, history *vo.History) (vo.Header, []error)
	ExecuteQueryModifiers(modifiers []vo.Modifier, query vo.Query, request *vo.HTTPRequest, history *vo.History) (vo.Query, []error)
	ExecuteBodyModifiers(modifiers []vo.Modifier, body *vo.Body, request *vo.HTTPRequest, history *vo.History) (*vo.Body, []error)

	ModifyURLPath(modifier *vo.Modifier, urlPath vo.URLPath, request *vo.HTTPRequest, history *vo.History) (vo.URLPath, error)
	ModifyHeader(modifier *vo.Modifier, header vo.Header, request *vo.HTTPRequest, history *vo.History) (vo.Header, error)
	ModifyQuery(modifier *vo.Modifier, query vo.Query, request *vo.HTTPRequest, history *vo.History) (vo.Query, error)
	ModifyBody(modifier *vo.Modifier, body *vo.Body, request *vo.HTTPRequest, history *vo.History) (*vo.Body, error)
}

func NewModifier(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Modifier {
	return modifierService{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (s modifierService) ExecuteBodyModifiers(modifiers []vo.Modifier, body *vo.Body, request *vo.HTTPRequest,
	history *vo.History) (*vo.Body, []error) {
	var allErrs []error
	for _, modifier := range modifiers {
		var err error

		body, err = s.ModifyBody(&modifier, body, request, history)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "modifier failed: op=execute kind=body action=%s key=%s value=%s",
				modifier.Action(), modifier.Key(), modifier.Value()))
		}
	}

	return body, allErrs
}

func (s modifierService) ExecuteURLPathModifiers(modifiers []vo.Modifier, urlPath vo.URLPath, request *vo.HTTPRequest,
	history *vo.History) (vo.URLPath, []error) {
	var allErrs []error

	for _, modifier := range modifiers {
		var err error

		urlPath, err = s.ModifyURLPath(&modifier, urlPath, request, history)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "modifier failed: op=execute kind=url-path action=%s key=%s value=%s",
				modifier.Action(), modifier.Key(), modifier.Value()))
		}
	}

	return urlPath, allErrs
}

func (s modifierService) ExecuteHeaderModifiers(modifiers []vo.Modifier, header vo.Header, request *vo.HTTPRequest,
	history *vo.History) (vo.Header, []error) {
	var allErrs []error

	for _, modifier := range modifiers {
		var err error

		header, err = s.ModifyHeader(&modifier, header, request, history)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "modifier failed: op=execute kind=header action=%s key=%s value=%s",
				modifier.Action(), modifier.Key(), modifier.Value()))
		}
	}

	return header, allErrs
}

func (s modifierService) ExecuteQueryModifiers(modifiers []vo.Modifier, query vo.Query, request *vo.HTTPRequest,
	history *vo.History) (vo.Query, []error) {
	var allErrs []error

	for _, modifier := range modifiers {
		var err error

		query, err = s.ModifyQuery(&modifier, query, request, history)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "modifier failed: op=execute kind=query action=%s key=%s value=%s",
				modifier.Action(), modifier.Key(), modifier.Value()))
		}
	}

	return query, allErrs
}

func (s modifierService) ModifyURLPath(modifier *vo.Modifier, urlPath vo.URLPath, request *vo.HTTPRequest,
	history *vo.History) (vo.URLPath, error) {
	shouldRun, err := s.evalModifierGuards("url path", modifier, request, history)
	if checker.NonNil(err) {
		return urlPath, s.wrapModifierErr("url-path", "eval-guards", modifier, err, "")
	} else if !shouldRun {
		return urlPath, nil
	}

	action := modifier.Action()
	key := modifier.Key()
	value, errs := s.dynamicValueService.Get(modifier.Value(), request, history)
	if checker.IsNotEmpty(errs) {
		return urlPath, s.joinDynamicValueErr("url-path", modifier, errs)
	}

	switch action {
	case enum.ModifierActionSet:
		return s.setURLPath(urlPath, key, value)
	case enum.ModifierActionRpl:
		return s.replaceURLPath(urlPath, key, value)
	case enum.ModifierActionDel:
		return s.deleteURLPath(urlPath, key)
	default:
		return urlPath, mapper.NewErrModifierActionNotImplemented("url-path", action)
	}
}

func (s modifierService) ModifyHeader(modifier *vo.Modifier, header vo.Header, request *vo.HTTPRequest,
	history *vo.History) (vo.Header, error) {
	shouldRun, err := s.evalModifierGuards("header", modifier, request, history)
	if checker.NonNil(err) {
		return header, s.wrapModifierErr("header", "eval-guards", modifier, err, "")
	} else if !shouldRun {
		return header, nil
	}

	action := modifier.Action()
	key := modifier.Key()
	values, errs := s.dynamicValueService.GetAsSliceOfString(modifier.Value(), request, history)
	if checker.IsNotEmpty(errs) {
		return header, s.joinDynamicValueErr("header", modifier, errs)
	}

	switch action {
	case enum.ModifierActionAdd:
		return s.addHeader(header, key, values)
	case enum.ModifierActionApd:
		return s.appendHeader(header, key, values)
	case enum.ModifierActionSet:
		return s.setHeader(header, key, values)
	case enum.ModifierActionRpl:
		return s.replaceHeader(header, key, values)
	case enum.ModifierActionDel:
		return s.deleteHeader(header, key)
	default:
		return header, mapper.NewErrModifierActionNotImplemented("header", action)
	}
}

func (s modifierService) ModifyQuery(modifier *vo.Modifier, query vo.Query, request *vo.HTTPRequest, history *vo.History,
) (vo.Query, error) {
	shouldRun, err := s.evalModifierGuards("query", modifier, request, history)
	if checker.NonNil(err) {
		return query, s.wrapModifierErr("query", "eval-guards", modifier, err, "")
	} else if !shouldRun {
		return query, nil
	}

	action := modifier.Action()
	key := modifier.Key()
	values, errs := s.dynamicValueService.GetAsSliceOfString(modifier.Value(), request, history)
	if checker.IsNotEmpty(errs) {
		return query, s.joinDynamicValueErr("query", modifier, errs)
	}

	switch action {
	case enum.ModifierActionAdd:
		return s.addQuery(query, key, values)
	case enum.ModifierActionApd:
		return s.appendQuery(query, key, values)
	case enum.ModifierActionSet:
		return s.setQuery(query, key, values)
	case enum.ModifierActionRpl:
		return s.replaceQuery(query, key, values)
	case enum.ModifierActionDel:
		return s.deleteQuery(query, key)
	default:
		return query, mapper.NewErrModifierActionNotImplemented("query", action)
	}
}

func (s modifierService) ModifyBody(modifier *vo.Modifier, body *vo.Body, request *vo.HTTPRequest, history *vo.History,
) (*vo.Body, error) {
	if checker.IsNil(body) {
		return nil, nil
	}

	shouldRun, err := s.evalModifierGuards("body", modifier, request, history)
	if checker.NonNil(err) {
		return body, s.wrapModifierErr("body", "eval-guards", modifier, err, "")
	} else if !shouldRun {
		return body, nil
	}

	action := modifier.Action()
	key := modifier.Key()
	value, errs := s.dynamicValueService.Get(modifier.Value(), request, history)
	if checker.IsNotEmpty(errs) {
		return body, s.joinDynamicValueErr("body", modifier, errs)
	}

	switch action {
	case enum.ModifierActionAdd:
		return s.addBody(body, key, value)
	case enum.ModifierActionApd:
		return s.appendBody(body, key, value)
	case enum.ModifierActionSet:
		return s.setBody(body, key, value)
	case enum.ModifierActionRpl:
		return s.replaceBody(body, key, value)
	case enum.ModifierActionDel:
		return s.deleteBody(body, key)
	default:
		return body, mapper.NewErrModifierActionNotImplemented("body", action)
	}
}

func (s modifierService) evalModifierGuards(kind string, modifier *vo.Modifier, request *vo.HTTPRequest, history *vo.History,
) (bool, error) {
	shouldRun, _, errs := s.dynamicValueService.EvalGuards(modifier.OnlyIf(), modifier.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInheritf(errs, ", ", "modifier failed: op=eval-guards kind=%s action=%s key=%s",
			kind, modifier.Action(), modifier.Key())
	}
	return shouldRun, nil
}

func (s modifierService) setURLPath(urlPath vo.URLPath, key, value string) (vo.URLPath, error) {
	path := urlPath.Raw()
	paramValues := urlPath.Params().Copy()

	paramValues[key] = value
	if checker.NotContains(path, fmt.Sprintf("/:%s", key)) {
		path = fmt.Sprintf("%s/:%s", path, key)
	}

	return vo.NewURLPath(path, paramValues), nil
}

func (s modifierService) replaceURLPath(urlPath vo.URLPath, key, value string) (vo.URLPath, error) {
	if urlPath.NotExists(key) {
		return urlPath, nil
	}

	return s.setURLPath(urlPath, key, value)
}

func (s modifierService) deleteURLPath(urlPath vo.URLPath, key string) (vo.URLPath, error) {
	path := strings.ReplaceAll(urlPath.Raw(), fmt.Sprintf("/:%s", key), "")

	paramValues := urlPath.Params().Copy()
	delete(paramValues, key)

	return vo.NewURLPath(path, paramValues), nil
}

func (s modifierService) addHeader(header vo.Header, key string, value []string) (vo.Header, error) {
	if mapper.IsHeaderMandatoryKey(key) {
		return header, nil
	}

	values := header.Copy()
	values[key] = append(header.GetAll(key), value...)

	return vo.NewHeader(values), nil
}

func (s modifierService) appendHeader(header vo.Header, key string, value []string) (vo.Header, error) {
	if mapper.IsHeaderMandatoryKey(key) || header.NotExists(key) {
		return header, nil
	}

	values := header.Copy()
	values[key] = append(header.GetAll(key), value...)

	return vo.NewHeader(values), nil
}

func (s modifierService) setHeader(header vo.Header, key string, value []string) (vo.Header, error) {
	if mapper.IsHeaderMandatoryKey(key) {
		return header, nil
	}

	values := header.Copy()
	values[key] = value

	return vo.NewHeader(values), nil
}

func (s modifierService) replaceHeader(header vo.Header, key string, value []string) (vo.Header, error) {
	if mapper.IsHeaderMandatoryKey(key) || header.NotExists(key) {
		return header, nil
	}

	return s.setHeader(header, key, value)
}

func (s modifierService) deleteHeader(header vo.Header, key string) (vo.Header, error) {
	if mapper.IsHeaderMandatoryKey(key) {
		return header, nil
	}

	values := header.Copy()
	delete(values, key)

	return vo.NewHeader(values), nil
}

func (s modifierService) addQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	values := query.Copy()
	values[key] = append(query.GetAll(key), value...)

	return vo.NewQuery(values), nil
}

func (s modifierService) appendQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	if query.NotExists(key) {
		return query, nil
	}

	values := query.Copy()
	values[key] = append(query.GetAll(key), value...)

	return vo.NewQuery(values), nil
}

func (s modifierService) setQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	values := query.Copy()
	values[key] = value

	return vo.NewQuery(values), nil
}

func (s modifierService) replaceQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	if query.NotExists(key) {
		return query, nil
	}

	return s.setQuery(query, key, value)
}

func (s modifierService) deleteQuery(query vo.Query, key string) (vo.Query, error) {
	values := query.Copy()

	delete(values, key)

	return vo.NewQuery(values), nil
}

func (s modifierService) addBody(body *vo.Body, key, value string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.addBodyText(body, value)
	} else if body.ContentType().IsJSON() {
		return s.addBodyJson(body, key, value)
	} else {
		return body, mapper.NewErrModifierIncompatibleBodyType("add-body", body.ContentType().String())
	}
}

func (s modifierService) addBodyText(body *vo.Body, value string) (*vo.Body, error) {
	bodyStr, err := body.String()
	if checker.NonNil(err) {
		return body, err
	}

	modifiedBodyText := fmt.Sprintf("%s%s", bodyStr, value)
	return s.newBodyByString(body, modifiedBodyText)
}

func (s modifierService) addBodyJson(body *vo.Body, key, value string) (*vo.Body, error) {
	bodyRaw, err := body.Raw()
	if checker.NonNil(err) {
		return body, err
	}

	modifiedBodyJson, err := s.jsonPath.Add(bodyRaw, key, value)
	if checker.NonNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) appendBody(body *vo.Body, key, value string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.appendBodyText(body, value)
	} else if body.ContentType().IsJSON() {
		return s.appendBodyJson(body, key, value)
	}
	return body, mapper.NewErrModifierIncompatibleBodyType("append-body", body.ContentType().String())
}

func (s modifierService) appendBodyText(body *vo.Body, value string) (*vo.Body, error) {
	bodyStr, err := body.String()
	if checker.NonNil(err) {
		return body, err
	}

	modifiedBodyText := fmt.Sprintf("%s\n%s", bodyStr, value)
	return s.newBodyByString(body, modifiedBodyText)
}

func (s modifierService) appendBodyJson(body *vo.Body, key, value string) (*vo.Body, error) {
	bodyRaw, err := body.Raw()
	if checker.NonNil(err) {
		return body, err
	}

	if s.jsonPath.Parse(bodyRaw).Get(key).NotExists() {
		return body, nil
	}

	modifiedBodyJson, err := s.jsonPath.Add(bodyRaw, key, value)
	if checker.NonNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) setBody(body *vo.Body, key, value string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.setBodyText(body, value)
	} else if body.ContentType().IsJSON() {
		return s.setBodyJson(body, key, value)
	}
	return body, mapper.NewErrModifierIncompatibleBodyType("set-body", body.ContentType().String())
}

func (s modifierService) setBodyText(body *vo.Body, value string) (*vo.Body, error) {
	return s.newBodyByString(body, value)
}

func (s modifierService) setBodyJson(body *vo.Body, key, value string) (*vo.Body, error) {
	bodyRaw, err := body.Raw()
	if checker.NonNil(err) {
		return body, err
	}

	modifiedBodyJson, err := s.jsonPath.Set(bodyRaw, key, value)
	if checker.NonNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) replaceBody(body *vo.Body, key, value string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.replaceBodyText(body, key, value)
	} else if body.ContentType().IsJSON() {
		return s.replaceBodyJson(body, key, value)
	} else {
		return body, mapper.NewErrModifierIncompatibleBodyType("replace-body", body.ContentType().String())
	}
}

func (s modifierService) replaceBodyText(body *vo.Body, key, value string) (*vo.Body, error) {
	bodyStr, err := body.String()
	if checker.NonNil(err) {
		return body, err
	}

	modifiedBodyText := strings.ReplaceAll(bodyStr, key, value)
	return s.newBodyByString(body, modifiedBodyText)
}

func (s modifierService) replaceBodyJson(body *vo.Body, key, value string) (*vo.Body, error) {
	bodyRaw, err := body.Raw()
	if checker.NonNil(err) {
		return body, err
	}

	if s.jsonPath.Parse(bodyRaw).Get(key).NotExists() {
		return body, nil
	}

	modifiedBodyJson, err := s.jsonPath.Set(bodyRaw, key, value)
	if checker.NonNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) deleteBody(body *vo.Body, key string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.deleteBodyText(body, key)
	} else if body.ContentType().IsJSON() {
		return s.deleteBodyJson(body, key)
	}
	return body, mapper.NewErrModifierIncompatibleBodyType("delete-body", body.ContentType().String())
}

func (s modifierService) deleteBodyText(body *vo.Body, key string) (*vo.Body, error) {
	bodyStr, err := body.String()
	if checker.NonNil(err) {
		return body, err
	}

	modifiedBodyText := strings.ReplaceAll(bodyStr, key, "")
	return s.newBodyByString(body, modifiedBodyText)
}

func (s modifierService) deleteBodyJson(body *vo.Body, key string) (*vo.Body, error) {
	bodyRaw, err := body.Raw()
	if checker.NonNil(err) {
		return body, err
	}

	modifiedBodyJson, err := s.jsonPath.Delete(bodyRaw, key)
	if checker.NonNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) newBodyByString(body *vo.Body, modifiedBodyJson string) (*vo.Body, error) {
	buffer, err := converter.ToBufferWithErr(modifiedBodyJson)
	if checker.NonNil(err) {
		return body, err
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}

func (s modifierService) wrapModifierErr(kind, op string, modifier *vo.Modifier, err error, format string, args ...any) error {
	if checker.IsNil(err) {
		return nil
	}

	base := fmt.Sprintf(
		"modifier failed: op=%s kind=%s action=%s key=%s",
		op,
		kind,
		modifier.Action(),
		modifier.Key(),
	)

	if checker.IsNotEmpty(format) {
		base = fmt.Sprintf("%s %s", base, fmt.Sprintf(format, args...))
	}

	return errors.Inheritf(err, base)
}

func (s modifierService) joinDynamicValueErr(kind string, modifier *vo.Modifier, errs []error) error {
	if checker.IsEmpty(errs) {
		return nil
	}

	return errors.JoinInheritf(errs, ", ", "modifier failed: op=resolve-dynamic-value kind=%s action=%s key=%s value=%s",
		kind, modifier.Action(), modifier.Key(), modifier.Value())
}
