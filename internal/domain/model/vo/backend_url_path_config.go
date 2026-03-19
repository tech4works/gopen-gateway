package vo

import (
	"github.com/tech4works/checker"
)

type URLPathConfig struct {
	modifiers []ModifierConfig
}

func NewURLPathConfig(modifiers []ModifierConfig) *URLPathConfig {
	return &URLPathConfig{modifiers: modifiers}
}

func (b URLPathConfig) HasModifiers() bool {
	return checker.IsNotEmpty(b.modifiers)
}

func (b URLPathConfig) Modifiers() []ModifierConfig {
	return b.modifiers
}

func (b URLPathConfig) CountDataTransforms() int {
	if b.HasModifiers() {
		return len(b.Modifiers())
	}
	return 0
}
