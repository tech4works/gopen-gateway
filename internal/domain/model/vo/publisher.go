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
	provider        enum.PublisherProvider
	reference       string
	groupID         string
	deduplicationID string
	delay           Duration
}

func NewPublisher(provider enum.PublisherProvider, reference, groupID, deduplicationID string, delay Duration) Publisher {
	return Publisher{
		provider:        provider,
		reference:       reference,
		groupID:         groupID,
		deduplicationID: deduplicationID,
		delay:           delay,
	}
}

func (p Publisher) Provider() enum.PublisherProvider {
	return p.provider
}

func (p Publisher) Reference() string {
	return p.reference
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
