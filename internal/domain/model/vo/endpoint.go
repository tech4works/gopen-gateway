package vo

type ExecuteEndpoint struct {
	gopen    Gopen
	endpoint Endpoint
	request  Request
}

func NewExecuteEndpoint(gopenVO Gopen, endpointVO Endpoint, url, method string, header Header, params Params,
	query Query, body any) ExecuteEndpoint {
	return ExecuteEndpoint{
		gopen:    gopenVO,
		endpoint: endpointVO,
		request:  newRequest(url, method, header, params, query, body),
	}
}

func (e ExecuteEndpoint) Gopen() Gopen {
	return e.gopen
}

func (e ExecuteEndpoint) Endpoint() Endpoint {
	return e.endpoint
}

func (e ExecuteEndpoint) Request() Request {
	return e.request
}

func (e ExecuteEndpoint) Middleware(key string) (Backend, bool) {
	backend, ok := e.gopen.middlewares[key]
	if !ok {
		return Backend{}, false
	}
	return newMiddlewareBackend(backend, backendExtraConfig{
		omitResponse: true,
	}), true
}

func (e ExecuteEndpoint) Backends() []Backend {
	return e.endpoint.backends
}
