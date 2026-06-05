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

	"github.com/tech4works/checker"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type backendLog struct {
}

// NewBackend cria uma instância de BackendLog para mensagens de backends HTTP e publisher.
func NewBackend() app.BackendLog {
	return backendLog{}
}

func (b backendLog) PrintHTTPRequest(
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	request *vo.HTTPBackendRequest,
) {
	method := BuildTintText(request.Method())
	url := BuildURIText(request.URL())

	if DebugLevel.Allowed() {
		header := request.Header().String()
		body := "<nil>"
		if request.HasBody() {
			if s, err := request.Body().CompactString(); err == nil {
				body = s
			}
		}
		text := fmt.Sprintf("Backend HTTP request started method=%s url=%s header=%s body=%s",
			method, url, header, body)
		Print(DebugLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), text)
	} else {
		text := fmt.Sprintf("Backend HTTP request started method=%s url=%s", method, url)
		if request.HasBody() {
			body := request.Body()
			text += fmt.Sprintf(" content_type=%s body_size=%s",
				body.ContentType().String(),
				body.SizeInByteUnit(),
			)
		}
		Print(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), text)
	}
}

func (b backendLog) PrintPublisherRequest(
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	request *vo.PublisherBackendRequest,
) {
	broker := request.Broker().String()
	path := request.Path()

	if DebugLevel.Allowed() {
		body := "<nil>"
		if request.HasBody() {
			if s, err := request.Body().CompactString(); err == nil {
				body = s
			}
		}
		text := fmt.Sprintf("Backend publisher request started broker=%s path=%s body=%s", broker, path, body)
		if checker.NonNil(request.GroupID()) {
			text += fmt.Sprintf(" group_id=%s", *request.GroupID())
		}
		if checker.NonNil(request.DeduplicationID()) {
			text += fmt.Sprintf(" deduplication_id=%s", *request.DeduplicationID())
		}
		if checker.IsGreaterThan(request.Delay().Time().Milliseconds(), 0) {
			text += fmt.Sprintf(" delay=%dms", request.Delay().Time().Milliseconds())
		}
		Print(DebugLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), text)
	} else {
		text := fmt.Sprintf("Backend publisher request started broker=%s path=%s", broker, path)
		if request.HasBody() {
			text += fmt.Sprintf(" body_size=%s", request.Body().SizeInByteUnit())
		}
		if checker.NonNil(request.GroupID()) {
			text += fmt.Sprintf(" group_id=%s", *request.GroupID())
		}
		if checker.NonNil(request.DeduplicationID()) {
			text += fmt.Sprintf(" deduplication_id=%s", *request.DeduplicationID())
		}
		if checker.IsGreaterThan(request.Delay().Time().Milliseconds(), 0) {
			text += fmt.Sprintf(" delay=%dms", request.Delay().Time().Milliseconds())
		}
		Print(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), text)
	}
}

func (b backendLog) PrintResponse(
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	response *vo.BackendResponse,
) {
	statusCode := BuildStatusCodeText(response.Status())
	duration := response.Duration().Milliseconds()

	if DebugLevel.Allowed() {
		header := response.Metadata().String()
		body := "<nil>"
		if response.HasBody() {
			if s, err := response.Payload().CompactString(); err == nil {
				body = s
			}
		}
		text := fmt.Sprintf("Backend %s response received status_code=%s ok=%v duration=%dms header=%s body=%s",
			backend.Kind(), statusCode, response.OK(), duration, header, body)
		Print(DebugLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), text)
	} else {
		text := fmt.Sprintf("Backend %s response received status_code=%s ok=%v duration=%dms",
			backend.Kind(), statusCode, response.OK(), duration)
		if response.HasBody() {
			body := response.Payload()
			text += fmt.Sprintf(" content_type=%s body_size=%s",
				body.ContentType().String(),
				body.SizeInByteUnit(),
			)
		}
		Print(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), text)
	}
}

func (b backendLog) PrintInfof(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, format string, msg ...any) {
	Printf(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), format, msg...)
}

func (b backendLog) PrintInfo(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, msg ...any) {
	Print(InfoLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), msg...)
}

func (b backendLog) PrintWarnf(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, format string, msg ...any) {
	Printf(WarnLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), format, msg...)
}

func (b backendLog) PrintWarn(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, msg ...any) {
	Print(WarnLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), msg...)
}

func (b backendLog) PrintErrorf(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, format string, msg ...any) {
	Printf(ErrorLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), format, msg...)
}

func (b backendLog) PrintError(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, msg ...any) {
	Print(ErrorLevel, backend.Flow().Abbreviation(), b.prefix(executeData, backend), msg...)
}

func (b backendLog) prefix(executeData dto.ExecuteEndpoint, backend *vo.BackendConfig) string {
	id := backend.ID()
	traceID := BuildTraceIDText(executeData.Request.TraceID())
	ip := executeData.Request.ClientIP()

	var tintText string
	if backend.IsHTTP() {
		tintText = BuildTintText(backend.HTTP().Method())
	} else if backend.IsPublisher() {
		tintText = BuildTintText(backend.Publisher().Broker().String())
	}

	return fmt.Sprintf("[%s | %s | %s |%s]", id, ip, traceID, tintText)
}
