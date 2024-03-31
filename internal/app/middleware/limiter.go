package middleware

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/util"
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type limiter struct {
}

type Limiter interface {
	Do(endpointVO vo.Endpoint, rateLimiterProvider interfaces.RateLimiterProvider,
		sizeLimiterProvider interfaces.SizeLimiterProvider) gin.HandlerFunc
}

func NewLimiter() Limiter {
	return limiter{}
}

func (l limiter) Do(
	endpointVO vo.Endpoint,
	rateLimiterProvider interfaces.RateLimiterProvider,
	sizeLimiterProvider interfaces.SizeLimiterProvider,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// aqui ja verificamos se a chave hoje sendo ela o IP está permitida
		err := rateLimiterProvider.Allow(ctx.GetHeader(consts.XForwardedFor))
		if helper.IsNotNil(err) {
			util.RespondGatewayError(ctx, endpointVO.ResponseEncode(), http.StatusTooManyRequests, err)
			return
		}

		// verificamos o tamanho da requisição, e tratamos o erro logo em seguida
		err = sizeLimiterProvider.Allow(ctx.Request)
		if errors.Contains(err, domainmapper.ErrHeaderTooLarge) {
			util.RespondGatewayError(ctx, endpointVO.ResponseEncode(), http.StatusRequestHeaderFieldsTooLarge, err)
			return
		} else if helper.IsNotNil(err) {
			util.RespondGatewayError(ctx, endpointVO.ResponseEncode(), http.StatusRequestEntityTooLarge, err)
			return
		}

		// chamamos o próximo handler da requisição
		ctx.Next()
	}
}
