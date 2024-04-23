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

package infra

import (
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

// rateLimiterProvider represents a rate limiter provider that manages rate limiters for different keys.
// It has a keys map to store rate limiters for each key, a mutex for thread safety, and configuration parameters
// every (time.Duration), limit (rate.Limit), and burst (int).
type rateLimiterProvider struct {
	// keys represents a map that stores rate limiters for each key.
	// The keys are of type string and the values are of type *rate.Limiter.
	// It is a field of the rateLimiterProvider struct.
	// This map is used to store and manage rate limiters for different keys.
	//
	// Note: The keys map and other related structures should be properly initialized before accessing this field.
	// The rateLimiterProvider type should be used to access this field.
	// Other types should not have direct access to this field.
	keys map[string]*rate.Limiter
	// mutex is a pointer to a sync.RWMutex object. It is used for thread-safety
	// in the rateLimiterProvider struct. It should be locked and unlocked using
	// the Lock and Unlock methods respectively to protect concurrent accesses
	// to shared resources.
	mutex *sync.RWMutex
	// every represents the frequency of allowed requests in the rateLimiterProvider configuration for rate limiting.
	every time.Duration
	// limit represents the rate limit value used for rate limiting in the rateLimiterProvider struct.
	limit rate.Limit
	// burst represents the maximum number of requests that can be allowed in a given time period.
	// It is a field of the rateLimiterProvider struct, and it is used in the rate limiting process.
	//
	// Note: The rateLimiterProvider type should be used to access this field.
	// Other types should not have direct access to this field.
	burst int
}

type RateLimiterProvider interface {
	Allow(key string) error
}

// NewRateLimiterProvider creates a new instance of the RateLimiterProvider interface.
// It takes a time duration 'every' and an integer 'limit' as parameters.
// It returns a RateLimiterProvider object.
func NewRateLimiterProvider(endpointRateVO *vo.EndpointRate) RateLimiterProvider {
	return &rateLimiterProvider{
		keys:  map[string]*rate.Limiter{},
		mutex: &sync.RWMutex{},
		every: endpointRateVO.Every(),
		limit: rate.Every(endpointRateVO.Every()),
		burst: endpointRateVO.Capacity(),
	}
}

// Allow checks if the provided key is allowed to proceed or if it exceeds the limit.
// It locks the mutex to ensure thread safety and releases it at the end using defer.
//
// If the key does not exist in the keys map, the mutex is unlocked to allow other goroutines to access it.
// Then, the addKey() function is called to add the key to the keys map.
// Afterward, the mutex is locked again to protect the keys map from concurrent accesses.
//
// If the limiter associated with the key does not allow the request (limiter.Allow() returns false),
// the function returns a new instance of the error ErrTooManyRequests, passing the defined burst and every values.
// Otherwise, it returns nil to indicate that the request is allowed.
//
// Note: The keys map and other related structures should be properly initialized before calling this method.
// The rateLimiterProvider type should be used to access this method.
// Other types should not have direct access to this method.
func (r *rateLimiterProvider) Allow(key string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	limiter, exists := r.keys[key]
	if !exists {
		r.mutex.Unlock()
		limiter = r.addKey(key)
		r.mutex.Lock()
	}

	if !limiter.Allow() {
		return domainmapper.NewErrTooManyRequests(r.burst, r.every)
	}

	return nil
}

// addKey creates a new rate limiter for the provided key and adds it to the keys map.
// It locks the mutex to ensure thread safety and releases it at the end using defer.
// It creates a new instance of rate.Limiter using the limit and burst values from the rateLimiterProvider.
// Then, it adds the limiter to the keys map using the provided key.
// Finally, it returns the newly created limiter.
// Example usage:
//
//	provider := NewRateLimiterProvider(time.Second, 10)
//	limiter := provider.addKey("key")
//
// Note: The keys map and other related structures should be properly initialized before calling this method.
// The rateLimiterProvider type should be used to access this method.
// Other types should not have direct access to this method.
func (r *rateLimiterProvider) addKey(key string) *rate.Limiter {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	limiter := rate.NewLimiter(r.limit, r.burst)
	r.keys[key] = limiter

	return limiter
}
