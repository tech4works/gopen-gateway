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
	"fmt"
	"time"

	"github.com/tech4works/checker"
)

type EndpointConfig struct {
	execution     EndpointExecutionConfig
	path          string
	method        string
	timeout       Duration
	securityCors  *SecurityCorsConfig
	limiter       *LimiterConfig
	cache         *CacheConfig
	backends      []BackendConfig
	response      EndpointResponseConfig
	requestClient *RequestClientConfig
}

func NewEndpointConfig(
	execution EndpointExecutionConfig,
	path,
	method string,
	timeout Duration,
	securityCors *SecurityCorsConfig,
	limiter *LimiterConfig,
	cache *CacheConfig,
	backends []BackendConfig,
	response EndpointResponseConfig,
	requestClient *RequestClientConfig,
) EndpointConfig {
	return EndpointConfig{
		execution:     execution,
		path:          path,
		method:        method,
		timeout:       timeout,
		securityCors:  securityCors,
		limiter:       limiter,
		cache:         cache,
		backends:      backends,
		response:      response,
		requestClient: requestClient,
	}
}

func NewEndpointConfigStatic(path, method string) EndpointConfig {
	return EndpointConfig{
		path:    path,
		method:  method,
		timeout: NewDuration(10 * time.Second),
		limiter: NewLimiterConfigDefault(),
	}
}

func (e *EndpointConfig) Execution() EndpointExecutionConfig {
	return e.execution
}

func (e *EndpointConfig) Path() string {
	return e.path
}

func (e *EndpointConfig) Method() string {
	return e.method
}

func (e *EndpointConfig) Timeout() Duration {
	if checker.IsGreaterThan(e.timeout, 0) {
		return e.timeout
	}
	return NewDuration(30 * time.Second)
}

func (e *EndpointConfig) HasSecurityCode() bool {
	return checker.NonNil(e.SecurityCors())
}

func (e *EndpointConfig) SecurityCors() *SecurityCorsConfig {
	return e.securityCors
}

func (e *EndpointConfig) HasLimiter() bool {
	return checker.NonNil(e.Limiter())
}

func (e *EndpointConfig) Limiter() *LimiterConfig {
	return e.limiter
}

func (e *EndpointConfig) HasCache() bool {
	return checker.NonNil(e.Cache())
}

func (e *EndpointConfig) AllowCache() bool {
	return true
}

func (e *EndpointConfig) Cache() *CacheConfig {
	return e.cache
}

func (e *EndpointConfig) Backends() []BackendConfig {
	return e.backends
}

func (e *EndpointConfig) CountBeforewares() (count int) {
	for _, backend := range e.backends {
		if backend.IsBeforeware() {
			count++
		}
	}
	return count
}

func (e *EndpointConfig) CountAfterwares() (count int) {
	for _, backend := range e.backends {
		if backend.IsAfterware() {
			count++
		}
	}
	return count
}

func (e *EndpointConfig) CountBackends() (count int) {
	for _, backend := range e.backends {
		if backend.IsNormal() {
			count++
		}
	}
	return count
}

func (e *EndpointConfig) CountAllDataTransforms() (count int) {
	if e.Response().HasPayload() {
		count += e.Response().Payload().CountDataTransforms()
	}
	for _, backend := range e.backends {
		count += backend.CountAllDataTransforms()
	}
	return count
}

func (e *EndpointConfig) Response() EndpointResponseConfig {
	return e.response
}

func (e *EndpointConfig) RequestClient() *RequestClientConfig {
	return e.requestClient
}

func (e *EndpointConfig) Resume() string {
	return fmt.Sprintf("%s --> \"%s\" (beforeware:%v, backends:%v, afterware:%v, transformations:%v)",
		e.method, e.path, e.CountBeforewares(), e.CountAfterwares(), e.CountBackends(), e.CountAllDataTransforms())
}
