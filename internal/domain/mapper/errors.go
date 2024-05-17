/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mapper

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"time"
)

// MsgErrBadGateway represents the error message for a bad gateway error.
// The constant value is "bad gateway error:".
var MsgErrBadGateway = "bad gateway error:"

// MsgErrGatewayTimeout represents the error message for a gateway timeout error.
// The constant value is "gateway timeout error:".
var MsgErrGatewayTimeout = "gateway timeout error:"

// MsgErrPayloadTooLarge represents the error message for a payload that is
// too large.
// The constant value is "payload too large error:".
var MsgErrPayloadTooLarge = "payload too large error:"

// MsgErrHeaderTooLarge represents the error message for a header too large error.
// The constant value is "header too large error:".
var MsgErrHeaderTooLarge = "header too large error:"

// MsgErrTooManyRequests represents the error message for a too many requests error.
// The constant value is "too many requests error:".
var MsgErrTooManyRequests = "too many requests error:"

// MsgErrCacheNotFound is a string variable that holds the message "cache not found".
var MsgErrCacheNotFound = "cache not found"

// ErrBadGateway represents an error indicating a bad gateway.
var ErrBadGateway = errors.New(MsgErrBadGateway)

// ErrGatewayTimeout represents the error for a gateway timeout.
var ErrGatewayTimeout = errors.New(MsgErrGatewayTimeout)

// ErrPayloadTooLarge represents the payload too large error.
var ErrPayloadTooLarge = errors.New(MsgErrPayloadTooLarge)

// ErrHeaderTooLarge represents the error for a header that is too large.
var ErrHeaderTooLarge = errors.New(MsgErrHeaderTooLarge)

// ErrTooManyRequests represents the error for when there are too many requests.
var ErrTooManyRequests = errors.New(MsgErrTooManyRequests)

// ErrCacheNotFound is an error variable representing the "cache not found" error.
var ErrCacheNotFound = errors.New(MsgErrCacheNotFound)

// NewErrBadGateway creates a new domainmapper.ErrBadGateway error with the specified error as the cause.
func NewErrBadGateway(err error) error {
	ErrBadGateway = errors.NewSkipCaller(2, MsgErrBadGateway, err)
	return ErrBadGateway
}

// NewErrGatewayTimeoutByErr creates a new domainmapper.ErrGatewayTimeout error with the specified error as the cause.
// It takes an error as input and returns the corresponding error after handling it, if any. If the input error is not nil,
// it checks if it is an url.Error and if it has a timeout. If it has a timeout, it creates a new domainmapper.ErrGatewayTimeout
// error and returns it. For any other type of error, it returns the error as it is.
func NewErrGatewayTimeoutByErr(err error) error {
	ErrGatewayTimeout = errors.NewSkipCaller(2, MsgErrGatewayTimeout, err)
	return ErrGatewayTimeout
}

// NewErrPayloadTooLarge creates a new domainmapper.ErrPayloadTooLarge error with the specified limit as the permitted limit.
func NewErrPayloadTooLarge(limit string) error {
	ErrPayloadTooLarge = errors.NewSkipCaller(2, MsgErrPayloadTooLarge, "permitted limit is", limit)
	return ErrPayloadTooLarge
}

// NewErrHeaderTooLarge creates a new domainmapper.ErrHeaderTooLarge error with the specified limit as the permitted limit.
func NewErrHeaderTooLarge(limit string) error {
	ErrHeaderTooLarge = errors.NewSkipCaller(2, MsgErrHeaderTooLarge, "permitted limit is", limit)
	return ErrHeaderTooLarge
}

// NewErrTooManyRequests creates a new domainmapper.ErrTooManyRequests error with the specified capacity and every value.
func NewErrTooManyRequests(capacity int, every time.Duration) error {
	ErrTooManyRequests = errors.NewSkipCaller(2, MsgErrTooManyRequests, "permitted limit is", capacity,
		"every", every.String())
	return ErrTooManyRequests
}

// NewErrCacheNotFound creates a new error of type "ErrCacheNotFound".
// It sets the value of the global variable "ErrCacheNotFound" to the error
// created using "errors.NewSkipCaller" function with skip caller value 2 and
// the message "cache not found" stored in the variable "MsgErrCacheNotFound".
// It returns the error "ErrCacheNotFound".
func NewErrCacheNotFound() error {
	ErrCacheNotFound = errors.NewSkipCaller(2, MsgErrCacheNotFound)
	return ErrCacheNotFound
}
