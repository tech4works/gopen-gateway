package external

import (
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

type CacheProvider interface {
	CacheStrategyHandler() gin.HandlerFunc
}
