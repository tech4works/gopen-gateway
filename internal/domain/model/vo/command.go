package vo

import "github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"

type ExecuteBackend struct {
	backend  Backend
	request  Request
	response Response
}

type ExecuteModifier struct {
	context         enum.ModifierContext
	backendModifier BackendModifiers
	request         Request
	response        Response
}

type ExecuteEndpoint struct {
	gopen    GOpen
	endpoint Endpoint
	request  Request
}

func NewExecuteEndpoint(gopenVO GOpen, endpointVO Endpoint, url, method string, header Header, params Params,
	query Query, body Body) ExecuteEndpoint {
	return ExecuteEndpoint{
		gopen:    gopenVO,
		endpoint: endpointVO,
		request:  newRequest(url, method, header, params, query, body),
	}
}

func NewExecuteBackend(backendVO Backend, requestVO Request, responseVO Response) ExecuteBackend {
	return ExecuteBackend{
		backend:  backendVO,
		request:  requestVO,
		response: responseVO,
	}
}

func NewExecuteRequestModifier(backendVO Backend, requestVO Request, responseVO Response,
) ExecuteModifier {
	return ExecuteModifier{
		context:         enum.ModifierContextRequest,
		backendModifier: backendVO.modifiers,
		request:         requestVO,
		response:        responseVO,
	}
}

func NewExecuteResponseModifier(backendVO Backend, requestVO Request, responseVO Response,
) ExecuteModifier {
	return ExecuteModifier{
		context:         enum.ModifierContextResponse,
		backendModifier: backendVO.modifiers,
		request:         requestVO,
		response:        responseVO,
	}
}

func (e ExecuteEndpoint) Gopen() GOpen {
	return e.gopen
}

func (e ExecuteEndpoint) Endpoint() Endpoint {
	return e.endpoint
}

func (e ExecuteEndpoint) Request() Request {
	return e.request
}

func (e ExecuteBackend) Backend() Backend {
	return e.backend
}

func (e ExecuteBackend) Request() Request {
	return e.request
}

func (e ExecuteBackend) Response() Response {
	return e.response
}

func (e ExecuteModifier) Request() Request {
	return e.request
}

func (e ExecuteModifier) Context() enum.ModifierContext {
	return e.context
}

func (e ExecuteModifier) ModifierHeader() []Modifier {
	return e.backendModifier.header
}

func (e ExecuteModifier) ModifierParams() []Modifier {
	return e.backendModifier.params
}

func (e ExecuteModifier) ModifierQuery() []Modifier {
	return e.backendModifier.query
}

func (e ExecuteModifier) ModifierBody() []Modifier {
	return e.backendModifier.body
}

func (e ExecuteModifier) ModifierStatusCode() Modifier {
	return e.backendModifier.statusCode
}

func (e ExecuteModifier) Response() Response {
	return e.response
}
