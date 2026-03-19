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

package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type BackendPublisherConfig struct {
	broker          enum.BackendBroker
	path            string
	groupID         string
	deduplicationID string
	delay           Duration
	message         BackendPublisherMessageConfig
}

func NewBackendPublisherConfig(
	broker enum.BackendBroker,
	path,
	groupID,
	deduplicationID string,
	delay Duration,
	message BackendPublisherMessageConfig,
) *BackendPublisherConfig {
	return &BackendPublisherConfig{
		broker:          broker,
		path:            path,
		groupID:         groupID,
		deduplicationID: deduplicationID,
		delay:           delay,
		message:         message,
	}
}

func (p BackendPublisherConfig) Broker() enum.BackendBroker {
	return p.broker
}

func (p BackendPublisherConfig) Path() string {
	return p.path
}

func (p BackendPublisherConfig) HasGroupID() bool {
	return checker.IsNotEmpty(p.groupID)
}

func (p BackendPublisherConfig) GroupID() string {
	return p.groupID
}

func (p BackendPublisherConfig) HasDeduplicationID() bool {
	return checker.IsNotEmpty(p.deduplicationID)
}
func (p BackendPublisherConfig) DeduplicationID() string {
	return p.deduplicationID
}

func (p BackendPublisherConfig) Delay() Duration {
	return p.delay
}

func (p BackendPublisherConfig) CountAllDataTransforms() (count int) {
	if p.Message().HasBody() {
		count += p.Message().Body().CountDataTransforms()
	}
	if p.Message().HasAttributes() {
		count += len(p.Message().Attributes())
	}

	return count
}

func (p BackendPublisherConfig) Message() BackendPublisherMessageConfig {
	return p.message
}
