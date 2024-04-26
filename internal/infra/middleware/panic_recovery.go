package middleware

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
	"runtime/debug"
)

type panicRecovery struct {
}

// PanicRecovery represents a type that can recover from panics.
// It provides a method called Do that takes a *api.Context as a parameter and does the necessary recovery actions.
type PanicRecovery interface {
	// Do apply the necessary recovery actions for a given *api.Context parameter.
	Do(ctx *api.Context)
}

// NewPanicRecovery returns a PanicRecovery implementation.
// The returned PanicRecovery contains a Do method that recovers from panics and handles the recovery process.
// The Do method takes a *api.Context as a parameter and performs the necessary recovery actions.
// The recovery actions include logging the recovered panic and stack trace, and writing an error response to the context.
// The Do method also calls ctx.Next() to proceed to the next request handling.
func NewPanicRecovery() PanicRecovery {
	return panicRecovery{}
}

// Do recovers from panics and handles the recovery process.
// It takes a *api.Context as a parameter and performs the necessary recovery actions.
// The recovery actions include logging the recovered panic and stack trace,
// and writing an error response to the context.
// It also calls ctx.Next() to proceed to the next request handling.
func (p panicRecovery) Do(ctx *api.Context) {
	defer func() {
		if r := recover(); helper.IsNotNil(r) {
			logger.Errorf("%s:%s", r, string(debug.Stack()))
			err := errors.New("gateway panic error occurred! detail:", r)
			ctx.WriteError(http.StatusInternalServerError, err)
		}
	}()
	// processa a próxima manipulação de requisição
	ctx.Next()
}
