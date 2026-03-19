package vo

import (
	"github.com/tech4works/checker"
)

type MetadataConfig struct {
	omit      bool
	mapper    *MapperConfig
	projector *ProjectorConfig
	modifiers []ModifierConfig
}

func NewMetadataConfig(
	omit bool,
	mapper *MapperConfig,
	projector *ProjectorConfig,
	modifiers []ModifierConfig,
) *MetadataConfig {
	return &MetadataConfig{
		omit:      omit,
		mapper:    mapper,
		projector: projector,
		modifiers: modifiers,
	}
}

func (c MetadataConfig) Omit() bool {
	return c.omit
}

func (c MetadataConfig) HasMapper() bool {
	return checker.NonNil(c.mapper)
}

func (c MetadataConfig) Mapper() *MapperConfig {
	return c.mapper
}

func (c MetadataConfig) Projector() *ProjectorConfig {
	return c.projector
}

func (c MetadataConfig) HasProjector() bool {
	return checker.NonNil(c.projector)
}

func (c MetadataConfig) Modifiers() []ModifierConfig {
	return c.modifiers
}

func (c MetadataConfig) HasModifiers() bool {
	return checker.IsNotEmpty(c.modifiers)
}

func (c MetadataConfig) CountDataTransforms() (count int) {
	if c.Omit() {
		return 1
	}
	if c.HasMapper() {
		count += len(c.Mapper().Map().Keys())
	}
	if c.HasProjector() {
		count += len(c.Projector().Project().Keys())
	}
	if c.HasModifiers() {
		count += len(c.Modifiers())
	}
	return count
}
