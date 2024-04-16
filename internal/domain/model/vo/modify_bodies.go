package vo

import "github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"

// modifyBodies is a type that represents the combination of a modify struct
// and additional methods or properties defined within it. This type allows for
// modifying bodies of an object by utilizing the methods and properties of the
// modify struct it contains.
type modifyBodies struct {
	modify
}

// NewModifyBodies creates a new instance of modifyBodies strategy that applies modifications
// to the request and/or response bodies based on the specified modifierVO, requestVO, and responseVO.
// Parameters:
// - modifierVO: a pointer to the Modifier object that contains the modification details.
// - requestVO: a pointer to the Request object that represents the HTTP Request.
// - responseVO: a pointer to the Response object that represents the HTTP Response.
// Returns:
// - ModifierStrategy: a new instance of modifyBodies strategy.
func NewModifyBodies(modifierVO *Modifier, requestVO *Request, responseVO *Response) ModifierStrategy {
	return &modifyBodies{
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
func (m modifyBodies) Execute() (*Request, *Response) {
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

func (m modifyBodies) executeRequestScope() (*Request, *Response) {
	// chamamos o modify de bodies passando os bodies a ser modificado e o mesmo retorna modificados
	globalBody, localBody := m.bodies(m.globalRequestBody(), m.localRequestBody())

	// modificamos o body local
	backendRequestVO := m.modifyLocalRequest(localBody)

	// modificamos o body propagate e retornamos
	return m.modifyGlobalRequest(globalBody, backendRequestVO), m.response
}

// executeResponseScope modifies both the global response body and local response body.
// It first calls the bodies method of modifyBodies, passing the bodies to be modified,
// which in turn returns the modified bodies. The local response body is then modified.
// Finally, the method modifies the ResponseGlobal with the new backendResponseVO.
//
// Returns:
//
// Request: The original request.
//
// Response: The modified response.
func (m modifyBodies) executeResponseScope() (*Request, *Response) {
	// chamamos o modify de bodies passando os bodies a ser modificado e o mesmo retorna modificados
	_, localBody := m.bodies(m.globalResponseBody(), m.localResponseBody())

	// modificamos o header local
	backendResponseVO := m.modifyLocalResponse(localBody)

	// modificamos a resposta vo com o novo backendResponseVO
	return m.request, m.modifyGlobalResponse(backendResponseVO)
}

// globalRequestBody returns the propagate body of the request in the modifyBodies struct.
func (m modifyBodies) globalRequestBody() *Body {
	return m.request.Body()
}

// globalRequestBody returns the local body of the current backend request in the Request struct.
func (m modifyBodies) localRequestBody() *Body {
	return m.request.CurrentBackendRequest().Body()
}

// globalRequestBody returns the propagate body of the response in the modifyBodies struct.
func (m modifyBodies) globalResponseBody() *Body {
	return m.response.Body()
}

// globalRequestBody returns the local body of the last backend response in the Response struct.
func (m modifyBodies) localResponseBody() *Body {
	return m.response.LastBackendResponse().Body()
}

// modifyLocalRequest modifies the request body of the current backend request.
// It takes a localBody of type Body as input and returns the modified backendRequest.
//
// Parameters
// localBody Body: body to replace in the current backend request.
//
// Returns
// backendRequest: the backend request after its body has been modified.
func (m modifyBodies) modifyLocalRequest(localBody *Body) *backendRequest {
	return m.request.CurrentBackendRequest().ModifyBody(localBody)
}

// modifyGlobalRequest modifies the request globally using the specified `globalBody` and `backendRequestVO`.
// It makes use of the `modifyBodies` receiver's `request` field to perform the modification.
//
// Parameters:
// - `globalBody`: This is the body that is applied to the object value request.
// - `backendRequestVO`: This is the value object(virtual object) that represents the current backend request.
//
// Returns:
// - `Request`: This is the modified request after applying the body to the original request.
func (m modifyBodies) modifyGlobalRequest(globalBody *Body, backendRequestVO *backendRequest) *Request {
	return m.request.ModifyBody(globalBody, backendRequestVO)
}

// modifyLocalResponse modifies the Response Body of a backendResponse.
// It uses a modifyBodies receiver that contains the LastBackendResponse.
// It consumes a Body type 'localBody' which is used to modify the existing backendResponse's body.
// Returns the modified 'backendResponse'.
func (m modifyBodies) modifyLocalResponse(localBody *Body) *backendResponse {
	return m.response.LastBackendResponse().ModifyBody(localBody)
}

// modifyGlobalResponse modifies the last backend response of the response body
// by replacing it with the provided backend response object.
//
// Parameters:
// - backendResponseVO: the backend response object to replace the last backend response
//
// Returns:
// - Response: the response with the modified last backend response
func (m modifyBodies) modifyGlobalResponse(backendResponseVO *backendResponse) *Response {
	return m.response.ModifyLastBackendResponse(backendResponseVO)
}
