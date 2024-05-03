package vo

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"io"
	"net/http"
	"net/url"
)

// httpBackendRequest represents a httpRequest to be made to a backend server.
// It includes the host information, method, path, modifyHeaders, params, query fields, and the httpRequest body.
type httpBackendRequest struct {
	// host represents the host of the backend httpRequest.
	// It is a string field that contains the host information used for constructing the URL.
	host string
	// path is a string field that represents the path of the backend httpRequest.
	// It contains the path information used for constructing the URL.
	// The value of path can be modified using the ModifyParams() `method`.
	path UrlPath
	// method is a string field that represents the HTTP method to be used for the backend httpRequest.
	// It contains information about the desired HTTP method, such as GET, POST, PUT, DELETE, etc.
	// The value of the method field can be accessed using the Method() `method`.
	method string
	// header represents the header fields of a backend httpRequest.
	// The value of header can be modified using the ModifyHeaders() `method`.
	header Header
	// params represents the params fields of a backend httpRequest.
	// The value of params can be modified using the ModifyParams() `method`.
	params Params
	// query represents the query fields of a backend httpRequest.
	// The value of params can be modified using the ModifyQuery() `method`.
	query Query
	// body represents the body of a backend httpRequest.
	// The value of body can be modified using the ModifyBody() `method`.
	body *Body
}

func NewHttpBackendRequest(backendVO *Backend, balancedHost string, requestVO *HttpRequest) *httpBackendRequest {
	// inicializamos o backendRequest para ser utilizado para configurar o http httpRequest
	backendRequestVO := backendVO.Request()

	// inicializamos o header
	var header Header
	// inicializamos o query
	var query Query
	// inicializamos o body
	var body *Body

	// verificamos se o mesmo não quer ser omitido
	if helper.IsNotNil(backendRequestVO) {
		if !backendRequestVO.OmitHeader() {
			header = requestVO.Header().FilterByRequest(backendRequestVO.HeaderFilter())
		}
		if !backendRequestVO.OmitQuery() {
			query = requestVO.Query().Filter(backendRequestVO.QueryFilter())
		}
		if !backendRequestVO.OmitBody() {
			body = requestVO.Body().Filter(backendRequestVO.BodyFilter())
		}
	}

	// inicializamos os params
	params := NewParamsByUrlPath(backendVO.Path(), requestVO.Params())

	// montamos o objeto de valor
	return &httpBackendRequest{
		host:   balancedHost,
		path:   backendVO.Path(),
		method: backendVO.Method(),
		header: header,
		params: params,
		query:  query,
		body:   body,
	}
}

// ModifyHeader returns a new httpBackendRequest with the specified header modified.
// The input header is used to replace the existing header of the httpBackendRequest.
// The other fields of the httpBackendRequest remain unchanged.
func (b *httpBackendRequest) ModifyHeader(header Header) *httpBackendRequest {
	return &httpBackendRequest{
		host:   b.host,
		path:   b.path,
		method: b.method,
		header: header,
		params: b.params,
		query:  b.query,
		body:   b.body,
	}
}

// ModifyParams creates a new httpBackendRequest with modified path and params.
// It takes in the path string and params Params as arguments and returns a new httpBackendRequest
// with the original values for host, method, header, query, and body, but with the modified path and params.
func (b *httpBackendRequest) ModifyParams(path UrlPath, params Params) *httpBackendRequest {
	return &httpBackendRequest{
		host:   b.host,
		path:   path,
		method: b.method,
		header: b.header,
		params: params,
		query:  b.query,
		body:   b.body,
	}
}

// ModifyQuery returns a new httpBackendRequest instance with the provided query modified.
// The original httpBackendRequest instance remains unchanged.
func (b *httpBackendRequest) ModifyQuery(query Query) *httpBackendRequest {
	return &httpBackendRequest{
		host:   b.host,
		path:   b.path,
		method: b.method,
		header: b.header,
		params: b.params,
		query:  query,
		body:   b.body,
	}
}

// ModifyBody returns a new instance of httpBackendRequest with the provided body.
// The new instance has the same values for host, path, method, header, params, and query as the original httpBackendRequest,
// but with the updated body.
func (b *httpBackendRequest) ModifyBody(bodyVO *Body) *httpBackendRequest {
	return &httpBackendRequest{
		host:   b.host,
		path:   b.path,
		method: b.method,
		header: b.header,
		params: b.params,
		query:  b.query,
		body:   bodyVO,
	}
}

// Path returns the path string of the httpBackendRequest instance.
func (b *httpBackendRequest) Path() UrlPath {
	return b.path
}

// Url returns the fully constructed URL for the httpBackendRequest instance.
// It replaces any path parameters in the path with the corresponding values from the params map.
// It then concatenates the host and modified path to construct the URL.
// The returned URL is a string representation of the complete URL for the httpRequest.
func (b *httpBackendRequest) Url() string {
	// aqui vamos dar o replace nos keys /users/:key para /users/2
	path := b.path
	for key, value := range b.params {
		if path.ContainsParam(key) {
			path = path.FillParamValue(key, value)
		}
	}
	// retornamos a url para req
	return fmt.Sprint(b.host, path.String())
}

// Method returns the method of the httpBackendRequest instance.
func (b *httpBackendRequest) Method() string {
	return b.method
}

// Header returns the Header of the httpBackendRequest instance.
func (b *httpBackendRequest) Header() Header {
	return b.header
}

// Params returns the Params field of the httpBackendRequest instance.
func (b *httpBackendRequest) Params() Params {
	return b.params
}

// Query returns the query of the httpBackendRequest.
func (b *httpBackendRequest) Query() Query {
	return b.query
}

// RawQuery encodes the query parameters into a string representation.
// It returns the encoded query string that can be appended to the URL.
func (b *httpBackendRequest) RawQuery() string {
	return url.Values(b.query).Encode()
}

// Body returns the `Body` field of the `httpBackendRequest` instance.
// The `Body` field represents the httpRequest body.
// The `Body` method allows you to access the httpRequest body for further manipulation or inspection.
func (b *httpBackendRequest) Body() *Body {
	return b.body
}

// BodyToReadCloser returns the body to send as an `io.ReadCloser` interface.
// If `omitRequestBody` is set to `true` or `body` is `nil`, it returns `nil`.
//
// It converts the body to bytes using the desired encoding (XML, JSON, TEXT/PLAIN) based on `Content-Type` config.
//
// If there is an error during the conversion, it returns `nil`.
//
// Finally, it returns the `io.ReadCloser` interface with the bytes of the body.
func (b *httpBackendRequest) BodyToReadCloser() io.ReadCloser {
	// se ele quer omitir o body da solicitação ou o mesmo tiver vazio retornamos
	if helper.IsNil(b.body) {
		return nil
	}
	// retornamos o valor da interface com os bytes do body
	// todo: aqui podemos futuramente colocar encode de httpRequest customizado
	return io.NopCloser(b.body.Value())
}

// Http returns an HTTP httpRequest based on the httpBackendRequest instance.
// It constructs the httpRequest, sets the headers, fills in the queries,
// and returns the created HTTP httpRequest.
// If an error occurs during the construction of the httpRequest, nil and the error are returned.
func (b *httpBackendRequest) Http(ctx context.Context) (*http.Request, error) {
	// construímos o http httpRequest para fazer a requisição
	httpRequest, err := http.NewRequestWithContext(ctx, b.Method(), b.Url(), b.BodyToReadCloser())
	if helper.IsNotNil(err) {
		return nil, err
	}

	// preenchemos o header com o httpBackendRequest montado
	httpRequest.Header = b.Header().Http()
	// preenchemos as queries com o httpBackendRequest montado
	httpRequest.URL.RawQuery = b.RawQuery()

	// retornamos o http httpRequest criado
	return httpRequest, nil
}

// Eval returns a map[string]any with the evaluated values of the httpBackendRequest instance:
//   - "header": The header field of the httpBackendRequest.
//   - "params": The params field of the httpBackendRequest.
//   - "query": The query field of the httpBackendRequest.
//   - "body": The Interface method of the body field of the httpBackendRequest.
//
// The map is then returned.
func (b *httpBackendRequest) Eval() any {
	var evalBody any
	if helper.IsNotNil(evalBody) {
		evalBody = b.body.Interface()
	}
	return map[string]any{
		"header": b.header,
		"params": b.params,
		"query":  b.query,
		"body":   evalBody,
	}
}
