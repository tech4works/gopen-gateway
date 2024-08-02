package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type backendLog struct {
}

func NewBackend() app.BackendLog {
	return backendLog{}
}

func (b backendLog) PrintInfof(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	request *vo.HTTPBackendRequest, format string, msg ...any) {
	logger.InfoOptsf(format, b.buildLoggerOptions(executeData, backend, request), msg...)
}

func (b backendLog) PrintInfo(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	request *vo.HTTPBackendRequest, msg ...any) {
	logger.InfoOpts(b.buildLoggerOptions(executeData, backend, request), msg...)
}

func (b backendLog) PrintWarnf(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	request *vo.HTTPBackendRequest, format string, msg ...any) {
	logger.WarnOptsf(format, b.buildLoggerOptions(executeData, backend, request), msg...)
}

func (b backendLog) PrintWarn(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	request *vo.HTTPBackendRequest, msg ...any) {
	logger.WarnOpts(b.buildLoggerOptions(executeData, backend, request), msg...)
}

func (b backendLog) PrintErrorf(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	request *vo.HTTPBackendRequest, format string, msg ...any) {
	logger.ErrorOptsf(format, b.buildLoggerOptions(executeData, backend, request), msg...)
}

func (b backendLog) PrintError(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	request *vo.HTTPBackendRequest, msg ...any) {
	logger.ErrorOpts(b.buildLoggerOptions(executeData, backend, request), msg...)
}

func (b backendLog) buildLoggerOptions(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	request *vo.HTTPBackendRequest) logger.Options {
	tag := BuildTagText(backend.Type())
	path := backend.Path()
	traceID := BuildTraceIDText(executeData.TraceID)
	ip := executeData.ClientIP

	method := BuildMethodText(request.Method())
	url := BuildUriText(request.Url())

	return logger.Options{
		HideArgDatetime:       true,
		HideArgCaller:         true,
		CustomAfterPrefixText: fmt.Sprintf("[%s] (%s | %s | %s |%s| %s)", tag, path, traceID, ip, method, url),
	}
}
