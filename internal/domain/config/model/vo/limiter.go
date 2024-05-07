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
)

// Limiter represents the configuration for rate limiting in the Gopen application.
type Limiter struct {
	// maxHeaderSize represents the maximum size of the header in bytes for rate limiting.
	maxHeaderSize Bytes
	// maxBodySize represents the maximum size of the body in bytes for rate limiting.
	maxBodySize Bytes
	// maxMultipartMemorySize represents the maximum memory size for multipart httpRequest bodies.
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
	// maxMultipartMemorySize represents the maximum memory size for multipart httpRequest bodies.
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
	every Duration
}

// EndpointRate represents the configuration for rate limiting for an endpoint in the Gopen application.
// It includes the capacity, which represents the maximum number of allowed requests within a given time period,
// and the every field, which represents the frequency of allowed requests in the Rate configuration for rate limiting.
type EndpointRate struct {
	// capacity represents the maximum number of allowed requests within a given time period.
	capacity int
	// every represents the frequency of allowed requests in the Rate configuration for rate limiting.
	every Duration
}

// newEndpointLimiter creates a new instance of EndpointLimiter using the provided Limiter
// and EndpointLimiterJson objects. It sets the values of maxHeaderSize, maxBodySize, and
// maxMultipartMemorySize based on the Limiter object by default. If an EndpointLimiterJson
// object is provided, it overrides the default values. The rate configuration is set using
// the newEndpointRate function.
//
// If the EndpointLimiterJson object is not provided or has empty/zero values for maxHeaderSize,
// maxBodySize, and maxMultipartMemorySize, the values from the Limiter object are used as defaults.
//
// Returns a pointer to the newly created EndpointLimiter object.
func newEndpointLimiter(limiterVO Limiter, endpointLimiterVO *EndpointLimiterJson) *EndpointLimiter {
	// por padrão obtemos o limiter.max-header-size configurado na raiz, caso não informado um valor padrão é retornado
	maxHeaderSize := limiterVO.MaxHeaderSize()
	// por padrão obtemos o limiter.max-body-size configurado na raiz, caso não informado um valor padrão é retornado
	maxBodySize := limiterVO.MaxBodySize()
	// por padrão obtemos o limiter.max-multipart-form-size configurado na raiz, caso não informado um valor padrão é retornado
	maxMultipartForm := limiterVO.MaxMultipartMemorySize()
	// instanciamos o endpointRate
	var endpointRate *RateJson

	// caso informado no endpoint damos prioridade
	if helper.IsNotNil(endpointLimiterVO) {
		endpointRate = endpointLimiterVO.Rate
		if endpointLimiterVO.HasMaxHeaderSize() {
			maxHeaderSize = endpointLimiterVO.MaxHeaderSize
		}
		if endpointLimiterVO.HasMaxBodySize() {
			maxBodySize = endpointLimiterVO.MaxBodySize
		}
		if endpointLimiterVO.HasMaxMultipartMemorySize() {
			maxMultipartForm = endpointLimiterVO.MaxMultipartMemorySize
		}
	}

	//construímos o objeto de valor limiter
	return &EndpointLimiter{
		maxHeaderSize:          maxHeaderSize,
		maxBodySize:            maxBodySize,
		maxMultipartMemorySize: maxMultipartForm,
		rate:                   newEndpointRate(limiterVO.Rate(), endpointRate),
	}
}

func newEndpointLimiterStatic(limiterVO Limiter) *EndpointLimiter {
	maxHeaderSize := limiterVO.MaxHeaderSize()
	maxBodySize := limiterVO.MaxBodySize()
	maxMultipartForm := limiterVO.MaxMultipartMemorySize()
	//construímos o objeto de valor limiter para os endpoints estáticos
	return &EndpointLimiter{
		maxHeaderSize:          maxHeaderSize,
		maxBodySize:            maxBodySize,
		maxMultipartMemorySize: maxMultipartForm,
		rate:                   newEndpointRateStatic(limiterVO.Rate()),
	}
}

func newEndpointRate(rateVO Rate, endpointRateVO *RateJson) *EndpointRate {
	// por padrão obtemos o limiter.rate.every configurado na raiz, caso não informado um valor padrão é retornado
	every := rateVO.Every()
	// por padrão obtemos o limiter.rate.capacity configurado na raiz, caso não informado um valor padrão é retornado
	capacity := rateVO.Capacity()
	// caso informado no endpoint damos prioridade
	if helper.IsNotNil(endpointRateVO) {
		if endpointRateVO.HasEvery() {
			every = endpointRateVO.Every
		}
		if endpointRateVO.HasCapacity() {
			capacity = endpointRateVO.Capacity
		}
	}

	// montamos o objeto de valor
	return &EndpointRate{
		capacity: capacity,
		every:    every,
	}
}

func newEndpointRateStatic(rateVO Rate) *EndpointRate {
	// montamos o objeto de valor
	return &EndpointRate{
		capacity: rateVO.Capacity(),
		every:    rateVO.Every(),
	}
}

// newLimiter creates a new instance of Limiter using the provided LimiterJson object.
// If the LimiterJson object is nil, it returns a Limiter object with default values.
// Otherwise, it sets the values of maxHeaderSize, maxBodySize, maxMultipartMemorySize,
// and rate based on the values of the LimiterJson object.
// Returns a Limiter object with the configured values.
func newLimiter(limiterJsonVO *LimiterJson) Limiter {
	if helper.IsNil(limiterJsonVO) {
		return Limiter{}
	}
	return Limiter{
		maxHeaderSize:          limiterJsonVO.MaxHeaderSize,
		maxBodySize:            limiterJsonVO.MaxBodySize,
		maxMultipartMemorySize: limiterJsonVO.MaxMultipartMemorySize,
		rate:                   newRate(&limiterJsonVO.Rate),
	}
}

// newRate creates a new instance of Rate using the provided RateJson object.
// The Rate structs capacity and every field are set based on the values of the RateJson object.
func newRate(rateJsonVO *RateJson) Rate {
	return Rate{
		capacity: rateJsonVO.Capacity,
		every:    rateJsonVO.Every,
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
func (r Rate) Every() Duration {
	return r.every
}

// MaxHeaderSize returns the value of the maxHeaderSize field in the EndpointLimiter object.
func (e EndpointLimiter) MaxHeaderSize() Bytes {
	return e.maxHeaderSize
}

// MaxBodySize returns the value of the maxBodySize field in the EndpointLimiter object.
func (e EndpointLimiter) MaxBodySize() Bytes {
	return e.maxBodySize
}

// MaxMultipartMemorySize returns the value of the maxMultipartMemorySize field in the EndpointLimiter object.
func (e EndpointLimiter) MaxMultipartMemorySize() Bytes {
	return e.maxMultipartMemorySize
}

// Rate returns the value of the rate field in the EndpointLimiter object.
func (e EndpointLimiter) Rate() *EndpointRate {
	return e.rate
}

// Every returns the value of the every field in the EndpointRate struct.
func (e EndpointRate) Every() Duration {
	return e.every
}

// Capacity returns the value of the capacity field in the EndpointRate struct.
// It represents the maximum number of allowed requests within a given time period.
func (e EndpointRate) Capacity() int {
	return e.capacity
}
