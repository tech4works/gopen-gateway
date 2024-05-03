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

package mapper

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

// BuildExecuteServiceParams builds the parameters needed for executing a service.
// It takes a *api.Context as input and returns a context.Context and a vo.ExecuteEndpoint.
// The context.Context is obtained from the input request.
// The vo.ExecuteEndpoint is created using the Gopen, Endpoint, and HttpRequest from the input request.
func BuildExecuteServiceParams(ctx *api.Context) (context.Context, *vo.ExecuteEndpoint) {
	return ctx.Context(), vo.NewExecuteEndpoint(ctx.Gopen(), ctx.Endpoint(), ctx.HttpRequest())
}
