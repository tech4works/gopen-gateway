package vo

// modifyQueries represents a type that allows modification of queries.
// It is composed of the modify struct, which provides the necessary
// functionality for modifying queries.
type modifyQueries struct {
	modify
}

// NewModifyQueries creates a new instance of modifyQueries struct
// with the provided Modifier, Request, and Response.
func NewModifyQueries(modifierVO Modifier, requestVO Request, responseVO Response) ModifierStrategy {
	return modifyQueries{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// Execute runs the method executeRequestScope from the modifyQueries receiver
// and returns the corresponding Request and Response.
//
// The method is executed from the standard scope.
//
// Returns:
//
//	Request: The executed request.
//	Response: The response from the executed request.
func (m modifyQueries) Execute() (Request, Response) {
	// executamos a partir do escopo padrão
	return m.executeRequestScope()
}

// executeRequestScope makes modifications to the globalRequestQuery and the localRequestQuery.
// It modifies the local query by calling modifyLocalRequest and the global request by calling modifyGlobalRequest.
// It returns a modified global request and a response once the process is done.
//
// This method is a part of modifyQueries structure, and it encapsulates the logic to handle the execution
// of the request in its scope.
//
// The queries are modified by handling the global request query and local request params according
// to the particular requirements of the function.
//
// Returns:
//  1. Request: The modified global request
//  2. Response: The result of the executeRequestScope operation
func (m modifyQueries) executeRequestScope() (Request, Response) {
	// chamamos o modify de queries passando as queries a ser modificado e o mesmo retorna os modificados
	globalQuery, localQuery := m.queries(m.globalRequestQuery(), m.localRequestQuery())

	// modificamos o query local
	backendRequestVO := m.modifyLocalRequest(localQuery)

	// modificamos o params propagate e retornamos
	return m.modifyGlobalRequest(globalQuery, backendRequestVO), m.response
}

// globalRequestQuery returns the propagate query of the request.
func (m modifyQueries) globalRequestQuery() Query {
	return m.request.Query()
}

// localRequestQuery returns the query of the current backend request.
func (m modifyQueries) localRequestQuery() Query {
	return m.request.CurrentBackendRequest().Query()
}

// modifyLocalRequest modifies the local query of the current backend request and returns the modified backend request.
// It takes a localQuery of type Query as an argument and returns a backendRequest.
// The local query is modified using the ModifyQuery method of the current backend request.
// The modified backend request is created with the updated local query and the other existing properties of the current backend request.
// The modified backend request is then returned.
func (m modifyQueries) modifyLocalRequest(localQuery Query) backendRequest {
	return m.request.CurrentBackendRequest().ModifyQuery(localQuery)
}

// modifyGlobalRequest takes a propagate Query and a backendRequest as inputs.
// This function calls the ModifyQuery method on the m.request with globalQuery and backendRequestVO as parameters.
// It returns a modified Request.
func (m modifyQueries) modifyGlobalRequest(globalQuery Query, backendRequestVO backendRequest) Request {
	return m.request.ModifyQuery(globalQuery, backendRequestVO)
}
