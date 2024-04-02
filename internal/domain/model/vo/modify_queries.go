package vo

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

// executeRequestScope executes the request scope by calling the modifyQueries method to modify the queries and returns the modified requests and the response.
// It modifies the global query and the local query by calling the globalRequestQuery and localRequestParams methods respectively.
// Then, it modifies the local query by calling the modifyRequestLocal method.
// Finally, it modifies the global params and returns the modified request and the response by calling the modifyRequestGlobal method.
// It returns the modified request and the response.
func (m modifyQueries) executeRequestScope() (Request, Response) {
	// chamamos o modify de queries passando as queries a ser modificado e o mesmo retorna os modificados
	globalQuery, localQuery := m.queries(m.globalRequestQuery(), m.localRequestParams())

	// modificamos o query local
	backendRequestVO := m.modifyRequestLocal(localQuery)

	// modificamos o params global e retornamos
	return m.modifyRequestGlobal(globalQuery, backendRequestVO), m.response
}

// globalRequestQuery returns the global query of the request.
func (m modifyQueries) globalRequestQuery() Query {
	return m.request.Query()
}

// localRequestParams returns the query parameters of the current backend request.
func (m modifyQueries) localRequestParams() Query {
	return m.request.CurrentBackendRequest().Query()
}

// modifyRequestLocal modifies the local query of the current backend request and returns the modified backend request.
// It takes a localQuery of type Query as an argument and returns a backendRequest.
// The local query is modified using the ModifyQuery method of the current backend request.
// The modified backend request is created with the updated local query and the other existing properties of the current backend request.
// The modified backend request is then returned.
func (m modifyQueries) modifyRequestLocal(localQuery Query) backendRequest {
	return m.request.CurrentBackendRequest().ModifyQuery(localQuery)
}

// modifyRequestGlobal takes a global Query and a backendRequest as inputs.
// This function calls the ModifyQuery method on the m.request with globalQuery and backendRequestVO as parameters.
// It returns a modified Request.
func (m modifyQueries) modifyRequestGlobal(globalQuery Query, backendRequestVO backendRequest) Request {
	return m.request.ModifyQuery(globalQuery, backendRequestVO)
}
