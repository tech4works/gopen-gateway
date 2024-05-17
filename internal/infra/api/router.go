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

// HandlerFunc is a function type that represents an HTTP request handler.
// It takes a pointer to a Context object as its parameter and does not return any value.
type HandlerFunc func(ctx *Context)

// Handle registers a handler with the specified engine, gopen and endpoint.
// It takes a pointer to gin.Engine, vo.Gopen, vo.Endpoint and one or more HandlerFuncs as arguments.
// The interface references are used to configure the handler and register it with the engine.
// The endpoint's method and path, along with the gopen configuration and array of HandlerFuncs
// are used to set up the handler. The parseHandles function is called to process the gopen, endpoint,
// and HandlerFuncs and return the final handler to be registered with the engine.
func Handle(engine *gin.Engine, gopen *vo.Gopen, endpoint *vo.Endpoint, handles ...HandlerFunc) {
	engine.Handle(endpoint.Method(), endpoint.Path(), parseHandles(gopen, endpoint, handles)...)
}

// parseHandles takes a pointer to vo.Gopen, vo.Endpoint, and a slice of HandlerFuncs as arguments.
// It iterates over the slice of HandlerFuncs and calls the handle function for each one,
// passing in the gopen, endpoint, and HandlerFunc as arguments and appending the returned
// gin.HandlerFunc to a slice. The resulting slice of gin.HandlerFuncs is returned.
func parseHandles(gopen *vo.Gopen, endpoint *vo.Endpoint, handles []HandlerFunc) []gin.HandlerFunc {
	var ginHandler []gin.HandlerFunc
	for _, apiHandler := range handles {
		ginHandler = append(ginHandler, handle(gopen, endpoint, apiHandler))
	}
	return ginHandler
}

// handle is a function that takes a pointer to vo.Gopen, vo.Endpoint, and a HandlerFunc as arguments.
// It returns a gin.HandlerFunc. The returned handler function takes a gin.Context as its argument.
// Inside the returned handler function, the gin.Context is checked for a "context" value.
// If the "context" value is not found, a new context is built using the newContext function
// with the provided vo.Gopen and vo.Endpoint. The newly built context is then set as the "context" value in the gin.Context.
// Finally, the provided HandlerFunc is called with the context.(*Context) as its argument.
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
