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

package vo

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"time"
)

// errorBody represents the structure of a httpResponse body containing error details.
// It is used to serialize the error details to JSON format.
//
// This struct is typically used in conjunction with the newBodyByError function to generate a JSON httpResponse body
// with error details based on the Endpoint and error provided.
type errorBody struct {
	// File represents the file name or path where the error occurred.
	File string `json:"file"`
	// Line represents the line number where the error occurred
	Line int `json:"line"`
	// Endpoint represents the endpoint path where the error occurred.
	Endpoint string `json:"endpoint"`
	// Message represents the error message.
	Message string `json:"message"`
	// Timestamp represents the timestamp when the error occurred.
	Timestamp time.Time `json:"timestamp"`
}

// newErrorBody creates a new errorBody object based on the provided path and error.
// It extracts details from the error using the errors.Details function, and constructs
// an errorBody struct with the relevant information. If the details are empty, it returns nil.
//
// Parameters:
//   - path: the path where the error occurred.
//   - err: the error object.
//
// Returns:
//   - a pointer to the new errorBody object, or nil if the details are empty.
func newErrorBody(path string, err error) *errorBody {
	detailsErr := errors.Details(err)
	if helper.IsNil(detailsErr) {
		return nil
	}

	return &errorBody{
		File:      detailsErr.GetFile(),
		Line:      detailsErr.GetLine(),
		Endpoint:  path,
		Message:   detailsErr.GetMessage(),
		Timestamp: time.Now(),
	}
}
