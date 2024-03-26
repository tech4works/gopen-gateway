package mapper

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"time"
)

func BuildLimiterDTO(limiterDefaultDTO, limiterDTO dto.Limiter) dto.Limiter {
	// obtemos os valores padrões informados na raiz do DTO Gopen
	maxHeaderSize := limiterDefaultDTO.MaxHeaderSize
	maxBodySize := limiterDefaultDTO.MaxBodySize
	maxMultipartMemorySize := limiterDefaultDTO.MaxMultipartMemorySize

	capacityRate := limiterDefaultDTO.Rate.Capacity
	everyRate := limiterDefaultDTO.Rate.Every

	// se o limiter do parâmetro vindo do endpoint DTO, tem valores, damos prioridade ao mesmo
	if helper.IsGreaterThan(limiterDTO.MaxHeaderSize, 0) {
		maxHeaderSize = limiterDTO.MaxHeaderSize
	}
	if helper.IsGreaterThan(limiterDTO.MaxBodySize, 0) {
		maxBodySize = limiterDTO.MaxBodySize
	}
	if helper.IsGreaterThan(limiterDTO.MaxMultipartMemorySize, 0) {
		maxMultipartMemorySize = limiterDTO.MaxMultipartMemorySize
	}

	// para o rate também
	if helper.IsGreaterThan(limiterDTO.Rate.Capacity, 0) {
		capacityRate = limiterDTO.Rate.Capacity
	}
	if helper.IsGreaterThan(limiterDTO.Rate.Every, 0) {
		everyRate = limiterDTO.Rate.Every
	}

	// caso o maxHeaderSize não seja informado, setamos um valor padrão de 1MB
	if helper.IsLessThanOrEqual(maxHeaderSize, 0) {
		maxHeaderSize = vo.Bytes(helper.SimpleConvertByteUnitStrToFloat("1MB"))
	}
	// caso o maxBodySize não seja informado, setamos um valor padrão de 3MB
	if helper.IsLessThanOrEqual(maxBodySize, 0) {
		maxBodySize = vo.Bytes(helper.SimpleConvertByteUnitStrToFloat("3MB"))
	}
	// caso o maxMultipartMemorySize não seja informado, setamos um valor padrão de 5MB
	if helper.IsLessThanOrEqual(maxMultipartMemorySize, 0) {
		maxMultipartMemorySize = vo.Bytes(helper.SimpleConvertByteUnitStrToFloat("5MB"))
	}

	// caso o maxRate não seja informado, setamos o valor 1 como padrão
	if helper.IsLessThanOrEqual(capacityRate, 0) {
		capacityRate = 1
	}
	// caso o everyRate não seja informado, setamos o valor 1s como padrão
	if helper.IsLessThanOrEqual(everyRate, 0) {
		everyRate = time.Second
	}

	// montamos o limiter DTO com os dados preenchidos
	return dto.Limiter{
		MaxHeaderSize:          maxHeaderSize,
		MaxBodySize:            maxBodySize,
		MaxMultipartMemorySize: maxMultipartMemorySize,
		Rate: dto.Rate{
			Capacity: capacityRate,
			Every:    everyRate,
		},
	}
}
