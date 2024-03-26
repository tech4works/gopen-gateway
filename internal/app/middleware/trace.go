package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/external"
	"github.com/gin-gonic/gin"
)

type trace struct {
	traceProvider external.TraceProvider
}

type Trace interface {
	Do(ctx *gin.Context)
}

func NewTrace(traceProvider external.TraceProvider) Trace {
	return trace{
		traceProvider: traceProvider,
	}
}

func (t trace) Do(ctx *gin.Context) {
	// adicionamos na requisição o X-Forwarded-For
	ctx.Request.Header.Add("X-Forwarded-For", ctx.ClientIP())
	// caso não tenha trace id informado, setamos
	if helper.IsEmpty(ctx.GetHeader("X-Trace-Id")) {
		ctx.Request.Header.Set("X-Trace-Id", t.traceProvider.GenerateTraceId())
	}
	// seguimos para a próxima func da requisição
	ctx.Next()
}
