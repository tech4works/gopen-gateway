package middleware

import (
	"bytes"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/handler"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/gin-gonic/gin"
	"time"
)

type log struct {
	logUseCase usecase.Log
}

type Log interface {
	PreHandlerRequest(ctx *gin.Context)
}

func NewLog(loggerUseCase usecase.Log) Log {
	return log{
		logUseCase: loggerUseCase,
	}
}

func (l log) PreHandlerRequest(ctx *gin.Context) {
	writer := &handler.ResponseWriter{Body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
	ctx.Writer = writer

	startTime := time.Now()
	l.logUseCase.PrintLogRequest(ctx.Request)
	ctx.Next()
	l.logUseCase.PrintLogResponse(ctx.Request, *writer, startTime)
}
