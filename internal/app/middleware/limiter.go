package middleware

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/external"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/util"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type limiter struct {
}

type Limiter interface {
	Do(rateLimiterProvider external.RateLimiterProvider, sizeLimiterProvider external.SizeLimiterProvider) gin.HandlerFunc
}

func NewLimiter() Limiter {
	return limiter{}
}

func (l limiter) Do(rateLimiterProvider external.RateLimiterProvider, sizeLimiterProvider external.SizeLimiterProvider,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// aqui ja verificamos se a chave hoje sendo ela o IP está permitida
		err := rateLimiterProvider.Allow(ctx.GetHeader("X-Forwarded-For"))
		if helper.IsNotNil(err) {
			util.RespondCodeWithError(ctx, http.StatusTooManyRequests, err)
			return
		}

		// verificamos o tamanho da requisição, e tratamos o erro logo em seguida
		err = sizeLimiterProvider.Allow(ctx.Request)
		if errors.Contains(err, mapper.ErrHeaderTooLarge) {
			util.RespondCodeWithError(ctx, http.StatusRequestHeaderFieldsTooLarge, err)
			return
		} else if helper.IsNotNil(err) {
			util.RespondCodeWithError(ctx, http.StatusRequestEntityTooLarge, err)
			return
		}

		// chamamos o próximo handler da requisição
		ctx.Next()
	}
}
