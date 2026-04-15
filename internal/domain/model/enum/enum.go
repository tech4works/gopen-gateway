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

import "google.golang.org/grpc/codes"

type Protocol string

type ExecutionPolicy string

type ModifierAction string

type Nomenclature string

type BackendFlow string

type BackendKind string

type ProjectKind int

type ProjectValue int

type BackendBroker string

type TemplateMerge string

type ComponentMerge string

type JoinTargetOnMissing string

type JoinTargetPolicy string

type MapperPolicy string

type BackendOutcome string

type ExecutionMode string

type ExecutionOn string

type AttributeValueType string

type ResponseStatus string

type DegradationKind string

type CacheKind string

const (
	ProtocolHTTP      Protocol = "HTTP"
	ProtocolGRPC      Protocol = "GRPC"
	ProtocolWebSocket Protocol = "WEBSOCKET"
)
const (
	BackendOutcomeExecuted  BackendOutcome = "EXECUTED"
	BackendOutcomeIgnored   BackendOutcome = "IGNORED"
	BackendOutcomeCancelled BackendOutcome = "CANCELLED"
	BackendOutcomeError     BackendOutcome = "ERROR"
)
const (
	ModifierActionAdd     ModifierAction = "ADD"
	ModifierActionAppend  ModifierAction = "APPEND"
	ModifierActionSet     ModifierAction = "SET"
	ModifierActionReplace ModifierAction = "REPLACE"
	ModifierActionDelete  ModifierAction = "DELETE"
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
	BackendBrokerAwsSqs BackendBroker = "AWS/SQS"
	BackendBrokerAwsSns BackendBroker = "AWS/SNS"
)
const (
	TemplateMergeBase TemplateMerge = "BASE"
	TemplateMergeFull TemplateMerge = "FULL"
)
const (
	ComponentMergeExtend   ComponentMerge = "EXTEND"
	ComponentMergeOverride ComponentMerge = "OVERRIDE"
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
const (
	ExecutionModeFailFast   ExecutionMode = "FAIL_FAST"
	ExecutionModeBestEffort ExecutionMode = "BEST_EFFORT"
)
const (
	ExecutionOnBuild       ExecutionOn = "BUILD"
	ExecutionOnClientError ExecutionOn = "CLIENT_ERROR"
	ExecutionOnServerError ExecutionOn = "SERVER_ERROR"
)
const (
	AttributeValueTypeString AttributeValueType = "STRING"
	AttributeValueTypeNumber AttributeValueType = "NUMBER"
	AttributeValueTypeBinary AttributeValueType = "BINARY"
)
const (
	ResponseStatusUnknown            ResponseStatus = "UNKNOWN"
	ResponseStatusOK                 ResponseStatus = "OK"
	ResponseStatusCancelled          ResponseStatus = "CANCELLED"
	ResponseStatusInvalidArgument    ResponseStatus = "INVALID_ARGUMENT"
	ResponseStatusDeadlineExceeded   ResponseStatus = "DEADLINE_EXCEEDED"
	ResponseStatusNotFound           ResponseStatus = "NOT_FOUND"
	ResponseStatusAlreadyExists      ResponseStatus = "ALREADY_EXISTS"
	ResponseStatusPermissionDenied   ResponseStatus = "PERMISSION_DENIED"
	ResponseStatusUnauthenticated    ResponseStatus = "UNAUTHENTICATED"
	ResponseStatusResourceExhausted  ResponseStatus = "RESOURCE_EXHAUSTED"
	ResponseStatusPayloadTooLarge    ResponseStatus = "PAYLOAD_TOO_LARGE"
	ResponseStatusMetadataTooLarge   ResponseStatus = "METADATA_TOO_LARGE"
	ResponseStatusFailedPrecondition ResponseStatus = "FAILED_PRECONDITION"
	ResponseStatusAborted            ResponseStatus = "ABORTED"
	ResponseStatusOutOfRange         ResponseStatus = "OUT_OF_RANGE"
	ResponseStatusUnimplemented      ResponseStatus = "UNIMPLEMENTED"
	ResponseStatusInternalError      ResponseStatus = "INTERNAL_ERROR"
	ResponseStatusUnavailable        ResponseStatus = "UNAVAILABLE"
	ResponseStatusBadGateway         ResponseStatus = "BAD_GATEWAY"
	ResponseStatusDataLoss           ResponseStatus = "DATA_LOSS"
	ResponseStatusConflict           ResponseStatus = "CONFLICT"
)
const (
	DegradationKindMetadata        DegradationKind = "METADATA"
	DegradationKindQuery           DegradationKind = "QUERY"
	DegradationKindURLPath         DegradationKind = "URL_PATH"
	DegradationKindPayload         DegradationKind = "PAYLOAD"
	DegradationKindDeduplicationID DegradationKind = "DEDUPLICATION_ID"
	DegradationKindGroupID         DegradationKind = "GROUP_ID"
	DegradationKindAttributes      DegradationKind = "ATTRIBUTES"
)
const (
	CacheKindEndpoint CacheKind = "ENDPOINT"
	CacheKindBackend  CacheKind = "BACKEND"
)

func NewResponseStatusFromGRPC(code codes.Code) ResponseStatus {
	switch code {
	case codes.OK:
		return ResponseStatusOK
	case codes.Canceled:
		return ResponseStatusCancelled
	case codes.InvalidArgument:
		return ResponseStatusInvalidArgument
	case codes.DeadlineExceeded:
		return ResponseStatusDeadlineExceeded
	case codes.NotFound:
		return ResponseStatusNotFound
	case codes.AlreadyExists:
		return ResponseStatusAlreadyExists
	case codes.PermissionDenied:
		return ResponseStatusPermissionDenied
	case codes.Unauthenticated:
		return ResponseStatusUnauthenticated
	case codes.ResourceExhausted:
		return ResponseStatusResourceExhausted
	case codes.FailedPrecondition:
		return ResponseStatusFailedPrecondition
	case codes.Aborted:
		return ResponseStatusAborted
	case codes.OutOfRange:
		return ResponseStatusOutOfRange
	case codes.Unimplemented:
		return ResponseStatusUnimplemented
	case codes.Internal:
		return ResponseStatusInternalError
	case codes.Unavailable:
		return ResponseStatusUnavailable
	case codes.DataLoss:
		return ResponseStatusDataLoss
	default:
		return ResponseStatusUnknown
	}
}

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

func (b BackendOutcome) IsEnumValid() bool {
	switch b {
	case BackendOutcomeExecuted, BackendOutcomeIgnored, BackendOutcomeCancelled, BackendOutcomeError:
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
	case ModifierActionSet, ModifierActionAppend, ModifierActionReplace, ModifierActionAdd, ModifierActionDelete:
		return true
	}
	return false
}

func (d DegradationKind) IsEnumValid() bool {
	switch d {
	case DegradationKindMetadata, DegradationKindQuery, DegradationKindURLPath, DegradationKindPayload,
		DegradationKindDeduplicationID, DegradationKindGroupID, DegradationKindAttributes:
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

func (c ComponentMerge) IsEnumValid() bool {
	switch c {
	case ComponentMergeExtend, ComponentMergeOverride:
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

func (m ExecutionMode) IsEnumValid() bool {
	switch m {
	case ExecutionModeFailFast, ExecutionModeBestEffort:
		return true
	}
	return false
}

func (o ExecutionOn) IsEnumValid() bool {
	switch o {
	case ExecutionOnBuild, ExecutionOnClientError, ExecutionOnServerError:
		return true
	}
	return false
}

func (p AttributeValueType) IsEnumValid() bool {
	switch p {
	case AttributeValueTypeString, AttributeValueTypeNumber,
		AttributeValueTypeBinary:
		return true
	}
	return false
}

func (p AttributeValueType) String() string {
	return string(p)
}

func (p Protocol) IsEnumValid() bool {
	switch p {
	case ProtocolHTTP, ProtocolGRPC, ProtocolWebSocket:
		return true
	}
	return false
}

func (c CacheKind) IsEnumValid() bool {
	switch c {
	case CacheKindEndpoint, CacheKindBackend:
		return true
	}
	return false
}
