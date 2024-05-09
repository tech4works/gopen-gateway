package vo

import (
	"time"
)

// CacheResponse represents a cached HTTP httpResponse.
// It contains the status code, header, body, duration, and creation timestamp of the httpResponse.
// The duration specifies how long the httpResponse should be cached.
// The CreatedAt field indicates the timestamp of the httpResponse's creation.
type CacheResponse struct {
	// StatusCode is an integer field representing the status code of an HTTP httpResponse.
	// It is included in the CacheResponse struct and is used to store the status code of a cached httpResponse.
	StatusCode int `json:"statusCode"`
	// Header is a field representing the header of an HTTP httpResponse.
	// It is included in the CacheResponse struct and is used to store the header of a cached httpResponse.
	Header Header `json:"header"`
	// Body is a field representing the body of an HTTP httpResponse. (optional)
	// It is included in the CacheResponse struct and is used to store the body of a cached httpResponse.
	Body *CacheBody `json:"body,omitempty"`
	// Duration represents the duration for which the httpResponse should be cached.
	Duration Duration `json:"duration"`
	// CreatedAt is a field of the CacheResponse struct indicating the timestamp of the httpResponse's creation.
	CreatedAt time.Time `json:"createdAt"`
}

func NewCacheResponse(httpResponse *HttpResponse, duration Duration) *CacheResponse {
	return &CacheResponse{
		StatusCode: httpResponse.StatusCode(),
		Header:     httpResponse.Header(),
		Body:       newCacheBody(httpResponse.Body()),
		Duration:   duration,
		CreatedAt:  time.Now(),
	}
}

// TTL calculates the time to live (TTL) for the CacheResponse object.
// It subtracts the current time from the sum of the CreatedAt time and the Duration of the CacheResponse.
// Returns the TTL duration as a string representation.
func (c CacheResponse) TTL() string {
	timeDuration := c.Duration.Time()
	sub := c.CreatedAt.Add(timeDuration).Sub(time.Now())
	return sub.String()
}
