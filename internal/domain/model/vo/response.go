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
	completed  bool
	statusCode int
	header     Header
	body       Body
	err        error
	abort      bool
	history    responseHistory
}

func NewResponse(endpointVO Endpoint) Response {
	return Response{
		endpoint:   endpointVO,
		statusCode: http.StatusNoContent,
	}
}

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

func NewCacheResponse(writer dto.ResponseWriter, duration time.Duration) CacheResponse {
	return CacheResponse{
		StatusCode: writer.Status(),
		Header:     NewHeader(writer.Header()),
		Body:       newBody(writer.Body.Bytes()),
		Duration:   duration,
		CreatedAt:  time.Now(),
	}
}

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

func (r Response) ModifyStatusCode(statusCode int, backendResponseVO backendResponse) Response {
	history := r.history
	history[len(history)-1] = backendResponseVO

	// se o valor vindo do parâmetro vir igual, damos prioridade ao local
	if helper.Equals(r.statusCode, statusCode) {
		// atualizamos os dados a partir do histórico alterado
		return r.notify(history)
	}

	return Response{
		endpoint:   r.endpoint,
		completed:  r.completed,
		statusCode: statusCode,
		header:     r.header,
		body:       r.body,
		abort:      r.abort,
		history:    history,
	}
}

func (r Response) ModifyHeader(header Header, backendResponseVO backendResponse) Response {
	history := r.history
	history[len(history)-1] = backendResponseVO

	// se o valor vindo do parâmetro vir igual, damos prioridade ao local
	if helper.Equals(r.header, header) {
		// atualizamos os dados a partir do histórico alterado
		return r.notify(history)
	}

	return Response{
		endpoint:   r.endpoint,
		completed:  r.completed,
		statusCode: r.statusCode,
		header:     header,
		body:       r.body,
		abort:      r.abort,
		history:    history,
	}
}

func (r Response) ModifyBody(body Body, backendResponseVO backendResponse) Response {
	history := r.history
	history[len(history)-1] = backendResponseVO

	// se o valor vindo do parâmetro vir igual, damos prioridade ao local
	if helper.Equals(r.body, body) {
		// atualizamos os dados a partir do histórico alterado
		return r.notify(history)
	}

	return Response{
		endpoint:   r.endpoint,
		completed:  r.completed,
		statusCode: r.statusCode,
		header:     r.header,
		body:       body,
		abort:      r.abort,
		history:    history,
	}
}

func (r Response) Append(backendResponseVO backendResponse) Response {
	// adicionamos na nova lista de histórico
	history := r.history
	history = append(history, backendResponseVO)

	// atualizamos os dados a partir do histórico alterado
	return r.notify(history)
}

func (r Response) notify(history responseHistory) Response {
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

func (r Response) Error(err error) Response {
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
		header:     NewHeaderFailed(),
		err:        err,
		abort:      true,
		history:    r.history,
	}
}

func (r Response) Abort() bool {
	return r.abort
}

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

func (r Response) abortEndpoint() *Response {
	lastBackendResponseVO := r.LastBackendResponse()
	return &Response{
		statusCode: lastBackendResponseVO.statusCode,
		header:     lastBackendResponseVO.header,
		body:       lastBackendResponseVO.body,
		history:    r.history,
	}
}

func (r Response) LastBackendResponse() backendResponse {
	if helper.IsEmpty(r.history) {
		return backendResponse{}
	}
	return r.history[len(r.history)-1]
}

func (r Response) Endpoint() Endpoint {
	return r.endpoint
}

func (r Response) StatusCode() int {
	return r.statusCode
}

func (r Response) Header() Header {
	return r.header
}

func (r Response) Body() Body {
	return r.body
}

func (r Response) Err() error {
	return r.err
}

func (r Response) Eval() map[string]any {
	return map[string]any{
		"statusCode": r.statusCode,
		"header":     r.header,
		"body":       r.body.Interface(),
	}
}

func (c CacheResponse) TTL() string {
	sub := c.CreatedAt.Add(c.Duration).Sub(time.Now())
	return sub.String()
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

func (a aggregateBody) IsEmpty() bool {
	return helper.IsEmpty(a.value.Keys())
}

func (a aggregateBody) exists(key string) bool {
	_, ok := a.value.Get(key)
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
	a.value.Set(a.uniqueKey(key), value)
}

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
