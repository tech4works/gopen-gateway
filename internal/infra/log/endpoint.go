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
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type endpointLog struct {
	tag string
}

func NewEndpoint() app.EndpointLog {
	return endpointLog{
		tag: "END",
	}
}

func (e endpointLog) PrintInfof(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID, format string,
	msg ...any) {
	Printf(InfoLevel, e.tag, e.prefix(endpoint, request, clientIP, traceID), format, msg...)
}

func (e endpointLog) PrintInfo(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID string, msg ...any) {
	Print(InfoLevel, e.tag, e.prefix(endpoint, request, clientIP, traceID), msg...)
}

func (e endpointLog) PrintWarnf(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID, format string,
	msg ...any) {
	Printf(WarnLevel, e.tag, e.prefix(endpoint, request, clientIP, traceID), format, msg...)
}

func (e endpointLog) PrintWarn(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID string, msg ...any) {
	Print(WarnLevel, e.tag, e.prefix(endpoint, request, clientIP, traceID), msg...)
}

func (e endpointLog) PrintErrorf(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID, format string,
	msg ...any) {
	Printf(ErrorLevel, e.tag, e.prefix(endpoint, request, clientIP, traceID), format, msg...)
}

func (e endpointLog) PrintError(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID string, msg ...any) {
	Print(ErrorLevel, e.tag, e.prefix(endpoint, request, clientIP, traceID), msg...)
}

func (e endpointLog) prefix(endpoint *vo.Endpoint, request *vo.HTTPRequest, clientIP, traceID string) string {
	path := endpoint.Path()
	traceIDText := BuildTraceIDText(traceID)

	method := BuildMethodText(request.Method())
	url := BuildUriText(request.Url())

	return fmt.Sprintf("[%s | %s | %s |%s| %s]", path, clientIP, traceIDText, method, url)
}
