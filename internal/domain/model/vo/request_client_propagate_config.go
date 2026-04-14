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

// RequestClientPropagateConfig configures the header names for propagation.
type RequestClientPropagateConfig struct {
	request  string
	response string
}

func NewRequestClientPropagateConfig(request, response string) *RequestClientPropagateConfig {
	return &RequestClientPropagateConfig{
		request:  request,
		response: response,
	}
}

func (r *RequestClientPropagateConfig) Request() string {
	return r.request
}

func (r *RequestClientPropagateConfig) Response() string {
	return r.response
}
