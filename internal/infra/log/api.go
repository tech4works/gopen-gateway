package log

import (
	"fmt"
	"github.com/tech4works/gopen-gateway/internal/app"
)

type httpLog struct {
}

func NewHTTPLog() app.HTTPLog {
	return httpLog{}
}

func (a httpLog) PrintRequest(ctx app.Context) {
	headerSize := ctx.Request().Header().Size()
	params := ctx.Request().Params().Len()
	queries := ctx.Request().Query().Len()
	body := 0
	if ctx.Request().HasBody() {
		body = ctx.Request().Body().Len()
	}

	prefix := a.prefix(ctx)
	format := "header: %v | params: %v | query: %v | body: %v"

	Printf(InfoLevel, "REQ", prefix, format, headerSize, params, queries, body)
}

func (a httpLog) PrintResponse(ctx app.Context) {
	statusCode := BuildStatusCodeText(ctx.Response().StatusCode())
	latency := ctx.Latency().String()

	prefix := a.prefix(ctx)

	Printf(InfoLevel, "RES", prefix, "status-code:%s| latency: %s", statusCode, latency)
}

func (a httpLog) prefix(ctx app.Context) string {
	path := ctx.Endpoint().Path()
	traceID := BuildTraceIDText(ctx.TraceID())
	ip := ctx.ClientIP()

	method := BuildMethodText(ctx.Request().Method())
	url := BuildUriText(ctx.Request().Url())

	return fmt.Sprintf("%s (%s | %s |%s| %s)", path, ip, traceID, method, url)
}
