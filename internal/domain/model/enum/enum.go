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

type ExecutionPolicy string

type ModifierAction string

type Nomenclature string

type BackendFlow string

type BackendKind string

type ProjectKind int

type ProjectValue int

type ContentType string

type ContentEncoding string

type CacheControl string

type ProxyProvider string

type BackendBroker string

type TemplateMerge string

type JoinTargetOnMissing string

type JoinTargetPolicy string

type MapperPolicy string

type BackendOutcome string

const (
	ExecutionPolicyFailFastAll                  ExecutionPolicy = "FAIL_FAST_ALL"
	ExecutionPolicyFailFastStatus               ExecutionPolicy = "FAIL_FAST_STATUS"
	ExecutionPolicyFailFastInfra                ExecutionPolicy = "FAIL_FAST_INFRA"
	ExecutionPolicyFailFastInfraOnly            ExecutionPolicy = "FAIL_FAST_INFRA_ONLY"
	ExecutionPolicyBestEffortAll                ExecutionPolicy = "BEST_EFFORT_ALL"
	ExecutionPolicyBestEffortStatusFail         ExecutionPolicy = "BEST_EFFORT_STATUS_FAIL"
	ExecutionPolicyBestEffortInfraFail          ExecutionPolicy = "BEST_EFFORT_INFRA_FAIL"
	ExecutionPolicyBestEffortTransformFail      ExecutionPolicy = "BEST_EFFORT_TRANSFORM_FAIL"
	ExecutionPolicyBestEffortStatusAndInfraFail ExecutionPolicy = "BEST_EFFORT_STATUS_AND_INFRA_FAIL"
	ExecutionPolicyBestEffortStrict             ExecutionPolicy = "BEST_EFFORT_STRICT"
)
const (
	BackendOutcomeExecuted  BackendOutcome = "EXECUTED"
	BackendOutcomeIgnored   BackendOutcome = "IGNORED"
	BackendOutcomeCancelled BackendOutcome = "CANCELLED"
)
const (
	ModifierActionAdd ModifierAction = "ADD"
	ModifierActionApd ModifierAction = "APD"
	ModifierActionSet ModifierAction = "SET"
	ModifierActionRpl ModifierAction = "RPL"
	ModifierActionDel ModifierAction = "DEL"
)
const (
	ProjectKindAll       ProjectKind = iota
	ProjectKindAddition  ProjectKind = iota
	ProjectKindRejection ProjectKind = iota
)
const (
	ProjectValueAddition  ProjectValue = 1
	ProjectValueRejection ProjectValue = 0
)
const (
	BackendFlowNormal     BackendFlow = "NORMAL"
	BackendFlowBeforeware BackendFlow = "BEFOREWARE"
	BackendFlowAfterware  BackendFlow = "AFTERWARE"
)
const (
	BackendKindHTTP      BackendKind = "HTTP"
	BackendKindPublisher BackendKind = "PUBLISHER"
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
	CacheControlNoCache CacheControl = "no-cache"
	CacheControlNoStore CacheControl = "no-store"
)
const (
	ProxyProviderNgrok ProxyProvider = "NGROK"
)
const (
	BackendBrokerAwsSqs BackendBroker = "AWS/SQS"
	BackendBrokerAwsSns BackendBroker = "AWS/SNS"
)
const (
	TemplateMergeBase TemplateMerge = "BASE"
	TemplateMergeFull TemplateMerge = "FULL"
)
const (
	JoinTargetOnMissingNull  JoinTargetOnMissing = "NULL"
	JoinTargetOnMissingOmit  JoinTargetOnMissing = "OMIT"
	JoinTargetOnMissingError JoinTargetOnMissing = "ERROR"
)
const (
	JoinTargetKeepKey         JoinTargetPolicy = "KEEP_KEY"
	JoinTargetDropKeyOnMerged JoinTargetPolicy = "DROP_KEY_ON_MERGED"
	JoinTargetDropKeyAlways   JoinTargetPolicy = "DROP_KEY_ALWAYS"
)
const (
	MapperPolicyKeepUnmapped MapperPolicy = "KEEP_UNMAPPED"
	MapperPolicyDropUnmapped MapperPolicy = "DROP_UNMAPPED"
)

func (p BackendBroker) String() string {
	return string(p)
}

func (p BackendBroker) IsEnumValid() bool {
	switch p {
	case BackendBrokerAwsSqs, BackendBrokerAwsSns:
		return true
	}
	return false
}

func (e JoinTargetOnMissing) IsEnumValid() bool {
	switch e {
	case JoinTargetOnMissingNull, JoinTargetOnMissingOmit, JoinTargetOnMissingError:
		return true
	}
	return false
}

func (e JoinTargetPolicy) IsEnumValid() bool {
	switch e {
	case JoinTargetKeepKey, JoinTargetDropKeyOnMerged, JoinTargetDropKeyAlways:
		return true
	}
	return false
}

func (e MapperPolicy) IsEnumValid() bool {
	switch e {
	case MapperPolicyKeepUnmapped, MapperPolicyDropUnmapped:
		return true
	}
	return false
}

func (c ContentType) IsEnumValid() bool {
	switch c {
	case ContentTypePlainText, ContentTypeJson, ContentTypeXml:
		return true
	}
	return false
}

func (b BackendOutcome) IsEnumValid() bool {
	switch b {
	case BackendOutcomeExecuted, BackendOutcomeIgnored, BackendOutcomeCancelled:
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

func (m ModifierAction) IsEnumValid() bool {
	switch m {
	case ModifierActionSet, ModifierActionApd, ModifierActionRpl, ModifierActionAdd, ModifierActionDel:
		return true
	}
	return false
}

func (b BackendFlow) String() string {
	return string(b)
}

func (b BackendFlow) Abbreviation() string {
	switch b {
	case BackendFlowNormal:
		return "BKD"
	case BackendFlowBeforeware:
		return "BFW"
	case BackendFlowAfterware:
		return "AFW"
	}
	return ""
}

func (t TemplateMerge) IsEnumValid() bool {
	switch t {
	case TemplateMergeBase, TemplateMergeFull:
		return true
	}
	return false
}

func (b BackendKind) IsEnumValid() bool {
	switch b {
	case BackendKindHTTP, BackendKindPublisher:
		return true
	}
	return false
}
