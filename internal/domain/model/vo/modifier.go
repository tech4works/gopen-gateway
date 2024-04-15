package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

// Modifier represents a modification that can be applied to a request or response in the Gopen application.
type Modifier struct {
	// context represents the context in which a modification should be applied.
	// It is an enum.ModifierContext value.
	// Valid values for Context are "request" and "response".
	context enum.ModifierContext
	// scope represents the scope of a modification in the Backend or Endpoint.
	// It is an enum.ModifierScope value that specifies where the modification should be applied.
	// Valid values for Scope are "request" and "response".
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

// newModifier creates a new instance of Modifier struct
// with the provided modifierDTO.
// It sets the context, scope, action, propagate, key, and value fields of the Modifier struct
// to the corresponding fields of the modifierDTO parameter.
// Returns the created Modifier struct.
func newModifier(modifierDTO dto.Modifier) Modifier {
	return Modifier{
		context:   modifierDTO.Context,
		scope:     modifierDTO.Scope,
		action:    modifierDTO.Action,
		propagate: modifierDTO.Propagate,
		key:       modifierDTO.Key,
		value:     modifierDTO.Value,
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
