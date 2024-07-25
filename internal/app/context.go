package app

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/opentracing/opentracing-go"
	"time"
)

type Context interface {
	Context() context.Context
	WithContext(ctx context.Context)
	Done() <-chan struct{}
	Next()
	Latency() time.Duration
	Span() opentracing.Span
	TraceID() string
	ClientIP() string
	Gopen() *vo.Gopen
	Endpoint() *vo.Endpoint
	Request() *vo.HTTPRequest
	Response() *vo.HTTPResponse
	Write(response *vo.HTTPResponse)
	WriteCacheResponse(cacheResponse *vo.CacheResponse)
	WriteError(code int, err error)
}
