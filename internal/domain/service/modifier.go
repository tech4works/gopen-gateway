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
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type modifier struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Modifier interface {
	ExecuteURLPathModifiers(
		configs []vo.ModifierConfig,
		urlPath vo.URLPath,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (vo.URLPath, []error)
	ExecuteMetadataModifiers(
		configs []vo.ModifierConfig,
		metadata vo.Metadata,
		ignoreKeys []string,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (vo.Metadata, []error)
	ExecuteQueryModifiers(
		configs []vo.ModifierConfig,
		query vo.Query,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (vo.Query, []error)
	ExecutePayloadModifiers(
		configs []vo.ModifierConfig,
		payload *vo.Payload,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (*vo.Payload, []error)

	ModifyURLPath(config *vo.ModifierConfig, urlPath vo.URLPath, request *vo.EndpointRequest, history *aggregate.History,
	) (vo.URLPath, error)
	ModifyMetadata(
		config *vo.ModifierConfig,
		metadata vo.Metadata,
		ignoreKeys []string,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (vo.Metadata, error)
	ModifyQuery(config *vo.ModifierConfig, query vo.Query, request *vo.EndpointRequest, history *aggregate.History) (
		vo.Query, error)
	ModifyPayload(config *vo.ModifierConfig, payload *vo.Payload, request *vo.EndpointRequest, history *aggregate.History,
	) (*vo.Payload, error)
}

func NewModifier(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Modifier {
	return modifier{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (s modifier) ExecutePayloadModifiers(
	configs []vo.ModifierConfig,
	payload *vo.Payload,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.Payload, []error) {
	var allErrs []error
	for _, config := range configs {
		var err error

		payload, err = s.ModifyPayload(&config, payload, request, history)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "modifier failed: op=execute kind=payload action=%s key=%s value=%s",
				config.Action(), config.Key(), config.Value()))
		}
	}
	return payload, allErrs
}

func (s modifier) ExecuteURLPathModifiers(
	configs []vo.ModifierConfig,
	urlPath vo.URLPath,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.URLPath, []error) {
	var allErrs []error

	for _, config := range configs {
		var err error

		urlPath, err = s.ModifyURLPath(&config, urlPath, request, history)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "modifier failed: op=execute kind=url-path action=%s key=%s value=%s",
				config.Action(), config.Key(), config.Value()))
		}
	}

	return urlPath, allErrs
}

func (s modifier) ExecuteMetadataModifiers(
	configs []vo.ModifierConfig,
	metadata vo.Metadata,
	ignoreKeys []string,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Metadata, []error) {
	var allErrs []error

	for _, config := range configs {
		var err error

		metadata, err = s.ModifyMetadata(&config, metadata, ignoreKeys, request, history)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "modifier failed: op=execute kind=metadata action=%s key=%s value=%s",
				config.Action(), config.Key(), config.Value()))
		}
	}

	return metadata, allErrs
}

func (s modifier) ExecuteQueryModifiers(configs []vo.ModifierConfig, query vo.Query, request *vo.EndpointRequest,
	history *aggregate.History) (vo.Query, []error) {
	var allErrs []error

	for _, config := range configs {
		var err error

		query, err = s.ModifyQuery(&config, query, request, history)
		if checker.NonNil(err) {
			allErrs = append(allErrs, errors.Inheritf(err, "modifier failed: op=execute kind=query action=%s key=%s value=%s",
				config.Action(), config.Key(), config.Value()))
		}
	}

	return query, allErrs
}

func (s modifier) ModifyURLPath(config *vo.ModifierConfig, urlPath vo.URLPath, request *vo.EndpointRequest,
	history *aggregate.History) (vo.URLPath, error) {
	shouldRun, err := s.evalModifierGuards(config, "url-path", request, history)
	if checker.NonNil(err) {
		return urlPath, s.wrapModifierErr(config, "url-path", "eval-guards", err, "")
	} else if !shouldRun {
		return urlPath, nil
	}

	action := config.Action()
	key := config.Key()
	value, errs := s.dynamicValueService.Get(config.Value(), request, history)
	if checker.IsNotEmpty(errs) {
		return urlPath, s.joinDynamicValueErr("url-path", config, errs)
	}

	switch action {
	case enum.ModifierActionSet:
		return s.setURLPath(urlPath, key, value)
	case enum.ModifierActionReplace:
		return s.replaceURLPath(urlPath, key, value)
	case enum.ModifierActionDelete:
		return s.deleteURLPath(urlPath, key)
	default:
		return urlPath, domain.NewErrModifierActionNotImplemented("url-path", action)
	}
}

func (s modifier) ModifyMetadata(
	config *vo.ModifierConfig,
	metadata vo.Metadata,
	ignoreKeys []string,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Metadata, error) {
	if checker.Contains(ignoreKeys, config.Key()) {
		return metadata, nil
	}

	shouldRun, err := s.evalModifierGuards(config, "metadata", request, history)
	if checker.NonNil(err) {
		return metadata, s.wrapModifierErr(config, "metadata", "eval-guards", err, "")
	} else if !shouldRun {
		return metadata, nil
	}

	action := config.Action()
	key := config.Key()
	values, errs := s.dynamicValueService.GetAsSliceOfString(config.Value(), request, history)
	if checker.IsNotEmpty(errs) {
		return metadata, s.joinDynamicValueErr("metadata", config, errs)
	}

	switch action {
	case enum.ModifierActionAdd:
		return s.addMetadata(metadata, key, values)
	case enum.ModifierActionAppend:
		return s.appendMetadata(metadata, key, values)
	case enum.ModifierActionSet:
		return s.setMetadata(metadata, key, values)
	case enum.ModifierActionReplace:
		return s.replaceMetadata(metadata, key, values)
	case enum.ModifierActionDelete:
		return s.deleteMetadata(metadata, key)
	default:
		return metadata, domain.NewErrModifierActionNotImplemented("metadata", action)
	}
}

func (s modifier) ModifyQuery(
	config *vo.ModifierConfig,
	query vo.Query,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Query, error) {
	shouldRun, err := s.evalModifierGuards(config, "query", request, history)
	if checker.NonNil(err) {
		return query, s.wrapModifierErr(config, "query", "eval-guards", err, "")
	} else if !shouldRun {
		return query, nil
	}

	action := config.Action()
	key := config.Key()
	values, errs := s.dynamicValueService.GetAsSliceOfString(config.Value(), request, history)
	if checker.IsNotEmpty(errs) {
		return query, s.joinDynamicValueErr("query", config, errs)
	}

	switch action {
	case enum.ModifierActionAdd:
		return s.addQuery(query, key, values)
	case enum.ModifierActionAppend:
		return s.appendQuery(query, key, values)
	case enum.ModifierActionSet:
		return s.setQuery(query, key, values)
	case enum.ModifierActionReplace:
		return s.replaceQuery(query, key, values)
	case enum.ModifierActionDelete:
		return s.deleteQuery(query, key)
	default:
		return query, domain.NewErrModifierActionNotImplemented("query", action)
	}
}

func (s modifier) ModifyPayload(
	config *vo.ModifierConfig,
	payload *vo.Payload,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.Payload, error) {
	if checker.IsNil(payload) {
		return nil, nil
	}

	shouldRun, err := s.evalModifierGuards(config, "payload", request, history)
	if checker.NonNil(err) {
		return payload, s.wrapModifierErr(config, "payload", "eval-guards", err, "")
	} else if !shouldRun {
		return payload, nil
	}

	action := config.Action()
	key := config.Key()
	value, errs := s.dynamicValueService.Get(config.Value(), request, history)
	if checker.IsNotEmpty(errs) {
		return payload, s.joinDynamicValueErr("payload", config, errs)
	}

	switch action {
	case enum.ModifierActionAdd:
		return s.addPayload(payload, key, value)
	case enum.ModifierActionAppend:
		return s.appendPayload(payload, key, value)
	case enum.ModifierActionSet:
		return s.setPayload(payload, key, value)
	case enum.ModifierActionReplace:
		return s.replacePayload(payload, key, value)
	case enum.ModifierActionDelete:
		return s.deletePayload(payload, key)
	default:
		return payload, domain.NewErrModifierActionNotImplemented("payload", action)
	}
}

func (s modifier) evalModifierGuards(
	config *vo.ModifierConfig,
	kind string,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (bool, error) {
	shouldRun, _, errs := s.dynamicValueService.EvalGuards(config.OnlyIf(), config.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInheritf(errs, ", ", "modifier failed: op=eval-guards kind=%s action=%s key=%s",
			kind, config.Action(), config.Key())
	}
	return shouldRun, nil
}

func (s modifier) setURLPath(urlPath vo.URLPath, key, value string) (vo.URLPath, error) {
	path := urlPath.Raw()
	paramValues := urlPath.Params().Copy()

	paramValues[key] = value
	if checker.NotContains(path, fmt.Sprintf("/:%s", key)) {
		path = fmt.Sprintf("%s/:%s", path, key)
	}

	return vo.NewURLPath(path, paramValues), nil
}

func (s modifier) replaceURLPath(urlPath vo.URLPath, key, value string) (vo.URLPath, error) {
	if urlPath.NotExists(key) {
		return urlPath, nil
	}

	return s.setURLPath(urlPath, key, value)
}

func (s modifier) deleteURLPath(urlPath vo.URLPath, key string) (vo.URLPath, error) {
	path := strings.ReplaceAll(urlPath.Raw(), fmt.Sprintf("/:%s", key), "")

	paramValues := urlPath.Params().Copy()
	delete(paramValues, key)

	return vo.NewURLPath(path, paramValues), nil
}

func (s modifier) addMetadata(metadata vo.Metadata, key string, value []string) (vo.Metadata, error) {
	values := metadata.Copy()
	values[key] = append(metadata.GetAll(key), value...)

	return vo.NewMetadata(values), nil
}

func (s modifier) appendMetadata(metadata vo.Metadata, key string, value []string) (vo.Metadata, error) {
	if metadata.NotExists(key) {
		return metadata, nil
	}

	values := metadata.Copy()
	values[key] = append(metadata.GetAll(key), value...)

	return vo.NewMetadata(values), nil
}

func (s modifier) setMetadata(metadata vo.Metadata, key string, value []string) (vo.Metadata, error) {
	values := metadata.Copy()
	values[key] = value

	return vo.NewMetadata(values), nil
}

func (s modifier) replaceMetadata(metadata vo.Metadata, key string, value []string) (vo.Metadata, error) {
	if metadata.NotExists(key) {
		return metadata, nil
	}
	return s.setMetadata(metadata, key, value)
}

func (s modifier) deleteMetadata(metadata vo.Metadata, key string) (vo.Metadata, error) {
	values := metadata.Copy()
	delete(values, key)

	return vo.NewMetadata(values), nil
}

func (s modifier) addQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	values := query.Copy()
	values[key] = append(query.GetAll(key), value...)

	return vo.NewQuery(values), nil
}

func (s modifier) appendQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	if query.NotExists(key) {
		return query, nil
	}

	values := query.Copy()
	values[key] = append(query.GetAll(key), value...)

	return vo.NewQuery(values), nil
}

func (s modifier) setQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	values := query.Copy()
	values[key] = value

	return vo.NewQuery(values), nil
}

func (s modifier) replaceQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	if query.NotExists(key) {
		return query, nil
	}

	return s.setQuery(query, key, value)
}

func (s modifier) deleteQuery(query vo.Query, key string) (vo.Query, error) {
	values := query.Copy()

	delete(values, key)

	return vo.NewQuery(values), nil
}

func (s modifier) addPayload(payload *vo.Payload, key, value string) (*vo.Payload, error) {
	if payload.ContentType().IsPlainText() {
		return s.addPayloadPlainText(payload, value)
	} else if payload.ContentType().IsJSON() {
		return s.addPayloadJSON(payload, key, value)
	} else {
		return payload, domain.NewErrModifierIncompatibleContentType("add-payload", payload.ContentType().String())
	}
}

func (s modifier) addPayloadPlainText(payload *vo.Payload, value string) (*vo.Payload, error) {
	str, err := payload.String()
	if checker.NonNil(err) {
		return payload, err
	}

	modifiedPayloadText := fmt.Sprintf("%s%s", str, value)
	return s.newPayloadByString(payload, modifiedPayloadText)
}

func (s modifier) addPayloadJSON(payload *vo.Payload, key, value string) (*vo.Payload, error) {
	payloadRaw, err := payload.Raw()
	if checker.NonNil(err) {
		return payload, err
	}

	modifiedPayloadJSON, err := s.jsonPath.Add(payloadRaw, key, value)
	if checker.NonNil(err) {
		return payload, err
	}

	return s.newPayloadByString(payload, modifiedPayloadJSON)
}

func (s modifier) appendPayload(payload *vo.Payload, key, value string) (*vo.Payload, error) {
	if payload.ContentType().IsPlainText() {
		return s.appendPayloadPlainText(payload, value)
	} else if payload.ContentType().IsJSON() {
		return s.appendPayloadJSON(payload, key, value)
	} else {
		return payload, domain.NewErrModifierIncompatibleContentType("append-payload", payload.ContentType().String())
	}
}

func (s modifier) appendPayloadPlainText(payload *vo.Payload, value string) (*vo.Payload, error) {
	str, err := payload.String()
	if checker.NonNil(err) {
		return payload, err
	}

	modifiedPayloadText := fmt.Sprintf("%s\n%s", str, value)
	return s.newPayloadByString(payload, modifiedPayloadText)
}

func (s modifier) appendPayloadJSON(payload *vo.Payload, key, value string) (*vo.Payload, error) {
	raw, err := payload.Raw()
	if checker.NonNil(err) {
		return payload, err
	}

	if s.jsonPath.Parse(raw).Get(key).NotExists() {
		return payload, nil
	}

	modifiedPayloadJSON, err := s.jsonPath.Add(raw, key, value)
	if checker.NonNil(err) {
		return payload, err
	}

	return s.newPayloadByString(payload, modifiedPayloadJSON)
}

func (s modifier) setPayload(payload *vo.Payload, key, value string) (*vo.Payload, error) {
	if payload.ContentType().IsPlainText() {
		return s.setPayloadPlainText(payload, value)
	} else if payload.ContentType().IsJSON() {
		return s.setPayloadJSON(payload, key, value)
	} else {
		return payload, domain.NewErrModifierIncompatibleContentType("set-payload", payload.ContentType().String())
	}
}

func (s modifier) setPayloadPlainText(payload *vo.Payload, value string) (*vo.Payload, error) {
	return s.newPayloadByString(payload, value)
}

func (s modifier) setPayloadJSON(payload *vo.Payload, key, value string) (*vo.Payload, error) {
	raw, err := payload.Raw()
	if checker.NonNil(err) {
		return payload, err
	}

	modifiedPayloadJSON, err := s.jsonPath.Set(raw, key, value)
	if checker.NonNil(err) {
		return payload, err
	}

	return s.newPayloadByString(payload, modifiedPayloadJSON)
}

func (s modifier) replacePayload(payload *vo.Payload, key, value string) (*vo.Payload, error) {
	if payload.ContentType().IsPlainText() {
		return s.replacePayloadPlainText(payload, key, value)
	} else if payload.ContentType().IsJSON() {
		return s.replacePayloadJSON(payload, key, value)
	} else {
		return payload, domain.NewErrModifierIncompatibleContentType("replace-payload", payload.ContentType().String())
	}
}

func (s modifier) replacePayloadPlainText(payload *vo.Payload, key, value string) (*vo.Payload, error) {
	str, err := payload.String()
	if checker.NonNil(err) {
		return payload, err
	}

	modifiedPayloadText := strings.ReplaceAll(str, key, value)
	return s.newPayloadByString(payload, modifiedPayloadText)
}

func (s modifier) replacePayloadJSON(payload *vo.Payload, key, value string) (*vo.Payload, error) {
	raw, err := payload.Raw()
	if checker.NonNil(err) {
		return payload, err
	}

	if s.jsonPath.Parse(raw).Get(key).NotExists() {
		return payload, nil
	}

	modifiedPayloadJSON, err := s.jsonPath.Set(raw, key, value)
	if checker.NonNil(err) {
		return payload, err
	}

	return s.newPayloadByString(payload, modifiedPayloadJSON)
}

func (s modifier) deletePayload(payload *vo.Payload, key string) (*vo.Payload, error) {
	if payload.ContentType().IsPlainText() {
		return s.deletePayloadPlainText(payload, key)
	} else if payload.ContentType().IsJSON() {
		return s.deletePayloadJSON(payload, key)
	} else {
		return payload, domain.NewErrModifierIncompatibleContentType("delete-payload", payload.ContentType().String())
	}
}

func (s modifier) deletePayloadPlainText(payload *vo.Payload, key string) (*vo.Payload, error) {
	str, err := payload.String()
	if checker.NonNil(err) {
		return payload, err
	}

	modifiedPayloadText := strings.ReplaceAll(str, key, "")
	return s.newPayloadByString(payload, modifiedPayloadText)
}

func (s modifier) deletePayloadJSON(payload *vo.Payload, key string) (*vo.Payload, error) {
	raw, err := payload.Raw()
	if checker.NonNil(err) {
		return payload, err
	}

	modifiedPayloadJSON, err := s.jsonPath.Delete(raw, key)
	if checker.NonNil(err) {
		return payload, err
	}

	return s.newPayloadByString(payload, modifiedPayloadJSON)
}

func (s modifier) newPayloadByString(payload *vo.Payload, raw string) (*vo.Payload, error) {
	buffer, err := converter.ToBufferWithErr(raw)
	if checker.NonNil(err) {
		return payload, err
	}

	return vo.NewPayloadWithContentType(payload.ContentType(), buffer), nil
}

func (s modifier) wrapModifierErr(config *vo.ModifierConfig, kind, op string, err error, format string, args ...any) error {
	if checker.IsNil(err) {
		return nil
	}

	base := fmt.Sprintf(
		"modifier failed: op=%s kind=%s action=%s key=%s",
		op,
		kind,
		config.Action(),
		config.Key(),
	)

	if checker.IsNotEmpty(format) {
		base = fmt.Sprintf("%s %s", base, fmt.Sprintf(format, args...))
	}

	return errors.Inheritf(err, base)
}

func (s modifier) joinDynamicValueErr(kind string, modifier *vo.ModifierConfig, errs []error) error {
	if checker.IsEmpty(errs) {
		return nil
	}

	return errors.JoinInheritf(errs, ", ",
		"modifier failed: op=resolve-dynamic-value kind=%s action=%s key=%s value=%s",
		kind, modifier.Action(), modifier.Key(), modifier.Value())
}
