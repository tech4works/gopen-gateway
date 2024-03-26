package middleware

import (
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/external"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/gin-gonic/gin"
	"time"
)

type log struct {
	logProvider external.LogProvider
}

type Log interface {
	Do(ctx *gin.Context)
}

func NewLog() Log {
	return log{}
}

func (l log) Do(ctx *gin.Context) {
	// mantemos o tempo que a requisição começou
	startTime := time.Now()

	// inicializamos o writer de resposta
	responseWriter := mapper.BuildResponseWriter(ctx)
	ctx.Writer = responseWriter

	// inicializamos a logger options global, com o traceId e XForwardedFor
	l.logProvider.InitializeLoggerOptions(ctx)

	// imprimimos o log de start
	logger.Info("Start!", l.logProvider.BuildInitialRequestMessage(ctx))

	// chamamos o próximo handler da requisição
	ctx.Next()

	// imprimimos o log de finish
	logger.Info("Finish!", l.logProvider.BuildFinishRequestMessage(*responseWriter, startTime))
}
