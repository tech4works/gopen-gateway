package infra

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/external"
	appmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
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

func NewRateLimiterProvider(every time.Duration, limit int) external.RateLimiterProvider {
	return &rateLimiterProvider{
		keys:  map[string]*rate.Limiter{},
		mutex: &sync.RWMutex{},
		every: every,
		limit: rate.Every(every),
		burst: limit,
	}
}

func (r *rateLimiterProvider) Allow(key string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	limiter, exists := r.keys[key]
	if !exists {
		r.mutex.Unlock()
		limiter = r.addKey(key)
	}

	if !limiter.Allow() {
		return appmapper.NewErrTooManyRequests(r.burst, r.every)
	}

	return nil
}

func (r *rateLimiterProvider) addKey(key string) *rate.Limiter {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	limiter := rate.NewLimiter(r.limit, r.burst)
	r.keys[key] = limiter

	return limiter
}
