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

type BackendConfig struct {
	flow         enum.BackendFlow
	onlyIf       []string
	ignoreIf     []string
	id           string
	execution    BackendExecutionConfig
	dependencies *BackendDependenciesConfig
	kind         enum.BackendKind
	cache        *CacheConfig
	http         *BackendHTTPConfig
	publisher    *BackendPublisherConfig
	response     *BackendResponseConfig
}

func NewBackendConfig(
	flow enum.BackendFlow,
	onlyIf []string,
	ignoreIf []string,
	id string,
	execution BackendExecutionConfig,
	dependencies *BackendDependenciesConfig,
	kind enum.BackendKind,
	cache *CacheConfig,
	http *BackendHTTPConfig,
	publisher *BackendPublisherConfig,
	response *BackendResponseConfig,
) BackendConfig {
	return BackendConfig{
		flow:         flow,
		onlyIf:       onlyIf,
		ignoreIf:     ignoreIf,
		id:           id,
		execution:    execution,
		dependencies: dependencies,
		kind:         kind,
		cache:        cache,
		http:         http,
		publisher:    publisher,
		response:     response,
	}
}

func (b *BackendConfig) OnlyIf() []string {
	return b.onlyIf
}

func (b *BackendConfig) IgnoreIf() []string {
	return b.ignoreIf
}

func (b *BackendConfig) ID() string {
	return b.id
}

func (b *BackendConfig) Execution() BackendExecutionConfig {
	return b.execution
}

func (b *BackendConfig) HasDependencies() bool {
	return checker.NonNil(b.dependencies)
}

func (b *BackendConfig) Dependencies() *BackendDependenciesConfig {
	return b.dependencies
}

func (b *BackendConfig) HasCache() bool {
	return checker.NonNil(b.cache)
}

func (b *BackendConfig) AllowCache() bool {
	return b.IsHTTP()
}

func (b *BackendConfig) Cache() *CacheConfig {
	return b.cache
}

func (b *BackendConfig) HTTP() *BackendHTTPConfig {
	return b.http
}

func (b *BackendConfig) Publisher() *BackendPublisherConfig {
	return b.publisher
}

func (b *BackendConfig) Kind() enum.BackendKind {
	return b.kind
}

func (b *BackendConfig) IsHTTP() bool {
	return checker.Equals(b.kind, enum.BackendKindHTTP)
}

func (b *BackendConfig) IsPublisher() bool {
	return checker.Equals(b.kind, enum.BackendKindPublisher)
}

func (b *BackendConfig) IsBeforeware() bool {
	return checker.Equals(b.flow, enum.BackendFlowBeforeware)
}

func (b *BackendConfig) IsNormal() bool {
	return checker.Equals(b.flow, enum.BackendFlowNormal)
}

func (b *BackendConfig) IsAfterware() bool {
	return checker.Equals(b.flow, enum.BackendFlowAfterware)
}

func (b *BackendConfig) IsMiddleware() bool {
	return b.IsBeforeware() || b.IsAfterware()
}

func (b *BackendConfig) Flow() enum.BackendFlow {
	return b.flow
}

func (b *BackendConfig) HasResponse() bool {
	return checker.NonNil(b.response)
}

func (b *BackendConfig) Response() *BackendResponseConfig {
	return b.response
}

func (b *BackendConfig) CountResponseDataTransforms() (count int) {
	if b.HasResponse() {
		count += b.Response().CountAllDataTransforms()
	}
	return count
}

func (b *BackendConfig) CountAllDataTransforms() (count int) {
	switch b.Kind() {
	case enum.BackendKindHTTP:
		count += b.HTTP().CountAllDataTransforms()
	case enum.BackendKindPublisher:
		count += b.Publisher().CountAllDataTransforms()
	}
	count += b.CountResponseDataTransforms()
	return count
}
