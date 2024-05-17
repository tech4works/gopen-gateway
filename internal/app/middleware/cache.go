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

package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

// cacheMiddleware represents a middleware that handles caching for an API endpoint.
type cacheMiddleware struct {
	cacheService service.Cache
	logger       interfaces.LoggerProvider
}

// Cache represents an interface for caching operations. Implementations of this interface
// should provide a Do method that takes a context as input and is responsible for performing
// the caching operation based on the provided context.
type Cache interface {
	// Do perform a caching operation based on the provided context.
	//
	// The context parameter is a pointer to an api.Context object that holds the information
	// needed to perform the caching operation. This method should be implemented by
	// implementations of the Cache interface.
	Do(ctx *api.Context)
}

// NewCache initializes a new cacheMiddleware struct based on the provided cacheService.
// The cacheMiddleware struct implements the Cache interface, and is responsible for handling
// caching operations for an API endpoint. The NewCache function returns a Cache interface.
func NewCache(cacheService service.Cache, loggerProvider interfaces.LoggerProvider) Cache {
	return cacheMiddleware{
		cacheService: cacheService,
		logger:       loggerProvider,
	}
}

// Do handle caching for an API endpoint. It first checks if the endpoint has caching enabled.
// If caching is disabled, it calls the next handler in the chain.
// If caching is enabled, it retrieves the cache configuration and strategy key.
// Then, it calls the cache service to read the cache based on the provided cache configuration and request.
// If there is an error during cache reading, it logs a warning.
// If there is a cache response available, it writes the cache response to the context and returns.
// After that, it calls the next handler in the chain.
// Finally, it calls the cache service to write the cache based on the provided cache configuration, request, and response.
// If there is an error during cache writing, it logs a warning.
func (c cacheMiddleware) Do(ctx *api.Context) {
	if ctx.Endpoint().NoCache() {
		ctx.Next()
		return
	}

	cache := ctx.Endpoint().Cache()
	httpRequest := ctx.HttpRequest()
	strategyKey := cache.StrategyKey(httpRequest)

	cacheResponse, err := c.cacheService.Read(ctx.Context(), cache, httpRequest)
	if helper.IsNotNil(err) {
		c.logger.PrintEndpointWarnf(ctx, "Error read cache key: %s err: %s", strategyKey, err)
	} else if helper.IsNotNil(cacheResponse) {
		ctx.WriteCacheResponse(cacheResponse)
		return
	}

	ctx.Next()

	err = c.cacheService.Write(ctx.Context(), cache, httpRequest, ctx.HttpResponse())
	if helper.IsNotNil(err) {
		c.logger.PrintEndpointWarnf(ctx, "Error write cache key: %s err: %s", strategyKey, err)
	}
}
