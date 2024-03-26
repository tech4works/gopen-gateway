package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	jsoniter "github.com/json-iterator/go"
	"github.com/ohler55/ojg/jp"
	"regexp"
	"strings"
)

type modify struct {
	scope    enum.ModifierScope
	action   enum.ModifierAction
	key      string
	value    string
	request  Request
	response Response
}

type modifyHeaders struct {
	modify
	globalHeader *Header
	localHeader  *Header
}

type modifyParams struct {
	modify
	globalParams *Params
	localParams  *Params
	localPath    *string
}

type modifyQueries struct {
	modify
	globalQuery *Query
	localQuery  *Query
}

type modifyBodies struct {
	modify
	globalBody any
	localBody  any
}

type modifyStatusCodes struct {
	modify
	globalStatusCode *int
	localStatusCode  *int
}

func newModify(modifierVO Modifier, requestVO Request, responseVO Response) modify {
	return modify{
		scope:    modifierVO.Scope,
		action:   modifierVO.Action,
		key:      modifierVO.Key,
		value:    modifierVO.Value,
		request:  requestVO,
		response: responseVO,
	}
}

func NewModifyHeaders(modifierVO Modifier, globalHeader, localHeader *Header, requestVO Request, responseVO Response,
) modifyHeaders {
	return modifyHeaders{
		modify:       newModify(modifierVO, requestVO, responseVO),
		globalHeader: globalHeader,
		localHeader:  localHeader,
	}
}

func NewModifyParams(modifierVO Modifier, globalParams, localParams *Params, localPath *string, requestVO Request,
	responseVO Response) modifyParams {
	return modifyParams{
		modify:       newModify(modifierVO, requestVO, responseVO),
		globalParams: globalParams,
		localParams:  localParams,
		localPath:    localPath,
	}
}

func NewModifyQueries(modifierVO Modifier, globalQuery, localQuery *Query, requestVO Request, responseVO Response,
) modifyQueries {
	return modifyQueries{
		modify:      newModify(modifierVO, requestVO, responseVO),
		globalQuery: globalQuery,
		localQuery:  localQuery,
	}
}

func NewModifyBodies(modifierVO Modifier, globalBody, localBody any, requestVO Request, responseVO Response,
) modifyBodies {
	return modifyBodies{
		modify:     newModify(modifierVO, requestVO, responseVO),
		globalBody: globalBody,
		localBody:  localBody,
	}
}

func NewModifyStatusCodes(modifierVO Modifier, globalStatusCode, localStatusCode *int, requestVO Request,
	responseVO Response) modifyStatusCodes {
	return modifyStatusCodes{
		modify:           newModify(modifierVO, requestVO, responseVO),
		globalStatusCode: globalStatusCode,
		localStatusCode:  localStatusCode,
	}
}

func (m modifyHeaders) Execute() {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.ValueStr()

	switch m.action {
	case enum.ModifierActionSet:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.isGlobalScope() {
			*m.globalHeader = m.globalHeader.Set(m.key, modifierValue)
		}
		*m.localHeader = m.localHeader.Set(m.key, modifierValue)
		break
	case enum.ModifierActionAdd:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.isGlobalScope() {
			*m.globalHeader = m.globalHeader.Add(m.key, modifierValue)
		}
		*m.localHeader = m.localHeader.Add(m.key, modifierValue)
		break
	case enum.ModifierActionDel:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.isGlobalScope() {
			*m.globalHeader = m.globalHeader.Del(m.key)
		}
		*m.localHeader = m.localHeader.Del(m.key)
		break
	}
}

func (m modifyHeaders) isGlobalScope() bool {
	return helper.Equals(m.scope, enum.ModifierScopeGlobal) && helper.IsNotNil(m.globalHeader)
}

func (m modifyParams) Execute() {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.ValueStr()

	// obtemos o valor a ser substituído do parâmetro
	paramUrl := fmt.Sprintf("/:%s", m.key)

	// obtemos o valor antigo caso exista
	oldLocalParam := m.localParams.Get(m.key)

	// se tiver valor, damos prioridade para substituir ou remover o valor
	if helper.IsNotEmpty(oldLocalParam) {
		paramUrl = fmt.Sprintf("/%s", oldLocalParam)
	}

	switch m.action {
	case enum.ModifierActionAdd, enum.ModifierActionSet:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.isGlobalScope() {
			*m.globalParams = m.globalParams.Set(m.key, modifierValue)
		}
		*m.localParams = m.localParams.Set(m.key, modifierValue)

		// adicionamos ou modificamos o valor da url caso tenha o parâmetro
		*m.localPath = strings.ReplaceAll(*m.localPath, paramUrl, fmt.Sprintf("/%s", modifierValue))
		break
	case enum.ModifierActionDel:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.isGlobalScope() {
			*m.globalParams = m.globalParams.Del(m.key)
		}
		*m.localParams = m.localParams.Del(m.key)

		// removemos o param de url no backend atual
		*m.localPath = strings.ReplaceAll(*m.localPath, paramUrl, "")
		break
	}
}

func (m modifyParams) isGlobalScope() bool {
	return helper.Equals(m.scope, enum.ModifierScopeGlobal) && helper.IsNotNil(m.globalParams)
}

func (m modifyQueries) Execute() {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.ValueStr()

	switch m.action {
	case enum.ModifierActionSet:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.isGlobalScope() {
			*m.globalQuery = m.globalQuery.Set(m.key, modifierValue)
		}
		*m.localQuery = m.localQuery.Set(m.key, modifierValue)
		break
	case enum.ModifierActionAdd:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.isGlobalScope() {
			*m.globalQuery = m.globalQuery.Add(m.key, modifierValue)
		}
		*m.localQuery = m.localQuery.Add(m.key, modifierValue)
		break
	case enum.ModifierActionDel:
		// se o escopo da modificação é global, modificamos o mesmo
		if m.isGlobalScope() {
			*m.globalQuery = m.globalQuery.Del(m.key)
		}
		*m.localQuery = m.localQuery.Del(m.key)
		break
	}
}

func (m modifyQueries) isGlobalScope() bool {
	return helper.Equals(m.scope, enum.ModifierScopeGlobal) && helper.IsNotNil(m.globalQuery)
}

func (m modifyBodies) Execute() {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.Value()

	// modificamos o body atual pelo tipo de dado
	if helper.IsJson(m.localBody) {
		m.executeJson(m.localBody, modifierValue)
	} else if helper.IsString(m.localBody) {
		m.executeString(m.localBody, modifierValue)
	}

	// caso seja em um escopo global, modificamos pelo tipo de dado também
	if m.isGlobalScope() {
		if helper.IsJson(m.globalBody) {
			m.executeJson(m.globalBody, modifierValue)
		} else if helper.IsString(m.globalBody) {
			m.executeString(m.globalBody, modifierValue)
		}
	}
}

func (m modifyBodies) isGlobalScope() bool {
	return helper.Equals(m.scope, enum.ModifierScopeGlobal) && helper.IsNotNil(m.globalBody)
}

func (m modifyBodies) executeJson(body, modifierValue any) {
	// damos o parse string da chave que eu quero modificar
	expr, err := jp.ParseString(m.key)
	if helper.IsNotNil(err) {
		return
	}

	// abaixo verificamos qual ação desejada para modificar o valor body
	switch m.action {
	case enum.ModifierActionSet, enum.ModifierActionAdd, enum.ModifierActionReplace:
		_ = expr.Set(body, modifierValue)
		break
	case enum.ModifierActionDel:
		_ = expr.Del(body)
		break
	case enum.ModifierActionRename:
		values := expr.Get(body)
		if helper.IsNotEmpty(values) {
			expr, err = jp.ParseString(m.value)
			if helper.IsNil(err) {
				_ = expr.Del(body)
				_ = expr.Set(body, values[len(values)-1])
			}
		}
		break
	}
}

func (m modifyBodies) executeString(body any, modifierValue any) {
	// convertemos o valor a ser modificado em str
	modifierValueStr := helper.SimpleConvertToString(modifierValue)
	// convertemos o body para string para garantir
	valueToModify := helper.SimpleConvertToString(body)

	// inicializamos o valor a ser modificado
	modifiedValue := valueToModify

	// modificamos a string com base no action fornecido
	switch m.action {
	case enum.ModifierActionAdd:
		modifiedValue = valueToModify + modifierValueStr
		break
	case enum.ModifierActionSet:
		modifiedValue = strings.ReplaceAll(valueToModify, m.key, modifierValueStr)
		break
	case enum.ModifierActionDel:
		modifiedValue = strings.ReplaceAll(valueToModify, m.key, "")
		break
	case enum.ModifierActionReplace:
		modifiedValue = modifierValueStr
		break
	}

	// preenchemos o body atual no valor modificado
	helper.SimpleConvertToDest(modifiedValue, body)
}

func (m modifyStatusCodes) Execute() {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.ValueInt()

	// se nao tiver valor nao fazemos nada
	if helper.IsEmpty(modifierValue) {
		return
	}

	// setamos o valor no ponteiro
	*m.localStatusCode = modifierValue

	// se for em scope global setamos o valor
	if m.isGlobalScope() {
		*m.globalStatusCode = modifierValue
	}
}

func (m modifyStatusCodes) isGlobalScope() bool {
	return helper.Equals(m.scope, enum.ModifierScopeGlobal) && helper.IsNotNil(m.globalStatusCode)
}

func (m modify) ValueInt() int {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.Value()

	// retornamos o valor modificado ou não
	return helper.SimpleConvertToInt(modifierValue)
}

func (m modify) ValueStr() string {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.Value()

	// retornamos o valor modificado ou não
	return helper.SimpleConvertToString(modifierValue)
}

func (m modify) Value() any {
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
	for _, word := range find { //responses[0].body.token or //requests[0].body.auth.token
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
		if helper.Contains(split[0], "requests") {
			evalValue = m.requestValueByEval(m.request, eval)
		} else if helper.Contains(split[0], "responses") {
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
	expr, err := jp.ParseString(strings.Replace(eval, "request", "", 1))
	if helper.IsNil(err) {
		values := expr.Get(requestVO)
		if helper.IsNotEmpty(values) {
			return values[len(values)-1]
		}
	}
	return nil
}

func (m modify) responseValueByEval(responseVO Response, eval string) any {
	expr, err := jp.ParseString(strings.Replace(eval, "response", "", 1))
	if helper.IsNil(err) {
		values := expr.Get(responseVO)
		if helper.IsNotEmpty(values) {
			return values[len(values)-1]
		}
	}
	return nil
}
