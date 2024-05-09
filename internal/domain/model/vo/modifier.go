/*
 * Copyright 2024 Gabriel Cataldo
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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

// Modifier represents a modification that can be applied to a httpRequest or httpResponse in the Gopen application.
type Modifier struct {
	// action represents the action to be performed in the Modifier struct.
	// It is an enum.ModifierAction value and can be one of the following values:
	// - ModifierActionSet: to set a value.
	// - ModifierActionAdd: to add a value.
	// - ModifierActionDel: to delete a value.
	// - ModifierActionReplace: to replace a value.
	action enum.ModifierAction
	// propagate represents a boolean flag that indicates whether the modification should be propagated to subsequent
	// Backend requests.
	propagate bool
	// key represents a string value that serves as the key for a modification in the Modifier structure.
	// Indicates the field that you want to modify.
	key string
	// value represents a string value in the Modifier struct.
	// It is used as a field to store the value of a modification.
	value DynamicValue
}

func newModifier(modifierJson ModifierJson) Modifier {
	return Modifier{
		action:    modifierJson.Action,
		propagate: modifierJson.Propagate,
		key:       modifierJson.Key,
		value:     NewDynamicValue(modifierJson.Value),
	}
}

// Action returns the value of the action field in the Modifier struct.
func (m Modifier) Action() enum.ModifierAction {
	return m.action
}

// Propagate returns the value of the propagate field in the Modifier struct.
func (m Modifier) Propagate() bool {
	return m.propagate
}

// Key returns the value of the key field in the Modifier struct.
func (m Modifier) Key() string {
	return m.key
}

func (m Modifier) ValueAsString(httpRequest *HttpRequest, httpResponse *HttpResponse) string {
	return m.value.AsString(httpRequest, httpResponse)
}

func (m Modifier) ValueAsInt(httpRequest *HttpRequest, httpResponse *HttpResponse) int {
	return m.value.AsInt(httpRequest, httpResponse)
}

func (m Modifier) ValueAsSliceOfString(httpRequest *HttpRequest, httpResponse *HttpResponse) []string {
	return m.value.AsSliceOfString(httpRequest, httpResponse)
}

// Valid checks if a Modifier is valid.
// A Modifier is considered valid if both the Modifier and its value are not empty.
func (m Modifier) Valid() bool {
	return helper.IsNotEmpty(m)
}

// Invalid checks if the Modifier is invalid. It returns true if the Modifier is not valid, otherwise it returns false.
func (m Modifier) Invalid() bool {
	return !m.Valid()
}
