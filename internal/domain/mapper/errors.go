package mapper

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"time"
)

var MsgErrBadGateway = "error bad gateway:"
var MsgErrGatewayTimeout = "error gateway timeout:"
var MsgErrPayloadTooLarge = "error payload too large:"
var MsgErrHeaderTooLarge = "error header too large:"
var MsgErrTooManyRequests = "error too many requests:"

var ErrBadGateway = errors.New(MsgErrBadGateway)
var ErrGatewayTimeout = errors.New(MsgErrGatewayTimeout)
var ErrPayloadTooLarge = errors.New(MsgErrPayloadTooLarge)
var ErrHeaderTooLarge = errors.New(MsgErrHeaderTooLarge)
var ErrTooManyRequests = errors.New(MsgErrTooManyRequests)

func NewErrBadGateway(err error) error {
	ErrBadGateway = errors.NewSkipCaller(2, MsgErrBadGateway, err)
	return ErrBadGateway
}

func NewErrGatewayTimeout(duration time.Duration) error {
	ErrGatewayTimeout = errors.NewSkipCaller(2, MsgErrGatewayTimeout, duration.String())
	return ErrGatewayTimeout
}

func NewErrGatewayTimeoutByErr(err error) error {
	ErrGatewayTimeout = errors.NewSkipCaller(2, MsgErrGatewayTimeout, err)
	return ErrGatewayTimeout
}

func NewErrPayloadTooLarge(limit string) error {
	ErrPayloadTooLarge = errors.NewSkipCaller(2, MsgErrPayloadTooLarge, "permitted limit is", limit)
	return ErrPayloadTooLarge
}

func NewErrHeaderTooLarge(limit string) error {
	ErrHeaderTooLarge = errors.NewSkipCaller(2, MsgErrHeaderTooLarge, "permitted limit is", limit)
	return ErrHeaderTooLarge
}

func NewErrTooManyRequests(capacity int, every time.Duration) error {
	ErrTooManyRequests = errors.NewSkipCaller(2, MsgErrTooManyRequests, "permitted limit is", capacity,
		"every", every.String())
	return ErrTooManyRequests
}
