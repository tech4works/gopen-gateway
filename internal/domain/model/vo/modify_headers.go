package vo

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

type headers struct {
	modify
}

// NewHeaders creates a new instance of modifyHeaders struct
// with the provided Modifier, Request, and Response.
func NewHeaders(modifierVO Modifier, requestVO Request, responseVO Response) ModifierStrategy {
	return headers{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// Execute executes the modifyHeaders functionality.
// It determines the scope and calls the appropriate method to execute the modification.
// If the scope is enum.ModifierScopeRequest, it calls executeRequestScope and returns the modified request and original response.
// If the scope is enum.ModifierScopeResponse, it calls executeResponseScope and returns the original request and modified response.
// If the scope is neither enum.ModifierScopeRequest nor enum.ModifierScopeResponse, it returns the original request and response.
func (m headers) Execute() (Request, Response) {
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

// executeRequestScope modifies both the global and local request headers as well as the propagate header,
// and returns the modified requests and the response.
//
// The function first calls the headers method with the global and local request header as input.
// It then changes the local request header and modifies the propagate header based on the backendRequestVO.
//
// Return values:
//
// - Request: The modified global request, which includes the propagate header.
//
// - Response: The original response from the `headers` struct.
func (m headers) executeRequestScope() (Request, Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	globalHeader, localHeader := m.headers(m.globalRequestHeader(), m.localRequestHeader())

	// modificamos o header local
	backendRequestVO := m.modifyLocalRequest(localHeader)

	// modificamos o header propagate e retornamos
	return m.modifyGlobalRequest(globalHeader, backendRequestVO), m.response
}

// executeResponseScope executes the modifyHeaders functionality for the response scope.
// It calls the headers method to get the global and local response headers.
// It then modifies the local header using the modifyLocalResponse method.
// Finally, it modifies the response using the modifyGlobalResponse method with the modified backendResponseVO.
// It returns the original request and the modified response.
func (m headers) executeResponseScope() (Request, Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	_, localHeader := m.headers(m.globalResponseHeader(), m.localResponseHeader())

	// modificamos o header local
	backendResponseVO := m.modifyLocalResponse(localHeader)

	// modificamos a resposta vo com o novo backendResponseVO
	return m.request, m.modifyGlobalResponse(backendResponseVO)
}

// globalRequestHeader returns the header of the request in the headers struct.
func (m headers) globalRequestHeader() Header {
	return m.request.Header()
}

// localRequestHeader returns the header of the current backend request in the headers struct.
func (m headers) localRequestHeader() Header {
	return m.request.CurrentBackendRequest().Header()
}

// globalResponseHeader returns the header of the response in the headers struct.
func (m headers) globalResponseHeader() Header {
	return m.response.Header()
}

// localResponseHeader returns the header of the last backend response in the headers struct.
func (m headers) localResponseHeader() Header {
	return m.response.LastBackendResponse().Header()
}

// modifyLocalRequest modifies the local request header of the backend request by applying the provided local header.
// It creates a new instance of the backendRequest struct with the modified header and returns it.
func (m headers) modifyLocalRequest(localHeader Header) backendRequest {
	return m.request.CurrentBackendRequest().ModifyHeader(localHeader)
}

// modifyGlobalRequest modifies the propagate request by modifying the propagate header and the backend request.
// It returns a new modified Request.
func (m headers) modifyGlobalRequest(globalHeader Header, backendRequestVO backendRequest) Request {
	return m.request.ModifyHeader(globalHeader, backendRequestVO)
}

// modifyLocalResponse modifies the local header of the response by applying changes from the given localHeader.
// It returns a new backendResponse object with the modified header.
func (m headers) modifyLocalResponse(localHeader Header) backendResponse {
	return m.response.LastBackendResponse().ModifyHeader(localHeader)
}

// modifyGlobalResponse modifies the global response by applying the modifications specified in the backendResponseVO.
// It calls the ModifyLastBackendResponse method of the response with the backendResponseVO as the argument.
// The modified response is then returned.
func (m headers) modifyGlobalResponse(backendResponseVO backendResponse) Response {
	return m.response.ModifyLastBackendResponse(backendResponseVO)
}
