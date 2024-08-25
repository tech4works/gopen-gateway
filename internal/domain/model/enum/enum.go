/*
 * Copyright 2024 Tech4Works
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

type ModifierScope string

type ModifierAction string

type Nomenclature string

type BackendType string

type BackendResponseApply string

type ProjectionType int

type ProjectionValue int

type ContentType string

type ContentEncoding string

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
	BackendTypeNormal     BackendType = "BACKEND"
	BackendTypeBeforeware BackendType = "BEFOREWARE"
	BackendTypeAfterware  BackendType = "AFTERWARE"
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
	BackendResponseApplyEarly BackendResponseApply = "EARLY"
	BackendResponseApplyLate  BackendResponseApply = "LATE"
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
	CacheControlNoCache CacheControl = "no-cache"
	CacheControlNoStore CacheControl = "no-store"
)

func (c ContentType) IsEnumValid() bool {
	switch c {
	case ContentTypePlainText, ContentTypeJson, ContentTypeXml:
		return true
	}
	return false
}

func (r ContentEncoding) IsEnumValid() bool {
	switch r {
	case ContentEncodingNone, ContentEncodingGzip, ContentEncodingDeflate:
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

func (n Nomenclature) IsEnumValid() bool {
	switch n {
	case NomenclatureCamel, NomenclatureLowerCamel, NomenclatureSnake, NomenclatureKebab, NomenclatureScreamingSnake,
		NomenclatureScreamingKebab:
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

func (m ModifierScope) IsEnumValid() bool {
	switch m {
	case ModifierScopeRequest, ModifierScopeResponse:
		return true
	}
	return false
}

func (m ModifierAction) IsEnumValid() bool {
	switch m {
	case ModifierActionSet, ModifierActionApd, ModifierActionRpl, ModifierActionAdd, ModifierActionDel:
		return true
	}
	return false
}

func (b BackendType) String() string {
	return string(b)
}

func (b BackendType) Abbreviation() string {
	switch b {
	case BackendTypeNormal:
		return "BKD"
	case BackendTypeBeforeware:
		return "BFW"
	case BackendTypeAfterware:
		return "AFW"
	}
	return ""
}
