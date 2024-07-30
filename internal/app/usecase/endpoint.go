package usecase

import (
	"context"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/factory"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type endpointUseCase struct {
	httpBackendFactory  factory.HTTPBackend
	httpResponseFactory factory.HTTPResponse
	httpClient          app.HTTPClient
}

type Endpoint interface {
	Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.HTTPResponse
}

func NewEndpoint(backendFactory factory.HTTPBackend, responseFactory factory.HTTPResponse, httpClient app.HTTPClient,
) Endpoint {
	return endpointUseCase{
		httpBackendFactory:  backendFactory,
		httpResponseFactory: responseFactory,
		httpClient:          httpClient,
	}
}

func (e endpointUseCase) Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.HTTPResponse {
	history := vo.NewEmptyHistory()

	abort := e.processMiddlewares(ctx, executeData, executeData.Endpoint.Beforewares(), history)
	if abort {
		return e.buildAbortedHTTPResponse(executeData, history)
	}

	abort = e.processBackends(ctx, executeData, history)
	if abort {
		return e.buildAbortedHTTPResponse(executeData, history)
	}

	abort = e.processMiddlewares(ctx, executeData, executeData.Endpoint.Afterwares(), history)
	if abort {
		return e.buildAbortedHTTPResponse(executeData, history)
	}

	return e.buildHTTPResponse(executeData, history)
}

func (e endpointUseCase) makeBackendRequest(ctx context.Context, endpoint *vo.Endpoint, request *vo.HTTPRequest,
	history *vo.History, backend *vo.Backend) *vo.HTTPBackendResponse {
	httpBackendRequest := e.buildHTTPBackendRequest(request, history, backend)

	return e.httpClient.MakeRequest(ctx, endpoint, httpBackendRequest)
}

func (e endpointUseCase) processMiddlewares(ctx context.Context, executeData dto.ExecuteEndpoint, beforewares []string,
	history *vo.History) bool {

	for _, middlewareKey := range beforewares {
		middleware, ok := executeData.Gopen.Middleware(middlewareKey)
		if !ok {
			// todo: imprimir um log ou estourar um erro?
			//	 return nil, errors.New(middlewareType, middlewareKey, "not configured on middlewares field!"))
			continue
		}

		httpBackendResponse := e.makeBackendRequest(ctx, executeData.Endpoint, executeData.Request, history, middleware)

		history = history.Add(middleware, httpBackendResponse)
		if e.checkAbortBackendResponse(executeData.Endpoint, httpBackendResponse) {
			return true
		}
	}

	return false
}

func (e endpointUseCase) processBackends(ctx context.Context, executeData dto.ExecuteEndpoint, history *vo.History) bool {
	for _, backendElem := range executeData.Endpoint.Backends() {
		httpBackendResponse := e.makeBackendRequest(ctx, executeData.Endpoint, executeData.Request, history, &backendElem)

		history = history.Add(&backendElem, httpBackendResponse)
		if e.checkAbortBackendResponse(executeData.Endpoint, httpBackendResponse) {
			return true
		}
	}
	return false
}

func (e endpointUseCase) checkAbortBackendResponse(endpoint *vo.Endpoint, response *vo.HTTPBackendResponse) bool {
	statusCode := response.StatusCode()
	return (endpoint.HasAbortStatusCodes() && helper.Contains(endpoint.AbortIfStatusCodes(), statusCode.Code())) ||
		(!endpoint.HasAbortStatusCodes() && statusCode.Failed())
}

func (e endpointUseCase) buildHTTPBackendRequest(request *vo.HTTPRequest, history *vo.History, backend *vo.Backend,
) *vo.HTTPBackendRequest {
	httpBackendRequest, errs := e.httpBackendFactory.BuildRequest(backend, request, history)
	for range errs {
		// todo: printar os logs vinculando ao endpoint e backend
	}
	return httpBackendRequest
}

func (e endpointUseCase) buildAbortedHTTPResponse(executeData dto.ExecuteEndpoint, history *vo.History) *vo.HTTPResponse {
	return e.httpResponseFactory.BuildAbortedResponse(executeData.Endpoint, history)
}

func (e endpointUseCase) buildHTTPResponse(executeData dto.ExecuteEndpoint, history *vo.History) *vo.HTTPResponse {
	httpResponse, errors := e.httpResponseFactory.BuildResponse(executeData.Endpoint, executeData.Request, history)
	for range errors {
		// todo: printar os logs vinculando ao endpoint e backend
	}
	return httpResponse
}
