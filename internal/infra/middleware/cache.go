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
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

// cache represents a Cache implementation that uses the provided infra.CacheStore for caching operations.
type cache struct {
	cacheStore infra.CacheStore
}

// Cache represents an interface for caching operations.
// The Do method takes a cacheVO object and returns a HandlerFunc function
// that can be used as an HTTP route handler.
type Cache interface {
	// Do takes a endpointCacheVO object and returns a HandlerFunc function
	// that can be used as an HTTP route handler.
	//
	// The cacheVO object contains information about cache configuration
	// such as duration, strategy headers, allowed status codes, and allowed methods.
	//
	// The returned HandlerFunc is responsible for handling the HTTP request,
	// implementing cache-related logic based on the provided cache configuration.
	Do(endpointCacheVO *vo.EndpointCache) api.HandlerFunc
}

// NewCache returns a Cache implementation that uses the provided CacheStore for caching operations.
func NewCache(cacheStore infra.CacheStore) Cache {
	return cache{
		cacheStore: cacheStore,
	}
}

// Do execute the cache logic based on the provided endpoint cache value object and returns a HandlerFunc.
// It initializes the cache key based on the strategy, checks if the cache can be read, and responds with the cached value if available.
// If the cache cannot be read or is not found, it proceeds to the next handler.
// After the next handler is executed, it checks if the response can be cached, sets the cache value, and logs any errors.
func (c cache) Do(endpointCacheVO *vo.EndpointCache) api.HandlerFunc {
	return func(ctx *api.Context) {
		// se for nil vamos para o próximo
		if helper.IsNil(endpointCacheVO) {
			ctx.Next()
			return
		}

		// inicializamos a chave que vai ser utilizada
		key := endpointCacheVO.StrategyKey(ctx.Request())

		// verificamos se ele permite ler o cache
		if endpointCacheVO.CanRead(ctx.Request()) {
			// inicializamos o valor a ser obtido
			var cacheResponse vo.CacheResponse

			// obtemos através do cache store se a chave exists respondemos, se não seguimos normalmente
			err := c.cacheStore.Get(ctx.Context(), key, &cacheResponse)
			if helper.IsNil(err) {
				ctx.WriteCacheResponse(&cacheResponse)
				return
			} else if errors.IsNot(err, mapper.ErrCacheNotFound) {
				logger.Warning("Error read cache key:", key, "err:", err)
			}
		}

		// damos próximo no handler
		ctx.Next()

		// verificamos se podemos gravar a resposta
		if endpointCacheVO.CanWrite(ctx.Request(), ctx.Response()) {
			// instanciamos a duração
			duration := endpointCacheVO.Duration()

			// construímos o valor a ser setado no cache
			cacheResponse := vo.NewCacheResponse(ctx.Response(), duration)

			// transformamos em cacheResponse e setamos
			err := c.cacheStore.Set(ctx.Context(), key, cacheResponse, duration.Time())
			if helper.IsNotNil(err) {
				logger.Warning("Error write cache key:", key, "err:", err)
			}
		}
	}
}
