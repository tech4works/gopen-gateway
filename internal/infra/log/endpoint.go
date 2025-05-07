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
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
)

type endpointLog struct {
	tag string
}

func NewEndpoint() app.EndpointLog {
	return endpointLog{
		tag: "END",
	}
}

func (e endpointLog) PrintInfof(executeData dto.ExecuteEndpoint, format string, msg ...any) {
	Printf(InfoLevel, e.tag, e.prefix(executeData), format, msg...)
}

func (e endpointLog) PrintInfo(executeData dto.ExecuteEndpoint, msg ...any) {
	Print(InfoLevel, e.tag, e.prefix(executeData), msg...)
}

func (e endpointLog) PrintWarnf(executeData dto.ExecuteEndpoint, format string, msg ...any) {
	Printf(WarnLevel, e.tag, e.prefix(executeData), format, msg...)
}

func (e endpointLog) PrintWarn(executeData dto.ExecuteEndpoint, msg ...any) {
	Print(WarnLevel, e.tag, e.prefix(executeData), msg...)
}

func (e endpointLog) PrintErrorf(executeData dto.ExecuteEndpoint, format string, msg ...any) {
	Printf(ErrorLevel, e.tag, e.prefix(executeData), format, msg...)
}

func (e endpointLog) PrintError(executeData dto.ExecuteEndpoint, msg ...any) {
	Print(ErrorLevel, e.tag, e.prefix(executeData), msg...)
}

func (e endpointLog) prefix(executeData dto.ExecuteEndpoint) string {
	path := executeData.Endpoint.Path()
	traceIDText := BuildTraceIDText(executeData.TraceID)

	method := BuildMethodText(executeData.Endpoint.Method())
	url := BuildUriText(executeData.Request.URL())

	return fmt.Sprintf("[%s | %s | %s |%s| %s]", path, executeData.ClientIP, traceIDText, method, url)
}
