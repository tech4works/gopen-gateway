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

// RequestClientConfig encapsulates all request-client configuration.
type RequestClientConfig struct {
	requestID        *RequestClientValueConfig
	ip               *RequestClientIPConfig
	transportHeaders *RequestClientTransportHeadersConfig
}

func NewRequestClientConfig(
	requestID *RequestClientValueConfig,
	ip *RequestClientIPConfig,
	transportHeaders *RequestClientTransportHeadersConfig,
) *RequestClientConfig {
	return &RequestClientConfig{
		requestID:        requestID,
		ip:               ip,
		transportHeaders: transportHeaders,
	}
}

func (r *RequestClientConfig) RequestID() *RequestClientValueConfig {
	if checker.IsNil(r) {
		return nil
	}
	return r.requestID
}

// IP returns the IP config, or a zero-value &RequestClientIPConfig{} when r or r.ip is nil.
func (r *RequestClientConfig) IP() *RequestClientIPConfig {
	if checker.IsNil(r) || checker.IsNil(r.ip) {
		return &RequestClientIPConfig{}
	}
	return r.ip
}

func (r *RequestClientConfig) TransportHeaders() *RequestClientTransportHeadersConfig {
	if checker.IsNil(r) {
		return nil
	}
	return r.transportHeaders
}

// TransportHeadersRequest returns the request transport headers config.
// Returns a config with all nil pointers (all enabled) when RequestClientConfig or TransportHeaders are nil.
func (r *RequestClientConfig) TransportHeadersRequest() RequestClientTransportHeadersRequestConfig {
	if checker.IsNil(r) || checker.IsNil(r.transportHeaders) {
		return NewRequestClientTransportHeadersRequestConfig(nil, nil)
	}
	return r.transportHeaders.Request()
}

// TransportHeadersResponse returns the response transport headers config.
// Returns a config with all nil pointers (all enabled) when RequestClientConfig or TransportHeaders are nil.
func (r *RequestClientConfig) TransportHeadersResponse() RequestClientTransportHeadersResponseConfig {
	if checker.IsNil(r) || checker.IsNil(r.transportHeaders) {
		return NewRequestClientTransportHeadersResponseConfig(nil, nil, nil, nil)
	}
	return r.transportHeaders.Response()
}
