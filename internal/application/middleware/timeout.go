package middleware

import (
	"context"
	"github.com/GabrielHCataldo/go-error-detail/errors"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/handler"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"time"
)

type timeout struct {
	duration time.Duration
}

type Timeout interface {
	PreHandlerRequest(ctx *gin.Context)
}

func NewTimeout(duration time.Duration) Timeout {
	return timeout{duration: duration}
}

func (t timeout) PreHandlerRequest(ctx *gin.Context) {
	ctxTimeout, cancel := context.WithTimeout(ctx.Request.Context(), t.duration)
	defer cancel()
	ctx.Request = ctx.Request.WithContext(ctxTimeout)
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				logger.Errorf("%s: %s", p, debug.Stack())
				panicChan <- p
			}
		}()
		ctx.Next()
		finish <- struct{}{}
	}()
	select {
	case <-panicChan:
		handler.RespondCode(ctx, http.StatusInternalServerError)
	case <-finish:
	case <-ctxTimeout.Done():
		handler.RespondCodeWithError(ctx, http.StatusGatewayTimeout, errors.New(
			"gateway timeout configured value:", t.duration.String(),
		))
	}
}
