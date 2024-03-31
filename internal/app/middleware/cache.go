package middleware

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/util"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
)

type cache struct {
}

type Cache interface {
	Do(cacheStore interfaces.CacheStore, endpointVO vo.Endpoint, cacheVO vo.Cache) gin.HandlerFunc
}

func NewCache() Cache {
	return cache{}
}

func (c cache) Do(cacheStore interfaces.CacheStore, endpointVO vo.Endpoint, cacheVO vo.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// inicializamos a chave que vai ser utilizada
		key := cacheVO.StrategyKey(ctx)

		// verificamos se ele permite ler o cache
		if cacheVO.CanRead(ctx) {
			// inicializamos o valor a ser obtido
			var cacheResponse vo.CacheResponse

			// obtemos através do cache store se a chave exists respondemos, se não seguimos normalmente
			err := cacheStore.Get(ctx.Request.Context(), key, &cacheResponse)
			if helper.IsNil(err) {
				util.RespondGateway(ctx, vo.NewResponseByCache(endpointVO, cacheResponse))
				return
			} else if errors.IsNot(err, mapper.ErrCacheNotFound) {
				logger.Warning("Error read cache key:", key, "err:", err)
			}
		}

		// damos próximo no handler
		ctx.Next()

		// verificamos se podemos gravar a resposta
		if cacheVO.CanWrite(ctx) {
			// instanciamos a duração
			duration := cacheVO.Duration()

			// obtemos o response writer
			responseWriter := util.GetResponseWriter(ctx)
			// transformamos em cacheResponse e setamos
			err := cacheStore.Set(ctx, key, vo.NewCacheResponse(responseWriter, duration), duration)
			if helper.IsNotNil(err) {
				logger.Warning("Error write cache key:", key, "err:", err)
			}
		}
	}
}
