package vo

import (
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type AttributeValueConfig struct {
	mType enum.AttributeValueType
	value string
}

func NewAttributeValueConfig(mType enum.AttributeValueType, value string,
) AttributeValueConfig {
	return AttributeValueConfig{
		mType: mType,
		value: value,
	}
}

func (a AttributeValueConfig) Type() enum.AttributeValueType {
	return a.mType
}

func (a AttributeValueConfig) Value() string {
	return a.value
}
