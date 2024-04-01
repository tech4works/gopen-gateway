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

// NewExecuteEndpoint creates a new ExecuteEndpoint using the provided GOpen, Endpoint, and Request objects.
func NewExecuteEndpoint(gopenVO GOpen, endpointVO Endpoint, requestVO Request) ExecuteEndpoint {
	return ExecuteEndpoint{
		gopen:    gopenVO,
		endpoint: endpointVO,
		request:  requestVO,
	}
}

// NewExecuteBackend creates a new ExecuteBackend using the provided Backend, Request, and Response objects.
func NewExecuteBackend(backendVO Backend, requestVO Request, responseVO Response) ExecuteBackend {
	return ExecuteBackend{
		backend:  backendVO,
		request:  requestVO,
		response: responseVO,
	}
}

// NewExecuteRequestModifier creates a new ExecuteModifier using the provided Backend, Request, and Response objects.
func NewExecuteRequestModifier(backendVO Backend, requestVO Request, responseVO Response,
) ExecuteModifier {
	return ExecuteModifier{
		context:         enum.ModifierContextRequest,
		backendModifier: backendVO.modifiers,
		request:         requestVO,
		response:        responseVO,
	}
}

// NewExecuteResponseModifier creates a new ExecuteModifier using the provided Backend, Request, and Response objects.
// The ExecuteModifier modifies the response in the context of enum.ModifierContextResponse, using the backend modifier
// functions. It returns the newly created ExecuteModifier.
func NewExecuteResponseModifier(backendVO Backend, requestVO Request, responseVO Response,
) ExecuteModifier {
	return ExecuteModifier{
		context:         enum.ModifierContextResponse,
		backendModifier: backendVO.modifiers,
		request:         requestVO,
		response:        responseVO,
	}
}

// GOpen returns the GOpen object associated with the ExecuteEndpoint object.
func (e ExecuteEndpoint) GOpen() GOpen {
	return e.gopen
}

// Endpoint returns the Endpoint object associated with the ExecuteEndpoint object.
func (e ExecuteEndpoint) Endpoint() Endpoint {
	return e.endpoint
}

// Request returns the Request object associated with the ExecuteEndpoint object.
func (e ExecuteEndpoint) Request() Request {
	return e.request
}

// Backend returns the Backend object associated with the ExecuteBackend object.
func (e ExecuteBackend) Backend() Backend {
	return e.backend
}

// Request returns the Request object associated with the ExecuteBackend object.
func (e ExecuteBackend) Request() Request {
	return e.request
}

// Response returns the Response object associated with the ExecuteBackend object.
func (e ExecuteBackend) Response() Response {
	return e.response
}

// Request returns the Request object associated with the ExecuteModifier object.
func (e ExecuteModifier) Request() Request {
	return e.request
}

// Context returns the enum.ModifierContext object associated with the ExecuteModifier object.
func (e ExecuteModifier) Context() enum.ModifierContext {
	return e.context
}

// ModifierHeader returns the header modifiers associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierHeader() []Modifier {
	return e.backendModifier.header
}

// ModifierParams returns the params modifiers associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierParams() []Modifier {
	return e.backendModifier.params
}

// ModifierQuery returns the query modifiers associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierQuery() []Modifier {
	return e.backendModifier.query
}

// ModifierBody returns the body modifiers associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierBody() []Modifier {
	return e.backendModifier.body
}

// ModifierStatusCode returns the status code modifier associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierStatusCode() Modifier {
	return e.backendModifier.statusCode
}

// Response returns the Response modifier associated with the ExecuteModifier object.
func (e ExecuteModifier) Response() Response {
	return e.response
}
