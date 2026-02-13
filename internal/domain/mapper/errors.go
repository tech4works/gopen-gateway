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
	"time"

	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

const (
	msgErrValueNotFound        = "dynamic value not found by syntax:"
	msgErrInvalidAction        = "Invalid action modifier, action:"
	msgErrEmptyKey             = "Modifier empty key!"
	msgErrEmptyValue           = "Modifier empty value!"
	msgErrIncompatibleBodyType = "Incompatible body type to modify:"
	msgErrBadGateway           = "bad gateway error:"
	msgErrGatewayTimeout       = "gateway timeout error:"
	msgErrPayloadTooLarge      = "payload too large error:"
	msgErrHeaderTooLarge       = "header too large error:"
	msgErrTooManyRequests      = "too many requests error:"
	msgErrCacheNotFound        = "cache not found"
	msgErrConcurrentCanceled   = "concurrent context canceled"
	msgErrMapperIgnored        = "mapper ignored by expression:"
)

var (
	ErrBadGateway           = errors.New(msgErrBadGateway)
	ErrGatewayTimeout       = errors.New(msgErrGatewayTimeout)
	ErrPayloadTooLarge      = errors.New(msgErrPayloadTooLarge)
	ErrHeaderTooLarge       = errors.New(msgErrHeaderTooLarge)
	ErrTooManyRequests      = errors.New(msgErrTooManyRequests)
	ErrCacheNotFound        = errors.New(msgErrCacheNotFound)
	ErrValueNotFound        = errors.New(msgErrValueNotFound)
	ErrInvalidAction        = errors.New(msgErrInvalidAction)
	ErrEmptyKey             = errors.New(msgErrEmptyKey)
	ErrEmptyValue           = errors.New(msgErrEmptyValue)
	ErrIncompatibleBodyType = errors.New(msgErrIncompatibleBodyType)
	ErrConcurrentCanceled   = errors.New(msgErrConcurrentCanceled)
	ErrMapperIgnored        = errors.New(msgErrMapperIgnored)
)

func NewErrBadGateway(err error) error {
	return errors.NewWithSkipCaller(2, msgErrBadGateway, err)
}

func NewErrGatewayTimeoutByErr(err error) error {
	return errors.NewWithSkipCaller(2, msgErrGatewayTimeout, err)
}

func NewErrConcurrentCanceled() error {
	return errors.NewWithSkipCaller(2, msgErrConcurrentCanceled)
}

func NewErrPayloadTooLarge(limit string) error {
	return errors.NewWithSkipCaller(2, msgErrPayloadTooLarge, "permitted limit is", limit)
}

func NewErrHeaderTooLarge(limit string) error {
	return errors.NewWithSkipCaller(2, msgErrHeaderTooLarge, "permitted limit is", limit)
}

func NewErrTooManyRequests(capacity int, every time.Duration) error {
	return errors.NewWithSkipCaller(2, msgErrTooManyRequests, "permitted limit is", capacity, "every", every.String())
}

func NewErrCacheNotFound() error {
	return errors.NewWithSkipCaller(2, msgErrCacheNotFound)
}

func NewErrValueNotFound(syntax string) error {
	return errors.NewWithSkipCaller(2, msgErrValueNotFound, syntax)
}

func NewErrInvalidAction(modifierName string, action enum.ModifierAction) error {
	return errors.NewWithSkipCaller(2, msgErrInvalidAction, modifierName, action)
}

func NewErrEmptyKey() error {
	return errors.NewWithSkipCaller(2, msgErrEmptyKey)
}

func NewErrEmptyValue() error {
	return errors.NewWithSkipCaller(2, msgErrEmptyValue)
}

func NewErrIncompatibleBodyType(contentType string) error {
	return errors.NewWithSkipCaller(2, msgErrIncompatibleBodyType, contentType)
}

func NewErrMapperIgnored(expression string) error {
	return errors.NewWithSkipCaller(2, msgErrMapperIgnored, expression)
}
