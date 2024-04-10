package enum

type ModifierContext string
type ModifierScope string
type ModifierAction string
type CacheControl string
type ResponseEncode string
type ContentType string

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
const (
	ContentTypeJson ContentType = "JSON"
	ContentTypeXml  ContentType = "XML"
	ContentTypeYml  ContentType = "YML"
	ContentTypeText ContentType = "TEXT"
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

func (r ResponseEncode) ContentType() ContentType {
	switch r {
	case ResponseEncodeJson:
		return ContentTypeJson
	case ResponseEncodeXml:
		return ContentTypeXml
	case ResponseEncodeYaml:
		return ContentTypeYml
	default:
		return ContentTypeText
	}
}

func (c ContentType) IsEnumValid() bool {
	switch c {
	case ContentTypeText, ContentTypeJson, ContentTypeXml, ContentTypeYml:
		return true
	}
	return false
}

func (c ContentType) String() string {
	switch c {
	case ContentTypeJson:
		return "application/json"
	case ContentTypeXml:
		return "application/xml"
	case ContentTypeYml:
		return "application/x-yaml"
	default:
		return "text/plain"
	}
}
