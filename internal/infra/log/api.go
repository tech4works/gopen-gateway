package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
)

type httpLog struct {
}

func NewHTTPLog() app.HTTPLog {
	return httpLog{}
}

func (a httpLog) PrintRequest(ctx app.Context) {
	contentType := ctx.Request().Header().Get(mapper.ContentType)
	contentEncoding := ctx.Request().Header().Get(mapper.ContentEncoding)
	headerSize := ctx.Request().Header().Size()
	params := ctx.Request().Params().Len()
	queries := ctx.Request().Query().Len()
	body := 0
	if ctx.Request().HasBody() {
		body = ctx.Request().Body().Len()
	}

	opts := buildLoggerOptions(ctx, "REQUEST")

	logger.InfoOptsf("content-type: %s | content-encoding: %s | header: %v | params: %v | query: %v | body: %s",
		opts, contentType, contentEncoding, headerSize, params, queries, body)
}

func (a httpLog) PrintResponse(ctx app.Context) {
	statusCode := BuildStatusCodeText(ctx.Response().StatusCode())
	latency := ctx.Latency().String()

	opts := buildLoggerOptions(ctx, "RESPONSE")

	logger.InfoOptsf("status-code:%s| latency: %s", opts, statusCode, latency)
}

func buildLoggerOptions(ctx app.Context, nameTag string) logger.Options {
	tag := BuildTagText(nameTag)
	path := ctx.Endpoint().Path()
	traceID := BuildTraceIDText(ctx.TraceID())
	ip := ctx.ClientIP()

	method := BuildMethodText(ctx.Request().Method())
	url := BuildUriText(ctx.Request().Url())

	prefix := fmt.Sprintf("[%s] (%s | %s | %s |%s| %s)", tag, path, ip, traceID, method, url)
	return logger.Options{
		HideArgDatetime:       true,
		HideArgCaller:         true,
		CustomAfterPrefixText: prefix,
	}
}
