package enum

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/enum"
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

type MiddlewareType string

type BackendResponseApply string

type ProjectionType int

type ProjectionValue int

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
	case NomenclatureCamel, NomenclatureLowerCamel, NomenclatureSnake, NomenclatureKebab:
		return true
	}
	return false
}

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
func (r Encode) ContentType() enum.ContentType {
	switch r {
	case EncodeText:
		return enum.ContentTypeText
	case EncodeJson:
		return enum.ContentTypeJson
	case EncodeXml:
		return enum.ContentTypeXml
	case EncodeYaml:
		return enum.ContentTypeYml
	}
	return ""
}

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
