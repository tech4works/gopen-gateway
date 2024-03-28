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

func NewRestTemplate() interfaces.RestTemplate {
	return restTemplate{}
}

func (r restTemplate) MakeRequest(httpRequest *http.Request) (*http.Response, error) {
	// fazemos a requisição http
	httpResponse, err := http.DefaultClient.Do(httpRequest)
	// caso ocorra um erro, tratamos, caso contrario retornamos a resposta
	return httpResponse, r.treatHttpClientErr(err)
}

func (r restTemplate) treatHttpClientErr(err error) error {
	// se tiver nil, retornamos nil
	if helper.IsNil(err) {
		return nil
	}

	// caso ocorra algum erro, tratamos
	if errors.Contains(err, syscall.ECONNREFUSED) {
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
