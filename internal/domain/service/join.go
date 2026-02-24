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

type joinService struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Join interface {
	ExecuteBodyJoins(joins []vo.Join, body *vo.Body, request *vo.HTTPRequest, history *aggregate.History) (*vo.Body, []error)
	JoinBody(join vo.Join, body *vo.Body, request *vo.HTTPRequest, history *aggregate.History) (*vo.Body, []error)
}

func NewJoin(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Join {
	return joinService{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (s joinService) ExecuteBodyJoins(
	joins []vo.Join,
	body *vo.Body,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (*vo.Body, []error) {
	if checker.IsNil(body) {
		return nil, nil
	}

	var allErrs []error

	for _, join := range joins {
		newBody, errs := s.JoinBody(join, body, request, history)
		if checker.IsNotEmpty(errs) {
			allErrs = append(allErrs, errs...)
		}
		body = newBody
	}

	return body, allErrs
}

func (s joinService) JoinBody(
	join vo.Join,
	body *vo.Body,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (*vo.Body, []error) {
	if checker.IsNil(body) || body.ContentType().IsNotJSON() {
		return body, nil
	}

	shouldRun, err := s.evalJoinGuards(join, request, history)
	if checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "join failed: op=eval-guards source.path=%s target.path=%s",
			join.Source().Path(), join.Target().Path())
	}
	if !shouldRun {
		return body, nil
	}

	source, err := s.getSourceByPath(join, request, history)
	if checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "join failed: op=get-source source.path=%s",
			join.Source().Path())
	}

	target, err := s.getTargetByPath(join, body)
	if checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "join failed: op=get-target target.path=%s",
			join.Target().Path())
	}

	if err = s.assertPathsAndTypes(join, source, target); checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "join failed: op=assert-paths-types source.path=%s target.path=%s",
			join.Source().Path(), join.Target().Path())
	}

	var allErrs []error

	newBody, errs := s.joinByShape(join, body, source, target)
	if checker.IsNotEmpty(errs) {
		for _, e := range errs {
			allErrs = append(allErrs, errors.Inheritf(e,
				"join failed: op=join-by-shape source.path=%s source.key=%s target.path=%s target.key=%s target.as=%s",
				join.Source().Path(), join.Source().Key(), join.Target().Path(), join.Target().Key(),
				join.Target().As()))
		}
	}

	newBody, errs = s.applyTargetKeyPolicy(join.Target(), newBody)
	if checker.NonNil(err) {
		for _, e := range errs {
			allErrs = append(allErrs, errors.Inheritf(e,
				"join failed: op=apply-target-key-policy source.path=%s source.key=%s target.path=%s target.key=%s target.as=%s",
				join.Source().Path(), join.Source().Key(), join.Target().Path(), join.Target().Key(),
				join.Target().As()))
		}
	}

	return newBody, allErrs
}

func (s joinService) joinByShape(
	join vo.Join,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	switch {
	case source.Exists() && source.IsArray() && target.Exists() && target.IsArray():
		return s.joinArrayToArray(join, body, source, target)
	case source.Exists() && source.IsObject() && target.Exists() && target.IsArray():
		return s.joinObjectToArray(join, body, source, target)
	case source.Exists() && source.IsArray() && target.Exists() && target.IsObject():
		return s.joinArrayToObject(join, body, source, target)
	case source.Exists() && source.IsObject() && target.Exists() && target.IsObject():
		return s.joinObjectToObject(join, body, source, target)
	default:
		return body, nil
	}
}

func (s joinService) joinArrayToArray(
	join vo.Join,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	index := s.buildIndexFromArray(source, join.Source().Key())
	out := target.Raw()

	var errs []error
	var i int
	target.ForEach(func(_ string, item domain.JSONValue) bool {
		var err error

		setPath := s.targetJSONPathOnArray(join.Target(), i)
		targetKey := item.Get(join.Target().Key())

		out, err = s.applyMatchOnArrayItem(join, out, setPath, i, targetKey, index)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err,
				"join failed: op=apply-match array idx=%d setPath=%s target.key=%s", i, setPath,
				join.Target().Key()))
		}

		i++
		return true
	})

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=buffer array out"))
	}

	return newBody, errs
}

func (s joinService) joinObjectToArray(
	join vo.Join,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	sourceKey := strings.TrimSpace(source.Get(join.Source().Key()).String())
	if checker.IsEmpty(sourceKey) {
		if join.Target().IsOnMissingError() {
			return body, errors.NewAsSlicef("join failed: source key missing/empty: source.path=%s; source.key=%s",
				join.Source().Path(), join.Source().Key())
		}
		return body, nil
	}

	out := target.Raw()

	var errs []error
	var i int
	target.ForEach(func(_ string, item domain.JSONValue) bool {
		var err error

		setPath := s.targetJSONPathOnArray(join.Target(), i)
		targetKey := item.Get(join.Target().Key())

		out, err = s.applyObjectMatchOnArrayItem(join, out, setPath, i, targetKey, sourceKey, source.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "join failed: op=apply-match array idx=%d setPath=%s target.key=%s",
				i, setPath, join.Target().Key()))
		}

		i++
		return true
	})

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=buffer array out"))
	}

	return newBody, errs
}

func (s joinService) joinArrayToObject(
	join vo.Join,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	index := s.buildIndexFromArray(source, join.Source().Key())
	setPath := s.targetJSONPathOnObject(join.Target())
	tgtKey := target.Get(join.Target().Key())

	if tgtKey.NotExists() {
		return s.onMissingAndRebuildBody(body, target.Raw(), setPath, join, -1, "", "target-key-not-found")
	} else if tgtKey.IsArray() {
		return s.joinObjectTargetFromIDArray(join, body, target.Raw(), setPath, tgtKey, index)
	} else {
		return s.joinObjectTargetFromSingleID(join, body, target.Raw(), setPath, tgtKey, index)
	}
}

func (s joinService) joinObjectTargetFromIDArray(
	join vo.Join,
	body *vo.Body,
	raw string,
	setPath string,
	tgtKey domain.JSONValue,
	index map[string]string,
) (*vo.Body, []error) {
	jsonArr, errs := s.mapTargetIDsToJSONArray(join, -1, tgtKey, index)

	out, err := s.jsonPath.Set(raw, setPath, jsonArr)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=set setPath=%s mode=array target.path=%s target.as=%s",
			setPath, join.Target().Path(), join.Target().As()))
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-body array->object"))
	}

	return newBody, errs
}

func (s joinService) joinObjectTargetFromSingleID(
	join vo.Join,
	body *vo.Body,
	raw,
	setPath string,
	tgtKey domain.JSONValue,
	index map[string]string,
) (*vo.Body, []error) {
	id := strings.TrimSpace(tgtKey.String())
	if checker.IsEmpty(id) {
		return s.onMissingAndRebuildBody(body, raw, setPath, join, -1, id, "target-key-empty")
	}

	srcRaw, ok := index[id]
	if !ok {
		return s.onMissingAndRebuildBody(body, raw, setPath, join, -1, id, "source-not-found")
	}

	var errs []error

	out, err := s.jsonPath.Set(raw, setPath, srcRaw)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=set setPath=%s mode=single id=%s target.path=%s target.as=%s",
			setPath, id, join.Target().Path(), join.Target().As()))
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-body array->object"))
	}

	return newBody, errs
}

func (s joinService) joinObjectToObject(
	join vo.Join,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	setPath := s.targetJSONPathOnObject(join.Target())

	srcID := strings.TrimSpace(source.Get(join.Source().Key()).String())
	tgtID := strings.TrimSpace(target.Get(join.Target().Key()).String())

	raw := target.Raw()

	if checker.IsEmpty(srcID) || checker.IsEmpty(tgtID) {
		return s.onMissingAndRebuildBody(body, raw, setPath, join, -1, tgtID, "source-or-target-key-empty")
	} else if checker.NotEquals(srcID, tgtID) {
		return s.onMissingAndRebuildBody(body, raw, setPath, join, -1, tgtID, "target-not-equals-source")
	}

	var errs []error

	out, err := s.jsonPath.Set(raw, setPath, source.Raw())
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err,
			"join failed: op=set reason=target-equals-source setPath=%s id=%s target.as=%s",
			setPath, tgtID, join.Target().As()))
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-body object->object"))
	}

	return newBody, errs
}

func (s joinService) onMissingAndRebuildBody(
	body *vo.Body,
	raw,
	setPath string,
	join vo.Join,
	targetIndex int,
	targetKeyValue,
	reason string,
) (*vo.Body, []error) {
	var errs []error

	out, err := s.applyOnMissingSingle(raw, setPath, join, targetIndex, targetKeyValue, reason)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=apply-missing-single"))
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-body"))
	}

	return newBody, errs
}

func (s joinService) applyMatchOnArrayItem(
	join vo.Join,
	raw string,
	setPath string,
	targetIndex int,
	targetKey domain.JSONValue,
	index map[string]string,
) (string, error) {
	if targetKey.NotExists() {
		return s.applyOnMissingSingle(raw, setPath, join, targetIndex, "", "target-key-not-found")
	} else if targetKey.IsArray() {
		jsonArr, errs := s.mapTargetIDsToJSONArray(join, targetIndex, targetKey, index)

		out, err := s.jsonPath.Set(raw, setPath, jsonArr)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "join failed: op=set-json-path"))
		}

		return out, errors.JoinInheritf(errs, ", ", "join failed: op=set-target-array")
	}

	id := strings.TrimSpace(targetKey.String())
	if checker.IsEmpty(id) {
		return s.applyOnMissingSingle(raw, setPath, join, targetIndex, id, "target-key-empty")
	}

	if srcRaw, ok := index[id]; ok {
		out, err := s.jsonPath.Set(raw, setPath, srcRaw)
		if checker.NonNil(err) {
			return raw, err
		}
		return out, nil
	}

	return s.applyOnMissingSingle(raw, setPath, join, targetIndex, id, "source-not-found")
}

func (s joinService) applyObjectMatchOnArrayItem(
	join vo.Join,
	raw string,
	setPath string,
	targetIndex int,
	targetKey domain.JSONValue,
	sourceKey,
	sourceRaw string,
) (string, error) {
	if targetKey.NotExists() {
		return s.applyOnMissingSingle(raw, setPath, join, targetIndex, "", "target-key-not-found")
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

		return s.applyOnMissingSingle(raw, setPath, join, targetIndex, sourceKey, "target-array-not-contains-source")
	}

	id := strings.TrimSpace(targetKey.String())
	if checker.IsEmpty(id) {
		return s.applyOnMissingSingle(raw, setPath, join, targetIndex, id, "target-key-empty")
	}

	if checker.Equals(id, sourceKey) {
		out, err := s.jsonPath.Set(raw, setPath, sourceRaw)
		if checker.NonNil(err) {
			return raw, err
		}
		return out, nil
	}

	return s.applyOnMissingSingle(raw, setPath, join, targetIndex, id, "target-not-equals-source")
}

func (s joinService) applyTargetKeyPolicy(joinTarget vo.JoinTarget, body *vo.Body) (*vo.Body, []error) {
	if joinTarget.KeepKey() {
		return body, nil
	}

	raw, err := body.Raw()
	if checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "join failed: op=raw path=%s key=%s", joinTarget.Path(), joinTarget.Key())
	}

	target := s.jsonPath.Get(raw, joinTarget.Path())
	if target.NotExists() {
		return body, nil
	}

	out := raw

	var errs []error
	switch joinTarget.Policy() {
	case enum.JoinTargetDropKeyAlways:
		out, errs = s.removeTargetKeyAll(joinTarget, raw, target)
	case enum.JoinTargetDropKeyOnMerged:
		out, errs = s.removeTargetKeyOnlyMerged(joinTarget, raw, target)
	default:
		return body, nil
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "join failed: op=build-body"))
	}

	return newBody, errs
}

func (s joinService) getTargetByPath(join vo.Join, body *vo.Body) (domain.JSONValue, error) {
	targetRaw, err := body.Raw()
	if checker.NonNil(err) {
		return nil, errors.Inheritf(err, "join failed: op=raw")
	}

	return s.jsonPath.Get(targetRaw, join.Target().Path()), nil
}

func (s joinService) getSourceByPath(
	join vo.Join,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (domain.JSONValue, error) {
	sourceValue, sourceErrs := s.dynamicValueService.Get(join.Source().Path(), request, history)
	if checker.IsNotEmpty(sourceErrs) {
		return nil, errors.JoinInheritf(sourceErrs, ", ", "join failed: op=resolve-source path=%s key=%s",
			join.Source().Path(), join.Source().Key())
	}
	return s.jsonPath.Parse(strings.TrimSpace(sourceValue)), nil
}

func (s joinService) mapTargetIDsToJSONArray(
	join vo.Join,
	targetIndex int,
	targetIDs domain.JSONValue,
	index map[string]string,
) (string, []error) {
	var out []string
	var errs []error

	targetIDs.ForEach(func(_ string, v domain.JSONValue) bool {
		id := strings.TrimSpace(v.String())

		if checker.IsEmpty(id) {
			appendIt, value, err := s.applyOnMissingCollection(join, targetIndex, id, "target-key-item-empty")
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

		appendIt, value, err := s.applyOnMissingCollection(join, targetIndex, id, "source-not-found")
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "join failed: source not found: id=%s", id))
		} else if appendIt {
			out = append(out, value)
		}
		return true
	})

	return s.buildJSONArrayRaw(out), errs
}

func (s joinService) removeTargetKeyAll(
	joinTarget vo.JoinTarget,
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

func (s joinService) removeTargetKeyOnlyMerged(
	joinTarget vo.JoinTarget,
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

func (s joinService) buildJSONArrayRaw(items []string) string {
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

func (s joinService) targetJSONPathOnArray(target vo.JoinTarget, i int) string {
	return fmt.Sprintf("%s.%d.%s", target.Path(), i, target.As())
}

func (s joinService) targetJSONPathOnObject(target vo.JoinTarget) string {
	return fmt.Sprintf("%s.%s", target.Path(), target.As())
}

func (s joinService) targetKeyJSONPathOnArray(target vo.JoinTarget, i int) string {
	return fmt.Sprintf("%s.%d.%s", target.Path(), i, target.Key())
}

func (s joinService) targetKeyJSONPathOnObject(target vo.JoinTarget) string {
	return fmt.Sprintf("%s.%s", target.Path(), target.Key())
}

func (s joinService) assertPathsAndTypes(join vo.Join, source, target domain.JSONValue) error {
	if !join.Target().IsOnMissingError() {
		return nil
	} else if !source.Exists() {
		return errors.Newf("join failed: source path not found: source.path=%s", join.Source().Path())
	} else if !target.Exists() {
		return errors.Newf("join failed: target path not found: target.path=%s", join.Target().Path())
	}

	if !(source.IsArray() || source.IsObject()) || !(target.IsArray() || target.IsObject()) {
		return errors.Newf(
			"join failed: source/target must be array or object: source.type=%s; target.type=%s; source.path=%s; target.path=%s",
			source.Type(),
			target.Type(),
			join.Source().Path(),
			join.Target().Path(),
		)
	}

	return nil
}

func (s joinService) evalJoinGuards(
	join vo.Join,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (bool, error) {
	shouldRun, _, errs := s.dynamicValueService.EvalGuards(join.OnlyIf(), join.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInherit(errs, ", ", "failed to evaluate guard for join")
	}
	return shouldRun, nil
}

func (s joinService) buildIndexFromArray(arr domain.JSONValue, key string) map[string]string {
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

func (s joinService) applyOnMissingSingle(
	raw,
	setPath string,
	join vo.Join,
	targetIndex int,
	targetKeyValue,
	reason string,
) (string, error) {
	switch join.Target().OnMissing() {
	case enum.JoinTargetOnMissingOmit:
		return raw, nil
	case enum.JoinTargetOnMissingError:
		return raw, errors.Newf(
			"join failed (missing data): reason=%s; target.path=%s; target.key=%s; target.index=%d; target.value=%s; source.path=%s; source.key=%s; target.as=%s",
			reason,
			join.Target().Path(),
			join.Target().Key(),
			targetIndex,
			targetKeyValue,
			join.Source().Path(),
			join.Source().Key(),
			join.Target().As(),
		)
	default:
		return s.jsonPath.Set(raw, setPath, "null")
	}
}

func (s joinService) applyOnMissingCollection(
	join vo.Join,
	targetIndex int,
	targetKeyValue string,
	reason string,
) (bool, string, error) {
	switch join.Target().OnMissing() {
	case enum.JoinTargetOnMissingOmit:
		return false, "", nil
	case enum.JoinTargetOnMissingNull:
		return true, "null", nil
	case enum.JoinTargetOnMissingError:
		return false, "", errors.Newf(
			"join failed (missing data): reason=%s; target.path=%s; target.key=%s; target.index=%d; target.value=%s; source.path=%s; source.key=%s; target.as=%s",
			reason,
			join.Target().Path(),
			join.Target().Key(),
			targetIndex,
			targetKeyValue,
			join.Source().Path(),
			join.Source().Key(),
			join.Target().As(),
		)
	default:
		return true, "null", nil
	}
}

func (s joinService) newBodyByString(body *vo.Body, modifiedBodyJson string) (*vo.Body, error) {
	buffer, err := converter.ToBufferWithErr(modifiedBodyJson)
	if checker.NonNil(err) {
		return body, err
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}
