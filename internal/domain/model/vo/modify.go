/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

// modify represents a modification operation to be performed on a request or response.
// It contains fields such as action, scope, propagate, key, value, request, and response,
// which define the details of the modification operation.
type modify struct {
	// action represents the action to be performed.
	action enum.ModifierAction
	// scope represents the scope of a modifier.
	scope enum.ModifierScope
	// propagate is a boolean variable used to control whether a certain operation should be propagated or not.
	propagate bool
	// key represents the key of the field to be modified
	key string
	// value represents the value to be inserted to modify the object
	value string
	// request represents an HTTP `request` object.
	request *Request
	// response represents an HTTP `response` object.
	response *Response
}

// ModifierStrategy represents a strategy for executing a modification operation on a Request and Response.
// It defines a single method Execute() which takes no arguments and returns a Request and Response.
// Implementations of this interface should provide their own Execute() method implementation.
type ModifierStrategy interface {
	// Execute is a method of the ModifierStrategy interface which represents a strategy for executing a modification
	// operation on a Request and Response.
	// It takes no arguments and returns a Request and Response.
	// The implementation of this method should define the logic for executing the modification operation.
	Execute() (*Request, *Response)
}

// NewModifyVOFunc represents a function type that takes a Modifier, Request, and Response as input parameters,
// and returns a ModifierStrategy.
type NewModifyVOFunc func(modifierVO *Modifier, requestVO *Request, responseVO *Response) ModifierStrategy

// newModify initializes a new modify object based on the given inputs.
// It first retrieves the scope from modifierVO.
// If the scope is empty, it sets it based on the context from modifierVO.
// Then it constructs the modify object with the following properties:
// - action: the action from modifierVO
// - scope: the scope determined previously
// - propagate: the propagate flag from modifierVO
// - key: the key from modifierVO
// - value: the value from modifierVO
// - request: the requestVO
// - response: the responseVO
// The modify object is then returned.
func newModify(modifierVO *Modifier, requestVO *Request, responseVO *Response) modify {
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
		action:    modifierVO.Action(),
		scope:     scope,
		propagate: modifierVO.Propagate(),
		key:       modifierVO.Key(),
		value:     modifierVO.Value(),
		request:   requestVO,
		response:  responseVO,
	}
}

// statusCode modifies the statusCode based on the receiver 'm' of the type modify.
// It obtains the value to be used for modification using the m.valueInt() method.
// If the modifierValue is empty, it returns the original statusCode without any modifications.
// Otherwise, it sets the modifierValue as the new statusCode and returns it.
func (m modify) statusCode(statusCode int) int {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.intEvalValue()

	// se nao tiver valor nao fazemos nada
	if helper.IsEmpty(modifierValue) {
		return statusCode
	}

	// retornamos o valor modificado
	return modifierValue
}

// headers modifies the globalHeader and localHeader based on the receiver 'm' of the type modify.
// It obtains the slice of string values to be used for modification using the m.sliceOfStrEvalValue() method.
// If m.propagate is true, it modifies the globalHeader by calling the m.header method with globalHeader and modifierValue.
// Finally, it returns the modified globalHeader and the result of calling the m.header method with localHeader and modifierValue.
func (m modify) headers(globalHeader, localHeader Header) (Header, Header) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.sliceOfStrEvalValue()

	// caso seja em um escopo propagate, modificamos o header global
	if m.propagate {
		globalHeader = m.header(globalHeader, modifierValue)
	}

	// retornamos os objetos de valores modificados, ou não
	return globalHeader, m.header(localHeader, modifierValue)
}

// header modifies the provided header based on the action specified in the receiver 'm' of the type modify.
// It accepts a Header object and an array of string values as parameters.
// It performs different modifications on the header based on the action value:
// - If the action is enum.ModifierActionAdd, it adds all the modifier values to the header under the specified key.
// - If the action is enum.ModifierActionApd, it appends the modifier values to the existing values associated with the key.
// - If the action is enum.ModifierActionSet, it sets the modifier values as the new values associated with the key.
// - If the action is enum.ModifierActionRpl, it replaces the existing values associated with the key with the modifier values.
// - If the action is enum.ModifierActionRen, it renames the key to the last value in the modifier values array.
// - If the action is enum.ModifierActionDel, it deletes the entry for the specified key from the header.
// If the action is not any of the above, it returns the original header unchanged.
// Finally, it returns the modified header.
func (m modify) header(header Header, modifierValue []string) Header {
	// alteramos os headers conforme o action indicada
	switch m.action {
	case enum.ModifierActionAdd:
		return header.AddAll(m.key, modifierValue)
	case enum.ModifierActionApd:
		return header.Append(m.key, modifierValue)
	case enum.ModifierActionSet:
		return header.SetAll(m.key, modifierValue)
	case enum.ModifierActionRpl:
		return header.Replace(m.key, modifierValue)
	case enum.ModifierActionRen:
		return header.Rename(m.key, modifierValue[len(modifierValue)-1])
	case enum.ModifierActionDel:
		return header.Delete(m.key)
	default:
		return header
	}
}

// params modifies the localPath, globalParams, and localParams based on the receiver 'm' of the type modify.
// It obtains the value to be used for modification using the m.strEvalValue() method.
// If m.propagate is true, it modifies the globalParams by calling the m.param() method with the modifierValue.
// It then returns the modified localPath, globalParams, and localParams
// by calling the m.urlPath(), m.param(), and m.param() methods respectively with the modifierValue.
//
// Returns:
//   - UrlPath: The modified local request path.
//   - Params: The modified global request parameters.
//   - Params: The modified local request parameters.
func (m modify) params(localPath UrlPath, globalParams, localParams Params) (UrlPath, Params, Params) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.strEvalValue()

	// caso seja em um escopo propagate, modificamos o params global
	if m.propagate {
		globalParams = m.param(globalParams, modifierValue)
	}

	// retornamos o path possívelmente alterado, o globa e local params possívelmente alterados
	return m.urlPath(localPath, modifierValue), globalParams, m.param(localParams, modifierValue)
}

// urlPath modifies the given `urlPath` based on the receiver `m` of the type `modify`.
// It performs different actions based on the value of `m.action`:
// - enum.ModifierActionSet or enum.ModifierActionRpl: It sets the parameter key `m.key` in the `urlPath`.
// - enum.ModifierActionRen: It renames the parameter key `m.key` to `modifierValue` in the `urlPath`.
// - enum.ModifierActionDel: It deletes the parameter key `m.key` from the `urlPath`.
// If `m.action` is not any of the above, it returns the unchanged `urlPath`.
// It returns the modified `urlPath`.
func (m modify) urlPath(urlPath UrlPath, modifierValue string) UrlPath {
	// alteramos o parâmetro pela action indicada
	switch m.action {
	case enum.ModifierActionSet, enum.ModifierActionRpl:
		return urlPath.SetParamKey(m.key)
	case enum.ModifierActionRen:
		return urlPath.RenameParamKey(m.key, modifierValue)
	case enum.ModifierActionDel:
		return urlPath.DeleteParamKey(m.key)
	default:
		return urlPath
	}
}

// param modifies the `params` based on the `action`, `key`, and `modifierValue` provided.
// It performs the following operations based on the `action`:
//   - If the `action` is `enum.ModifierActionSet`, it sets the `modifierValue` for the
//     specified `key` in the `params` and returns the updated `params`.
//   - If the `action` is `enum.ModifierActionRpl`, it replaces the value of the specified `key`
//     in the `params` with the `modifierValue` and returns the updated `params`.
//   - If the `action` is `enum.ModifierActionRen`, it renames the specified `key` to the
//     `modifierValue` in the `params` and returns the updated `params`.
//   - If the `action` is `enum.ModifierActionDel`, it deletes the specified `key` from the `params`
//     and returns the updated `params`.
//   - If the `action` does not match any of the above cases, it returns the original `params` as is.
func (m modify) param(params Params, modifierValue string) Params {
	// alteramos o parâmetro pela action indicada
	switch m.action {
	case enum.ModifierActionSet:
		return params.Set(m.key, modifierValue)
	case enum.ModifierActionRpl:
		return params.Replace(m.key, modifierValue)
	case enum.ModifierActionRen:
		return params.Rename(m.key, modifierValue)
	case enum.ModifierActionDel:
		return params.Delete(m.key)
	default:
		return params
	}
}

// queries modifies the globalQuery and localQuery based on the receiver 'm' of the type modify.
// It obtains the value to be used for modification using the m.sliceOfStrEvalValue() method.
// If m.propagate is true, it modifies the globalQuery by calling the m.query() method.
// It then returns the modified globalQuery and localQuery by calling the m.query() method.
func (m modify) queries(globalQuery, localQuery Query) (Query, Query) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.sliceOfStrEvalValue()

	// caso seja em um escopo propagate, modificamos o query global
	if m.propagate {
		globalQuery = m.query(globalQuery, modifierValue)
	}

	// retornamos a query global e local possivelmente alteradas
	return globalQuery, m.query(localQuery, modifierValue)
}

// query modifies the given query based on the action specified in the receiver 'm' of the type modify.
// It takes a Query and modifierValue as parameters and returns a modified Query.
//
// If the action is ModifierActionAdd, it calls the Add method of the Query with the specified key and modifierValue.
// If the action is ModifierActionApd, it calls the Append method of the Query with the specified key and modifierValue.
// If the action is ModifierActionSet, it calls the Set method of the Query with the specified key and modifierValue.
// If the action is ModifierActionRpl, it calls the Replace method of the Query with the specified key and modifierValue.
// If the action is ModifierActionRen, it calls the Rename method of the Query with the specified key and the last element of modifierValue.
// If the action is ModifierActionDel, it calls the Delete method of the Query with the specified key.
//
// If the action is not one of the predefined actions above, it returns the original Query without any modifications.
//
// Returns:
//
//	A modified Query based on the action and key specified in the receiver 'm'.
//
// Note:
//
//	The original Query is not modified.
func (m modify) query(query Query, modifierValue []string) Query {
	// alteramos o query pela action indicada
	switch m.action {
	case enum.ModifierActionAdd:
		return query.Add(m.key, modifierValue)
	case enum.ModifierActionApd:
		return query.Append(m.key, modifierValue)
	case enum.ModifierActionSet:
		return query.Set(m.key, modifierValue)
	case enum.ModifierActionRpl:
		return query.Replace(m.key, modifierValue)
	case enum.ModifierActionRen:
		return query.Rename(m.key, modifierValue[len(modifierValue)-1])
	case enum.ModifierActionDel:
		return query.Delete(m.key)
	default:
		return query
	}
}

// bodies modifies the globalBody and localBody based on the receiver 'm' of the type modifyBodies.
// It obtains the value to be used for modification using the m.evalValue() method.
// If m.propagate is true, it calls the m.body method on globalBody and assigns the modified result back to globalBody.
// The method then calls m.body method on localBody and assigns the modified result back to localBody.
// The modified globalBody and localBody are then returned.
func (m modify) bodies(globalBody, localBody *Body) (*Body, *Body) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.evalValue()

	// caso seja em um escopo propagate, modificamos pelo tipo de dado também
	if m.propagate {
		globalBody = m.body(globalBody, modifierValue)
	}

	// retornamos o body global e local possivelmente alterados
	return globalBody, m.body(localBody, modifierValue)
}

// body modifies the body based on the receiver 'm' of the type modify.
// It takes a pointer to a Body object and a modifierValue of any type.
// If the body pointer is nil, it returns nil.
// If the action is enum.ModifierActionAdd, it calls the Add method of the body object and assigns the modifiedBody and
// error to modifiedBody and err, respectively.
// If the action is enum.ModifierActionApd, it calls the Append method of the body object and assigns the modifiedBody
// and error to modifiedBody and err, respectively.
// If the action is enum.ModifierActionSet, it calls the Set method of the body object and assigns the modifiedBody and
// error to modifiedBody and err, respectively.
// If the action is enum.ModifierActionRpl, it calls the Replace method of the body object and assigns the modifiedBody
// and error to modifiedBody and err, respectively.
// If the action is enum.ModifierActionRen, it calls the Rename method of the body object and assigns the modifiedBody
// and error to modifiedBody and err, respectively.
// If the action is enum.ModifierActionDel, it calls the Delete method of the body object and assigns the modifiedBody
// and error to modifiedBody and err, respectively.
// If an error occurs during the modification, it is handled but not logged.
// It returns the modifiedBody, which is the body object after the modification.
func (m modify) body(body *Body, modifierValue any) *Body {
	// se for nil ja retornamo
	if helper.IsNil(body) {
		return nil
	}

	// instânciamos o body a ser modificado
	var modifiedBody = body
	var err error

	// abaixo verificamos qual ação desejada para modificar o valor body
	switch m.action {
	case enum.ModifierActionAdd:
		modifiedBody, err = body.Add(m.key, modifierValue)
	case enum.ModifierActionApd:
		modifiedBody, err = body.Append(m.key, modifierValue)
	case enum.ModifierActionSet:
		modifiedBody, err = body.Set(m.key, modifierValue)
	case enum.ModifierActionRpl:
		modifiedBody, err = body.Replace(m.key, modifierValue)
	case enum.ModifierActionRen:
		modifiedBody, err = body.Rename(m.key, modifierValue)
	case enum.ModifierActionDel:
		modifiedBody, err = body.Delete(m.key)
	default:
		return body
	}

	// tratamos o erro e retornamos o próprio body
	if helper.IsNotNil(err) {
		// todo: imprimir log?
	}

	// caso tenha dado tudo certo retornamos o body modificado
	return modifiedBody
}

// intEvalValue method in the modify struct initializes the modifier value by calling
// the evalValue method, and then returns either the modified value or the original value.
// The return value is converted to an integer using the SimpleConvertToInt helper function.
func (m modify) intEvalValue() int {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.evalValue()
	// retornamos o valor modificado ou não
	return helper.SimpleConvertToInt(modifierValue)
}

// strEvalValue returns the modified value as a string.
// The value is obtained by evaluating the `valueEval()` method, which initializes the value to be used for modification.
// The modified value is then converted to a string using the `helper.SimpleConvertToString()` function.
// The function returns the modified value as a string.
func (m modify) strEvalValue() string {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.evalValue()
	// retornamos o valor modificado ou não
	return helper.SimpleConvertToString(modifierValue)
}

// sliceOfStrEvalValue returns a slice of string values to be used for modification.
// It initializes the modifierValue by calling the evalValue method of the receiver 'm'.
// If the modifierValue is a slice, it converts it to a []string using the SimpleConvertToDest method.
// If the converted slice is not empty, it returns the slice.
// Otherwise, it converts the modifierValue to a string using the SimpleConvertToString method,
// and returns it as a single-element []string.
func (m modify) sliceOfStrEvalValue() []string {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.evalValue()
	// verificamos se o mesmo é um slice
	if helper.IsSliceType(modifierValue) {
		var ss []string
		helper.SimpleConvertToDest(modifierValue, &ss)
		if helper.IsNotEmpty(ss) {
			return ss
		}
	}
	return []string{helper.SimpleConvertToString(modifierValue)}
}

// evalValue evaluates the modifierValue based on the receiver 'm' of the type modify.
// If the action is "DEL", it returns nil.
// Otherwise, it iterates through the values and processes them based on eval syntax.
// Finally, it parses the modifierValue to the appropriate data type and returns it.
func (m modify) evalValue() any {
	// caso a action seja DEL retornamos nil
	if helper.Equals(m.action, enum.ModifierActionDel) {
		return nil
	}
	// inicializamos o valor a ser modificado ou não
	modifierValue := m.value
	// iteramos os valores com base na sintaxe eval
	for _, word := range m.findAllByEvalSintaxe(modifierValue) {
		// processamos o palavra para converter em um valor eval
		modifierValue = m.processEvalWord(modifierValue, word)
	}
	// damos o parse do valor em string caso tenha um valor do tipo não string
	return m.parseModifierValueToRealType(modifierValue)
}

// findAllByEvalSintaxe searches for all values in 'value' that match the expected evaluation regex.
// It creates the regex pattern for the expected evaluation syntax to obtain the value.
// It returns an array with all values found in 'value' that match the evaluation syntax.
func (m modify) findAllByEvalSintaxe(value string) []string {
	// criamos o regex de evaluation esperado para obter o valor
	regex := regexp.MustCompile(`\B#[a-zA-Z0-9_.\[\]]+`)
	// buscamos todos os valores no modifierValue com esse valor eval
	return regex.FindAllString(value, -1)
}

// processEvalWord takes a modifierValue and a word as input.
// It searches for the evalValue based on the word using the evalValueByWord method of the receiver 'm' of type modify.
// If the evalValue is not found, it returns the modifierValue as it is.
// Otherwise, it converts the evalValue to a string using the SimpleConvertToString helper function and returns it.
func (m modify) processEvalWord(modifierValue, word string) string {
	// obtemos o valor pela palavra
	evalValue := m.evalValueByWord(word)

	// caso o valor não encontrado, vamos para próximo
	if helper.IsNil(evalValue) {
		return modifierValue
	}

	// retornamos o valor substituído pelo valor do eval encontrado
	evalValueString := helper.SimpleConvertToString(evalValue)
	return strings.Replace(modifierValue, word, evalValueString, 1)
}

// evalValueByWord evaluates the value associated with the given word.
// It removes all instances of "#" from the word.
// It then splits the modified word by "." to obtain individual components.
// If the split result is empty, it returns nil.
// It extracts the value from either the request or response based on the first component of split.
// If the first component contains "request", it calls m.requestValueByEval() with m.request and eval as arguments.
// If the first component contains "response", it calls m.responseValueByEval() with m.response and eval as arguments.
// Otherwise, it sets evalValue to nil.
// It returns the obtained evalValue.
func (m modify) evalValueByWord(word string) any {
	// limpamos a #
	eval := strings.ReplaceAll(word, "#", "")
	// damos o split pela pontuação
	split := strings.Split(eval, ".")
	// caso esteja vazio vamos para o próximo
	if helper.IsEmpty(split) {
		return nil
	}

	// obtemos o valor da eval vindo pela requests or responses
	var evalValue any
	if helper.Contains(split[0], "request") {
		evalValue = m.requestValueByEval(m.request, eval)
	} else if helper.Contains(split[0], "response") {
		evalValue = m.responseValueByEval(m.response, eval)
	} else {
		evalValue = nil
	}
	// retornamos o valor obtido
	return evalValue
}

// requestValueByEval obtains the value from the Request object based on the evaluation string 'eval'.
// It replaces the "request." substring with an empty string in the evaluation string to form the expression.
// Using the gjson.Get method, it retrieves the value from the evaluation string in the Request object.
// If the value exists, it returns the result. Otherwise, it returns nil.
func (m modify) requestValueByEval(requestVO *Request, eval string) any {
	expr := strings.Replace(eval, "request.", "", 1)
	result := gjson.Get(requestVO.Eval(), expr)
	if result.Exists() {
		return result.Value()
	}
	return nil
}

// responseValueByEval obtains the value from the responseVO object by evaluating the expression given by the eval string.
// The expression is modified by replacing "response." with an empty string using strings.Replace() method.
// The modified expression is then used to extract the value from the response using gjson.Get() method.
// If the value exists, it is returned as result.Value(). Otherwise, it returns nil.
func (m modify) responseValueByEval(responseVO *Response, eval string) any {
	expr := strings.Replace(eval, "response.", "", 1)
	result := gjson.Get(responseVO.Eval(), expr)
	if result.Exists() {
		return result.Value()
	}
	return nil
}

// parseModifierValueToRealType converts the modifierValue to its corresponding real type.
// If the modifierValue is of type json, int, float, bool, or time, it converts it to that type.
// It uses the ConvertToDest method from the helper package to perform the conversion.
// If the conversion is successful, it returns the converted value.
// Otherwise, it returns the original modifierValue.
func (m modify) parseModifierValueToRealType(modifierValue string) any {
	// caso seja do tipo json, int, float, bool, time convertemos a esses tipos
	if helper.IsJson(modifierValue) || helper.IsInt(modifierValue) || helper.IsFloat(modifierValue) ||
		helper.IsBool(modifierValue) || helper.IsTime(modifierValue) {
		var obj any
		err := helper.ConvertToDest(modifierValue, &obj)
		if helper.IsNil(err) {
			return obj
		}
	}
	// retornamos o valor modificado
	return modifierValue
}
