package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	ginCache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

type cache struct {
}

type Cache interface {
	Do(cacheStore persist.CacheStore, cacheVO vo.Cache) gin.HandlerFunc
}

func NewCache() Cache {
	return cache{}
}

func (c cache) Do(cacheStore persist.CacheStore, cacheVO vo.Cache) gin.HandlerFunc {
	return ginCache.Cache(
		cacheStore,
		cacheVO.Duration(),
		ginCache.WithCacheStrategyByRequest(cacheStrategyByRequest(cacheVO)),
		ginCache.WithBeforeReplyWithCache(cacheBeforeReplyStrategy),
	)
}

func cacheStrategyByRequest(cacheVO vo.Cache) ginCache.GetCacheStrategyByRequest {
	return func(ctx *gin.Context) (bool, ginCache.Strategy) {
		return cacheVO.ShouldCache(ctx), ginCache.Strategy{CacheKey: cacheVO.StrategyKey(ctx)}
	}
}

func cacheBeforeReplyStrategy(ctx *gin.Context, cache *ginCache.ResponseCache) {
	// setamos o X-Gateway-Cache indicando que a resposta veio de um cache
	ctx.Header(consts.XGOpenCache, helper.SimpleConvertToString(true))
}
