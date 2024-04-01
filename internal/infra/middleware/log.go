package middleware

import (
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"time"
)

type log struct {
	logProvider infra.LogProvider
}

type Log interface {
	Do(req *api.Request)
}

func NewLog(logProvider infra.LogProvider) Log {
	return log{
		logProvider: logProvider,
	}
}

func (l log) Do(req *api.Request) {
	// mantemos o tempo que a requisição começou
	startTime := time.Now()

	// inicializamos a logger options global, com o traceId e XForwardedFor
	l.logProvider.InitializeLoggerOptions(req)

	// imprimimos o log de start
	logger.Info("Start!", l.logProvider.BuildInitialRequestMessage(req))

	// chamamos o próximo handler da requisição
	req.Next()

	// imprimimos o log de finish
	logger.Info("Finish!", l.logProvider.BuildFinishRequestMessage(req.Writer(), startTime))
}
