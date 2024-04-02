package infra

import (
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type rateLimiterProvider struct {
	keys  map[string]*rate.Limiter
	mutex *sync.RWMutex
	every time.Duration
	limit rate.Limit
	burst int
}

type RateLimiterProvider interface {
	Allow(key string) error
}

// NewRateLimiterProvider creates a new instance of the RateLimiterProvider interface.
// It takes a time duration 'every' and an integer 'limit' as parameters.
// It returns a RateLimiterProvider object.
func NewRateLimiterProvider(every time.Duration, limit int) RateLimiterProvider {
	return &rateLimiterProvider{
		keys:  map[string]*rate.Limiter{},
		mutex: &sync.RWMutex{},
		every: every,
		limit: rate.Every(every),
		burst: limit,
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
// Example usage:
//
//	provider := NewRateLimiterProvider(time.Second, 10)
//	err := provider.Allow("key")
//	if err != nil {
//	  // handle error
//	}
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
