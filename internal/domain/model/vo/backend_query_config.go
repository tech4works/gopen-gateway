package vo

import (
	"github.com/tech4works/checker"
)

type QueryConfig struct {
	omit      bool
	mapper    *MapperConfig
	projector *ProjectorConfig
	modifiers []ModifierConfig
}

func NewQueryConfig(
	omit bool,
	mapper *MapperConfig,
	projector *ProjectorConfig,
	modifiers []ModifierConfig,
) *QueryConfig {
	return &QueryConfig{
		omit:      omit,
		mapper:    mapper,
		projector: projector,
		modifiers: modifiers,
	}
}

func (b QueryConfig) Omit() bool {
	return b.omit
}

func (b QueryConfig) Projector() *ProjectorConfig {
	return b.projector
}

func (b QueryConfig) Modifiers() []ModifierConfig {
	return b.modifiers
}

func (b QueryConfig) Mapper() *MapperConfig {
	return b.mapper
}

func (b QueryConfig) CountDataTransforms() (count int) {
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
	return count
}

func (b QueryConfig) HasMapper() bool {
	return checker.NonNil(b.mapper)
}

func (b QueryConfig) HasProjector() bool {
	return checker.NonNil(b.projector)
}

func (b QueryConfig) HasModifiers() bool {
	return checker.IsNotEmpty(b.modifiers)
}
