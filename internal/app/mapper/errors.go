package mapper

import "github.com/GabrielHCataldo/go-errors/errors"

// MsgErrCacheNotFound is a string variable that holds the message "cache not found".
var MsgErrCacheNotFound = "cache not found"

// ErrCacheNotFound is an error variable representing the "cache not found" error.
var ErrCacheNotFound = errors.New(MsgErrCacheNotFound)

// NewErrCacheNotFound creates a new error of type "ErrCacheNotFound".
// It sets the value of the global variable "ErrCacheNotFound" to the error
// created using "errors.NewSkipCaller" function with skip caller value 2 and
// the message "cache not found" stored in the variable "MsgErrCacheNotFound".
// It returns the error "ErrCacheNotFound".
func NewErrCacheNotFound() error {
	ErrCacheNotFound = errors.NewSkipCaller(2, MsgErrCacheNotFound)
	return ErrCacheNotFound
}
