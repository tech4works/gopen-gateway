package usecase

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"net/http"
	"runtime/debug"
	"time"
)

type timeout struct {
	handlerTimeout time.Duration
}

type Timeout interface {
	SetRequestContextTimeout(request *http.Request, funcCall func()) error
}

func NewTimeout(handlerTimeout time.Duration) Timeout {
	return timeout{
		handlerTimeout: handlerTimeout,
	}
}

func (t timeout) SetRequestContextTimeout(request *http.Request, funcAsync func()) error {
	ctx, cancel := context.WithTimeout(context.TODO(), t.handlerTimeout)
	defer cancel()
	*request = *request.WithContext(ctx)
	finish := make(chan interface{}, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); helper.IsNotNil(p) {
				logger.Errorf("%s:%s", p, debug.Stack())
				panicChan <- p
			}
		}()
		funcAsync()
		finish <- struct{}{}
	}()
	select {
	case <-finish:
		return nil
	case <-panicChan:
		return errors.New("panic error occurred")
	case <-ctx.Done():
		return errGatewayTimeout(t.handlerTimeout)
	}
}
