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
