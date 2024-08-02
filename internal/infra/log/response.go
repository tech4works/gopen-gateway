package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/tech4works/gopen-gateway/internal/app"
)

func PrintResponse(ctx app.Context) {
	tag := fmt.Sprint(logger.StyleBold, "RESPONSE", logger.StyleReset)
	path := ctx.Endpoint().Path()
	traceID := BuildTraceIDText(ctx.TraceID())
	ip := ctx.ClientIP()
	statusCode := BuildStatusCodeText(ctx.Response().StatusCode())
	latency := ctx.Latency().String()
	method := BuildMethodText(ctx.Request().Method())
	url := BuildUriText(ctx.Request().Url())

	m := fmt.Sprintf("[%s] %s | %s | %s |%s| %s |%s| %s", tag, path, traceID, ip, statusCode, latency, method, url)
	logger.InfoOpts(logger.Options{HideAllArgs: true}, m)
}
