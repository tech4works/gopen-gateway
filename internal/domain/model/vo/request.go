package vo

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/gin-gonic/gin"
	"io"
)

// Request represents an HTTP request and contains information such as the request path, URL, method, header,
// parameters, query parameters, request body, and request history.
type Request struct {
	// path represents the URI of the Request.
	// It is a string field in the Request struct.
	path string
	// url represents the URL of the Request.
	url string
	// method represents the HTTP method of the Request.
	// It is a string field in the Request struct.
	method string
	// header represents the request header of a Request.
	header Header
	// params represents the parameters associated with the request.
	params Params
	// query represents the query parameter map of the Request.
	query Query
	// Body represents the body of an HTTP request.
	// It is a field in the Request struct.
	body *Body
	// history represents the history of backend requests made by the Request object.
	// It is a slice of backendRequest objects.
	history []*backendRequest
}

// NewRequest creates a new Request object from a gin.Context object.
//
// The function extracts the URL path and query parameters from the gin.Context object
// and prepares the URL field by concatenating the path and query parameters.
// It then reads the bytes from the request body and sets up the Request object
// with the necessary fields such as URL, method, header, params, query, and body.
//
// Parameters:
// gin - The gin.Context object from which to extract the request information.
//
// Returns:
// Request - A new Request object with the extracted request information.
func NewRequest(gin *gin.Context) *Request {
	// instanciamos o query VO para obter funções de montagem da url por ele
	query := NewQuery(gin.Request.URL.Query())

	// preparamos a url ordenando as chaves de busca
	url := gin.Request.URL.Path
	if helper.IsNotEmpty(gin.Request.URL.RawQuery) {
		url += "?" + query.Encode()
	}

	// obtemos os bytes da requisição
	bodyBytes, _ := io.ReadAll(gin.Request.Body)
	bodyBuffer := bytes.NewBuffer(bodyBytes)
	gin.Request.Body = io.NopCloser(bodyBuffer)

	// montamos o VO de requisição
	return &Request{
		path:   gin.Request.URL.Path,
		url:    url,
		method: gin.Request.Method,
		header: NewHeader(gin.Request.Header),
		params: NewParams(gin.Params),
		query:  query,
		body:   NewBody(gin.GetHeader("Content-Type"), bodyBuffer),
	}
}

// SetHeader takes a Header object as an argument and returns a new Request with the provided header.
// The other fields of the new Request remain unchanged.
//
// Parameters:
// header - The Header object to be set in the new Request.
//
// Returns:
// Request - A new Request instance with the updated header and the original values for the other fields.
func (r *Request) SetHeader(header Header) *Request {
	return &Request{
		path:    r.path,
		url:     r.url,
		method:  r.method,
		header:  header,
		params:  r.params,
		query:   r.query,
		body:    r.body,
		history: r.history,
	}
}

// ModifyHeader creates a new Request from an existing one with modifications to the request header.
// Also modifies the request history by adding a new backendRequestVO at the end.
// 'header' which is an instance of Header will replace the existing request header.
// 'backendRequestVO' is an instance of a backend request object to be changed in the history array.
// The function returns a new instance of Request with the modified header and history.
func (r *Request) ModifyHeader(header Header, backendRequestVO *backendRequest) *Request {
	history := r.history
	history[len(history)-1] = backendRequestVO

	return &Request{
		path:    r.path,
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
func (r *Request) ModifyParams(params Params, backendRequestVO *backendRequest) *Request {
	history := r.history
	history[len(history)-1] = backendRequestVO

	return &Request{
		path:    r.path,
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
func (r *Request) ModifyQuery(query Query, backendRequestVO *backendRequest) *Request {
	history := r.history
	history[len(history)-1] = backendRequestVO

	return &Request{
		path:    r.path,
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
func (r *Request) ModifyBody(body *Body, backendRequestVO *backendRequest) *Request {
	history := r.history
	history[len(history)-1] = backendRequestVO
	return &Request{
		path:    r.path,
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
func (r *Request) Append(backendRequest *backendRequest) *Request {
	return &Request{
		path:    r.path,
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
func (r *Request) CurrentBackendRequest() *backendRequest {
	return r.history[len(r.history)-1]
}

// Url returns the URL of the Request.
// It retrieves the value of the `url` field from the Request struct.
func (r *Request) Url() string {
	return r.url
}

// Path returns the URI of the Request.
// It retrieves the value of the `path` field from the Request struct.
func (r *Request) Path() string {
	return r.path
}

// Method returns the HTTP method of the Request.
// It retrieves the value of the `method` field from the Request struct.
func (r *Request) Method() string {
	return r.method
}

// Header returns the request header of a Request.
// It returns an instance of Header.
func (r *Request) Header() Header {
	return r.header
}

// Params returns the parameters associated with the request.
// The returned value is of type Params.
func (r *Request) Params() Params {
	return r.params
}

// Query returns the query parameter map of the Request.
func (r *Request) Query() Query {
	return r.query
}

// Body returns the body of the request.
func (r *Request) Body() *Body {
	return r.body
}

// Eval is a method of the Request type.
//
// This method returns a map where the keys are strings (specifically "header", "params", "query", "body")
// and the values are of any type. These values correspond to data within the Request: header data,
// parameters, query information, and request body, respectively.
func (r *Request) Eval() string {
	mapEval := map[string]any{
		"header": r.header,
		"params": r.params,
		"query":  r.query,
		"body":   r.body.String(),
	}
	return helper.SimpleConvertToString(mapEval)
}
