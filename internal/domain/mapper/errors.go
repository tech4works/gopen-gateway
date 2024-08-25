/*
 * Copyright 2024 Tech4Works
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
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"time"
)

const msgErrValueNotFound = "dynamic value not found by syntax: %s"
const msgErrInvalidAction = "Invalid action modifier %s! action: %s"
const msgErrEmptyKey = "Modifier empty key!"
const msgErrEmptyValue = "Modifier empty value!"
const msgErrIncompatibleBodyType = "Incompatible body type %s to modify!"
const msgErrBadGateway = "bad gateway error:"
const msgErrGatewayTimeout = "gateway timeout error:"
const msgErrPayloadTooLarge = "payload too large error:"
const msgErrHeaderTooLarge = "header too large error:"
const msgErrTooManyRequests = "too many requests error:"
const msgErrCacheNotFound = "cache not found"

var ErrBadGateway = errors.New(msgErrBadGateway)
var ErrGatewayTimeout = errors.New(msgErrGatewayTimeout)
var ErrPayloadTooLarge = errors.New(msgErrPayloadTooLarge)
var ErrHeaderTooLarge = errors.New(msgErrHeaderTooLarge)
var ErrTooManyRequests = errors.New(msgErrTooManyRequests)
var ErrCacheNotFound = errors.New(msgErrCacheNotFound)
var ErrValueNotFound = errors.New(msgErrValueNotFound)
var ErrInvalidAction = errors.New(msgErrInvalidAction)
var ErrEmptyKey = errors.New(msgErrEmptyKey)
var ErrEmptyValue = errors.New(msgErrEmptyValue)
var ErrIncompatibleBodyType = errors.New(msgErrIncompatibleBodyType)

func NewErrBadGateway(err error) error {
	ErrBadGateway = errors.NewSkipCaller(2, msgErrBadGateway, err)
	return ErrBadGateway
}

func NewErrGatewayTimeoutByErr(err error) error {
	ErrGatewayTimeout = errors.NewSkipCaller(2, msgErrGatewayTimeout, err)
	return ErrGatewayTimeout
}

func NewErrPayloadTooLarge(limit string) error {
	ErrPayloadTooLarge = errors.NewSkipCaller(2, msgErrPayloadTooLarge, "permitted limit is", limit)
	return ErrPayloadTooLarge
}

func NewErrHeaderTooLarge(limit string) error {
	ErrHeaderTooLarge = errors.NewSkipCaller(2, msgErrHeaderTooLarge, "permitted limit is", limit)
	return ErrHeaderTooLarge
}

func NewErrTooManyRequests(capacity int, every time.Duration) error {
	ErrTooManyRequests = errors.NewSkipCaller(2, msgErrTooManyRequests, "permitted limit is", capacity,
		"every", every.String())
	return ErrTooManyRequests
}

func NewErrCacheNotFound() error {
	ErrCacheNotFound = errors.NewSkipCaller(2, msgErrCacheNotFound)
	return ErrCacheNotFound
}

func NewErrValueNotFound(syntax string) error {
	ErrValueNotFound = errors.NewSkipCallerf(2, msgErrValueNotFound, syntax)
	return ErrValueNotFound
}

func NewErrInvalidAction(modifierName string, action enum.ModifierAction) error {
	ErrInvalidAction = errors.NewSkipCallerf(2, msgErrInvalidAction, modifierName, action)
	return ErrInvalidAction
}

func NewErrEmptyKey() error {
	ErrEmptyKey = errors.NewSkipCallerf(2, msgErrEmptyKey)
	return ErrEmptyKey
}

func NewErrEmptyValue() error {
	ErrEmptyValue = errors.NewSkipCaller(2, msgErrEmptyValue)
	return ErrEmptyValue
}

func NewErrIncompatibleBodyType(contentType string) error {
	ErrIncompatibleBodyType = errors.NewSkipCallerf(2, msgErrIncompatibleBodyType, contentType)
	return ErrIncompatibleBodyType
}
