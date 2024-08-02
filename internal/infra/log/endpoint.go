package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/todo/interfaces/log"
)

type endpointConsole struct {
	endpoint *vo.Endpoint
	traceID  string
	clientIP string
}

func newEndpointConsole(endpoint *vo.Endpoint, traceID, clientIP string) log.Console {
	return endpointConsole{
		endpoint: endpoint,
		traceID:  traceID,
		clientIP: clientIP,
	}
}

func (e endpointConsole) Infof(format string, msg ...any) {
	logger.Info(e.buildDefaultMsg(fmt.Sprintf(format, msg...)))
}

func (e endpointConsole) Info(msg ...any) {
	logger.Info(e.buildDefaultMsg(helper.Sprintln(msg...)))
}

func (e endpointConsole) Warnf(format string, msg ...any) {
	logger.Warn(e.buildDefaultMsg(fmt.Sprintf(format, msg...)))
}

func (e endpointConsole) Warn(msg ...any) {
	logger.Warn(e.buildDefaultMsg(helper.Sprintln(msg...)))
}

func (e endpointConsole) Errorf(format string, msg ...any) {
	logger.Error(e.buildDefaultMsg(fmt.Sprintf(format, msg...)))
}

func (e endpointConsole) Error(msg ...any) {
	logger.Error(e.buildDefaultMsg(helper.Sprintln(msg...)))
}

func (e endpointConsole) buildDefaultMsg(msg string) string {
	tag := fmt.Sprint(logger.StyleBold, "ENDPOINT", logger.StyleReset)
	path := e.endpoint.Path()
	traceID := BuildTraceIDText(e.httpRequest.TraceID())
	ip := e.httpRequest.ClientIP()

	return fmt.Sprintf("[%s] %s | %s | %s | %s", tag, path, traceID, ip, msg)
}
