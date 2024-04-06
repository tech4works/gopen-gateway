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

// NewLimiter creates a new instance of Limiter.
// It returns a Limiter value which implements the Limiter interface.
func NewLimiter() Limiter {
	return limiter{}
}

// Do execute a handler that implements the api.HandlerFunc interface, providing rate limiting and size limiting functionality.
// It takes a RateLimiterProvider instance and a SizeLimiterProvider instance as input parameters.
// The rateLimiterProvider is used to check whether the request is allowed based on the rate limit.
// The sizeLimiterProvider is used to check whether the request size is within the allowed limit.
// If the request is not allowed, it returns an error and writes an error response to the request.
// If the request is allowed, it calls the Next() method of the Request object to execute the next handler in the chain.
func (l limiter) Do(rateLimiterProvider infra.RateLimiterProvider, sizeLimiterProvider infra.SizeLimiterProvider,
) api.HandlerFunc {
	return func(ctx *api.Context) {
		// aqui ja verificamos se a chave hoje sendo ela o IP está permitida
		err := rateLimiterProvider.Allow(ctx.HeaderValue(consts.XForwardedFor))
		if helper.IsNotNil(err) {
			ctx.WriteError(http.StatusTooManyRequests, err)
			return
		}

		// verificamos o tamanho da requisição, e tratamos o erro logo em seguida
		err = sizeLimiterProvider.Allow(ctx.Http())
		if errors.Contains(err, domainmapper.ErrHeaderTooLarge) {
			ctx.WriteError(http.StatusRequestHeaderFieldsTooLarge, err)
			return
		} else if helper.IsNotNil(err) {
			ctx.WriteError(http.StatusRequestEntityTooLarge, err)
			return
		}

		// chamamos o próximo handler da requisição
		ctx.Next()
	}
}
