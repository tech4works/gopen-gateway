package vo

type ExecuteBackend struct {
	backend  Backend
	request  Request
	response Response
}

type ExecuteModifierInRequestContext struct {
	backendModifier backendModifier
	request         Request
	response        Response
}

type ExecuteModifierInResponseContext struct {
	backendModifier backendModifier
	request         Request
	response        Response
}

func NewExecuteBackend(backendVO Backend, requestVO Request, responseVO Response) ExecuteBackend {
	return ExecuteBackend{
		backend:  backendVO,
		request:  requestVO,
		response: responseVO,
	}
}

func NewExecuteModifierInRequestContext(executeBackendVO ExecuteBackend, requestVO Request,
) ExecuteModifierInRequestContext {
	return ExecuteModifierInRequestContext{
		backendModifier: executeBackendVO.backend.modifier,
		request:         requestVO,
		response:        executeBackendVO.response,
	}
}

func NewExecuteModifierInResponseContext(executeBackendVO ExecuteBackend, responseVO Response,
) ExecuteModifierInResponseContext {
	return ExecuteModifierInResponseContext{
		backendModifier: executeBackendVO.backend.modifier,
		request:         executeBackendVO.request,
		response:        responseVO,
	}
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

func (e ExecuteModifierInRequestContext) Request() Request {
	return e.request
}

func (e ExecuteModifierInRequestContext) ModifierHeader() []Modifier {
	return e.backendModifier.header
}

func (e ExecuteModifierInRequestContext) ModifierParams() []Modifier {
	return e.backendModifier.params
}

func (e ExecuteModifierInRequestContext) ModifierQuery() []Modifier {
	return e.backendModifier.query
}

func (e ExecuteModifierInRequestContext) ModifierBody() []Modifier {
	return e.backendModifier.body
}

func (e ExecuteModifierInRequestContext) Response() Response {
	return e.response
}

func (e ExecuteModifierInResponseContext) Request() Request {
	return e.request
}

func (e ExecuteModifierInResponseContext) Response() Response {
	return e.response
}

func (e ExecuteModifierInResponseContext) ModifierHeader() []Modifier {
	return e.backendModifier.header
}

func (e ExecuteModifierInResponseContext) ModifierBody() []Modifier {
	return e.backendModifier.body
}

func (e ExecuteModifierInResponseContext) ModifierStatusCode() Modifier {
	return e.backendModifier.statusCode
}
