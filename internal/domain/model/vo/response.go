package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"net/http"
	"time"
)

type responseHistory []backendResponse

type aggregateBody map[string]any

type Response struct {
	endpoint   Endpoint
	completed  bool
	statusCode int
	header     Header
	body       any
	abort      bool
	history    responseHistory
}

type BodyErrorResponse struct {
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
	Endpoint  string    `json:"endpoint,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func NewResponse(endpointVO Endpoint) Response {
	return Response{
		endpoint:   endpointVO,
		statusCode: http.StatusNoContent,
	}
}

func NewBodyErrorResponse(endpoint string, err error) BodyErrorResponse {
	errDetails := errors.Details(err)
	return BodyErrorResponse{
		File:      errDetails.GetFile(),
		Line:      errDetails.GetLine(),
		Endpoint:  endpoint,
		Message:   errDetails.GetMessage(),
		Timestamp: time.Now(),
	}
}

func newAggregateBody(index int, backendResponseVO backendResponse) aggregateBody {
	newInstance := aggregateBody{
		"ok":   backendResponseVO.Ok(),
		"code": backendResponseVO.statusCode,
	}
	newInstance.aggregate(index, backendResponseVO)
	return newInstance
}

func (r Response) Modify(statusCode int, header Header, body any, backendResponseVO backendResponse) Response {
	history := r.history
	history[len(history)-1] = backendResponseVO
	return Response{
		endpoint:   r.endpoint,
		completed:  r.completed,
		statusCode: statusCode,
		header:     header,
		body:       body,
		abort:      r.abort,
		history:    history,
	}
}

func (r Response) Append(backendResponseVO backendResponse) Response {
	// adicionamos na nova lista de histórico
	history := r.history
	history = append(history, backendResponseVO)

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
		completed:  completed,
		statusCode: statusCodeByHistory,
		header:     header,
		body:       bodyByHistory,
		history:    history,
	}
}

func (r Response) Err(requestUrl string, err error) Response {
	// construímos o statusCode de resposta a partir do erro recebido
	statusCode := http.StatusInternalServerError
	if errors.Contains(err, mapper.ErrBadGateway) {
		statusCode = http.StatusBadGateway
	} else if errors.Contains(err, mapper.ErrGatewayTimeout) {
		statusCode = http.StatusGatewayTimeout
	}

	// construímos a resposta de erro padrão do gateway
	return Response{
		statusCode: statusCode,
		header:     newResponseHeaderFailed(),
		body:       NewBodyErrorResponse(requestUrl, err),
		abort:      true,
		history:    r.history,
	}
}

func (r Response) Abort() *Response {
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

func (r Response) abortEndpoint() *Response {
	lastBackendResponseVO := r.LastBackendResponse()
	return &Response{
		statusCode: lastBackendResponseVO.statusCode,
		header:     newResponseHeaderFailed(),
		body:       lastBackendResponseVO.body,
		history:    r.history,
	}
}

func (r Response) LastBackendResponse() backendResponse {
	return r.history[len(r.history)-1]
}

func (r Response) StatusCode() int {
	return r.statusCode
}

func (r Response) Header() Header {
	return r.header
}

func (r Response) Body() any {
	return r.body
}

func (r responseHistory) Size() int {
	return len(r)
}

func (r responseHistory) Success() bool {
	for _, backendResponseVO := range r {
		if helper.IsGreaterThanOrEqual(backendResponseVO.statusCode, http.StatusBadRequest) {
			return false
		}
	}
	return true
}

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

func (r responseHistory) SingleResponse() bool {
	return helper.IsLessThanOrEqual(r.Size(), 1)
}

func (r responseHistory) MultipleResponse() bool {
	return helper.IsGreaterThan(r.Size(), 1)
}

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

func (r responseHistory) Body(aggregateResponses bool) any {
	if r.MultipleResponse() {
		if aggregateResponses {
			return r.aggregateBody()
		} else {
			return r.aggregatedBodies()
		}
	} else if r.SingleResponse() {
		return r[0].body
	}

	return nil
}

func (r responseHistory) aggregateBody() (body aggregateBody) {
	for index, backendResponseVO := range r {
		// se tiver nil pulamos para o próximo
		if helper.IsNil(backendResponseVO.body) {
			continue
		}
		// agregamos o backendResponse
		body.aggregate(index, backendResponseVO)
	}

	// se o agregado for vazio, retornamos nil
	if helper.IsEmpty(body) {
		return nil
	}
	// se tudo ocorreu bem retornamos o corpo agregado
	return body
}

func (r responseHistory) aggregatedBodies() (bodies []aggregateBody) {
	for index, backendResponseVO := range r {
		// se tiver nil pulamos para o próximo
		if helper.IsNil(backendResponseVO.body) {
			continue
		}

		// inicializamos o body agregado
		bodyGateway := newAggregateBody(index, backendResponseVO)

		// inserimos na lista de retorno
		bodies = append(bodies, bodyGateway)
	}
	// se tudo ocorrer bem, teremos o body agregado
	return bodies
}

func (a aggregateBody) exists(key string) bool {
	_, ok := a[key]
	return ok
}

func (a aggregateBody) notExists(key string) bool {
	return !a.exists(key)
}

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

func (a aggregateBody) append(key string, value any) {
	// chamamos o generateKey para gerar uma chave a partir da informada que nao existe no body ainda
	a[a.uniqueKey(key)] = value
}

func (a aggregateBody) aggregate(index int, backendResponseVO backendResponse) {
	// obtemos a chave do backendResponse pelo index
	key := backendResponseVO.Key(index)

	// agregamos pelo tipo de dado que o body é
	if helper.IsSlice(backendResponseVO.body) || helper.IsString(backendResponseVO.body) {
		// se o body de resposta for um slice ou string, agrupamos
		a.append(key, backendResponseVO.body)
	} else if helper.IsJson(backendResponseVO.body) {
		// se o body for um map
		bodyMap := backendResponseVO.body.(map[string]any)
		// se ele quer agrupar a resposta, então fazemos isso
		if backendResponseVO.groupResponse {
			a.append(key, bodyMap)
		} else {
			// iteramos para agregar no corpo principal
			for mKey, value := range bodyMap {
				a.append(mKey, value)
			}
		}
	}
}
