package enum

type ModifierScope string
type ModifierAction string

const (
	ModifierScopeRequest  ModifierScope = "REQUEST"
	ModifierScopeResponse ModifierScope = "RESPONSE"
)
const (
	ModifierActionSet     ModifierAction = "SET"
	ModifierActionAdd     ModifierAction = "ADD"
	ModifierActionDel     ModifierAction = "DEL"
	ModifierActionReplace ModifierAction = "REPLACE"
)
