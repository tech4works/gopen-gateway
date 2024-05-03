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

type modifyBodies struct {
	modify
}

func NewModifyBodies(modifierVO *Modifier, httpRequestVO *HttpRequest, httpResponseVO *HttpResponse) ModifierStrategy {
	return &modifyBodies{
		modify: newModify(modifierVO, httpRequestVO, httpResponseVO),
	}
}

func (m modifyBodies) Execute() (*HttpRequest, *HttpResponse) {
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

func (m modifyBodies) executeRequestScope() (*HttpRequest, *HttpResponse) {
	// chamamos o modify de bodies passando os bodies a ser modificado e o mesmo retorna modificados
	httpBodyVO, httpBackendBodyVO := m.bodies(m.httpRequestBody(), m.httpBackendRequestBody())

	// modificamos o body local
	httpBackendRequestVO := m.modifyHttpRequestBackend(httpBackendBodyVO)

	// modificamos o body propagate e retornamos
	return m.modifyHttpRequest(httpBodyVO, httpBackendRequestVO), m.httpResponse
}

func (m modifyBodies) executeResponseScope() (*HttpRequest, *HttpResponse) {
	// se tiver com o histórico vazio, retornamos os mesmos
	if !m.httpResponse.HasHistory() {
		return m.httpRequest, m.httpResponse
	}

	// chamamos o modify de bodies passando os bodies a ser modificado e o mesmo retorna modificados
	_, httpBackendBodyVO := m.bodies(m.httpResponseBody(), m.httpBackendResponseBody())

	// modificamos o header local
	httpBackendResponseVO := m.modifyHttpBackendResponse(httpBackendBodyVO)

	// modificamos a resposta vo com o novo httpBackendResponseVO
	return m.httpRequest, m.modifyHttpResponse(httpBackendResponseVO)
}

func (m modifyBodies) httpRequestBody() *Body {
	return m.httpRequest.Body()
}

func (m modifyBodies) httpBackendRequestBody() *Body {
	return m.httpRequest.LastHttpBackendRequest().Body()
}

func (m modifyBodies) httpResponseBody() *Body {
	return m.httpResponse.Body()
}

func (m modifyBodies) httpBackendResponseBody() *Body {
	return m.httpResponse.LastHttpBackendResponse().Body()
}

func (m modifyBodies) modifyHttpRequestBackend(bodyVO *Body) *httpBackendRequest {
	return m.httpRequest.LastHttpBackendRequest().ModifyBody(bodyVO)
}

func (m modifyBodies) modifyHttpRequest(bodyVO *Body, httpBackendRequestVO *httpBackendRequest) *HttpRequest {
	return m.httpRequest.ModifyBody(bodyVO, httpBackendRequestVO)
}

func (m modifyBodies) modifyHttpBackendResponse(bodyVO *Body) *httpBackendResponse {
	return m.httpResponse.LastHttpBackendResponse().ModifyBody(bodyVO)
}

func (m modifyBodies) modifyHttpResponse(httpBackendResponseVO *httpBackendResponse) *HttpResponse {
	return m.httpResponse.ModifyLastHttpBackendResponse(httpBackendResponseVO)
}
