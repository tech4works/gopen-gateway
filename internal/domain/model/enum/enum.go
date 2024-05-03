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

// Encode represents the encoding format for the API endpoint response.
type Encode string

// Nomenclature represents the case format for text values.
type Nomenclature string

// ContentType represents the format of the content.
type ContentType string

type ContentEncoding string

type MiddlewareType string

type BackendResponseApply string

const (
	ContentEncodingUnzip ContentEncoding = "unzip"
)
const (
	BackendResponseApplyEarly BackendResponseApply = "EARLY"
	BackendResponseApplyLate  BackendResponseApply = "LATE"
)
const (
	Beforewares MiddlewareType = "beforewares"
	Afterwares  MiddlewareType = "afterwares"
)
const (
	NomenclatureCamel      Nomenclature = "CAMEL"
	NomenclatureLowerCamel Nomenclature = "LOWER_CAMEL"
	NomenclatureSnake      Nomenclature = "SNAKE"
	NomenclatureKebab      Nomenclature = "KEBAB"
)
const (
	EncodeText Encode = "TEXT"
	EncodeJson Encode = "JSON"
	EncodeXml  Encode = "XML"
	EncodeYaml Encode = "YAML"
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
	ModifierActionAdd ModifierAction = "ADD"
	ModifierActionApd ModifierAction = "APD"
	ModifierActionSet ModifierAction = "SET"
	ModifierActionRpl ModifierAction = "RPL"
	ModifierActionRen ModifierAction = "REN"
	ModifierActionDel ModifierAction = "DEL"
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

func ContentEncodingFromString(s string) ContentEncoding {
	if helper.ContainsIgnoreCase(s, ContentEncodingUnzip) {
		return ContentEncodingUnzip
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
// It returns true if the ModifierAction is either ModifierActionSet, ModifierActionApd,
// ModifierActionRpl, ModifierActionAdd, ModifierActionDel, or ModifierActionRen,
// otherwise it returns false.
func (m ModifierAction) IsEnumValid() bool {
	switch m {
	case ModifierActionSet, ModifierActionApd, ModifierActionRpl, ModifierActionAdd, ModifierActionDel, ModifierActionRen:
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

// IsEnumValid checks if the Encode is a valid enumeration value.
// It returns true if the Encode is either EncodeText,
// EncodeJson or EncodeXml, otherwise it returns false.
func (r Encode) IsEnumValid() bool {
	switch r {
	case EncodeText, EncodeJson, EncodeXml:
		return true
	}
	return false
}

// IsEnumValid checks if the Nomenclature is a valid enumeration value. Returns true if
// Nomenclature is either NomenclatureCamel, NomenclatureLowerCamel, NomenclatureSnake, or NomenclatureKebab,
// otherwise it returns false.
func (c Nomenclature) IsEnumValid() bool {
	switch c {
	case NomenclatureCamel, NomenclatureLowerCamel, NomenclatureSnake, NomenclatureKebab:
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

func (c ContentEncoding) IsEnumValid() bool {
	switch c {
	case ContentEncodingUnzip:
		return true
	}
	return false
}

// todo
func (b BackendResponseApply) IsEnumValid() bool {
	switch b {
	case BackendResponseApplyEarly, BackendResponseApplyLate:
		return true
	}
	return false
}

// ContentType returns the format of the content based on the Encode value.
func (r Encode) ContentType() ContentType {
	switch r {
	case EncodeJson:
		return ContentTypeJson
	case EncodeXml:
		return ContentTypeXml
	case EncodeYaml:
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
