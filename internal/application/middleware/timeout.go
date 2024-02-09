package middleware

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/handler"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type timeout struct {
	timeoutUseCase usecase.Timeout
}

type Timeout interface {
	PreHandlerRequest(ctx *gin.Context)
}

func NewTimeout(timeoutUseCase usecase.Timeout) Timeout {
	return timeout{
		timeoutUseCase: timeoutUseCase,
	}
}

func (t timeout) PreHandlerRequest(ctx *gin.Context) {
	err := t.timeoutUseCase.SetRequestContextTimeout(ctx.Request, func() {
		ctx.Next()
	})
	if errors.Contains(err, usecase.ErrGatewayTimeout) {
		handler.RespondCodeWithError(ctx, http.StatusGatewayTimeout, err)
	} else if helper.IsNotNil(err) {
		handler.RespondCodeWithError(ctx, http.StatusInternalServerError, err)
	}
}
