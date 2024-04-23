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

package interfaces

import (
	"net/http"
)

// RestTemplate is an interface that represents a template for making HTTP requests.
// It provides a method MakeRequest for sending an HTTP request and returning the corresponding
// HTTP response and an error, if any.
type RestTemplate interface {
	// MakeRequest sends an HTTP request and returns the corresponding HTTP response and error.
	// It takes an HTTP request object as a parameter.
	// The function's steps are as follows:
	//
	//  1. The function sends the HTTP request using a REST client.
	//  2. If the operation fails, the function returns an error.
	//  3. Otherwise, it returns the HTTP response.
	//
	// Parameters:
	// httpRequest: the HTTP request to be sent.
	//
	// Returns:
	// An HTTP response object and an error.
	MakeRequest(httpRequest *http.Request) (*http.Response, error)
}
