package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type endpointLog struct {
}

func NewEndpoint() app.EndpointLog {
	return endpointLog{}
}

func (e endpointLog) PrintInfof(endpoint *vo.Endpoint, traceID, clientIP, format string, msg ...any) {
	logger.InfoOptsf(format, e.buildLoggerOptions(endpoint, traceID, clientIP), msg...)
}

func (e endpointLog) PrintInfo(endpoint *vo.Endpoint, traceID, clientIP string, msg ...any) {
	logger.InfoOpts(e.buildLoggerOptions(endpoint, traceID, clientIP), msg...)
}

func (e endpointLog) PrintWarnf(endpoint *vo.Endpoint, traceID, clientIP, format string, msg ...any) {
	logger.WarnOptsf(format, e.buildLoggerOptions(endpoint, traceID, clientIP), msg...)
}

func (e endpointLog) PrintWarn(endpoint *vo.Endpoint, traceID, clientIP string, msg ...any) {
	logger.WarnOpts(e.buildLoggerOptions(endpoint, traceID, clientIP), msg...)
}

func (e endpointLog) PrintErrorf(endpoint *vo.Endpoint, traceID, clientIP, format string, msg ...any) {
	logger.ErrorOptsf(format, e.buildLoggerOptions(endpoint, traceID, clientIP), msg...)
}

func (e endpointLog) PrintError(endpoint *vo.Endpoint, traceID, clientIP string, msg ...any) {
	logger.ErrorOpts(e.buildLoggerOptions(endpoint, traceID, clientIP), msg...)
}

func (e endpointLog) buildLoggerOptions(endpoint *vo.Endpoint, traceID, clientIP string) logger.Options {
	tag := fmt.Sprint(logger.StyleBold, "ENDPOINT", logger.StyleReset)
	path := endpoint.Path()

	return logger.Options{
		HideAllArgs:           true,
		CustomAfterPrefixText: fmt.Sprintf("[%s] %s | %s | %s |", tag, path, traceID, clientIP),
	}
}
