package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/handler"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type cors struct {
	corsUseCase usecase.Cors
}

type Cors interface {
	PreHandlerRequest(ctx *gin.Context)
}

func NewCors(corsUseCase usecase.Cors) Cors {
	return cors{
		corsUseCase: corsUseCase,
	}
}

func (c cors) PreHandlerRequest(ctx *gin.Context) {
	//check allow-origins
	if err := c.corsUseCase.ValidateOrigins(ctx.ClientIP()); helper.IsNotNil(err) {
		handler.RespondCodeWithError(ctx, http.StatusForbidden, err)
		return
	}
	//check allow-methods
	if err := c.corsUseCase.ValidateMethods(ctx.Request.Method); helper.IsNotNil(err) {
		handler.RespondCodeWithError(ctx, http.StatusForbidden, err)
		return
	}
	//check allow-headers
	if err := c.corsUseCase.ValidateHeaders(ctx.Request.Header); helper.IsNotNil(err) {
		handler.RespondCodeWithError(ctx, http.StatusForbidden, err)
		return
	}
}
