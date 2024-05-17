package interfaces

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

type LoggerProvider interface {
	PrintBackendErrorf(backend *vo.Backend, format string, msg ...any)
	PrintBackendResponseInfo(backend *vo.Backend, httpBackendResponse *vo.HttpBackendResponse)
	PrintEndpointWarnf(ctx *api.Context, format string, msg ...any)
	PrintEndpointErrorf(ctx *api.Context, format string, msg ...any)
	PrintEndpointResponseInfo(ctx *api.Context)
	PrintHttpResponseInfo(ctx *api.Context)
}
