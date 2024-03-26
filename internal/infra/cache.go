package infra

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/external"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/util"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	cache "github.com/chenyahui/gin-cache"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type cacheProvider struct {
	cacheVO vo.Cache
}

func NewCacheProvider(cacheVO vo.Cache) external.CacheProvider {
	return cacheProvider{
		cacheVO: cacheVO,
	}
}

func (c cacheProvider) CacheStrategyHandler() gin.HandlerFunc {
	return cache.Cache(
		c.cacheVO.MemoryStore,
		c.cacheVO.Duration,
		cache.WithCacheStrategyByRequest(cacheStrategyByRequest(c.cacheVO)),
		cache.WithBeforeReplyWithCache(cacheBeforeReplyStrategy),
	)
}

func cacheStrategyByRequest(cacheVO vo.Cache) cache.GetCacheStrategyByRequest {
	return func(ctx *gin.Context) (bool, cache.Strategy) {
		// obtemos o Cache-Control para verificar se ele quer o resultado fresco, se o allow estiver como true
		var cacheControl enum.CacheControl
		// caso esteja permitido o cache control obtemos do header
		if cacheVO.AllowCacheControl {
			cacheControl = enum.CacheControl(ctx.GetHeader("Cache-Control"))
		}

		// verificamos se no header Cache-Control veio o valor "no-cache"
		isNoCache := helper.Equals(enum.CacheControlNoCache, cacheControl)
		// obtemos o método da requisição
		method := ctx.Request.Method

		return isNoCache && helper.Equals(method, http.MethodGet), cache.Strategy{
			CacheKey: buildKeyByStrategy(ctx, cacheVO),
		}
	}
}

func cacheBeforeReplyStrategy(ctx *gin.Context, cache *cache.ResponseCache) {
	// setamos o X-Gateway-Cache indicando que a resposta veio de um cache
	ctx.Header(consts.XGatewayCache, helper.SimpleConvertToString(true))
}

func buildKeyByStrategy(ctx *gin.Context, cacheVO vo.Cache) string {
	// obtemos os dados que precisamos da requisição para setar como chave de cache
	requestUri := util.GetRequestUri(ctx)
	method := ctx.Request.Method

	// construímos a chave inicialmente com os valores de requisição
	key := fmt.Sprintf("%s:%s", requestUri, method)
	// caso key-strategy seja informado, construímos a chave a partir deles
	var valueStrategyKeys []string
	// iteramos as chaves para obter os valores
	for _, strategyKey := range cacheVO.StrategyKeys {
		valueByStrategyKey := ctx.GetHeader(strategyKey)
		if helper.IsNotEmpty(valueByStrategyKey) {
			valueStrategyKeys = append(valueStrategyKeys, valueByStrategyKey)
		}
	}

	// caso tenha encontrado valores, separamos os mesmos
	valueStrategyKey := strings.Join(valueStrategyKeys, ":")

	// caso o valor não esteja vazio retornamos o key padrão com a estratégia imposto no objeto de valor
	if helper.IsNotEmpty(valueStrategyKey) {
		return fmt.Sprintf("%s:%s", key, valueStrategyKey)
	}

	// se nao retornamos a key padrão
	return key
}
