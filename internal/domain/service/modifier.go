package service

import (
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/enum"
	"github.com/GabrielHCataldo/martini-gateway/internal/domain/model/valueobject"
	"github.com/iancoleman/orderedmap"
	"github.com/ohler55/ojg/jp"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type ModifierRequest struct {
	Host     string
	Endpoint string
	Url      string
	Method   string
	Header   http.Header
	Query    url.Values
	Params   map[string]string
	Body     any
}

type ModifierResponse struct {
	StatusCode int
	Header     http.Header
	Body       any
}

type ExecuteRequestScopeInput struct {
	Request         ModifierRequest
	Headers         []valueobject.Modifier
	Params          []valueobject.Modifier
	Queries         []valueobject.Modifier
	Body            []valueobject.Modifier
	RequestHistory  []ModifierRequest
	ResponseHistory []ModifierResponse
}

type ExecuteResponseScopeInput struct {
	Response        ModifierResponse
	Headers         []valueobject.Modifier
	Params          []valueobject.Modifier
	Queries         []valueobject.Modifier
	Body            []valueobject.Modifier
	RequestHistory  []ModifierRequest
	ResponseHistory []ModifierResponse
}

type executeModifierInput struct {
	Scope     enum.ModifierScope
	Headers   []valueobject.Modifier
	Params    []valueobject.Modifier
	Queries   []valueobject.Modifier
	Body      []valueobject.Modifier
	Requests  []ModifierRequest
	Responses []ModifierResponse
}

type modifier struct {
}

type Modifier interface {
	ExecuteRequestScope(input ExecuteRequestScopeInput) (*ModifierRequest, error)
	ExecuteResponseScope(input ExecuteResponseScopeInput) (*ModifierResponse, error)
}

func NewModifier() Modifier {
	return modifier{}
}

func (m modifier) ExecuteRequestScope(input ExecuteRequestScopeInput) (*ModifierRequest, error) {
	input.RequestHistory = append(input.RequestHistory, input.Request)
	inputExecModifier := executeModifierInput{
		Scope:     enum.ModifierScopeRequest,
		Headers:   input.Headers,
		Params:    input.Params,
		Queries:   input.Queries,
		Body:      input.Body,
		Requests:  input.RequestHistory,
		Responses: input.ResponseHistory,
	}
	err := m.executeModifier(inputExecModifier)
	if helper.IsNotNil(err) {
		return nil, err
	}
	request := m.getCurrentRequest(inputExecModifier.Requests)
	return &request, nil
}

func (m modifier) ExecuteResponseScope(input ExecuteResponseScopeInput) (*ModifierResponse, error) {
	input.ResponseHistory = append(input.ResponseHistory, input.Response)
	inputExecModifier := executeModifierInput{
		Scope:     enum.ModifierScopeResponse,
		Headers:   input.Headers,
		Params:    input.Params,
		Queries:   input.Queries,
		Body:      input.Body,
		Requests:  input.RequestHistory,
		Responses: input.ResponseHistory,
	}
	err := m.executeModifier(inputExecModifier)
	if helper.IsNotNil(err) {
		return nil, err
	}
	response := m.getCurrentResponse(inputExecModifier.Responses)
	return &response, nil
}

func (m modifier) executeModifier(input executeModifierInput) error {
	for _, modifierHeader := range input.Headers {
		m.modifierHeader(input.Scope, modifierHeader, input.Requests, input.Responses)
	}
	for _, modifierParam := range input.Params {
		m.modifierParams(input.Scope, modifierParam, input.Requests, input.Responses)
	}
	for _, modifierQuery := range input.Queries {
		m.modifierQueries(input.Scope, modifierQuery, input.Requests, input.Responses)
	}
	for _, modifierBody := range input.Body {
		m.modifierBody(input.Scope, modifierBody, input.Requests, input.Responses)
	}
	return nil
}

func (m modifier) modifierHeader(scope enum.ModifierScope, modifier valueobject.Modifier, requests []ModifierRequest,
	responses []ModifierResponse) {
	if helper.NotContains(modifier.Scope, scope) && helper.NotContains(modifier.Scope, "*") {
		return
	}
	var value string
	if helper.IsNotEqualTo(modifier.Action, enum.ModifierActionDel) {
		valueEval := m.getValueEval(modifier.Value, requests, responses)
		value = helper.SimpleConvertToString(valueEval)
	}
	var header http.Header
	switch scope {
	case enum.ModifierScopeRequest:
		header = m.getCurrentRequest(requests).Header
		break
	case enum.ModifierScopeResponse:
		header = m.getCurrentResponse(responses).Header
		break
	}
	switch modifier.Action {
	case enum.ModifierActionSet:
		header.Set(modifier.Key, value)
		break
	case enum.ModifierActionAdd:
		header.Add(modifier.Key, value)
		break
	case enum.ModifierActionDel:
		header.Del(modifier.Key)
		break
	}
}

func (m modifier) modifierParams(scope enum.ModifierScope, modifier valueobject.Modifier, requests []ModifierRequest,
	responses []ModifierResponse) {
	if helper.NotContains(modifier.Scope, scope) && helper.NotContains(modifier.Scope, "*") {
		return
	} else if helper.IsNotEqualTo(scope, enum.ModifierScopeRequest) && helper.IsNotEqualTo(scope, "*") {
		panic(errors.New("params modifier not allow on scope:", scope))
	}
	var value string
	if helper.IsNotEqualTo(modifier.Action, enum.ModifierActionDel) {
		valueEval := m.getValueEval(modifier.Value, requests, responses)
		value = helper.SimpleConvertToString(valueEval)
	}
	request := m.getCurrentRequest(requests)
	if helper.Equals(modifier.Action, enum.ModifierActionAdd) || helper.Equals(modifier.Action, enum.ModifierActionSet) {
		request.Params[modifier.Key] = value
	} else if helper.Equals(modifier.Action, enum.ModifierActionDel) {
		delete(request.Params, modifier.Key)
		rKey := fmt.Sprint("/:", modifier.Key)
		request.Endpoint = strings.Replace(request.Endpoint, rKey, "", 1)
		request.Url = fmt.Sprint(request.Host, request.Endpoint)
	}
}

func (m modifier) modifierQueries(scope enum.ModifierScope, modifier valueobject.Modifier, requests []ModifierRequest,
	responses []ModifierResponse) {
	if helper.NotContains(modifier.Scope, scope) && helper.NotContains(modifier.Scope, "*") {
		return
	} else if helper.IsNotEqualTo(scope, enum.ModifierScopeRequest) && helper.IsNotEqualTo(scope, "*") {
		panic(errors.New("queries modifier not allow on scope:" + scope))
	}
	var value string
	if helper.IsNotEqualTo(modifier.Action, enum.ModifierActionDel) {
		valueEval := m.getValueEval(modifier.Value, requests, responses)
		value = helper.SimpleConvertToString(valueEval)
	}
	request := m.getCurrentRequest(requests)
	switch modifier.Action {
	case enum.ModifierActionSet:
		request.Query.Set(modifier.Key, value)
		break
	case enum.ModifierActionAdd:
		request.Query.Add(modifier.Key, value)
		break
	case enum.ModifierActionDel:
		request.Query.Del(modifier.Key)
		break
	}
}

func (m modifier) modifierBody(scope enum.ModifierScope, modifier valueobject.Modifier, requests []ModifierRequest,
	responses []ModifierResponse) {
	if helper.NotContains(modifier.Scope, scope) && helper.NotContains(modifier.Scope, "*") {
		return
	}
	var value any
	if helper.IsNotEqualTo(modifier.Action, enum.ModifierActionDel) &&
		helper.IsNotEqualTo(modifier.Action, enum.ModifierActionReplace) {
		value = m.getValueEval(modifier.Value, requests, responses)
	}
	var body any
	var bodyToModifier any
	switch scope {
	case enum.ModifierScopeRequest:
		requestsEval := m.convertRequestsToEval(requests)
		reqBodyEval := m.getCurrentRequest(requestsEval).Body
		bodyToModifier = reqBodyEval
		body = m.getCurrentRequest(requests).Body
		break
	case enum.ModifierScopeResponse:
		responsesEval := m.convertResponsesToEval(responses)
		resBodyEval := m.getCurrentResponse(responsesEval).Body
		bodyToModifier = resBodyEval
		body = m.getCurrentResponse(responses).Body
		break
	}
	if helper.IsJson(bodyToModifier) {
		m.modifierBodyJson(modifier, body, bodyToModifier, value)
	} else if helper.IsNotNil(bodyToModifier) {
		m.modifierBodyString(modifier, bodyToModifier, value)
	}
}

func (m modifier) modifierBodyJson(modifier valueobject.Modifier, body, bodyToModifier any, value any) {
	x, err := jp.ParseString(modifier.Key)
	if helper.IsNotNil(err) {
		return
	}
	switch modifier.Action {
	case enum.ModifierActionSet, enum.ModifierActionAdd:
		_ = x.Set(bodyToModifier, value)
		break
	case enum.ModifierActionDel:
		_ = x.Del(bodyToModifier)
		break
	case enum.ModifierActionReplace:
		results := x.Get(bodyToModifier)
		if helper.IsNotEmpty(results) {
			y, err := jp.ParseString(modifier.Value)
			if helper.IsNil(err) {
				_ = x.Del(bodyToModifier)
				_ = y.Set(bodyToModifier, results[0])
			}
		}
		break
	}
	m.persistModifierOnBody(modifier, body, bodyToModifier)
}

func (m modifier) modifierBodyString(modifier valueobject.Modifier, body any, value any) {
	modifierValue := helper.SimpleConvertToString(value)
	valueToReplace := helper.SimpleConvertToString(body)
	switch modifier.Action {
	case enum.ModifierActionAdd:
		valueToReplace = valueToReplace + modifierValue
		break
	case enum.ModifierActionSet:
		valueToReplace = strings.ReplaceAll(valueToReplace, modifier.Key, modifierValue)
		break
	case enum.ModifierActionDel:
		valueToReplace = strings.ReplaceAll(valueToReplace, modifier.Key, "")
		break
	case enum.ModifierActionReplace:
		valueToReplace = modifierValue
		break
	}
	helper.SimpleConvertToDest(value, body)
}

func (m modifier) getValueEval(valueModifier string, requests []ModifierRequest, responses []ModifierResponse) any {
	regex := regexp.MustCompile(`\B#[a-zA-Z0-9_.\[\]]+`)
	find := regex.FindAllString(valueModifier, -1)
	for _, word := range find {
		eval := strings.ReplaceAll(word, "#", "") //responses[0].body.token or //requests[0].body.auth.token
		splitDot := strings.Split(eval, ".")
		if helper.IsEmpty(splitDot) {
			continue
		}
		var valueGet any
		if helper.Contains(splitDot[0], "requests") {
			valueGet = m.getRequestsValueByEval(requests, strings.Replace(eval, "requests", "", 1))
		} else if helper.Contains(splitDot[0], "responses") {
			valueGet = m.getResponsesValueByEval(responses, strings.Replace(eval, "responses", "", 1))
		}
		if helper.IsNil(valueGet) {
			continue
		}
		if helper.Equals(word, valueModifier) {
			return valueGet
		} else {
			valueGetString := helper.SimpleConvertToString(valueGet)
			valueModifier = strings.Replace(valueModifier, word, valueGetString, 1)
		}
	}
	valueModifierBytes := []byte(valueModifier)
	if json.Valid(valueModifierBytes) {
		var obj any
		err := json.Unmarshal(valueModifierBytes, &obj)
		if helper.IsNil(err) {
			return obj
		}
	}
	return valueModifier
}

func (m modifier) getRequestsValueByEval(requests []ModifierRequest, eval string) any {
	x, err := jp.ParseString(eval)
	if helper.IsNil(err) {
		requestsEval := m.convertRequestsToEval(requests)
		resultJsonPath := x.Get(requestsEval)
		if helper.IsNotEmpty(resultJsonPath) {
			return resultJsonPath[0]
		}
	}
	return nil
}

func (m modifier) getResponsesValueByEval(responses []ModifierResponse, eval string) any {
	x, err := jp.ParseString(eval)
	if helper.IsNil(err) {
		responsesEval := m.convertResponsesToEval(responses)
		resultJsonPath := x.Get(responsesEval)
		if helper.IsNotEmpty(resultJsonPath) {
			return resultJsonPath[0]
		}
	}
	return nil
}

func (m modifier) convertRequestsToEval(requests []ModifierRequest) []ModifierRequest {
	var requestsEval []ModifierRequest
	for _, request := range requests {
		request.Body = m.convertBodyToEval(request.Body)
		requestsEval = append(requestsEval, request)
	}
	return requestsEval
}

func (m modifier) convertResponsesToEval(responses []ModifierResponse) []ModifierResponse {
	var responsesEval []ModifierResponse
	for _, response := range responses {
		response.Body = m.convertBodyToEval(response.Body)
		responsesEval = append(responsesEval, response)
	}
	return responsesEval
}

func (m modifier) convertBodyToEval(body any) any {
	if helper.IsNil(body) {
		return nil
	}
	var valueAny any
	if helper.IsJson(body) {
		helper.SimpleConvertToDest(body, &valueAny)
	} else {
		valueAny = body
	}
	return valueAny
}

func (m modifier) persistModifierOnBody(modifier valueobject.Modifier, body, bodyToModifier any) {
	//aqui modificamos o body real e mantemos a ordem
	if helper.IsStruct(body) {
		bodyOrderedMap := body.(*orderedmap.OrderedMap)
		bodyModified := bodyToModifier.(map[string]any)
		m.modifierBodyMapByModifiedMap(modifier.Key, bodyOrderedMap, bodyModified)
	} else if helper.IsSlice(body) {
		bodySliceOrderedMap := body.([]*orderedmap.OrderedMap)
		bodySliceModified := bodyToModifier.([]map[string]any)
		for i, itemBody := range bodySliceOrderedMap {
			for i2, itemBodyModifier := range bodySliceModified {
				if helper.Equals(i, i2) {
					m.modifierBodyMapByModifiedMap(modifier.Key, itemBody, itemBodyModifier)
					break
				}
			}
		}
	}
}

func (m modifier) modifierBodyMapByModifiedMap(
	modifierKey string,
	body *orderedmap.OrderedMap,
	mMap map[string]any,
) {
	// alteramos oq ja tem
	for _, k := range body.Keys() {
		v, ok := mMap[k]
		if !ok {
			v, ok = mMap[modifierKey]
			if ok {
				k = modifierKey
				body.Delete(k)
			}
		}
		if ok {
			body.Set(k, v)
		} else {
			body.Delete(k)
		}
	}
	// adicionamos oq nao tem
	for k, v := range mMap {
		if _, ok := body.Get(k); !ok {
			body.Set(k, v)
		}
	}
}

func (m modifier) getCurrentRequest(requests []ModifierRequest) ModifierRequest {
	return requests[len(requests)-1]
}

func (m modifier) getCurrentResponse(responses []ModifierResponse) ModifierResponse {
	return responses[len(responses)-1]
}
