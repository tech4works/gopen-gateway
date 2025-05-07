/*
 * Copyright 2024 Tech4Works
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package usecase

import (
	"context"
	berrors "errors"
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/factory"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"go.elastic.co/apm/v2"
	"net/url"
	"time"
)

type endpointUseCase struct {
	httpBackendFactory  factory.HTTPBackend
	httpResponseFactory factory.HTTPResponse
	httpClient          app.HTTPClient
	publishClient       app.PublisherClient
	endpointLog         app.EndpointLog
	backendLog          app.BackendLog
	publisherLog        app.PublisherLog
}

type Endpoint interface {
	Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.HTTPResponse
}

func NewEndpoint(
	backendFactory factory.HTTPBackend,
	responseFactory factory.HTTPResponse,
	httpClient app.HTTPClient,
	publishClient app.PublisherClient,
	endpointLog app.EndpointLog,
	backendLog app.BackendLog,
	publisherLog app.PublisherLog,
) Endpoint {
	return endpointUseCase{
		httpBackendFactory:  backendFactory,
		httpResponseFactory: responseFactory,
		httpClient:          httpClient,
		publishClient:       publishClient,
		endpointLog:         endpointLog,
		backendLog:          backendLog,
		publisherLog:        publisherLog,
	}
}

func (e endpointUseCase) Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.HTTPResponse {
	history := vo.NewEmptyHistory()

	for _, backend := range executeData.Endpoint.Backends() {
		httpBackendRequest := e.buildHTTPBackendRequest(ctx, executeData, &backend, history)

		var httpBackendResponse *vo.HTTPBackendResponse
		if backend.HasRequest() && backend.Request().IsConcurrent() {
			httpBackendResponse = e.makeConcurrentBackendRequest(ctx, &backend, executeData, httpBackendRequest)
		} else {
			httpBackendResponse = e.makeBackendRequest(ctx, executeData, &backend, httpBackendRequest)
		}

		history = history.Add(&backend, httpBackendRequest, httpBackendResponse)
		if e.checkAbortBackendResponse(executeData.Endpoint, httpBackendResponse) {
			return e.buildAbortedHTTPResponse(executeData, history)
		}
	}

	err := e.executePublishers(ctx, executeData)
	if checker.NonNil(err) {
		return e.httpResponseFactory.BuildErrorResponse(executeData.Endpoint, err)
	}

	return e.buildHTTPResponse(ctx, executeData, history)
}

func (e endpointUseCase) makeConcurrentBackendRequest(
	ctx context.Context,
	backend *vo.Backend,
	executeData dto.ExecuteEndpoint,
	httpBackendRequest *vo.HTTPBackendRequest,
) *vo.HTTPBackendResponse {
	concurrentCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	responseChan := make(chan *vo.HTTPBackendResponse)
	for i := 0; i < backend.Request().Concurrent(); i++ {
		go func() {
			httpBackendResponse := e.makeBackendRequest(concurrentCtx, executeData, backend, httpBackendRequest)
			responseChan <- httpBackendResponse
		}()
	}

	select {
	case httpBackendResponse := <-responseChan:
		return httpBackendResponse
	}
}

func (e endpointUseCase) makeBackendRequest(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	httpBackendRequest *vo.HTTPBackendRequest,
) *vo.HTTPBackendResponse {
	timeout, ok := ctx.Deadline()
	if !ok {
		return e.httpBackendFactory.BuildTemporaryResponseByErr(executeData.Endpoint, context.DeadlineExceeded)
	}

	requestCtx, cancel := context.WithTimeout(context.Background(), time.Until(timeout))
	defer cancel()

	e.backendLog.PrintRequest(executeData, backend, httpBackendRequest)

	startTime := time.Now()
	httpResponse, err := e.httpClient.MakeRequest(requestCtx, httpBackendRequest)
	duration := time.Since(startTime)

	var httpBackendResponse *vo.HTTPBackendResponse
	if err = e.treatHTTPClientErr(err); checker.NonNil(err) {
		httpBackendResponse = e.httpBackendFactory.BuildTemporaryResponseByErr(executeData.Endpoint, err)
	} else {
		httpBackendResponse = e.httpBackendFactory.BuildTemporaryResponse(httpResponse)
	}

	e.backendLog.PrintResponse(executeData, backend, httpBackendRequest, httpBackendResponse, duration)

	return httpBackendResponse
}

func (e endpointUseCase) publish(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	publisher *vo.Publisher,
	message *vo.Message,
) error {
	e.publisherLog.PrintRequest(executeData, publisher, message)
	return e.publishClient.Publish(ctx, publisher, message)
}

func (e endpointUseCase) treatHTTPClientErr(err error) error {
	if checker.IsNil(err) {
		return nil
	}

	var urlErr *url.Error
	berrors.As(err, &urlErr)
	if berrors.Is(urlErr.Err, context.Canceled) {
		return mapper.NewErrConcurrentCanceled()
	} else if urlErr.Timeout() {
		return mapper.NewErrGatewayTimeoutByErr(err)
	}

	return mapper.NewErrBadGateway(err)
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
		e.endpointLog.PrintWarn(executeData, err)
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

func (e endpointUseCase) executePublishers(ctx context.Context, executeData dto.ExecuteEndpoint) error {
	if !executeData.Endpoint.HasPublishers() {
		return nil
	}

	if executeData.Request.HasBody() {
		// todo: add edition flow of message sending
		body, err := executeData.Request.Body().String()
		if checker.NonNil(err) {
			return err
		}
		message := vo.NewMessage(body)

		for _, publisher := range executeData.Endpoint.Publishers() {
			err = e.publish(ctx, executeData, &publisher, &message)
			if checker.NonNil(err) {
				return err
			}
		}
	} else {
		e.endpointLog.PrintWarn(executeData, "Ignore publishers because request body is empty!")
	}

	return nil
}
