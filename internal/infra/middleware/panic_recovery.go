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

type PanicRecovery interface {
	Do(ctx *api.Context)
}

func NewPanicRecovery() PanicRecovery {
	return panicRecovery{}
}

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
