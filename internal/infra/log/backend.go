package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/todo/interfaces/log"
)

type backendConsole struct {
	backend     *vo.Backend
	httpRequest *vo.HttpRequest
}

func newBackendConsole(backend *vo.Backend, httpRequest *vo.HttpRequest) log.Console {
	return backendConsole{
		backend:     backend,
		httpRequest: httpRequest,
	}
}

func (b backendConsole) Infof(format string, msg ...any) {
	logger.Info(b.buildDefaultMsg(fmt.Sprintf(format, msg...)))
}

func (b backendConsole) Info(msg ...any) {
	logger.Info(b.buildDefaultMsg(helper.Sprintln(msg...)))
}

func (b backendConsole) Warnf(format string, msg ...any) {
	logger.Warn(b.buildDefaultMsg(fmt.Sprintf(format, msg...)))
}

func (b backendConsole) Warn(msg ...any) {
	logger.Warn(b.buildDefaultMsg(helper.Sprintln(msg...)))
}

func (b backendConsole) Errorf(format string, msg ...any) {
	logger.Error(b.buildDefaultMsg(fmt.Sprintf(format, msg...)))
}

func (b backendConsole) Error(msg ...any) {
	logger.Error(b.buildDefaultMsg(helper.Sprintln(msg...)))
}

func (b backendConsole) buildDefaultMsg(msg string) string {
	tag := fmt.Sprint(logger.StyleBold, "BACKEND", logger.StyleReset)
	path := b.backend.Path()
	traceID := BuildTraceIDText(b.httpRequest.TraceID())
	ip := b.httpRequest.ClientIP()

	return fmt.Sprintf("[%s] %s | %s | %s | %s", tag, path, traceID, ip, msg)
}
