package vo

type modifyStatusCodes struct {
	modify
}

// NewModifyStatusCodes creates a new instance of modifyStatusCodes struct
// with the provided Modifier, Request, and Response.
func NewModifyStatusCodes(modifierVO Modifier, requestVO Request, responseVO Response) ModifierStrategy {
	return modifyStatusCodes{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// Execute executes the modifyStatusCodes by calling the executeResponseScope method and returns the modified Request
// and Response.
// It starts execution from the default scope.
func (m modifyStatusCodes) Execute() (Request, Response) {
	// executamos a partir do escopo padrão
	return m.executeResponseScope()
}

// executeResponseScope executes the response scope by modifying the status codes and returning the modified request and response.
// It first calls the statusCodes method of the modifyStatusCodes instance to get the modified global and local status codes.
// Then it calls the modifyResponseLocal method to modify the local status code and returns the modified backendResponse.
// Finally, it calls the modifyResponseGlobal method to modify the global status code and returns the modified response
// with the modified backend response.
// The modified request and response are returned as a tuple.
// Note: This method is called within the Execute method of the modifyStatusCodes instance.
func (m modifyStatusCodes) executeResponseScope() (Request, Response) {
	// chamamos o modify de status code passando os status codes a ser modificado e o mesmo retorna modificados
	globalStatusCode, localStatusCode := m.statusCodes(m.globalResponseStatusCode(), m.localResponseStatusCode())

	// modificamos o status code local
	backendResponseVO := m.modifyResponseLocal(localStatusCode)

	// modificamos o status code global e retornamos
	return m.request, m.modifyResponseGlobal(globalStatusCode, backendResponseVO)
}

// globalResponseStatusCode returns the status code of the response.
func (m modifyStatusCodes) globalResponseStatusCode() int {
	return m.response.StatusCode()
}

// localResponseStatusCode returns the status code of the last backend response.
func (m modifyStatusCodes) localResponseStatusCode() int {
	return m.response.LastBackendResponse().StatusCode()
}

// modifyResponseLocal modifies the status code of the last backend response in the history of the response object to
// the given statusCode.
// It returns a new backendResponse object with the modified status code.
func (m modifyStatusCodes) modifyResponseLocal(statusCode int) backendResponse {
	return m.response.LastBackendResponse().ModifyStatusCode(statusCode)
}

// modifyResponseGlobal is a method of the modifyStatusCodes struct. It takes an
// integer representing a status code and an instance of backendResponse as
// parameters. The method will return a Response after modifying the status
// code of the input backend response.
//
// Parameters:
//
//	statusCode - An integer representation of an HTTP status code.
//	backendResponseVO - An instance of backendResponse containing the server response.
//
// Returns:
//
//	Response - The response with a modified status code based on the input.
func (m modifyStatusCodes) modifyResponseGlobal(statusCode int, backendResponseVO backendResponse) Response {
	return m.response.ModifyStatusCode(statusCode, backendResponseVO)
}
