package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	jsoniter "github.com/json-iterator/go"
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

func NewModifyStatusCodes(modifierVO Modifier, requestVO Request, responseVO Response) modifyStatusCodes {
	return modifyStatusCodes{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

func NewModifyHeaders(modifierVO Modifier, requestVO Request, responseVO Response) modifyHeaders {
	return modifyHeaders{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

func NewModifyParams(modifierVO Modifier, requestVO Request, responseVO Response) modifyParams {
	return modifyParams{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

func NewModifyQueries(modifierVO Modifier, requestVO Request, responseVO Response) modifyQueries {
	return modifyQueries{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

func NewModifyBodies(modifierVO Modifier, requestVO Request, responseVO Response) modifyBodies {
	return modifyBodies{
		modify: newModify(modifierVO, requestVO, responseVO),
	}
}

func (m modifyStatusCodes) globalResponseStatusCode() int {
	return m.response.StatusCode()
}

func (m modifyStatusCodes) localResponseStatusCode() int {
	return m.response.LastBackendResponse().StatusCode()
}

func (m modifyStatusCodes) modifyResponseLocal(statusCode int) backendResponse {
	return m.response.LastBackendResponse().ModifyStatusCode(statusCode)
}

func (m modifyStatusCodes) modifyResponseGlobal(statusCode int, backendResponseVO backendResponse) Response {
	return m.response.ModifyStatusCode(statusCode, backendResponseVO)
}

func (m modifyStatusCodes) executeResponseScope() (Request, Response) {
	// chamamos o modify de status code passando os status codes a ser modificado e o mesmo retorna modificados
	globalStatusCode, localStatusCode := m.statusCodes(m.globalResponseStatusCode(), m.localResponseStatusCode())

	// modificamos o status code local
	backendResponseVO := m.modifyResponseLocal(localStatusCode)

	// modificamos o status code global e retornamos
	return m.request, m.modifyResponseGlobal(globalStatusCode, backendResponseVO)
}

func (m modifyStatusCodes) Execute() (Request, Response) {
	// executamos a partir do escopo padrão
	return m.executeResponseScope()
}

func (m modifyHeaders) globalRequestHeader() Header {
	return m.request.Header()
}

func (m modifyHeaders) localRequestHeader() Header {
	return m.request.CurrentBackendRequest().Header()
}

func (m modifyHeaders) globalResponseHeader() Header {
	return m.response.Header()
}

func (m modifyHeaders) localResponseHeader() Header {
	return m.response.LastBackendResponse().Header()
}

func (m modifyHeaders) modifyRequestLocal(localHeader Header) backendRequest {
	return m.request.CurrentBackendRequest().ModifyHeader(localHeader)
}

func (m modifyHeaders) modifyRequestGlobal(globalHeader Header, backendRequestVO backendRequest) Request {
	return m.request.ModifyHeader(globalHeader, backendRequestVO)
}

func (m modifyHeaders) modifyResponseLocal(localHeader Header) backendResponse {
	return m.response.LastBackendResponse().ModifyHeader(localHeader)
}

func (m modifyHeaders) modifyResponseGlobal(globalHeader Header, backendResponseVO backendResponse) Response {
	return m.response.ModifyHeader(globalHeader, backendResponseVO)
}

func (m modifyHeaders) executeRequestScope() (Request, Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	globalHeader, localHeader := m.headers(m.globalRequestHeader(), m.localRequestHeader())

	// modificamos o header local
	backendRequestVO := m.modifyRequestLocal(localHeader)

	// modificamos o header global e retornamos
	return m.modifyRequestGlobal(globalHeader, backendRequestVO), m.response
}

func (m modifyHeaders) executeResponseScope() (Request, Response) {
	// chamamos o modify de headers passando o headers a ser modificado e o mesmo retorna os mesmo modificados
	globalHeader, localHeader := m.headers(m.globalResponseHeader(), m.localResponseHeader())

	// modificamos o header local
	backendResponseVO := m.modifyResponseLocal(localHeader)

	// modificamos o header global e retornamos
	return m.request, m.modifyResponseGlobal(globalHeader, backendResponseVO)
}

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

func (m modifyParams) localRequestPath() string {
	return m.request.CurrentBackendRequest().Path()
}

func (m modifyParams) globalRequestParams() Params {
	return m.request.Params()
}

func (m modifyParams) localRequestParams() Params {
	return m.request.CurrentBackendRequest().Params()
}

func (m modifyParams) modifyRequestLocal(localPath string, localParams Params) backendRequest {
	return m.request.CurrentBackendRequest().ModifyParams(localPath, localParams)
}

func (m modifyParams) modifyRequestGlobal(globalParams Params, backendRequestVO backendRequest) Request {
	return m.request.ModifyParams(globalParams, backendRequestVO)
}

func (m modifyParams) executeRequestScope() (Request, Response) {
	// chamamos o modify de params passando o path e params a ser modificado e o mesmo retorna os mesmo modificados
	localPath, globalParams, localParams := m.params(m.localRequestPath(), m.globalRequestParams(), m.localRequestParams())

	// modificamos o params local
	backendRequestVO := m.modifyRequestLocal(localPath, localParams)

	// modificamos o params global e retornamos
	return m.modifyRequestGlobal(globalParams, backendRequestVO), m.response
}

func (m modifyParams) Execute() (Request, Response) {
	// executamos a partir do escopo padrão
	return m.executeRequestScope()
}

func (m modifyQueries) globalRequestQuery() Query {
	return m.request.Query()
}

func (m modifyQueries) localRequestParams() Query {
	return m.request.CurrentBackendRequest().Query()
}

func (m modifyQueries) modifyRequestLocal(localQuery Query) backendRequest {
	return m.request.CurrentBackendRequest().ModifyQuery(localQuery)
}

func (m modifyQueries) modifyRequestGlobal(globalQuery Query, backendRequestVO backendRequest) Request {
	return m.request.ModifyQuery(globalQuery, backendRequestVO)
}

func (m modifyQueries) executeRequestScope() (Request, Response) {
	// chamamos o modify de queries passando as queries a ser modificado e o mesmo retorna os modificados
	globalQuery, localQuery := m.queries(m.globalRequestQuery(), m.localRequestParams())

	// modificamos o query local
	backendRequestVO := m.modifyRequestLocal(localQuery)

	// modificamos o params global e retornamos
	return m.modifyRequestGlobal(globalQuery, backendRequestVO), m.response
}

func (m modifyQueries) Execute() (Request, Response) {
	// executamos a partir do escopo padrão
	return m.executeRequestScope()
}

func (m modifyBodies) globalRequestBody() Body {
	return m.request.Body()
}

func (m modifyBodies) localRequestBody() Body {
	return m.request.CurrentBackendRequest().Body()
}

func (m modifyBodies) globalResponseBody() Body {
	return m.response.Body()
}

func (m modifyBodies) localResponseBody() Body {
	return m.response.LastBackendResponse().Body()
}

func (m modifyBodies) modifyRequestLocal(localBody Body) backendRequest {
	return m.request.CurrentBackendRequest().ModifyBody(localBody)
}

func (m modifyBodies) modifyRequestGlobal(globalBody Body, backendRequestVO backendRequest) Request {
	return m.request.ModifyBody(globalBody, backendRequestVO)
}

func (m modifyBodies) modifyResponseLocal(localBody Body) backendResponse {
	return m.response.LastBackendResponse().ModifyBody(localBody)
}

func (m modifyBodies) modifyResponseGlobal(globalBody Body, backendResponseVO backendResponse) Response {
	return m.response.ModifyBody(globalBody, backendResponseVO)
}

func (m modifyBodies) executeRequestScope() (Request, Response) {
	// chamamos o modify de bodies passando os bodies a ser modificado e o mesmo retorna modificados
	globalBody, localBody := m.bodies(m.globalRequestBody(), m.localRequestBody())

	// modificamos o body local
	backendRequestVO := m.modifyRequestLocal(localBody)

	// modificamos o body global e retornamos
	return m.modifyRequestGlobal(globalBody, backendRequestVO), m.response
}

func (m modifyBodies) executeResponseScope() (Request, Response) {
	// chamamos o modify de bodies passando os bodies a ser modificado e o mesmo retorna modificados
	globalBody, localBody := m.bodies(m.globalResponseBody(), m.localResponseBody())

	// modificamos o header local
	backendResponseVO := m.modifyResponseLocal(localBody)

	// modificamos o header global e retornamos
	return m.request, m.modifyResponseGlobal(globalBody, backendResponseVO)
}

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

func (m modify) valueInt() int {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.valueEval()

	// retornamos o valor modificado ou não
	return helper.SimpleConvertToInt(modifierValue)
}

func (m modify) valueStr() string {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.valueEval()

	// retornamos o valor modificado ou não
	return helper.SimpleConvertToString(modifierValue)
}

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
	if jsoniter.Valid(modifierValueBytes) {
		var obj any
		err := jsoniter.Unmarshal(modifierValueBytes, &obj)
		if helper.IsNil(err) {
			return obj
		}
	}

	// retornamos o valor modificado
	return modifierValue
}

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
