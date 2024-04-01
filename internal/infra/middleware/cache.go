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

type cache struct {
	cacheStore infra.CacheStore
}

type Cache interface {
	Do(cacheVO vo.Cache) api.HandlerFunc
}

func NewCache(cacheStore infra.CacheStore) Cache {
	return cache{
		cacheStore: cacheStore,
	}
}

func (c cache) Do(cacheVO vo.Cache) api.HandlerFunc {
	return func(req *api.Request) {
		// inicializamos a chave que vai ser utilizada
		key := cacheVO.StrategyKey(req.Method(), req.Url(), req.Header())

		// verificamos se ele permite ler o cache
		if cacheVO.CanRead(req.Method(), req.Header()) {
			// inicializamos o valor a ser obtido
			var cacheResponse vo.CacheResponse

			// obtemos através do cache store se a chave exists respondemos, se não seguimos normalmente
			err := c.cacheStore.Get(req.Context(), key, &cacheResponse)
			if helper.IsNil(err) {
				req.WriteCacheResponse(cacheResponse)
				return
			} else if errors.IsNot(err, mapper.ErrCacheNotFound) {
				logger.Warning("Error read cache key:", key, "err:", err)
			}
		}

		// damos próximo no handler
		req.Next()

		// verificamos se podemos gravar a resposta
		if cacheVO.CanWrite(req.Method(), req.Header()) {
			// instanciamos a duração
			duration := cacheVO.Duration()

			// construímos o valor a ser setado no cache
			cacheResponse := vo.NewCacheResponse(req.Writer(), duration)

			// transformamos em cacheResponse e setamos
			err := c.cacheStore.Set(req.Context(), key, cacheResponse, duration)
			if helper.IsNotNil(err) {
				logger.Warning("Error write cache key:", key, "err:", err)
			}
		}
	}
}
