package usecase

import (
	"fmt"
	"github.com/GabrielHCataldo/go-error-detail/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/dto"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/factory"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/handler"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type endpoint struct {
	config          dto.Config
	backendService  service.Backend
	backendFactory  factory.Backend
	modifierFactory factory.Modifier
}

type Endpoint interface {
	Execute(ctx *gin.Context)
}

func NewEndpoint(
	configDto dto.Config,
	backendService service.Backend,
	backendFactory factory.Backend,
	modifierFactory factory.Modifier,
) Endpoint {
	return endpoint{
		config:          configDto,
		backendService:  backendService,
		backendFactory:  backendFactory,
		modifierFactory: modifierFactory,
	}
}

func (e endpoint) Execute(ctx *gin.Context) {
	for _, item := range e.config.Endpoints {
		if (item.Endpoint == ctx.Request.URL.Path ||
			item.Endpoint == ctx.FullPath()) &&
			item.Method == ctx.Request.Method {
			responses := &[]dto.BackendResponse{}
			requests := &[]dto.BackendRequest{}
			for _, backend := range item.Backends {
				if ctx.IsAborted() {
					return
				}
				e.processBackendAuthorizations(ctx, item, backend, requests, responses)
				e.processBackend(ctx, item, backend, requests, responses)
			}
			e.replyGateway(ctx, item, responses)
			return
		}
	}
	handler.RespondCodeWithError(ctx, http.StatusInternalServerError, errors.New(
		"endpoint is not response nothing:", ctx.Request.RequestURI, "method:", ctx.Request.Method,
	))
}

func (e endpoint) processBackendAuthorizations(
	ctx *gin.Context,
	endpoint dto.Endpoint,
	backend dto.Backend,
	requests *[]dto.BackendRequest,
	responses *[]dto.BackendResponse,
) {
	for _, authKey := range backend.Authorizations {
		authBackend, ok := e.config.ExtraConfig.Authorizations[authKey]
		if !ok {
			handler.RespondCodeWithError(ctx, http.StatusInternalServerError, errors.New(
				"authorization", "'"+authKey+"'", "not configured on extra-config.authorizations",
			))
			return
		}
		//preparamos a requisição nas factories
		err := e.prepareBackendRequest(ctx, authBackend, true, requests, responses)
		if err != nil {
			handler.RespondCodeWithError(ctx, http.StatusInternalServerError, errors.NewByErr(err))
			return
		}
		//executamos a requisition backend
		resp, err := e.backendService.Execute(ctx, authBackend, requests)
		if err != nil {
			handler.RespondCodeWithError(ctx, http.StatusBadGateway, errors.NewByErr(err))
			return
		}
		//preparamos a resposta anexando na lista responses e usaremos a ferramenta modifier
		err = e.prepareBackendResponse(resp, authBackend, requests, responses)
		if err != nil {
			handler.RespondCodeWithError(ctx, http.StatusInternalServerError, errors.NewByErr(err))
			return
		}
		//qualquer erro na requisição de autorização ja paramos e retornamos o mesmo
		if resp.StatusCode != http.StatusOK {
			e.replyGateway(ctx, endpoint, responses)
		}
	}
}

func (e endpoint) processBackend(
	ctx *gin.Context,
	endpoint dto.Endpoint,
	backend dto.Backend,
	requests *[]dto.BackendRequest,
	responses *[]dto.BackendResponse,
) {
	//preparamos a requisição para o backend
	err := e.prepareBackendRequest(ctx, backend, false, requests, responses)
	if err != nil {
		handler.RespondCodeWithError(ctx, http.StatusInternalServerError, errors.NewByErr(err))
		return
	}
	//processamos o backend enviando solicitação
	resp, err := e.backendService.Execute(ctx, backend, requests)
	if err != nil {
		handler.RespondCodeWithError(ctx, http.StatusBadGateway, errors.NewByErr(err))
		return
	}
	//preparamos a resposta anexando na lista responses e usaremos a ferramenta modifier
	err = e.prepareBackendResponse(resp, backend, requests, responses)
	if err != nil {
		handler.RespondCodeWithError(ctx, http.StatusInternalServerError, errors.NewByErr(err))
		return
	}
	//qualquer erro na requisição e teríamos mais de 1 backend e o campo abortSequential estiver true,
	//abortamos as demais requisições
	if len(endpoint.Backends) > 1 && resp.StatusCode >= http.StatusBadRequest && endpoint.AbortSequential {
		currentResp := (*responses)[len(*responses)-1]
		responses = &[]dto.BackendResponse{currentResp}
		e.replyGateway(ctx, endpoint, responses)
	}
}

func (e endpoint) prepareBackendRequest(
	ctx *gin.Context,
	backend dto.Backend,
	isAuthBackend bool,
	requests *[]dto.BackendRequest,
	responses *[]dto.BackendResponse,
) error {
	e.backendFactory.CreateBackendRequest(ctx, backend, requests, isAuthBackend)
	err := e.modifierFactory.ExecuteModifier("request", backend, *requests, *responses)
	if err != nil {
		return err
	}
	return nil
}

func (e endpoint) prepareBackendResponse(
	resp *http.Response,
	backend dto.Backend,
	requests *[]dto.BackendRequest,
	responses *[]dto.BackendResponse,
) error {
	err := e.backendFactory.CreateBackendResponse(resp, backend, responses)
	if err != nil {
		return err
	}
	err = e.modifierFactory.ExecuteModifier("response", backend, *requests, *responses)
	if err != nil {
		return err
	}
	return nil
}

func (e endpoint) replyGateway(ctx *gin.Context, endpoint dto.Endpoint, responses *[]dto.BackendResponse) {
	if ctx.IsAborted() {
		return
	}
	statusCode, header, body := e.prepareGatewayResponse(endpoint, *responses)
	for k, v := range header {
		ctx.Header(k, v)
	}
	if helper.IsString(body) && helper.IsStringNotEmpty(body) {
		ctx.String(statusCode, helper.ConvertToString(body))
	} else if helper.IsJson(body) && helper.IsJsonNotEmpty(body) {
		ctx.JSON(statusCode, body)
	} else {
		ctx.Status(statusCode)
	}
	ctx.Abort()
}

func (e endpoint) prepareGatewayResponse(endpoint dto.Endpoint, responses []dto.BackendResponse) (
	statusCode int, header map[string]string, body any) {
	gatewayCompleted, responsesCleaned := e.cleanResponses(responses)
	header = e.prepareHeaderResponse(gatewayCompleted, responses)
	if len(endpoint.Backends) > 1 && len(responsesCleaned) > 1 {
		statusCode = http.StatusOK
		if endpoint.AggregateResponses {
			body = e.prepareBodyAggregateResponses(responsesCleaned)
		} else {
			body = e.prepareBodyResponses(responsesCleaned)
		}
	} else if len(responsesCleaned) > 0 {
		statusCode = responsesCleaned[0].StatusCode
		body = responsesCleaned[0].Body
	} else {
		gatewayCompleted = false
		statusCode = http.StatusInternalServerError
	}
	return statusCode, header, body
}

func (e endpoint) cleanResponses(responses []dto.BackendResponse) (gatewayCompleted bool, result []dto.BackendResponse) {
	gatewayCompleted = true
	for _, respBackend := range responses {
		if respBackend.Remove && len(responses) > 1 {
			continue
		} else if respBackend.StatusCode >= http.StatusBadRequest {
			gatewayCompleted = false
		}
		result = append(result, respBackend)
	}
	return gatewayCompleted, result
}

func (e endpoint) prepareHeaderResponse(gatewayCompleted bool, responses []dto.BackendResponse) map[string]string {
	header := map[string]string{
		"X-Gateway-Completed": strconv.FormatBool(gatewayCompleted),
	}
	for _, respBackend := range responses {
		for k, v := range respBackend.Header {
			if k != "Content-Length" && k != "Content-Type" {
				header[k] = strings.Join(v, ", ")
			}
		}
	}
	return header
}

func (e endpoint) prepareBodyResponses(responses []dto.BackendResponse) []dto.BackendResponse {
	var result []dto.BackendResponse
	for _, respBackend := range responses {
		respBackend.Header = nil
		respBackend.Remove = false
		result = append(result, respBackend)
	}
	return result
}

func (e endpoint) prepareBodyAggregateResponses(responses []dto.BackendResponse) map[string]any {
	result := map[string]any{}
	for index, respBackend := range responses {
		if respBackend.Body == nil {
			continue
		}
		genericKey := "response" + strconv.Itoa(index)
		if helper.IsSlice(respBackend.Body) {
			key := genericKey
			if helper.IsNotEmpty(respBackend.Group) {
				key = respBackend.Group
			}
			result[key] = respBackend.Body
		} else if helper.IsJson(respBackend.Body) {
			bodyResult := map[string]any{}
			for k, v := range respBackend.Body.(map[string]any) {
				if _, ok := result[k]; ok {
					bodyResult[k+strconv.Itoa(index)] = v
				} else {
					bodyResult[k] = v
				}
			}
			if helper.IsNotEmpty(respBackend.Group) {
				result[respBackend.Group] = bodyResult
			} else {
				result = bodyResult
			}
		} else if helper.IsString(respBackend.Body) {
			key := genericKey
			if helper.IsNotEmpty(respBackend.Group) {
				key = respBackend.Group
			}
			result[key] = fmt.Sprintf("%s", respBackend.Body)
		}
	}
	if helper.IsEmpty(result) {
		return nil
	}
	return result
}
