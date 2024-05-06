package vo

import (
	"bytes"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"io"
	"net/http"
)

// httpBackendResponse represents a httpResponse from a backend service.
type httpBackendResponse struct {
	applied bool
	config  *BackendResponse
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

func NewHttpBackendResponse(backendVO *Backend, httpResponse *http.Response) *httpBackendResponse {
	// fazemos o parse dos bytes da resposta em para uma interface
	bodyBytes, _ := io.ReadAll(httpResponse.Body)

	// convertemos em body VO a partir dos bytes e do content-type
	contentType := httpResponse.Header.Get("Content-Type")
	contentEncoding := httpResponse.Header.Get("Content-Encoding")
	body := NewBody(contentType, contentEncoding, bytes.NewBuffer(bodyBytes))

	// construímos o objeto de valor do backend httpResponse
	httpBackendResponseVO := &httpBackendResponse{
		config:     backendVO.Response(),
		statusCode: httpResponse.StatusCode,
		header:     NewHeader(httpResponse.Header),
		body:       body,
	}

	// tentamos aplicar as configurações no moment early
	return httpBackendResponseVO.ApplyConfig(enum.BackendResponseApplyEarly)
}

func (b *httpBackendResponse) ModifyStatusCode(statusCode int) *httpBackendResponse {
	return &httpBackendResponse{
		config:     b.config,
		statusCode: statusCode,
		header:     b.header,
		body:       b.body,
	}
}

func (b *httpBackendResponse) ModifyHeader(header Header) *httpBackendResponse {
	return &httpBackendResponse{
		config:     b.config,
		statusCode: b.statusCode,
		header:     header,
		body:       b.body,
	}
}

// ModifyBody returns a new instance of httpBackendResponse with the given body modified.
// The method creates a copy of the original httpBackendResponse and sets the body to the provided value.
// The other fields are copied from the original httpBackendResponse.
func (b *httpBackendResponse) ModifyBody(body *Body) *httpBackendResponse {
	return &httpBackendResponse{
		config:     b.config,
		statusCode: b.statusCode,
		header:     b.header,
		body:       body,
	}
}

// Ok returns a boolean indicating if the statusCode of the httpBackendResponse instance is within the range 200-299.
func (b *httpBackendResponse) Ok() bool {
	return helper.IsGreaterThanOrEqual(b.statusCode, 200) && helper.IsLessThanOrEqual(b.statusCode, 299)
}

// Key returns the key of the httpBackendResponse for aggregation.
// The key is composed of the string "backend" and the index, if it is greater than or equal to zero.
// If the httpBackendResponse has a name, the key is set to the name.
func (b *httpBackendResponse) Key(index int) (key string) {
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
func (b *httpBackendResponse) StatusCode() int {
	return b.statusCode
}

// Header returns the `header` of the `httpBackendResponse` instance.
func (b *httpBackendResponse) Header() Header {
	return b.header
}

// Body returns the `body` of the `httpBackendResponse` instance.
func (b *httpBackendResponse) Body() *Body {
	return b.body
}

// GroupByType returns true if the httpBackendResponse instance should be grouped,
// either by setting the groupResponse field to true or by the value of the body being a text or a slice.
// Otherwise, it returns false.
func (b *httpBackendResponse) GroupByType() bool {
	body := b.Body()
	return helper.IsNotNil(body) && body.IsText() || helper.IsSlice(body.Bytes())
}

// Map returns a map containing the evaluated fields of the httpBackendResponse instance.
// The returned map includes the "statusCode" field, which represents the HTTP statusCode of the httpResponse,
// the "header" field, which represents the body fields of the httpResponse, and the "body" field, which represents
// the body of the httpResponse as an interface{} type.
func (b *httpBackendResponse) Map() any {
	var evalBody any
	if helper.IsNotNil(evalBody) {
		evalBody = b.body.Interface()
	}
	return map[string]any{
		"statusCode": b.statusCode,
		"header":     b.header,
		"body":       evalBody,
	}
}

func (b *httpBackendResponse) ApplyConfig(momentToApply enum.BackendResponseApply) *httpBackendResponse {
	// instanciamos o httpResponse config do backend
	backendResponseVO := b.Config()

	// se o backend ja foi aplicado, ou a config for nil, ou o momento não é o ideal
	if b.Applied() || helper.IsNil(backendResponseVO) || helper.IsNotEqualTo(momentToApply, backendResponseVO.Apply()) {
		return b
	}

	// ele quer ser omitido retornamos nil
	if backendResponseVO.Omit() {
		return nil
	}

	// instanciamos o novo header
	var header Header
	// verificamos se o header não quer ser omitido, caso não, preenchemos o mesmo mapeado e projetado segundo o json config
	if !backendResponseVO.OmitHeader() {
		header = b.Header()
		header = header.MapToResponse(backendResponseVO.HeaderMapper())
		header = header.ProjectionToResponse(backendResponseVO.HeaderProjection())
	}

	// instanciamos o novo body
	var body *Body
	// verificamos se o body não quer ser omitido e não for nil, caso isso aconteça, preenchemos o mesmo mapeado e projetado
	if !backendResponseVO.OmitBody() && helper.IsNotNil(b.Body()) {
		body = b.Body()
		body = body.Map(backendResponseVO.BodyMapper())
		body = body.Projection(backendResponseVO.BodyProjection())
		if backendResponseVO.HasGroup() {
			body = newBodyAggregateByKey(backendResponseVO.Group(), body)
		}
	}

	// construímos o novo httpBackendResponse aplicado
	return &httpBackendResponse{
		config:     b.config,
		applied:    true,
		statusCode: b.statusCode,
		header:     header,
		body:       body,
	}
}

func (b *httpBackendResponse) Applied() bool {
	return b.applied
}

func (b *httpBackendResponse) Config() *BackendResponse {
	return b.config
}
