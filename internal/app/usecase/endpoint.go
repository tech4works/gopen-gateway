package usecase

import (
	"context"
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/factory"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"go.elastic.co/apm/v2"
	"time"
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
		httpBackendRequest := e.buildHTTPBackendRequest(ctx, executeData, &backend, history)

		httpBackendResponse := e.makeBackendRequest(ctx, executeData, &backend, httpBackendRequest)

		history = history.Add(&backend, httpBackendRequest, httpBackendResponse)
		if e.checkAbortBackendResponse(executeData.Endpoint, httpBackendResponse) {
			return e.buildAbortedHTTPResponse(executeData, history)
		}
	}

	return e.buildHTTPResponse(ctx, executeData, history)
}

func (e endpointUseCase) makeBackendRequest(ctx context.Context, executeData dto.ExecuteEndpoint, backend *vo.Backend,
	httpBackendRequest *vo.HTTPBackendRequest) *vo.HTTPBackendResponse {

	e.backendLog.PrintRequest(executeData, backend, httpBackendRequest)

	startTime := time.Now()
	httpBackendResponse := e.httpClient.MakeRequest(ctx, executeData.Endpoint, httpBackendRequest)
	latency := time.Since(startTime)

	e.backendLog.PrintResponse(executeData, backend, httpBackendRequest, httpBackendResponse, latency)

	return httpBackendResponse
}

func (e endpointUseCase) checkAbortBackendResponse(endpoint *vo.Endpoint, response *vo.HTTPBackendResponse) bool {
	statusCode := response.StatusCode()
	return (endpoint.HasAbortStatusCodes() && checker.Contains(endpoint.AbortIfStatusCodes(), statusCode.Code())) ||
		(!endpoint.HasAbortStatusCodes() && statusCode.Failed())
}

func (e endpointUseCase) buildHTTPBackendRequest(ctx context.Context, executeData dto.ExecuteEndpoint, backend *vo.Backend,
	history *vo.History) *vo.HTTPBackendRequest {
	span, _ := apm.StartSpan(ctx, "Backend request", "factory")
	if checker.NonNil(span) {
		span.Context.SetLabel("transformations", backend.CountRequestDataTransforms())
		defer span.End()
	}

	httpBackendRequest, errs := e.httpBackendFactory.BuildRequest(backend, executeData.Request, history)
	for _, err := range errs {
		e.backendLog.PrintWarn(executeData, backend, httpBackendRequest, err)
	}
	return httpBackendRequest
}

func (e endpointUseCase) buildHTTPBackendResponse(executeData dto.ExecuteEndpoint, backend *vo.Backend,
	httpBackendRequest *vo.HTTPBackendRequest, httpBackendResponse *vo.HTTPBackendResponse, history *vo.History,
) *vo.HTTPBackendResponse {
	if !backend.HasResponse() {
		return httpBackendResponse
	}

	httpBackendResponse, errors := e.httpBackendFactory.BuildResponse(backend, httpBackendResponse, executeData.Request, history)
	for _, err := range errors {
		e.backendLog.PrintWarn(executeData, backend, httpBackendRequest, err)
	}

	return httpBackendResponse
}

func (e endpointUseCase) buildAbortedHTTPResponse(executeData dto.ExecuteEndpoint, history *vo.History) *vo.HTTPResponse {
	return e.httpResponseFactory.BuildAbortedResponse(executeData.Endpoint, history)
}

func (e endpointUseCase) buildHTTPResponse(ctx context.Context, executeData dto.ExecuteEndpoint, history *vo.History,
) *vo.HTTPResponse {
	filteredHistory := e.filterHistory(ctx, executeData, history)
	httpResponse, errs := e.httpResponseFactory.BuildResponse(executeData.Endpoint, filteredHistory)

	for _, err := range errs {
		e.printEndpointWarn(executeData, err)
	}

	return httpResponse
}

func (e endpointUseCase) filterHistory(ctx context.Context, executeData dto.ExecuteEndpoint, history *vo.History,
) *vo.History {
	span, _ := apm.StartSpan(ctx, "Response", "factory")
	if checker.NonNil(span) {
		defer span.End()
	}

	var backends []*vo.Backend
	var requests []*vo.HTTPBackendRequest
	var responses []*vo.HTTPBackendResponse

	for i := 0; i < history.Size(); i++ {
		backend, httpBackendRequest, httpBackendTemporaryResponse := history.Get(i)

		httpBackendResponse := e.buildHTTPBackendResponse(executeData, backend, httpBackendRequest,
			httpBackendTemporaryResponse, history)

		if checker.NonNil(httpBackendResponse) {
			backends = append(backends, backend)
			requests = append(requests, httpBackendRequest)
			responses = append(responses, httpBackendResponse)
		}
	}

	return vo.NewHistory(backends, requests, responses)
}

func (e endpointUseCase) printEndpointWarn(executeData dto.ExecuteEndpoint, err error) {
	e.endpointLog.PrintWarn(executeData.Endpoint, executeData.Request, executeData.ClientIP, executeData.TraceID, err)
}
