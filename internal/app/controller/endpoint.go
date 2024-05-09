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

package controller

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

// endpointController represents a controller that handles endpoint requests.
// It implements the Endpoint interface and internally uses an instance of the service.Endpoint interface.
// The Execute method of the endpointController delegates the execution to the endpointService and writes the HTTP
// request and response to the context.
type endpointController struct {
	endpointService service.Endpoint
}

// Endpoint represents an interface for executing endpoint requests.
// The Execute method takes a Context parameter and is responsible for handling the request.
type Endpoint interface {
	// Execute executes the endpoint request by taking a Context parameter and handling the request.
	// It updates the current HTTP request in the context and writes the HTTP response to the client.
	// This method is part of an interface called Endpoint, which represents an interface for executing endpoint requests.
	Execute(ctx *api.Context)
}

// NewEndpoint returns a new instance of the `Endpoint` interface,
// implemented by the `endpointController` struct, which internally uses
// an instance of the `service.Endpoint` interface as a dependency.
func NewEndpoint(endpointService service.Endpoint) Endpoint {
	return endpointController{
		endpointService: endpointService,
	}
}

// Execute executes the endpoint by delegating the execution to the endpointService.
// It takes the context as a parameter and calls the BuildExecuteServiceParams function
// to build the necessary parameters for the execution. The resulting HTTP request
// and response are then written to the context using the Write method.
func (e endpointController) Execute(ctx *api.Context) {
	httpRequest, httpResponse := e.endpointService.Execute(mapper.BuildExecuteServiceParams(ctx))
	ctx.Write(httpRequest, httpResponse)
}
