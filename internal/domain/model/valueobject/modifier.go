package valueobject

import "github.com/GabrielHCataldo/martini-gateway/internal/domain/enum"

type Modifier struct {
	Scope  []enum.ModifierScope
	Action enum.ModifierAction
	Key    string
	Value  string
}
