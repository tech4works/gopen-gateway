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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

type modifyHeaders struct {
	modify
}

func NewHeaders(modifierVO *Modifier, httpRequestVO *HttpRequest, httpResponseVO *HttpResponse) ModifierStrategy {
	return modifyHeaders{
		modify: newModify(modifierVO, httpRequestVO, httpResponseVO),
	}
}

func (m modifyHeaders) Execute() (*HttpRequest, *HttpResponse) {
	// executamos a partir do escopo configurado
	switch m.scope {
	case enum.ModifierScopeRequest:
		return m.executeRequestScope()
	case enum.ModifierScopeResponse:
		return m.executeResponseScope()
	default:
		return m.httpRequest, m.httpResponse
	}
}

func (m modifyHeaders) executeRequestScope() (*HttpRequest, *HttpResponse) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	httpRequestHeaderVO, httpBackendRequestHeaderVO := m.headers(m.httpRequestHeader(), m.httpBackendRequestHeader())

	// modificamos o header local
	httpBackendRequestVO := m.modifyHttpBackendRequest(httpBackendRequestHeaderVO)

	// modificamos o header propagate e retornamos
	return m.modifyHttpRequest(httpRequestHeaderVO, httpBackendRequestVO), m.httpResponse
}

func (m modifyHeaders) executeResponseScope() (*HttpRequest, *HttpResponse) {
	// se tiver com o histórico vazio, retornamos os mesmos
	if !m.httpResponse.HasHistory() {
		return m.httpRequest, m.httpResponse
	}

	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	_, httpBackendHeaderVO := m.headers(m.httpResponseHeader(), m.httpBackendResponseHeader())

	// modificamos o header local
	httpBackendResponseVO := m.modifyHttpBackendResponse(httpBackendHeaderVO)

	// modificamos a resposta vo com o novo httpBackendResponseVO
	return m.httpRequest, m.modifyHttpResponse(httpBackendResponseVO)
}

func (m modifyHeaders) httpRequestHeader() Header {
	return m.httpRequest.Header()
}

func (m modifyHeaders) httpBackendRequestHeader() Header {
	return m.httpRequest.LastHttpBackendRequest().Header()
}

func (m modifyHeaders) httpResponseHeader() Header {
	return m.httpResponse.Header()
}

func (m modifyHeaders) httpBackendResponseHeader() Header {
	return m.httpResponse.LastHttpBackendResponse().Header()
}

func (m modifyHeaders) modifyHttpRequest(headerVO Header, httpBackendRequestVO *httpBackendRequest) *HttpRequest {
	return m.httpRequest.ModifyHeader(headerVO, httpBackendRequestVO)
}

func (m modifyHeaders) modifyHttpBackendRequest(headerVO Header) *httpBackendRequest {
	return m.httpRequest.LastHttpBackendRequest().ModifyHeader(headerVO)
}

func (m modifyHeaders) modifyHttpBackendResponse(headerVO Header) *httpBackendResponse {
	return m.httpResponse.LastHttpBackendResponse().ModifyHeader(headerVO)
}

func (m modifyHeaders) modifyHttpResponse(httpBackendResponseVO *httpBackendResponse) *HttpResponse {
	return m.httpResponse.ModifyLastHttpBackendResponse(httpBackendResponseVO)
}
