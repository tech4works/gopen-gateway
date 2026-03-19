package vo

import (
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type EndpointExecutionConfig struct {
	BaseExecutionConfig
}

func NewEndpointExecutionConfigDefault() EndpointExecutionConfig {
	return EndpointExecutionConfig{
		BaseExecutionConfig: NewBaseExecutionConfig(enum.ExecutionModeBestEffort, nil),
	}
}

func NewEndpointExecutionConfig(mode enum.ExecutionMode, on []enum.ExecutionOn) EndpointExecutionConfig {
	return EndpointExecutionConfig{
		BaseExecutionConfig: NewBaseExecutionConfig(mode, on),
	}
}
