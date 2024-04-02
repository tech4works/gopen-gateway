package vo

import (
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/ohler55/ojg/jp"
	"regexp"
	"strings"
)

type modify struct {
	action   enum.ModifierAction
	scope    enum.ModifierScope
	global   bool
	key      string
	value    string
	request  Request
	response Response
}

type ModifierStrategy interface {
	Execute() (Request, Response)
}

// NewModifyVOFunc represents a function type that takes a Modifier, Request, and Response as input parameters,
// and returns a ModifierStrategy.
type NewModifyVOFunc func(modifierVO Modifier, requestVO Request, responseVO Response) ModifierStrategy

// newModify initializes a new modify object based on the given inputs.
// It first retrieves the scope from modifierVO.
// If the scope is empty, it sets it based on the context from modifierVO.
// Then it constructs the modify object with the following properties:
// - action: the action from modifierVO
// - scope: the scope determined previously
// - global: the global flag from modifierVO
// - key: the key from modifierVO
// - value: the value from modifierVO
// - request: the requestVO
// - response: the responseVO
// The modify object is then returned.
func newModify(modifierVO Modifier, requestVO Request, responseVO Response) modify {
	// inicializamos o escopo padrão
	scope := modifierVO.Scope()

	// caso ele esteja vazio, setamos com base no context obrigatório fornecido
	if helper.IsEmpty(scope) {
		switch modifierVO.Context() {
		case enum.ModifierContextRequest:
			scope = enum.ModifierScopeRequest
		case enum.ModifierContextResponse:
			scope = enum.ModifierScopeResponse
		}
	}

	// construímos o objeto de valor para modificar
	return modify{
		action:   modifierVO.Action(),
		scope:    scope,
		global:   modifierVO.Global(),
		key:      modifierVO.Key(),
		value:    modifierVO.Value(),
		request:  requestVO,
		response: responseVO,
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
// and then modifies the body according to the provided action. The actions can Add, Set, Del, or Replace.
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
