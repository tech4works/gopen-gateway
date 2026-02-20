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

type Publisher struct {
	onlyIf          []string
	ignoreIf        []string
	broker          enum.BackendBroker
	path            string
	groupID         string
	deduplicationID string
	delay           Duration
	async           bool
	message         *PublisherMessage
}

type PublisherMessage struct {
	continueOnError bool
	onlyIf          []string
	ignoreIf        []string
	attributes      map[string]PublisherMessageAttribute
	body            *PublisherMessageBody
}

type PublisherMessageAttribute struct {
	dataType string
	value    string
}

type PublisherMessageBody struct {
	omitEmpty bool
	mapper    *Mapper
	projector *Projector
	modifiers []Modifier
}

func NewPublisher(
	onlyIf,
	ignoreIf []string,
	broker enum.BackendBroker,
	path,
	groupID,
	deduplicationID string,
	delay Duration,
	async bool,
	message *PublisherMessage,
) Publisher {
	return Publisher{
		onlyIf:          onlyIf,
		ignoreIf:        ignoreIf,
		broker:          broker,
		path:            path,
		groupID:         groupID,
		deduplicationID: deduplicationID,
		delay:           delay,
		async:           async,
		message:         message,
	}
}

func NewPublisherMessage(
	continueOnError bool,
	onlyIf,
	ignoreIf []string,
	attributes map[string]PublisherMessageAttribute,
	body *PublisherMessageBody,
) *PublisherMessage {
	return &PublisherMessage{
		continueOnError: continueOnError,
		onlyIf:          onlyIf,
		ignoreIf:        ignoreIf,
		attributes:      attributes,
		body:            body,
	}
}

func NewPublisherMessageAttribute(dataType, value string) PublisherMessageAttribute {
	return PublisherMessageAttribute{
		dataType: dataType,
		value:    value,
	}
}

func NewPublisherMessageBody(omitEmpty bool, mapper *Mapper, projector *Projector, modifiers []Modifier,
) *PublisherMessageBody {
	return &PublisherMessageBody{
		omitEmpty: omitEmpty,
		mapper:    mapper,
		projector: projector,
		modifiers: modifiers,
	}
}

func (p Publisher) Broker() enum.BackendBroker {
	return p.broker
}

func (p Publisher) Path() string {
	return p.path
}

func (p Publisher) OnlyIf() []string {
	return p.onlyIf
}

func (p Publisher) IgnoreIf() []string {
	return p.ignoreIf
}

func (p Publisher) GroupID() string {
	return p.groupID
}

func (p Publisher) DeduplicationID() string {
	return p.deduplicationID
}

func (p Publisher) Delay() Duration {
	return p.delay
}

func (p Publisher) HasGroupID() bool {
	return checker.IsNotEmpty(p.groupID)
}

func (p Publisher) HasOnlyIf() bool {
	return checker.IsNotEmpty(p.onlyIf)
}

func (p Publisher) HasIgnoreIf() bool {
	return checker.IsNotEmpty(p.ignoreIf)
}

func (p Publisher) CountAllDataTransforms() (count int) {
	if !p.HasMessage() {
		return 0
	}

	if p.Message().HasBody() {
		count += p.Message().Body().CountBodyDataTransforms()
	}
	if p.Message().HasAttributes() {
		count += len(p.Message().Attributes())
	}

	return count
}

func (p Publisher) HasMessage() bool {
	return checker.NonNil(p.message)
}

func (p Publisher) Message() *PublisherMessage {
	return p.message
}

func (p Publisher) Async() bool {
	return p.async
}

func (m PublisherMessage) ContinueOnError() bool {
	return m.continueOnError
}

func (m PublisherMessage) HasBody() bool {
	return checker.NonNil(m.body)
}

func (m PublisherMessage) Body() *PublisherMessageBody {
	return m.body
}

func (m PublisherMessage) HasAttributes() bool {
	return checker.IsNotNilOrEmpty(m.attributes)
}

func (m PublisherMessage) Attributes() map[string]PublisherMessageAttribute {
	return m.attributes
}

func (a PublisherMessageAttribute) DataType() string {
	return a.dataType
}

func (a PublisherMessageAttribute) Value() string {
	return a.value
}

func (p PublisherMessageBody) HasMapper() bool {
	return checker.NonNil(p.mapper)
}

func (p PublisherMessageBody) HasProjector() bool {
	return checker.NonNil(p.projector)
}

func (p PublisherMessageBody) HasModifiers() bool {
	return checker.IsNotEmpty(p.modifiers)
}

func (p PublisherMessageBody) OmitEmpty() bool {
	return p.omitEmpty
}

func (p PublisherMessageBody) Mapper() *Mapper {
	return p.mapper
}

func (p PublisherMessageBody) Projector() *Projector {
	return p.projector
}

func (p PublisherMessageBody) Modifiers() []Modifier {
	return p.modifiers
}

func (p PublisherMessageBody) CountBodyDataTransforms() (count int) {
	if p.HasMapper() {
		count += len(p.Mapper().Map().Keys())
	}
	if p.HasProjector() {
		count += len(p.Projector().Project().Keys())
	}
	if p.HasModifiers() {
		count += len(p.Modifiers())
	}
	return count
}
