package middleware

import (
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/enum"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

type header struct {
	traceUseCase usecase.Trace
}

type Header interface {
	PreHandlerRequest(ctx *gin.Context)
}

func NewHeader(traceUseCase usecase.Trace) Log {
	return header{
		traceUseCase: traceUseCase,
	}
}

func (h header) PreHandlerRequest(ctx *gin.Context) {
	ctx.Header(enum.XTraceId, h.traceUseCase.GenerateTraceId())
	ctx.Header(enum.XForwardedFor, ctx.ClientIP())
	ctx.Request.Header.Set(enum.XTraceId, h.traceUseCase.GenerateTraceId())
	ctx.Request.Header.Set(enum.XForwardedFor, ctx.ClientIP())
}
