package infra

import (
	berrors "errors"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"net/http"
	"net/url"
	"syscall"
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

// treatHttpClientErr handles an error that occurred during an HTTP request made by the restTemplate.
// It takes an error as input and returns the corresponding error after handling it, if any.
// If the input error is nil, it returns nil.
// If the input error is a connection refused error or host down error, it creates a new domainmapper.ErrBadGateway error and returns it.
// If the input error is not nil, it checks if it is an url.Error and if it has a timeout.
// If it has a timeout, it creates a new domainmapper.ErrGatewayTimeout error and returns it.
// For any other type of error, it returns the error as it is.
func (r restTemplate) treatHttpClientErr(err error) error {
	// se tiver nil, retornamos nil
	if helper.IsNil(err) {
		return nil
	}

	// caso ocorra algum erro, tratamos
	if errors.Contains(err, syscall.ECONNREFUSED) || errors.Contains(err, syscall.EHOSTDOWN) {
		err = domainmapper.NewErrBadGateway(err)
	} else if helper.IsNotNil(err) {
		var urlErr *url.Error
		berrors.As(err, &urlErr)
		if urlErr.Timeout() {
			err = domainmapper.NewErrGatewayTimeoutByErr(err)
		}
	}

	// retornamos o erro tratado, ou não
	return err
}
