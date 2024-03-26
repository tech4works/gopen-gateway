package enum

type ModifierContext string
type ModifierScope string
type ModifierAction string
type CacheControl string

const (
	CacheControlNoCache CacheControl = "no-cache"
)

const (
	ModifierScopeLocal  ModifierScope = "LOCAL"
	ModifierScopeGlobal ModifierScope = "GLOBAL"
)
const (
	ModifierContextRequest  ModifierContext = "REQUEST"
	ModifierContextResponse ModifierContext = "RESPONSE"
)
const (
	ModifierActionSet     ModifierAction = "SET"
	ModifierActionAdd     ModifierAction = "ADD"
	ModifierActionDel     ModifierAction = "DEL"
	ModifierActionReplace ModifierAction = "REPLACE"
	ModifierActionRename  ModifierAction = "RENAME"
)

func (m ModifierScope) IsEnumValid() bool {
	switch m {
	case ModifierScopeLocal, ModifierScopeGlobal:
		return true
	}
	return false
}

func (m ModifierContext) IsEnumValid() bool {
	switch m {
	case ModifierContextRequest, ModifierContextResponse:
		return true
	}
	return false
}

func (m ModifierAction) IsEnumValid() bool {
	switch m {
	case ModifierActionSet, ModifierActionAdd, ModifierActionDel, ModifierActionReplace, ModifierActionRename:
		return true
	}
	return false
}

func (c CacheControl) IsEnumValid() bool {
	switch c {
	case CacheControlNoCache:
		return true
	}
	return false
}
