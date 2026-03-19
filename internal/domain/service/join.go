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

type join struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Join interface {
	ExecutePayloadJoins(
		configs []vo.JoinConfig,
		payload *vo.Payload,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (*vo.Payload, []error)
	JoinPayload(
		config *vo.JoinConfig,
		payload *vo.Payload,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (*vo.Payload, []error)
}

func NewJoin(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Join {
	return join{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (s join) ExecutePayloadJoins(
	configs []vo.JoinConfig,
	payload *vo.Payload,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.Payload, []error) {
	if checker.IsNil(payload) {
		return nil, nil
	}

	var allErrs []error

	for _, config := range configs {
		newPayload, errs := s.JoinPayload(&config, payload, request, history)
		if checker.IsNotEmpty(errs) {
			allErrs = append(allErrs, errs...)
		}
		payload = newPayload
	}

	return payload, allErrs
}

func (s join) JoinPayload(
	config *vo.JoinConfig,
	payload *vo.Payload,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.Payload, []error) {
	if checker.IsNil(payload) || payload.ContentType().IsNotJSON() {
		return payload, nil
	}

	shouldRun, err := s.evalJoinGuards(config, request, history)
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlicef(err, "join failed: op=eval-guards source.path=%s target.path=%s",
			config.Source().Path(), config.Target().Path())
	} else if !shouldRun {
		return payload, nil
	}

	source, err := s.getSourceByPath(config, request, history)
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlicef(err, "join failed: op=get-source source.path=%s",
			config.Source().Path())
	}

	target, err := s.getTargetByPath(config, payload)
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlicef(err, "join failed: op=get-target target.path=%s",
			config.Target().Path())
	}

	if err = s.assertPathsAndTypes(config, source, target); checker.NonNil(err) {
		return payload, errors.InheritAsSlicef(err, "join failed: op=assert-paths-types source.path=%s target.path=%s",
			config.Source().Path(), config.Target().Path())
	}

	var allErrs []error

	newPayload, errs := s.joinByShape(config, payload, source, target)
	if checker.IsNotEmpty(errs) {
		for _, e := range errs {
			allErrs = append(allErrs, errors.Inheritf(e,
				"join failed: op=join-by-shape source.path=%s source.key=%s target.path=%s target.key=%s target.as=%s",
				config.Source().Path(), config.Source().Key(), config.Target().Path(), config.Target().Key(),
				config.Target().As()))
		}
	}

	newPayload, errs = s.applyTargetKeyPolicy(config.Target(), newPayload)
	if checker.NonNil(err) {
		for _, e := range errs {
			allErrs = append(allErrs, errors.Inheritf(e,
				"join failed: op=apply-target-key-policy source.path=%s source.key=%s target.path=%s target.key=%s target.as=%s",
				config.Source().Path(), config.Source().Key(), config.Target().Path(), config.Target().Key(),
				config.Target().As()))
		}
	}

	return newPayload, allErrs
}

func (s join) joinByShape(
	config *vo.JoinConfig,
	payload *vo.Payload,
	source,
	target domain.JSONValue,
) (*vo.Payload, []error) {
	switch {
	case source.Exists() && source.IsArray() && target.Exists() && target.IsArray():
		return s.joinArrayToArray(config, payload, source, target)
	case source.Exists() && source.IsObject() && target.Exists() && target.IsArray():
		return s.joinObjectToArray(config, payload, source, target)
	case source.Exists() && source.IsArray() && target.Exists() && target.IsObject():
		return s.joinArrayToObject(config, payload, source, target)
	case source.Exists() && source.IsObject() && target.Exists() && target.IsObject():
		return s.joinObjectToObject(config, payload, source, target)
	default:
		return payload, nil
	}
}

func (s join) joinArrayToArray(
	config *vo.JoinConfig,
	payload *vo.Payload,
	source,
	target domain.JSONValue,
) (*vo.Payload, []error) {
	index := s.buildIndexFromArray(source, config.Source().Key())
	out := target.Raw()

	var errs []error
	var i int
	target.ForEach(func(_ string, item domain.JSONValue) bool {
		var err error

		setPath := s.targetJSONPathOnArray(config.Target(), i)
		targetKey := item.Get(config.Target().Key())

		out, err = s.applyMatchOnArrayItem(config, out, setPath, i, targetKey, index)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err,
				"join failed: op=apply-match array idx=%d setPath=%s target.key=%s", i, setPath,
				config.Target().Key()))
		}

		i++
		return true
	})

	newPayload, err := s.newPayloadByString(payload, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=buffer array out"))
	}

	return newPayload, errs
}

func (s join) joinObjectToArray(
	config *vo.JoinConfig,
	payload *vo.Payload,
	source,
	target domain.JSONValue,
) (*vo.Payload, []error) {
	sourceKey := strings.TrimSpace(source.Get(config.Source().Key()).String())
	if checker.IsEmpty(sourceKey) {
		if config.Target().IsOnMissingError() {
			return payload, errors.NewAsSlicef("join failed: source key missing/empty: source.path=%s; source.key=%s",
				config.Source().Path(), config.Source().Key())
		}
		return payload, nil
	}

	out := target.Raw()

	var errs []error
	var i int
	target.ForEach(func(_ string, item domain.JSONValue) bool {
		var err error

		setPath := s.targetJSONPathOnArray(config.Target(), i)
		targetKey := item.Get(config.Target().Key())

		out, err = s.applyObjectMatchOnArrayItem(config, out, setPath, i, targetKey, sourceKey, source.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "join failed: op=apply-match array idx=%d setPath=%s target.key=%s",
				i, setPath, config.Target().Key()))
		}

		i++
		return true
	})

	newPayload, err := s.newPayloadByString(payload, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=buffer array out"))
	}

	return newPayload, errs
}

func (s join) joinArrayToObject(
	config *vo.JoinConfig,
	payload *vo.Payload,
	source,
	target domain.JSONValue,
) (*vo.Payload, []error) {
	index := s.buildIndexFromArray(source, config.Source().Key())
	setPath := s.targetJSONPathOnObject(config.Target())
	tgtKey := target.Get(config.Target().Key())

	if tgtKey.NotExists() {
		return s.onMissingAndRebuildPayload(config, payload, target.Raw(), setPath, -1, "", "target-key-not-found")
	} else if tgtKey.IsArray() {
		return s.joinObjectTargetFromIDArray(config, payload, target.Raw(), setPath, tgtKey, index)
	} else {
		return s.joinObjectTargetFromSingleID(config, payload, target.Raw(), setPath, tgtKey, index)
	}
}

func (s join) joinObjectTargetFromIDArray(
	config *vo.JoinConfig,
	payload *vo.Payload,
	raw string,
	setPath string,
	tgtKey domain.JSONValue,
	index map[string]string,
) (*vo.Payload, []error) {
	jsonArr, errs := s.mapTargetIDsToJSONArray(config, -1, tgtKey, index)

	out, err := s.jsonPath.Set(raw, setPath, jsonArr)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=set setPath=%s mode=array target.path=%s target.as=%s",
			setPath, config.Target().Path(), config.Target().As()))
	}

	newPayload, err := s.newPayloadByString(payload, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-payload array->object"))
	}

	return newPayload, errs
}

func (s join) joinObjectTargetFromSingleID(
	config *vo.JoinConfig,
	payload *vo.Payload,
	raw,
	setPath string,
	tgtKey domain.JSONValue,
	index map[string]string,
) (*vo.Payload, []error) {
	id := strings.TrimSpace(tgtKey.String())
	if checker.IsEmpty(id) {
		return s.onMissingAndRebuildPayload(config, payload, raw, setPath, -1, id, "target-key-empty")
	}

	srcRaw, ok := index[id]
	if !ok {
		return s.onMissingAndRebuildPayload(config, payload, raw, setPath, -1, id, "source-not-found")
	}

	var errs []error

	out, err := s.jsonPath.Set(raw, setPath, srcRaw)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=set setPath=%s mode=single id=%s target.path=%s target.as=%s",
			setPath, id, config.Target().Path(), config.Target().As()))
	}

	newPayload, err := s.newPayloadByString(payload, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-payload array->object"))
	}

	return newPayload, errs
}

func (s join) joinObjectToObject(
	config *vo.JoinConfig,
	payload *vo.Payload,
	source,
	target domain.JSONValue,
) (*vo.Payload, []error) {
	setPath := s.targetJSONPathOnObject(config.Target())

	srcID := strings.TrimSpace(source.Get(config.Source().Key()).String())
	tgtID := strings.TrimSpace(target.Get(config.Target().Key()).String())

	raw := target.Raw()

	if checker.IsEmpty(srcID) || checker.IsEmpty(tgtID) {
		return s.onMissingAndRebuildPayload(config, payload, raw, setPath, -1, tgtID, "source-or-target-key-empty")
	} else if checker.NotEquals(srcID, tgtID) {
		return s.onMissingAndRebuildPayload(config, payload, raw, setPath, -1, tgtID, "target-not-equals-source")
	}

	var errs []error

	out, err := s.jsonPath.Set(raw, setPath, source.Raw())
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err,
			"join failed: op=set reason=target-equals-source setPath=%s id=%s target.as=%s",
			setPath, tgtID, config.Target().As()))
	}

	newPayload, err := s.newPayloadByString(payload, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-payload object->object"))
	}

	return newPayload, errs
}

func (s join) onMissingAndRebuildPayload(
	config *vo.JoinConfig,
	payload *vo.Payload,
	raw,
	setPath string,
	targetIndex int,
	targetKeyValue,
	reason string,
) (*vo.Payload, []error) {
	var errs []error

	out, err := s.applyOnMissingSingle(config, raw, setPath, targetIndex, targetKeyValue, reason)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=apply-missing-single"))
	}

	newPayload, err := s.newPayloadByString(payload, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-payload"))
	}

	return newPayload, errs
}

func (s join) applyMatchOnArrayItem(
	config *vo.JoinConfig,
	raw string,
	setPath string,
	targetIndex int,
	targetKey domain.JSONValue,
	index map[string]string,
) (string, error) {
	if targetKey.NotExists() {
		return s.applyOnMissingSingle(config, raw, setPath, targetIndex, "", "target-key-not-found")
	} else if targetKey.IsArray() {
		jsonArr, errs := s.mapTargetIDsToJSONArray(config, targetIndex, targetKey, index)

		out, err := s.jsonPath.Set(raw, setPath, jsonArr)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "join failed: op=set-json-path"))
		}

		return out, errors.JoinInheritf(errs, ", ", "join failed: op=set-target-array")
	}

	id := strings.TrimSpace(targetKey.String())
	if checker.IsEmpty(id) {
		return s.applyOnMissingSingle(config, raw, setPath, targetIndex, id, "target-key-empty")
	}

	if srcRaw, ok := index[id]; ok {
		out, err := s.jsonPath.Set(raw, setPath, srcRaw)
		if checker.NonNil(err) {
			return raw, err
		}
		return out, nil
	}

	return s.applyOnMissingSingle(config, raw, setPath, targetIndex, id, "source-not-found")
}

func (s join) applyObjectMatchOnArrayItem(
	config *vo.JoinConfig,
	raw string,
	setPath string,
	targetIndex int,
	targetKey domain.JSONValue,
	sourceKey,
	sourceRaw string,
) (string, error) {
	if targetKey.NotExists() {
		return s.applyOnMissingSingle(config, raw, setPath, targetIndex, "", "target-key-not-found")
	}

	if targetKey.IsArray() {
		found := false
		targetKey.ForEach(func(_ string, v domain.JSONValue) bool {
			if checker.Equals(strings.TrimSpace(v.String()), sourceKey) {
				found = true
				return false
			}
			return true
		})

		if found {
			out, err := s.jsonPath.Set(raw, setPath, sourceRaw)
			if checker.NonNil(err) {
				return raw, err
			}
			return out, nil
		}

		return s.applyOnMissingSingle(config, raw, setPath, targetIndex, sourceKey, "target-array-not-contains-source")
	}

	id := strings.TrimSpace(targetKey.String())
	if checker.IsEmpty(id) {
		return s.applyOnMissingSingle(config, raw, setPath, targetIndex, id, "target-key-empty")
	}

	if checker.Equals(id, sourceKey) {
		out, err := s.jsonPath.Set(raw, setPath, sourceRaw)
		if checker.NonNil(err) {
			return raw, err
		}
		return out, nil
	}

	return s.applyOnMissingSingle(config, raw, setPath, targetIndex, id, "target-not-equals-source")
}

func (s join) applyTargetKeyPolicy(targetConfig vo.JoinConfigTarget, payload *vo.Payload) (*vo.Payload, []error) {
	if targetConfig.KeepKey() {
		return payload, nil
	}

	raw, err := payload.Raw()
	if checker.NonNil(err) {
		return payload, errors.InheritAsSlicef(err, "join failed: op=raw path=%s key=%s", targetConfig.Path(),
			targetConfig.Key())
	}

	target := s.jsonPath.Get(raw, targetConfig.Path())
	if target.NotExists() {
		return payload, nil
	}

	out := raw

	var errs []error
	switch targetConfig.Policy() {
	case enum.JoinTargetDropKeyAlways:
		out, errs = s.removeTargetKeyAll(targetConfig, raw, target)
	case enum.JoinTargetDropKeyOnMerged:
		out, errs = s.removeTargetKeyOnlyMerged(targetConfig, raw, target)
	default:
		return payload, nil
	}

	newPayload, err := s.newPayloadByString(payload, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-payload"))
	}

	return newPayload, errs
}

func (s join) getTargetByPath(config *vo.JoinConfig, payload *vo.Payload) (domain.JSONValue, error) {
	targetRaw, err := payload.Raw()
	if checker.NonNil(err) {
		return nil, errors.Inheritf(err, "join failed: op=raw")
	}

	return s.jsonPath.Get(targetRaw, config.Target().Path()), nil
}

func (s join) getSourceByPath(
	config *vo.JoinConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (domain.JSONValue, error) {
	sourceValue, sourceErrs := s.dynamicValueService.Get(config.Source().Path(), request, history)
	if checker.IsNotEmpty(sourceErrs) {
		return nil, errors.JoinInheritf(sourceErrs, ", ", "join failed: op=resolve-source path=%s key=%s",
			config.Source().Path(), config.Source().Key())
	}
	return s.jsonPath.Parse(strings.TrimSpace(sourceValue)), nil
}

func (s join) mapTargetIDsToJSONArray(
	config *vo.JoinConfig,
	targetIndex int,
	targetIDs domain.JSONValue,
	index map[string]string,
) (string, []error) {
	var out []string
	var errs []error

	targetIDs.ForEach(func(_ string, v domain.JSONValue) bool {
		id := strings.TrimSpace(v.String())

		if checker.IsEmpty(id) {
			appendIt, value, err := s.applyOnMissingCollection(config, targetIndex, id, "target-key-item-empty")
			if checker.NonNil(err) {
				errs = append(errs, errors.Inherit(err, "join failed: missing collection item (empty)"))
			} else if appendIt {
				out = append(out, value)
			}
			return true
		}

		if srcRaw, ok := index[id]; ok {
			out = append(out, srcRaw)
			return true
		}

		appendIt, value, err := s.applyOnMissingCollection(config, targetIndex, id, "source-not-found")
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "join failed: source not found: id=%s", id))
		} else if appendIt {
			out = append(out, value)
		}
		return true
	})

	return s.buildJSONArrayRaw(out), errs
}

func (s join) removeTargetKeyAll(
	joinTarget vo.JoinConfigTarget,
	raw string,
	target domain.JSONValue,
) (string, []error) {
	var errs []error
	if target.IsArray() {
		out := raw

		var i int
		target.ForEach(func(_ string, _ domain.JSONValue) bool {
			var err error

			path := s.targetKeyJSONPathOnArray(joinTarget, i)

			out, err = s.jsonPath.Delete(out, path)
			if checker.NonNil(err) {
				errs = append(errs, errors.Inheritf(err, "join failed: on array op=delete-target-key path=%s idx=%d",
					path, i))
			}

			i++
			return true
		})

		return out, errs
	}

	path := s.targetKeyJSONPathOnObject(joinTarget)
	out, err := s.jsonPath.Delete(raw, path)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: on object op=delete-target-key path=%s", path))
	}

	return out, errs
}

func (s join) removeTargetKeyOnlyMerged(
	joinTarget vo.JoinConfigTarget,
	raw string,
	target domain.JSONValue,
) (string, []error) {
	var errs []error
	if target.IsArray() {
		out := raw

		var i int
		target.ForEach(func(_ string, item domain.JSONValue) bool {
			if item.Get(joinTarget.As()).Exists() {
				var err error

				path := s.targetKeyJSONPathOnArray(joinTarget, i)

				out, err = s.jsonPath.Delete(out, path)
				if checker.NonNil(err) {
					errs = append(errs, errors.Inheritf(err, "join failed: on array op=delete-target-key path=%s idx=%d",
						path, i))
				}
			}

			i++
			return true
		})

		return out, errs
	} else if target.Get(joinTarget.As()).NotExists() {
		return raw, nil
	}

	path := s.targetKeyJSONPathOnObject(joinTarget)
	out, err := s.jsonPath.Delete(raw, path)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: on object op=delete-target-key path=%s", path))
	}
	return out, errs
}

func (s join) buildJSONArrayRaw(items []string) string {
	var b strings.Builder
	b.WriteByte('[')

	first := true
	for _, it := range items {
		it = strings.TrimSpace(it)
		if checker.IsEmpty(it) {
			continue
		}
		if !first {
			b.WriteByte(',')
		}
		first = false
		b.WriteString(it)
	}

	b.WriteByte(']')
	return b.String()
}

func (s join) targetJSONPathOnArray(target vo.JoinConfigTarget, i int) string {
	return fmt.Sprintf("%s.%d.%s", target.Path(), i, target.As())
}

func (s join) targetJSONPathOnObject(target vo.JoinConfigTarget) string {
	return fmt.Sprintf("%s.%s", target.Path(), target.As())
}

func (s join) targetKeyJSONPathOnArray(target vo.JoinConfigTarget, i int) string {
	return fmt.Sprintf("%s.%d.%s", target.Path(), i, target.Key())
}

func (s join) targetKeyJSONPathOnObject(target vo.JoinConfigTarget) string {
	return fmt.Sprintf("%s.%s", target.Path(), target.Key())
}

func (s join) assertPathsAndTypes(config *vo.JoinConfig, source, target domain.JSONValue) error {
	if !config.Target().IsOnMissingError() {
		return nil
	} else if !source.Exists() {
		return errors.Newf("join failed: source path not found: source.path=%s", config.Source().Path())
	} else if !target.Exists() {
		return errors.Newf("join failed: target path not found: target.path=%s", config.Target().Path())
	}

	if !(source.IsArray() || source.IsObject()) || !(target.IsArray() || target.IsObject()) {
		return errors.Newf(
			"join failed: source/target must be array or object: source.type=%s; target.type=%s; source.path=%s; target.path=%s",
			source.Type(),
			target.Type(),
			config.Source().Path(),
			config.Target().Path(),
		)
	}

	return nil
}

func (s join) evalJoinGuards(
	config *vo.JoinConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (bool, error) {
	shouldRun, _, errs := s.dynamicValueService.EvalGuards(config.OnlyIf(), config.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInherit(errs, ", ", "failed to evaluate guard for join")
	}
	return shouldRun, nil
}

func (s join) buildIndexFromArray(arr domain.JSONValue, key string) map[string]string {
	index := map[string]string{}
	if arr.NotExists() || !arr.IsArray() {
		return index
	}

	arr.ForEach(func(_ string, item domain.JSONValue) bool {
		k := item.Get(key)
		if k.NotExists() {
			return true
		}

		id := strings.TrimSpace(k.String())
		if checker.IsEmpty(id) {
			return true
		}

		if _, ok := index[id]; !ok {
			index[id] = item.Raw()
		}
		return true
	})

	return index
}

func (s join) applyOnMissingSingle(
	config *vo.JoinConfig,
	raw,
	setPath string,
	targetIndex int,
	targetKeyValue,
	reason string,
) (string, error) {
	switch config.Target().OnMissing() {
	case enum.JoinTargetOnMissingOmit:
		return raw, nil
	case enum.JoinTargetOnMissingError:
		return raw, errors.Newf(
			"join failed (missing data): reason=%s; target.path=%s; target.key=%s; target.index=%d; target.value=%s; source.path=%s; source.key=%s; target.as=%s",
			reason,
			config.Target().Path(),
			config.Target().Key(),
			targetIndex,
			targetKeyValue,
			config.Source().Path(),
			config.Source().Key(),
			config.Target().As(),
		)
	default:
		return s.jsonPath.Set(raw, setPath, "null")
	}
}

func (s join) applyOnMissingCollection(
	config *vo.JoinConfig,
	targetIndex int,
	targetKeyValue string,
	reason string,
) (bool, string, error) {
	switch config.Target().OnMissing() {
	case enum.JoinTargetOnMissingOmit:
		return false, "", nil
	case enum.JoinTargetOnMissingNull:
		return true, "null", nil
	case enum.JoinTargetOnMissingError:
		return false, "", errors.Newf(
			"join failed (missing data): reason=%s; target.path=%s; target.key=%s; target.index=%d; target.value=%s; source.path=%s; source.key=%s; target.as=%s",
			reason,
			config.Target().Path(),
			config.Target().Key(),
			targetIndex,
			targetKeyValue,
			config.Source().Path(),
			config.Source().Key(),
			config.Target().As(),
		)
	default:
		return true, "null", nil
	}
}

func (s join) newPayloadByString(payload *vo.Payload, modifiedPayload string) (*vo.Payload, error) {
	buffer, err := converter.ToBufferWithErr(modifiedPayload)
	if checker.NonNil(err) {
		return payload, err
	}

	return vo.NewPayloadWithContentType(payload.ContentType(), buffer), nil
}
