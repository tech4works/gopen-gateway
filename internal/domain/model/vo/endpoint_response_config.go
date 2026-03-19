package vo

import "github.com/tech4works/checker"

type EndpointResponseConfig struct {
	metadata *MetadataConfig
	payload  *PayloadConfig
}

func NewEndpointResponseConfig(metadata *MetadataConfig, payload *PayloadConfig) EndpointResponseConfig {
	return EndpointResponseConfig{
		metadata: metadata,
		payload:  payload,
	}
}

func (e EndpointResponseConfig) HasMetadata() bool {
	return checker.NonNil(e.metadata)
}

func (e EndpointResponseConfig) Metadata() *MetadataConfig {
	return e.metadata
}

func (e EndpointResponseConfig) HasPayload() bool {
	return checker.NonNil(e.payload)
}

func (e EndpointResponseConfig) Payload() *PayloadConfig {
	return e.payload
}
