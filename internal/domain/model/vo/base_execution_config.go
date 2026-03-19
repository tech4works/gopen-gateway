package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type BaseExecutionConfig struct {
	mode enum.ExecutionMode
	on   []enum.ExecutionOn
}

func NewBaseExecutionConfig(mode enum.ExecutionMode, on []enum.ExecutionOn) BaseExecutionConfig {
	if !mode.IsEnumValid() {
		mode = enum.ExecutionModeBestEffort
	}
	if checker.IsNilOrEmpty(on) {
		on = map[enum.ExecutionMode][]enum.ExecutionOn{
			enum.ExecutionModeBestEffort: {},
			enum.ExecutionModeFailFast:   {enum.ExecutionOnBuild, enum.ExecutionOnServerError},
		}[mode]
	}
	return BaseExecutionConfig{
		mode: mode,
		on:   on,
	}
}

func (b BaseExecutionConfig) HasOn(v enum.ExecutionOn) bool {
	return checker.Contains(b.on, v)
}

func (b BaseExecutionConfig) IsBestEffort() bool {
	return checker.Equals(b.mode, enum.ExecutionModeBestEffort)
}

func (b BaseExecutionConfig) IsFailFast() bool {
	return checker.Equals(b.mode, enum.ExecutionModeFailFast)
}

func (b BaseExecutionConfig) UseFallback(ev enum.ExecutionOn) bool {
	return b.IsFailFast() && b.ContinueOn(ev)
}

func (b BaseExecutionConfig) ContinueOn(ev enum.ExecutionOn) bool {
	return !b.ShouldAbortOn(ev)
}

func (b BaseExecutionConfig) ShouldAbortOn(ev enum.ExecutionOn) bool {
	if b.IsBestEffort() {
		return false
	} else if b.IsFailFast() {
		return b.HasOn(ev)
	}
	return false
}
