package enum

import "github.com/GabrielHCataldo/go-helper/helper"

// ModifierContext represents the context in which a modification should be applied.
// It is a string value that can be either "REQUEST" or "RESPONSE".
type ModifierContext string

// ModifierScope represents the scope of a modification in the Backend or Endpoint.
type ModifierScope string

// ModifierAction represents the action to be performed in the Modifier struct.
type ModifierAction string

// CacheControl represents the header value of cache control.
type CacheControl string

// ResponseEncode represents the encoding format for the API endpoint response.
type ResponseEncode string

// ContentType represents the format of the content.
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
	ModifierActionSet ModifierAction = "SET"
	ModifierActionAdd ModifierAction = "ADD"
	ModifierActionDel ModifierAction = "DEL"
	ModifierActionRen ModifierAction = "REN"
)
const (
	ContentTypeJson ContentType = "JSON"
	ContentTypeXml  ContentType = "XML"
	ContentTypeYml  ContentType = "YML"
	ContentTypeText ContentType = "TEXT"
)

// ContentTypeFromString converts a string representation of a content type
// to its corresponding ContentType value. It checks if the given string
// contains the string representation of ContentTypeJson or ContentTypeText case-insensitively.
// If the string contains ContentTypeJson, it returns ContentTypeJson.
// If the string contains ContentTypeText, it returns ContentTypeText.
// Otherwise, it returns an empty string.
// This function is used to convert a string content type to the ContentType enumeration value.
func ContentTypeFromString(s string) ContentType {
	// todo: aqui podemos ter XML, YAML, form-data
	if helper.ContainsIgnoreCase(s, ContentTypeJson.String()) {
		return ContentTypeJson
	} else if helper.ContainsIgnoreCase(s, ContentTypeText.String()) {
		return ContentTypeText
	}
	return ""
}

// IsEnumValid checks if the ModifierScope is a valid enumeration value.
// It returns true if the ModifierScope is either ModifierScopeRequest or ModifierScopeResponse,
// otherwise it returns false.
func (m ModifierScope) IsEnumValid() bool {
	switch m {
	case ModifierScopeRequest, ModifierScopeResponse:
		return true
	}
	return false
}

// IsEnumValid checks if the ModifierContext is a valid enumeration value.
// It returns true if the ModifierContext is either ModifierContextRequest or ModifierContextResponse,
// otherwise it returns false.
func (m ModifierContext) IsEnumValid() bool {
	switch m {
	case ModifierContextRequest, ModifierContextResponse:
		return true
	}
	return false
}

// IsEnumValid checks if the ModifierAction is a valid enumeration value.
// It returns true if the ModifierAction is either ModifierActionSet, ModifierActionAdd,
// ModifierActionDel, ModifierActionReplace, or ModifierActionRen, otherwise it returns false.
func (m ModifierAction) IsEnumValid() bool {
	switch m {
	case ModifierActionSet, ModifierActionAdd, ModifierActionDel, ModifierActionRen:
		return true
	}
	return false
}

// IsEnumValid checks if the CacheControl is a valid enumeration value.
// It returns true if the CacheControl is either CacheControlNoCache or CacheControlNoStore,
// otherwise it returns false.
func (c CacheControl) IsEnumValid() bool {
	switch c {
	case CacheControlNoCache, CacheControlNoStore:
		return true
	}
	return false
}

// IsEnumValid checks if the ResponseEncode is a valid enumeration value.
// It returns true if the ResponseEncode is either ResponseEncodeText,
// ResponseEncodeJson or ResponseEncodeXml, otherwise it returns false.
func (r ResponseEncode) IsEnumValid() bool {
	switch r {
	case ResponseEncodeText, ResponseEncodeJson, ResponseEncodeXml:
		return true
	}
	return false
}

// IsEnumValid checks if the ContentType is a valid enumeration value.
// It returns true if the ContentType is either ContentTypeText, ContentTypeJson,
// ContentTypeXml, or ContentTypeYml, otherwise it returns false.
func (c ContentType) IsEnumValid() bool {
	switch c {
	case ContentTypeText, ContentTypeJson, ContentTypeXml, ContentTypeYml:
		return true
	}
	return false
}

// ContentType returns the format of the content based on the ResponseEncode value.
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

// String returns the string representation of the ContentType value.
// It returns "application/json" if c is ContentTypeJson, "application/xml" if c is ContentTypeXml,
// "application/x-yaml" if c is ContentTypeYml, and "text/plain" for any other value of c.
// This method is used to convert the ContentType value to its corresponding MIME type string representation.
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
