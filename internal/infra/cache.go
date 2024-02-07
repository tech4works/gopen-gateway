package infra

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CacheHandler(memoryStore *persist.RedisStore, cacheDuration time.Duration) gin.HandlerFunc {
	return cache.Cache(
		memoryStore,
		cacheDuration,
		cache.WithCacheStrategyByRequest(cacheStrategyHandler),
		cache.WithBeforeReplyWithCache(cacheBeforeReply),
	)
}

func cacheStrategyHandler(ctx *gin.Context) (bool, cache.Strategy) {
	path := ctx.Request.RequestURI
	method := ctx.Request.Method
	device := ctx.GetHeader("Device")
	ip := ctx.ClientIP()
	key := ip
	if helper.IsNotEmpty(device) {
		key = fmt.Sprint(key, ":", device)
	}
	key = fmt.Sprint(key, ":", path, ":", method)
	withoutCacheHeader := ctx.GetHeader("Without-Cache")
	shouldCache := true
	if helper.IsNotEmpty(withoutCacheHeader) {
		withoutCache := helper.SimpleConvertToBool(withoutCacheHeader)
		shouldCache = withoutCache && helper.Equals(method, http.MethodGet)
	}
	return shouldCache, cache.Strategy{CacheKey: key}
}

func cacheBeforeReply(ctx *gin.Context, cache *cache.ResponseCache) {
	logger.Info("Cache response request:", ctx.FullPath(), "statusCode:", cache.Status)
	ctx.Header("X-Gateway-Cache", helper.SimpleConvertToString(true))
}
