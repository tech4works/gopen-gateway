package factory

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/dto"
	"github.com/ohler55/ojg/jp"
	"net/http"
	"regexp"
	"slices"
	"strings"
)

type modifier struct {
}

type Modifier interface {
	ExecuteModifier(
		scope string,
		backend dto.Backend,
		requests []dto.BackendRequest,
		responses []dto.BackendResponse,
	) error
}

func NewModifier() Modifier {
	return modifier{}
}

func (m modifier) ExecuteModifier(
	scope string,
	backend dto.Backend,
	requests []dto.BackendRequest,
	responses []dto.BackendResponse,
) error {
	for _, modifierHeader := range backend.Headers {
		err := m.modifierHeader(scope, modifierHeader, requests, responses)
		if err != nil {
			return err
		}
	}
	for _, modifierParam := range backend.Params {
		err := m.modifierParams(scope, modifierParam, requests, responses)
		if err != nil {
			return err
		}
	}
	for _, modifierQuery := range backend.Queries {
		err := m.modifierQueries(scope, modifierQuery, requests, responses)
		if err != nil {
			return err
		}
	}
	for _, modifierBody := range backend.Body {
		err := m.modifierBody(scope, modifierBody, requests, responses)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m modifier) modifierHeader(
	scope string,
	modifier dto.Modifier,
	requests []dto.BackendRequest,
	responses []dto.BackendResponse,
) error {
	if slices.Contains(modifier.Scope, scope) || slices.Contains(modifier.Scope, "*") {
		var value string
		if modifier.Action != "del" {
			valueEval := m.getValueEval(modifier.Value, requests, responses)
			value = fmt.Sprintf("%v", valueEval)
		}
		var header http.Header
		if scope == "request" {
			header = requests[len(requests)-1].Header
		} else if scope == "response" {
			header = responses[len(responses)-1].Header
		} else {
			return errors.New("unknown scope:" + scope)
		}
		if modifier.Action == "set" {
			header.Set(modifier.Key, value)
		} else if modifier.Action == "add" {
			header.Add(modifier.Key, value)
		} else if modifier.Action == "del" {
			header.Del(modifier.Key)
		}
	}
	return nil
}

func (m modifier) modifierParams(
	scope string,
	modifier dto.Modifier,
	requests []dto.BackendRequest,
	responses []dto.BackendResponse,
) error {
	if slices.Contains(modifier.Scope, scope) || slices.Contains(modifier.Scope, "*") {
		var value string
		if modifier.Action != "del" {
			valueEval := m.getValueEval(modifier.Value, requests, responses)
			value = fmt.Sprintf("%v", valueEval)
		}
		currentRequest := requests[len(requests)-1]
		if scope != "request" && scope != "*" {
			return errors.New("params modifier not allow on scope:" + scope)
		}
		if modifier.Action == "add" || modifier.Action == "set" {
			currentRequest.Params[modifier.Key] = value
		} else if modifier.Action == "del" {
			delete(currentRequest.Params, modifier.Key)
			currentRequest.Endpoint = strings.Replace(currentRequest.Endpoint, "/:"+modifier.Key, "", 1)
		}
	}
	return nil
}

func (m modifier) modifierQueries(
	scope string,
	modifier dto.Modifier,
	requests []dto.BackendRequest,
	responses []dto.BackendResponse,
) error {
	if slices.Contains(modifier.Scope, scope) || slices.Contains(modifier.Scope, "*") {
		var value string
		if modifier.Action != "del" {
			valueEval := m.getValueEval(modifier.Value, requests, responses)
			value = fmt.Sprintf("%v", valueEval)
		}
		currentRequest := requests[len(requests)-1]
		if scope != "request" && scope != "*" {
			return errors.New("queries modifier not allow on scope:" + scope)
		}
		if modifier.Action == "set" {
			currentRequest.Query.Set(modifier.Key, value)
		} else if modifier.Action == "add" {
			currentRequest.Query.Add(modifier.Key, value)
		} else if modifier.Action == "del" {
			currentRequest.Query.Del(modifier.Key)
		}
	}
	return nil
}

func (m modifier) modifierBody(
	scope string,
	modifier dto.Modifier,
	requests []dto.BackendRequest,
	responses []dto.BackendResponse,
) error {
	if slices.Contains(modifier.Scope, scope) || slices.Contains(modifier.Scope, "*") {
		var value any
		if modifier.Action != "del" {
			value = m.getValueEval(modifier.Value, requests, responses)
		}
		var bodyToModifier *any
		if scope == "request" {
			bodyToModifier = &requests[len(requests)-1].Body
		} else if scope == "response" {
			bodyToModifier = &responses[len(responses)-1].Body
		} else {
			return errors.New("unknown scope:" + scope)
		}
		x, err := jp.ParseString(modifier.Key)
		if err != nil {
			return err
		}
		if helper.IsJson(bodyToModifier) {
			if modifier.Action == "set" || modifier.Action == "add" {
				_ = x.Set(bodyToModifier, value)
			} else if modifier.Action == "del" {
				_ = x.Del(bodyToModifier)
			} else if modifier.Action == "replace" {
				//{"ID":"1"} to {"teste": 1}
				results := x.Get(bodyToModifier)
				if len(results) > 0 {
					y, err := jp.ParseString(modifier.Value)
					if err == nil {
						_ = x.Del(bodyToModifier)
						_ = y.Set(bodyToModifier, results[0])
					}
				}
			}
		} else if bodyToModifier != nil {
			modifierValue := fmt.Sprintf("%v", value)
			valueToReplace := fmt.Sprintf("%v", *bodyToModifier)
			if modifier.Action == "add" {
				valueToReplace = valueToReplace + modifierValue
			} else if modifier.Action == "set" {
				valueToReplace = strings.ReplaceAll(valueToReplace, modifier.Key, modifierValue)
			} else if modifier.Action == "del" {
				valueToReplace = strings.ReplaceAll(valueToReplace, modifier.Key, "")
			} else if modifier.Action == "replace" {
				valueToReplace = modifierValue
			}
			*bodyToModifier = valueToReplace
		}
	}
	return nil
}

func (m modifier) getValueEval(valueModifier string, requests []dto.BackendRequest, responses []dto.BackendResponse) any {
	regex := regexp.MustCompile(`\B\$[a-zA-Z0-9_.\[\]]+`)
	resultFind := regex.FindAllString(valueModifier, -1)
	for _, word := range resultFind {
		evalValue := strings.ReplaceAll(word, "$", "") //responses[0].body.token or //requests[0].body.auth.token
		splitDot := strings.Split(evalValue, ".")
		if len(splitDot) > 0 {
			var valueGet any
			if strings.Contains(splitDot[0], "responses") {
				valueGet = m.getValueByObject(responses, strings.Replace(evalValue, "responses", "", 1))
			} else if strings.Contains(splitDot[0], "requests") {
				valueGet = m.getValueByObject(requests, strings.Replace(evalValue, "requests", "", 1))
			}
			if valueGet != nil {
				if word == valueModifier {
					return valueGet
				} else {
					if helper.IsJson(valueGet) {
						bytesJSON, err := json.Marshal(valueGet)
						if err == nil {
							valueGet = string(bytesJSON)
						}
					}
					valueGetString := fmt.Sprintf("%v", valueGet)
					valueModifier = strings.Replace(valueModifier, word, valueGetString, 1)
				}
			}
		}
	}
	valueModifierBytes := []byte(valueModifier)
	if json.Valid(valueModifierBytes) {
		var obj any
		err := json.Unmarshal(valueModifierBytes, &obj)
		if err == nil {
			return obj
		}
	}
	return valueModifier
}

func (m modifier) getValueByObject(obj any, eval string) any {
	x, err := jp.ParseString(eval)
	if err == nil {
		resultJsonPath := x.Get(obj)
		if len(resultJsonPath) > 0 {
			return resultJsonPath[0]
		}
	}
	return nil
}
