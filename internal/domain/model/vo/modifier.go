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
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type Modifier struct {
	action    enum.ModifierAction
	propagate bool
	key       string
	value     string
}

func NewModifier(action enum.ModifierAction, propagate bool, key, value string) Modifier {
	return Modifier{
		action:    action,
		propagate: propagate,
		key:       key,
		value:     value,
	}
}

func (m Modifier) Action() enum.ModifierAction {
	return m.action
}

func (m Modifier) Propagate() bool {
	return m.propagate
}

func (m Modifier) Key() string {
	return m.key
}

func (m Modifier) Value() string {
	return m.value
}
