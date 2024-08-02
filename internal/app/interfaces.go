package app

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"net/http"
	"time"
)

type Boot interface {
	Init() string
	Start(env string)
	Stop()
}

type BootLog interface {
	PrintLogo()
	PrintTitle(title string)
	PrintInfo(msg ...any)
	PrintInfof(format string, msg ...any)
	PrintWarn(msg ...any)
	PrintWarnf(format string, msg ...any)
	PrintError(msg ...any)
}

type HandlerFunc func(ctx Context)

type Router interface {
	Engine() http.Handler
	Handle(gopen *vo.Gopen, endpoint *vo.Endpoint, handles ...HandlerFunc)
}

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
	WriteString(code int, s string)
	WriteJson(code int, a any)
	WriteStatusCode(code int)
}

type HTTPClient interface {
	MakeRequest(ctx context.Context, endpoint *vo.Endpoint, request *vo.HTTPBackendRequest) *vo.HTTPBackendResponse
}

type EndpointLog interface {
	PrintInfof(endpoint *vo.Endpoint, traceID, clientIP, format string, msg ...any)
	PrintInfo(endpoint *vo.Endpoint, traceID, clientIP string, msg ...any)
	PrintWarnf(endpoint *vo.Endpoint, traceID, clientIP, format string, msg ...any)
	PrintWarn(endpoint *vo.Endpoint, traceID, clientIP string, msg ...any)
	PrintErrorf(endpoint *vo.Endpoint, traceID, clientIP, format string, msg ...any)
	PrintError(endpoint *vo.Endpoint, traceID, clientIP string, msg ...any)
}

type BackendLog interface {
	PrintInfof(executeData dto.ExecuteEndpoint, backend *vo.Backend, format string, msg ...any)
	PrintInfo(executeData dto.ExecuteEndpoint, backend *vo.Backend, msg ...any)
	PrintWarnf(executeData dto.ExecuteEndpoint, backend *vo.Backend, format string, msg ...any)
	PrintWarn(executeData dto.ExecuteEndpoint, backend *vo.Backend, msg ...any)
	PrintErrorf(executeData dto.ExecuteEndpoint, backend *vo.Backend, format string, msg ...any)
	PrintError(executeData dto.ExecuteEndpoint, backend *vo.Backend, msg ...any)
}
