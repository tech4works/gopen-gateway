package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"time"
)

// Limiter represents the configuration for rate limiting in the Gopen application.
type Limiter struct {
	// maxHeaderSize represents the maximum size of the header in bytes for rate limiting.
	maxHeaderSize Bytes
	// maxBodySize represents the maximum size of the body in bytes for rate limiting.
	maxBodySize Bytes
	// maxMultipartMemorySize represents the maximum memory size for multipart request bodies.
	maxMultipartMemorySize Bytes
	// rate represents the configuration for `rate` limiting in the Limiter struct. It specifies the capacity and
	// frequency of allowed requests.
	rate Rate
}

// EndpointLimiter represents the configuration for rate limiting for an endpoint in the Gopen application.
// It includes the maximum sizes for the header, body, and multipart memory, as well as the rate configuration for limiting requests.
type EndpointLimiter struct {
	// maxHeaderSize represents the maximum size of the header in bytes for rate limiting.
	maxHeaderSize Bytes
	// maxBodySize represents the maximum size of the body in bytes for rate limiting.
	maxBodySize Bytes
	// maxMultipartMemorySize represents the maximum memory size for multipart request bodies.
	maxMultipartMemorySize Bytes
	// rate represents the configuration for `rate` limiting in the Limiter struct. It specifies the capacity and
	// frequency of allowed requests.
	rate EndpointRate
}

// Rate represents the configuration for rate limiting. It specifies the capacity
// and frequency of allowed requests.
type Rate struct {
	// capacity represents the maximum number of allowed requests within a given time period.
	capacity int
	// every represents the frequency of allowed requests in the Rate configuration for rate limiting.
	every time.Duration
}

// EndpointRate represents the configuration for rate limiting for an endpoint in the Gopen application.
// It includes the capacity, which represents the maximum number of allowed requests within a given time period,
// and the every field, which represents the frequency of allowed requests in the Rate configuration for rate limiting.
type EndpointRate struct {
	// capacity represents the maximum number of allowed requests within a given time period.
	capacity int
	// every represents the frequency of allowed requests in the Rate configuration for rate limiting.
	every time.Duration
}

// newEndpointLimiter creates a new instance of EndpointLimiter based on the provided Limiter and EndpointLimiter.
// It sets the maxHeaderSize, maxBodySize, maxMultipartMemorySize fields of EndpointLimiter
// based on values from the Limiter and EndpointLimiter arguments.
// If the maxHeaderSize field is specified in EndpointLimiter, it takes precedence over the maxHeaderSize field in Limiter.
// If the maxBodySize field is specified in EndpointLimiter, it takes precedence over the maxBodySize field in Limiter.
// If the maxMultipartMemorySize field is specified in EndpointLimiter,
// it takes precedence over the maxMultipartMemorySize field in Limiter.
// It also sets the rate field of EndpointLimiter by calling the newEndpointRate function
// with the rate field from Limiter and EndpointLimiter as arguments.
func newEndpointLimiter(limiterVO Limiter, endpointLimiterVO EndpointLimiter) EndpointLimiter {
	// por padrão obtemos o limiter.max-header-size configurado na raiz, caso não informado um valor padrão é retornado
	maxHeaderSize := limiterVO.MaxHeaderSize()
	// caso informado no endpoint damos prioridade
	if endpointLimiterVO.HasMaxHeaderSize() {
		maxHeaderSize = endpointLimiterVO.MaxHeaderSize()
	}

	// por padrão obtemos o limiter.max-body-size configurado na raiz, caso não informado um valor padrão é retornado
	maxBodySize := limiterVO.MaxBodySize()
	// caso informado no endpoint damos prioridade
	if endpointLimiterVO.HasMaxBodySize() {
		maxBodySize = endpointLimiterVO.MaxBodySize()
	}

	// por padrão obtemos o limiter.max-multipart-form-size configurado na raiz, caso não informado um valor padrão é retornado
	maxMultipartForm := limiterVO.MaxMultipartMemorySize()
	// caso informado no endpoint damos prioridade
	if endpointLimiterVO.HasMaxMultipartFormSize() {
		maxMultipartForm = endpointLimiterVO.MaxMultipartMemorySize()
	}

	//construímos o limiter vo
	return EndpointLimiter{
		maxHeaderSize:          maxHeaderSize,
		maxBodySize:            maxBodySize,
		maxMultipartMemorySize: maxMultipartForm,
		rate:                   newEndpointRate(limiterVO.Rate(), endpointLimiterVO.Rate()),
	}
}

// newEndpointRate creates a new instance of EndpointRate based on the provided Rate and EndpointRate.
// It sets the capacity and every field of EndpointRate based on values from the Rate and EndpointRate arguments.
// If the every field is specified in EndpointRate, it takes precedence over the every field in Rate.
// If the capacity field is specified in EndpointRate, it takes precedence over the capacity field in Rate.
func newEndpointRate(rateVO Rate, endpointRateVO EndpointRate) EndpointRate {
	// por padrão obtemos o limiter.rate.every configurado na raiz, caso não informado um valor padrão é retornado
	every := rateVO.Every()
	// caso informado no endpoint damos prioridade
	if endpointRateVO.HasEvery() {
		every = endpointRateVO.Every()
	}

	// por padrão obtemos o limiter.rate.capacity configurado na raiz, caso não informado um valor padrão é retornado
	capacity := rateVO.Capacity()
	// caso informado no endpoint damos prioridade
	if endpointRateVO.HasCapacity() {
		capacity = endpointRateVO.Capacity()
	}

	// montamos o objeto de valor
	return EndpointRate{
		capacity: capacity,
		every:    every,
	}
}

// newLimiterFromDTO creates a new instance of Limiter based on the provided limiterDTO.
// It initializes the fields of Limiter based on values from limiterDTO and sets default values for empty fields.
func newLimiterFromDTO(limiterDTO dto.Limiter) Limiter {
	return Limiter{
		maxHeaderSize:          NewBytes(limiterDTO.MaxHeaderSize),
		maxBodySize:            NewBytes(limiterDTO.MaxBodySize),
		maxMultipartMemorySize: NewBytes(limiterDTO.MaxMultipartMemorySize),
		rate:                   newRateFromDTO(limiterDTO.Rate),
	}
}

// newRateFromDTO creates a new instance of Rate based on the provided rateDTO.
// It initializes the fields of Rate based on values from rateDTO and sets default values for empty fields.
func newRateFromDTO(rateDTO dto.Rate) Rate {
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

// newEndpointLimiterFromDTO creates a new instance of EndpointLimiter based on the provided limiterDTO.
// It initializes the fields of EndpointLimiter based on values from limiterDTO and sets default values for empty fields.
func newEndpointLimiterFromDTO(limiterDTO dto.Limiter) EndpointLimiter {
	return EndpointLimiter{
		maxHeaderSize:          NewBytes(limiterDTO.MaxHeaderSize),
		maxBodySize:            NewBytes(limiterDTO.MaxBodySize),
		maxMultipartMemorySize: NewBytes(limiterDTO.MaxMultipartMemorySize),
		rate:                   newEndpointRateFromDTO(limiterDTO.Rate),
	}
}

// newEndpointRateFromDTO creates a new instance of EndpointRate based on the provided rateDTO.
// It initializes the fields of EndpointRate based on values from rateDTO and sets default values for empty fields.
func newEndpointRateFromDTO(rateDTO dto.Rate) EndpointRate {
	var every time.Duration
	var err error
	if helper.IsNotEmpty(rateDTO.Every) {
		every, err = time.ParseDuration(rateDTO.Every)
		if helper.IsNotNil(err) {
			logger.Warning("Parse duration limiter.rate.every err:", err)
		}
	}
	return EndpointRate{
		capacity: rateDTO.Capacity,
		every:    every,
	}
}

// MaxHeaderSize returns the value of the maxHeaderSize field in the Limiter object.
// If the maxHeaderSize field is greater than 0, it returns the value of that field.
// Otherwise, it returns a new Bytes value with the string "1MB" as the byte unit.
func (l Limiter) MaxHeaderSize() Bytes {
	if helper.IsGreaterThan(l.maxHeaderSize, 0) {
		return l.maxHeaderSize
	}
	return NewBytes("1MB")
}

// MaxBodySize returns the value of the maxBodySize field in the Limiter object.
// If the maxBodySize field is greater than 0, it returns the value of that field.
// Otherwise, it returns a new Bytes value with the string "3MB" as the byte unit.
func (l Limiter) MaxBodySize() Bytes {
	if helper.IsGreaterThan(l.maxBodySize, 0) {
		return l.maxBodySize
	}
	return NewBytes("3MB")
}

// MaxMultipartMemorySize returns the value of the maxMultipartMemorySize field in the Limiter object.
// If the maxMultipartMemorySize field is greater than 0, it returns the value of that field.
// Otherwise, it returns a new Bytes value with the string "5MB" as the byte unit.
func (l Limiter) MaxMultipartMemorySize() Bytes {
	if helper.IsGreaterThan(l.maxMultipartMemorySize, 0) {
		return l.maxMultipartMemorySize
	}
	return NewBytes("5MB")
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

// HasMaxHeaderSize returns true if the maxHeaderSize field in the EndpointLimiter object is greater than 0.
// Otherwise, it returns false.
func (e EndpointLimiter) HasMaxHeaderSize() bool {
	return helper.IsGreaterThan(e.maxHeaderSize, 0)
}

// MaxHeaderSize returns the value of the maxHeaderSize field in the EndpointLimiter object.
func (e EndpointLimiter) MaxHeaderSize() Bytes {
	return e.maxHeaderSize
}

// HasMaxBodySize returns true if the maxBodySize field in the Limiter object is greater than 0.
// Otherwise, it returns false.
func (e EndpointLimiter) HasMaxBodySize() bool {
	return helper.IsGreaterThan(e.maxBodySize, 0)
}

// MaxBodySize returns the value of the maxBodySize field in the EndpointLimiter object.
func (e EndpointLimiter) MaxBodySize() Bytes {
	return e.maxBodySize
}

// HasMaxMultipartFormSize returns true if the maxMultipartMemorySize field in the EndpointLimiter object is greater than 0.
// Otherwise, it returns false.
func (e EndpointLimiter) HasMaxMultipartFormSize() bool {
	return helper.IsGreaterThan(e.maxMultipartMemorySize, 0)
}

// MaxMultipartMemorySize returns the value of the maxMultipartMemorySize field in the EndpointLimiter object.
func (e EndpointLimiter) MaxMultipartMemorySize() Bytes {
	return e.maxMultipartMemorySize
}

// Rate returns the value of the rate field in the EndpointLimiter object.
func (e EndpointLimiter) Rate() EndpointRate {
	return e.rate
}

// HasEvery determines if the every field in the EndpointRate struct is greater than 0.
// It returns true if the every field is greater than 0, otherwise it returns false.
func (e EndpointRate) HasEvery() bool {
	return helper.IsGreaterThan(e.every, 0)
}

// Every returns the value of the every field in the EndpointRate struct.
func (e EndpointRate) Every() time.Duration {
	return e.every
}

// HasCapacity determines if the capacity field in the EndpointRate struct is greater than 0.
// It returns true if the capacity field is greater than 0, otherwise it returns false.
func (e EndpointRate) HasCapacity() bool {
	return helper.IsGreaterThan(e.capacity, 0)
}

// Capacity returns the value of the capacity field in the EndpointRate struct.
// It represents the maximum number of allowed requests within a given time period.
func (e EndpointRate) Capacity() int {
	return e.capacity
}
