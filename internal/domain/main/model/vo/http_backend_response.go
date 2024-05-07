package vo

import (
	"bytes"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	configEnum "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/enum"
	configVO "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
	"io"
	"net/http"
)

type HttpBackendResponse struct {
	written bool
	config  *configVO.BackendResponse
	// statusCode represents HTTP statusCode of a backend httpResponse.
	// The value of statusCode can be modified using the ModifyStatusCode() `method`.
	statusCode int
	// header represents the body fields of a backend httpResponse.
	// The value of header can be modified using the ModifyHeader() `method`.
	header Header
	// body represents the body of a backend httpResponse.
	// The value of body can be modified using the ModifyBody() `method`.
	body *Body
}

func NewHttpBackendResponse(backend *configVO.Backend, netHttpResponse *http.Response, httpRequest *HttpRequest,
	httpResponse *HttpResponse) *HttpBackendResponse {
	// construimos o header com base no netHttpResponse
	header := NewHeader(netHttpResponse.Header)

	// fazemos o parse dos bytes da resposta em para uma interface
	bodyBytes, _ := io.ReadAll(netHttpResponse.Body)

	// convertemos em body VO a partir dos bytes e do content-type
	contentType := netHttpResponse.Header.Get("Content-Type")
	contentEncoding := netHttpResponse.Header.Get("Content-Encoding")
	body := NewBody(contentType, contentEncoding, bytes.NewBuffer(bodyBytes))

	// construímos o objeto de valor do backend httpResponse
	httpBackendResponse := &HttpBackendResponse{
		config:     backend.Response(),
		statusCode: netHttpResponse.StatusCode,
		header:     header,
		body:       body,
	}

	// chamamos o applyConfig no momento EARLY
	return httpBackendResponse.ApplyConfig(configEnum.BackendResponseApplyEarly, httpRequest, httpResponse)
}

// Ok returns a boolean indicating if the statusCode of the httpBackendResponse instance is within the range 200-299.
func (b *HttpBackendResponse) Ok() bool {
	return helper.IsGreaterThanOrEqual(b.statusCode, 200) && helper.IsLessThanOrEqual(b.statusCode, 299)
}

// Key returns the key of the httpBackendResponse for aggregation.
// The key is composed of the string "backend" and the index, if it is greater than or equal to zero.
// If the httpBackendResponse has a name, the key is set to the name.
func (b *HttpBackendResponse) Key(index int) (key string) {
	// montamos o key do backend para agregar
	key = "backend"
	if helper.IsGreaterThanOrEqual(index, 0) {
		key = fmt.Sprintf("%s-%v", key, index)
	}
	// se o backend tiver informado group value
	if helper.IsNotNil(b.Config()) && b.Config().HasGroup() {
		key = b.Config().Group()
	}
	return key
}

// StatusCode returns the `statusCode` of the `httpBackendResponse` instance.
func (b *HttpBackendResponse) StatusCode() int {
	return b.statusCode
}

// Header returns the `header` of the `httpBackendResponse` instance.
func (b *HttpBackendResponse) Header() Header {
	return b.header
}

// Body returns the `body` of the `httpBackendResponse` instance.
func (b *HttpBackendResponse) Body() *Body {
	return b.body
}

func (b *HttpBackendResponse) Config() *configVO.BackendResponse {
	return b.config
}

// GroupByType returns true if the httpBackendResponse instance should be grouped,
// either by setting the groupResponse field to true or by the value of the body being a text or a slice.
// Otherwise, it returns false.
func (b *HttpBackendResponse) GroupByType() bool {
	body := b.Body()
	return helper.IsNotNil(body) && body.IsText() || helper.IsSlice(body.Bytes())
}

func (b *HttpBackendResponse) Map() any {
	var body any
	if helper.IsNotNil(body) {
		body = b.body.Interface()
	}
	return map[string]any{
		"statusCode": b.statusCode,
		"header":     b.header,
		"body":       body,
	}
}

func (b *HttpBackendResponse) ApplyConfig(momentToApply configEnum.BackendResponseApply, httpRequest *HttpRequest,
	httpResponse *HttpResponse) *HttpBackendResponse {
	// instanciamos o httpResponse config do backend
	backendResponse := b.Config()

	// se o backend ja foi aplicado, ou a config for nil, ou o momento não é o que foi configurado
	if b.Applied() || helper.IsNil(backendResponse) || helper.IsNotEqualTo(momentToApply, backendResponse.Apply()) {
		return b
	}

	// ele quer ser omitido retornamos nil
	if backendResponse.Omit() {
		return nil
	}

	// construímos o header com base nas configs do backend
	header := b.buildHeaderByConfig(backendResponse, httpRequest, httpResponse)
	// construímos o body com base nas configs do backend
	body := b.buildBodyByConfig(backendResponse, httpRequest, httpResponse)

	// construímos o novo httpBackendResponse aplicado
	return &HttpBackendResponse{
		statusCode: b.statusCode,
		header:     header,
		body:       body,
		written:    true,
	}
}

func (b *HttpBackendResponse) Applied() bool {
	return b.written
}

func (b *HttpBackendResponse) buildHeaderByConfig(backendResponse *configVO.BackendResponse, httpRequest *HttpRequest,
	httpResponse *HttpResponse) Header {
	// se ele quer omitir retornamos o header vazio
	if backendResponse.OmitHeader() {
		return NewEmptyHeader()
	}

	// primeiro obtemos o header da própria resposta do backend
	header := b.Header()
	// mapeamos o header atual com base na config
	header = header.Map(backendResponse.HeaderMapper())
	// projetamos o header atual com base na config
	header = header.Projection(backendResponse.HeaderProjection())
	// modificamos o header atual com base na config
	for _, modifier := range backendResponse.HeaderModifiers() {
		header = header.Modify(NewModify(&modifier, httpRequest, httpResponse))
	}

	// modificamos o header atual com base na config
	return header
}

func (b *HttpBackendResponse) buildBodyByConfig(backendResponse *configVO.BackendResponse, httpRequest *HttpRequest,
	httpResponse *HttpResponse) *Body {
	// se o backend quer omitir o body ou não tem body retornamos
	if backendResponse.OmitBody() || helper.IsNil(b.Body()) {
		return nil
	}

	// primeiro obtemos o body da própria resposta do backend
	body := b.Body()
	// mapeamos o body atual com base na config
	body = body.Map(backendResponse.BodyMapper())
	// projetamos o body atual com base na config
	body = body.Projection(backendResponse.BodyProjection())
	// modificamos o body atual com base na config
	for _, modifier := range backendResponse.BodyModifiers() {
		body = body.Modify(NewModify(&modifier, httpRequest, httpResponse))
	}
	// por fim, agrupamos caso tenha o campo "group"
	if backendResponse.HasGroup() {
		body = NewBodyAggregateByKey(backendResponse.Group(), body)
	}

	//retornamos o body modificado segundo a config de resposta do backend
	return body
}
