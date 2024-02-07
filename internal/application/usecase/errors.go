package usecase

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"time"
)

var MsgErrBadGateway = "error bad gateway:"
var MsgErrGatewayTimeout = "error gateway timeout:"

var ErrBadGateway = errors.New(MsgErrBadGateway)
var ErrGatewayTimeout = errors.New(MsgErrGatewayTimeout)

func errBadGateway(err error) error {
	ErrBadGateway = errors.NewSkipCaller(2, MsgErrBadGateway, err)
	return ErrBadGateway
}

func errGatewayTimeout(duration time.Duration) error {
	ErrGatewayTimeout = errors.NewSkipCaller(2, MsgErrGatewayTimeout, duration.String())
	return ErrGatewayTimeout
}
