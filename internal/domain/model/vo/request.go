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

// NewRequest creates a new Request with the specified URL, method, header, params, query, and body.
// 'url' is the URL of the request.
// 'method' is the HTTP method of the request.
// 'header' is an instance of Header that represents the request header.
// 'params' is an instance of Params that represents the request parameters.
// 'query' is an instance of Query that represents the query parameters.
// 'body' is the body of the request.
// The function returns a new instance of Request initialized with the provided values.
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

// ModifyHeader creates a new Request from an existing one with modifications to the request header.
// Also modifies the request history by adding a new backendRequestVO at the end.
// 'header' which is an instance of Header will replace the existing request header.
// 'backendRequestVO' is an instance of a backend request object to be changed in the history array.
// The function returns a new instance of Request with the modified header and history.
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

// ModifyParams takes a Params object and a backend request object of type backendRequest.
// It modifies the history of the request by replacing the last element with the provided backendRequest object,
// and updates the params of the original request.
// It returns a new Request with the updated history and params.
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

// ModifyQuery takes a Query and a backendRequest as arguments and modifies the Request's history field.
// It returns a new Request object with updated values.
//
// It allocates a new slice, then copies the values of the original Request's fields into the new one,
// and modifies its history field by updating the last element with the provided backendRequest.
// The query field is also updated with the provided query argument.
//
// The original Request is not modified. All modifications are performed on the new Request object.
//
// Parameters:
// query - A Query object that is used to update the new Request's query field.
// backendRequestVO - A backendRequest object that is used to update the last element of the history field in the new Request.
//
// Returns:
// Request - a new Request object with updated query and history fields.
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

// ModifyBody takes a Body and a backendRequestVO as parameters,
// replaces the last item in the history slice of the Request receiver with backendRequestVO, and
// returns a new Request with the modified history, preserving all other fields.
//
// Parameters:
//
//	body - The new body for the Request.
//	backendRequestVO - The backend request value object to be added in the request history.
//
// Returns:
// A new Request instance with the updated history.
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

// Append is a method for the Request type. It adds a backendRequest to the history of the Request
// and returns the updated Request.
//
// The method receives a backendRequest as an argument.
// It constructs a new Request with the same parameters as the original request (url, method, header, params, query, body),
// but with the backendRequest added to the history.
//
// Parameters:
//
//	backendRequest : The request to append to the history of the Request.
//
// Returns:
//
//	Request - A new Request with the backendRequest added to its history.
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

// CurrentBackendRequest returns the last backendRequest object in the request's history array.
func (r Request) CurrentBackendRequest() backendRequest {
	return r.history[len(r.history)-1]
}

// Url returns the URL of the Request.
// It retrieves the value of the `url` field from the Request struct.
func (r Request) Url() string {
	return r.url
}

// Header returns the request header of a Request.
// It returns an instance of Header.
func (r Request) Header() Header {
	return r.header
}

// Params returns the parameters associated with the request.
// The returned value is of type Params.
func (r Request) Params() Params {
	return r.params
}

// Query returns the query parameter map of the Request.
func (r Request) Query() Query {
	return r.query
}

// Body returns the body of the request.
func (r Request) Body() Body {
	return r.body
}

// Eval is a method of the Request type.
//
// This method returns a map where the keys are strings (specifically "header", "params", "query", "body")
// and the values are of any type. These values correspond to data within the Request: header data,
// parameters, query information, and request body, respectively.
func (r Request) Eval() map[string]any {
	return map[string]any{
		"header": r.header,
		"params": r.params,
		"query":  r.query,
		"body":   r.body.Interface(),
	}
}
