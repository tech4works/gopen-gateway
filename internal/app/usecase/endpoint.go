package usecase

import (
	"context"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/factory"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type endpointUseCase struct {
	httpBackendFactory  factory.HTTPBackend
	httpResponseFactory factory.HTTPResponse
	httpClient          app.HTTPClient
	endpointLog         app.EndpointLog
	backendLog          app.BackendLog
}

type Endpoint interface {
	Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.HTTPResponse
}

func NewEndpoint(backendFactory factory.HTTPBackend, responseFactory factory.HTTPResponse, httpClient app.HTTPClient,
	endpointLog app.EndpointLog, backendLog app.BackendLog) Endpoint {
	return endpointUseCase{
		httpBackendFactory:  backendFactory,
		httpResponseFactory: responseFactory,
		httpClient:          httpClient,
		endpointLog:         endpointLog,
		backendLog:          backendLog,
	}
}

func (e endpointUseCase) Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.HTTPResponse {
	history := vo.NewEmptyHistory()

	for _, backend := range executeData.Endpoint.Backends() {
		httpBackendResponse := e.makeBackendRequest(ctx, executeData, &backend, history)

		history = history.Add(&backend, httpBackendResponse)
		if e.checkAbortBackendResponse(executeData.Endpoint, httpBackendResponse) {
			return e.buildAbortedHTTPResponse(executeData, history)
		}
	}

	return e.buildHTTPResponse(executeData, history)
}

func (e endpointUseCase) makeBackendRequest(ctx context.Context, executeData dto.ExecuteEndpoint, backend *vo.Backend,
	history *vo.History) *vo.HTTPBackendResponse {
	httpBackendRequest := e.buildHTTPBackendRequest(executeData, backend, history)

	return e.httpClient.MakeRequest(ctx, executeData.Endpoint, httpBackendRequest)
}

func (e endpointUseCase) checkAbortBackendResponse(endpoint *vo.Endpoint, response *vo.HTTPBackendResponse) bool {
	statusCode := response.StatusCode()
	return (endpoint.HasAbortStatusCodes() && helper.Contains(endpoint.AbortIfStatusCodes(), statusCode.Code())) ||
		(!endpoint.HasAbortStatusCodes() && statusCode.Failed())
}

func (e endpointUseCase) buildHTTPBackendRequest(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	history *vo.History) *vo.HTTPBackendRequest {
	httpBackendRequest, errs := e.httpBackendFactory.BuildRequest(backend, executeData.Request, history)
	for _, err := range errs {
		e.backendLog.PrintWarn(executeData, backend, err)
	}
	return httpBackendRequest
}

func (e endpointUseCase) buildHTTPBackendResponse(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	httpBackendResponse *vo.HTTPBackendResponse, history *vo.History) *vo.HTTPBackendResponse {
	if !backend.HasResponse() {
		return httpBackendResponse
	}

	httpBackendResponse, errors := e.httpBackendFactory.BuildResponse(backend, httpBackendResponse, executeData.Request, history)
	for _, err := range errors {
		e.backendLog.PrintWarn(executeData, backend, err)
	}

	return httpBackendResponse
}

func (e endpointUseCase) buildAbortedHTTPResponse(executeData dto.ExecuteEndpoint, history *vo.History) *vo.HTTPResponse {
	return e.httpResponseFactory.BuildAbortedResponse(executeData.Endpoint, history)
}

func (e endpointUseCase) buildHTTPResponse(executeData dto.ExecuteEndpoint, history *vo.History) *vo.HTTPResponse {
	filteredHistory := e.filterHistory(executeData, history)

	httpResponse, errs := e.httpResponseFactory.BuildResponse(executeData.Endpoint, filteredHistory)

	for _, err := range errs {
		e.endpointLog.PrintWarn(executeData.Endpoint, executeData.TraceID, executeData.ClientIP, err)
	}

	return httpResponse
}

func (e endpointUseCase) filterHistory(executeData dto.ExecuteEndpoint, history *vo.History) *vo.History {
	var backends []*vo.Backend
	var responses []*vo.HTTPBackendResponse

	for i := 0; i < history.Size(); i++ {
		backend, httpBackendTemporaryResponse := history.Get(i)

		httpBackendResponse := e.buildHTTPBackendResponse(executeData, backend, httpBackendTemporaryResponse, history)

		if helper.IsNotNil(httpBackendResponse) {
			backends = append(backends, backend)
			responses = append(responses, httpBackendResponse)
		}
	}

	return vo.NewHistory(backends, responses)
}
