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
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"strings"
)

// Params is a type alias for a map of string key-value pairs that represents additional parameters for a URL path.
type Params = map[string]string

// UrlPath is a struct that represents a URL path with additional parameter mappings.
// It consists of a value, which is the actual URL path, and params, which is a map
// of string key-value pairs representing additional parameters for the URL path.
type UrlPath struct {
	// value is a string field that represents the actual URL path of a UrlPath object.
	// It is used to store the URL path value and to perform modifications and
	// parameter replacements on the URL path.
	value string
	// Params is a type alias for a map of string key-value pairs that represents additional parameters for a URL path.
	params Params
}

// NewUrlPath creates a new instance of UrlPath.
// It initializes the UrlPath with the specified value and params.
// The function filters the params that are present in the value and returns
// a new UrlPath object with the filtered params.
func NewUrlPath(value string, params Params) UrlPath {
	filteredParams := Params{}
	for k, v := range params {
		paramPathKey := patternParamPathKey(k)
		if helper.Contains(value, paramPathKey) {
			filteredParams[k] = v
		}
	}
	return UrlPath{
		value:  value,
		params: filteredParams,
	}
}

// Params returns the `Params` field of the `UrlPath` instance.
// `Params` represents a map of string key-value pairs that store additional parameters for the URL path.
// The `Params` method allows you to access and manipulate these parameters.
// The returned `Params` field is of type `map[string]string`.
func (u UrlPath) Params() Params {
	return u.params
}

// String returns the string representation of the UrlPath instance.
// It replaces all parameter placeholders in the value with their corresponding parameter values.
// The parameter placeholders are identified by a leading colon (":").
// The returned string is the modified value with the parameter placeholders replaced.
func (u UrlPath) String() string {
	urlPathString := u.value
	for k, v := range u.params {
		paramPathKey := patternParamPathKey(k)
		urlPathString = strings.ReplaceAll(urlPathString, paramPathKey, v)
	}
	return urlPathString
}

// SetParam sets a parameter with the given key and value in the UrlPath.
// It creates a copy of the current parameter map, adds the new key-value pair,
// and updates the URL path value if the key is not already present.
// The method returns a new UrlPath instance with the updated parameters and URL path value.
func (u UrlPath) SetParam(key, value string) UrlPath {
	urlPathParams := u.copyParams()
	urlPathParams[key] = value

	urlPathValue := u.value
	paramPathKey := patternParamPathKey(key)
	if helper.NotContains(urlPathValue, paramPathKey) {
		urlPathValue = fmt.Sprintf("%s/%s", urlPathValue, paramPathKey)
	}

	return UrlPath{
		value:  urlPathValue,
		params: urlPathParams,
	}
}

// ReplaceParam replaces the value of a parameter with the given key in the UrlPath.
// If the parameter does not exist, it returns the current UrlPath unchanged.
// It calls the SetParam method to update the parameter value and returns a new UrlPath instance.
func (u UrlPath) ReplaceParam(key, value string) UrlPath {
	if u.NotExistsParam(key) {
		return u
	}
	return u.SetParam(key, value)
}

// DeleteParam deletes the parameter with the given key from the UrlPath instance.
// It creates a copy of the current parameter map, removes the key-value pair,
// and updates the URL path value by replacing the corresponding parameter placeholder with an empty string.
// The method returns a new UrlPath instance with the updated parameters and URL path value.
func (u UrlPath) DeleteParam(key string) UrlPath {
	urlPathParams := u.copyParams()
	delete(urlPathParams, key)

	urlPathValue := strings.ReplaceAll(u.value, patternParamUrlKey(key), "")

	return UrlPath{
		value:  urlPathValue,
		params: urlPathParams,
	}
}

// ExistsParam checks if a parameter with the given key exists in the UrlPath instance.
// It returns true if the key exists in the parameters map and if the parameter placeholder
// exists in the UrlPath value. Otherwise, it returns false.
func (u UrlPath) ExistsParam(key string) bool {
	_, ok := u.params[key]
	return ok && helper.Contains(u.value, patternParamPathKey(key))
}

// NotExistsParam checks if a parameter with the given key does not exist in the UrlPath instance.
// It negates the result of calling the ExistsParam method and returns true if the parameter does not exist.
func (u UrlPath) NotExistsParam(key string) bool {
	return !u.ExistsParam(key)
}

// Modify modifies the UrlPath based on the provided Modifier, HttpRequest, and HttpResponse.
// It returns a new modified UrlPath object based on the action specified in the Modifier.
// If the action is set, it sets the parameter value in the UrlPath.
// If the action is replace, it replaces the parameter value in the UrlPath.
// If the action is delete, it deletes the parameter from the UrlPath.
// If the action is not recognized, it returns the original UrlPath unchanged.
func (u UrlPath) Modify(modifier *Modifier, httpRequest *HttpRequest, httpResponse *HttpResponse) UrlPath {
	newValue := modifier.ValueAsString(httpRequest, httpResponse)
	switch modifier.Action() {
	case enum.ModifierActionSet:
		return u.SetParam(modifier.Key(), newValue)
	case enum.ModifierActionRpl:
		return u.ReplaceParam(modifier.Key(), newValue)
	case enum.ModifierActionDel:
		return u.DeleteParam(modifier.Key())
	default:
		return u
	}
}

// copyParams creates a deep copy of the current parameter map in UrlPath.
// It iterates over the key-value pairs in the parameter map and creates a new map
// with the same key-value pairs.
// The method returns the copied parameter map.
func (u UrlPath) copyParams() map[string]string {
	copiedParams := make(map[string]string)
	for k, v := range u.params {
		copiedParams[k] = v
	}
	return copiedParams
}

// patternParamPathKey returns a string with a leading colon (":") followed by the given key.
func patternParamPathKey(k string) string {
	return fmt.Sprintf(":%s", k)
}

// patternParamUrlKey takes a string parameter `k` and prepends it with "/:".
// It returns the formatted string with the parameter key prefixed with "/:".
func patternParamUrlKey(k string) string {
	return fmt.Sprintf("/:%s", k)
}
