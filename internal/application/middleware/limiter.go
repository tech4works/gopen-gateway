package middleware

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/util"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type limiter struct {
	limiterUseCase usecase.Limiter
}

type Limiter interface {
	PreHandlerRequest(ctx *gin.Context)
}

func NewLimiter(limiterUseCase usecase.Limiter) Limiter {
	return limiter{
		limiterUseCase: limiterUseCase,
	}
}

func (l limiter) PreHandlerRequest(ctx *gin.Context) {
	// checa rate limite de requisição pelo ip
	if err := l.limiterUseCase.ValidateRate(ctx.ClientIP()); helper.IsNotNil(err) {
		util.RespondCodeWithError(ctx, http.StatusTooManyRequests, err)
		return
	}
	// checa tamanho do body da requisição
	headerContentType := ctx.GetHeader("Content-Type")
	body := ctx.Request.Body
	bodyBytes, err := l.limiterUseCase.ValidateRequestSize(headerContentType, body)
	if helper.IsNotNil(err) {
		util.RespondCodeWithError(ctx, http.StatusTooManyRequests, err)
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	ctx.Next()
}
