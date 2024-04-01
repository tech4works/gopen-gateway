package vo

type Request struct {
	url     string
	method  string
	header  Header
	params  Params
	query   Query
	body    Body
	history []backendRequest
}

func NewRequest(url, method string, header Header, params Params, query Query, body Body) Request {
	return Request{
		url:    url,
		method: method,
		header: header,
		params: params,
		query:  query,
		body:   body,
	}
}

func (r Request) ModifyHeader(header Header, backendRequestVO backendRequest) Request {
	history := r.history
	history[len(history)-1] = backendRequestVO

	return Request{
		url:     r.url,
		method:  r.method,
		header:  header,
		params:  r.params,
		query:   r.query,
		body:    r.body,
		history: history,
	}
}

func (r Request) ModifyParams(params Params, backendRequestVO backendRequest) Request {
	history := r.history
	history[len(history)-1] = backendRequestVO

	return Request{
		url:     r.url,
		method:  r.method,
		header:  r.header,
		params:  params,
		query:   r.query,
		body:    r.body,
		history: history,
	}
}

func (r Request) ModifyQuery(query Query, backendRequestVO backendRequest) Request {
	history := r.history
	history[len(history)-1] = backendRequestVO

	return Request{
		url:     r.url,
		method:  r.method,
		header:  r.header,
		params:  r.params,
		query:   query,
		body:    r.body,
		history: history,
	}
}

func (r Request) ModifyBody(body Body, backendRequestVO backendRequest) Request {
	history := r.history
	history[len(history)-1] = backendRequestVO

	return Request{
		url:     r.url,
		method:  r.method,
		header:  r.header,
		params:  r.params,
		query:   r.query,
		body:    body,
		history: history,
	}
}

func (r Request) Append(backendRequest backendRequest) Request {
	return Request{
		url:     r.url,
		method:  r.method,
		header:  r.header,
		params:  r.params,
		query:   r.query,
		body:    r.body,
		history: append(r.history, backendRequest),
	}
}

func (r Request) CurrentBackendRequest() backendRequest {
	return r.history[len(r.history)-1]
}

func (r Request) Url() string {
	return r.url
}

func (r Request) Header() Header {
	return r.header
}

func (r Request) Params() Params {
	return r.params
}

func (r Request) Query() Query {
	return r.query
}

func (r Request) Body() Body {
	return r.body
}

func (r Request) Eval() map[string]any {
	return map[string]any{
		"header": r.header,
		"params": r.params,
		"query":  r.query,
		"body":   r.body.Interface(),
	}
}
