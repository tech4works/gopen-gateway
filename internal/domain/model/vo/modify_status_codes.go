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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

type modifyStatusCodes struct {
	modify
}

func NewModifyStatusCodes(statusCodeValue int, httpRequestVO *HttpRequest, httpResponseVO *HttpResponse) ModifierStrategy {
	statusCodeStr := helper.SimpleConvertToString(statusCodeValue)
	return modifyStatusCodes{
		modify: newModify(newModifierFromValue(enum.ModifierContextResponse, statusCodeStr), httpRequestVO, httpResponseVO),
	}
}

func (m modifyStatusCodes) Execute() (*HttpRequest, *HttpResponse) {
	// executamos a partir do escopo padrão
	return m.executeResponseScope()
}

func (m modifyStatusCodes) executeResponseScope() (*HttpRequest, *HttpResponse) {
	// se tiver com o histórico vazio, retornamos os mesmos
	if !m.httpResponse.HasHistory() {
		return m.httpRequest, m.httpResponse
	}

	// chamamos o modify de status code passando os status codes a ser modificado e o mesmo retorna modificados
	statusCode := m.statusCode(m.httpResponse.LastHttpBackendResponse().StatusCode())

	// modificamos o status code local
	httpBackendResponseVO := m.modifyHttpBackendResponse(statusCode)

	// modificamos o httpResponse global com o novo backendResponseVO
	return m.httpRequest, m.modifyHttpResponse(httpBackendResponseVO)
}

func (m modifyStatusCodes) modifyHttpBackendResponse(statusCode int) *httpBackendResponse {
	return m.httpResponse.LastHttpBackendResponse().ModifyStatusCode(statusCode)
}

func (m modifyStatusCodes) modifyHttpResponse(httpBackendResponseVO *httpBackendResponse) *HttpResponse {
	return m.httpResponse.ModifyLastHttpBackendResponse(httpBackendResponseVO)
}
