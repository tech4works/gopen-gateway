package mapper

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

func BuildModifiersVO(modifiersDTO []dto.Modifier) (modifiersVO []vo.Modifier) {
	for _, modifierDTO := range modifiersDTO {
		modifiersVO = append(modifiersVO, BuildModifierVO(modifierDTO))
	}
	return modifiersVO
}

func BuildModifierVO(modifierDTO dto.Modifier) vo.Modifier {
	return vo.Modifier{
		Context: modifierDTO.Context,
		Scope:   modifierDTO.Scope,
		Action:  modifierDTO.Action,
		Key:     modifierDTO.Key,
		Value:   modifierDTO.Value,
	}
}
