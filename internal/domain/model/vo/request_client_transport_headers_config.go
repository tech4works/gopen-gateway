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

import "github.com/tech4works/checker"

// RequestClientTransportHeadersRequestConfig controla headers injetados no request ao backend.
// Todos os métodos retornam true quando o ponteiro correspondente é nil (default habilitado).
type RequestClientTransportHeadersRequestConfig struct {
	degradation *bool
	timeout     *bool
}

func NewRequestClientTransportHeadersRequestConfig(
	degradation *bool,
	timeout *bool,
) RequestClientTransportHeadersRequestConfig {
	return RequestClientTransportHeadersRequestConfig{
		degradation: degradation,
		timeout:     timeout,
	}
}

func (r RequestClientTransportHeadersRequestConfig) DegradationEnabled() bool {
	if checker.IsNil(r.degradation) {
		return true
	}
	return *r.degradation
}

func (r RequestClientTransportHeadersRequestConfig) TimeoutEnabled() bool {
	if checker.IsNil(r.timeout) {
		return true
	}
	return *r.timeout
}

// RequestClientTransportHeadersResponseConfig controla headers injetados no response ao client.
// Todos os métodos retornam true quando o ponteiro correspondente é nil (default habilitado).
type RequestClientTransportHeadersResponseConfig struct {
	cache           *bool
	backendCache    *bool
	executionStatus *bool
	degradation     *bool
}

func NewRequestClientTransportHeadersResponseConfig(
	cache *bool,
	backendCache *bool,
	executionStatus *bool,
	degradation *bool,
) RequestClientTransportHeadersResponseConfig {
	return RequestClientTransportHeadersResponseConfig{
		cache:           cache,
		backendCache:    backendCache,
		executionStatus: executionStatus,
		degradation:     degradation,
	}
}

func (r RequestClientTransportHeadersResponseConfig) CacheEnabled() bool {
	if checker.IsNil(r.cache) {
		return true
	}
	return *r.cache
}

func (r RequestClientTransportHeadersResponseConfig) BackendCacheEnabled() bool {
	if checker.IsNil(r.backendCache) {
		return false // default disabled
	}
	return *r.backendCache
}

func (r RequestClientTransportHeadersResponseConfig) ExecutionStatusEnabled() bool {
	if checker.IsNil(r.executionStatus) {
		return true
	}
	return *r.executionStatus
}

func (r RequestClientTransportHeadersResponseConfig) DegradationEnabled() bool {
	if checker.IsNil(r.degradation) {
		return true
	}
	return *r.degradation
}

// RequestClientTransportHeadersConfig controla quais grupos de transport headers são injetados.
type RequestClientTransportHeadersConfig struct {
	request  RequestClientTransportHeadersRequestConfig
	response RequestClientTransportHeadersResponseConfig
}

func NewRequestClientTransportHeadersConfig(
	request RequestClientTransportHeadersRequestConfig,
	response RequestClientTransportHeadersResponseConfig,
) *RequestClientTransportHeadersConfig {
	return &RequestClientTransportHeadersConfig{
		request:  request,
		response: response,
	}
}

func (r *RequestClientTransportHeadersConfig) Request() RequestClientTransportHeadersRequestConfig {
	return r.request
}

func (r *RequestClientTransportHeadersConfig) Response() RequestClientTransportHeadersResponseConfig {
	return r.response
}
