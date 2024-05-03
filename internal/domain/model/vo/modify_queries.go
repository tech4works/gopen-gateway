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

type modifyQueries struct {
	modify
}

func NewModifyQueries(modifierVO *Modifier, httpRequestVO *HttpRequest, httpResponseVO *HttpResponse) ModifierStrategy {
	return modifyQueries{
		modify: newModify(modifierVO, httpRequestVO, httpResponseVO),
	}
}

func (m modifyQueries) Execute() (*HttpRequest, *HttpResponse) {
	// executamos a partir do escopo padrão
	return m.executeRequestScope()
}

func (m modifyQueries) executeRequestScope() (*HttpRequest, *HttpResponse) {
	// chamamos o modify de queries passando as queries a ser modificado e o mesmo retorna os modificados
	httpRequestQuery, httpBackendRequestQuery := m.queries(m.httpRequestQuery(), m.httpBackendRequestQuery())

	// modificamos o query local
	httpBackendRequestVO := m.modifyHttpBackendRequest(httpBackendRequestQuery)

	// modificamos o params propagate e retornamos
	return m.modifyHttpRequest(httpRequestQuery, httpBackendRequestVO), m.httpResponse
}

func (m modifyQueries) httpRequestQuery() Query {
	return m.httpRequest.Query()
}

func (m modifyQueries) httpBackendRequestQuery() Query {
	return m.httpRequest.LastHttpBackendRequest().Query()
}

func (m modifyQueries) modifyHttpBackendRequest(queryVO Query) *httpBackendRequest {
	return m.httpRequest.LastHttpBackendRequest().ModifyQuery(queryVO)
}

func (m modifyQueries) modifyHttpRequest(queryVO Query, httpBackendRequestVO *httpBackendRequest) *HttpRequest {
	return m.httpRequest.ModifyQuery(queryVO, httpBackendRequestVO)
}
