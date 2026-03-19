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

type JoinConfig struct {
	onlyIf   []string
	ignoreIf []string
	source   JoinConfigSource
	target   JoinConfigTarget
}

type JoinConfigSource struct {
	path string
	key  string
}

type JoinConfigTarget struct {
	policy    enum.JoinTargetPolicy
	path      string
	key       string
	as        string
	onMissing enum.JoinTargetOnMissing
}

func NewJoin(
	onlyIf,
	ignoreIf []string,
	source JoinConfigSource,
	target JoinConfigTarget,
) JoinConfig {
	return JoinConfig{
		onlyIf:   onlyIf,
		ignoreIf: ignoreIf,
		source:   source,
		target:   target,
	}
}

func NewJoinSource(path, key string) JoinConfigSource {
	return JoinConfigSource{
		path: path,
		key:  key,
	}
}

func NewJoinTarget(policy enum.JoinTargetPolicy, path, key, as string, onMissing enum.JoinTargetOnMissing) JoinConfigTarget {
	return JoinConfigTarget{
		policy:    policy,
		path:      path,
		key:       key,
		as:        as,
		onMissing: onMissing,
	}
}

func (j JoinConfig) OnlyIf() []string {
	return j.onlyIf
}

func (j JoinConfig) IgnoreIf() []string {
	return j.ignoreIf
}

func (j JoinConfig) Source() JoinConfigSource {
	return j.source
}

func (j JoinConfig) Target() JoinConfigTarget {
	return j.target
}

func (s JoinConfigSource) Path() string {
	return s.path
}

func (s JoinConfigSource) Key() string {
	return s.key
}

func (e JoinConfigTarget) Policy() enum.JoinTargetPolicy {
	return e.policy
}

func (e JoinConfigTarget) Path() string {
	return e.path
}

func (e JoinConfigTarget) Key() string {
	return e.key
}

func (e JoinConfigTarget) As() string {
	return e.as
}

func (e JoinConfigTarget) OnMissing() enum.JoinTargetOnMissing {
	return e.onMissing
}

func (e JoinConfigTarget) IsOnMissingError() bool {
	return checker.Equals(e.OnMissing(), enum.JoinTargetOnMissingError)
}

func (e JoinConfigTarget) KeepKey() bool {
	return checker.Equals(e.policy, enum.JoinTargetKeepKey)
}
