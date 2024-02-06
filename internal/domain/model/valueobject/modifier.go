package valueobject

import (
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/enum"
)

type Modifier struct {
	Scope  []enum.ModifierScope
	Action enum.ModifierAction
	Key    string
	Value  string
}
