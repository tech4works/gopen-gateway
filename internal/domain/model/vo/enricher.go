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

package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type Enricher struct {
	onlyIf   []string
	ignoreIf []string
	source   EnricherSource
	target   EnricherTarget
}

type EnricherSource struct {
	path string
	key  string
}

type EnricherTarget struct {
	policy    enum.EnrichTargetPolicy
	path      string
	key       string
	as        string
	onMissing enum.EnrichTargetOnMissing
}

func NewEnricher(
	onlyIf,
	ignoreIf []string,
	source EnricherSource,
	target EnricherTarget,
) Enricher {
	return Enricher{
		onlyIf:   onlyIf,
		ignoreIf: ignoreIf,
		source:   source,
		target:   target,
	}
}

func NewEnricherSource(path, key string) EnricherSource {
	return EnricherSource{
		path: path,
		key:  key,
	}
}

func NewEnricherTarget(policy enum.EnrichTargetPolicy, path, key, as string, onMissing enum.EnrichTargetOnMissing) EnricherTarget {
	return EnricherTarget{
		policy:    policy,
		path:      path,
		key:       key,
		as:        as,
		onMissing: onMissing,
	}
}

func (e Enricher) OnlyIf() []string {
	return e.onlyIf
}

func (e Enricher) IgnoreIf() []string {
	return e.ignoreIf
}

func (e Enricher) Source() EnricherSource {
	return e.source
}

func (e Enricher) Target() EnricherTarget {
	return e.target
}

func (s EnricherSource) Path() string {
	return s.path
}

func (s EnricherSource) Key() string {
	return s.key
}

func (e EnricherTarget) Policy() enum.EnrichTargetPolicy {
	return e.policy
}

func (e EnricherTarget) Path() string {
	return e.path
}

func (e EnricherTarget) Key() string {
	return e.key
}

func (e EnricherTarget) As() string {
	return e.as
}

func (e EnricherTarget) OnMissing() enum.EnrichTargetOnMissing {
	return e.onMissing
}

func (e EnricherTarget) IsOnMissingError() bool {
	return checker.Equals(e.OnMissing(), enum.EnrichTargetOnMissingError)
}

func (e EnricherTarget) KeepKey() bool {
	return checker.Equals(e.policy, enum.EnrichTargetKeyKeep)
}

func (e EnricherTarget) ShouldDropKeyAlways() bool {
	return checker.Equals(e.policy, enum.EnrichTargetKeyDropAlways)
}

func (e EnricherTarget) ShouldDropKeyOnEnrich() bool {
	return checker.Equals(e.policy, enum.EnrichTargetKeyDropOnEnrich)
}
