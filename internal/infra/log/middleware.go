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

type middlewareLog struct {
	tag string
}

// NewMiddleware cria uma instância de MiddlewareLog para mensagens de interceptors.
func NewMiddleware() app.MiddlewareLog {
	return middlewareLog{
		tag: "MID",
	}
}

func (e middlewareLog) PrintInfof(ctx app.Context, format string, msg ...any) {
	PrintfCtx(ctx.Context(), InfoLevel, e.tag, e.prefix(ctx), format, msg...)
}

func (e middlewareLog) PrintInfo(ctx app.Context, msg ...any) {
	PrintCtx(ctx.Context(), InfoLevel, e.tag, e.prefix(ctx), msg...)
}

func (e middlewareLog) PrintWarnf(ctx app.Context, format string, msg ...any) {
	PrintfCtx(ctx.Context(), WarnLevel, e.tag, e.prefix(ctx), format, msg...)
}

func (e middlewareLog) PrintWarn(ctx app.Context, msg ...any) {
	PrintCtx(ctx.Context(), WarnLevel, e.tag, e.prefix(ctx), msg...)
}

func (e middlewareLog) PrintErrorf(ctx app.Context, format string, msg ...any) {
	PrintfCtx(ctx.Context(), ErrorLevel, e.tag, e.prefix(ctx), format, msg...)
}

func (e middlewareLog) PrintError(ctx app.Context, msg ...any) {
	PrintCtx(ctx.Context(), ErrorLevel, e.tag, e.prefix(ctx), msg...)
}

func (e middlewareLog) prefix(ctx app.Context) string {
	path := ctx.Endpoint().Path()
	traceIDText := BuildTraceIDText(ctx.Request().TraceID())
	method := BuildTintText(ctx.Endpoint().Method())
	url := BuildURIText(ctx.Request().Route())

	return fmt.Sprintf("[%s | %s | %s |%s| %s]", path, ctx.Request().ClientIP(), traceIDText, method, url)
}
