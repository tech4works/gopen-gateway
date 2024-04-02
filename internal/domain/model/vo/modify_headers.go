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

// executeRequestScope executes the request scope of the headers method.
// It calls the headers method, passing the global and local request headers, and returns the modified headers.
// Then, it modifies the local header by calling the modifyRequestLocal method.
// Finally, it modifies the global header by calling the modifyRequestGlobal method and returns the modified global
// header and the response.
// Returns:
// - Request: Modified global header and backend request object.
// - Response: The same response object.
func (m headers) executeRequestScope() (Request, Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	globalHeader, localHeader := m.headers(m.globalRequestHeader(), m.localRequestHeader())

	// modificamos o header local
	backendRequestVO := m.modifyRequestLocal(localHeader)

	// modificamos o header global e retornamos
	return m.modifyRequestGlobal(globalHeader, backendRequestVO), m.response
}

// executeResponseScope executes the response scope of the headers method.
// It calls the headers method, passing the global and local response headers, and returns the modified headers.
// Then, it modifies the local header by calling the modifyResponseLocal method.
// Finally, it modifies the global header by calling the modifyResponseGlobal method and returns the same request and the
// modified global header response.
// Returns:
// - Request: The same request object.
// - Response: Modified global header and backend response object.
func (m headers) executeResponseScope() (Request, Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	globalHeader, localHeader := m.headers(m.globalResponseHeader(), m.localResponseHeader())

	// modificamos o header local
	backendResponseVO := m.modifyResponseLocal(localHeader)

	// modificamos o header global e retornamos
	return m.request, m.modifyResponseGlobal(globalHeader, backendResponseVO)
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

// modifyRequestLocal modifies the local request header of the backend request by applying the provided local header.
// It creates a new instance of the backendRequest struct with the modified header and returns it.
func (m headers) modifyRequestLocal(localHeader Header) backendRequest {
	return m.request.CurrentBackendRequest().ModifyHeader(localHeader)
}

// modifyRequestGlobal modifies the global request by modifying the global header and the backend request.
// It returns a new modified Request.
func (m headers) modifyRequestGlobal(globalHeader Header, backendRequestVO backendRequest) Request {
	return m.request.ModifyHeader(globalHeader, backendRequestVO)
}

// modifyResponseLocal modifies the local header of the response by applying changes from the given localHeader.
// It returns a new backendResponse object with the modified header.
func (m headers) modifyResponseLocal(localHeader Header) backendResponse {
	return m.response.LastBackendResponse().ModifyHeader(localHeader)
}

// modifyResponseGlobal modifies the global response by applying the given global header and backend response.
// It returns the modified response.
func (m headers) modifyResponseGlobal(globalHeader Header, backendResponseVO backendResponse) Response {
	return m.response.ModifyHeader(globalHeader, backendResponseVO)
}
