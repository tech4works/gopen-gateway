package infra

import (
	"golang.org/x/time/rate"
	"sync"
)

type IpRateLimiter struct {
	ips   map[string]*rate.Limiter
	mutex *sync.RWMutex
	r     rate.Limit
	b     int
}

func NewIpRateLimiter(r rate.Limit, b int) *IpRateLimiter {
	i := &IpRateLimiter{
		ips:   make(map[string]*rate.Limiter),
		mutex: &sync.RWMutex{},
		r:     r,
		b:     b,
	}
	return i
}

func (i *IpRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

func (i *IpRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mutex.Lock()
	limiter, exists := i.ips[ip]
	if !exists {
		i.mutex.Unlock()
		return i.AddIP(ip)
	}
	i.mutex.Unlock()
	return limiter
}
