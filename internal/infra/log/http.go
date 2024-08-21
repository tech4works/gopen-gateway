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
	header := ctx.Request().Header()

	text := fmt.Sprintf("header.user-agent: %s | header.size: %s", header.Get("User-Agent"), header.SizeStr())
	if ctx.Request().HasBody() {
		body := ctx.Request().Body()
		text += fmt.Sprintf(" | body.content-type: %s | body.size: %s", body.ContentType().String(), body.SizeInByteUnit())
	}

	Print(InfoLevel, "REQ", a.prefix(ctx), text)
}

func (a httpLog) PrintResponse(ctx app.Context) {
	statusCode := BuildStatusCodeText(ctx.Response().StatusCode())
	duration := ctx.Duration().Milliseconds()

	prefix := a.prefix(ctx)

	Printf(InfoLevel, "RES", prefix, "status-code:%s| duration: %vms", statusCode, duration)
}

func (a httpLog) prefix(ctx app.Context) string {
	path := ctx.Request().Path().Raw()
	traceID := BuildTraceIDText(ctx.TraceID())
	ip := ctx.ClientIP()

	method := BuildMethodText(ctx.Request().Method())
	url := BuildUriText(ctx.Request().Url())

	return fmt.Sprintf("[%s | %s | %s |%s| %s]", path, ip, traceID, method, url)
}
