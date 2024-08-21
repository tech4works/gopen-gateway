package log

import (
	"fmt"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"time"
)

type backendLog struct {
}

func NewBackend() app.BackendLog {
	return backendLog{}
}

func (b backendLog) PrintRequest(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest) {
	text := fmt.Sprintf("REQ header.user-agent: %s | header.size: %s", request.Header().Get("User-Agent"),
		request.Header().SizeStr())
	if request.HasBody() {
		body := request.Body()
		text += fmt.Sprintf(" | body.content-type: %s | body.size: %s", body.ContentType().String(), body.SizeInByteUnit())
	}

	Printf(InfoLevel, backend.Type().Abbreviation(), b.prefix(executeData, backend, request), text)
}

func (b backendLog) PrintResponse(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest,
	response *vo.HTTPBackendResponse, duration time.Duration) {
	statusCode := response.StatusCode()
	statusCodeText := BuildStatusCodeText(statusCode)

	Printf(InfoLevel, backend.Type().Abbreviation(), b.prefix(executeData, backend, request),
		"RES status-code:%v| duration: %vms", statusCodeText, duration.Milliseconds())
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
	url := BuildUriText(request.FullPath())

	return fmt.Sprintf("[%s | %s | %s |%s| %s]", path, ip, traceID, method, url)
}
