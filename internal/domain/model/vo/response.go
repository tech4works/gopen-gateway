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

type CacheResponse struct {
	StatusCode int        `json:"statusCode"`
	Header     Header     `json:"header"`
	Body       *CacheBody `json:"body,omitempty"`
	Duration   string     `json:"duration"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type Response struct {
	endpoint   Endpoint
	statusCode int
	header     Header
	body       Body
	abort      bool
	history    responseHistory
}

type responseHistory []backendResponse

type errorResponseBody struct {
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
	Endpoint  string    `json:"endpoint,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
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

	var body string
	if helper.IsNotNil(cacheResponseVO.Body) {
		body = cacheResponseVO.Body.value
	}

	return Response{
		endpoint:   endpointVO,
		statusCode: cacheResponseVO.StatusCode,
		header:     header,
		body:       NewBodyFromString(body),
	}
}

func NewResponseByErr(endpointVO Endpoint, statusCode int, err error) Response {
	return Response{
		endpoint:   endpointVO,
		statusCode: statusCode,
		header:     newHeaderFailed(),
		body:       newErrorBody(endpointVO, err),
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
func NewCacheResponse(responseVO Response, duration time.Duration) CacheResponse {
	return CacheResponse{
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
func (r Response) ModifyLastBackendResponse(backendResponseVO backendResponse) Response {
	history := r.history
	history[len(history)-1] = backendResponseVO

	// atualizamos os dados a partir do histórico alterado
	return r.notifyDataChanged(history)
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
		body:       newErrorBody(r.endpoint, err),
		abort:      true,
	}
}

// Abort returns the value of the `abort` property of the Response object.
// If `abort` is true, it indicates that the response should be aborted.
// Returns a boolean value representing the `abort` property.
func (r Response) Abort() bool {
	return r.abort
}

// IsAbortResponse if the abort field of the response structure is true or func of the AbortSequential endpoint returns
// true, we return true, otherwise it returns false.
func (r Response) IsAbortResponse() bool {
	return r.abort || r.endpoint.AbortSequencial(r)
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

func (r Response) ContentType() enum.ContentType {
	responseEncode := r.endpoint.ResponseEncode()
	if responseEncode.IsEnumValid() {
		return responseEncode.ContentType()
	}
	return r.Body().ContentType()
}

// Body returns the body of the Response object.
func (r Response) Body() Body {
	return r.body
}

func (r Response) BodyBytes() []byte {
	// instanciamos o response encode do endpoint
	responseEncode := r.endpoint.ResponseEncode()
	// retornamos pelo responseEncode caso ele seja valido
	if responseEncode.IsEnumValid() {
		return r.body.BytesByContentType(responseEncode.ContentType())
	}
	// se não respondemos pelo tipo do body
	return r.body.Bytes()
}

// Eval returns a map representation of the Response object.
// The map contains the following key-value pairs:
// - "statusCode": the integer status code of the response.
// - "header": the Header object of the response.
// - "body": the interface representation of the Body object.
func (r Response) Eval() string {
	mapEval := map[string]any{
		"statusCode": r.statusCode,
		"header":     r.header,
		"body":       r.body.Interface(),
	}
	return helper.SimpleConvertToString(mapEval)
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
	// instanciamos o ultimo backend response, que é para ser abortado
	lastBackendResponseVO := r.LastBackendResponse()

	// instanciamos o novo header
	header := newResponseHeader(false, false)
	// agregamos o header do backend abortado
	header = header.Aggregate(lastBackendResponseVO.Header())

	// construímos o response com os dados do backend abortado
	return &Response{
		statusCode: lastBackendResponseVO.statusCode,
		header:     header,
		body:       lastBackendResponseVO.body,
		history:    r.history,
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
			return r.body()
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
	bodyHistory := Body{
		contentType: enum.ContentTypeJson,
		value:       "{}",
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

// body constructs an aggregated body by looping through the response history.
// For each backend response in the history, it checks if the body is nil, and if it's not,
// initializes an aggregate body and appends it to a list.
// If everything goes well, the function returns the aggregated bodies.
//
// Returns:
//
//	Body: A new slice of aggregate bodies constructed from the response history.
func (r responseHistory) body() Body {
	// instanciamos o valor a ser construído
	var bodies []Body
	// iteramos o histórico para listar os bodies de resposta
	for index, backendResponseVO := range r {
		// se tiver nil pulamos para o próximo
		if helper.IsNil(backendResponseVO.body) {
			continue
		}
		// inicializamos o body agregado
		bodyBackendResponse := newBodyFromBackendResponse(index, backendResponseVO)
		// inserimos na lista de retorno
		bodies = append(bodies, bodyBackendResponse)
	}
	// se tudo ocorrer bem, teremos o body agregado
	return newSliceBody(bodies)
}
