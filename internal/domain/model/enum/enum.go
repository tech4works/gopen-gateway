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
	"github.com/iancoleman/strcase"
)

// ModifierScope represents the scope of a modification in the Backend or Endpoint.
type ModifierScope string

// ModifierAction represents the action to be performed in the Modifier struct.
type ModifierAction string

// ContentType represents the encoding format. It is a string type.
// The valid values for ContentType are "PLAIN_TEXT", "JSON", and "XML".
type ContentType string

// ContentEncoding represents the encoding format.
type ContentEncoding string

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
	ContentTypePlainText ContentType = "PLAIN_TEXT"
	ContentTypeJson      ContentType = "JSON"
	ContentTypeXml       ContentType = "XML"
)
const (
	ContentEncodingNone    ContentEncoding = "NONE"
	ContentEncodingGzip    ContentEncoding = "GZIP"
	ContentEncodingDeflate ContentEncoding = "DEFLATE"
)
const (
	BackendResponseApplyEarly BackendResponseApply = "EARLY"
	BackendResponseApplyLate  BackendResponseApply = "LATE"
)
const (
	CacheControlNoCache CacheControl = "no-cache"
	CacheControlNoStore CacheControl = "no-store"
)

// IsEnumValid checks if the ContentType is a valid enumeration value.
// It returns true if the ContentType is either ContentTypePlainText,
// ContentTypeJson or ContentTypeXml, otherwise it returns false.
func (r ContentType) IsEnumValid() bool {
	switch r {
	case ContentTypePlainText, ContentTypeJson, ContentTypeXml:
		return true
	}
	return false
}

// IsEnumValid checks if the ContentEncoding is a valid enumeration value.
// It returns true if the ContentEncoding is either ContentEncodingNone,
// ContentEncodingGzip or ContentEncodingDeflate, otherwise it returns false.
func (r ContentEncoding) IsEnumValid() bool {
	switch r {
	case ContentEncodingNone, ContentEncodingGzip, ContentEncodingDeflate:
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
