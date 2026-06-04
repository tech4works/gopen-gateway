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

// NewHTTPLog cria uma instância de HTTPLog para mensagens de request/response HTTP.
func NewHTTPLog() app.HTTPLog {
	return httpLog{}
}

func (a httpLog) PrintRequest(ctx app.Context) {
	header := ctx.Request().Metadata()

	text := fmt.Sprintf(
		"Server received request userAgent=%s headerSize=%s",
		header.Get("User-Agent"),
		header.SizeStr(),
	)
	if ctx.Request().HasPayload() {
		payload := ctx.Request().Payload()
		text += fmt.Sprintf(" contentType=%s bodySize=%s",
			payload.ContentType().String(),
			payload.SizeInByteUnit(),
		)
	}

	PrintfCtx(ctx.Context(), InfoLevel, "REQ", a.prefix(ctx), "%s", text)
}

func (a httpLog) PrintResponse(ctx app.Context) {
	statusCode := BuildStatusCodeText(ctx.Response().Status())
	duration := ctx.Duration().Milliseconds()

	text := fmt.Sprintf("Server responded request statusCode=%s duration=%dms", statusCode, duration)
	PrintfCtx(ctx.Context(), InfoLevel, "RES", a.prefix(ctx), "%s", text)
}

func (a httpLog) prefix(ctx app.Context) string {
	path := ctx.Request().Path().Raw()
	traceID := BuildTraceIDText(ctx.Request().TraceID())
	ip := ctx.Request().ClientIP()

	method := BuildTintText(ctx.Request().Operation())
	url := BuildURIText(ctx.Request().Route())

	return fmt.Sprintf("[%s | %s | %s |%s| %s]", path, ip, traceID, method, url)
}
