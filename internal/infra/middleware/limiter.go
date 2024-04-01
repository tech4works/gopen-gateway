package middleware

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

type limiter struct {
}

type Limiter interface {
	Do(rateLimiterProvider infra.RateLimiterProvider, sizeLimiterProvider infra.SizeLimiterProvider) api.HandlerFunc
}

func NewLimiter() Limiter {
	return limiter{}
}

func (l limiter) Do(rateLimiterProvider infra.RateLimiterProvider, sizeLimiterProvider infra.SizeLimiterProvider,
) api.HandlerFunc {
	return func(req *api.Request) {
		// aqui ja verificamos se a chave hoje sendo ela o IP está permitida
		err := rateLimiterProvider.Allow(req.HeaderValue(consts.XForwardedFor))
		if helper.IsNotNil(err) {
			req.WriteError(http.StatusTooManyRequests, err)
			return
		}

		// verificamos o tamanho da requisição, e tratamos o erro logo em seguida
		err = sizeLimiterProvider.Allow(req.Http())
		if errors.Contains(err, domainmapper.ErrHeaderTooLarge) {
			req.WriteError(http.StatusRequestHeaderFieldsTooLarge, err)
			return
		} else if helper.IsNotNil(err) {
			req.WriteError(http.StatusRequestEntityTooLarge, err)
			return
		}

		// chamamos o próximo handler da requisição
		req.Next()
	}
}
