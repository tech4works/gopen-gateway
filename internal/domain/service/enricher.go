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
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type enricherService struct {
	jsonPath            domain.JSONPath
	dynamicValueService DynamicValue
}

type Enricher interface {
	ExecuteBodyEnrichers(enrichers []vo.Enricher, body *vo.Body, request *vo.HTTPRequest, history *vo.History) (*vo.Body, []error)
	EnrichBody(enricher vo.Enricher, body *vo.Body, request *vo.HTTPRequest, history *vo.History) (*vo.Body, []error)
}

func NewEnricher(jsonPath domain.JSONPath, dynamicValueService DynamicValue) Enricher {
	return enricherService{
		jsonPath:            jsonPath,
		dynamicValueService: dynamicValueService,
	}
}

func (s enricherService) ExecuteBodyEnrichers(
	enrichers []vo.Enricher,
	body *vo.Body,
	request *vo.HTTPRequest,
	history *vo.History,
) (*vo.Body, []error) {
	if checker.IsNil(body) {
		return nil, nil
	}

	var allErrs []error

	for _, enr := range enrichers {
		newBody, errs := s.EnrichBody(enr, body, request, history)
		if checker.IsNotEmpty(errs) {
			allErrs = append(allErrs, errs...)
		}
		body = newBody
	}

	return body, allErrs
}

func (s enricherService) EnrichBody(
	enricher vo.Enricher,
	body *vo.Body,
	request *vo.HTTPRequest,
	history *vo.History,
) (*vo.Body, []error) {
	if checker.IsNil(body) || body.ContentType().IsNotJSON() {
		return body, nil
	}

	shouldRun, err := s.evalEnricherGuards(enricher, request, history)
	if checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "enricher failed: op=eval-guards source.path=%s target.path=%s",
			enricher.Source().Path(), enricher.Target().Path())
	}
	if !shouldRun {
		return body, nil
	}

	source, err := s.getSourceByPath(enricher, request, history)
	if checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "enricher failed: op=get-source source.path=%s",
			enricher.Source().Path())
	}

	target, err := s.getTargetByPath(enricher, body)
	if checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "enricher failed: op=get-target target.path=%s",
			enricher.Target().Path())
	}

	if err = s.assertPathsAndTypes(enricher, source, target); checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "enricher failed: op=assert-paths-types source.path=%s target.path=%s",
			enricher.Source().Path(), enricher.Target().Path())
	}

	var allErrs []error

	newBody, errs := s.enrichByShape(enricher, body, source, target)
	if checker.IsNotEmpty(errs) {
		for _, e := range errs {
			allErrs = append(allErrs, errors.Inheritf(e,
				"enricher failed: op=enrich source.path=%s source.key=%s target.path=%s target.key=%s target.as=%s",
				enricher.Source().Path(), enricher.Source().Key(), enricher.Target().Path(), enricher.Target().Key(),
				enricher.Target().As()))
		}
	}

	newBody, errs = s.applyTargetKeyPolicy(enricher.Target(), newBody)
	if checker.NonNil(err) {
		for _, e := range errs {
			allErrs = append(allErrs, errors.Inheritf(e,
				"enricher failed: op=apply-target-key-policy source.path=%s source.key=%s target.path=%s target.key=%s target.as=%s",
				enricher.Source().Path(), enricher.Source().Key(), enricher.Target().Path(), enricher.Target().Key(),
				enricher.Target().As()))
		}
	}

	return newBody, allErrs
}

func (s enricherService) enrichByShape(
	enricher vo.Enricher,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	switch {
	case source.Exists() && source.IsArray() && target.Exists() && target.IsArray():
		return s.enrichArrayToArray(enricher, body, source, target)
	case source.Exists() && source.IsObject() && target.Exists() && target.IsArray():
		return s.enrichObjectToArray(enricher, body, source, target)
	case source.Exists() && source.IsArray() && target.Exists() && target.IsObject():
		return s.enrichArrayToObject(enricher, body, source, target)
	case source.Exists() && source.IsObject() && target.Exists() && target.IsObject():
		return s.enrichObjectToObject(enricher, body, source, target)
	default:
		return body, nil
	}
}

func (s enricherService) enrichArrayToArray(
	enricher vo.Enricher,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	index := s.buildIndexFromArray(source, enricher.Source().Key())
	out := target.Raw()

	var errs []error
	var i int
	target.ForEach(func(_ string, item domain.JSONValue) bool {
		var err error

		setPath := s.targetJSONPathOnArray(enricher.Target(), i)
		targetKey := item.Get(enricher.Target().Key())

		out, err = s.applyMatchOnArrayItem(enricher, out, setPath, i, targetKey, index)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err,
				"enricher failed: op=apply-match array idx=%d setPath=%s target.key=%s", i, setPath,
				enricher.Target().Key()))
		}

		i++
		return true
	})

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "enricher failed: op=buffer array out"))
	}

	return newBody, errs
}

func (s enricherService) enrichObjectToArray(
	enricher vo.Enricher,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	sourceKey := strings.TrimSpace(source.Get(enricher.Source().Key()).String())
	if checker.IsEmpty(sourceKey) {
		if enricher.Target().IsOnMissingError() {
			return body, errors.NewAsSlicef("enrich failed: source key missing/empty: source.path=%s; source.key=%s",
				enricher.Source().Path(), enricher.Source().Key())
		}
		return body, nil
	}

	out := target.Raw()

	var errs []error
	var i int
	target.ForEach(func(_ string, item domain.JSONValue) bool {
		var err error

		setPath := s.targetJSONPathOnArray(enricher.Target(), i)
		targetKey := item.Get(enricher.Target().Key())

		out, err = s.applyObjectMatchOnArrayItem(enricher, out, setPath, i, targetKey, sourceKey, source.Raw())
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err,
				"enricher failed: op=apply-match array idx=%d setPath=%s target.key=%s",
				i, setPath, enricher.Target().Key()))
		}

		i++
		return true
	})

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "enricher failed: op=buffer array out"))
	}

	return newBody, errs
}

func (s enricherService) enrichArrayToObject(
	enricher vo.Enricher,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	index := s.buildIndexFromArray(source, enricher.Source().Key())
	setPath := s.targetJSONPathOnObject(enricher.Target())
	tgtKey := target.Get(enricher.Target().Key())

	if tgtKey.NotExists() {
		return s.onMissingAndRebuildBody(body, target.Raw(), setPath, enricher, -1, "", "target-key-not-found")
	} else if tgtKey.IsArray() {
		return s.enrichObjectTargetFromIDArray(enricher, body, target.Raw(), setPath, tgtKey, index)
	} else {
		return s.enrichObjectTargetFromSingleID(enricher, body, target.Raw(), setPath, tgtKey, index)
	}
}

func (s enricherService) enrichObjectTargetFromIDArray(
	enricher vo.Enricher,
	body *vo.Body,
	raw string,
	setPath string,
	tgtKey domain.JSONValue,
	index map[string]string,
) (*vo.Body, []error) {
	jsonArr, errs := s.mapTargetIDsToJSONArray(enricher, -1, tgtKey, index)

	out, err := s.jsonPath.Set(raw, setPath, jsonArr)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "enricher failed: op=set setPath=%s mode=array target.path=%s target.as=%s",
			setPath, enricher.Target().Path(), enricher.Target().As()))
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "enricher failed: op=build-body array->object"))
	}

	return newBody, errs
}

func (s enricherService) enrichObjectTargetFromSingleID(
	enricher vo.Enricher,
	body *vo.Body,
	raw,
	setPath string,
	tgtKey domain.JSONValue,
	index map[string]string,
) (*vo.Body, []error) {
	id := strings.TrimSpace(tgtKey.String())
	if checker.IsEmpty(id) {
		return s.onMissingAndRebuildBody(body, raw, setPath, enricher, -1, id, "target-key-empty")
	}

	srcRaw, ok := index[id]
	if !ok {
		return s.onMissingAndRebuildBody(body, raw, setPath, enricher, -1, id, "source-not-found")
	}

	var errs []error

	out, err := s.jsonPath.Set(raw, setPath, srcRaw)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err,
			"enricher failed: op=set setPath=%s mode=single id=%s target.path=%s target.as=%s",
			setPath, id, enricher.Target().Path(), enricher.Target().As()))
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "enricher failed: op=build-body array->object"))
	}

	return newBody, errs
}

func (s enricherService) enrichObjectToObject(
	enricher vo.Enricher,
	body *vo.Body,
	source,
	target domain.JSONValue,
) (*vo.Body, []error) {
	setPath := s.targetJSONPathOnObject(enricher.Target())

	srcID := strings.TrimSpace(source.Get(enricher.Source().Key()).String())
	tgtID := strings.TrimSpace(target.Get(enricher.Target().Key()).String())

	raw := target.Raw()

	if checker.IsEmpty(srcID) || checker.IsEmpty(tgtID) {
		return s.onMissingAndRebuildBody(body, raw, setPath, enricher, -1, tgtID, "source-or-target-key-empty")
	} else if checker.NotEquals(srcID, tgtID) {
		return s.onMissingAndRebuildBody(body, raw, setPath, enricher, -1, tgtID, "target-not-equals-source")
	}

	var errs []error

	out, err := s.jsonPath.Set(raw, setPath, source.Raw())
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err,
			"enricher failed: op=set reason=target-equals-source setPath=%s id=%s target.as=%s",
			setPath, tgtID, enricher.Target().As()))
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "enricher failed: op=build-body object->object"))
	}

	return newBody, errs
}

func (s enricherService) onMissingAndRebuildBody(
	body *vo.Body,
	raw,
	setPath string,
	enricher vo.Enricher,
	targetIndex int,
	targetKeyValue,
	reason string,
) (*vo.Body, []error) {
	var errs []error

	out, err := s.applyOnMissingSingle(raw, setPath, enricher, targetIndex, targetKeyValue, reason)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "enricher failed: op=apply-missing-single"))
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "enricher failed: op=build-body"))
	}

	return newBody, errs
}

func (s enricherService) applyMatchOnArrayItem(
	enricher vo.Enricher,
	raw string,
	setPath string,
	targetIndex int,
	targetKey domain.JSONValue,
	index map[string]string,
) (string, error) {
	if targetKey.NotExists() {
		return s.applyOnMissingSingle(raw, setPath, enricher, targetIndex, "", "target-key-not-found")
	} else if targetKey.IsArray() {
		jsonArr, errs := s.mapTargetIDsToJSONArray(enricher, targetIndex, targetKey, index)

		out, err := s.jsonPath.Set(raw, setPath, jsonArr)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "enricher failed: op=set-json-path"))
		}

		return out, errors.JoinInheritf(errs, ", ", "enricher failed: op=set-target-array")
	}

	id := strings.TrimSpace(targetKey.String())
	if checker.IsEmpty(id) {
		return s.applyOnMissingSingle(raw, setPath, enricher, targetIndex, id, "target-key-empty")
	}

	if srcRaw, ok := index[id]; ok {
		out, err := s.jsonPath.Set(raw, setPath, srcRaw)
		if checker.NonNil(err) {
			return raw, err
		}
		return out, nil
	}

	return s.applyOnMissingSingle(raw, setPath, enricher, targetIndex, id, "source-not-found")
}

func (s enricherService) applyObjectMatchOnArrayItem(
	enricher vo.Enricher,
	raw string,
	setPath string,
	targetIndex int,
	targetKey domain.JSONValue,
	sourceKey,
	sourceRaw string,
) (string, error) {
	if targetKey.NotExists() {
		return s.applyOnMissingSingle(raw, setPath, enricher, targetIndex, "", "target-key-not-found")
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

		return s.applyOnMissingSingle(raw, setPath, enricher, targetIndex, sourceKey, "target-array-not-contains-source")
	}

	id := strings.TrimSpace(targetKey.String())
	if checker.IsEmpty(id) {
		return s.applyOnMissingSingle(raw, setPath, enricher, targetIndex, id, "target-key-empty")
	}

	if checker.Equals(id, sourceKey) {
		out, err := s.jsonPath.Set(raw, setPath, sourceRaw)
		if checker.NonNil(err) {
			return raw, err
		}
		return out, nil
	}

	return s.applyOnMissingSingle(raw, setPath, enricher, targetIndex, id, "target-not-equals-source")
}

func (s enricherService) applyTargetKeyPolicy(enricherTarget vo.EnricherTarget, body *vo.Body) (*vo.Body, []error) {
	if enricherTarget.KeepKey() {
		return body, nil
	}

	raw, err := body.Raw()
	if checker.NonNil(err) {
		return body, errors.InheritAsSlicef(err, "enricher failed: applyTargetKeyPolicy op=raw path=%s key=%s",
			enricherTarget.Path(), enricherTarget.Key())
	}

	target := s.jsonPath.Get(raw, enricherTarget.Path())
	if target.NotExists() {
		return body, nil
	}

	out := raw

	var errs []error
	switch enricherTarget.Policy() {
	case enum.EnrichTargetKeyDropAlways:
		out, errs = s.removeTargetKeyAll(enricherTarget, raw, target)
	case enum.EnrichTargetKeyDropOnEnrich:
		out, errs = s.removeTargetKeyOnlyEnriched(enricherTarget, raw, target)
	default:
		return body, nil
	}

	newBody, err := s.newBodyByString(body, out)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err, "enricher failed: applyTargetKeyPolicy op=build-body"))
	}

	return newBody, errs
}

func (s enricherService) getTargetByPath(enricher vo.Enricher, body *vo.Body) (domain.JSONValue, error) {
	targetRaw, err := body.Raw()
	if checker.NonNil(err) {
		return nil, errors.Inheritf(err, "enricher failed: getTargetByPath op=raw")
	}

	return s.jsonPath.Get(targetRaw, enricher.Target().Path()), nil
}

func (s enricherService) getSourceByPath(
	enricher vo.Enricher,
	request *vo.HTTPRequest,
	history *vo.History,
) (domain.JSONValue, error) {
	sourceValue, sourceErrs := s.dynamicValueService.Get(enricher.Source().Path(), request, history)
	if checker.IsNotEmpty(sourceErrs) {
		return nil, errors.JoinInheritf(sourceErrs, ", ", "enricher failed: op=resolve-source path=%s key=%s",
			enricher.Source().Path(), enricher.Source().Key())
	}
	return s.jsonPath.Parse(strings.TrimSpace(sourceValue)), nil
}

func (s enricherService) mapTargetIDsToJSONArray(
	enricher vo.Enricher,
	targetIndex int,
	targetIDs domain.JSONValue,
	index map[string]string,
) (string, []error) {
	var out []string
	var errs []error

	targetIDs.ForEach(func(_ string, v domain.JSONValue) bool {
		id := strings.TrimSpace(v.String())

		if checker.IsEmpty(id) {
			appendIt, value, err := s.applyOnMissingCollection(enricher, targetIndex, id, "target-key-item-empty")
			if checker.NonNil(err) {
				errs = append(errs, errors.Inherit(err, "enricher failed: missing collection item (empty)"))
			} else if appendIt {
				out = append(out, value)
			}
			return true
		}

		if srcRaw, ok := index[id]; ok {
			out = append(out, srcRaw)
			return true
		}

		appendIt, value, err := s.applyOnMissingCollection(enricher, targetIndex, id, "source-not-found")
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "enricher failed: source not found: id=%s", id))
		} else if appendIt {
			out = append(out, value)
		}
		return true
	})

	return s.buildJSONArrayRaw(out), errs
}

func (s enricherService) removeTargetKeyAll(
	enricherTarget vo.EnricherTarget,
	raw string,
	target domain.JSONValue,
) (string, []error) {
	var errs []error
	if target.IsArray() {
		out := raw

		var i int
		target.ForEach(func(_ string, _ domain.JSONValue) bool {
			var err error

			path := s.targetKeyJSONPathOnArray(enricherTarget, i)

			out, err = s.jsonPath.Delete(out, path)
			if checker.NonNil(err) {
				errs = append(errs, errors.Inheritf(err,
					"enricher failed: removeTargetKeyAll on array op=delete-target-key path=%s idx=%d", path, i))
			}

			i++
			return true
		})

		return out, errs
	}

	path := s.targetKeyJSONPathOnObject(enricherTarget)
	out, err := s.jsonPath.Delete(raw, path)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err,
			"enricher failed: removeTargetKeyAll on object op=delete-target-key path=%s", path))
	}

	return out, errs
}

func (s enricherService) removeTargetKeyOnlyEnriched(
	enricherTarget vo.EnricherTarget,
	raw string,
	target domain.JSONValue,
) (string, []error) {
	var errs []error
	if target.IsArray() {
		out := raw

		var i int
		target.ForEach(func(_ string, item domain.JSONValue) bool {
			if item.Get(enricherTarget.As()).Exists() {
				var err error

				path := s.targetKeyJSONPathOnArray(enricherTarget, i)

				out, err = s.jsonPath.Delete(out, path)
				if checker.NonNil(err) {
					errs = append(errs, errors.Inheritf(err,
						"enricher failed: removeTargetKeyOnlyEnriched on array op=delete-target-key path=%s idx=%d",
						path, i))
				}
			}

			i++
			return true
		})

		return out, errs
	} else if target.Get(enricherTarget.As()).NotExists() {
		return raw, nil
	}

	path := s.targetKeyJSONPathOnObject(enricherTarget)
	out, err := s.jsonPath.Delete(raw, path)
	if checker.NonNil(err) {
		errs = append(errs, errors.Inheritf(err,
			"enricher failed: removeTargetKeyOnlyEnriched on object op=delete-target-key path=%s", path))
	}
	return out, errs
}

func (s enricherService) buildJSONArrayRaw(items []string) string {
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

func (s enricherService) targetJSONPathOnArray(enricherTarget vo.EnricherTarget, i int) string {
	return fmt.Sprintf("%s.%d.%s", enricherTarget.Path(), i, enricherTarget.As())
}

func (s enricherService) targetJSONPathOnObject(enricherTarget vo.EnricherTarget) string {
	return fmt.Sprintf("%s.%s", enricherTarget.Path(), enricherTarget.As())
}

func (s enricherService) targetKeyJSONPathOnArray(enricherTarget vo.EnricherTarget, i int) string {
	return fmt.Sprintf("%s.%d.%s", enricherTarget.Path(), i, enricherTarget.Key())
}

func (s enricherService) targetKeyJSONPathOnObject(enricherTarget vo.EnricherTarget) string {
	return fmt.Sprintf("%s.%s", enricherTarget.Path(), enricherTarget.Key())
}

func (s enricherService) assertPathsAndTypes(enricher vo.Enricher, source, target domain.JSONValue) error {
	if !enricher.Target().IsOnMissingError() {
		return nil
	} else if !source.Exists() {
		return errors.Newf("enrich failed: source path not found: source.path=%s", enricher.Source().Path())
	} else if !target.Exists() {
		return errors.Newf("enrich failed: target path not found: target.path=%s", enricher.Target().Path())
	}

	if !(source.IsArray() || source.IsObject()) || !(target.IsArray() || target.IsObject()) {
		return errors.Newf(
			"enrich failed: source/target must be array or object: source.type=%s; target.type=%s; source.path=%s; target.path=%s",
			source.Type(),
			target.Type(),
			enricher.Source().Path(),
			enricher.Target().Path(),
		)
	}

	return nil
}

func (s enricherService) evalEnricherGuards(
	enricher vo.Enricher,
	request *vo.HTTPRequest,
	history *vo.History,
) (bool, error) {
	shouldRun, _, errs := s.dynamicValueService.EvalGuards(enricher.OnlyIf(), enricher.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInherit(errs, ", ", "failed to evaluate guard for enrich")
	}
	return shouldRun, nil
}

func (s enricherService) buildIndexFromArray(arr domain.JSONValue, key string) map[string]string {
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

func (s enricherService) applyOnMissingSingle(
	raw,
	setPath string,
	enricher vo.Enricher,
	targetIndex int,
	targetKeyValue,
	reason string,
) (string, error) {
	switch enricher.Target().OnMissing() {
	case enum.EnrichTargetOnMissingOmit:
		return raw, nil
	case enum.EnrichTargetOnMissingError:
		return raw, errors.Newf(
			"enrich failed (missing data): reason=%s; target.path=%s; target.key=%s; target.index=%d; target.value=%s; source.path=%s; source.key=%s; target.as=%s",
			reason,
			enricher.Target().Path(),
			enricher.Target().Key(),
			targetIndex,
			targetKeyValue,
			enricher.Source().Path(),
			enricher.Source().Key(),
			enricher.Target().As(),
		)
	default:
		return s.jsonPath.Set(raw, setPath, "null")
	}
}

func (s enricherService) applyOnMissingCollection(
	enricher vo.Enricher,
	targetIndex int,
	targetKeyValue string,
	reason string,
) (bool, string, error) {
	switch enricher.Target().OnMissing() {
	case enum.EnrichTargetOnMissingOmit:
		return false, "", nil
	case enum.EnrichTargetOnMissingNull:
		return true, "null", nil
	case enum.EnrichTargetOnMissingError:
		return false, "", errors.Newf(
			"enrich failed (missing data): reason=%s; target.path=%s; target.key=%s; target.index=%d; target.value=%s; source.path=%s; source.key=%s; target.as=%s",
			reason,
			enricher.Target().Path(),
			enricher.Target().Key(),
			targetIndex,
			targetKeyValue,
			enricher.Source().Path(),
			enricher.Source().Key(),
			enricher.Target().As(),
		)
	default:
		return true, "null", nil
	}
}

func (s enricherService) newBodyByString(body *vo.Body, modifiedBodyJson string) (*vo.Body, error) {
	buffer, err := converter.ToBufferWithErr(modifiedBodyJson)
	if checker.NonNil(err) {
		return body, err
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}
