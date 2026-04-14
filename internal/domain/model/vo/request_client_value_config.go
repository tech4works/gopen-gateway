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
	"net/http"

	"github.com/google/uuid"

	"github.com/tech4works/checker"
)

// RequestClientValueConfig configures the resolution and propagation of a request-client value
// (request-id or trace) from external headers, with configurable fallback and propagation.
type RequestClientValueConfig struct {
	headers   []string
	fallback  *bool
	propagate *RequestClientPropagateConfig
}

func NewRequestClientValueConfig(
	headers []string,
	fallback *bool,
	propagate *RequestClientPropagateConfig,
) *RequestClientValueConfig {
	return &RequestClientValueConfig{
		headers:   headers,
		fallback:  fallback,
		propagate: propagate,
	}
}

func (r *RequestClientValueConfig) Headers() []string {
	return r.headers
}

// Fallback returns true when the fallback pointer is nil (default enabled),
// otherwise returns the pointed value.
func (r *RequestClientValueConfig) Fallback() bool {
	if checker.IsNil(r.fallback) {
		return true
	}
	return *r.fallback
}

func (r *RequestClientValueConfig) Propagate() *RequestClientPropagateConfig {
	return r.propagate
}

// HasHeaders returns true when headers is non-nil and non-empty.
func (r *RequestClientValueConfig) HasHeaders() bool {
	return checker.NonNil(r.headers) && checker.IsNotEmpty(r.headers)
}

// HasPropagateRequest returns true when propagate is non-nil and propagate.Request() is non-empty.
func (r *RequestClientValueConfig) HasPropagateRequest() bool {
	return checker.NonNil(r.propagate) && checker.IsNotEmpty(r.propagate.Request())
}

// HasPropagateResponse returns true when propagate is non-nil and propagate.Response() is non-empty.
func (r *RequestClientValueConfig) HasPropagateResponse() bool {
	return checker.NonNil(r.propagate) && checker.IsNotEmpty(r.propagate.Response())
}

// ResolveRequestID resolves the request-id from the request headers.
// It iterates over r.headers and returns the first non-empty value found.
// If no header is found and Fallback() is true, returns a new UUID string.
// If no header is found and Fallback() is false, returns "".
func (r *RequestClientValueConfig) ResolveRequestID(header http.Header) string {
	for _, h := range r.headers {
		if val := header.Get(h); checker.IsNotEmpty(val) {
			return val
		}
	}
	if r.Fallback() {
		return uuid.New().String()
	}
	return ""
}
