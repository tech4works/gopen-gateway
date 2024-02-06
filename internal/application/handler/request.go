package handler

import (
	"github.com/GabrielHCataldo/go-error-detail/errors"
	"github.com/gin-gonic/gin"
)

func RespondCode(ctx *gin.Context, code int) {
	if ctx.IsAborted() {
		return
	}
	ctx.Status(code)
	ctx.Abort()
}

func RespondCodeWithError(ctx *gin.Context, code int, err *errors.ErrorDetail) {
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(code, fillErrorObject(ctx, err))
	ctx.Abort()
}

func fillErrorObject(ctx *gin.Context, err *errors.ErrorDetail) errors.ErrorDetail {
	err.Endpoint = ctx.Request.RequestURI
	return *err
}
