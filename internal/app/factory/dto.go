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

package factory

import (
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
)

func BuildExecuteEndpoint(ctx app.Context) dto.ExecuteEndpoint {
	return dto.ExecuteEndpoint{
		TraceID:  ctx.TraceID(),
		ClientIP: ctx.ClientIP(),
		Gopen:    ctx.Gopen(),
		Endpoint: ctx.Endpoint(),
		Request:  ctx.Request(),
	}
}
