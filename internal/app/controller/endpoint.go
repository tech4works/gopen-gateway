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

// endpoint represents an endpoint implementation that uses the endpointService to
// execute the service logic and write the response to the request.
type endpoint struct {
	endpointService service.Endpoint
}

// Endpoint represents an interface for executing the service logic and responding to a request.
// The Execute method takes a context object as a parameter and is responsible for handling the request.
type Endpoint interface {
	// Execute executes the service logic and responds to a request.
	// It takes a context object as a parameter.
	// --
	// Note: This method is part of the Endpoint interface.
	Execute(ctx *api.Context)
}

// NewEndpoint creates a new endpoint instance using the provided endpointService.
// It returns an Endpoint object.
func NewEndpoint(endpointService service.Endpoint) Endpoint {
	return endpoint{
		endpointService: endpointService,
	}
}

// Execute executes the service to process the endpoint.
// It builds the service parameters using mapper.BuildExecuteServiceParams.
// It invokes the endpointService.Execute method passing the built parameters.
// It writes the response to the request using ctx.Write method.
func (e endpoint) Execute(ctx *api.Context) {
	// executamos o serviço de dominío para processar o endpoint
	responseVO := e.endpointService.Execute(mapper.BuildExecuteServiceParams(ctx))
	// respondemos a requisição a partir do objeto de valor recebido
	ctx.Write(responseVO)
}
