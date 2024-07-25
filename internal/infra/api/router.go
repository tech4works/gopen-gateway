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

package api

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *Context)

func Handle(engine *gin.Engine, gopen *vo.Gopen, endpoint *vo.Endpoint, handles ...HandlerFunc) {
	engine.Handle(endpoint.Method(), endpoint.Path(), parseHandles(gopen, endpoint, handles)...)
}

func parseHandles(gopen *vo.Gopen, endpoint *vo.Endpoint, handles []HandlerFunc) []gin.HandlerFunc {
	var ginHandler []gin.HandlerFunc
	for _, apiHandler := range handles {
		ginHandler = append(ginHandler, handle(gopen, endpoint, apiHandler))
	}
	return ginHandler
}

func handle(gopen *vo.Gopen, endpoint *vo.Endpoint, handle HandlerFunc) gin.HandlerFunc {
	return func(gin *gin.Context) {
		ctx, ok := gin.Get("context")
		if !ok {
			ctx = newContext(gin, gopen, endpoint)
			gin.Set("context", ctx)
		}
		handle(ctx.(*Context))
	}
}
