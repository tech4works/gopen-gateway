package middleware

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"time"
)

type timeout struct {
}

type Timeout interface {
	Do(timeoutDuration time.Duration) gin.HandlerFunc
}

func NewTimeout() Timeout {
	return timeout{}
}

func (t timeout) Do(timeoutDuration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// inicializamos o context com timeout fornecido na config do gateway
		timeoutContext, cancel := context.WithTimeout(ctx.Request.Context(), timeoutDuration)
		defer cancel()

		// setamos esse context na request atual para propagar para os outros manipuladores
		ctx.Request = ctx.Request.WithContext(timeoutContext)

		// criamos os canais de alerta
		finishChan := make(chan interface{}, 1)
		panicChan := make(chan interface{}, 1)

		go func() {
			// chamamos o panic recovery caso ocorra
			defer func() {
				if p := recover(); helper.IsNotNil(p) {
					logger.Errorf("%s: %s", p, string(debug.Stack()))
					panicChan <- p
				}
			}()
			// chamamos o próximo handler na requisição
			ctx.Next()
			// se finalizou a tempo, chamamos o channel para seguir normalmente
			finishChan <- struct{}{}
		}()

		// inicializamos as variáveis para serem utilizadas ou não no futuro
		var statusCode int
		var err error

		// seguramos o goroutine principal aguardando os canais ou o context serem notificados
		select {
		case <-finishChan:
			break
		case <-panicChan:
			statusCode = http.StatusInternalServerError
			err = errors.New("panic error occurred")
			break
		case <-ctx.Done():
			statusCode = http.StatusGatewayTimeout
			err = errors.New("gateway timeout: ", timeoutDuration.String())
			break
		}

		// caso tenha passado nos dois fluxos de timeout ou de erro, respondemos à requisição
		if helper.IsGreaterThan(statusCode, 0) {
			util.RespondCodeWithError(ctx, statusCode, err)
		}
	}
}

//func (t timeout) buildTimeoutDuration(timeoutDuration time.Duration) time.Duration {
//	// caso o valor padrão não tenha sido informado, setamos 30s como
//	if helper.IsLessThanOrEqual(t.timeoutDurationDefault, 0) {
//		t.timeoutDurationDefault = 30 * time.Second
//	}
//	// caso timeout no DTO do endpoint seja vazio, usamos o valor padrão
//	if helper.IsLessThanOrEqual(timeoutDuration, 0) {
//		return t.timeoutDurationDefault
//	}
//	return timeoutDuration
//}
