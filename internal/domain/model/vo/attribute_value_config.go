package vo

import (
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type AttributeValueConfig struct {
	dataType enum.AttributeValueDataType
	value    string
}

func NewAttributeValueConfig(dataType enum.AttributeValueDataType, value string,
) AttributeValueConfig {
	return AttributeValueConfig{
		dataType: dataType,
		value:    value,
	}
}

func (a AttributeValueConfig) DataType() enum.AttributeValueDataType {
	return a.dataType
}

func (a AttributeValueConfig) Value() string {
	return a.value
}
