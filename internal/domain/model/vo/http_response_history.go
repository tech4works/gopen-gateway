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
	"net/http"
)

// httpResponseHistory type represents a collection of `HttpBackendResponse` objects.
// It can be used to store and manipulate the history of HTTP responses from a backend.
// The `httpResponseHistory` type provides methods to perform various operations on the history,
// such as filtering, size calculation, checking for success, and getting the most common status code and header.
type httpResponseHistory []*HttpBackendResponse

// Size returns the number of elements in the httpResponseHistory.
func (h httpResponseHistory) Size() int {
	return len(h)
}

// Success returns a boolean value indicating whether all backend responses in the httpResponseHistory are successful.
// It iterates through each backend httpResponse in the httpResponseHistory.
// If any backend httpResponse is not successful (Ok() returns false), it returns false.
// Otherwise, it returns true, indicating that all backend responses are successful.
func (h httpResponseHistory) Success() bool {
	for _, httpBackendResponse := range h {
		if !httpBackendResponse.Ok() {
			return false
		}
	}
	return true
}

// Filter applies a filter to the httpResponseHistory based on the provided httpRequest and httpResponse.
// It iterates through each httpBackendResponse in the httpResponseHistory and applies the config using the ApplyConfig method.
// If the ApplyConfig method returns a non-nil value, it appends the applied httpBackendResponse to the filteredHistory list.
// Returns the filteredHistory httpResponseHistory.
func (h httpResponseHistory) Filter(httpRequest *HttpRequest, httpResponse *HttpResponse) (
	filteredHistory httpResponseHistory) {
	for _, httpBackendResponse := range h {
		applied := httpBackendResponse.ApplyConfig(enum.BackendResponseApplyLate, httpRequest, httpResponse)
		if helper.IsNotNil(applied) {
			filteredHistory = append(filteredHistory, applied)
		}
	}
	return filteredHistory
}

// SingleResponse checks if the httpResponse history contains only one httpResponse.
// Returns true if the httpResponse history size is less than or equal to 1, false otherwise.
func (h httpResponseHistory) SingleResponse() bool {
	return helper.Equals(h.Size(), 1)
}

// MultipleResponse returns true if the size of the httpResponse history is greater than 1, indicating multiple responses.
// Otherwise, it returns false.
func (h httpResponseHistory) MultipleResponse() bool {
	return helper.IsGreaterThan(h.Size(), 1)
}

func (h httpResponseHistory) StatusCode() StatusCode {
	// se tiver mais de 1 resposta obtemos o código de status mais frequente
	if h.MultipleResponse() {
		return h.mostFrequentStatusCode()
	}
	return h.statusCode()
}

// Header aggregates the headers of all the HttpBackendResponse objects in the httpResponseHistory.
// It creates an empty Header object to store the aggregated headers.
// Then, it iterates through each HttpBackendResponse in the httpResponseHistory.
// For each non-empty header in a HttpBackendResponse, it adds the header to the aggregated Header object using the
// Aggregate method.
// Finally, it returns the aggregated Header object.
func (h httpResponseHistory) Header() Header {
	historyHeader := Header{}
	for _, httpBackendResponse := range h {
		if helper.IsNotNil(httpBackendResponse.Header()) && helper.IsNotEmpty(httpBackendResponse.Header()) {
			historyHeader = historyHeader.Aggregate(httpBackendResponse.Header())
		}
	}
	return historyHeader
}

// Body returns the Body object of the last httpBackendResponse in the httpResponseHistory list.
// If the httpResponseHistory contains multiple responses, it checks if aggregation is required.
// If aggregation is required, it returns the aggregated Body by key or the aggregated Body.
// If aggregation is not required, it returns a slice of Body objects from the backend responses.
// If the httpResponseHistory contains a single response, it returns the Body of that response.
// If the httpResponseHistory is empty, it returns nil.
func (h httpResponseHistory) Body(aggregate bool) *Body {
	if h.MultipleResponse() {
		return h.multipleBody(aggregate)
	}
	return h.body()
}

// Map iterates over each HttpBackendResponse in the httpResponseHistory and calls the Map method
// on each HttpBackendResponse to get a mapped representation of the response.
// The mapped representations are then added to a string slice called mappedHistory.
// Finally, the mappedHistory slice is converted to a string using the helper.SimpleConvertToString function
// and returned as the result of the Map method.
func (h httpResponseHistory) Map() string {
	var mappedHistory []any
	for _, httpBackendResponse := range h {
		mappedHistory = append(mappedHistory, httpBackendResponse.Map())
	}
	return helper.SimpleConvertToString(mappedHistory)
}

// last returns the last HttpBackendResponse in the httpResponseHistory list.
func (h httpResponseHistory) last() *HttpBackendResponse {
	return h[len(h)-1]
}

// statusCode returns the status code of the last HttpBackendResponse in the httpResponseHistory list.
// If the last response is nil, it returns http.StatusNoContent.
// Otherwise, it returns the status code of the last response.
//
// If there are multiple responses in the history, it returns the most frequent status code.
func (h httpResponseHistory) statusCode() StatusCode {
	if helper.IsNil(h.last()) {
		return http.StatusNoContent
	}
	return h.last().StatusCode()
}

// mostFrequentStatusCode calculates and returns the most frequent status code
// from the httpResponseHistory. It iterates through the history and counts
// the occurrences of each status code by using a map. Then, it finds the
// status code with the maximum count and returns it. If there are multiple
// status codes with the same maximum count, it returns the first one found.
func (h httpResponseHistory) mostFrequentStatusCode() StatusCode {
	statusCodes := make(map[StatusCode]int)
	for _, httpBackendResponse := range h {
		statusCodes[httpBackendResponse.StatusCode()]++
	}

	var mostFrequentCode StatusCode = http.StatusNoContent
	var maxCount int
	for code, count := range statusCodes {
		if count >= maxCount {
			mostFrequentCode = code
			maxCount = count
		}
	}

	return mostFrequentCode
}

// body returns the Body object of the last httpBackendResponse in the  httpResponseHistory list.
// If the httpResponseHistory contains multiple responses, it checks if aggregation is required.
// If aggregation is required, it returns the aggregated Body by key or the aggregated Body.
// If aggregation is not required, it returns a slice of Body objects from the backend responses.
// If the httpResponseHistory contains a single response, it returns the Body of that response.
// If the httpResponseHistory is empty, it returns nil.
func (h httpResponseHistory) body() *Body {
	if helper.IsNil(h.last()) {
		return nil
	}
	return h.last().Body()
}

// aggregateBody iterates over the httpResponseHistory and aggregates the bodies of the HTTP backend responses.
// If the body is nil, it is skipped. If the backend response is grouped by type, the body is aggregated using the key.
// Otherwise, the body is aggregated without using the key. The aggregated body is returned as a *Body.
// Note: The method assumes that all HTTP backend responses in the httpResponseHistory have valid bodies.
func (h httpResponseHistory) aggregateBody() *Body {
	historyBody := NewEmptyBodyJson()

	for index, httpBackendResponse := range h {
		if helper.IsNil(httpBackendResponse.Body()) {
			continue
		}
		body := httpBackendResponse.Body()
		if httpBackendResponse.GroupByType() {
			historyBody = historyBody.AggregateByKey(httpBackendResponse.Key(index), body)
		} else {
			historyBody = historyBody.Aggregate(body)
		}
	}

	return historyBody
}

// sliceOfBodies iterates through the httpResponseHistory and creates a new slice of Body objects.
// It skips over any HttpBackendResponse that has a nil body.
// For each non-nil body, it creates a new Body object using the NewBodyByHttpBackendResponse function.
// Finally, it returns a new Body object created by aggregating the slice of Body objects
// using the NewBodyBySlice function.
func (h httpResponseHistory) sliceOfBodies() *Body {
	var bodies []*Body

	for index, httpBackendResponse := range h {
		if helper.IsNil(httpBackendResponse.Body()) {
			continue
		}
		body := NewBodyByHttpBackendResponse(index, httpBackendResponse)
		bodies = append(bodies, body)
	}

	return NewBodyBySlice(bodies)
}

// multipleBody returns the response body of the httpResponseHistory. If aggregation is needed,
// it aggregates the body of all the responses using the aggregateBody method. Otherwise, it returns a
// slice of body objects from the backend responses using the sliceOfBodies method. If the httpResponseHistory
// is empty, it returns nil.
func (h httpResponseHistory) multipleBody(aggregate bool) *Body {
	if aggregate {
		return h.aggregateBody()
	}
	return h.sliceOfBodies()
}
