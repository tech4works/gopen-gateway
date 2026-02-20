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
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type backendLog struct {
}

func NewBackend() app.BackendLog {
	return backendLog{}
}

func (b backendLog) PrintHTTPRequest(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.HTTPBackendRequest) {
	text := fmt.Sprintf("HTTP REQ url: %s | header.user-agent: %s | header.size: %s",
		BuildUriText(request.FullPath()), request.Header().Get("User-Agent"), request.Header().SizeStr())
	if request.HasBody() {
		body := request.Body()
		text += fmt.Sprintf(" | body.content-type: %s | body.size: %s", body.ContentType().String(), body.SizeInByteUnit())
	}

	Printf(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), text)
}

func (b backendLog) PrintHTTPResponse(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	response *vo.HTTPBackendResponse, duration time.Duration) {
	if checker.IsNil(response) {
		return
	}

	statusCode := response.StatusCode()
	statusCodeText := BuildStatusCodeText(statusCode)

	Printf(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend),
		"HTTP RES status-code:%v| duration: %vms", statusCodeText, duration.Milliseconds())
}

func (b backendLog) PrintPublisherRequest(executeData dto.ExecuteEndpoint, backend *vo.Backend, request *vo.PublisherBackendRequest) {
	text := fmt.Sprintf("PUBLISHER REQ body.size: %s", vo.NewBytesByInt(len([]byte(request.Body()))).String())
	if checker.NonNil(request.GroupID()) {
		text += fmt.Sprintf(" | group-id: %s", *request.GroupID())
	}
	if checker.NonNil(request.DeduplicationID()) {
		text += fmt.Sprintf(" | deduplication-id: %s", *request.DeduplicationID())
	}
	if checker.IsGreaterThan(request.Delay().Time().Milliseconds(), 0) {
		text += fmt.Sprintf(" | delay: %vms", request.Delay().Time().Milliseconds())
	}

	Print(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), text)
}

func (b backendLog) PrintPublisherResponse(
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	response *vo.PublisherBackendResponse,
	duration time.Duration,
) {
	text := fmt.Sprintf("PUBLISHER RES ok: %v | duration: %vms", response.OK(), duration.Milliseconds())

	Printf(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), text)
}

func (b backendLog) PrintInfof(executeData dto.ExecuteEndpoint, backend *vo.Backend, format string, msg ...any) {
	Printf(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), format, msg...)
}

func (b backendLog) PrintInfo(executeData dto.ExecuteEndpoint, backend *vo.Backend, msg ...any) {
	Print(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), msg...)
}

func (b backendLog) PrintWarnf(executeData dto.ExecuteEndpoint, backend *vo.Backend, format string, msg ...any) {
	Printf(WarnLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), format, msg...)
}

func (b backendLog) PrintWarn(executeData dto.ExecuteEndpoint, backend *vo.Backend, msg ...any) {
	Print(WarnLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), msg...)
}

func (b backendLog) PrintErrorf(executeData dto.ExecuteEndpoint, backend *vo.Backend, format string, msg ...any) {
	Printf(ErrorLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), format, msg...)
}

func (b backendLog) PrintError(executeData dto.ExecuteEndpoint, backend *vo.Backend, msg ...any) {
	Print(ErrorLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), msg...)
}

func (b backendLog) prefix(executeData dto.ExecuteEndpoint, backend *vo.Backend) string {
	id := backend.ID()

	var tintText string
	if backend.IsHTTP() {
		tintText = backend.HTTP().Method()
	} else if backend.IsPublisher() {
		tintText = backend.Publisher().Broker().String()
	}

	traceID := BuildTraceIDText(executeData.TraceID)
	ip := executeData.ClientIP
	tintText = BuildTintText(tintText)

	return fmt.Sprintf("[%s | %s | %s |%s]", id, ip, traceID, tintText)
}
