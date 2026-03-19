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

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tech4works/gopen-gateway/internal/app"
	vo2 "github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"go.elastic.co/apm/module/apmgin/v2"
)

type router struct {
	engine *gin.Engine
}

func NewRouter() app.Router {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(apmgin.Middleware(engine))

	return router{
		engine: engine,
	}
}

func (r router) Engine() http.Handler {
	return r.engine
}

func (r router) Handle(gopen *vo2.GopenConfig, endpoint *vo2.EndpointConfig, handles ...app.HandlerFunc) {
	r.engine.Handle(endpoint.Method(), endpoint.Path(), r.buildEngineHandles(gopen, endpoint, handles)...)
}

func (r router) buildEngineHandles(gopen *vo2.GopenConfig, endpoint *vo2.EndpointConfig, handles []app.HandlerFunc) []gin.HandlerFunc {
	var ginHandler []gin.HandlerFunc
	for _, handler := range handles {
		ginHandler = append(ginHandler, r.buildEngineHandle(gopen, endpoint, handler))
	}
	return ginHandler
}

func (r router) buildEngineHandle(gopen *vo2.GopenConfig, endpoint *vo2.EndpointConfig, handle app.HandlerFunc) gin.HandlerFunc {
	return func(gin *gin.Context) {
		ctx, ok := gin.Get("context")
		if !ok {
			ctx = newHTTPContext(gin, gopen, endpoint)
			gin.Set("context", ctx)
		}
		handle(ctx.(*Context))
	}
}
