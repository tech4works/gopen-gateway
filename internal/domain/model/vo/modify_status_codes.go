package vo

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"strconv"
)

// modifyStatusCodes is a struct that extends the functionality of the modify type.
// It is used to modify status codes in HTTP responses.
type modifyStatusCodes struct {
	modify
}

// NewModifyStatusCodes creates a new instance of modifyStatusCodes struct
// with the provided Modifier, Request, and Response.
func NewModifyStatusCodes(statusCodeValue int, requestVO *Request, responseVO *Response) ModifierStrategy {
	statusCodeStr := strconv.Itoa(statusCodeValue)
	return modifyStatusCodes{
		modify: newModify(newModifierFromValue(enum.ModifierContextResponse, statusCodeStr), requestVO, responseVO),
	}
}

// Execute executes the modifyStatusCodes by calling the executeResponseScope method and returns the original Request
// and modified Response.
// It starts execution from the default scope.
func (m modifyStatusCodes) Execute() (*Request, *Response) {
	// executamos a partir do escopo padrão
	return m.executeResponseScope()
}

// executeResponseScope is a method of the modifyStatusCodes structure.
// It modifies both the local and global response status codes based on provided
// status codes and returns both updated request and response.
//
// It first calls the modify method of status code with the last backend response's
// status code, which gets modified.
//
// Then, it modifies the local response status code.
// Afterward, it modifies the global response with the new backendResponseVO.
// Finally, it returns the modified request and response.
//
// Returns:
//   - Request: The original (possibly modified) request.
//   - Response: The modified response after operations.
func (m modifyStatusCodes) executeResponseScope() (*Request, *Response) {
	// chamamos o modify de status code passando os status codes a ser modificado e o mesmo retorna modificados
	statusCode := m.statusCode(m.response.LastBackendResponse().StatusCode())

	// modificamos o status code local
	backendResponseVO := m.modifyLocalResponse(statusCode)

	// modificamos o response global com o novo backendResponseVO
	return m.request, m.modifyGlobalResponse(backendResponseVO)
}

// modifyLocalResponse modifies the status code of the last backend response in the history of the response object to
// the given statusCode.
// It returns a new backendResponse object with the modified status code.
func (m modifyStatusCodes) modifyLocalResponse(statusCode int) *backendResponse {
	return m.response.LastBackendResponse().ModifyStatusCode(statusCode)
}

// modifyGlobalResponse modifies the global Response by calling the ModifyLastBackendResponse method on the Response object.
// It takes a backendResponseVO parameter and returns the modified Response.
// The modified Response is obtained by calling the ModifyLastBackendResponse method on the Response object,
// passing the backendResponseVO as the argument.
func (m modifyStatusCodes) modifyGlobalResponse(backendResponseVO *backendResponse) *Response {
	return m.response.ModifyLastBackendResponse(backendResponseVO)
}
