package vo

import "github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"

type modifyBodies struct {
	modify
}

// NewModifyBodies creates a new instance of modifyBodies struct
// with the provided Modifier, Request, and Response.
func NewModifyBodies(modifierVO Modifier, requestVO Request, responseVO Response) ModifierStrategy {
	return modifyBodies{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// Execute applies modifications to the request and/or response bodies
// based on the specified ModifyScope.
//
// For a ModifierScopeRequest, it will execute changes in the context of the request scope.
// For a ModifierScopeResponse, it will execute changes in the context of the response scope.
// If the ModifierScope is not recognized, it will return the original request and response.
//
// Returns:
//   - Request: the modified (or original) HTTP Request
//   - Response: the modified (or original) HTTP Response
func (m modifyBodies) Execute() (Request, Response) {
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

// executeRequestScope executes the body modifications in the scope of a request.
// This function first modifies the local and global bodies of the request.
// It then modifies the request's backend locally and globally.
// The function finally returns the globally modified request and the response.
//
// Returns:
// - Request: The modified request after applying the global modifications.
// - Response: The unmodified response as is.
//
// Note: The nature of modifications are defined within the 'bodies', 'modifyRequestLocal', and
// 'modifyRequestGlobal' methods of the modifyBodies receiver.
func (m modifyBodies) executeRequestScope() (Request, Response) {
	// chamamos o modify de bodies passando os bodies a ser modificado e o mesmo retorna modificados
	globalBody, localBody := m.bodies(m.globalRequestBody(), m.localRequestBody())

	// modificamos o body local
	backendRequestVO := m.modifyRequestLocal(localBody)

	// modificamos o body global e retornamos
	return m.modifyRequestGlobal(globalBody, backendRequestVO), m.response
}

// executeResponseScope modifies the bodies of the global and local responses and returns the modified responses along
// with the request.
// It first gets the global and local body of the response by calling the modify bodies function.
// Then, it modifies the local header and the global header by using the respective modify functions.
// The modified global response body, local response body, and the request is returned as a result.
//
// Returns:
// 1) Request : Original request object
// 2) Response: Modified global and local response body
func (m modifyBodies) executeResponseScope() (Request, Response) {
	// chamamos o modify de bodies passando os bodies a ser modificado e o mesmo retorna modificados
	globalBody, localBody := m.bodies(m.globalResponseBody(), m.localResponseBody())

	// modificamos o header local
	backendResponseVO := m.modifyResponseLocal(localBody)

	// modificamos o header global e retornamos
	return m.request, m.modifyResponseGlobal(globalBody, backendResponseVO)
}

// globalRequestBody returns the global body of the request in the modifyBodies struct.
func (m modifyBodies) globalRequestBody() Body {
	return m.request.Body()
}

// globalRequestBody returns the local body of the current backend request in the Request struct.
func (m modifyBodies) localRequestBody() Body {
	return m.request.CurrentBackendRequest().Body()
}

// globalRequestBody returns the global body of the response in the modifyBodies struct.
func (m modifyBodies) globalResponseBody() Body {
	return m.response.Body()
}

// globalRequestBody returns the local body of the last backend response in the Response struct.
func (m modifyBodies) localResponseBody() Body {
	return m.response.LastBackendResponse().Body()
}

// modifyRequestLocal modifies the request body of the current backend request.
// It takes a localBody of type Body as input and returns the modified backendRequest.
//
// Parameters
// localBody Body: body to replace in the current backend request.
//
// Returns
// backendRequest: the backend request after its body has been modified.
func (m modifyBodies) modifyRequestLocal(localBody Body) backendRequest {
	return m.request.CurrentBackendRequest().ModifyBody(localBody)
}

// modifyRequestGlobal modifies the request globally using the specified `globalBody` and `backendRequestVO`.
// It makes use of the `modifyBodies` receiver's `request` field to perform the modification.
//
// Parameters:
// - `globalBody`: This is the body that is applied to the object value request.
// - `backendRequestVO`: This is the value object(virtual object) that represents the current backend request.
//
// Returns:
// - `Request`: This is the modified request after applying the body to the original request.
func (m modifyBodies) modifyRequestGlobal(globalBody Body, backendRequestVO backendRequest) Request {
	return m.request.ModifyBody(globalBody, backendRequestVO)
}

// modifyResponseLocal modifies the Response Body of a backendResponse.
// It uses a modifyBodies receiver that contains the LastBackendResponse.
// It consumes a Body type 'localBody' which is used to modify the existing backendResponse's body.
// Returns the modified 'backendResponse'.
func (m modifyBodies) modifyResponseLocal(localBody Body) backendResponse {
	return m.response.LastBackendResponse().ModifyBody(localBody)
}

// modifyResponseGlobal is a method of the modifyBodies type. It accepts two parameters:
// globalBody of type Body and backendResponseVO of type backendResponse.
// The method invokes the ModifyBody method of the response field in the modifyBodies receiver,
// passing the globalBody and backendResponseVO as arguments.
// It returns a Response, the outcome of the ModifyBody method.
func (m modifyBodies) modifyResponseGlobal(globalBody Body, backendResponseVO backendResponse) Response {
	return m.response.ModifyBody(globalBody, backendResponseVO)
}
