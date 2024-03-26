package vo

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"io"
	"net/http"
	"net/url"
)

type Backend struct {
	name           string
	host           []string
	path           string
	method         string
	forwardHeaders []string
	forwardQueries []string
	modifier       backendModifier
	extraConfig    backendExtraConfig
}

type backendModifier struct {
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
	params          map[string]string
	query           Query
	body            any
}

type backendResponse struct {
	name          string
	omitResponse  bool
	groupResponse bool
	statusCode    int
	header        Header
	body          any
}

func newMiddlewareBackend(backendVO Backend, backendExtraConfigVO backendExtraConfig) Backend {
	return Backend{
		name:           backendVO.name,
		host:           backendVO.host,
		path:           backendVO.path,
		method:         backendVO.method,
		forwardHeaders: backendVO.forwardHeaders,
		forwardQueries: backendVO.forwardQueries,
		modifier:       backendVO.modifier,
		extraConfig:    backendExtraConfigVO,
	}
}

func NewBackendRequest(backendVO Backend, balancedHost string, requestVO Request) backendRequest {
	// inicializamos o header a ser utilizado na construção do VO filtrado pelo forward-headers
	header := requestVO.Header().FilterByForwarded(backendVO.forwardHeaders)

	// inicializamos a query a ser utilizado na construção do VO filtrada pelo forward-queries
	query := requestVO.Query().FilterByForwarded(backendVO.forwardQueries)

	// inicializamos os params
	path, params := NewParamsByPath(backendVO.path, requestVO.params)

	// montamos o objeto de valor
	return backendRequest{
		omitRequestBody: backendVO.extraConfig.omitRequestBody,
		host:            balancedHost,
		path:            path,
		method:          backendVO.method,
		header:          header,
		params:          params,
		query:           query,
		body:            requestVO.Body(),
	}
}

func NewBackendResponse(backendVO Backend, httpResponse *http.Response) backendResponse {
	// fazemos o parse dos bytes da resposta em para uma interface
	var body any
	bodyBytes, _ := io.ReadAll(httpResponse.Body)

	// se não for vazio, fazemos o parse para interface
	if helper.IsNotEmpty(bodyBytes) {
		helper.SimpleConvertToDest(bodyBytes, &body)
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

func (b Backend) Host() string {
	// todo: aqui obtemos o host (correto é criar um domínio chamado balancer aonde ele vai retornar o host
	//  disponível pegando como base, se ele esta de pé ou não, e sua config de porcentagem)
	return b.host[helper.RandomNumber(0, len(b.host)-1)]
}

func (b backendRequest) Modify(path string, header Header, params Params, query Query, body any) backendRequest {
	return backendRequest{
		host:   b.host,
		path:   path,
		method: b.method,
		header: header,
		params: params,
		query:  query,
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
	return fmt.Sprint(b.host, b.path)
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

func (b backendRequest) Body() any {
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

func (b backendResponse) Modify(statusCode int, header Header, body any) backendResponse {
	return backendResponse{
		name:          b.name,
		omitResponse:  b.omitResponse,
		groupResponse: b.groupResponse,
		statusCode:    statusCode,
		header:        header,
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

func (b backendResponse) Body() any {
	return b.body
}
