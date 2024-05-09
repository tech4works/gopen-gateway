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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	timeRate "golang.org/x/time/rate"
	"io"
	"net/http"
	"sync"
	"time"
)

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

type Rate struct {
	// keys represents a map that stores rate limiters for each key.
	// The keys are of type string and the values are of type *rate.Limiter.
	// It is a field of the rateLimiter struct.
	// This map is used to store and manage rate limiters for different keys.
	//
	// Note: The keys map and other related structures should be properly initialized before accessing this field.
	// The rateLimiter type should be used to access this field.
	// Other types should not have direct access to this field.
	keys map[string]*timeRate.Limiter
	// mutex is a pointer to a sync.RWMutex object. It is used for thread-safety
	// in the rateLimiter struct. It should be locked and unlocked using
	// the Lock and Unlock methods respectively to protect concurrent accesses
	// to shared resources.
	mutex *sync.RWMutex
	// capacity represents the maximum number of allowed requests within a given time period.
	capacity int
	// every represents the frequency of allowed requests in the Rate configuration for rate limiting.
	every Duration
}

func newLimiter(limiterJson *LimiterJson, endpointLimiterJson *EndpointLimiterJson) *Limiter {
	// instanciamos o valor vazio
	var maxHeaderSize Bytes
	var maxBodySize Bytes
	var maxMultipartForm Bytes

	var limiterRateJson *RateJson
	var endpointLimiterRateJson *RateJson

	// caso informado no topo do json de configuração inserimos inicialmente
	if helper.IsNotNil(limiterJson) {
		maxHeaderSize = limiterJson.MaxHeaderSize
		maxBodySize = limiterJson.MaxBodySize
		maxMultipartForm = limiterJson.MaxMultipartMemorySize
		limiterRateJson = limiterJson.Rate
	}

	// caso informado no endpoint damos prioridade
	if helper.IsNotNil(endpointLimiterJson) {
		if endpointLimiterJson.HasMaxHeaderSize() {
			maxHeaderSize = endpointLimiterJson.MaxHeaderSize
		}
		if endpointLimiterJson.HasMaxBodySize() {
			maxBodySize = endpointLimiterJson.MaxBodySize
		}
		if endpointLimiterJson.HasMaxMultipartMemorySize() {
			maxMultipartForm = endpointLimiterJson.MaxMultipartMemorySize
		}
		endpointLimiterRateJson = endpointLimiterJson.Rate
	}

	//construímos o objeto de valor limiter
	return &Limiter{
		maxHeaderSize:          maxHeaderSize,
		maxBodySize:            maxBodySize,
		maxMultipartMemorySize: maxMultipartForm,
		rate:                   newRate(limiterRateJson, endpointLimiterRateJson),
	}
}

func newLimiterDefault() *Limiter {
	return &Limiter{
		rate: newRateDefault(),
	}
}

func newRate(rateJson *RateJson, endpointRateJson *RateJson) Rate {
	// instanciamos os valores vazios
	var every Duration
	var capacity int

	// caso informado no topo do json preenchemos
	if helper.IsNotNil(rateJson) {
		every = rateJson.Every
		capacity = rateJson.Capacity
	}

	// caso informado no endpoint damos prioridade
	if helper.IsNotNil(endpointRateJson) {
		if endpointRateJson.HasEvery() {
			every = endpointRateJson.Every
		}
		if endpointRateJson.HasCapacity() {
			capacity = endpointRateJson.Capacity
		}
	}

	// montamos o objeto de valor
	return Rate{
		keys:     map[string]*timeRate.Limiter{},
		mutex:    &sync.RWMutex{},
		capacity: capacity,
		every:    every,
	}
}

func newRateDefault() Rate {
	return Rate{
		keys:  map[string]*timeRate.Limiter{},
		mutex: &sync.RWMutex{},
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

func (l Limiter) Allow(httpRequest *HttpRequest) (err error) {
	// checamos primeiramente o tamanho do header
	err = l.allowHeader(httpRequest)
	// se não ocorreu nenhum erro ao validar o header e ele tem body, validamos o body
	if helper.IsNil(err) && helper.IsNotNil(httpRequest.Body()) {
		err = l.allowBody(httpRequest)
	}
	// retornamos um possível erro
	return err
}

func (r Rate) HasData() bool {
	return helper.IsGreaterThan(r.Capacity(), 0) && helper.IsGreaterThan(r.Every(), 0)
}

func (r Rate) NoData() bool {
	return !r.HasData()
}

func (r Rate) Capacity() int {
	return r.capacity
}

func (r Rate) Every() Duration {
	return r.every
}

func (r Rate) EveryTime() time.Duration {
	return r.every.Time()
}

func (r Rate) Allow(key string) error {
	// verificamos se a dados para seguir com a validação, caso não retornamos nil
	if r.NoData() {
		return nil
	}

	// bloqueamos o mutex
	r.mutex.Lock()
	// deixamos desbloqueado apenas quando a func terminar
	defer r.mutex.Unlock()

	// verificamos se a chave ja existe
	timeRateLimiter, exists := r.keys[key]
	if !exists {
		// caso não exista, adicionamos e preenchemos a variável
		timeRateLimiter = timeRate.NewLimiter(timeRate.Every(r.EveryTime()), r.Capacity())
		r.keys[key] = timeRateLimiter
	}

	// verificamos com base no objeto obtido pela chave se ele esta permitido
	if !timeRateLimiter.Allow() {
		return mapper.NewErrTooManyRequests(r.Capacity(), r.EveryTime())
	}

	// se tudo ocorrer bem ele retorna nil
	return nil
}

func (l Limiter) allowHeader(httpRequest *HttpRequest) (err error) {
	// instanciamos o valor máximo permitido
	maxSizeAllowed := l.MaxHeaderSize()

	// checamos se o header é maior que o permitido
	if helper.IsGreaterThan(httpRequest.HeaderSize(), maxSizeAllowed) {
		err = mapper.NewErrHeaderTooLarge(maxSizeAllowed.String())
	}

	// retornamos um possível erro
	return err
}

func (l Limiter) allowBody(httpRequest *HttpRequest) (err error) {
	// verificamos qual Content-Type fornecido, para obter a config real
	maxSizeAllowed := l.MaxBodySize()
	if helper.ContainsIgnoreCase(httpRequest.Header().Get("Content-Type"), "multipart/form-data") {
		maxSizeAllowed = l.MaxMultipartMemorySize()
	}

	// instanciamos o body
	bodyBuffer := httpRequest.Body().Buffer()
	// verificamos o tamanho utilizando o maxBytesReader, e logo em seguida se der certo, voltamos o body para requisição
	readCloser := http.MaxBytesReader(nil, io.NopCloser(bodyBuffer), int64(maxSizeAllowed))
	// lemos o body read closer
	_, err = io.ReadAll(readCloser)
	// verificamos se ocorreu algum erro
	if helper.IsNotNil(err) {
		err = mapper.NewErrPayloadTooLarge(maxSizeAllowed.String())
	}

	// retornamos um possível erro
	return err
}
