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

type Backend struct {
	name           string
	hosts          []string
	path           string
	method         string
	forwardHeaders []string
	forwardQueries []string
	modifiers      BackendModifiers
	extraConfig    backendExtraConfig
}

type BackendModifiers struct {
	statusCode Modifier
	header     []Modifier
	params     []Modifier
	query      []Modifier
	body       []Modifier
}

type backendExtraConfig struct {
	groupResponse   bool
	omitRequestBody bool
	omitResponse    bool
}

type backendRequest struct {
	omitRequestBody bool
	host            string
	path            string
	method          string
	header          Header
	params          Params
	query           Query
	body            Body
}

type backendResponse struct {
	name          string
	omitResponse  bool
	groupResponse bool
	statusCode    int
	header        Header
	body          Body
}

// NewBackendRequest creates a new instance of backendRequest based on the provided parameters.
// It initializes the header to be used in the construction of the filtered VO by forward-headers.
// It initializes the query to be used in the construction of the filtered VO by forward-queries.
// It initializes the params using the NewParamsByPath function, passing the path and requestVO.params parameters.
// It constructs the backendRequest object and returns it.
func NewBackendRequest(backendVO Backend, balancedHost string, requestVO Request) backendRequest {
	// inicializamos o header a ser utilizado na construção do VO filtrado pelo forward-headers
	header := requestVO.Header().FilterByForwarded(backendVO.forwardHeaders)

	// inicializamos a query a ser utilizado na construção do VO filtrada pelo forward-queries
	query := requestVO.Query().FilterByForwarded(backendVO.forwardQueries)

	// inicializamos os params
	params := NewParamsByPath(backendVO.path, requestVO.params)

	// montamos o objeto de valor
	return backendRequest{
		omitRequestBody: backendVO.extraConfig.omitRequestBody,
		host:            balancedHost,
		path:            backendVO.path,
		method:          backendVO.method,
		header:          header,
		params:          params,
		query:           query,
		body:            requestVO.Body(),
	}
}

// NewBackendResponse creates a new backendResponse instance based on the provided Backend and httpResponse.
// It parses the bytes of the response into an interface and assigns it to the body field of backendResponse.
// If the response body is not empty, it parses the body bytes into a Body instance, maintaining the JSON order.
// The function initializes the backendResponse instance with the following fields:
// - name: taken from the name field of Backend
// - omitResponse: taken from the omitResponse field of Backend's extraConfig
// - groupResponse: taken from the groupResponse field of Backend's extraConfig
// - statusCode: taken from the StatusCode field of httpResponse
// - header: a new Header instance created from httpResponse's Header
// - body: the parsed Body instance
// The function returns the created backendResponse instance.
func NewBackendResponse(backendVO Backend, httpResponse *http.Response) backendResponse {
	// fazemos o parse dos bytes da resposta em para uma interface
	var body Body
	bodyBytes, _ := io.ReadAll(httpResponse.Body)

	// se não for vazio, fazemos o parse para Body mantendo sempre a ordenação do json
	if helper.IsNotEmpty(bodyBytes) {
		body = newBody(bodyBytes)
	}

	// construímos o objeto de valor do backend response
	return backendResponse{
		name:          backendVO.name,
		omitResponse:  backendVO.extraConfig.omitResponse,
		groupResponse: backendVO.extraConfig.groupResponse,
		statusCode:    httpResponse.StatusCode,
		header:        NewHeader(httpResponse.Header),
		body:          body,
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
		modifiers:      newBackendModifier(helper.IfNilReturns(backendDTO.Modifiers, dto.BackendModifiers{})),
		extraConfig:    newBackendExtraConfig(helper.IfNilReturns(backendDTO.ExtraConfig, dto.BackendExtraConfig{})),
	}
}

// newMiddlewareBackend creates a new Backend instance based on the provided backendVO and backendExtraConfigVO.
// It takes the fields from backendVO and assigns them to the corresponding fields in the Backend struct.
// It assigns the backendExtraConfigVO parameter to the extraConfig field of the Backend struct.
// The function returns the created Backend instance.
func newMiddlewareBackend(backendVO Backend, backendExtraConfigVO backendExtraConfig) Backend {
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

// newBackendModifier creates a new BackendModifiers instance based on the provided backendModifierDTO.
// It initializes the header, params, query, and body fields by iterating over the corresponding lists in the
// backendModifierDTO and calling the newModifier function to create a new Modifier instance for each element in the list.
// The function returns the created BackendModifiers instance.
func newBackendModifier(backendModifierDTO dto.BackendModifiers) BackendModifiers {
	var header []Modifier
	for _, modifierDTO := range backendModifierDTO.Header {
		header = append(header, newModifier(modifierDTO))
	}
	var params []Modifier
	for _, modifierDTO := range backendModifierDTO.Params {
		params = append(params, newModifier(modifierDTO))
	}
	var query []Modifier
	for _, modifierDTO := range backendModifierDTO.Query {
		query = append(query, newModifier(modifierDTO))
	}
	var body []Modifier
	for _, modifierDTO := range backendModifierDTO.Body {
		body = append(body, newModifier(modifierDTO))
	}

	return BackendModifiers{
		statusCode: newModifier(helper.IfNilReturns(backendModifierDTO.StatusCode, dto.Modifier{})),
		header:     header,
		params:     params,
		query:      query,
		body:       body,
	}
}

// newBackendExtraConfig creates a new `backendExtraConfig` instance based on the provided `extraConfigDTO`.
// It takes the `extraConfigDTO` fields and assigns them to the corresponding fields in the `backendExtraConfig` struct.
// The function returns the created `backendExtraConfig` instance.
func newBackendExtraConfig(extraConfigDTO dto.BackendExtraConfig) backendExtraConfig {
	return backendExtraConfig{
		groupResponse:   extraConfigDTO.GroupResponse,
		omitRequestBody: extraConfigDTO.OmitRequestBody,
		omitResponse:    extraConfigDTO.OmitResponse,
	}
}

// Name returns the name of the Backend instance.
func (b Backend) Name() string {
	return b.name
}

// Hosts returns the host array of the Backend instance.
func (b Backend) Hosts() []string {
	return b.hosts
}

// BalancedHost returns a balanced host from the Backend instance. If there is only one host, it is returned directly.
// Otherwise, a random host is selected from the available ones.
func (b Backend) BalancedHost() string {
	// todo: aqui obtemos o host (correto é criar um domínio chamado balancer aonde ele vai retornar o host
	//  disponível pegando como base, se ele esta de pé ou não, e sua config de porcentagem)
	if helper.EqualsLen(b.hosts, 1) {
		return b.hosts[0]
	}
	return b.hosts[helper.RandomNumber(0, len(b.hosts)-1)]
}

// Path returns the path of the Backend instance.
func (b Backend) Path() string {
	return b.path
}

// Method returns the method of the Backend instance.
func (b Backend) Method() string {
	return b.method
}

// ForwardHeaders returns the slice of strings representing the forward headers of the Backend instance.
func (b Backend) ForwardHeaders() []string {
	return b.forwardHeaders
}

// ForwardQueries returns the slice of forward queries of the Backend instance.
func (b Backend) ForwardQueries() []string {
	return b.forwardQueries
}

// BackendModifiers returns the BackendModifiers instance associated with the Backend.
// BackendModifiers contains methods to access and modify the status code, header, params,
// query, and body of the Backend.
//
// Note: This method returns a copy of the BackendModifiers instance, any modifications made to it will not affect the
// original Backend instance.
func (b Backend) BackendModifiers() BackendModifiers {
	return b.modifiers
}

// CountModifiers returns the number of modifiers present in the Backend instance.
// If the modifiers field is not nil, it counts all the modifiers using the CountAll() method of BackendModifiers.
// Otherwise, it returns 0.
func (b Backend) CountModifiers() int {
	if helper.IsNotNil(b.modifiers) {
		return b.modifiers.CountAll()
	}
	return 0
}

// StatusCode returns the status code Modifier of the BackendModifiers instance.
func (b BackendModifiers) StatusCode() Modifier {
	return b.statusCode
}

// Header returns the header modifiers of the BackendModifiers instance.
func (b BackendModifiers) Header() []Modifier {
	return b.header
}

// Params returns an array of Modifier instances that represent the params of the BackendModifiers instance.
func (b BackendModifiers) Params() []Modifier {
	return b.params
}

// Query returns the list of modifiers for the query of the BackendModifiers instance.
func (b BackendModifiers) Query() []Modifier {
	return b.query
}

// Body returns the list of modifiers for the body of the BackendModifiers instance.
func (b BackendModifiers) Body() []Modifier {
	return b.body
}

// CountAll returns the total count of modifiers for a BackendModifiers instance.
// It counts the number of valid `statusCode` and the length of `header`, `params`, `query`, and `body` slices,
// and adds them up to get the total count.
func (b BackendModifiers) CountAll() (count int) {
	if helper.IsNotNil(b.statusCode) {
		count++
	}
	count += len(b.header) + len(b.params) + len(b.query) + len(b.body)
	return count
}

// ModifyHeader returns a new backendRequest with the specified header modified.
// The input header is used to replace the existing header of the backendRequest.
// The other fields of the backendRequest remain unchanged.
func (b backendRequest) ModifyHeader(header Header) backendRequest {
	return backendRequest{
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
func (b backendRequest) ModifyParams(path string, params Params) backendRequest {
	return backendRequest{
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
func (b backendRequest) ModifyQuery(query Query) backendRequest {
	return backendRequest{
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
func (b backendRequest) ModifyBody(body Body) backendRequest {
	return backendRequest{
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
func (b backendRequest) Host() string {
	return b.host
}

// Path returns the path string of the backendRequest instance.
func (b backendRequest) Path() string {
	return b.path
}

// Url constructs the URL string for the backend request instance.
// It replaces the path placeholders with their corresponding values in the params map.
// The final URL is formed by concatenating the host and the modified path.
func (b backendRequest) Url() string {
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

func (b backendRequest) Method() string {
	return b.method
}

// Header returns the Header of the backendRequest instance.
func (b backendRequest) Header() Header {
	return b.header
}

// Params returns the Params field of the backendRequest instance.
func (b backendRequest) Params() Params {
	return b.params
}

// Query returns the query of the backendRequest.
func (b backendRequest) Query() Query {
	return b.query
}

// RawQuery encodes the query parameters into a string representation.
// It returns the encoded query string that can be appended to the URL.
func (b backendRequest) RawQuery() string {
	return url.Values(b.query).Encode()
}

// Body returns the `Body` field of the `backendRequest` instance.
// The `Body` field represents the request body.
// The `Body` method allows you to access the request body for further manipulation or inspection.
func (b backendRequest) Body() Body {
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
func (b backendRequest) BodyToRead() io.ReadCloser {
	// se ele quer omitir o body da solicitação ou o mesmo tiver vazio retornamos
	if b.omitRequestBody || helper.IsNil(b.body.value) {
		return nil
	}

	// convertemos o body para bytes
	// todo: aqui vamos obter o encode desejado XML, JSON, TEXT/PLAIN como um CONTENT-TYPE config
	bytesBody, err := helper.ConvertToBytes(b.body.ToRead())
	if helper.IsNotNil(err) {
		// todo: log?
		return nil
	}

	// retornamos o valor da interface com os bytes do body
	return io.NopCloser(bytes.NewReader(bytesBody))
}

// Http returns an HTTP request based on the backendRequest instance.
// It constructs the request, sets the headers, fills in the queries,
// and returns the created HTTP request.
// If an error occurs during the construction of the request, nil and the error are returned.
func (b backendRequest) Http(ctx context.Context) (*http.Request, error) {
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
func (b backendResponse) ModifyStatusCode(statusCode int) backendResponse {
	return backendResponse{
		name:          b.name,
		omitResponse:  b.omitResponse,
		groupResponse: b.groupResponse,
		statusCode:    statusCode,
		header:        b.header,
		body:          b.body,
	}
}

// ModifyHeader returns a new instance of backendResponse with the given header modified.
// The method creates a copy of the original backendResponse and sets the header to the provided value.
// The other fields are copied from the original backendResponse.
func (b backendResponse) ModifyHeader(header Header) backendResponse {
	return backendResponse{
		name:          b.name,
		omitResponse:  b.omitResponse,
		groupResponse: b.groupResponse,
		statusCode:    b.statusCode,
		header:        header,
		body:          b.body,
	}
}

// ModifyBody returns a new instance of backendResponse with the given body modified.
// The method creates a copy of the original backendResponse and sets the body to the provided value.
// The other fields are copied from the original backendResponse.
func (b backendResponse) ModifyBody(body Body) backendResponse {
	return backendResponse{
		name:          b.name,
		omitResponse:  b.omitResponse,
		groupResponse: b.groupResponse,
		statusCode:    b.statusCode,
		header:        b.header,
		body:          body,
	}
}

// Ok returns a boolean indicating whether the statusCode of the backendResponse is less than http.StatusOK.
// It uses the helper.IsLessThan function to compare the statusCode.
func (b backendResponse) Ok() bool {
	return helper.IsLessThan(b.statusCode, http.StatusOK)
}

// Key returns the key of the backendResponse for aggregation.
// The key is composed of the string "backend" and the index, if it is greater than or equal to zero.
// If the backendResponse has a name, the key is set to the name.
func (b backendResponse) Key(index int) (key string) {
	// montamos o key do backend para agregar
	key = "backend"
	if helper.IsGreaterThanOrEqual(index, 0) {
		key = fmt.Sprintf("%s #%v", key, index)
	}
	// se o backend tiver nome, damos prioridade
	if helper.IsNotEmpty(b.name) {
		key = b.name
	}
	return key
}

// StatusCode returns the `statusCode` of the `backendResponse` instance.
func (b backendResponse) StatusCode() int {
	return b.statusCode
}

// Header returns the `header` of the `backendResponse` instance.
func (b backendResponse) Header() Header {
	return b.header
}

// Body returns the `body` of the `backendResponse` instance.
func (b backendResponse) Body() Body {
	return b.body
}
