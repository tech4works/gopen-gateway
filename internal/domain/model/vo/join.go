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

type Join struct {
	onlyIf   []string
	ignoreIf []string
	source   JoinSource
	target   JoinTarget
}

type JoinSource struct {
	path string
	key  string
}

type JoinTarget struct {
	policy    enum.JoinTargetPolicy
	path      string
	key       string
	as        string
	onMissing enum.JoinTargetOnMissing
}

func NewJoin(
	onlyIf,
	ignoreIf []string,
	source JoinSource,
	target JoinTarget,
) Join {
	return Join{
		onlyIf:   onlyIf,
		ignoreIf: ignoreIf,
		source:   source,
		target:   target,
	}
}

func NewJoinSource(path, key string) JoinSource {
	return JoinSource{
		path: path,
		key:  key,
	}
}

func NewJoinTarget(policy enum.JoinTargetPolicy, path, key, as string, onMissing enum.JoinTargetOnMissing) JoinTarget {
	return JoinTarget{
		policy:    policy,
		path:      path,
		key:       key,
		as:        as,
		onMissing: onMissing,
	}
}

func (j Join) OnlyIf() []string {
	return j.onlyIf
}

func (j Join) IgnoreIf() []string {
	return j.ignoreIf
}

func (j Join) Source() JoinSource {
	return j.source
}

func (j Join) Target() JoinTarget {
	return j.target
}

func (s JoinSource) Path() string {
	return s.path
}

func (s JoinSource) Key() string {
	return s.key
}

func (e JoinTarget) Policy() enum.JoinTargetPolicy {
	return e.policy
}

func (e JoinTarget) Path() string {
	return e.path
}

func (e JoinTarget) Key() string {
	return e.key
}

func (e JoinTarget) As() string {
	return e.as
}

func (e JoinTarget) OnMissing() enum.JoinTargetOnMissing {
	return e.onMissing
}

func (e JoinTarget) IsOnMissingError() bool {
	return checker.Equals(e.OnMissing(), enum.JoinTargetOnMissingError)
}

func (e JoinTarget) KeepKey() bool {
	return checker.Equals(e.policy, enum.JoinTargetKeepKey)
}
