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
	configVO "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/vo"
	"github.com/gin-gonic/gin"
	"sync"
)

type HandlerFunc func(ctx *Context)

func Handle(engine *gin.Engine, gopen *configVO.Gopen, endpoint *configVO.Endpoint, handles ...HandlerFunc) {
	engine.Handle(endpoint.Method(), endpoint.Path(), parseHandles(gopen, endpoint, handles)...)
}

func parseHandles(gopen *configVO.Gopen, endpoint *configVO.Endpoint, handles []HandlerFunc) []gin.HandlerFunc {
	var ginHandler []gin.HandlerFunc
	for _, apiHandler := range handles {
		ginHandler = append(ginHandler, handle(gopen, endpoint, apiHandler))
	}
	return ginHandler
}

func handle(gopen *configVO.Gopen, endpoint *configVO.Endpoint, handle HandlerFunc) gin.HandlerFunc {
	return func(gin *gin.Context) {
		// verificamos se esse contexto ja foi construído
		ctx, ok := gin.Get("context")
		if !ok {
			// construímos o contexto da requisição através dos objetos de valores e o gin
			ctx = buildContext(gin, gopen, endpoint)
			// setamos o contexto criado da requisição
			gin.Set("context", ctx)
		}
		// chamamos a função persistida
		handle(ctx.(*Context))
	}
}

func buildContext(gin *gin.Context, gopen *configVO.Gopen, endpoint *configVO.Endpoint) *Context {
	// o contexto da requisição é criado
	return &Context{
		mutex:       &sync.RWMutex{},
		framework:   gin,
		gopen:       gopen,
		endpoint:    endpoint,
		httpRequest: vo.NewHttpRequest(gin),
	}
}
