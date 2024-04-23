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
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"net/http"
	"time"
)

// CacheResponse represents a cached HTTP response.
// It contains the status code, header, body, duration, and creation timestamp of the response.
// The duration specifies how long the response should be cached.
// The CreatedAt field indicates the timestamp of the response's creation.
type CacheResponse struct {
	// StatusCode is an integer field representing the status code of an HTTP response.
	// It is included in the CacheResponse struct and is used to store the status code of a cached response.
	StatusCode int `json:"statusCode"`
	// Header is a field representing the header of an HTTP response.
	// It is included in the CacheResponse struct and is used to store the header of a cached response.
	Header Header `json:"header"`
	// Body is a field representing the body of an HTTP response. (optional)
	// It is included in the CacheResponse struct and is used to store the body of a cached response.
	Body *CacheBody `json:"body,omitempty"`
	// Duration represents the duration for which the response should be cached.
	Duration string `json:"duration"`
	// CreatedAt is a field of the CacheResponse struct indicating the timestamp of the response's creation.
	CreatedAt time.Time `json:"createdAt"`
}

// Response represents the gateway HTTP response.
type Response struct {
	// endpoint represents a gateway endpoint.
	endpoint *Endpoint
	// statusCode stores the integer HTTP status code of the Response object.
	statusCode int
	// header represents the header of the Response object.
	header Header
	// Body represents the body of the gateway HTTP response.
	body *Body
	// abort bool is a flag in the Response object that indicates whether the response should be aborted.
	// If abort is set to true, it means that an error has occurred and the response should not be processed further.
	// The Abort method returns the value of the abort flag.
	// The AbortResponse method checks if the abort flag is set to true.
	abort bool
	// history represents the history of backend responses in the Response object.
	history responseHistory
}

// responseHistory represents the history of backend responses.
// It is a slice of backendResponse, which represents the responses from a backend service.
// The response history can be filtered and modified based on certain conditions.
// It also provides methods to retrieve information about the response history, such as size, success, and status code.
type responseHistory []*backendResponse

// errorResponseBody represents the structure of a response body containing error details.
// It is used to serialize the error details to JSON format.
//
// This struct is typically used in conjunction with the newErrorBody function to generate a JSON response body
// with error details based on the Endpoint and error provided.
type errorResponseBody struct {
	// File represents the file name or path where the error occurred.
	File string `json:"file"`
	// Line represents the line number where the error occurred
	Line int `json:"line"`
	// Endpoint represents the endpoint path where the error occurred.
	Endpoint string `json:"endpoint"`
	// Message represents the error message.
	Message string `json:"message"`
	// Timestamp represents the timestamp when the error occurred.
	Timestamp time.Time `json:"timestamp"`
}

// NewResponse creates a new Response object with the given Endpoint.
// The StatusCode is set to http.StatusNoContent.
// Returns the newly created Response object.
func NewResponse(endpointVO *Endpoint) *Response {
	return &Response{
		endpoint:   endpointVO,
		statusCode: http.StatusNoContent,
	}
}

// NewResponseByCache creates a new Response object with the given endpoint and cache response.
// The header of the cache response is modified to include the XGopenCache and XGopenCacheTTL headers.
// Returns the newly created Response object.
func NewResponseByCache(endpointVO *Endpoint, cacheResponseVO *CacheResponse) *Response {
	header := cacheResponseVO.Header
	header = header.Set(consts.XGopenCache, helper.SimpleConvertToString(true))
	header = header.Set(consts.XGopenCacheTTL, cacheResponseVO.TTL())
	return &Response{
		endpoint:   endpointVO,
		statusCode: cacheResponseVO.StatusCode,
		header:     header,
		body:       newBodyFromCacheBody(cacheResponseVO.Body),
	}
}

// NewResponseByErr creates a new Response object with the given Endpoint, statusCode, and err.
// It sets the header to newHeaderFailed() and the body to newErrorBody(endpointVO.path, err).
// Returns a pointer to the newly created Response object.
func NewResponseByErr(endpointVO *Endpoint, statusCode int, err error) *Response {
	return &Response{
		endpoint:   endpointVO,
		statusCode: statusCode,
		header:     newHeaderFailed(),
		abort:      true,
		body:       newErrorBody(endpointVO.path, err),
	}
}

// NewCacheResponse takes in a Response and a duration, then returns a new CacheResponse.
// The CacheResponse contains the StatusCode, Header, Body from the original Response,
// as well as the duration for which the response should be cached.
// The CreatedAt time is set to the current time.
//
// Parameters:
//   - responseVO : The original Response to be cached.
//   - duration : Duration for which the Response should be cached.
//
// Returns:
//   - A new CacheResponse containing the provided data and the current time of creation.
func NewCacheResponse(responseVO *Response, duration time.Duration) *CacheResponse {
	return &CacheResponse{
		StatusCode: responseVO.StatusCode(),
		Header:     responseVO.Header(),
		Body:       newCacheBody(responseVO.Body()),
		Duration:   duration.String(),
		CreatedAt:  time.Now(),
	}
}

// ModifyLastBackendResponse modifies the last backendResponse in the history list of the Response object.
// Replaces the last backendResponse with the provided backendResponseVO.
// Returns the modified Response object with the updated history.
// Does not modify other properties of the Response object.
func (r *Response) ModifyLastBackendResponse(backendResponseVO *backendResponse) *Response {
	history := r.history
	history[len(history)-1] = backendResponseVO

	// atualizamos os dados a partir do histórico alterado
	return r.notifyDataChanged(history)
}

// Append appends the backendResponseVO to the history list of the Response object.
// Returns the modified Response object with updated history.
// Does not modify other properties of the Response object.
func (r *Response) Append(backendResponseVO *backendResponse) *Response {
	// verificamos se o backendResponse precisa ser abortado
	if r.endpoint.AbortSequencial(backendResponseVO.StatusCode()) {
		return r.AbortResponse(backendResponseVO)
	}

	// adicionamos na nova lista de histórico
	history := r.history
	history = append(history, backendResponseVO)

	// atualizamos os dados a partir do histórico alterado
	return r.notifyDataChanged(history)
}

// Error constructs a standard gateway error response based on the received error.
// It builds the response's status code from the received error.
// If the error contains mapper.ErrBadGateway, the status code is set to http.StatusBadGateway.
// If the error contains mapper.ErrGatewayTimeout, the status code is set to http.StatusGatewayTimeout.
// Otherwise, the status code is set to http.StatusInternalServerError.
// It constructs the default gateway error response by setting the status code, header, body, and abort properties.
// Returns the constructed Response object representing the gateway error response.
func (r *Response) Error(path string, err error) *Response {
	// construímos o statusCode de resposta a partir do erro recebido
	var statusCode int
	if errors.Contains(err, mapper.ErrBadGateway) {
		statusCode = http.StatusBadGateway
	} else if errors.Contains(err, mapper.ErrGatewayTimeout) {
		statusCode = http.StatusGatewayTimeout
	} else {
		statusCode = http.StatusInternalServerError
	}
	// construímos a resposta de erro padrão do gateway
	return &Response{
		endpoint:   r.endpoint,
		statusCode: statusCode,
		header:     newHeaderFailed(),
		body:       newErrorBody(path, err),
		abort:      true,
	}
}

// Abort returns the value of the `abort` property of the Response object.
// If `abort` is true, it indicates that the response should be aborted.
// Returns a boolean value representing the `abort` property.
func (r *Response) Abort() bool {
	return r.abort
}

// AbortResponse creates a new Response object with the provided backendResponseVO.
// It sets the abort flag to true, indicating that the response should be aborted.
// The method combines the headers of the Response object and the backendResponseVO to create a new header.
// The created header is set in the new Response object along with the endpoint, status code, body,
// and history from the original Response object.
// Returns the modified Response object representing the aborted response.
func (r *Response) AbortResponse(backendResponseVO *backendResponse) *Response {
	// instanciamos o novo header
	header := newResponseHeader(false, false)
	// agregamos o header do backend abortado
	header = header.Aggregate(backendResponseVO.Header())

	// construímos o response com os dados do backend abortado
	return &Response{
		endpoint:   r.endpoint,
		statusCode: backendResponseVO.statusCode,
		header:     header,
		body:       backendResponseVO.body,
		abort:      true,
		history:    r.history,
	}
}

// LastBackendResponse returns the last backendResponse in the history list of the Response object.
// If the history is empty, it returns nil.
// Returns the last backendResponse from the history list.
func (r *Response) LastBackendResponse() *backendResponse {
	if helper.IsEmpty(r.history) {
		return nil
	}
	return r.history[len(r.history)-1]
}

// StatusCode returns the status code of the Response object.
func (r *Response) StatusCode() int {
	return r.statusCode
}

// Header returns the header of the Response object.
func (r *Response) Header() Header {
	return r.header
}

func (r *Response) ContentType() enum.ContentType {
	responseEncode := r.endpoint.ResponseEncode()
	if responseEncode.IsEnumValid() {
		return responseEncode.ContentType()
	} else if helper.IsNotNil(r.Body()) {
		return r.Body().ContentType()
	}
	return ""
}

// Body returns the body of the Response object.
func (r *Response) Body() *Body {
	return r.body
}

// BytesBody returns the body of the gateway HTTP response as a byte slice.
// It checks the response encoding and returns the body bytes based on the content type.
// If the response encoding is valid, it returns the body bytes by content type.
// If the response encoding is not valid and the body is not nil, it returns the body bytes.
// If the body is nil, it returns nil.
func (r *Response) BytesBody() []byte {
	// instanciamos o response encode do endpoint
	responseEncode := r.endpoint.ResponseEncode()
	// retornamos pelo responseEncode caso ele seja valido
	if responseEncode.IsEnumValid() {
		return r.body.BytesByContentType(responseEncode.ContentType())
	} else if helper.IsNotNil(r.Body()) {
		// se não tiver nil respondemos pelo tipo do body
		return r.body.Bytes()
	}
	// se nao tiver respondemos nil
	return nil
}

// Eval converts the Response object to a string representation.
// It creates a map of the Response object's properties, including:
// - statusCode: the integer HTTP status code
// - header: the header of the Response object
// - body: the body of the Response object as an interface{}
// - history: the history of backend responses as a string representation
// It then uses the helper.SimpleConvertToString function to convert the map to a string.
// Returns the string representation of the Response object.
func (r *Response) Eval() string {
	mapEval := map[string]any{
		"statusCode": r.statusCode,
		"header":     r.header,
		"body":       r.body.Interface(),
		"history":    r.history.Eval(),
	}
	return helper.SimpleConvertToString(mapEval)
}

// notifyDataChanged updates the Response object based on the modified history.
// It checks if the endpoint has reached its completion and returns the completion status.
// Then it filters the history based on the completion status to determine if it was a success or not.
// The method retrieves the status code from the filtered history and creates a new Response header
// by aggregating the completion and success values.
// It also aggregates the modifyHeaders from the filtered history.
// The method retrieves the body from the filtered history based on the endpoint's aggregateResponses flag.
// Finally, it constructs and returns a new Response object with the updated values.
func (r *Response) notifyDataChanged(history responseHistory) *Response {
	// checamos se ele chegou ao final para o valor padrão
	completed := r.endpoint.Completed(history.Size())

	// filtramos o histórico, com ele recebemos se sucesso ou não o histórico
	filteredHistory := history.Filter(completed)

	// obtemos o status code a partir do histórico filtrado
	statusCodeByHistory := filteredHistory.StatusCode()

	// criamos o header a partir dos valores complete e success
	header := newResponseHeader(completed, filteredHistory.Success())
	// agregamos os headers do histórico filtrado
	header = header.Aggregate(filteredHistory.Header())

	// obtemos o body a partir do histórico filtrado
	bodyByHistory := filteredHistory.Body(r.endpoint.AggregateResponses())

	// construímos o novo objeto de valor
	return &Response{
		endpoint:   r.endpoint,
		statusCode: statusCodeByHistory,
		header:     header,
		body:       bodyByHistory,
		history:    history,
	}
}

// TTL calculates the time to live (TTL) for the CacheResponse object.
// It subtracts the current time from the sum of the CreatedAt time and the Duration of the CacheResponse.
// Returns the TTL duration as a string representation.
func (c CacheResponse) TTL() string {
	duration, _ := time.ParseDuration(c.Duration)
	sub := c.CreatedAt.Add(duration).Sub(time.Now())
	return sub.String()
}

// Size returns the number of elements in the responseHistory.
func (r responseHistory) Size() int {
	return len(r)
}

// Success returns a boolean value indicating whether all backend responses in the response history were successful.
// It checks if the status code of each backend response is greater than or equal to http.StatusBadRequest.
// Returns true if all backend responses are successful, otherwise returns false.
func (r responseHistory) Success() bool {
	for _, backendResponseVO := range r {
		if helper.IsGreaterThanOrEqual(backendResponseVO.statusCode, http.StatusBadRequest) {
			return false
		}
	}
	return true
}

// Filter filters the response history based on the completed flag.
// It iterates over the history and filters out the responses that need to be omitted and have passed through all backends.
// The filtered responses are then added to the filtered history.
// Returns the filtered response history.
func (r responseHistory) Filter(completed bool) (filteredHistory responseHistory) {
	// iteramos o histórico para ser filtrado
	for _, backendResponseVO := range r {
		if backendResponseVO.omit && completed {
			//se a resposta do histórico quer ser omitida, e passou por todos os backends, pulamos ela
			continue
		}
		// setamos a resposta filtrada
		filteredHistory = append(filteredHistory, backendResponseVO)
	}
	return filteredHistory
}

// SingleResponse checks if the response history contains only one response.
// Returns true if the response history size is less than or equal to 1, false otherwise.
func (r responseHistory) SingleResponse() bool {
	return helper.Equals(r.Size(), 1)
}

// MultipleResponse returns true if the size of the response history is greater than 1, indicating multiple responses.
// Otherwise, it returns false.
func (r responseHistory) MultipleResponse() bool {
	return helper.IsGreaterThan(r.Size(), 1)
}

// StatusCode function provides HTTP Status code based on the response history.
// It checks if there are multiple or single responses, and returns the relevant HTTP status code.
// If there is more than one response, it returns HTTP status code 200 (OK).
// If there is a single response, it returns the HTTP status code of that particular response.
// If there are no responses, it returns HTTP status code 204 (No Content).
func (r responseHistory) StatusCode() int {
	// se tiver mais de 1 resposta obtemos o código de status mais frequente
	if r.MultipleResponse() {
		return r.mostFrequentStatusCode()
	} else if r.SingleResponse() {
		return r[0].statusCode
	}
	// resposta padrão de sucesso
	return http.StatusNoContent
}

// Header iterates over the responseHistory and aggregates non-nil headers from each backendResponseVO.
// It constructs a final Header value object, which is an aggregation of all the individual non-nil headers.
// The method returns this final aggregated header.
//
// This function wouldn't consider any backendResponseVO whose header is nil.
// Also, with every aggregation, a new Header value object is created.
func (r responseHistory) Header() Header {
	h := Header{}
	for _, backendResponseVO := range r {
		// se tiver nil pulamos para o próximo
		if helper.IsNil(backendResponseVO.header) {
			continue
		}
		// agregamos os valores ao header de resultado, gerando um novo objeto de valor a cada agregação
		h = h.Aggregate(backendResponseVO.header)
	}
	return h
}

// Body function takes 'aggregateResponses' boolean as argument.
// If there are multiple responses present, it aggregates them based on
// the boolean parameter. If the boolean parameter is true, it returns the result
// of 'aggregateBody' function else 'aggregatedBodies' function is called.
// For a single response case, the body of the single response is returned.
// If no responses are available, returns an empty Body struct.
func (r responseHistory) Body(aggregateResponses bool) *Body {
	if r.MultipleResponse() {
		if aggregateResponses {
			return r.aggregateBody()
		} else {
			return r.sliceOfBodies()
		}
	} else if r.SingleResponse() {
		return r.body()
	}
	return nil
}

// Eval iterates over the responseHistory and calls the Eval method on each backendResponseVO.
// It creates a new list of the return values from each Eval call and returns it.
// Returns a new list []any containing the results of the Eval calls on each backendResponseVO.
func (r responseHistory) Eval() []any {
	var evalHistory []any
	for _, backendResponseVO := range r {
		evalHistory = append(evalHistory, backendResponseVO.Eval())
	}
	return evalHistory
}

// last returns the last backendResponse in the responseHistory list.
// Returns the last backendResponse object.
// Does not modify the responseHistory.
func (r responseHistory) last() *backendResponse {
	return r[len(r)-1]
}

// body returns the Body object of the last backendResponse in the responseHistory list.
// Creates a new Body object using the last backendResponse in the responseHistory.
// Returns the newly created Body object.
func (r responseHistory) body() *Body {
	return newBodyFromBackendResponse(r.last())
}

// aggregateBody aggregates the body from each backend response in the response history.
// It creates an initial Body object with empty JSON content.
// Then, it iterates through each backend response,
// skipping responses with a nil body.
// If the backend response has a group response flag set to true,
// it aggregates the body by key using the AggregateByKey method of the Body object.
// Otherwise, it aggregates all the JSON fields into the body using the Aggregate method of the Body object.
// Returns the final aggregated body.
// If there are no non-nil bodies in the response history, returns an empty Body struct.
func (r responseHistory) aggregateBody() *Body {
	// instanciamos primeiro o aggregate body para retornar
	bodyHistory := &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer("{}"),
	}

	// iteramos o histórico de backends response
	for index, backendResponseVO := range r {
		// se tiver nil pulamos para o próximo
		if helper.IsNil(backendResponseVO.body) {
			continue
		}
		// instânciamos o body
		body := backendResponseVO.Body()
		// caso seja string ou slice agregamos na chave, caso contrario, iremos agregar todos os campos json no bodyHistory
		if backendResponseVO.GroupResponse() {
			bodyHistory = bodyHistory.AggregateByKey(backendResponseVO.Key(index), body)
		} else {
			bodyHistory = bodyHistory.Aggregate(body)
		}
	}

	// se tudo ocorreu bem retornamos o corpo agregado
	return bodyHistory
}

// sliceOfBodies iterates over the response history and constructs a slice of bodies from the backend responses.
// If a backend response has an empty body, it is skipped.
// For each backend response with a non-empty body, a bodyBackendResponse object is created and added to the list of bodies.
// Returns a new Body object that contains the aggregated list of bodies from the response history.
func (r responseHistory) sliceOfBodies() *Body {
	// instanciamos o valor a ser construído
	var bodies []*Body
	// iteramos o histórico para listar os bodies de resposta
	for index, backendResponseVO := range r {
		// se tiver vazio pulamos para o próximo
		if helper.IsNil(backendResponseVO.body) {
			continue
		}
		// inicializamos o body agregado
		bodyBackendResponse := newBodyFromIndexAndBackendResponse(index, backendResponseVO)
		// inserimos na lista de retorno
		bodies = append(bodies, bodyBackendResponse)
	}
	// se tudo ocorrer bem, teremos o body agregado
	return newSliceBody(bodies)
}

func (r responseHistory) mostFrequentStatusCode() int {
	statusCodes := make(map[int]int)
	for _, backendResponseVO := range r {
		statusCodes[backendResponseVO.statusCode]++
	}

	maxCount := 0
	mostFrequentCode := 0
	for code, count := range statusCodes {
		if count >= maxCount {
			mostFrequentCode = code
			maxCount = count
		}
	}

	return mostFrequentCode
}
