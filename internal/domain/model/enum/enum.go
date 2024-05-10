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

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/iancoleman/strcase"
)

// ModifierScope represents the scope of a modification in the Backend or Endpoint.
type ModifierScope string

// ModifierAction represents the action to be performed in the Modifier struct.
type ModifierAction string

// Encode represents the encoding format for the API endpoint response.
type Encode string

// Nomenclature represents the case format for text values.
type Nomenclature string

// MiddlewareType represents the type of middleware in the Gopen application.
type MiddlewareType string

// BackendResponseApply represents the scope of applying a BackendResponse.
// It is used to indicate whether the response should be applied early or late.
// The possible values are "EARLY" and "LATE".
type BackendResponseApply string

// ProjectionType represents the type of projection in the Backend or Endpoint.
type ProjectionType int

// ProjectionValue represents the value of a projection in the Backend or Endpoint.
type ProjectionValue int

// CacheControl represents the header value of cache control.
type CacheControl string

// ContentType represents the format of the content.
type ContentType string

// ContentEncoding represents the encoding used for content.
type ContentEncoding string

const (
	ModifierScopeRequest  ModifierScope = "REQUEST"
	ModifierScopeResponse ModifierScope = "RESPONSE"
)
const (
	ModifierActionAdd ModifierAction = "ADD"
	ModifierActionApd ModifierAction = "APD"
	ModifierActionSet ModifierAction = "SET"
	ModifierActionRpl ModifierAction = "RPL"
	ModifierActionDel ModifierAction = "DEL"
)
const (
	ProjectionTypeAll       ProjectionType = iota
	ProjectionTypeAddition  ProjectionType = iota
	ProjectionTypeRejection ProjectionType = iota
)
const (
	ProjectionValueAddition  ProjectionValue = 1
	ProjectionValueRejection ProjectionValue = 0
)
const (
	Beforewares MiddlewareType = "beforewares"
	Afterwares  MiddlewareType = "afterwares"
)
const (
	NomenclatureCamel          Nomenclature = "CAMEL"
	NomenclatureLowerCamel     Nomenclature = "LOWER_CAMEL"
	NomenclatureSnake          Nomenclature = "SNAKE"
	NomenclatureScreamingSnake Nomenclature = "SCREAMING_SNAKE"
	NomenclatureKebab          Nomenclature = "KEBAB"
	NomenclatureScreamingKebab Nomenclature = "SCREAMING_KEBAB"
)
const (
	EncodeText Encode = "TEXT"
	EncodeJson Encode = "JSON"
	EncodeXml  Encode = "XML"
	EncodeYaml Encode = "YAML"
)
const (
	BackendResponseApplyEarly BackendResponseApply = "EARLY"
	BackendResponseApplyLate  BackendResponseApply = "LATE"
)
const (
	ContentEncodingGzip ContentEncoding = "gzip"
)
const (
	CacheControlNoCache CacheControl = "no-cache"
	CacheControlNoStore CacheControl = "no-store"
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
	} else if helper.ContainsIgnoreCase(s, ContentTypeXml.String()) {
		return ContentTypeXml
	} else if helper.ContainsIgnoreCase(s, ContentTypeText.String()) {
		return ContentTypeText
	}
	return ""
}

// ContentEncodingFromString converts a string representation of a content encoding
// to its corresponding ContentEncoding value. It checks if the given string
// contains ContentEncodingGzip case-insensitively. If the string contains
// ContentEncodingGzip, it returns ContentEncodingGzip. Otherwise, it returns an empty string.
// This function is used to convert a string content encoding to the ContentEncoding enumeration value.
func ContentEncodingFromString(s string) ContentEncoding {
	if helper.ContainsIgnoreCase(s, ContentEncodingGzip) {
		return ContentEncodingGzip
	}
	return ""
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
func (n Nomenclature) IsEnumValid() bool {
	switch n {
	case NomenclatureCamel, NomenclatureLowerCamel, NomenclatureSnake, NomenclatureKebab, NomenclatureScreamingSnake,
		NomenclatureScreamingKebab:
		return true
	}
	return false
}

// IsEnumValid checks if the BackendResponseApply is a valid enumeration value.
// It returns true if the BackendResponseApply is either BackendResponseApplyEarly
// or BackendResponseApplyLate, otherwise it returns false.
func (b BackendResponseApply) IsEnumValid() bool {
	switch b {
	case BackendResponseApplyEarly, BackendResponseApplyLate:
		return true
	}
	return false
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

// IsEnumValid checks if the ModifierAction is a valid enumeration value.
// It returns true if the ModifierAction is either ModifierActionSet, ModifierActionApd,
// ModifierActionRpl, ModifierActionAdd, ModifierActionDel, or ModifierActionRen,
// otherwise it returns false.
func (m ModifierAction) IsEnumValid() bool {
	switch m {
	case ModifierActionSet, ModifierActionApd, ModifierActionRpl, ModifierActionAdd, ModifierActionDel:
		return true
	}
	return false
}

// ContentType returns the format of the content based on the Encode value.
func (r Encode) ContentType() ContentType {
	switch r {
	case EncodeText:
		return ContentTypeText
	case EncodeJson:
		return ContentTypeJson
	case EncodeXml:
		return ContentTypeXml
	case EncodeYaml:
		return ContentTypeYml
	}
	return ""
}

// Parse takes a key string and converts it to the specified case format based on the value of Nomenclature.
// It returns the key string converted to the specified case format as per the Nomenclature value.
// If the Nomenclature value is not one of the predefined cases, it returns the key as is.
func (n Nomenclature) Parse(key string) string {
	switch n {
	case NomenclatureCamel:
		return strcase.ToCamel(key)
	case NomenclatureLowerCamel:
		return strcase.ToLowerCamel(key)
	case NomenclatureSnake:
		return strcase.ToSnake(key)
	case NomenclatureScreamingSnake:
		return strcase.ToScreamingSnake(key)
	case NomenclatureKebab:
		return strcase.ToKebab(key)
	case NomenclatureScreamingKebab:
		return strcase.ToScreamingKebab(key)
	default:
		return key
	}
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

// IsEnumValid checks if the ContentEncoding is a valid enumeration value.
// It returns true if the ContentEncoding is ContentEncodingGzip,
// otherwise it returns false.
func (c ContentEncoding) IsEnumValid() bool {
	switch c {
	case ContentEncodingGzip:
		return true
	}
	return false
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
