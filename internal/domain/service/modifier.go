package service

import (
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/enum"
	"github.com/GabrielHCataldo/martini-gateway/internal/domain/model/valueobject"
	"github.com/ohler55/ojg/jp"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Request struct {
	Host     string
	Endpoint string
	Url      string
	Method   string
	Header   http.Header
	Query    url.Values
	Params   map[string]string
	Body     any
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       any
	Group      string
	Hide       bool
}

type ExecuteRequestScopeInput struct {
	Request         Request
	Headers         []valueobject.Modifier
	Params          []valueobject.Modifier
	Queries         []valueobject.Modifier
	Body            []valueobject.Modifier
	RequestHistory  []Request
	ResponseHistory []Response
}

type ExecuteResponseScopeInput struct {
	Response        Response
	Headers         []valueobject.Modifier
	Params          []valueobject.Modifier
	Queries         []valueobject.Modifier
	Body            []valueobject.Modifier
	RequestHistory  []Request
	ResponseHistory []Response
}

type executeModifierInput struct {
	Scope     enum.ModifierScope
	Headers   []valueobject.Modifier
	Params    []valueobject.Modifier
	Queries   []valueobject.Modifier
	Body      []valueobject.Modifier
	Requests  []Request
	Responses []Response
}

type modifier struct {
}

type Modifier interface {
	ExecuteRequestScope(input ExecuteRequestScopeInput) (*Request, error)
	ExecuteResponseScope(input ExecuteResponseScopeInput) (*Response, error)
}

func NewModifier() Modifier {
	return modifier{}
}

func (m modifier) ExecuteRequestScope(input ExecuteRequestScopeInput) (*Request, error) {
	input.RequestHistory = append([]Request{input.Request}, input.RequestHistory...)
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

func (m modifier) ExecuteResponseScope(input ExecuteResponseScopeInput) (*Response, error) {
	input.ResponseHistory = append([]Response{input.Response}, input.ResponseHistory...)
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
		err := m.modifierParams(input.Scope, modifierParam, input.Requests, input.Responses)
		if helper.IsNotNil(err) {
			return err
		}
	}
	for _, modifierQuery := range input.Queries {
		err := m.modifierQueries(input.Scope, modifierQuery, input.Requests, input.Responses)
		if helper.IsNotNil(err) {
			return err
		}
	}
	for _, modifierBody := range input.Body {
		err := m.modifierBody(input.Scope, modifierBody, input.Requests, input.Responses)
		if helper.IsNotNil(err) {
			return err
		}
	}
	return nil
}

func (m modifier) modifierHeader(scope enum.ModifierScope, modifier valueobject.Modifier, requests []Request,
	responses []Response) {
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

func (m modifier) modifierParams(scope enum.ModifierScope, modifier valueobject.Modifier, requests []Request,
	responses []Response) error {
	if helper.NotContains(modifier.Scope, scope) && helper.NotContains(modifier.Scope, "*") {
		return nil
	} else if helper.IsNotEqualTo(scope, enum.ModifierScopeRequest) && helper.IsNotEqualTo(scope, "*") {
		return errors.New("params modifier not allow on scope:", scope)
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
	return nil
}

func (m modifier) modifierQueries(scope enum.ModifierScope, modifier valueobject.Modifier, requests []Request,
	responses []Response) error {
	if helper.NotContains(modifier.Scope, scope) && helper.NotContains(modifier.Scope, "*") {
		return nil
	} else if helper.IsNotEqualTo(scope, enum.ModifierScopeRequest) && helper.IsNotEqualTo(scope, "*") {
		return errors.New("queries modifier not allow on scope:" + scope)
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
	return nil
}

func (m modifier) modifierBody(scope enum.ModifierScope, modifier valueobject.Modifier, requests []Request,
	responses []Response) error {
	if helper.NotContains(modifier.Scope, scope) && helper.NotContains(modifier.Scope, "*") {
		return nil
	}
	var value any
	if helper.IsNotEqualTo(modifier.Action, enum.ModifierActionDel) {
		value = m.getValueEval(modifier.Value, requests, responses)
	}
	var bodyToModifier *any
	switch scope {
	case enum.ModifierScopeRequest:
		reqBody := m.getCurrentRequest(requests).Body
		bodyToModifier = &reqBody
		break
	case enum.ModifierScopeResponse:
		resBody := m.getCurrentResponse(responses).Body
		bodyToModifier = &resBody
		break
	}
	if helper.IsJson(bodyToModifier) {
		return m.modifierBodyJson(modifier, bodyToModifier, value)
	} else if helper.IsNotNil(bodyToModifier) {
		return m.modifierBodyString(modifier, bodyToModifier, value)
	}
	return nil
}

func (m modifier) modifierBodyJson(modifier valueobject.Modifier, body *any, value any) error {
	x, err := jp.ParseString(modifier.Key)
	if helper.IsNotNil(err) {
		return err
	}
	switch modifier.Action {
	case enum.ModifierActionSet, enum.ModifierActionAdd:
		_ = x.Set(body, value)
		break
	case enum.ModifierActionDel:
		_ = x.Del(body)
		break
	case enum.ModifierActionReplace:
		results := x.Get(body)
		if helper.IsNotEmpty(results) {
			y, err := jp.ParseString(modifier.Value)
			if helper.IsNil(err) {
				_ = x.Del(body)
				_ = y.Set(body, results[0])
			}
		}
		break
	}
	return nil
}

func (m modifier) modifierBodyString(modifier valueobject.Modifier, body *any, value any) error {
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
	*body = valueToReplace
	return nil
}

func (m modifier) getValueEval(valueModifier string, requests, responses any) any {
	regex := regexp.MustCompile(`\B\$[a-zA-Z0-9_.\[\]]+`)
	find := regex.FindAllString(valueModifier, -1)
	for _, word := range find {
		evalValue := strings.ReplaceAll(word, "$", "") //responses[0].body.token or //requests[0].body.auth.token
		splitDot := strings.Split(evalValue, ".")
		if helper.IsEmpty(splitDot) {
			continue
		}
		var valueGet any
		if helper.Contains(splitDot[0], "responses") {
			valueGet = m.getValueByStruct(responses, strings.Replace(evalValue, "responses", "", 1))
		} else if helper.Contains(splitDot[0], "requests") {
			valueGet = m.getValueByStruct(requests, strings.Replace(evalValue, "requests", "", 1))
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

func (m modifier) getValueByStruct(v any, eval string) any {
	x, err := jp.ParseString(eval)
	if helper.IsNil(err) {
		resultJsonPath := x.Get(v)
		if helper.IsNotEmpty(resultJsonPath) {
			return resultJsonPath[0]
		}
	}
	return nil
}

func (m modifier) getCurrentRequest(requests []Request) Request {
	return requests[len(requests)-1]
}

func (m modifier) getCurrentResponse(responses []Response) Response {
	return responses[len(responses)-1]
}
