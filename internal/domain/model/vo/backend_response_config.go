package vo

import (
	"github.com/tech4works/checker"
)

type BackendResponseConfig struct {
	omit     bool
	metadata *MetadataConfig
	payload  *PayloadConfig
}

func NewBackendResponseConfigForMiddleware(omit bool, metadata *MetadataConfig) *BackendResponseConfig {
	return &BackendResponseConfig{
		omit:     omit,
		metadata: metadata,
		payload: &PayloadConfig{
			omit: true,
		},
	}
}

func NewBackendResponseConfig(omit bool, metadata *MetadataConfig, payload *PayloadConfig,
) *BackendResponseConfig {
	return &BackendResponseConfig{
		omit:     omit,
		metadata: metadata,
		payload:  payload,
	}
}

func (b BackendResponseConfig) Omit() bool {
	return b.omit
}

func (b BackendResponseConfig) CountAllDataTransforms() (count int) {
	if b.Omit() {
		return 1
	}
	if b.HasMetadata() {
		count += b.Metadata().CountDataTransforms()
	}
	if b.HasPayload() {
		count += b.Payload().CountDataTransforms()
	}
	return count
}

func (b BackendResponseConfig) HasMetadata() bool {
	return checker.NonNil(b.metadata)
}

func (b BackendResponseConfig) Metadata() *MetadataConfig {
	return b.metadata
}

func (b BackendResponseConfig) HasPayload() bool {
	return checker.NonNil(b.payload)
}

func (b BackendResponseConfig) Payload() *PayloadConfig {
	return b.payload
}
