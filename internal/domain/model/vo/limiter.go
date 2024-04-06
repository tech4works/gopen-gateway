package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"time"
)

type Limiter struct {
	maxHeaderSize          Bytes
	maxBodySize            Bytes
	maxMultipartMemorySize Bytes
	rate                   Rate
}

type Rate struct {
	capacity int
	every    time.Duration
}

// NewLimiterFromEndpoint takes a Gopen object and an Endpoint object and returns a Limiter object.
// It initializes the fields of the Limiter object based on values from the Gopen and Endpoint objects.
// If a value is provided in the Endpoint object, it takes priority over the value from the Gopen object.
func NewLimiterFromEndpoint(gopenVO Gopen, endpointVO Endpoint) Limiter {
	// por padrão obtemos o limiter.max-header-size configurado na raiz, caso não informado um valor padrão é retornado
	maxHeaderSize := gopenVO.LimiterMaxHeaderSize()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterMaxHeaderSize() {
		maxHeaderSize = endpointVO.LimiterMaxHeaderSize()
	}

	// por padrão obtemos o limiter.max-body-size configurado na raiz, caso não informado um valor padrão é retornado
	maxBodySize := gopenVO.LimiterMaxBodySize()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterMaxBodySize() {
		maxBodySize = endpointVO.LimiterMaxBodySize()
	}

	// por padrão obtemos o limiter.max-multipart-form-size configurado na raiz, caso não informado um valor padrão é retornado
	maxMultipartForm := gopenVO.LimiterMaxMultipartMemorySize()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterMaxMultipartFormSize() {
		maxMultipartForm = endpointVO.LimiterMaxMultipartMemorySize()
	}

	//construímos o limiter vo
	return Limiter{
		maxHeaderSize:          maxHeaderSize,
		maxBodySize:            maxBodySize,
		maxMultipartMemorySize: maxMultipartForm,
		rate:                   newRateFromEndpoint(gopenVO, endpointVO),
	}
}

// newRateFromEndpoint creates a new instance of Rate based on the provided Gopen and Endpoint objects.
// It retrieves the rate.every and rate.capacity values from the Gopen and Endpoint objects.
// If the rate.every value is provided in the Endpoint object, it takes priority over the value in the Gopen object.
// If the rate.capacity value is provided in the Endpoint object, it takes priority over the value in the Gopen object.
// The function returns a Rate object initialized with the capacity and every value.
func newRateFromEndpoint(gopenVO Gopen, endpointVO Endpoint) Rate {
	// por padrão obtemos o limiter.rate.every configurado na raiz, caso não informado um valor padrão é retornado
	every := gopenVO.LimiterRateEvery()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterRateEvery() {
		every = endpointVO.LimiterRateEvery()
	}

	// por padrão obtemos o limiter.rate.capacity configurado na raiz, caso não informado um valor padrão é retornado
	capacity := gopenVO.LimiterRateCapacity()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterRateCapacity() {
		capacity = endpointVO.LimiterRateCapacity()
	}

	// montamos o objeto de valor
	return Rate{
		capacity: capacity,
		every:    every,
	}
}

// newLimiter creates a new instance of Limiter based on the provided limiterDTO.
// It initializes the fields of Limiter based on values from limiterDTO and sets default values for empty fields.
func newLimiter(limiterDTO dto.Limiter) Limiter {
	return Limiter{
		maxHeaderSize:          NewBytes(limiterDTO.MaxHeaderSize),
		maxBodySize:            NewBytes(limiterDTO.MaxBodySize),
		maxMultipartMemorySize: NewBytes(limiterDTO.MaxMultipartMemorySize),
		rate:                   newRate(helper.IfNilReturns(limiterDTO.Rate, dto.Rate{})),
	}
}

// newRate creates a new instance of Rate based on the provided rateDTO.
// It initializes the fields of Rate based on values from rateDTO and sets default values for empty fields.
func newRate(rateDTO dto.Rate) Rate {
	var every time.Duration
	var err error
	if helper.IsNotEmpty(rateDTO.Every) {
		every, err = time.ParseDuration(rateDTO.Every)
		if helper.IsNotNil(err) {
			logger.Warning("Parse duration limiter.rate.every err:", err)
		}
	}
	return Rate{
		capacity: rateDTO.Capacity,
		every:    every,
	}
}

// MaxHeaderSize returns the value of the maxHeaderSize field in the Limiter object.
func (l Limiter) MaxHeaderSize() Bytes {
	return l.maxHeaderSize
}

// MaxBodySize returns the value of the maxBodySize field in the Limiter object.
func (l Limiter) MaxBodySize() Bytes {
	return l.maxBodySize
}

// MaxMultipartMemorySize returns the value of the maxMultipartMemorySize field in the Limiter object.
func (l Limiter) MaxMultipartMemorySize() Bytes {
	return l.maxMultipartMemorySize
}

// Rate returns the value of the rate field in the Limiter object.
func (l Limiter) Rate() Rate {
	return l.rate
}

// Capacity returns the value of the capacity field in the Rate struct.
func (r Rate) Capacity() int {
	return r.capacity
}

// Every returns the value of the every field in the Rate struct.
func (r Rate) Every() time.Duration {
	return r.every
}
