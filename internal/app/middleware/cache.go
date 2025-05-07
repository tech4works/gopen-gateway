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

package middleware

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type cacheMiddleware struct {
	service service.Cache
	log     app.MiddlewareLog
}

type Cache interface {
	Do(ctx app.Context)
}

func NewCache(service service.Cache, log app.MiddlewareLog) Cache {
	return cacheMiddleware{
		service: service,
		log:     log,
	}
}

func (c cacheMiddleware) Do(ctx app.Context) {
	if ctx.Endpoint().NoCache() {
		ctx.Next()
		return
	}

	response, err := c.service.Read(ctx.Context(), ctx.Endpoint().Cache(), ctx.Request())
	if checker.NonNil(err) {
		c.printWarnf(ctx, "Error read cache err: %s", err)
	} else if checker.NonNil(response) {
		ctx.WriteCacheResponse(response)
		return
	}

	ctx.Next()

	err = c.service.Write(ctx.Context(), ctx.Endpoint().Cache(), ctx.Request(), ctx.Response())
	if checker.NonNil(err) {
		c.printWarnf(ctx, "Error write cache err: %s", err)
	}
}

func (c cacheMiddleware) printWarnf(ctx app.Context, format string, msg ...any) {
	c.log.PrintWarnf(ctx, format, msg...)
}
