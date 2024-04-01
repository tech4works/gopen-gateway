package vo

import (
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/ohler55/ojg/jp"
	"regexp"
	"strings"
)

type Modifier struct {
	context enum.ModifierContext
	scope   enum.ModifierScope
	action  enum.ModifierAction
	global  bool
	key     string
	value   string
}

type modify struct {
	action   enum.ModifierAction
	scope    enum.ModifierScope
	global   bool
	key      string
	value    string
	request  Request
	response Response
}

type modifyHeaders struct {
	modify
}

type modifyParams struct {
	modify
}

type modifyQueries struct {
	modify
}

type modifyBodies struct {
	modify
}

type modifyStatusCodes struct {
	modify
}

// NewModifyStatusCodes creates a new instance of modifyStatusCodes struct
// with the provided Modifier, Request, and Response.
func NewModifyStatusCodes(modifierVO Modifier, requestVO Request, responseVO Response) modifyStatusCodes {
	return modifyStatusCodes{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// NewModifyHeaders creates a new instance of modifyHeaders struct
// with the provided Modifier, Request, and Response.
func NewModifyHeaders(modifierVO Modifier, requestVO Request, responseVO Response) modifyHeaders {
	return modifyHeaders{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// NewModifyParams creates a new instance of modifyParams struct
// with the provided Modifier, Request, and Response.
func NewModifyParams(modifierVO Modifier, requestVO Request, responseVO Response) modifyParams {
	return modifyParams{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// NewModifyQueries creates a new instance of modifyQueries struct
// with the provided Modifier, Request, and Response.
func NewModifyQueries(modifierVO Modifier, requestVO Request, responseVO Response) modifyQueries {
	return modifyQueries{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// NewModifyBodies creates a new instance of modifyBodies struct
// with the provided Modifier, Request, and Response.
func NewModifyBodies(modifierVO Modifier, requestVO Request, responseVO Response) modifyBodies {
	return modifyBodies{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

// newModifier creates a new instance of Modifier struct
// with the provided modifierDTO.
// It sets the context, scope, action, global, key, and value fields of the Modifier struct
// to the corresponding fields of the modifierDTO parameter.
// Returns the created Modifier struct.
func newModifier(modifierDTO dto.Modifier) Modifier {
	return Modifier{
		context: modifierDTO.Context,
		scope:   modifierDTO.Scope,
		action:  modifierDTO.Action,
		global:  modifierDTO.Global,
		key:     modifierDTO.Key,
		value:   modifierDTO.Value,
	}
}

// newModify creates a new instance of the modify struct
// based on the provided Modifier, Request, and Response.
// If the scope is empty, it sets the default value based on the context.
// The default value for the request context is "request",
// and the default value for the response context is "response".
func newModify(modifierVO Modifier, requestVO Request, responseVO Response) modify {
	scope := modifierVO.scope
	// se o escopo ta vazio então setamos o valor padrão do context
	if helper.IsEmpty(scope) {
		switch modifierVO.context {
		case enum.ModifierContextRequest:
			scope = enum.ModifierScopeRequest
			break
		case enum.ModifierContextResponse:
			scope = enum.ModifierScopeResponse
			break
		}
	}

	return modify{
		action:   modifierVO.action,
		scope:    scope,
		global:   modifierVO.global,
		key:      modifierVO.key,
		value:    modifierVO.value,
		request:  requestVO,
		response: responseVO,
	}
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

// Execute executes the modifyStatusCodes by calling the executeResponseScope method and returns the modified Request
// and Response.
// It starts execution from the default scope.
func (m modifyStatusCodes) Execute() (Request, Response) {
	// executamos a partir do escopo padrão
	return m.executeResponseScope()
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

// modifyRequestLocal modifies the local request header of the backend request by applying the provided local header.
// It creates a new instance of the backendRequest struct with the modified header and returns it.
func (m modifyHeaders) modifyRequestLocal(localHeader Header) backendRequest {
	return m.request.CurrentBackendRequest().ModifyHeader(localHeader)
}

// modifyRequestGlobal modifies the global request by modifying the global header and the backend request.
// It returns a new modified Request.
func (m modifyHeaders) modifyRequestGlobal(globalHeader Header, backendRequestVO backendRequest) Request {
	return m.request.ModifyHeader(globalHeader, backendRequestVO)
}

// modifyResponseLocal modifies the local header of the response by applying changes from the given localHeader.
// It returns a new backendResponse object with the modified header.
func (m modifyHeaders) modifyResponseLocal(localHeader Header) backendResponse {
	return m.response.LastBackendResponse().ModifyHeader(localHeader)
}

// modifyResponseGlobal modifies the global response by applying the given global header and backend response.
// It returns the modified response.
func (m modifyHeaders) modifyResponseGlobal(globalHeader Header, backendResponseVO backendResponse) Response {
	return m.response.ModifyHeader(globalHeader, backendResponseVO)
}

// executeRequestScope executes the request scope of the modifyHeaders method.
// It calls the headers method, passing the global and local request headers, and returns the modified headers.
// Then, it modifies the local header by calling the modifyRequestLocal method.
// Finally, it modifies the global header by calling the modifyRequestGlobal method and returns the modified global
// header and the response.
// Returns:
// - Request: Modified global header and backend request object.
// - Response: The same response object.
func (m modifyHeaders) executeRequestScope() (Request, Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	globalHeader, localHeader := m.headers(m.globalRequestHeader(), m.localRequestHeader())

	// modificamos o header local
	backendRequestVO := m.modifyRequestLocal(localHeader)

	// modificamos o header global e retornamos
	return m.modifyRequestGlobal(globalHeader, backendRequestVO), m.response
}

// executeResponseScope executes the response scope of the modifyHeaders method.
// It calls the headers method, passing the global and local response headers, and returns the modified headers.
// Then, it modifies the local header by calling the modifyResponseLocal method.
// Finally, it modifies the global header by calling the modifyResponseGlobal method and returns the same request and the
// modified global header response.
// Returns:
// - Request: The same request object.
// - Response: Modified global header and backend response object.
func (m modifyHeaders) executeResponseScope() (Request, Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	globalHeader, localHeader := m.headers(m.globalResponseHeader(), m.localResponseHeader())

	// modificamos o header local
	backendResponseVO := m.modifyResponseLocal(localHeader)

	// modificamos o header global e retornamos
	return m.request, m.modifyResponseGlobal(globalHeader, backendResponseVO)
}

// Execute executes the modifyHeaders functionality.
// It determines the scope and calls the appropriate method to execute the modification.
// If the scope is enum.ModifierScopeRequest, it calls executeRequestScope and returns the modified request and original response.
// If the scope is enum.ModifierScopeResponse, it calls executeResponseScope and returns the original request and modified response.
// If the scope is neither enum.ModifierScopeRequest nor enum.ModifierScopeResponse, it returns the original request and response.
func (m modifyHeaders) Execute() (Request, Response) {
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

// Execute executes the modifyParams by calling the executeRequestScope method and returns the modified Request and Response.
// The execution starts from the default scope.
func (m modifyParams) Execute() (Request, Response) {
	// executamos a partir do escopo padrão
	return m.executeRequestScope()
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

// statusCodes modifies and returns the global and local status codes based on the condition defined.
// It first retrieves the value to be used for modification from the method m.valueInt().
// If the modifierValue is empty, the function exits early, and the global and local status codes remain unchanged.
// Else, the localStatusCode is set to the modifierValue.
// If the scope is global (m.global), the globalStatusCode is also set to modifierValue.
// It returns the possibly modified globalStatusCode and localStatusCode.
func (m modify) statusCodes(globalStatusCode, localStatusCode int) (int, int) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.valueInt()

	// se nao tiver valor nao fazemos nada
	if helper.IsEmpty(modifierValue) {
		return globalStatusCode, localStatusCode
	}

	// setamos o valor
	localStatusCode = modifierValue

	// se for em scope global setamos o valor
	if m.global {
		globalStatusCode = modifierValue
	}

	return globalStatusCode, localStatusCode
}

// Function headers modifies both global and local headers based on the action.
// It accepts two parameters - globalHeader and localHeader of type Header.
// The function performs the action (set, add, delete) on both global and local headers based on the modifierValue.
// The action and the key for the operation are based on the receiver 'm' of the type modifier.
// If the global field of the receiver 'm' is true, the action is applied on the globalHeader, else only on the localHeader.
//
// The function returns the modified globalHeader and localHeader.
func (m modify) headers(globalHeader, localHeader Header) (Header, Header) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.valueStr()

	switch m.action {
	case enum.ModifierActionSet:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.global {
			globalHeader = globalHeader.Set(m.key, modifierValue)
		}
		localHeader = localHeader.Set(m.key, modifierValue)
		break
	case enum.ModifierActionAdd:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.global {
			globalHeader = globalHeader.Add(m.key, modifierValue)
		}
		localHeader = localHeader.Add(m.key, modifierValue)
		break
	case enum.ModifierActionDel:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.global {
			globalHeader = globalHeader.Del(m.key)
		}
		localHeader = localHeader.Del(m.key)
		break
	}

	return globalHeader, localHeader
}

// params is a method on the 'modify' struct which returns the modified local path, global parameters, and local parameters.
// It modifies the parameters based on the action specified in the 'modify' struct (Set, Add, Del).
// When the global scope is set in the modify struct, it performs the modification on the global parameters as well.
// If the key provided in 'modify' struct does not exist in current local path, it is added to the path.
// The method also handles the removal of keys from the path.
// It takes in a local path (string), global parameters(map), and local parameters (map) as input.
// It returns a string (New local path), map(global parameters), map(local parameters).
func (m modify) params(localPath string, globalParams, localParams Params) (string, Params, Params) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.valueStr()

	// construímos o valor do key com padrão a ser modificado caso não exista
	paramUrl := fmt.Sprintf("/:%s", m.key)

	switch m.action {
	case enum.ModifierActionSet, enum.ModifierActionAdd:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.global {
			globalParams = globalParams.Set(m.key, modifierValue)
		}
		localParams = localParams.Set(m.key, modifierValue)

		// se o parâmetro não conte no path atual, adicionamos
		if !strings.Contains(localPath, paramUrl) {
			// checamos se no fim da url tem o /
			if localPath[len(localPath)-1] == '/' {
				localPath = fmt.Sprintf("%s:%s", localPath, m.key)
			} else {
				localPath = fmt.Sprintf("%s/:%s", localPath, m.key)
			}
		}
		break
	case enum.ModifierActionDel:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.global {
			globalParams = globalParams.Del(m.key)
		}
		localParams = localParams.Del(m.key)

		// removemos o param de url no backend atual
		localPath = strings.ReplaceAll(localPath, paramUrl, "")
		break
	}

	return localPath, globalParams, localParams
}

// queries modifies the queries based on the modifier action set (set, add, del),
// the key and value specified in the Modify object, and whether the scope is global or not.
// If the scope of the modification is global, it also modifies the global query.
// It returns modified global and local queries.
func (m modify) queries(globalQuery, localQuery Query) (Query, Query) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.valueStr()

	switch m.action {
	case enum.ModifierActionSet:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.global {
			globalQuery = globalQuery.Set(m.key, modifierValue)
		}
		localQuery = localQuery.Set(m.key, modifierValue)
		break
	case enum.ModifierActionAdd:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.global {
			globalQuery = globalQuery.Add(m.key, modifierValue)
		}
		localQuery = localQuery.Add(m.key, modifierValue)
		break
	case enum.ModifierActionDel:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.global {
			globalQuery = globalQuery.Del(m.key)
		}
		localQuery = localQuery.Del(m.key)
		break
	}

	return globalQuery, localQuery
}

// bodies modifies the global and local bodies based on the modifier action and value.
// It evaluates the modifier value and modifies the current body accordingly.
// If the local body is of JSON type, it is modified using the bodyJson() method.
// If the local body is of string type, it is modified using the bodyString() method.
// If the scope is global, the global body is also modified in the same way.
// The modified global and local bodies are returned.
func (m modify) bodies(globalBody, localBody Body) (Body, Body) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.valueEval()

	// modificamos o body atual pelo tipo de dado
	if helper.IsJsonType(localBody) {
		localBody = m.bodyJson(localBody, modifierValue)
	} else if helper.IsStringType(localBody) {
		localBody = m.bodyString(localBody, modifierValue)
	}

	// caso seja em um escopo global, modificamos pelo tipo de dado também
	if m.global {
		if helper.IsJsonType(globalBody) {
			globalBody = m.bodyJson(globalBody, modifierValue)
		} else if helper.IsStringType(globalBody) {
			globalBody = m.bodyString(globalBody, modifierValue)
		}
	}

	return globalBody, localBody
}

// bodyJson takes a body and a modifierValue of any type and returns a modified Body.
// It uses the modify field key to determine which part of the body to modify.
// Based on the modify action field, it performs different actions:
// - For ModifierActionSet, ModifierActionAdd, ModifierActionReplace it sets the new value to the body.
// - For ModifierActionDel, it deletes the key from the body.
// - For ModifierActionRename, it changes the key of the body to the new value, retaining the original value.
// After modification, it uses the body's Modify method to apply the changes.
func (m modify) bodyJson(body Body, modifierValue any) Body {
	// damos o parse string da chave que eu quero modificar
	expr, err := jp.ParseString(m.key)
	if helper.IsNotNil(err) {
		return body
	}

	// instanciamos a interface do body para ser modificada
	bodyToModify := body.Interface()

	// abaixo verificamos qual ação desejada para modificar o valor body
	switch m.action {
	case enum.ModifierActionSet, enum.ModifierActionAdd, enum.ModifierActionReplace:
		_ = expr.Set(bodyToModify, modifierValue)
		break
	case enum.ModifierActionDel:
		_ = expr.Del(bodyToModify)
		break
	case enum.ModifierActionRename:
		values := expr.Get(bodyToModify)
		if helper.IsNotEmpty(values) {
			exprValue, errValue := jp.ParseString(m.value)
			if helper.IsNil(errValue) {
				_ = expr.Del(bodyToModify)
				_ = exprValue.Set(bodyToModify, values[len(values)-1])
				m.key = m.value
			}
		}
		break
	}

	// chamamos modify do body objeto de valor para ele alterar os dados sem perder a ordenação
	return body.Modify(m.key, bodyToModify)
}

// bodyString modifies the body based on the provided action and returns the modified body.
//
// It converts the provided modifierValue to a string, ensures the body is also a string,
// and then modifies the body according to the provided action. The actions can be Add, Set, Del, or Replace.
//
// For the Add action, modifierValue is appended to the body.
// For the Set action, all instances of the key in the body are replaced with the modifierValue.
// For the Del action, all instances of the key in the body are deleted.
// For the Replace action, the body is replaced entirely with the modifierValue.
//
// Parameters:
// body - the original body to be modified.
// modifierValue - the value to be used for modification.
//
// Returns:
// The modified body.
func (m modify) bodyString(body Body, modifierValue any) Body {
	// convertemos o valor a ser modificado em str
	modifierValueStr := helper.SimpleConvertToString(modifierValue)
	// convertemos o body para string para garantir
	bodyToModify := helper.SimpleConvertToString(body.Interface())

	// inicializamos o valor a ser modificado
	modifiedValue := bodyToModify

	// modificamos a string com base no action fornecido
	switch m.action {
	case enum.ModifierActionAdd:
		modifiedValue = bodyToModify + modifierValueStr
		break
	case enum.ModifierActionSet:
		modifiedValue = strings.ReplaceAll(bodyToModify, m.key, modifierValueStr)
		break
	case enum.ModifierActionDel:
		modifiedValue = strings.ReplaceAll(bodyToModify, m.key, "")
		break
	case enum.ModifierActionReplace:
		modifiedValue = modifierValueStr
		break
	}

	return newBodyByAny(modifiedValue)
}

// valueInt method in the modify struct initializes the modifier value by calling
// the valueEval method, and then returns either the modified value or the original value.
// The return value is converted to an integer using the SimpleConvertToInt helper function.
func (m modify) valueInt() int {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.valueEval()

	// retornamos o valor modificado ou não
	return helper.SimpleConvertToInt(modifierValue)
}

// valueStr returns the modified value as a string.
// The value is obtained by evaluating the `valueEval()` method, which initializes the value to be used for modification.
// The modified value is then converted to a string using the `helper.SimpleConvertToString()` function.
// The function returns the modified value as a string.
func (m modify) valueStr() string {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.valueEval()

	// retornamos o valor modificado ou não
	return helper.SimpleConvertToString(modifierValue)
}

// valueEval method of the modify struct performs value evaluation.
// It initializes the value to be potentially modified.
// If the action is DEL, it returns nil.
// It uses a regex to find all the eval values within modifierValue.
// Iterates over these values and performs various operations based on them.
// Checks if evalValue comes from requests or responses.
// If the evalValue doesn't exist, it skips to the next one.
// If the value found equals the pre-defined word, it returns, otherwise replacing the eval key with the value obtained.
// Trying to parse the modifierValue string to bytes to check if it's JSON.
// If it is, it transforms it into an object.
// It finally returns the modified value.
// Note: Uses helper functions and requires encoding/json for json operations.
func (m modify) valueEval() any {
	// inicializamos o valor a ser modificado ou não
	modifierValue := m.value

	// caso a action seja DEL retornamos nil
	if helper.Equals(m.action, enum.ModifierActionDel) {
		return nil
	}

	// criamos o regex de evaluation esperado para obter o valor
	regex := regexp.MustCompile(`\B#[a-zA-Z0-9_.\[\]]+`)
	// buscamos todos os valores no modifierValue com esse valor eval
	find := regex.FindAllString(modifierValue, -1)
	// iteramos os valores eval
	for _, word := range find { //response.body.token or //request.body.auth.token
		// limpamos a #
		eval := strings.ReplaceAll(word, "#", "")

		// damos o split pela pontuação
		split := strings.Split(eval, ".")
		// caso esteja vazio vamos para o próximo
		if helper.IsEmpty(split) {
			continue
		}

		// obtemos o valor da eval vindo pela requests or responses
		var evalValue any
		if helper.Contains(split[0], "request") {
			evalValue = m.requestValueByEval(m.request, eval)
		} else if helper.Contains(split[0], "response") {
			evalValue = m.responseValueByEval(m.response, eval)
		}
		// caso o valor não encontrado, vamos para próximo
		if helper.IsNil(evalValue) {
			continue
		}

		// se a palavra é igual ao valor prescrito ja retornamos, caso contrário damos o replace do eval key pelo valor obtido
		if helper.Equals(word, modifierValue) {
			return evalValue
		} else {
			evalValueString := helper.SimpleConvertToString(evalValue)
			modifierValue = strings.Replace(modifierValue, word, evalValueString, 1)
		}
	}

	// damos o parse da string para bytes para verificarmos se o valor é um json, se for, transformamos o mesmo em objeto
	modifierValueBytes := []byte(modifierValue)
	if json.Valid(modifierValueBytes) {
		var obj any
		err := json.Unmarshal(modifierValueBytes, &obj)
		if helper.IsNil(err) {
			return obj
		}
	}

	// retornamos o valor modificado
	return modifierValue
}

// requestValueByEval is a method associated with the modify struct. This method evaluates a string input,
// retrieves a value from the provided Request object based on the evaluation.
//
// Parameters:
// requestVO (type Request): This is used as the source for the `eval` evaluation.
// eval (type string): This is evaluated after replacing the "request." prefix with an empty string.
//
// Procedure:
// First, the `eval` string is parsed into a JSONPath expression using the jp.ParseString method after replacing "request.".
// If there is an error during parsing, the method returns nil.
// If parsing succeeds, the expression is then applied to the requestVO to fetch values.
// If the result yields multiple values, only the last value is returned.
//
// The function returns a single value of any type, or nil on encountering parsing errors.
func (m modify) requestValueByEval(requestVO Request, eval string) any {
	expr, err := jp.ParseString(strings.Replace(eval, "request.", "", 1))
	if helper.IsNil(err) {
		values := expr.Get(requestVO.Eval())
		if helper.IsNotEmpty(values) {
			return values[len(values)-1]
		}
	}
	return nil
}

// responseValueByEval takes a Response object and an 'eval' string as parameters. It attempts to parse the eval argument,
// replacing occurrences of "response." with an empty string. If parsing is successful and doesn't return an error,
// it executes a Get method on the returned expression using Eval method of Response as an argument.
// If the retrieved values are not empty, the last value in the values slice is returned. Otherwise, or if an error occurs
// during parsing, it returns nil.
//
// The function expects:
// - responseVO: A Response struct.
// - eval: A string representing an eval field.
//
// It returns:
//   - An interface that contains either the last value of a slice or nil if an error occurs during parsing or the values
//     are empty.
func (m modify) responseValueByEval(responseVO Response, eval string) any {
	expr, err := jp.ParseString(strings.Replace(eval, "response.", "", 1))
	if helper.IsNil(err) {
		values := expr.Get(responseVO.Eval())
		if helper.IsNotEmpty(values) {
			return values[len(values)-1]
		}
	}
	return nil
}
