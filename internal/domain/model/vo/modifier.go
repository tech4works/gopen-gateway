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
	// context represents the context in which a modification should be applied.
	// It is an enum.ModifierContext value.
	// Valid values for Context are "httpRequest" and "httpResponse".
	context enum.ModifierContext
	// scope represents the scope of a modification in the Backend or Endpoint.
	// It is an enum.ModifierScope value that specifies where the modification should be applied.
	// Valid values for Scope are "httpRequest" and "httpResponse".
	scope enum.ModifierScope
	// action represents the action to be performed in the Modifier struct.
	// It is an enum.ModifierAction value and can be one of the following values:
	// - ModifierActionSet: to set a value.
	// - ModifierActionAdd: to add a value.
	// - ModifierActionDel: to delete a value.
	// - ModifierActionReplace: to replace a value.
	// - ModifierActionRename: to rename a value.
	action enum.ModifierAction
	// propagate represents a boolean flag that indicates whether the modification should be propagated to subsequent
	// Backend requests.
	propagate bool
	// key represents a string value that serves as the key for a modification in the Modifier structure.
	// Indicates the field that you want to modify.
	key string
	// value represents a string value in the Modifier struct.
	// It is used as a field to store the value of a modification.
	value string
}

// newModifier creates a new instance of Modifier based on the provided ModifierJson.
// If the provided ModifierJson is nil, it returns nil.
// It initializes a new Modifier with the values from the ModifierJson and returns a pointer to it.
// The context field of the Modifier struct is populated with the context value from the ModifierJson,
// the scope field is populated with the scope value,
// the action field is populated with the action value,
// the propagate field is populated with the propagate flag,
// the key field is populated with the key value,
// and the value field is populated with the value.
// Other fields of the Modifier struct are not populated and will have their zero values.
func newModifier(modifierJsonVO *ModifierJson) *Modifier {
	if helper.IsNil(modifierJsonVO) {
		return nil
	}
	return &Modifier{
		context:   modifierJsonVO.Context,
		scope:     modifierJsonVO.Scope,
		action:    modifierJsonVO.Action,
		propagate: modifierJsonVO.Propagate,
		key:       modifierJsonVO.Key,
		value:     modifierJsonVO.Value,
	}
}

// newModifierFromValue creates a new instance of Modifier with the provided context and value.
// It initializes a new Modifier with the given context and value, and returns a pointer to it.
// The context field of the Modifier struct is populated with the provided context value,
// and the value field is populated with the provided value.
// Other fields of the Modifier struct are not populated and will have their zero values.
func newModifierFromValue(context enum.ModifierContext, value string) *Modifier {
	return &Modifier{
		context: context,
		value:   value,
	}
}

// EqualsContext checks if the context of the Modifier is equal to the given enum.ModifierContext.
// It returns true if the context is empty or if it is equal to the given context, otherwise it returns false.
func (m Modifier) EqualsContext(context enum.ModifierContext) bool {
	return helper.IsEmpty(m.context) || helper.Equals(m.context, context)
}

// NotEqualsContext returns `true` if the `Modifier` context is not equal to the specified `enum.ModifierContext`, otherwise `false`.
// It uses the `EqualsContext` method to check for equality.
func (m Modifier) NotEqualsContext(context enum.ModifierContext) bool {
	return !m.EqualsContext(context)
}

// Context returns the value of the context field in the Modifier struct.
func (m Modifier) Context() enum.ModifierContext {
	return m.context
}

// Scope returns the value of the scope field in the Modifier struct.
func (m Modifier) Scope() enum.ModifierScope {
	return m.scope
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

// Value returns the value of the value field in the Modifier struct.
func (m Modifier) Value() string {
	return m.value
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
