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
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

// modify represents a modification operation to be performed on a httpRequest or httpResponse.
// It contains fields such as action, scope, propagate, key, value, httpRequest, and httpResponse,
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
	// httpRequest represents an HTTP `httpRequest` object.
	httpRequest *HttpRequest
	// httpResponse represents an HTTP `httpResponse` object.
	httpResponse *HttpResponse
}

// ModifierStrategy represents a strategy for executing a modification operation on a HttpRequest and HttpResponse.
// It defines a single method Execute() which takes no arguments and returns a HttpRequest and HttpResponse.
// Implementations of this interface should provide their own Execute() method implementation.
type ModifierStrategy interface {
	// Execute is a method of the ModifierStrategy interface which represents a strategy for executing a modification
	// operation on a HttpRequest and HttpResponse.
	// It takes no arguments and returns a HttpRequest and HttpResponse.
	// The implementation of this method should define the logic for executing the modification operation.
	Execute() (*HttpRequest, *HttpResponse)
}

// NewModifyVOFunc represents a function type that takes a Modifier, HttpRequest, and HttpResponse as input parameters,
// and returns a ModifierStrategy.
type NewModifyVOFunc func(modifierVO *Modifier, httpRequestVO *HttpRequest, httpResponseVO *HttpResponse) ModifierStrategy

// newModify initializes a new modify object based on the given inputs.
// It first retrieves the scope from modifierVO.
// If the scope is empty, it sets it based on the context from modifierVO.
// Then it constructs the modify object with the following properties:
// - action: the action from modifierVO
// - scope: the scope determined previously
// - propagate: the propagate flag from modifierVO
// - key: the key from modifierVO
// - value: the value from modifierVO
// - httpRequest: the requestVO
// - httpResponse: the responseVO
// The modify object is then returned.
func newModify(modifierVO *Modifier, httpRequestVO *HttpRequest, httpResponseVO *HttpResponse) modify {
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
		action:       modifierVO.Action(),
		scope:        scope,
		propagate:    modifierVO.Propagate(),
		key:          modifierVO.Key(),
		value:        modifierVO.Value(),
		httpRequest:  httpRequestVO,
		httpResponse: httpResponseVO,
	}
}

// statusCode modifies the statusCode based on the receiver 'm' of the type modify.
// It obtains the value to be used for modification using the m.valueInt() method.
// If the modifierValueAsInt is empty, it returns the original statusCode without any modifications.
// Otherwise, it sets the modifierValueAsInt as the new statusCode and returns it.
func (m modify) statusCode(statusCode int) int {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.modifierValueAsInt()

	// se nao tiver valor nao fazemos nada
	if helper.IsEmpty(modifierValue) {
		return statusCode
	}

	// retornamos o valor modificado
	return modifierValue
}

func (m modify) headers(httpHeaderVO, httpBackendHeader Header) (Header, Header) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.modifierValueAsSliceOfStr()

	// caso seja em um escopo propagate, modificamos o header global
	if m.propagate {
		httpHeaderVO = m.header(httpHeaderVO, modifierValue)
	}

	// retornamos os objetos de valores modificados, ou não
	return httpHeaderVO, m.header(httpBackendHeader, modifierValue)
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
// If m.propagate is true, it modifies the globalParams by calling the m.param() method with the modifierValueAsInt.
// It then returns the modified localPath, globalParams, and localParams
// by calling the m.urlPath(), m.param(), and m.param() methods respectively with the modifierValueAsInt.
//
// Returns:
//   - UrlPath: The modified local httpRequest path.
//   - Params: The modified global httpRequest parameters.
//   - Params: The modified local httpRequest parameters.
func (m modify) params(localPath UrlPath, globalParams, localParams Params) (UrlPath, Params, Params) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.modifierValue()

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
// - enum.ModifierActionRen: It renames the parameter key `m.key` to `modifierValueAsInt` in the `urlPath`.
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

// param modifies the `params` based on the `action`, `key`, and `modifierValueAsInt` provided.
// It performs the following operations based on the `action`:
//   - If the `action` is `enum.ModifierActionSet`, it sets the `modifierValueAsInt` for the
//     specified `key` in the `params` and returns the updated `params`.
//   - If the `action` is `enum.ModifierActionRpl`, it replaces the value of the specified `key`
//     in the `params` with the `modifierValueAsInt` and returns the updated `params`.
//   - If the `action` is `enum.ModifierActionRen`, it renames the specified `key` to the
//     `modifierValueAsInt` in the `params` and returns the updated `params`.
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
	modifierValue := m.modifierValueAsSliceOfStr()

	// caso seja em um escopo propagate, modificamos o query global
	if m.propagate {
		globalQuery = m.query(globalQuery, modifierValue)
	}

	// retornamos a query global local possivelmente alteradas
	return globalQuery, m.query(localQuery, modifierValue)
}

// query modifies the given query based on the action specified in the receiver 'm' of the type modify.
// It takes a Query and modifierValueAsInt as parameters and returns a modified Query.
//
// If the action is ModifierActionAdd, it calls the Add method of the Query with the specified key and modifierValueAsInt.
// If the action is ModifierActionApd, it calls the Append method of the Query with the specified key and modifierValueAsInt.
// If the action is ModifierActionSet, it calls the Set method of the Query with the specified key and modifierValueAsInt.
// If the action is ModifierActionRpl, it calls the Replace method of the Query with the specified key and modifierValueAsInt.
// If the action is ModifierActionRen, it calls the Rename method of the Query with the specified key and the last element of modifierValueAsInt.
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

func (m modify) bodies(httpBodyVO, httpBackendBodyVO *Body) (*Body, *Body) {
	// obtemos o valor a ser usado para modificar
	modifierValue := m.modifierValue()

	// caso seja em um escopo propagate, modificamos pelo tipo de dado também
	if m.propagate {
		httpBodyVO = m.body(httpBodyVO, modifierValue)
	}

	// retornamos o body global e local possivelmente alterados
	return httpBodyVO, m.body(httpBackendBodyVO, modifierValue)
}

// body modifies the body based on the receiver 'm' of the type modify.
// It takes a pointer to a Body object and a modifierValueAsInt of any type.
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
func (m modify) body(body *Body, modifierValue string) *Body {
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
		logger.Warning("Error modify body:", err)
	}

	// caso tenha dado tudo certo retornamos o body modificado
	return modifiedBody
}

// modifierValueAsInt method in the modify struct initializes the modifier value by calling
// the modifierValue method, and then returns either the modified value or the original value.
// The return value is converted to an integer using the SimpleConvertToInt helper function.
func (m modify) modifierValueAsInt() int {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.modifierValue()
	// retornamos o valor modificado ou não
	return helper.SimpleConvertToInt(modifierValue)
}

// modifierValueAsSliceOfStr returns a slice of string values to be used for modification.
// It initializes the modifierValueAsInt by calling the modifierValue method of the receiver 'm'.
// If the modifierValueAsInt is a slice, it converts it to a []string using the SimpleConvertToDest method.
// If the converted slice is not empty, it returns the slice.
// Otherwise, it converts the modifierValueAsInt to a string using the SimpleConvertToString method,
// and returns it as a single-element []string.
func (m modify) modifierValueAsSliceOfStr() []string {
	// inicializamos o valor a ser usado para modificar
	modifierValue := m.modifierValue()
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

func (m modify) modifierValue() string {
	// caso a action seja DEL retornamos nil
	if helper.Equals(m.action, enum.ModifierActionDel) {
		return ""
	}
	// inicializamos o valor a ser modificado ou não
	modifierValue := m.value
	// iteramos os valores com base na sintaxe eval
	for _, dynamicValueWord := range m.findAllByDynamicValueSyntax(modifierValue) {
		// processamos o palavra para converter em um valor dinâmico
		modifierValue = m.processDynamicValueWord(modifierValue, dynamicValueWord)
	}
	// damos o parse do valor em string caso tenha um valor do tipo não string
	return modifierValue
}

func (m modify) findAllByDynamicValueSyntax(value string) []string {
	// criamos o regex de evaluation esperado para obter o valor
	regex := regexp.MustCompile(`\B#[a-zA-Z0-9_.\-\[\]]+`)
	// buscamos todos os valores no modifierValueAsInt com esse valor eval
	return regex.FindAllString(value, -1)
}

func (m modify) processDynamicValueWord(modifierValue, dynamicValueWord string) string {
	// obtemos o valor pela palavra
	dynamicValue := m.getDynamicValuePerWord(dynamicValueWord)

	// caso o valor não encontrado, vamos para próximo
	if helper.IsEmpty(dynamicValue) {
		return modifierValue
	}

	// substituímos a palavra pelo eval
	return strings.Replace(modifierValue, dynamicValueWord, dynamicValue, 1)
}

func (m modify) getDynamicValuePerWord(dynamicValueWord string) string {
	// limpamos a #
	cleanSintaxe := strings.ReplaceAll(dynamicValueWord, "#", "")
	// damos o split pela pontuação
	dotSplit := strings.Split(cleanSintaxe, ".")
	// caso esteja vazio vamos para o próximo
	if helper.IsEmpty(dotSplit) {
		return ""
	}

	// obtemos o valor da eval vindo pela requests or responses
	var dynamicValue string
	if helper.Contains(dotSplit[0], "request") {
		dynamicValue = m.getHttpRequestValueByJsonPath(m.httpRequest, cleanSintaxe)
	} else if helper.Contains(dotSplit[0], "response") {
		dynamicValue = m.getHttpResponseValueByEval(m.httpResponse, cleanSintaxe)
	} else {
		dynamicValue = ""
	}
	// retornamos o valor obtido
	return dynamicValue
}

func (m modify) getHttpRequestValueByJsonPath(requestVO *HttpRequest, eval string) string {
	expr := strings.Replace(eval, "request.", "", 1)
	result := gjson.Get(requestVO.Json(), expr)
	if result.Exists() {
		return result.String()
	}
	return ""
}

func (m modify) getHttpResponseValueByEval(responseVO *HttpResponse, eval string) string {
	expr := strings.Replace(eval, "response.", "", 1)
	result := gjson.Get(responseVO.Map(), expr)
	if result.Exists() {
		return result.String()
	}
	return ""
}
