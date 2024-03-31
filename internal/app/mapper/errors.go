package mapper

import "github.com/GabrielHCataldo/go-errors/errors"

var MsgErrCacheNotFound = "cache not found"

var ErrCacheNotFound = errors.New(MsgErrCacheNotFound)

func NewErrCacheNotFound() error {
	ErrCacheNotFound = errors.NewSkipCaller(2, MsgErrCacheNotFound)
	return ErrCacheNotFound
}
