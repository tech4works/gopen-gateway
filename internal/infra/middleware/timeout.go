package middleware

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
	"runtime/debug"
	"time"
)

type timeout struct {
}

type Timeout interface {
	Do(timeoutDuration time.Duration) api.HandlerFunc
}

// NewTimeout returns a new instance of the `timeout` type that implements the `Timeout` interface.
func NewTimeout() Timeout {
	return timeout{}
}

// Do takes a timeoutDuration of type time.Duration and returns a function of type api.HandlerFunc.
// The returned function initializes a context with the provided timeout duration from the gateway's config.
// It sets this context in the current request to propagate it to other handlers.
// It creates two channels for alerting - finishChan and panicChan.
// It then spawns a goroutine that calls the panic recovery function in case of a panic and calls the next handler in the request.
// If the request finishes before the timeout, it sends a signal to the finishChan.
// If a panic occurs, it sends the panic value to the panicChan.
// If the context timeout occurs before the request finishes, it sets the statusCode to http.StatusGatewayTimeout and the error to "gateway timeout: <timeoutDuration>".
// It then waits for one of the three channels to receive a signal by using a select statement.
// Finally, it checks whether the statusCode is greater than 0 and writes the error response using ctx.WriteError if true.
func (t timeout) Do(timeoutDuration time.Duration) api.HandlerFunc {
	return func(ctx *api.Context) {
		// inicializamos o context com timeout fornecido na config do gateway
		timeoutContext, cancel := context.WithTimeout(ctx.Context(), timeoutDuration)
		defer cancel()

		// setamos esse context na request atual para propagar para os outros manipuladores
		ctx.SetRequestContext(timeoutContext)

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
		case <-ctx.Context().Done():
			statusCode = http.StatusGatewayTimeout
			err = errors.New("gateway timeout:", timeoutDuration.String())
			break
		}

		// caso tenha passado nos dois fluxos de timeout ou de erro, respondemos à requisição
		if helper.IsGreaterThan(statusCode, 0) {
			ctx.WriteError(statusCode, err)
		}
	}
}
