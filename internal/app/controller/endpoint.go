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

type endpointController struct {
	endpointService service.Endpoint
}

type Endpoint interface {
	Execute(ctx *api.Context)
}

func NewEndpoint(endpointService service.Endpoint) Endpoint {
	return endpointController{
		endpointService: endpointService,
	}
}

func (e endpointController) Execute(ctx *api.Context) {
	httpRequest, httpResponse := e.endpointService.Execute(mapper.BuildExecuteServiceParams(ctx))
	ctx.Write(httpRequest, httpResponse)
}
