package vo

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Backend is a type that represents a backend server configuration.
type Backend struct {
	// name represents the name of the Backend instance.
	name string
	// hosts is an array of host addresses.
	hosts []string
	// path is a string that represents the path of the backend server configuration.
	path string
	// method is the HTTP method to be used for requests to the backend.
	method string
	// forwardHeaders is a slice of strings representing the modifyHeaders to be forwarded.
	forwardHeaders []string
	// forwardQueries is a slice of strings representing the query parameters to be forwarded.
	forwardQueries []string
	// modifiers is an instance of BackendModifiers containing modifiers for the backend request and response.
	modifiers *BackendModifiers
	// extraConfig is an instance of BackendExtraConfig containing extra configuration options for the backend.
	extraConfig *BackendExtraConfig
}

// BackendModifiers is a type that represents the set of modifiers for a backend configuration.
// It contains fields for the status code, header, params, query, and body modifiers.
type BackendModifiers struct {
	// statusCode represents the status code modifier for a BackendModifiers instance.
	// It is an integer value that specifies the desired status code for a response.
	statusCode int
	// header represents an array of Modifier instances that modify the modifyHeaders of a request or response
	// from Endpoint or only current backend.
	header []Modifier
	// param is a field in the BackendModifiers struct.
	// It represents an array of Modifier instances that modify the parameters of a request from Endpoint or only current
	// backend.
	param []Modifier
	// `query` is a field in the `BackendModifiers` struct. It represents an array of `Modifier` instances that
	// modify the query parameters of a request from `Endpoint` or only the current backend.
	query []Modifier
	// body represents an array of Modifier instances that modify the body of a request or response
	// from Endpoint or only current backend.
	body []Modifier
}

// BackendExtraConfig is a type that represents additional configuration options for a backend in the Gopen application.
type BackendExtraConfig struct {
	// groupResponse is a boolean flag indicating whether the backend should group response.
	// The default value is false.
	groupResponse bool
	// omitRequestBody represents a boolean flag indicating whether the backend should omit the request body in request.
	// If set to true, the backend will not include the request body in the request.
	// If set to false, the request body will be included in the request. The default value is false.
	omitRequestBody bool
	// omitResponse represents a boolean flag indicating whether the backend should omit the response in the incoming request.
	// If set to true, the backend will not include the response in the incoming request.
	// If set to false, the response will be included in the incoming request.
	// The default value is false.
	omitResponse bool
}

// backendRequest represents a request to be made to a backend server.
// It includes the host information, method, path, modifyHeaders, params, query fields, and the request body.
type backendRequest struct {
	// omitBody is a boolean field that determines whether the request body should be omitted.
	// If set to true, the request body will not be included in the backend request.
	// If set to false, the request body will be included in the backend request.
	// The default value is false.
	omitBody bool
	// host represents the host of the backend request.
	// It is a string field that contains the host information used for constructing the URL.
	host string
	// path is a string field that represents the path of the backend request.
	// It contains the path information used for constructing the URL.
	// The value of path can be modified using the ModifyParams() `method`.
	path string
	// method is a string field that represents the HTTP method to be used for the backend request.
	// It contains information about the desired HTTP method, such as GET, POST, PUT, DELETE, etc.
	// The value of the method field can be accessed using the Method() `method`.
	method string
	// header represents the header fields of a backend request.
	// The value of header can be modified using the ModifyHeaders() `method`.
	header Header
	// params represents the params fields of a backend request.
	// The value of params can be modified using the ModifyParams() `method`.
	params Params
	// query represents the query fields of a backend request.
	// The value of params can be modified using the ModifyQuery() `method`.
	query Query
	// body represents the body of a backend request.
	// The value of body can be modified using the ModifyBody() `method`.
	body *Body
}

// backendResponse represents a response from a backend service.
type backendResponse struct {
	// name represents the name of the Backend instance.
	name string
	// omit represents a boolean flag that indicates whether to omit the response in the incoming request.
	omit bool
	// group is a boolean flag indicating whether they should group response body.
	group bool
	// statusCode represents HTTP statusCode of a backend response.
	// The value of statusCode can be modified using the ModifyStatusCode() `method`.
	statusCode int
	// header represents the body fields of a backend response.
	// The value of header can be modified using the ModifyHeader() `method`.
	header Header
	// body represents the body of a backend response.
	// The value of body can be modified using the ModifyBody() `method`.
	body *Body
}

// NewBackendRequest creates a new instance of backendRequest based on the provided parameters.
// It initializes the header to be used in the construction of the filtered VO by forward-headers.
// It initializes the query to be used in the construction of the filtered VO by forward-queries.
// It initializes the params using the NewParamsByPath function, passing the path and requestVO.params parameters.
// It constructs the backendRequest object and returns it.
func NewBackendRequest(backendVO *Backend, balancedHost string, requestVO *Request) *backendRequest {
	// inicializamos o header a ser utilizado na construção do VO filtrado pelo forward-headers
	header := requestVO.Header().FilterByForwarded(backendVO.forwardHeaders)

	// inicializamos a query a ser utilizado na construção do VO filtrada pelo forward-queries
	query := requestVO.Query().FilterByForwarded(backendVO.forwardQueries)

	// inicializamos os params
	params := NewParamsByPath(backendVO.path, requestVO.params)

	// inicializamos o omitRequestBody como false
	var omitBody bool
	if helper.IsNotNil(backendVO.ExtraConfig()) {
		omitBody = backendVO.ExtraConfig().OmitRequestBody()
	}

	// montamos o objeto de valor
	return &backendRequest{
		omitBody: omitBody,
		host:     balancedHost,
		path:     backendVO.Path(),
		method:   backendVO.Method(),
		header:   header,
		params:   params,
		query:    query,
		body:     requestVO.Body(),
	}
}

// NewBackendResponse creates a new instance of backendResponse based on the provided parameters.
// It parses the bytes of the response body into an interface.
// It converts the bytes and content-type into a body VO.
// It constructs the backendResponse object and returns it.
func NewBackendResponse(backendVO *Backend, httpResponse *http.Response) *backendResponse {
	// fazemos o parse dos bytes da resposta em para uma interface
	bodyBytes, _ := io.ReadAll(httpResponse.Body)

	// convertemos em body VO a partir dos bytes e do content-type
	body := NewBody(httpResponse.Header.Get("Content-Type"), bytes.NewBuffer(bodyBytes))

	// instanciamos o omit e group
	var omit bool
	var group bool

	// se tiver extraConfig preenchemos os valores
	if helper.IsNotNil(backendVO.ExtraConfig()) {
		omit = backendVO.ExtraConfig().OmitResponse()
		group = backendVO.ExtraConfig().GroupResponse()
	}

	// construímos o objeto de valor do backend response
	return &backendResponse{
		name:       backendVO.Name(),
		omit:       omit,
		group:      group,
		statusCode: httpResponse.StatusCode,
		header:     NewHeader(httpResponse.Header),
		body:       body,
	}
}

// newBackend creates a new Backend instance based on the provided backendDTO.
// It takes the backendDTO fields and assigns them to the corresponding fields in the Backend struct.
// It also creates a new BackendModifiers instance by calling the newBackendModifier function,
// passing the modifiers field from the backendDTO.
// It creates a new BackendExtraConfig instance by calling the newBackendExtraConfig function,
// passing the extraConfig field from the backendDTO.
// The function returns the created Backend instance.
func newBackend(backendDTO dto.Backend) Backend {
	return Backend{
		name:           backendDTO.Name,
		hosts:          backendDTO.Hosts,
		path:           backendDTO.Path,
		method:         backendDTO.Method,
		forwardHeaders: backendDTO.ForwardHeaders,
		forwardQueries: backendDTO.ForwardQueries,
		modifiers:      newBackendModifier(backendDTO.Modifiers),
		extraConfig:    newBackendExtraConfig(backendDTO.ExtraConfig),
	}
}

// newMiddlewareBackend creates a new Backend instance based on the provided backendVO and backendExtraConfigVO.
// It takes the fields from backendVO and assigns them to the corresponding fields in the Backend struct.
// It assigns the backendExtraConfigVO parameter to the extraConfig field of the Backend struct.
// The function returns the created Backend instance.
func newMiddlewareBackend(backendVO *Backend, backendExtraConfigVO *BackendExtraConfig) Backend {
	return Backend{
		name:           backendVO.name,
		hosts:          backendVO.hosts,
		path:           backendVO.path,
		method:         backendVO.method,
		forwardHeaders: backendVO.forwardHeaders,
		forwardQueries: backendVO.forwardQueries,
		modifiers:      backendVO.modifiers,
		extraConfig:    backendExtraConfigVO,
	}
}

// newBackendModifier creates a new instance of BackendModifiers based on the provided backendModifierDTO.
// If the backendModifierDTO is nil, it returns nil.
// Otherwise, it initializes a new BackendModifiers object and populates its fields with the values from the
// backendModifierDTO. The header, params, query, and body fields of the BackendModifiers struct are populated from the
// corresponding fields in the backendModifierDTO. It uses the newModifier function to create a new Modifier for each
// element in the backendModifierDTO slice. The newModifier function initializes the context, scope, action, propagate,
// key, and value fields of the Modifier struct. Once all the Modifiers are created, the BackendModifiers object is
// returned. The StatusCode field of the BackendModifiers struct is set to the value of the StatusCode field in the
// backendModifierDTO.
func newBackendModifier(backendModifierDTO *dto.BackendModifiers) *BackendModifiers {
	if helper.IsNil(backendModifierDTO) {
		return nil
	}

	var header []Modifier
	for _, modifierDTO := range backendModifierDTO.Header {
		header = append(header, *newModifier(&modifierDTO))
	}
	var params []Modifier
	for _, modifierDTO := range backendModifierDTO.Param {
		params = append(params, *newModifier(&modifierDTO))
	}
	var query []Modifier
	for _, modifierDTO := range backendModifierDTO.Query {
		query = append(query, *newModifier(&modifierDTO))
	}
	var body []Modifier
	for _, modifierDTO := range backendModifierDTO.Body {
		body = append(body, *newModifier(&modifierDTO))
	}

	return &BackendModifiers{
		statusCode: backendModifierDTO.StatusCode,
		header:     header,
		param:      params,
		query:      query,
		body:       body,
	}
}

func newBackendExtraConfig(extraConfigDTO *dto.BackendExtraConfig) *BackendExtraConfig {
	if helper.IsNil(extraConfigDTO) {
		return nil
	}
	return &BackendExtraConfig{
		groupResponse:   extraConfigDTO.GroupResponse,
		omitRequestBody: extraConfigDTO.OmitRequestBody,
		omitResponse:    extraConfigDTO.OmitResponse,
	}
}

// Name returns the name of the Backend instance.
func (b *Backend) Name() string {
	return b.name
}

// Hosts returns the host array of the Backend instance.
func (b *Backend) Hosts() []string {
	return b.hosts
}

// BalancedHost returns a balanced host from the Backend instance. If there is only one host, it is returned directly.
// Otherwise, a random host is selected from the available ones.
func (b *Backend) BalancedHost() string {
	// todo: aqui obtemos o host (correto é criar um domínio chamado balancer aonde ele vai retornar o host
	//  disponível pegando como base, se ele esta de pé ou não, e sua config de porcentagem)
	if helper.EqualsLen(b.hosts, 1) {
		return b.hosts[0]
	}
	return b.hosts[helper.RandomNumber(0, len(b.hosts)-1)]
}

// Path returns the path of the Backend instance.
func (b *Backend) Path() string {
	return b.path
}

// Method returns the method of the Backend instance.
func (b *Backend) Method() string {
	return b.method
}

// ForwardHeaders returns the slice of strings representing the forward headers of the Backend instance.
func (b *Backend) ForwardHeaders() []string {
	return b.forwardHeaders
}

// ForwardQueries returns the slice of forward queries of the Backend instance.
func (b *Backend) ForwardQueries() []string {
	return b.forwardQueries
}

// BackendModifiers returns the BackendModifiers instance associated with the Backend.
// BackendModifiers contains methods to access and modify the status code, header, params,
// query, and body of the Backend.
//
// Note: This method returns a copy of the BackendModifiers instance, any modifications made to it will not affect the
// original Backend instance.
func (b *Backend) BackendModifiers() *BackendModifiers {
	return b.modifiers
}

// ExtraConfig returns the extra configuration options for the Backend instance.
// It returns an instance of BackendExtraConfig that contains additional configuration options
// such as grouping response, omitting request body, and omitting response.
// This method returns a copy of the BackendExtraConfig instance, any modifications made to it will not affect the
// original Backend instance.
func (b *Backend) ExtraConfig() *BackendExtraConfig {
	return b.extraConfig
}

// CountModifiers returns the number of modifiers present in the Backend instance.
// If the modifiers field is not nil, it counts all the modifiers using the CountAll() method of BackendModifiers.
// Otherwise, it returns 0.
func (b *Backend) CountModifiers() int {
	if helper.IsNotNil(b.modifiers) {
		return b.modifiers.CountAll()
	}
	return 0
}

// StatusCode returns the status code Modifier of the BackendModifiers instance.
func (b *BackendModifiers) StatusCode() int {
	return b.statusCode
}

// Header returns the header modifiers of the BackendModifiers instance.
func (b *BackendModifiers) Header() []Modifier {
	return b.header
}

// Param returns an array of Modifier instances that represent the params of the BackendModifiers instance.
func (b *BackendModifiers) Param() []Modifier {
	return b.param
}

// Query returns the list of modifiers for the query of the BackendModifiers instance.
func (b *BackendModifiers) Query() []Modifier {
	return b.query
}

// Body returns the list of modifiers for the body of the BackendModifiers instance.
func (b *BackendModifiers) Body() []Modifier {
	return b.body
}

// CountAll returns the total count of modifiers for a BackendModifiers instance.
// It counts the number of valid `statusCode` and the length of `header`, `params`, `query`, and `body` slices,
// and adds them up to get the total count.
func (b *BackendModifiers) CountAll() (count int) {
	if helper.IsNotNil(b.statusCode) {
		count++
	}
	count += len(b.header) + len(b.param) + len(b.query) + len(b.body)
	return count
}

// GroupResponse returns a boolean flag indicating whether the backend should group response.
// The default value is false.
func (b *BackendExtraConfig) GroupResponse() bool {
	return b.groupResponse
}

// OmitRequestBody returns a boolean flag indicating whether the backend should omit the request body in request.
// If set to true, the backend will not include the request body in the request.
// If set to false, the request body will be included in the request. The default value is false.
func (b *BackendExtraConfig) OmitRequestBody() bool {
	return b.omitRequestBody
}

// OmitResponse returns a boolean flag indicating whether the backend should omit the response body in the incoming request.
// If set to true, the backend will not include the response body in the incoming request.
// If set to false, the response body will be included in the incoming request.
// The default value is false.
func (b *BackendExtraConfig) OmitResponse() bool {
	return b.omitResponse
}

// ModifyHeader returns a new backendRequest with the specified header modified.
// The input header is used to replace the existing header of the backendRequest.
// The other fields of the backendRequest remain unchanged.
func (b *backendRequest) ModifyHeader(header Header) *backendRequest {
	return &backendRequest{
		host:   b.host,
		path:   b.path,
		method: b.method,
		header: header,
		params: b.params,
		query:  b.query,
		body:   b.body,
	}
}

// ModifyParams creates a new backendRequest with modified path and params.
// It takes in the path string and params Params as arguments and returns a new backendRequest
// with the original values for host, method, header, query, and body, but with the modified path and params.
func (b *backendRequest) ModifyParams(path string, params Params) *backendRequest {
	return &backendRequest{
		host:   b.host,
		path:   path,
		method: b.method,
		header: b.header,
		params: params,
		query:  b.query,
		body:   b.body,
	}
}

// ModifyQuery returns a new backendRequest instance with the provided query modified.
// The original backendRequest instance remains unchanged.
func (b *backendRequest) ModifyQuery(query Query) *backendRequest {
	return &backendRequest{
		host:   b.host,
		path:   b.path,
		method: b.method,
		header: b.header,
		params: b.params,
		query:  query,
		body:   b.body,
	}
}

// ModifyBody returns a new instance of backendRequest with the provided body.
// The new instance has the same values for host, path, method, header, params, and query as the original backendRequest,
// but with the updated body.
func (b *backendRequest) ModifyBody(body *Body) *backendRequest {
	return &backendRequest{
		host:   b.host,
		path:   b.path,
		method: b.method,
		header: b.header,
		params: b.params,
		query:  b.query,
		body:   body,
	}
}

// Host returns the value of the `host` attribute of the backendRequest instance.
func (b *backendRequest) Host() string {
	return b.host
}

// Path returns the path string of the backendRequest instance.
func (b *backendRequest) Path() string {
	return b.path
}

// Url constructs the URL string for the backend request instance.
// It replaces the path placeholders with their corresponding values in the params map.
// The final URL is formed by concatenating the host and the modified path.
func (b *backendRequest) Url() string {
	// aqui vamos dar o replace nos keys /users/:key para /users/2
	path := b.path
	for key, value := range b.params {
		paramUrl := fmt.Sprintf(":%s", key)
		if helper.Contains(b.path, paramUrl) {
			path = strings.ReplaceAll(path, paramUrl, value)
		}
	}

	// retornamos a url para req
	return fmt.Sprint(b.host, path)
}

func (b *backendRequest) Method() string {
	return b.method
}

// Header returns the Header of the backendRequest instance.
func (b *backendRequest) Header() Header {
	return b.header
}

// Params returns the Params field of the backendRequest instance.
func (b *backendRequest) Params() Params {
	return b.params
}

// Query returns the query of the backendRequest.
func (b *backendRequest) Query() Query {
	return b.query
}

// RawQuery encodes the query parameters into a string representation.
// It returns the encoded query string that can be appended to the URL.
func (b *backendRequest) RawQuery() string {
	return url.Values(b.query).Encode()
}

// Body returns the `Body` field of the `backendRequest` instance.
// The `Body` field represents the request body.
// The `Body` method allows you to access the request body for further manipulation or inspection.
func (b *backendRequest) Body() *Body {
	return b.body
}

// BodyToRead returns the body to send as an `io.ReadCloser` interface.
// If `omitRequestBody` is set to `true` or `body` is `nil`, it returns `nil`.
//
// It converts the body to bytes using the desired encoding (XML, JSON, TEXT/PLAIN) based on `Content-Type` config.
//
// If there is an error during the conversion, it returns `nil`.
//
// Finally, it returns the `io.ReadCloser` interface with the bytes of the body.
func (b *backendRequest) BodyToRead() io.ReadCloser {
	// se ele quer omitir o body da solicitação ou o mesmo tiver vazio retornamos
	if b.omitBody || helper.IsNil(b.body) {
		return nil
	}
	// retornamos o valor da interface com os bytes do body
	// todo: aqui podemos futuramente colocar encode de request customizado
	return io.NopCloser(b.body.Value())
}

// Http returns an HTTP request based on the backendRequest instance.
// It constructs the request, sets the headers, fills in the queries,
// and returns the created HTTP request.
// If an error occurs during the construction of the request, nil and the error are returned.
func (b *backendRequest) Http(ctx context.Context) (*http.Request, error) {
	// construímos o http request para fazer a requisição
	httpRequest, err := http.NewRequestWithContext(ctx, b.Method(), b.Url(), b.BodyToRead())
	if helper.IsNotNil(err) {
		return nil, err
	}

	// preenchemos o header com o backendRequest montado
	httpRequest.Header = b.Header().Http()
	// preenchemos as queries com o backendRequest montado
	httpRequest.URL.RawQuery = b.RawQuery()

	// retornamos o http request criado
	return httpRequest, nil
}

// ModifyStatusCode returns a new instance of backendResponse with the given statusCode modified.
// The method creates a copy of the original backendResponse and sets the statusCode to the provided value.
// The other fields are copied from the original backendResponse.
func (b *backendResponse) ModifyStatusCode(statusCode int) *backendResponse {
	return &backendResponse{
		name:       b.name,
		omit:       b.omit,
		group:      b.group,
		statusCode: statusCode,
		header:     b.header,
		body:       b.body,
	}
}

// ModifyHeader returns a new instance of backendResponse with the given header modified.
// The method creates a copy of the original backendResponse and sets the header to the provided value.
// The other fields are copied from the original backendResponse.
func (b *backendResponse) ModifyHeader(header Header) *backendResponse {
	return &backendResponse{
		name:       b.name,
		omit:       b.omit,
		group:      b.group,
		statusCode: b.statusCode,
		header:     header,
		body:       b.body,
	}
}

// ModifyBody returns a new instance of backendResponse with the given body modified.
// The method creates a copy of the original backendResponse and sets the body to the provided value.
// The other fields are copied from the original backendResponse.
func (b *backendResponse) ModifyBody(body *Body) *backendResponse {
	return &backendResponse{
		name:       b.name,
		omit:       b.omit,
		group:      b.group,
		statusCode: b.statusCode,
		header:     b.header,
		body:       body,
	}
}

func (b *backendResponse) Ok() bool {
	return helper.IsGreaterThanOrEqual(b.statusCode, 200) && helper.IsLessThanOrEqual(b.statusCode, 299)
}

// Key returns the key of the backendResponse for aggregation.
// The key is composed of the string "backend" and the index, if it is greater than or equal to zero.
// If the backendResponse has a name, the key is set to the name.
func (b *backendResponse) Key(index int) (key string) {
	// montamos o key do backend para agregar
	key = "backend"
	if helper.IsGreaterThanOrEqual(index, 0) {
		key = fmt.Sprintf("%s-%v", key, index)
	}
	// se o backend tiver nome, damos prioridade
	if helper.IsNotEmpty(b.name) {
		key = b.name
	}
	return key
}

// StatusCode returns the `statusCode` of the `backendResponse` instance.
func (b *backendResponse) StatusCode() int {
	return b.statusCode
}

// Header returns the `header` of the `backendResponse` instance.
func (b *backendResponse) Header() Header {
	return b.header
}

// Body returns the `body` of the `backendResponse` instance.
func (b *backendResponse) Body() *Body {
	return b.body
}

// GroupResponseByType returns true if the body of the backendResponse is text or if the value of the body is a slice.
// Otherwise, it returns false.
func (b *backendResponse) GroupResponseByType() bool {
	body := b.Body()
	bodyValue := body.Bytes()
	return body.IsText() || helper.IsSlice(bodyValue)
}

// GroupResponse returns true if the backendResponse instance should be grouped, either by setting the groupResponse
// field to true or by the value of the body being a text or a slice. Otherwise, it returns false.
func (b *backendResponse) GroupResponse() bool {
	return b.group || b.GroupResponseByType()
}
