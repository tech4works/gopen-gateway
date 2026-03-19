package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type PayloadConfig struct {
	aggregate       bool
	omit            bool
	omitEmpty       bool
	group           string
	contentType     ContentType
	contentEncoding ContentEncoding
	nomenclature    enum.Nomenclature
	mapper          *MapperConfig
	projector       *ProjectorConfig
	modifiers       []ModifierConfig
	joins           []JoinConfig
}

func NewPayloadConfig(
	aggregate bool,
	omit bool,
	omitEmpty bool,
	group,
	contentType,
	contentEncoding string,
	nomenclature enum.Nomenclature,
	mapper *MapperConfig,
	projector *ProjectorConfig,
	modifiers []ModifierConfig,
	joins []JoinConfig,
) *PayloadConfig {
	return &PayloadConfig{
		aggregate:       aggregate,
		omit:            omit,
		omitEmpty:       omitEmpty,
		group:           group,
		contentType:     NewContentType(contentType),
		contentEncoding: NewContentEncoding(contentEncoding),
		nomenclature:    nomenclature,
		mapper:          mapper,
		projector:       projector,
		modifiers:       modifiers,
		joins:           joins,
	}
}

func (b PayloadConfig) Aggregate() bool {
	return b.aggregate
}

func (b PayloadConfig) Omit() bool {
	return b.omit
}

func (b PayloadConfig) OmitEmpty() bool {
	return b.omitEmpty
}

func (b PayloadConfig) HasGroup() bool {
	return checker.IsNotEmpty(b.group)
}

func (b PayloadConfig) Group() string {
	return b.group
}

func (b PayloadConfig) HasContentType() bool {
	return b.contentType.IsSupported()
}

func (b PayloadConfig) HasContentEncoding() bool {
	return b.contentEncoding.IsSupported()
}

func (b PayloadConfig) ContentType() ContentType {
	return b.contentType
}

func (b PayloadConfig) ContentEncoding() ContentEncoding {
	return b.contentEncoding
}

func (b PayloadConfig) HasNomenclature() bool {
	return b.nomenclature.IsEnumValid()
}

func (b PayloadConfig) Nomenclature() enum.Nomenclature {
	return b.nomenclature
}

func (b PayloadConfig) Projector() *ProjectorConfig {
	return b.projector
}

func (b PayloadConfig) Mapper() *MapperConfig {
	return b.mapper
}

func (b PayloadConfig) Modifiers() []ModifierConfig {
	return b.modifiers
}

func (b PayloadConfig) Joins() []JoinConfig {
	return b.joins
}

func (b PayloadConfig) HasMapper() bool {
	return checker.NonNil(b.mapper)
}

func (b PayloadConfig) HasProjector() bool {
	return checker.NonNil(b.projector)
}

func (b PayloadConfig) HasModifiers() bool {
	return checker.IsNotEmpty(b.modifiers)
}

func (b PayloadConfig) HasJoins() bool {
	return checker.IsNotEmpty(b.joins)
}

func (b PayloadConfig) CountDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	if b.HasMapper() {
		count += len(b.Mapper().Map().Keys())
	}
	if b.HasProjector() {
		count += len(b.Projector().Project().Keys())
	}
	if b.HasModifiers() {
		count += len(b.Modifiers())
	}
	if b.HasJoins() {
		count += len(b.Joins())
	}
	return count
}
