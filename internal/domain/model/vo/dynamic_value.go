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
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

// DynamicValue represents a dynamic value that can be processed and converted into different types.
type DynamicValue struct {
	value string
}

// NewDynamicValue creates a new DynamicValue object with the given value.
// The DynamicValue can be processed and converted into different types.
// It is used in the Modifier struct to store a string value for modifications.
func NewDynamicValue(value string) DynamicValue {
	return DynamicValue{
		value: value,
	}
}

// AsInt converts the DynamicValue object to an integer value.
// It takes an HttpRequest and HttpResponse object as arguments to resolve dynamic values.
// First, it converts the DynamicValue to a string using the AsString method.
// Then, it uses the SimpleConvertToInt function from the helper package to convert the string to an integer.
// Finally, it returns the converted integer value.
func (d DynamicValue) AsInt(httpRequest *HttpRequest, httpResponse *HttpResponse) int {
	dynamicValueStr := d.AsString(httpRequest, httpResponse)
	return helper.SimpleConvertToInt(dynamicValueStr)
}

// AsSliceOfString converts the DynamicValue object to a slice of strings.
// It takes an HttpRequest and HttpResponse object as arguments to resolve dynamic values.
// First, it converts the DynamicValue to a string using the AsString method.
// Then, it checks if the converted string is a slice type using the helper.IsSliceType function.
// If it is a slice type, it creates an empty string slice (ss) and uses the helper.SimpleConvertToDest
// function to convert the string to the slice.
// If the converted slice is not empty, it returns the slice.
// If the converted string is not a slice type, it returns a string slice with the converted string as its only element.
// If the DynamicValue object is empty, it returns an empty string slice.
func (d DynamicValue) AsSliceOfString(httpRequest *HttpRequest, httpResponse *HttpResponse) []string {
	dynamicValueStr := d.AsString(httpRequest, httpResponse)
	if helper.IsSliceType(dynamicValueStr) {
		var ss []string
		helper.SimpleConvertToDest(dynamicValueStr, &ss)
		if helper.IsNotEmpty(ss) {
			return ss
		}
	}
	return []string{dynamicValueStr}
}

// AsString converts the DynamicValue object to a string value.
// It takes an HttpRequest and HttpResponse object as arguments to resolve dynamic values.
// First, it initializes the value variable with the initial value of the DynamicValue object.
// Then, it iterates through the words identified by the findAllByDynamicValueSyntax method.
// For each word, it processes the dynamic value by calling the processDynamicValueWord method.
// The processDynamicValueWord method replaces the word with the corresponding dynamic value from the HttpRequest or HttpResponse.
// Finally, it returns the resulting string value.
func (d DynamicValue) AsString(httpRequest *HttpRequest, httpResponse *HttpResponse) string {
	value := d.value
	for _, word := range d.findAllByDynamicValueSyntax() {
		value = d.processDynamicValueWord(word, httpRequest, httpResponse)
	}
	return value
}

// findAllByDynamicValueSyntax finds all dynamic values in the given DynamicValue string.
// It uses a regular expression pattern to match the dynamic value syntax: \B#[a-zA-Z0-9_.\-\[\]]+.
// The method returns a slice of strings containing all dynamic values found in the string.
func (d DynamicValue) findAllByDynamicValueSyntax() []string {
	regex := regexp.MustCompile(`\B#[a-zA-Z0-9_.\-\[\]]+`)
	return regex.FindAllString(d.value, -1)
}

// processDynamicValueWord processes a dynamic value word in the given DynamicValue string.
// It takes a word, an HttpRequest object, and an HttpResponse object as arguments.
// First, it calls the getDynamicValuePerWord method to get the corresponding dynamic value for the word.
// If the dynamic value is empty, it returns the initial value of the DynamicValue object.
// Otherwise, it replaces the word with the dynamic value in the string using the strings.Replace method,
// and returns the resulting string value.
func (d DynamicValue) processDynamicValueWord(word string, httpRequest *HttpRequest, httpResponse *HttpResponse) string {
	dynamicValue := d.getDynamicValuePerWord(word, httpRequest, httpResponse)
	if helper.IsEmpty(dynamicValue) {
		return d.value
	}
	return strings.Replace(d.value, word, dynamicValue, 1)
}

// getDynamicValuePerWord returns the dynamic value for the given word.
// It takes a word, an HttpRequest object, and an HttpResponse object as arguments.
// First, it cleans the word by removing the '#' character using strings.ReplaceAll method.
// Then, it splits the cleaned word by the '.' character using strings.Split method.
// If the resulting slice is empty, it returns an empty string.
// If the first element of the split slice contains the substring "request", it calls the
// getHttpRequestValueByJsonPath method to get the dynamic value from the HttpRequest object.
// If the first element of the split slice contains the substring "response", it calls the
// getHttpResponseValueByJsonPath method to get the dynamic value from the HttpResponse object.
// If none of the above conditions are met, it returns an empty string.
func (d DynamicValue) getDynamicValuePerWord(word string, httpRequest *HttpRequest, httpResponse *HttpResponse) string {
	cleanSintaxe := strings.ReplaceAll(word, "#", "")
	dotSplit := strings.Split(cleanSintaxe, ".")
	if helper.IsEmpty(dotSplit) {
		return ""
	}

	if helper.Contains(dotSplit[0], "request") {
		return d.getHttpRequestValueByJsonPath(cleanSintaxe, httpRequest)
	} else if helper.Contains(dotSplit[0], "response") {
		return d.getHttpResponseValueByJsonPath(cleanSintaxe, httpResponse)
	}
	return ""
}

// getHttpRequestValueByJsonPath returns the value from the HttpRequest object based on the given JSON path.
// It takes a jsonPath string and an httpRequest object as arguments.
// First, it replaces the "request." substring in jsonPath with an empty string.
// Then, it uses the gjson.Get method to retrieve the value from the httpRequest.Map() using the jsonPath.
// Finally, it returns the retrieved value as a string.
func (d DynamicValue) getHttpRequestValueByJsonPath(jsonPath string, httpRequest *HttpRequest) string {
	jsonPath = strings.Replace(jsonPath, "request.", "", 1)
	result := gjson.Get(httpRequest.Map(), jsonPath)
	return result.String()
}

// getHttpResponseValueByJsonPath retrieves the value from the HttpResponse object based on the given JSON path.
// It takes a jsonPath string and an httpResponse object as arguments.
// First, it removes the "response." substring in jsonPath with an empty string.
// Then, it uses the gjson.Get method to retrieve the value from the httpResponse.Map() using the jsonPath.
// Finally, it returns the retrieved value as a string.
func (d DynamicValue) getHttpResponseValueByJsonPath(jsonPath string, httpResponse *HttpResponse) string {
	jsonPath = strings.Replace(jsonPath, "response.", "", 1)
	result := gjson.Get(httpResponse.Map(), jsonPath)
	return result.String()
}
