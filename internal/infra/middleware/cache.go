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
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

// cacheMiddleware represents a Cache implementation that uses the provided infra.CacheStore for caching operations.
type cacheMiddleware struct {
	cacheService service.Cache
}

type Cache interface {
	Do(ctx *api.Context)
}

func NewCache(cacheService service.Cache) Cache {
	return cacheMiddleware{
		cacheService: cacheService,
	}
}

func (c cacheMiddleware) Do(ctx *api.Context) {
	// se meu endpoint não tem cache, ignoramos e chamamos o próximo
	if ctx.Endpoint().NoCache() {
		ctx.Next()
		return
	}

	// instanciamos o cache vo
	cache := ctx.Endpoint().Cache()
	// instanciar httpRequest para manter o mesmo do inicio ao fim
	httpRequest := ctx.HttpRequest()
	// instanciamos a chave utilizada nos logs
	strategyKey := cache.StrategyKey(httpRequest)

	// chamamos o serviço de dominio esperando um possível response caso configurado e caso tenha
	cacheResponse, err := c.cacheService.Read(ctx.Context(), cache, httpRequest)
	if helper.IsNotNil(err) {
		logger.Warning("Error read cache key:", strategyKey, "err:", err)
	} else if helper.IsNotNil(cacheResponse) {
		ctx.WriteCacheResponse(cacheResponse)
		return
	}

	// damos próximo no handler
	ctx.Next()

	// chamamos o serviço de dominio para gravar o cache caso configurado
	err = c.cacheService.Write(ctx.Context(), cache, httpRequest, ctx.HttpResponse())
	if helper.IsNotNil(err) {
		logger.Warning("Error write cache key:", strategyKey, "err:", err)
	}
}
