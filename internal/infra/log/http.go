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
	method := ctx.Request().Operation()
	url := ctx.Request().Route()

	if DebugLevel.Allowed() {
		header := ctx.Request().Metadata().String()
		body := "<nil>"
		if ctx.Request().HasPayload() {
			if s, err := ctx.Request().Payload().CompactString(); err == nil {
				body = s
			}
		}
		text := fmt.Sprintf("Server received request method=%s url=%s header=%s body=%s",
			BuildTintText(method), BuildURIText(url), header, body)
		PrintfCtx(ctx.Context(), DebugLevel, "REQ", a.prefix(ctx), "%s", text)
	} else {
		text := fmt.Sprintf("Server received request method=%s url=%s",
			BuildTintText(method), BuildURIText(url))
		if ctx.Request().HasPayload() {
			payload := ctx.Request().Payload()
			text += fmt.Sprintf(" content_type=%s body_size=%s",
				payload.ContentType().String(),
				payload.SizeInByteUnit(),
			)
		}
		PrintfCtx(ctx.Context(), InfoLevel, "REQ", a.prefix(ctx), "%s", text)
	}
}

func (a httpLog) PrintResponse(ctx app.Context) {
	method := ctx.Request().Operation()
	url := ctx.Request().Route()
	status := ctx.Response().Status()
	duration := ctx.Duration().Milliseconds()

	if DebugLevel.Allowed() {
		header := ctx.Response().Metadata().String()
		body := "<nil>"
		if ctx.Response().HasPayload() {
			if s, err := ctx.Response().Payload().CompactString(); err == nil {
				body = s
			}
		}
		text := fmt.Sprintf("Server responded request method=%s url=%s status_code=%s duration=%dms header=%s body=%s",
			BuildTintText(method), BuildURIText(url), BuildStatusCodeText(status), duration, header, body)
		PrintfCtx(ctx.Context(), DebugLevel, "RES", a.prefix(ctx), "%s", text)
	} else {
		text := fmt.Sprintf("Server responded request method=%s url=%s status_code=%s duration=%dms",
			BuildTintText(method), BuildURIText(url), BuildStatusCodeText(status), duration)
		if ctx.Response().HasPayload() {
			payload := ctx.Response().Payload()
			text += fmt.Sprintf(" content_type=%s body_size=%s",
				payload.ContentType().String(),
				payload.SizeInByteUnit(),
			)
		}
		PrintfCtx(ctx.Context(), InfoLevel, "RES", a.prefix(ctx), "%s", text)
	}
}

func (a httpLog) prefix(ctx app.Context) string {
	path := ctx.Request().Path().Raw()
	traceID := BuildTraceIDText(ctx.Request().TraceID())
	ip := ctx.Request().ClientIP()

	return fmt.Sprintf("[%s | %s | %s]", path, ip, traceID)
}
