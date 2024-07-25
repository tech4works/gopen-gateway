package vo

import (
	"time"
)

type CacheResponse struct {
	StatusCode StatusCode `json:"statusCode"`
	Header     Header     `json:"header"`
	Body       *Body      `json:"body,omitempty"`
	Duration   Duration   `json:"duration"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func NewCacheResponse(cacheConfig *Cache, response *HTTPResponse) *CacheResponse {
	return &CacheResponse{
		StatusCode: response.StatusCode(),
		Header:     response.Header(),
		Body:       response.Body(),
		Duration:   cacheConfig.Duration(),
		CreatedAt:  time.Now(),
	}
}

func (r CacheResponse) TTL() string {
	timeDuration := r.Duration.Time()
	sub := r.CreatedAt.Add(timeDuration).Sub(time.Now())
	return sub.String()
}
