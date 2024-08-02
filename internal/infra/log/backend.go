package log

import (
	"fmt"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type backendLog struct {
}

func NewBackend() app.BackendLog {
	return backendLog{}
}

func (b backendLog) PrintInfof(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest,
	format string, msg ...any) {
	Printf(InfoLevel, backend.Type().Abbreviation(), b.prefix(executeData, backend, request), format, msg...)
}

func (b backendLog) PrintInfo(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest,
	msg ...any) {
	Print(InfoLevel, backend.Type().Abbreviation(), b.prefix(executeData, backend, request), msg...)
}

func (b backendLog) PrintWarnf(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest,
	format string, msg ...any) {
	Printf(WarnLevel, backend.Type().Abbreviation(), b.prefix(executeData, backend, request), format, msg...)
}

func (b backendLog) PrintWarn(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest,
	msg ...any) {
	Print(WarnLevel, backend.Type().Abbreviation(), b.prefix(executeData, backend, request), msg...)
}

func (b backendLog) PrintErrorf(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest,
	format string, msg ...any) {
	Printf(ErrorLevel, backend.Type().Abbreviation(), b.prefix(executeData, backend, request), format, msg...)
}

func (b backendLog) PrintError(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest,
	msg ...any) {
	Print(ErrorLevel, backend.Type().Abbreviation(), b.prefix(executeData, backend, request), msg...)
}

func (b backendLog) prefix(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest) string {
	path := backend.Path()
	traceID := BuildTraceIDText(executeData.TraceID)
	ip := executeData.ClientIP

	method := BuildMethodText(request.Method())
	url := BuildUriText(request.Url())

	return fmt.Sprintf("%s (%s | %s |%s| %s)", path, traceID, ip, method, url)
}
