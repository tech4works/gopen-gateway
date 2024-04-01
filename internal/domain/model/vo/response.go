package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/iancoleman/orderedmap"
	"net/http"
	"time"
)

type responseHistory []backendResponse

type aggregateBody struct {
	value orderedmap.OrderedMap
}

type CacheResponse struct {
	StatusCode int           `json:"statusCode"`
	Header     Header        `json:"header"`
	Body       Body          `json:"body"`
	Duration   time.Duration `json:"duration"`
	CreatedAt  time.Time     `json:"createdAt"`
}

type Response struct {
	endpoint   Endpoint
	statusCode int
	header     Header
	body       Body
	abort      bool
	history    responseHistory
}

// NewResponse creates a new Response object with the given endpoint.
// The status code of the Response object is set to http.StatusNoContent.
// Returns the newly created Response object.
func NewResponse(endpointVO Endpoint) Response {
	return Response{
		endpoint:   endpointVO,
		statusCode: http.StatusNoContent,
	}
}

// NewResponseByCache creates a new Response object with the given endpoint and cache response.
// The header of the cache response is modified to include the XGOpenCache and XGOpenCacheTTL headers.
// Returns the newly created Response object.
func NewResponseByCache(endpointVO Endpoint, cacheResponseVO CacheResponse) Response {
	header := cacheResponseVO.Header
	header = header.Set(consts.XGOpenCache, helper.SimpleConvertToString(true))
	header = header.Set(consts.XGOpenCacheTTL, cacheResponseVO.TTL())
	return Response{
		endpoint:   endpointVO,
		statusCode: cacheResponseVO.StatusCode,
		header:     header,
		body:       cacheResponseVO.Body,
	}
}

// NewResponseByErr creates a new Response object with the given endpoint, status code, and error.
// The Response object has a header set to newHeaderFailed() and a body set to newBodyByErr(err).
// Returns the newly created Response object.
func NewResponseByErr(endpointVO Endpoint, statusCode int, err error) Response {
	return Response{
		endpoint:   endpointVO,
		statusCode: statusCode,
		header:     newHeaderFailed(),
		body:       newBodyByErr(err),
	}
}

// NewCacheResponse creates a new CacheResponse object with the provided writer and duration.
// The StatusCode of the CacheResponse object is set to the writer's status code.
// The Header of the CacheResponse object is initialized using NewHeader(writer.Header()).
// The Body of the CacheResponse object is created by parsing the writer's body bytes using the newBody function.
// The Duration of the CacheResponse object is set to the provided duration.
// The CreatedAt field of the CacheResponse object set to the current time.
// Returns the newly created CacheResponse object.
func NewCacheResponse(writer dto.Writer, duration time.Duration) CacheResponse {
	return CacheResponse{
		StatusCode: writer.Status(),
		Header:     NewHeader(writer.Header()),
		Body:       newBody(writer.Body.Bytes()),
		Duration:   duration,
		CreatedAt:  time.Now(),
	}
}

// newAggregateBody creates a new aggregateBody object with the given index and backendResponse.
// It initializes an orderedmap, sets the "ok" and "code" fields of the value to the corresponding values from backendResponseVO,
// and creates a new instance of aggregateBody with the value set.
// Then, it calls the aggregate method on the new instance to aggregate the data based on the index and backendResponseVO.
// Finally, it returns the newly created aggregateBody object.
func newAggregateBody(index int, backendResponseVO backendResponse) aggregateBody {
	value := orderedmap.New()
	value.Set("ok", backendResponseVO.Ok())
	value.Set("code", backendResponseVO.StatusCode())

	newInstance := aggregateBody{
		value: *value,
	}
	newInstance.aggregate(index, backendResponseVO)
	return newInstance
}

// ModifyStatusCode modifies the status code of the Response object.
// If the statusCode parameter is equal to the current status code, gives priority to the local value and updates the
// data by modified history. Otherwise, create a new Response object with the modified status code and the same values
// for other properties.
// Returns the modified Response object.
func (r Response) ModifyStatusCode(statusCode int, backendResponseVO backendResponse) Response {
	history := r.history
	history[len(history)-1] = backendResponseVO

	// se o valor vindo do parâmetro vir igual, damos prioridade ao local
	if helper.Equals(r.statusCode, statusCode) {
		// atualizamos os dados a partir do histórico alterado
		return r.notifyDataChanged(history)
	}

	return Response{
		endpoint:   r.endpoint,
		statusCode: statusCode,
		header:     r.header,
		body:       r.body,
		abort:      r.abort,
		history:    history,
	}
}

// ModifyHeader modifies the header of the Response object.
// If the header parameter is equal to the current header, gives priority to the local value and updates the
// data by modified history. Otherwise, create a new Response object with the modified header and the same values
// for other properties.
// Returns the modified Response object.
func (r Response) ModifyHeader(header Header, backendResponseVO backendResponse) Response {
	history := r.history
	history[len(history)-1] = backendResponseVO

	// se o valor vindo do parâmetro vir igual, damos prioridade ao local
	if helper.Equals(r.header, header) {
		// atualizamos os dados a partir do histórico alterado
		return r.notifyDataChanged(history)
	}

	return Response{
		endpoint:   r.endpoint,
		statusCode: r.statusCode,
		header:     header,
		body:       r.body,
		abort:      r.abort,
		history:    history,
	}
}

// ModifyBody modifies the body property of the Response object.
// If the body parameter is equal to the current body, gives priority to the local value and updates the
// data by modified history. Otherwise, create a new Response object with the modified body and the same values
// for other properties.
// Returns the modified Response object.
func (r Response) ModifyBody(body Body, backendResponseVO backendResponse) Response {
	history := r.history
	history[len(history)-1] = backendResponseVO

	// se o valor vindo do parâmetro vir igual, damos prioridade ao local
	if helper.Equals(r.body, body) {
		// atualizamos os dados a partir do histórico alterado
		return r.notifyDataChanged(history)
	}

	return Response{
		endpoint:   r.endpoint,
		statusCode: r.statusCode,
		header:     r.header,
		body:       body,
		abort:      r.abort,
		history:    history,
	}
}

// Append appends the backendResponseVO to the history list of the Response object.
// Returns the modified Response object with updated history.
// Does not modify other properties of the Response object.
func (r Response) Append(backendResponseVO backendResponse) Response {
	// adicionamos na nova lista de histórico
	history := r.history
	history = append(history, backendResponseVO)

	// atualizamos os dados a partir do histórico alterado
	return r.notifyDataChanged(history)
}

// Error modifies the Response object to represent an error.
// It builds the status code of the response based on the received error.
// If the error contains mapper.ErrBadGateway, sets the status code to http.StatusBadGateway.
// If the error contains mapper.ErrGatewayTimeout, sets the status code to http.StatusGatewayTimeout.
// Otherwise, sets the status code to http.StatusInternalServerError.
// Builds a default error response of the gateway by creating a new Response object with the modified status code,
// a new failed header, and a body constructed from the error.
// Returns the modified Response object.
func (r Response) Error(err error) Response {
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
	return Response{
		statusCode: statusCode,
		header:     newHeaderFailed(),
		body:       newBodyByErr(err),
		abort:      true,
	}
}

// Abort returns the value of the `abort` property of the Response object.
// If `abort` is true, it indicates that the response should be aborted.
// Returns a boolean value representing the `abort` property.
func (r Response) Abort() bool {
	return r.abort
}

// AbortResponse checks if the abort flag is set to true in the Response object.
// If it is true, returns a pointer to the same Response object because it indicates that an error has occurred.
// Otherwise, it checks if the AbortSequencial method in the Endpoint struct returns true for the current Response object.
// If it does, calls the abortEndpoint method and returns a pointer to the aborting Response object.
// If neither condition is met, returns nil.
func (r Response) AbortResponse() *Response {
	if r.abort {
		// caso a resposta vem como abort true, retornamos o mesmo, isso acontece, pois ocorreu um erro
		return &r
	} else if r.endpoint.AbortSequencial(r) {
		// caso o abort sequencial no endpoint vo retorne true, retornamos o endpoint abortando o mesmo
		return r.abortEndpoint()
	}
	// se não for para abortar retornamos nil
	return nil
}

// LastBackendResponse returns the last backendResponse object in the history of the Response object.
// If the history is empty, it returns an empty backendResponse object.
func (r Response) LastBackendResponse() backendResponse {
	if helper.IsEmpty(r.history) {
		return backendResponse{}
	}
	return r.history[len(r.history)-1]
}

// Endpoint returns the endpoint of the Response object.
func (r Response) Endpoint() Endpoint {
	return r.endpoint
}

// StatusCode returns the status code of the Response object.
func (r Response) StatusCode() int {
	return r.statusCode
}

// Header returns the header of the Response object.
func (r Response) Header() Header {
	return r.header
}

// Body returns the body of the Response object.
func (r Response) Body() Body {
	return r.body
}

// Eval returns a map representation of the Response object.
// The map contains the following key-value pairs:
// - "statusCode": the integer status code of the response.
// - "header": the Header object of the response.
// - "body": the interface representation of the Body object.
func (r Response) Eval() map[string]any {
	return map[string]any{
		"statusCode": r.statusCode,
		"header":     r.header,
		"body":       r.body.Interface(),
	}
}

// notifyDataChanged updates the Response object based on the modified history.
// It checks if the endpoint has reached its completion and returns the completion status.
// Then it filters the history based on the completion status to determine if it was a success or not.
// The method retrieves the status code from the filtered history and creates a new Response header
// by aggregating the completion and success values.
// It also aggregates the headers from the filtered history.
// The method retrieves the body from the filtered history based on the endpoint's aggregateResponses flag.
// Finally, it constructs and returns a new Response object with the updated values.
func (r Response) notifyDataChanged(history responseHistory) Response {
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
	bodyByHistory := filteredHistory.Body(r.endpoint.aggregateResponses)

	// construímos o novo objeto de valor
	return Response{
		endpoint:   r.endpoint,
		statusCode: statusCodeByHistory,
		header:     header,
		body:       bodyByHistory,
		history:    history,
	}
}

// abortEndpoint creates a new Response object with the status code, header, body, and history of the last backend response.
// Returns a pointer to the created Response object.
func (r Response) abortEndpoint() *Response {
	lastBackendResponseVO := r.LastBackendResponse()
	return &Response{
		statusCode: lastBackendResponseVO.statusCode,
		header:     lastBackendResponseVO.header,
		body:       lastBackendResponseVO.body,
		history:    r.history,
	}
}

// TTL calculates the time to live (TTL) for the CacheResponse object.
// It subtracts the current time from the sum of the CreatedAt time and the Duration of the CacheResponse.
// Returns the TTL duration as a string representation.
func (c CacheResponse) TTL() string {
	sub := c.CreatedAt.Add(c.Duration).Sub(time.Now())
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
		if backendResponseVO.omitResponse && completed {
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
	return helper.IsLessThanOrEqual(r.Size(), 1)
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
	// se tiver mais de 1 resposta
	if r.MultipleResponse() {
		return http.StatusOK
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
func (r responseHistory) Header() (h Header) {
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
func (r responseHistory) Body(aggregateResponses bool) (b Body) {
	if r.MultipleResponse() {
		if aggregateResponses {
			return r.aggregateBody()
		} else {
			return r.aggregatedBodies()
		}
	} else if r.SingleResponse() {
		return r[0].body
	}

	return b
}

// aggregateBody aggregates the response history into a single body. The function
// goes through each backend response in the response history. If the body of a response
// is nil, it is skipped. Otherwise, the response and its index are passed to the
// 'aggregate' method of aggregateBody. Once through with all responses, the function
// returns a new body created by the 'newBodyByAny' function, using the aggregated values.
//
// This function is a method of the responseHistory struct.
//
// Returns:
//
//	Body : The aggregated Body from the response history.
func (r responseHistory) aggregateBody() Body {
	// instanciamos primeiro o aggregate body para retornar
	var bodyHistory aggregateBody

	for index, backendResponseVO := range r {
		// se tiver nil pulamos para o próximo
		if helper.IsNil(backendResponseVO.body) {
			continue
		}
		// agregamos o backendResponse
		bodyHistory.aggregate(index, backendResponseVO)
	}

	// se tudo ocorreu bem retornamos o corpo agregado
	return newBodyByAny(bodyHistory.value)
}

// aggregatedBodies constructs an aggregated body by looping through the response history.
// For each backend response in the history, it checks if the body is nil, and if it's not,
// initializes an aggregate body and appends it to a list.
// If everything goes well, the function returns the aggregated bodies.
//
// Returns:
//
//	Body: A new slice of aggregate bodies constructed from the response history.
func (r responseHistory) aggregatedBodies() Body {
	// instanciamos o valor a ser construído
	var bodies []orderedmap.OrderedMap

	for index, backendResponseVO := range r {
		// se tiver nil pulamos para o próximo
		if helper.IsNil(backendResponseVO.body) {
			continue
		}

		// inicializamos o body agregado
		bodyGateway := newAggregateBody(index, backendResponseVO)

		// inserimos na lista de retorno
		bodies = append(bodies, bodyGateway.value)
	}
	// se tudo ocorrer bem, teremos o body agregado
	return newBodyByAny(bodies)
}

// exists checks if a given key exists in the 'aggregateBody' structure.
// It takes a string, 'key', as parameter which represents the key to be checked in the 'aggregateBody' structure.
// It returns a boolean value - 'true' if the key exists and 'false' otherwise.
func (a aggregateBody) exists(key string) bool {
	_, ok := a.value.Get(key)
	return ok
}

// notExists checks if the given key does not exist in the aggregateBody.
// It returns true if the key does not exist and false otherwise.
func (a aggregateBody) notExists(key string) bool {
	return !a.exists(key)
}

// uniqueKey checks if a given key exists in aggregateBody map.
// If the key exists, a number is appended to it, forming a unique key not existing in the map.
// This function checks the existence of the original key and also potential keys with appended number starting from 1
// until it finds a key that does not already exist.
// If the key is found to be unique, no changes are made.
//
// Parameters:
//
// key string: the candidate key to be checked/modified
//
// Returns:
//
// A string which is the unique key guaranteed to not exist so far in the aggregateBody.
func (a aggregateBody) uniqueKey(key string) string {
	// verificamos se o valor ja existe no mapper, se existir adicionamos um número
	if a.exists(key) {
		// faremos um count ate achar o campo que não existe no map ainda
		count := 1
		exists := a.exists(key)

		// iteramos até encontrar um campo que não exists
		for exists {
			tempKey := fmt.Sprintf("%s %v", key, count)
			exists = a.exists(tempKey)
			if !exists {
				key = tempKey
			}
			count++
		}
	}

	return key
}

// append is a function method that adds a new key value pair to the AggregateBody.
// It ensures that the key is unique to avoid overwriting existing values.
// 'key' is the string that will be used as the key in the key-value pair.
// 'value' is the value to be associated with the provided key.
// This function does not return any values.
func (a aggregateBody) append(key string, value any) {
	// chamamos o generateKey para gerar uma chave a partir da informada que nao existe no body ainda
	a.value.Set(a.uniqueKey(key), value)
}

// aggregate is a method on the aggregateBody type that aggregates backend responses.
// The aggregateBody type is an ordered map that collects and organizes backend responses.
//
// The method takes an index and a backendResponseVO (value object) as input parameters.
// It uses the index to get a key from the response and the value object to access the response body.
//
// The aggregate method determines the type of the response body and structures it accordingly.
// If the body is a slice or string type, or groupResponse condition is set to true, the body is appended to the key.
// If the body is a struct type, it's treated as an ordered map.
// Then, every value from the ordered map is appended to the corresponding key in the aggregate body,
// preserving the original order of the map.
//
// This method doesn't return a value. It modifies the aggregateBody in-place.
func (a aggregateBody) aggregate(index int, backendResponseVO backendResponse) {
	// obtemos a chave do backendResponse pelo index
	key := backendResponseVO.Key(index)

	// agregamos pelo tipo de dado que o body é ou forcamos se o groupResponse for true
	body := backendResponseVO.Body()
	bodyValue := body.Value()
	if backendResponseVO.groupResponse || helper.IsSliceType(bodyValue) || helper.IsStringType(bodyValue) {
		// se o body de resposta for um slice ou string, agrupamos
		a.append(key, bodyValue)
	} else if helper.IsStructType(bodyValue) {
		// se o body for uma estrutura quer dizer que ele é um ordered map
		orderedMap := body.OrderedMap()
		// iteramos para agregar no corpo principal mantendo a ordem
		for _, orderedKey := range orderedMap.Keys() {
			valueToAppend, exists := orderedMap.Get(orderedKey)
			if exists {
				a.append(orderedKey, valueToAppend)
			}
		}
	}
}
