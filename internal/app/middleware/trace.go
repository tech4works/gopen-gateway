package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/gin-gonic/gin"
)

type trace struct {
	traceProvider interfaces.TraceProvider
}

type Trace interface {
	Do(ctx *gin.Context)
}

func NewTrace(traceProvider interfaces.TraceProvider) Trace {
	return trace{
		traceProvider: traceProvider,
	}
}

func (t trace) Do(ctx *gin.Context) {
	// adicionamos na requisição o X-Forwarded-For
	ctx.Request.Header.Add(consts.XForwardedFor, ctx.ClientIP())
	// caso não tenha trace id informado, setamos
	if helper.IsEmpty(ctx.GetHeader(consts.XTraceId)) {
		ctx.Request.Header.Set(consts.XTraceId, t.traceProvider.GenerateTraceId())
	}
	// seguimos para a próxima func da requisição
	ctx.Next()
}
