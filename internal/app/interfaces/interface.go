package interfaces

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type TraceProvider interface {
	GenerateTraceId() string
}

type LogProvider interface {
	InitializeLoggerOptions(ctx *gin.Context)
	BuildInitialRequestMessage(ctx *gin.Context) string
	BuildFinishRequestMessage(writer dto.ResponseWriter, startTime time.Time) string
}

type RateLimiterProvider interface {
	Allow(key string) error
}

type SizeLimiterProvider interface {
	Allow(request *http.Request) error
}

type CacheStore interface {
	// Set sets an item to the Cache, replacing any existing item.
	Set(ctx context.Context, key string, value any, expire time.Duration) error
	// Del removes an item from the Cache. Does nothing if the key is not in the Cache.
	Del(ctx context.Context, key string) error
	// Get retrieves an item from the Cache. if key does not exist in the store, return ErrCacheMiss
	Get(ctx context.Context, key string, dest any) error
	// Close client connection
	Close() error
}
