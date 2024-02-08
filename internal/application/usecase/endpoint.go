package usecase

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/dto"
	"github.com/GabrielHCataldo/martini-gateway/internal/domain/model/valueobject"
	"github.com/GabrielHCataldo/martini-gateway/internal/domain/service"
	"github.com/iancoleman/orderedmap"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Host     string            `json:"host,omitempty"`
	Endpoint string            `json:"endpoint,omitempty"`
	Url      string            `json:"url,omitempty"`
	Method   string            `json:"method,omitempty"`
	Header   http.Header       `json:"header,omitempty"`
	Query    url.Values        `json:"query,omitempty"`
	Params   map[string]string `json:"params,omitempty"`
	Body     any               `json:"body,omitempty"`
}

type Response struct {
	StatusCode int         `json:"statusCode,omitempty"`
	Header     http.Header `json:"header,omitempty"`
	Body       any         `json:"body,omitempty"`
	Group      string      `json:"group,omitempty"`
	Hide       bool        `json:"hide,omitempty"`
}

type ExecuteInput struct {
	Martini       dto.Martini
	Endpoint      dto.Endpoint
	XForwardedFor string
	Header        http.Header
	Query         url.Values
	Params        map[string]string
	Body          any
}

type ExecuteOutput struct {
	StatusCode int
	Header     map[string]string
	Body       any
}

type endpoint struct {
	backendService  service.Backend
	modifierService service.Modifier
}

type Endpoint interface {
	Execute(ctx context.Context, input ExecuteInput) (*ExecuteOutput, error)
}

func NewEndpoint(backendService service.Backend, modifierService service.Modifier) Endpoint {
	return endpoint{
		backendService:  backendService,
		modifierService: modifierService,
	}
}

func (e endpoint) Execute(ctx context.Context, input ExecuteInput) (*ExecuteOutput, error) {
	var requests []Request
	var responses []Response
	// primeiro processamos as requisições de autorização configuradas
	for _, authorizationKey := range input.Endpoint.Authorizations {
		authBackend, ok := input.Martini.ExtraConfig.Authorizations[authorizationKey]
		if !ok {
			return nil, errors.New("authorization", authorizationKey, "not configured on extra-config.authorizations")
		}
		// montamos a request com base no authBackend
		request, err := e.buildBackendRequest(authBackend, input, requests, responses)
		if helper.IsNotNil(err) {
			return nil, err
		}
		requests = append(requests, *request)
		// processamos o backend de autorização
		response, err := e.makeBackendRequest(ctx, authBackend, *request, requests, responses)
		//qualquer erro na requisição de autorização ja paramos e retornamos o mesmo
		if helper.IsNotNil(err) {
			return nil, err
		} else if helper.IsNotEqualTo(response.StatusCode, http.StatusOK) {
			return e.buildExecuteOutput(input.Endpoint, []Response{*response})
		}
		responses = append(responses, *response)
	}
	// em seguida executamos o backend
	for _, backend := range input.Endpoint.Backends {
		// montamos a request com base no authBackend
		request, err := e.buildBackendRequest(backend, input, requests, responses)
		if helper.IsNotNil(err) {
			return nil, err
		}
		requests = append(requests, *request)
		// processamos o backend
		response, err := e.makeBackendRequest(ctx, backend, *request, requests, responses)
		//qualquer erro na requisição e teríamos mais de 1 backend e o campo abortSequential estiver true,
		//abortamos as demais requisições
		if helper.IsNotNil(err) {
			return nil, err
		} else if helper.IsGreaterThan(input.Endpoint.Backends, 1) &&
			helper.IsGreaterThanOrEqual(response.StatusCode, http.StatusBadRequest) &&
			input.Endpoint.AbortSequential {
			return e.buildExecuteOutput(input.Endpoint, []Response{*response})
		}
		responses = append(responses, *response)
	}
	return e.buildExecuteOutput(input.Endpoint, responses)
}

func (e endpoint) buildBackendRequest(
	backend dto.Backend,
	input ExecuteInput,
	requests []Request,
	responses []Response,
) (*Request, error) {
	//preparamos a requisição com base o dto backend
	// todo: aqui obtemos o host (correto é criar um domínio chamado balancer aonde ele vai retornar o host
	//  disponível pegando como base, se ele esta de pé ou não, e sua config de porcentagem)
	host := backend.Host[len(backend.Host)-1]
	backendRequest, err := e.backendService.BuildBackendRequest(service.BuildBackendRequestInput{
		Backend: valueobject.BuildBackend(backend),
		Host:    host,
		Header:  input.Header,
		Query:   input.Query,
		Params:  input.Params,
		Body:    input.Body,
	})
	if helper.IsNotNil(err) {
		return nil, err
	}
	//modificamos a requisição no escopo request no domínio modifier
	request, err := e.modifierRequest(backend, *backendRequest, requests, responses)
	if helper.IsNotNil(err) {
		return nil, err
	}
	//aqui forcamos o X-Forwarded-For no header da requisição
	request.Header.Set("X-Forwarded-For", input.XForwardedFor)
	//retornamos o objeto completo, pronto para prosseguir com a requisição
	return &Request{
		Host:     request.Host,
		Endpoint: request.Endpoint,
		Url:      request.Url,
		Method:   request.Method,
		Header:   request.Header,
		Query:    request.Query,
		Params:   request.Params,
		Body:     request.Body,
	}, nil
}

func (e endpoint) makeBackendRequest(
	ctx context.Context,
	backend dto.Backend,
	request Request,
	requests []Request,
	responses []Response,
) (*Response, error) {
	//executamos a requisição backend
	backendResponse, err := e.backendService.Execute(ctx, service.ExecuteBackendInput{
		Backend: valueobject.BuildBackend(backend),
		BackendRequest: service.BackendRequest{
			Host:     request.Host,
			Endpoint: request.Endpoint,
			Url:      request.Url,
			Method:   request.Method,
			Header:   request.Header,
			Query:    request.Query,
			Params:   request.Params,
			Body:     request.Body,
		},
	})
	if helper.IsNotNil(err) {
		return nil, errBadGateway(err)
	}
	//modificamos a response no escopo response no domínio modifier
	response, err := e.modifierResponse(backend, *backendResponse, requests, responses)
	if helper.IsNotNil(err) {
		return nil, err
	}
	return &Response{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Body:       response.Body,
		Group:      backend.Group,
		Hide:       backend.HideResponse,
	}, nil
}

func (e endpoint) modifierRequest(
	backend dto.Backend,
	backendRequest service.BackendRequest,
	requests []Request,
	responses []Response,
) (*service.ModifierRequest, error) {
	return e.modifierService.ExecuteRequestScope(service.ExecuteRequestScopeInput{
		Request:         BuildModifierRequestByBackendRequest(backendRequest),
		Headers:         valueobject.BuildModifiers(backend.Headers),
		Params:          valueobject.BuildModifiers(backend.Params),
		Queries:         valueobject.BuildModifiers(backend.Queries),
		Body:            valueobject.BuildModifiers(backend.Body),
		RequestHistory:  BuildModifierRequests(requests),
		ResponseHistory: BuildModifierResponses(responses),
	})
}

func (e endpoint) modifierResponse(
	backend dto.Backend,
	backendResponse service.BackendResponse,
	requests []Request,
	responses []Response,
) (*service.ModifierResponse, error) {
	return e.modifierService.ExecuteResponseScope(service.ExecuteResponseScopeInput{
		Response:        BuildModifierResponseByBackendResponse(backendResponse),
		Headers:         valueobject.BuildModifiers(backend.Headers),
		Params:          valueobject.BuildModifiers(backend.Params),
		Queries:         valueobject.BuildModifiers(backend.Queries),
		Body:            valueobject.BuildModifiers(backend.Body),
		RequestHistory:  BuildModifierRequests(requests),
		ResponseHistory: BuildModifierResponses(responses),
	})
}

func (e endpoint) buildExecuteOutput(endpoint dto.Endpoint, responses []Response) (*ExecuteOutput, error) {
	var statusCode int
	var body any
	gatewayCompleted, responsesFiltered := e.filterResponses(responses)
	if helper.IsGreaterThan(endpoint.Backends, 1) && helper.IsGreaterThan(responsesFiltered, 1) {
		statusCode = http.StatusOK
		if endpoint.AggregateResponses {
			body = e.buildAggregateResponsesBody(responsesFiltered)
		} else {
			body = responsesFiltered
		}
	} else if helper.IsGreaterThan(responsesFiltered, 0) {
		statusCode = responsesFiltered[0].StatusCode
		body = responsesFiltered[0].Body
	} else {
		return nil, errors.New("endpoint is not response nothing:", endpoint.Endpoint, "method:", endpoint.Method)
	}
	return &ExecuteOutput{
		StatusCode: statusCode,
		Header:     e.buildHeaderResponse(gatewayCompleted, responsesFiltered),
		Body:       body,
	}, nil
}

func (e endpoint) filterResponses(responses []Response) (gatewayCompleted bool, responsesFiltered []Response) {
	gatewayCompleted = true
	for _, response := range responses {
		if response.Hide && helper.IsGreaterThan(responses, 1) {
			continue
		} else if helper.IsGreaterThanOrEqual(response.StatusCode, http.StatusBadRequest) {
			gatewayCompleted = false
		}
		response.Header = nil
		response.Hide = false
		responsesFiltered = append(responsesFiltered, response)
	}
	return gatewayCompleted, responsesFiltered
}

func (e endpoint) buildHeaderResponse(gatewayCompleted bool, responses []Response) map[string]string {
	header := map[string]string{
		"X-Gateway-Completed": helper.SimpleConvertToString(gatewayCompleted),
	}
	for _, respBackend := range responses {
		for k, v := range respBackend.Header {
			if helper.IsNotEqualToIgnoreCase(k, "content-length") &&
				helper.IsNotEqualToIgnoreCase(k, "content-type") {
				header[k] = strings.Join(v, ", ")
			}
		}
	}
	return header
}

func (e endpoint) buildAggregateResponsesBody(responses []Response) *orderedmap.OrderedMap {
	result := orderedmap.New()
	for index, response := range responses {
		if helper.IsNil(response.Body) {
			continue
		}
		genericKey := fmt.Sprint("response", index)
		if helper.IsSlice(response.Body) {
			key := genericKey
			if helper.IsNotEmpty(response.Group) {
				key = response.Group
			}
			result.Set(key, response.Body)
		} else if helper.IsStruct(response.Body) { //todo -> pq ele ta como map e nao como ordered?
			body := response.Body.(orderedmap.OrderedMap)
			if helper.IsNotEmpty(response.Group) {
				bodyResult := orderedmap.New()
				for _, k := range body.Keys() {
					v, ok := body.Get(k)
					if ok {
						bodyResult.Set(k, v)
					}
				}
				result.Set(response.Group, bodyResult)
			} else {
				for _, k := range body.Keys() {
					if v, ok := result.Get(k); ok {
						result.Set(fmt.Sprint(k, index), v)
					} else {
						v, ok := body.Get(k)
						if ok {
							result.Set(k, v)
						}
					}
				}
			}
		} else if helper.IsString(response.Body) {
			key := genericKey
			if helper.IsNotEmpty(response.Group) {
				key = response.Group
			}
			result.Set(key, helper.SimpleConvertToString(response))
		}
	}
	if helper.IsEmpty(result) {
		return nil
	}
	return result
}
