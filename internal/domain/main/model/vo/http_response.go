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
	configVO "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/enum"
	"net/http"
)

// HttpResponse represents the gateway HTTP httpResponse.
type HttpResponse struct {
	// statusCode stores the integer HTTP status code of the HttpResponse object.
	statusCode int
	// header represents the header of the HttpResponse object.
	header Header
	// Body represents the body of the gateway HTTP httpResponse.
	body *Body
	// todo:
	abort bool
	// todo:
	written bool
	// history represents the history of backend responses in the HttpResponse object.
	history httpResponseHistory
}

// NewHttpResponse creates a new HttpResponse object with the given Endpoint.
// The StatusCode is set to http.StatusNoContent.
// Returns the newly created HttpResponse object.
func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		statusCode: http.StatusNoContent,
	}
}

func NewHttpResponseAborted(endpoint *configVO.Endpoint, httpBackendResponse *HttpBackendResponse) *HttpResponse {
	// construímos o httpResponse com os dados do backend abortado
	header := NewResponseHeader(endpoint.Completed(1), httpBackendResponse.Ok())
	header = header.Aggregate(httpBackendResponse.Header())
	return &HttpResponse{
		statusCode: httpBackendResponse.StatusCode(),
		header:     header,
		body:       httpBackendResponse.Body(),
		abort:      true,
	}
}

func NewHttpResponseByStatusCode(statusCode int) *HttpResponse {
	ok := helper.IsGreaterThanOrEqual(statusCode, 200) || helper.IsLessThanOrEqual(statusCode, 299)
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewResponseHeader(true, ok),
	}
}

func NewHttpResponseByString(statusCode int, body string) *HttpResponse {
	ok := helper.IsGreaterThanOrEqual(statusCode, 200) || helper.IsLessThanOrEqual(statusCode, 299)
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewResponseHeader(true, ok),
		body:       NewBodyByString(body),
	}
}

func NewHttpResponseByJson(statusCode int, body any) *HttpResponse {
	ok := helper.IsGreaterThanOrEqual(statusCode, 200) || helper.IsLessThanOrEqual(statusCode, 299)
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewResponseHeader(true, ok),
		body:       NewBodyByJson(body),
	}
}

func NewHttpResponseByCache(cacheResponseVO *CacheResponse) *HttpResponse {
	header := cacheResponseVO.Header
	header = header.Set(consts.XGopenCache, helper.SimpleConvertToString(true))
	header = header.Set(consts.XGopenCacheTTL, cacheResponseVO.TTL())
	return &HttpResponse{
		statusCode: cacheResponseVO.StatusCode,
		header:     header,
		body:       NewBodyByCache(cacheResponseVO.Body),
	}
}

func NewHttpResponseByErr(path string, statusCode int, err error) *HttpResponse {
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewHeaderFailed(),
		body:       NewBodyByError(path, err),
		abort:      true,
	}
}

func (r *HttpResponse) Append(httpBackendResponse *HttpBackendResponse) *HttpResponse {
	// se for nil quer dizer que ele quis ser omitido, então nem damos o append
	if helper.IsNil(httpBackendResponse) {
		return r
	}

	// adicionamos na nova lista de histórico
	history := r.history
	history = append(history, httpBackendResponse)
	return &HttpResponse{
		statusCode: r.StatusCode(),
		header:     r.Header(),
		body:       r.Body(),
		history:    history,
	}
}

// Error constructs a standard gateway error httpResponse based on the received error.
// It builds the httpResponse's status code from the received error.
// If the error contains mapper.ErrBadGateway, the status code is set to http.StatusBadGateway.
// If the error contains mapper.ErrGatewayTimeout, the status code is set to http.StatusGatewayTimeout.
// Otherwise, the status code is set to http.StatusInternalServerError.
// It constructs the default gateway error httpResponse by setting the status code, header, body, and abort properties.
// Returns the constructed HttpResponse object representing the gateway error httpResponse.
func (r *HttpResponse) Error(path string, err error) *HttpResponse {
	// construímos o statusCode de resposta a partir do erro recebido
	var statusCode int
	if errors.Contains(err, mapper.ErrBadGateway) {
		statusCode = http.StatusBadGateway
	} else if errors.Contains(err, mapper.ErrGatewayTimeout) {
		statusCode = http.StatusGatewayTimeout
	} else {
		statusCode = http.StatusInternalServerError
	}
	// construímos a resposta de erro padrão do gateway
	return &HttpResponse{
		statusCode: statusCode,
		header:     NewHeaderFailed(),
		body:       NewBodyByError(path, err),
		abort:      true,
	}
}

// Abort returns the value of the `abort` property of the HttpResponse object.
// If `abort` is true, it indicates that the httpResponse should be aborted.
// Returns a boolean value representing the `abort` property.
func (r *HttpResponse) Abort() bool {
	return r.abort
}

// Written returns a boolean value indicating whether the HttpResponse has been written.
// Returns true if the HttpResponse has been written, false otherwise.
func (r *HttpResponse) Written() bool {
	return r.written
}

// StatusCode returns the status code of the HttpResponse object.
func (r *HttpResponse) StatusCode() int {
	return r.statusCode
}

// Header returns the header of the HttpResponse object.
func (r *HttpResponse) Header() Header {
	return r.header
}

func (r *HttpResponse) ContentType() enum.ContentType {
	if helper.IsNotNil(r.Body()) {
		return r.Body().ContentType()
	}
	return ""
}

// Body returns the body of the HttpResponse object.
func (r *HttpResponse) Body() *Body {
	return r.body
}

// todo:
func (r *HttpResponse) BytesBody() []byte {
	// se o body for nil retornamos nil
	if helper.IsNil(r.body) {
		return nil
	}
	return r.Body().Bytes()
}

func (r *HttpResponse) History() httpResponseHistory {
	return r.history
}

// todo:
func (r *HttpResponse) Write(endpoint *configVO.Endpoint, httpRequest *HttpRequest, httpResponse *HttpResponse,
) *HttpResponse {
	// se resposta ja foi escrita retornamos a mesma
	if r.Written() {
		return r
	}

	// instanciamos o novo statusCode com o que ja existe
	statusCode := r.StatusCode()
	// instanciamos o novo header com o que ja existe
	header := r.Header()
	// instanciamos o novo body com o que ja existe
	body := r.Body()

	// verificamos se ele tem um histórico, caso tenha, iremos preencher esses dados com o histórico
	if r.HasHistory() {
		statusCode, header, body = r.writeByHistory(endpoint, httpRequest, httpResponse)
	}

	// escrevemos a resposta com base nas configurações do endpoint
	return r.writeByEndpointConfig(endpoint, statusCode, header, body)
}

// Map returns the evaluation result of the history list of the HttpResponse object.
func (r *HttpResponse) Map() string {
	return r.History().Map()
}

func (r *HttpResponse) HasHistory() bool {
	return helper.IsNotEmpty(r.History())
}

// todo:
func (r *HttpResponse) writeByHistory(endpoint *configVO.Endpoint, httpRequest *HttpRequest, httpResponse *HttpResponse) (
	statusCode int, header Header, body *Body) {
	// instanciamos a configuração de resposta do endpoint
	endpointResponse := endpoint.Response()

	// filtramos o histórico
	filteredHistory := r.History().Filter(httpRequest, httpResponse)

	// obtemos o status code a partir do histórico filtrado
	statusCode = filteredHistory.StatusCode()
	// criamos o header a partir dos valores complete e success
	header = NewResponseHeader(endpoint.Completed(filteredHistory.Size()), filteredHistory.Success())
	// agregamos os headers do histórico filtrado
	header = header.Aggregate(filteredHistory.Header())

	// verificamos se precisa agregar o body
	aggregate := false
	if helper.IsNotNil(endpointResponse) {
		aggregate = endpointResponse.Aggregate()
	}
	// obtemos o body a partir do histórico filtrado
	body = filteredHistory.Body(aggregate)

	// retornamos os valores preenchidos
	return statusCode, header, body
}

func (r *HttpResponse) writeByEndpointConfig(endpoint *configVO.Endpoint, statusCode int, header Header,
	body *Body) *HttpResponse {
	// escrevemos o body com base na config do endpoint
	body = r.writeBodyByEndpointConfig(endpoint, body)
	// montamos a resposta com base na configuração do endpoint
	return &HttpResponse{
		statusCode: statusCode,
		header:     header,
		body:       body,
		written:    true,
	}
}

func (r *HttpResponse) writeBodyByEndpointConfig(endpoint *configVO.Endpoint, body *Body) *Body {
	// instanciamos a config de resposta do endpoint
	endpointResponse := endpoint.Response()

	// se a config for nil ou o body ignoramos e retornamos o mesmo
	if helper.IsNil(endpointResponse) || helper.IsNil(body) {
		return body
	}

	// se omitEmpty for true omitimos os campos vazios ou nulos
	if endpointResponse.OmitEmpty() {
		body = body.OmitEmpty()
	}
	// se ele tem a nomenclatura desejada, fazemos a conversão
	if endpointResponse.HasNomenclature() {
		body = body.ToCase(endpointResponse.Nomenclature())
	}

	// obtemos o content-type para preparar o body pela content-type desejada
	var contentType enum.ContentType
	if helper.IsNotNil(endpoint.Response()) && endpoint.Response().HasEncode() {
		contentType = endpoint.Response().Encode().ContentType()
	} else if helper.IsNotNil(body) {
		contentType = body.ContentType()
	}

	// se content-type deseja for diferente do que ja temos, então criamos um novo body com o encode desejado
	if helper.IsNotEqualTo(contentType, body.ContentType()) {
		bs := body.BytesByContentType(contentType)
		body = NewBodyByContentType(contentType.String(), helper.SimpleConvertToBuffer(bs))
	}

	// retornamos o body escrito e possívelmente modificado
	return body
}
