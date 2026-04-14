/*
 * Copyright 2024 Tech4Works
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	header := ctx.Request().Metadata()

	text := fmt.Sprintf("header.user-agent: %s | header.size: %s", header.Get("User-Agent"), header.SizeStr())
	if ctx.Request().HasPayload() {
		payload := ctx.Request().Payload()
		text += fmt.Sprintf(" | body.content-type: %s | body.size: %s", payload.ContentType().String(),
			payload.SizeInByteUnit())
	}

	Print(InfoLevel, "REQ", a.prefix(ctx), text)
}

func (a httpLog) PrintResponse(ctx app.Context) {
	statusCode := BuildStatusCodeText(ctx.Response().Status())
	duration := ctx.Duration().Milliseconds()

	prefix := a.prefix(ctx)

	Printf(InfoLevel, "RES", prefix, "status:%s| duration: %vms", statusCode, duration)
}

func (a httpLog) prefix(ctx app.Context) string {
	path := ctx.Request().Path().Raw()
	traceID := BuildTraceIDText(ctx.Request().TraceID())
	ip := ctx.Request().ClientIP()

	method := BuildTintText(ctx.Request().Operation())
	url := BuildURIText(ctx.Request().Route())

	return fmt.Sprintf("[%s | %s | %s |%s| %s]", path, ip, traceID, method, url)
}
