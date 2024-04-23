/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	rate *EndpointRate
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

// newEndpointLimiter creates a new instance of EndpointLimiter based on the provided limiterVO  and endpointLimiterVO.
// It initializes the fields of EndpointLimiter based on values from limiterVO and endpointLimiterVO,
// giving priority to the values from endpointLimiterVO if they are present. The default values from limiterVO are used otherwise.
// The maxHeaderSize field is set to the value of limiterVO.MaxHeaderSize. If endpointLimiterVO.HasMaxHeaderSize()
// returns true, it sets the maxHeaderSize field to the value of endpointLimiterVO.MaxHeaderSize.
// The maxBodySize field is set to the value of limiterVO.MaxBodySize. If endpointLimiterVO.HasMaxBodySize()
// returns true, it sets the maxBodySize field to the value of endpointLimiterVO.MaxBodySize.
// The maxMultipartMemorySize field is set to the value of limiterVO.MaxMultipartMemorySize. If endpointLimiterVO.HasMaxMultipartFormSize()
// returns true, it sets the maxMultipartMemorySize field to the value of endpointLimiterVO.MaxMultipartMemorySize.
// The rate field is set to a new instance of EndpointRate, created by calling the newEndpointRate function with
// limiterVO.Rate and endpointLimiterVO.Rate as parameters.
// The function returns a pointer to the created EndpointLimiter object.
func newEndpointLimiter(limiterVO Limiter, endpointLimiterVO *EndpointLimiter) *EndpointLimiter {
	// por padrão obtemos o limiter.max-header-size configurado na raiz, caso não informado um valor padrão é retornado
	maxHeaderSize := limiterVO.MaxHeaderSize()
	// por padrão obtemos o limiter.max-body-size configurado na raiz, caso não informado um valor padrão é retornado
	maxBodySize := limiterVO.MaxBodySize()
	// por padrão obtemos o limiter.max-multipart-form-size configurado na raiz, caso não informado um valor padrão é retornado
	maxMultipartForm := limiterVO.MaxMultipartMemorySize()
	// instanciamos o endpointRate
	var endpointRate *EndpointRate

	// caso informado no endpoint damos prioridade
	if helper.IsNotNil(endpointLimiterVO) {
		endpointRate = endpointLimiterVO.Rate()
		if endpointLimiterVO.HasMaxHeaderSize() {
			maxHeaderSize = endpointLimiterVO.MaxHeaderSize()
		}
		if endpointLimiterVO.HasMaxBodySize() {
			maxBodySize = endpointLimiterVO.MaxBodySize()
		}
		if endpointLimiterVO.HasMaxMultipartFormSize() {
			maxMultipartForm = endpointLimiterVO.MaxMultipartMemorySize()
		}
	}

	//construímos o limiter vo
	return &EndpointLimiter{
		maxHeaderSize:          maxHeaderSize,
		maxBodySize:            maxBodySize,
		maxMultipartMemorySize: maxMultipartForm,
		rate:                   newEndpointRate(limiterVO.Rate(), endpointRate),
	}
}

func newEndpointRate(rateVO Rate, endpointRateVO *EndpointRate) *EndpointRate {
	// por padrão obtemos o limiter.rate.every configurado na raiz, caso não informado um valor padrão é retornado
	every := rateVO.Every()
	// por padrão obtemos o limiter.rate.capacity configurado na raiz, caso não informado um valor padrão é retornado
	capacity := rateVO.Capacity()
	// caso informado no endpoint damos prioridade
	if helper.IsNotNil(endpointRateVO) {
		if endpointRateVO.HasEvery() {
			every = endpointRateVO.Every()
		}
		if endpointRateVO.HasCapacity() {
			capacity = endpointRateVO.Capacity()
		}
	}

	// montamos o objeto de valor
	return &EndpointRate{
		capacity: capacity,
		every:    every,
	}
}

// newLimiterFromDTO creates a new instance of Limiter based on the provided limiterDTO.
// It initializes the fields of Limiter based on values from limiterDTO and sets default values for empty fields.
func newLimiterFromDTO(limiterDTO *dto.Limiter) Limiter {
	return Limiter{
		maxHeaderSize:          NewBytes(limiterDTO.MaxHeaderSize),
		maxBodySize:            NewBytes(limiterDTO.MaxBodySize),
		maxMultipartMemorySize: NewBytes(limiterDTO.MaxMultipartMemorySize),
		rate:                   newRateFromDTO(limiterDTO.Rate),
	}
}

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

func newEndpointLimiterFromDTO(limiterDTO *dto.EndpointLimiter) *EndpointLimiter {
	if helper.IsNil(limiterDTO) {
		return nil
	}
	return &EndpointLimiter{
		maxHeaderSize:          NewBytes(limiterDTO.MaxHeaderSize),
		maxBodySize:            NewBytes(limiterDTO.MaxBodySize),
		maxMultipartMemorySize: NewBytes(limiterDTO.MaxMultipartMemorySize),
		rate:                   newEndpointRateFromDTO(limiterDTO.Rate),
	}
}

func newEndpointRateFromDTO(rateDTO *dto.Rate) *EndpointRate {
	if helper.IsNil(rateDTO) {
		return nil
	}

	var every time.Duration
	var err error
	if helper.IsNotEmpty(rateDTO.Every) {
		every, err = time.ParseDuration(rateDTO.Every)
		if helper.IsNotNil(err) {
			logger.Warning("Parse duration limiter.rate.every err:", err)
		}
	}
	return &EndpointRate{
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
func (e EndpointLimiter) Rate() *EndpointRate {
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

func (e EndpointRate) EveryStr() string {
	if e.HasEvery() {
		return e.every.String()
	}
	return ""
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
