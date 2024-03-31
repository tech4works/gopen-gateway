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
	host           []string
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

func newBackend(backendDTO dto.Backend) Backend {
	return Backend{
		name:           backendDTO.Name,
		host:           backendDTO.Host,
		path:           backendDTO.Path,
		method:         backendDTO.Method,
		forwardHeaders: backendDTO.ForwardHeaders,
		forwardQueries: backendDTO.ForwardQueries,
		modifiers:      newBackendModifier(helper.IfNilReturns(backendDTO.Modifiers, dto.BackendModifiers{})),
		extraConfig:    newBackendExtraConfig(helper.IfNilReturns(backendDTO.ExtraConfig, dto.BackendExtraConfig{})),
	}
}

func newMiddlewareBackend(backendVO Backend, backendExtraConfigVO backendExtraConfig) Backend {
	return Backend{
		name:           backendVO.name,
		host:           backendVO.host,
		path:           backendVO.path,
		method:         backendVO.method,
		forwardHeaders: backendVO.forwardHeaders,
		forwardQueries: backendVO.forwardQueries,
		modifiers:      backendVO.modifiers,
		extraConfig:    backendExtraConfigVO,
	}
}

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

func newBackendExtraConfig(extraConfigDTO dto.BackendExtraConfig) backendExtraConfig {
	return backendExtraConfig{
		groupResponse:   extraConfigDTO.GroupResponse,
		omitRequestBody: extraConfigDTO.OmitRequestBody,
		omitResponse:    extraConfigDTO.OmitResponse,
	}
}

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

func (b Backend) Name() string {
	return b.name
}

func (b Backend) Host() []string {
	return b.host
}

func (b Backend) BalancedHost() string {
	// todo: aqui obtemos o host (correto é criar um domínio chamado balancer aonde ele vai retornar o host
	//  disponível pegando como base, se ele esta de pé ou não, e sua config de porcentagem)
	if helper.EqualsLen(b.host, 1) {
		return b.host[0]
	}
	return b.host[helper.RandomNumber(0, len(b.host)-1)]
}

func (b Backend) Path() string {
	return b.path
}

func (b Backend) Method() string {
	return b.method
}

func (b Backend) ForwardHeaders() []string {
	return b.forwardHeaders
}

func (b Backend) ForwardQueries() []string {
	return b.forwardQueries
}

func (b Backend) BackendModifiers() BackendModifiers {
	return b.modifiers
}

func (b Backend) CountModifiers() int {
	if helper.IsNotNil(b.modifiers) {
		return b.modifiers.CountAll()
	}
	return 0
}

func (b BackendModifiers) StatusCode() Modifier {
	return b.statusCode
}

func (b BackendModifiers) Header() []Modifier {
	return b.header
}

func (b BackendModifiers) Params() []Modifier {
	return b.params
}

func (b BackendModifiers) Query() []Modifier {
	return b.query
}

func (b BackendModifiers) Body() []Modifier {
	return b.body
}

func (b BackendModifiers) CountAll() (count int) {
	if helper.IsNotNil(b.statusCode) {
		count++
	}
	count += len(b.header) + len(b.params) + len(b.query) + len(b.body)
	return count
}

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

func (b backendRequest) Host() string {
	return b.host
}

func (b backendRequest) Path() string {
	return b.path
}

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

func (b backendRequest) Header() Header {
	return b.header
}

func (b backendRequest) Params() Params {
	return b.params
}

func (b backendRequest) Query() Query {
	return b.query
}

func (b backendRequest) RawQuery() string {
	return url.Values(b.query).Encode()
}

func (b backendRequest) Body() Body {
	return b.body
}

func (b backendRequest) BodyToSend() io.ReadCloser {
	// se ele quer omitir o body da solicitação ou o mesmo tiver vazio retornamos
	if b.omitRequestBody || helper.IsNil(b.body) {
		return nil
	}

	// convertemos o body para bytes
	// todo: aqui vamos obter o encode desejado XML, JSON, TEXT/PLAIN como um CONTENT-TYPE config
	bytesBody, err := helper.ConvertToBytes(b.body)
	if helper.IsNotNil(err) {
		// todo: log?
		return nil
	}

	// retornamos o valor da interface com os bytes do body
	return io.NopCloser(bytes.NewReader(bytesBody))
}

func (b backendRequest) Http(ctx context.Context) (*http.Request, error) {
	// construímos o http request para fazer a requisição
	httpRequest, err := http.NewRequestWithContext(ctx, b.Method(), b.Url(), b.BodyToSend())
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

func (b backendResponse) Ok() bool {
	return helper.IsLessThan(b.statusCode, http.StatusOK)
}

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

func (b backendResponse) StatusCode() int {
	return b.statusCode
}

func (b backendResponse) Header() Header {
	return b.header
}

func (b backendResponse) Body() Body {
	return b.body
}
