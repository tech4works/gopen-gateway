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
	"strings"
)

type publisherLog struct {
}

func NewPublisher() app.PublisherLog {
	return publisherLog{}
}

func (p publisherLog) PrintRequest(executeData dto.ExecuteEndpoint, publisher *vo.Publisher, message *vo.Message) {
	text := fmt.Sprintf("provider: %s | body.size: %s", publisher.Provider(),
		vo.NewBytesByInt(len([]byte(message.Body()))).String())
	Printf(InfoLevel, "PUB", p.prefix(executeData, publisher), text)
}

func (p publisherLog) prefix(executeData dto.ExecuteEndpoint, publisher *vo.Publisher) string {
	path := p.extractQueueName(publisher.Reference())
	if checker.IsEmpty(path) {
		path = p.extractTopicName(publisher.Reference())
	}

	traceID := BuildTraceIDText(executeData.TraceID)
	ip := executeData.ClientIP

	url := BuildUriText(publisher.Reference())

	return fmt.Sprintf("[%s | %s | %s | %s]", path, ip, traceID, url)
}

func (p publisherLog) extractTopicName(topicArn string) string {
	parts := strings.Split(topicArn, ":")
	if checker.IsEmpty(parts) {
		return ""
	}
	return parts[len(parts)-1]
}

func (p publisherLog) extractQueueName(queueURL string) string {
	parts := strings.Split(queueURL, "/")
	if checker.IsEmpty(parts) {
		return ""
	}
	return parts[len(parts)-1]
}
