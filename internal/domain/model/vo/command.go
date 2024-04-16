package vo

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

// ExecuteBackend is a type that represents the execution of a backend server request and response.
type ExecuteBackend struct {
	endpoint *Endpoint
	// Backend represents a backend server configuration.
	backend *Backend
	// Request represents an HTTP `request` object.
	request *Request
	// Response represents an HTTP `response` object.
	response *Response
}

// ExecuteModifier is a type that represents the execution of a modifier on a backend request and response.
// It contains fields for the modifier context, backend modifiers, request, and response.
type ExecuteModifier struct {
	// context represents the context in which a modification should be applied.
	context enum.ModifierContext
	// backendModifiers represents the set of modifiers for a backend configuration. It contains fields for the
	// status code, header, params, query, and body modifiers.
	backendModifiers *BackendModifiers
	// Request represents an HTTP `request` object.
	request *Request
	// Response represents an HTTP `response` object.
	response *Response
}

// ExecuteEndpoint represents the execution of a specific endpoint in the Gopen server.
// It contains the configuration for the Gopen server, the targeted endpoint, and the HTTP request for execution.
type ExecuteEndpoint struct {
	// Gopen represents the configuration for the Gopen server, including environment, version, hot reload status, port,
	// timeout duration, limiter, cache, security CORS, middlewares, and endpoints.
	gopen *Gopen
	// endpoint represents a specific endpoint in the Gopen server.
	endpoint *Endpoint
	// Request represents an HTTP `request` object.
	request *Request
}

// NewExecuteEndpoint creates a new ExecuteEndpoint using the provided Gopen, Endpoint, and Request objects.
func NewExecuteEndpoint(gopenVO *Gopen, endpointVO *Endpoint, requestVO *Request) *ExecuteEndpoint {
	return &ExecuteEndpoint{
		gopen:    gopenVO,
		endpoint: endpointVO,
		request:  requestVO,
	}
}

// NewExecuteBackend creates a new ExecuteBackend using the provided Endpoint, Backend, Request, and Response objects.
func NewExecuteBackend(endpointVO *Endpoint, backendVO *Backend, requestVO *Request, responseVO *Response) *ExecuteBackend {
	return &ExecuteBackend{
		endpoint: endpointVO,
		backend:  backendVO,
		request:  requestVO,
		response: responseVO,
	}
}

// NewExecuteRequestModifier creates a new ExecuteModifier using the provided Backend, Request, and Response objects.
func NewExecuteRequestModifier(backendVO *Backend, requestVO *Request, responseVO *Response,
) *ExecuteModifier {
	return &ExecuteModifier{
		context:          enum.ModifierContextRequest,
		backendModifiers: backendVO.modifiers,
		request:          requestVO,
		response:         responseVO,
	}
}

// NewExecuteResponseModifier creates a new ExecuteModifier using the provided Backend, Request, and Response objects.
// The ExecuteModifier modifies the response in the context of enum.ModifierContextResponse, using the backend modifier
// functions. It returns the newly created ExecuteModifier.
func NewExecuteResponseModifier(backendVO *Backend, requestVO *Request, responseVO *Response,
) *ExecuteModifier {
	return &ExecuteModifier{
		context:          enum.ModifierContextResponse,
		backendModifiers: backendVO.modifiers,
		request:          requestVO,
		response:         responseVO,
	}
}

// Gopen returns the Gopen object associated with the ExecuteEndpoint object.
func (e ExecuteEndpoint) Gopen() *Gopen {
	return e.gopen
}

// Endpoint returns the Endpoint object associated with the ExecuteEndpoint object.
func (e ExecuteEndpoint) Endpoint() *Endpoint {
	return e.endpoint
}

// Request returns the Request object associated with the ExecuteEndpoint object.
func (e ExecuteEndpoint) Request() *Request {
	return e.request
}

// Endpoint returns the Endpoint object associated with the ExecuteEndpoint object.
func (e ExecuteBackend) Endpoint() *Endpoint {
	return e.endpoint
}

// Backend returns the Backend object associated with the ExecuteBackend object.
func (e ExecuteBackend) Backend() *Backend {
	return e.backend
}

// Request returns the Request object associated with the ExecuteBackend object.
func (e ExecuteBackend) Request() *Request {
	return e.request
}

// Response returns the Response object associated with the ExecuteBackend object.
func (e ExecuteBackend) Response() *Response {
	return e.response
}

// Request returns the Request object associated with the ExecuteModifier object.
func (e ExecuteModifier) Request() *Request {
	return e.request
}

// Context returns the enum.ModifierContext object associated with the ExecuteModifier object.
func (e ExecuteModifier) Context() enum.ModifierContext {
	return e.context
}

// BackendModifiers returns the BackendModifiers object associated with the ExecuteModifier object.
func (e ExecuteModifier) BackendModifiers() *BackendModifiers {
	return e.backendModifiers
}

// ModifierHeader returns the header modifiers associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierHeader() []Modifier {
	return e.backendModifiers.header
}

// ModifierParams returns the params modifiers associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierParams() []Modifier {
	return e.backendModifiers.params
}

// ModifierQuery returns the query modifiers associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierQuery() []Modifier {
	return e.backendModifiers.query
}

// ModifierBody returns the body modifiers associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierBody() []Modifier {
	return e.backendModifiers.body
}

// ModifierStatusCode returns the status code modifier associated with the ExecuteModifier object.
func (e ExecuteModifier) ModifierStatusCode() *Modifier {
	return e.backendModifiers.statusCode
}

// Response returns the Response modifier associated with the ExecuteModifier object.
func (e ExecuteModifier) Response() *Response {
	return e.response
}
