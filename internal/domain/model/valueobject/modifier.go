package valueobject

import (
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/dto"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/enum"
)

type Modifier struct {
	Scope  []enum.ModifierScope
	Action enum.ModifierAction
	Key    string
	Value  string
}

func BuildModifier(modifier dto.Modifier) Modifier {
	return Modifier{
		Scope:  modifier.Scope,
		Action: modifier.Action,
		Key:    modifier.Key,
		Value:  modifier.Value,
	}
}

func BuildModifiers(modifiers []dto.Modifier) []Modifier {
	var result []Modifier
	for _, modifier := range modifiers {
		result = append(result, BuildModifier(modifier))
	}
	return result
}
