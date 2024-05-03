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

type modifyParam struct {
	modify
}

func NewModifyParam(modifierVO *Modifier, httpRequestVO *HttpRequest, httResponseVO *HttpResponse) ModifierStrategy {
	return modifyParam{
		modify: newModify(modifierVO, httpRequestVO, httResponseVO),
	}
}

func (m modifyParam) Execute() (*HttpRequest, *HttpResponse) {
	// executamos a partir do escopo padrão
	return m.executeRequestScope()
}

func (m modifyParam) executeRequestScope() (*HttpRequest, *HttpResponse) {
	// chamamos o modify de params passando o path e params a ser modificado e o mesmo retorna os mesmo modificados
	httpBackendRequestUrlPathVO, httpRequestParamsVO, httpBackendRequestParamsVO := m.params(
		m.httpBackendRequestUrlPath(), m.httpRequestParams(), m.httpBackendRequestParams())

	// modificamos o params local
	httpBackendRequestVO := m.modifyHttpBackendRequest(httpBackendRequestUrlPathVO, httpBackendRequestParamsVO)

	// modificamos o params propagate e retornamos
	return m.modifyHttpRequest(httpRequestParamsVO, httpBackendRequestVO), m.httpResponse
}

func (m modifyParam) httpBackendRequestUrlPath() UrlPath {
	return m.httpRequest.LastHttpBackendRequest().Path()
}

func (m modifyParam) httpRequestParams() Params {
	return m.httpRequest.Params()
}

func (m modifyParam) httpBackendRequestParams() Params {
	return m.httpRequest.LastHttpBackendRequest().Params()
}

func (m modifyParam) modifyHttpBackendRequest(urlPathVO UrlPath, paramsVO Params) *httpBackendRequest {
	return m.httpRequest.LastHttpBackendRequest().ModifyParams(urlPathVO, paramsVO)
}

func (m modifyParam) modifyHttpRequest(paramsVO Params, httpBackendRequestVO *httpBackendRequest) *HttpRequest {
	return m.httpRequest.ModifyParams(paramsVO, httpBackendRequestVO)
}
