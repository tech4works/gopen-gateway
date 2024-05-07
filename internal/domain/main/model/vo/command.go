/*
 * Copyright 2024 Gabriel Cataldo
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
	configVO "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
)

// ExecuteEndpoint represents the execution of a specific endpoint in the Gopen server.
// It contains the configuration for the Gopen server, the targeted endpoint, and the HTTP httpRequest for execution.
type ExecuteEndpoint struct {
	// Gopen represents the configuration for the Gopen server, including environment, version, hot reload status, port,
	// timeout duration, limiter, cache, security CORS, middlewares, and endpoints.
	gopen *configVO.Gopen
	// endpoint represents a specific endpoint in the Gopen server.
	endpoint *configVO.Endpoint
	// httpRequest represents an HTTP `httpRequest` object.
	httpRequest *HttpRequest
}

// ExecuteBackend is a type that represents the execution of a backend server httpRequest and httpResponse.
type ExecuteBackend struct {
	endpoint *configVO.Endpoint
	// Backend represents a backend server configuration.
	backend *configVO.Backend
	// httpRequest represents an HTTP `httpRequest` object.
	httpRequest *HttpRequest
	// httpResponse represents an HTTP `httpResponse` object.
	httpResponse *HttpResponse
}

// NewExecuteEndpoint creates a new ExecuteEndpoint using the provided Gopen, Endpoint, and HttpRequest objects.
func NewExecuteEndpoint(gopen *configVO.Gopen, endpoint *configVO.Endpoint, httpRequest *HttpRequest) *ExecuteEndpoint {
	return &ExecuteEndpoint{
		gopen:       gopen,
		endpoint:    endpoint,
		httpRequest: httpRequest,
	}
}

// NewExecuteBackend creates a new ExecuteBackend using the provided Endpoint, Backend, HttpRequest, and HttpResponse objects.
func NewExecuteBackend(endpoint *configVO.Endpoint, backend *configVO.Backend, httpRequest *HttpRequest,
	httpResponse *HttpResponse) *ExecuteBackend {
	return &ExecuteBackend{
		endpoint:     endpoint,
		backend:      backend,
		httpRequest:  httpRequest,
		httpResponse: httpResponse,
	}
}

// Gopen returns the Gopen object associated with the ExecuteEndpoint object.
func (e ExecuteEndpoint) Gopen() *configVO.Gopen {
	return e.gopen
}

// Endpoint returns the Endpoint object associated with the ExecuteEndpoint object.
func (e ExecuteEndpoint) Endpoint() *configVO.Endpoint {
	return e.endpoint
}

// HttpRequest returns the HttpRequest object associated with the ExecuteEndpoint object.
func (e ExecuteEndpoint) HttpRequest() *HttpRequest {
	return e.httpRequest
}

// Endpoint returns the Endpoint object associated with the ExecuteEndpoint object.
func (e ExecuteBackend) Endpoint() *configVO.Endpoint {
	return e.endpoint
}

// Backend returns the Backend object associated with the ExecuteBackend object.
func (e ExecuteBackend) Backend() *configVO.Backend {
	return e.backend
}

// HttpRequest returns the HttpRequest object associated with the ExecuteBackend object.
func (e ExecuteBackend) HttpRequest() *HttpRequest {
	return e.httpRequest
}

// HttpResponse returns the HttpResponse object associated with the ExecuteBackend object.
func (e ExecuteBackend) HttpResponse() *HttpResponse {
	return e.httpResponse
}
