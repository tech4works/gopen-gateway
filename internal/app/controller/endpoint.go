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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/usecase"
)

type endpointController struct {
	endpointUseCase usecase.Endpoint
}

type Endpoint interface {
	Execute(ctx app.Context)
}

func NewEndpoint(endpointUseCase usecase.Endpoint) Endpoint {
	return endpointController{
		endpointUseCase: endpointUseCase,
	}
}

func (e endpointController) Execute(ctx app.Context) {
	response := e.endpointUseCase.Execute(ctx.Context(), dto.ExecuteEndpoint{
		Gopen:    ctx.Gopen(),
		Endpoint: ctx.Endpoint(),
		Request:  ctx.Request(),
	})
	ctx.Write(response)
}
