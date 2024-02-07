package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/util"
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
	logger.Info("Start request on URI:", ctx.Request.RequestURI)
	//check allow-origins
	if err := c.corsUseCase.ValidateOrigins(ctx.ClientIP()); helper.IsNotNil(err) {
		util.RespondCodeWithError(ctx, http.StatusForbidden, err)
		return
	}
	//check allow-methods
	if err := c.corsUseCase.ValidateMethods(ctx.Request.Method); helper.IsNotNil(err) {
		util.RespondCodeWithError(ctx, http.StatusForbidden, err)
		return
	}
	//check allow-headers
	if err := c.corsUseCase.ValidateHeaders(ctx.Request.Header); helper.IsNotNil(err) {
		util.RespondCodeWithError(ctx, http.StatusForbidden, err)
		return
	}
	logger.Info("Finish cors validate on URI:", ctx.Request.RequestURI)
}
