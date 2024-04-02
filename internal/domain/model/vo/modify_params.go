package vo

type modifyParams struct {
	modify
}

// NewModifyParams creates a new instance of modifyParams struct
// with the provided Modifier, Request, and Response.
func NewModifyParams(modifierVO Modifier, requestVO Request, responseVO Response) ModifierStrategy {
	return modifyParams{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// Execute executes the modifyParams by calling the executeRequestScope method and returns the modified Request and Response.
// The execution starts from the default scope.
func (m modifyParams) Execute() (Request, Response) {
	// executamos a partir do escopo padrão
	return m.executeRequestScope()
}

// executeRequestScope executes the modifications on request parameters.
// It first modifies the local parameters by calling the modify method on the params,
// using the local request path and both global and local request parameters.
// It then modifies the local parameters and
// finally modifies the global parameters before returning the updated Request and Response.
//
// Returns:
// Request: the updated request after modifications
// Response: the unchanged response
func (m modifyParams) executeRequestScope() (Request, Response) {
	// chamamos o modify de params passando o path e params a ser modificado e o mesmo retorna os mesmo modificados
	localPath, globalParams, localParams := m.params(m.localRequestPath(), m.globalRequestParams(), m.localRequestParams())

	// modificamos o params local
	backendRequestVO := m.modifyRequestLocal(localPath, localParams)

	// modificamos o params global e retornamos
	return m.modifyRequestGlobal(globalParams, backendRequestVO), m.response
}

// localRequestPath returns the path of the current backend request in the modifyParams struct.
func (m modifyParams) localRequestPath() string {
	return m.request.CurrentBackendRequest().Path()
}

// globalRequestParams returns the global request parameters of the request object in the modifyParams struct.
func (m modifyParams) globalRequestParams() Params {
	return m.request.Params()
}

// localRequestParams returns the local request parameters of the current backend request object in the modifyParams struct.
func (m modifyParams) localRequestParams() Params {
	return m.request.CurrentBackendRequest().Params()
}

// modifyRequestLocal is a method on the modifyParams type.
// This method takes in a 'localPath' of type string and 'localParams' of type Params,
// and returns a modified backendRequest.
//
// The 'localPath' represents the local path of the backend request.
//
// The 'localParams' represents the new local parameters to be used for the backend request.
//
// This method is used to modify the parameters of a local backend request with new parameters.
func (m modifyParams) modifyRequestLocal(localPath string, localParams Params) backendRequest {
	return m.request.CurrentBackendRequest().ModifyParams(localPath, localParams)
}

// modifyRequestGlobal is a method on the modifyParams struct.
// It modifies global request parameters based on provided Params and backendRequest.
// It takes two arguments - globalParams of type Params and backendRequestVO of type backendRequest.
// It returns a Request which is the modified version of the original Request.
//
// Parameters:
//
// globalParams: The Params type global parameters that need to be modified.
//
// backendRequestVO: The backendRequest type which contains the request sent to the backend.
//
// Returns:
//
// The modified Request after applying changes based on globalParams and backendRequestVO.
func (m modifyParams) modifyRequestGlobal(globalParams Params, backendRequestVO backendRequest) Request {
	return m.request.ModifyParams(globalParams, backendRequestVO)
}
