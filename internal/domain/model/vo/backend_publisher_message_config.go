package vo

import "github.com/tech4works/checker"

type BackendPublisherMessageConfig struct {
	onlyIf     []string
	ignoreIf   []string
	attributes map[string]AttributeValueConfig
	body       *PayloadConfig
}

func NewBackendPublisherMessageConfig(
	onlyIf,
	ignoreIf []string,
	attributes map[string]AttributeValueConfig,
	body *PayloadConfig,
) BackendPublisherMessageConfig {
	return BackendPublisherMessageConfig{
		onlyIf:     onlyIf,
		ignoreIf:   ignoreIf,
		attributes: attributes,
		body:       body,
	}
}

func (m BackendPublisherMessageConfig) HasBody() bool {
	return checker.NonNil(m.body)
}

func (m BackendPublisherMessageConfig) Body() *PayloadConfig {
	return m.body
}

func (m BackendPublisherMessageConfig) HasAttributes() bool {
	return checker.IsNotNilOrEmpty(m.attributes)
}

func (m BackendPublisherMessageConfig) Attributes() map[string]AttributeValueConfig {
	return m.attributes
}
