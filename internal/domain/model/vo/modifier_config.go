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

type ModifierConfig struct {
	ignoreIf  []string
	onlyIf    []string
	action    enum.ModifierAction
	propagate bool
	key       string
	value     string
}

func NewModifierConfig(
	ignoreIf,
	onlyIf []string,
	action enum.ModifierAction,
	propagate bool,
	key,
	value string,
) ModifierConfig {
	return ModifierConfig{
		ignoreIf:  ignoreIf,
		onlyIf:    onlyIf,
		action:    action,
		propagate: propagate,
		key:       key,
		value:     value,
	}
}

func (m ModifierConfig) HasOnlyIf() bool {
	return checker.IsNotEmpty(m.onlyIf)
}

func (m ModifierConfig) HasIgnoreIf() bool {
	return checker.IsNotEmpty(m.ignoreIf)
}

func (m ModifierConfig) OnlyIf() []string {
	return m.onlyIf
}

func (m ModifierConfig) IgnoreIf() []string {
	return m.ignoreIf
}

func (m ModifierConfig) Action() enum.ModifierAction {
	return m.action
}

func (m ModifierConfig) Propagate() bool {
	return m.propagate
}

func (m ModifierConfig) Key() string {
	return m.key
}

func (m ModifierConfig) Value() string {
	return m.value
}
