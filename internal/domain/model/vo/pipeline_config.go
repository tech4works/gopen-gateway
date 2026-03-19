package vo

import (
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type OmitSpec interface {
	Omit() bool
}

type OmitEmptySpec interface {
	OmitEmpty() bool
}

type MapperSpec interface {
	Mapper() *MapperConfig
}

type ProjectorSpec interface {
	Projector() *ProjectorConfig
}

type ModifierSpec interface {
	Modifiers() []ModifierConfig
}

type JoinSpec interface {
	Joins() []JoinConfig
}

type GroupSpec interface {
	HasGroup() bool
	Group() string
}

type NomenclatureSpec interface {
	Nomenclature() enum.Nomenclature
}

type ContentTypeSpec interface {
	ContentType() ContentType
}

type ContentEncodingSpec interface {
	ContentEncoding() ContentEncoding
}

type HostsSpec interface {
	Hosts() []string
}

type GroupIDSpec interface {
	HasGroupID() bool
	GroupID() string
}

type DeduplicationIDSpec interface {
	HasDeduplicationID() bool
	DeduplicationID() string
}

type AttributesSpec interface {
	Attributes() map[string]AttributeValueConfig
}

type HostPipelineSpec interface {
	HostsSpec
}

type MetadataPipelineSpec interface {
	OmitSpec
	MapperSpec
	ProjectorSpec
	ModifierSpec
}

type URLPathPipelineSpec interface {
	ModifierSpec
}

type QueryPipelineSpec interface {
	OmitSpec
	MapperSpec
	ProjectorSpec
	ModifierSpec
}

type PayloadPipelineSpec interface {
	OmitSpec
	OmitEmptySpec
	GroupSpec
	NomenclatureSpec
	ContentTypeSpec
	ContentEncodingSpec
	MapperSpec
	ProjectorSpec
	ModifierSpec
	JoinSpec
}

type GroupIDPipelineSpec interface {
	GroupIDSpec
}

type DeduplicationIDPipelineSpec interface {
	DeduplicationIDSpec
}

type AttributesPipelineSpec interface {
	AttributesSpec
}
