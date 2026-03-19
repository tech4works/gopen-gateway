package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type BackendExecutionConfig struct {
	BaseExecutionConfig
	concurrent int
	async      bool
}

func NewBackendExecutionConfigDefault() BackendExecutionConfig {
	return NewBackendExecutionConfigWithModeDefault(0, false)
}

func NewBackendExecutionConfigWithModeDefault(concurrent int, async bool) BackendExecutionConfig {
	return NewBackendExecutionConfig(concurrent, async, "", nil)
}

func NewBackendExecutionConfig(concurrent int, async bool, mode enum.ExecutionMode, on []enum.ExecutionOn,
) BackendExecutionConfig {
	return BackendExecutionConfig{
		BaseExecutionConfig: NewBaseExecutionConfig(mode, on),
		concurrent:          concurrent,
		async:               async,
	}
}

func (b BackendExecutionConfig) IsConcurrent() bool {
	return checker.IsGreaterThanOrEqual(b.concurrent, 2)
}

func (b BackendExecutionConfig) Concurrent() int {
	return b.concurrent
}

func (b BackendExecutionConfig) Async() bool {
	return b.async
}

func (b BackendExecutionConfig) ShouldAbortOnResponseStatus(responseStatus ResponseStatus) bool {
	if responseStatus.ClientError() {
		return b.ShouldAbortOn(enum.ExecutionOnClientError)
	} else if responseStatus.ServerError() {
		return b.ShouldAbortOn(enum.ExecutionOnServerError)
	} else {
		return false
	}
}
