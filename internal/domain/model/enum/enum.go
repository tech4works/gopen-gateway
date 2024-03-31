package enum

type ModifierContext string
type ModifierScope string
type ModifierAction string
type CacheControl string
type ResponseEncode string

const (
	ResponseEncodeText ResponseEncode = "TEXT"
	ResponseEncodeJson ResponseEncode = "JSON"
	ResponseEncodeXml  ResponseEncode = "XML"
	ResponseEncodeYaml ResponseEncode = "YAML"
)
const (
	CacheControlNoCache CacheControl = "no-cache"
	CacheControlNoStore CacheControl = "no-store"
)
const (
	ModifierScopeRequest  ModifierScope = "REQUEST"
	ModifierScopeResponse ModifierScope = "RESPONSE"
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
	case ModifierScopeRequest, ModifierScopeResponse:
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
	case CacheControlNoCache, CacheControlNoStore:
		return true
	}
	return false
}

func (r ResponseEncode) IsEnumValid() bool {
	switch r {
	case ResponseEncodeText, ResponseEncodeJson, ResponseEncodeXml:
		return true
	}
	return false
}
