package vo

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"io"
	"net/http"
	"net/url"
)

type HttpBackendRequest struct {
	// host represents the host of the backend httpRequest.
	// It is a string field that contains the host information used for constructing the URL.
	host string
	// path is a string field that represents the path of the backend httpRequest.
	// It contains the path information used for constructing the URL.
	// The value of path can be modified using the ModifyUrlPath() `method`.
	path UrlPath
	// method is a string field that represents the HTTP method to be used for the backend httpRequest.
	// It contains information about the desired HTTP method, such as GET, POST, PUT, DELETE, etc.
	// The value of the method field can be accessed using the Method() `method`.
	method string
	// header represents the header fields of a backend httpRequest.
	// The value of header can be modified using the ModifyHeaders() `method`.
	header Header
	// query represents the query fields of a backend httpRequest.
	// The value of Params can be modified using the ModifyQuery() `method`.
	query Query
	// body represents the body of a backend httpRequest.
	// The value of body can be modified using the ModifyBody() `method`.
	body *Body
}

func NewHttpBackendRequest(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse,
) *HttpBackendRequest {
	// criamos o path de requisição backend com base no httpRequest e as config backend
	path := newBackendRequestPath(backend, httpRequest, httpResponse)
	// criamos o header de requisição backend com base no httpRequest e as config backend
	header := newBackendRequestHeader(backend, httpRequest, httpResponse)
	// criamos o query de requisição backend com base no httpRequest e as config backend
	query := newBackendRequestQuery(backend, httpRequest, httpResponse)
	// criamos o body de requisição backend com base no httpRequest e as config backend
	body := newBackendRequestBody(backend, httpRequest, httpResponse)
	// montamos o objeto de valor
	return &HttpBackendRequest{
		host:   backend.BalancedHost(),
		path:   path,
		method: backend.Method(),
		header: header,
		query:  query,
		body:   body,
	}
}

// ModifyHeader returns a new httpBackendRequest with the specified header modified.
// The input header is used to replace the existing header of the httpBackendRequest.
// The other fields of the httpBackendRequest remain unchanged.
func (b *HttpBackendRequest) ModifyHeader(header Header) *HttpBackendRequest {
	return &HttpBackendRequest{
		host:   b.host,
		path:   b.path,
		method: b.method,
		header: header,
		query:  b.query,
		body:   b.body,
	}
}

// ModifyUrlPath creates a new httpBackendRequest with modified path and params.
// It takes in the path string and params as arguments and returns a new httpBackendRequest
// with the original values for host, method, header, query, and body, but with the modified path and params.
func (b *HttpBackendRequest) ModifyUrlPath(path UrlPath) *HttpBackendRequest {
	return &HttpBackendRequest{
		host:   b.host,
		path:   path,
		method: b.method,
		header: b.header,
		query:  b.query,
		body:   b.body,
	}
}

// ModifyQuery returns a new httpBackendRequest instance with the provided query modified.
// The original httpBackendRequest instance remains unchanged.
func (b *HttpBackendRequest) ModifyQuery(query Query) *HttpBackendRequest {
	return &HttpBackendRequest{
		host:   b.host,
		path:   b.path,
		method: b.method,
		header: b.header,
		query:  query,
		body:   b.body,
	}
}

// ModifyBody returns a new instance of httpBackendRequest with the provided body.
// The new instance has the same values for host, path, method, header, params, and query as the original httpBackendRequest,
// but with the updated body.
func (b *HttpBackendRequest) ModifyBody(bodyVO *Body) *HttpBackendRequest {
	return &HttpBackendRequest{
		host:   b.host,
		path:   b.path,
		method: b.method,
		header: b.header,
		query:  b.query,
		body:   bodyVO,
	}
}

func (b *HttpBackendRequest) Path() UrlPath {
	return b.path
}

func (b *HttpBackendRequest) Url() string {
	return fmt.Sprint(b.host, b.Path().String())
}

func (b *HttpBackendRequest) Params() Params {
	return b.Path().Params()
}

// Method returns the method of the httpBackendRequest instance.
func (b *HttpBackendRequest) Method() string {
	return b.method
}

// Header returns the Header of the httpBackendRequest instance.
func (b *HttpBackendRequest) Header() Header {
	return b.header
}

// Query returns the query of the httpBackendRequest.
func (b *HttpBackendRequest) Query() Query {
	return b.query
}

// RawQuery encodes the query parameters into a string representation.
// It returns the encoded query string that can be appended to the URL.
func (b *HttpBackendRequest) RawQuery() string {
	return url.Values(b.query).Encode()
}

// Body returns the `Body` field of the `httpBackendRequest` instance.
// The `Body` field represents the httpRequest body.
// The `Body` method allows you to access the httpRequest body for further manipulation or inspection.
func (b *HttpBackendRequest) Body() *Body {
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
func (b *HttpBackendRequest) BodyToReadCloser() io.ReadCloser {
	// se ele quer omitir o body da solicitação ou o mesmo tiver vazio retornamos
	if helper.IsNil(b.body) {
		return nil
	}
	// retornamos o valor da interface com os bytes do body
	// todo: aqui podemos futuramente colocar encode de httpRequest customizado
	return io.NopCloser(b.body.Buffer())
}

// NetHttp returns an HTTP httpRequest based on the httpBackendRequest instance.
// It constructs the httpRequest, sets the headers, fills in the queries,
// and returns the created HTTP httpRequest.
// If an error occurs during the construction of the httpRequest, nil and the error are returned.
func (b *HttpBackendRequest) NetHttp(ctx context.Context) (*http.Request, error) {
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

func (b *HttpBackendRequest) Map() any {
	var body any
	if helper.IsNotNil(body) {
		body = b.Body().Interface()
	}
	return map[string]any{
		"header": b.Header(),
		"params": b.Params(),
		"query":  b.Query(),
		"body":   body,
	}
}

func newBackendRequestPath(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) UrlPath {
	// inicializamos o backendRequest para ser utilizado para configurar o path
	backendRequest := backend.Request()

	// instanciamos o path a partir das informações ja existentes no httpRequest
	path := NewUrlPath(backend.Path(), httpRequest.Params())

	// checamos se a config nao esta nil
	if helper.IsNotNil(backendRequest) {
		// rodamos a modificações do path
		for _, modifier := range backendRequest.ParamModifiers() {
			path = path.Modify(&modifier, httpRequest, httpResponse)
		}
	}

	// retornamos o path
	return path
}

func newBackendRequestHeader(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) Header {
	// inicializamos o backendRequest para ser utilizado para configurar o header
	backendRequest := backend.Request()

	// checamos se a config esta nil ou ele quer omitir o header
	if helper.IsNil(backendRequest) {
		// se a config for nil retornamos o header do httpRequest
		return httpRequest.Header()
	} else if backendRequest.OmitHeader() {
		// se ele quer omitir o header, retornamos vazio
		return NewEmptyHeader()
	}

	// instanciamos o header a partir das informações ja existentes no httpRequest
	header := httpRequest.Header()
	// mapeamos seguindo a config
	header = header.Map(backendRequest.HeaderMapper())
	// projetamos seguindo a config
	header = header.Projection(backendRequest.HeaderProjection())

	// rodamos a modificações do header
	for _, modifier := range backendRequest.HeaderModifiers() {
		header = header.Modify(&modifier, httpRequest, httpResponse)
	}

	// retornamos o header
	return header
}

func newBackendRequestQuery(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) Query {
	// inicializamos o backendRequest para ser utilizado para configurar o query
	backendRequest := backend.Request()

	// checamos se a config esta nil ou ele quer omitir o query
	if helper.IsNil(backendRequest) {
		// se a config for nil retornamos o query do httpRequest
		return httpRequest.Query()
	} else if backendRequest.OmitQuery() {
		// se ele quer omitir o query, retornamos vazio
		return NewEmptyQuery()
	}

	// instanciamos o query a partir das informações ja existentes no httpRequest
	query := httpRequest.Query()
	// mapeamos seguindo a config
	query = query.Map(backendRequest.QueryMapper())
	// projetamos seguindo a config
	query = query.Projection(backendRequest.QueryProjection())

	// rodamos a modificações do query
	for _, modifier := range backendRequest.QueryModifiers() {
		query = query.Modify(&modifier, httpRequest, httpResponse)
	}

	// retornamos o query
	return query
}

func newBackendRequestBody(backend *Backend, httpRequest *HttpRequest, httpResponse *HttpResponse) *Body {
	// inicializamos o backendRequest para ser utilizado para configurar o query
	backendRequest := backend.Request()

	// checamos se a config esta nil ou ele quer omitir o body
	if helper.IsNil(backendRequest) {
		// se o body for nil e a config for nil retornamos o body do httpRequest
		return httpRequest.Body()
	} else if helper.IsNil(httpRequest.Body()) || backendRequest.OmitBody() {
		// se o body for nil ou se ele quer omitir o body, retornamos nil
		return nil
	}

	// instanciamos o body a partir das informações ja existentes no httpRequest
	body := httpRequest.Body()
	// mapeamos seguindo a config
	body = body.Map(backendRequest.BodyMapper())
	// projetamos seguindo a config
	body = body.Projection(backendRequest.BodyProjection())

	// rodamos a modificações do body
	for _, modifier := range backendRequest.BodyModifiers() {
		body = body.Modify(&modifier, httpRequest, httpResponse)
	}

	// retornamos o body
	return body
}
