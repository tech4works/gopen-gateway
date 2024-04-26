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

package infra

import (
	berrors "errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"net/http"
	"net/url"
)

type restTemplate struct {
}

// NewRestTemplate returns a new instance of a restTemplate object.
// It implements the interfaces.RestTemplate interface.
func NewRestTemplate() interfaces.RestTemplate {
	return restTemplate{}
}

// MakeRequest sends an HTTP request using the restTemplate.
// It takes an *http.Request as input and returns the corresponding *http.Response and an error, if any.
// If there is an error during the HTTP request, it is handled before returning the response.
// The error is treated depending on the type of error that occurred. If the error is a connection refused error,
// then a domainmapper.ErrBadGateway error is created and returned. If the error is a timeout error,
// then a domainmapper.ErrGatewayTimeout error is created and returned. For any other type of error,
// the error is returned as it is.
func (r restTemplate) MakeRequest(httpRequest *http.Request) (*http.Response, error) {
	// fazemos a requisição http
	httpClient := http.Client{}
	httpResponse, err := httpClient.Do(httpRequest)
	// caso ocorra um erro, tratamos, caso contrario retornamos a resposta
	return httpResponse, r.treatHttpClientErr(err)
}

// treatHttpClientErr handles the error returned by httpClient.Do method.
// If the error is nil, it returns nil.
// If the error is not nil, it checks if it is an url.Error and if it has a timeout.
// If it has a timeout, it creates a new domainmapper.ErrGatewayTimeout error and returns it.
// Otherwise, it creates a new domainmapper.ErrBadGateway error and returns it.
// If the error is not an url.Error, it returns the error as it is.
func (r restTemplate) treatHttpClientErr(err error) error {
	// se tiver nil, retornamos nil
	if helper.IsNil(err) {
		return nil
	}

	// caso ocorra algum erro, tratamos
	if helper.IsNotNil(err) {
		var urlErr *url.Error
		berrors.As(err, &urlErr)
		if urlErr.Timeout() {
			err = domainmapper.NewErrGatewayTimeoutByErr(err)
		} else {
			err = domainmapper.NewErrBadGateway(err)
		}
	}

	// retornamos o erro tratado, ou não
	return err
}
