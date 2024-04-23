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

// modifyHeaders represents a type that allows modifying HTTP headers in a request or response.
// It is used to customize the headers for a request or response in an HTTP client or server.
// This type has a `modify` embedded field, which provides a set of methods to modify headers.
type modifyHeaders struct {
	modify
}

// NewHeaders creates a new instance of modifyHeaders struct
// with the provided Modifier, Request, and Response.
func NewHeaders(modifierVO *Modifier, requestVO *Request, responseVO *Response) ModifierStrategy {
	return modifyHeaders{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// Execute executes the modifyHeaders functionality.
// It determines the scope and calls the appropriate method to execute the modification.
// If the scope is enum.ModifierScopeRequest, it calls executeRequestScope and returns the modified request and original response.
// If the scope is enum.ModifierScopeResponse, it calls executeResponseScope and returns the original request and modified response.
// If the scope is neither enum.ModifierScopeRequest nor enum.ModifierScopeResponse, it returns the original request and response.
func (m modifyHeaders) Execute() (*Request, *Response) {
	// executamos a partir do escopo configurado
	switch m.scope {
	case enum.ModifierScopeRequest:
		return m.executeRequestScope()
	case enum.ModifierScopeResponse:
		return m.executeResponseScope()
	default:
		return m.request, m.response
	}
}

// executeRequestScope modifies both the global and local request modifyHeaders as well as the propagate header,
// and returns the modified requests and the response.
//
// The function first calls the modifyHeaders method with the global and local request header as input.
// It then changes the local request header and modifies the propagate header based on the backendRequestVO.
//
// Return values:
//
// - Request: The modified global request, which includes the propagate header.
//
// - Response: The original response from the `headers` struct.
func (m modifyHeaders) executeRequestScope() (*Request, *Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	globalHeader, localHeader := m.headers(m.globalRequestHeader(), m.localRequestHeader())

	// modificamos o header local
	backendRequestVO := m.modifyLocalRequest(localHeader)

	// modificamos o header propagate e retornamos
	return m.modifyGlobalRequest(globalHeader, backendRequestVO), m.response
}

// executeResponseScope executes the modifyHeaders functionality for the response scope.
// It calls the modifyHeaders method to get the global and local response modifyHeaders.
// It then modifies the local header using the modifyLocalResponse method.
// Finally, it modifies the response using the modifyGlobalResponse method with the modified backendResponseVO.
// It returns the original request and the modified response.
func (m modifyHeaders) executeResponseScope() (*Request, *Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	_, localHeader := m.headers(m.globalResponseHeader(), m.localResponseHeader())

	// modificamos o header local
	backendResponseVO := m.modifyLocalResponse(localHeader)

	// modificamos a resposta vo com o novo backendResponseVO
	return m.request, m.modifyGlobalResponse(backendResponseVO)
}

// globalRequestHeader returns the header of the request in the modifyHeaders struct.
func (m modifyHeaders) globalRequestHeader() Header {
	return m.request.Header()
}

// localRequestHeader returns the header of the current backend request in the modifyHeaders struct.
func (m modifyHeaders) localRequestHeader() Header {
	return m.request.CurrentBackendRequest().Header()
}

// globalResponseHeader returns the header of the response in the modifyHeaders struct.
func (m modifyHeaders) globalResponseHeader() Header {
	return m.response.Header()
}

// localResponseHeader returns the header of the last backend response in the modifyHeaders struct.
func (m modifyHeaders) localResponseHeader() Header {
	return m.response.LastBackendResponse().Header()
}

// modifyLocalRequest modifies the local request header of the backend request by applying the provided local header.
// It creates a new instance of the backendRequest struct with the modified header and returns it.
func (m modifyHeaders) modifyLocalRequest(localHeader Header) *backendRequest {
	return m.request.CurrentBackendRequest().ModifyHeader(localHeader)
}

// modifyGlobalRequest modifies the propagate request by modifying the propagate header and the backend request.
// It returns a new modified Request.
func (m modifyHeaders) modifyGlobalRequest(globalHeader Header, backendRequestVO *backendRequest) *Request {
	return m.request.ModifyHeader(globalHeader, backendRequestVO)
}

// modifyLocalResponse modifies the local header of the response by applying changes from the given localHeader.
// It returns a new backendResponse object with the modified header.
func (m modifyHeaders) modifyLocalResponse(localHeader Header) *backendResponse {
	return m.response.LastBackendResponse().ModifyHeader(localHeader)
}

// modifyGlobalResponse modifies the global response by applying the modifications specified in the backendResponseVO.
// It calls the ModifyLastBackendResponse method of the response with the backendResponseVO as the argument.
// The modified response is then returned.
func (m modifyHeaders) modifyGlobalResponse(backendResponseVO *backendResponse) *Response {
	return m.response.ModifyLastBackendResponse(backendResponseVO)
}
