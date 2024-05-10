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
	rateTime "golang.org/x/time/rate"
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
	keys map[string]*rateTime.Limiter
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

// newLimiter creates a new Limiter object by initializing its fields based on the provided limiterJson and
// endpointLimiterJson.
// If endpointLimiterJson is not nil, its values will be used to initialize the Limiter fields.
// If endpointLimiterJson has values for the maxHeaderSize, maxBodySize, and maxMultipartMemorySize fields,
// they will be assigned to the respective fields in the Limiter object.
// If endpointLimiterJson has a value for the Rate field, it will be assigned to the rate field in the Limiter object.
// If endpointLimiterJson is nil and limiterJson is not nil, the values from limiterJson will be assigned to the
// respective fields in the Limiter object.
// If limiterJson has a value for the Rate field, it will be assigned to the rate field in the Limiter object.
// The function returns a pointer to the newly created Limiter object.
func newLimiter(limiterJson *LimiterJson, endpointLimiterJson *EndpointLimiterJson) *Limiter {
	var maxHeaderSize Bytes
	var maxBodySize Bytes
	var maxMultipartForm Bytes

	var limiterRateJson *RateJson
	var endpointLimiterRateJson *RateJson

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
	} else if helper.IsNotNil(limiterJson) {
		maxHeaderSize = limiterJson.MaxHeaderSize
		maxBodySize = limiterJson.MaxBodySize
		maxMultipartForm = limiterJson.MaxMultipartMemorySize
		limiterRateJson = limiterJson.Rate
	}

	return &Limiter{
		maxHeaderSize:          maxHeaderSize,
		maxBodySize:            maxBodySize,
		maxMultipartMemorySize: maxMultipartForm,
		rate:                   newRate(limiterRateJson, endpointLimiterRateJson),
	}
}

// newLimiterDefault creates a new Limiter object by initializing its fields with default values.
// The rate field is initialized using newRateDefault().
// The function returns a pointer to the newly created Limiter object.
func newLimiterDefault() *Limiter {
	return &Limiter{
		rate: newRateDefault(),
	}
}

// newRate creates a new Rate object by initializing its fields based on the provided rateJson and
// endpointRateJson.
// If endpointRateJson is not nil, its values will be used to initialize the Rate fields.
// If endpointRateJson has a value for the Every field, it will be assigned to the every field in the Rate object.
// If endpointRateJson has a value for the Capacity field, it will be assigned to the capacity field in the Rate object.
// If endpointRateJson is nil and rateJson is not nil, the values from rateJson will be assigned to the
// respective fields in the Rate object.
// The function returns a Rate object with the initialized fields: keys, mutex, capacity, and every.
func newRate(rateJson *RateJson, endpointRateJson *RateJson) Rate {
	var every Duration
	var capacity int

	if helper.IsNotNil(endpointRateJson) {
		if endpointRateJson.HasEvery() {
			every = endpointRateJson.Every
		}
		if endpointRateJson.HasCapacity() {
			capacity = endpointRateJson.Capacity
		}
	} else if helper.IsNotNil(rateJson) {
		every = rateJson.Every
		capacity = rateJson.Capacity
	}

	return Rate{
		keys:     map[string]*rateTime.Limiter{},
		mutex:    &sync.RWMutex{},
		capacity: capacity,
		every:    every,
	}
}

// newRateDefault creates a default Rate object with empty keys and a mutex.
// The function returns the newly created Rate object.
func newRateDefault() Rate {
	return Rate{
		keys:  map[string]*rateTime.Limiter{},
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

// Allow checks if the header size of the httpRequest is within the allowed limit.
// If the header size exceeds the limit, it returns a new instance of ErrHeaderTooLarge.
// Otherwise, it continues to check if the httpRequest has a body. If so, it invokes
// the allowBody method to check if the body size is within the allowed limit.
// If the body size exceeds the limit, it returns a new instance of ErrPayloadTooLarge.
// The method returns an error if any of the checks fail, otherwise it returns nil.
func (l Limiter) Allow(httpRequest *HttpRequest) (err error) {
	err = l.allowHeader(httpRequest)
	if helper.IsNil(err) && helper.IsNotNil(httpRequest.Body()) {
		err = l.allowBody(httpRequest)
	}
	return err
}

// HasData returns a boolean value indicating if the Rate object has data.
// It checks if the capacity and every field of the Rate object are greater than 0.
func (r Rate) HasData() bool {
	return helper.IsGreaterThan(r.Capacity(), 0) && helper.IsGreaterThan(r.Every(), 0)
}

// NoData returns a boolean value indicating whether the Rate object has data.
// It checks if the Capacity and Every field of the Rate object are greater than 0.
func (r Rate) NoData() bool {
	return !r.HasData()
}

// Capacity returns the maximum number of allowed requests within a given time period.
//
// It retrieves the value of the capacity field in the Rate object.
// The capacity field represents the maximum number of allowed requests.
func (r Rate) Capacity() int {
	return r.capacity
}

// Every returns the value of the every field in the Rate object.
// The every field represents the frequency of allowed requests in the Rate configuration for rate limiting.
// It returns a Duration object, which is a wrapper around the time.Duration type.
func (r Rate) Every() Duration {
	return r.every
}

// EveryTime returns the time.Duration value of the every field in the Rate object.
// It calls the Time() method of the Duration object to convert the value to time.Duration.
// The resulting time.Duration value is returned.
func (r Rate) EveryTime() time.Duration {
	return r.every.Time()
}

// Allow checks if the Rate object has rate limiters for the given key.
// If the key exists, it checks if the rate limiter allows the request.
// If the key does not exist, it creates a new rate limiter with the Rate object's configuration and adds it to the keys map.
// If the rate limiter does not allow the request, it returns an error indicating too many requests.
// If the Rate object has no data, it returns nil.
//
// The key parameter is a string representing the unique identifier for the rate limiter.
//
// The mutex field of the Rate object is used for thread safety during rate limiter creation and access to the keys map.
// It should be locked with the Lock method before accessing the keys map and unlocked with the Unlock method after access.
//
// The Allow method returns an error if the rate limiter does not allow the request, otherwise it returns nil.
func (r Rate) Allow(key string) error {
	if r.NoData() {
		return nil
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	timeRateLimiter, exists := r.keys[key]
	if !exists {
		timeRateLimiter = rateTime.NewLimiter(rateTime.Every(r.EveryTime()), r.Capacity())
		r.keys[key] = timeRateLimiter
	}

	if !timeRateLimiter.Allow() {
		return mapper.NewErrTooManyRequests(r.Capacity(), r.EveryTime())
	}
	return nil
}

// allowHeader checks if the header size of the httpRequest is within the allowed limit.
// If the header size exceeds the limit, it returns a new instance of ErrHeaderTooLarge.
// The method returns an error if the check fails, otherwise it returns nil.
func (l Limiter) allowHeader(httpRequest *HttpRequest) (err error) {
	maxSizeAllowed := l.MaxHeaderSize()
	if helper.IsGreaterThan(httpRequest.HeaderSize(), maxSizeAllowed) {
		err = mapper.NewErrHeaderTooLarge(maxSizeAllowed.String())
	}
	return err
}

// allowBody checks the size of the body in the httpRequest to see if it is within the allowed limit.
// It first determines the actual size limit based on the provided Content-Type.
// If the Content-Type is "multipart/form-data", it uses the MaxMultipartMemorySize value.
// Otherwise, it uses the MaxBodySize value.
//
// It creates a readCloser with MaxBytesReader to limit the size of the body to maxSizeAllowed.
// It then reads the entire body from the readCloser using io.ReadAll.
// If there is an error while reading the body, it returns a new instance of ErrPayloadTooLarge
// with the maxSizeAllowed value as the payload size limit.
//
// The method returns an error if the body size exceeds the limit, otherwise it returns nil.
func (l Limiter) allowBody(httpRequest *HttpRequest) (err error) {
	maxSizeAllowed := l.MaxBodySize()
	if helper.ContainsIgnoreCase(httpRequest.Header().Get("Content-Type"), "multipart/form-data") {
		maxSizeAllowed = l.MaxMultipartMemorySize()
	}

	bodyBuffer := httpRequest.Body().Buffer()
	readCloser := http.MaxBytesReader(nil, io.NopCloser(bodyBuffer), int64(maxSizeAllowed))

	_, err = io.ReadAll(readCloser)
	if helper.IsNotNil(err) {
		err = mapper.NewErrPayloadTooLarge(maxSizeAllowed.String())
	}

	return err
}
