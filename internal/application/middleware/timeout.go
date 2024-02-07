package middleware

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/util"
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
		util.RespondCodeWithError(ctx, http.StatusGatewayTimeout, err)
	} else if helper.IsNotNil(err) {
		util.RespondCodeWithError(ctx, http.StatusInternalServerError, err)
	}
}
