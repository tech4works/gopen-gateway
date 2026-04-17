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
	"fmt"
	"net"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/tech4works/checker"
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/factory"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
	"github.com/tech4works/gopen-gateway/internal/infra/telemetry"
)

type endpointUseCase struct {
	dynamicValueService     service.DynamicValue
	cacheService            service.Cache
	backendRequestFactory   factory.BackendRequest
	backendResponseFactory  factory.BackendResponse
	endpointResponseFactory factory.EndpointResponse
	httpClient              app.HTTPClient
	publishClient           app.PublisherClient
	endpointLog             app.EndpointLog
	backendLog              app.BackendLog
}

type backendExecResult struct {
	i        int
	backend  *vo.BackendConfig
	response *vo.BackendResponse
}

type Endpoint interface {
	Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.EndpointResponse
}

func NewEndpoint(
	dynamicValueService service.DynamicValue,
	cacheService service.Cache,
	backendRequestFactory factory.BackendRequest,
	backendResponseFactory factory.BackendResponse,
	endpointResponseFactory factory.EndpointResponse,
	httpClient app.HTTPClient,
	publishClient app.PublisherClient,
	endpointLog app.EndpointLog,
	backendLog app.BackendLog,
) Endpoint {
	return endpointUseCase{
		dynamicValueService:     dynamicValueService,
		cacheService:            cacheService,
		backendRequestFactory:   backendRequestFactory,
		backendResponseFactory:  backendResponseFactory,
		endpointResponseFactory: endpointResponseFactory,
		httpClient:              httpClient,
		publishClient:           publishClient,
		endpointLog:             endpointLog,
		backendLog:              backendLog,
	}
}

func (e endpointUseCase) Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.EndpointResponse {
	cacheResponse := e.readEndpointResponseOnCacheIfNeeded(ctx, executeData)
	if checker.NonNil(cacheResponse) {
		return cacheResponse
	}

	history, aborted := e.executeAllBackends(ctx, executeData, executeData.Endpoint.Backends())

	var response *vo.EndpointResponse
	if aborted {
		response = e.endpointResponseFactory.BuildAbortedResponse(history)
	} else {
		response = e.buildEndpointResponse(ctx, executeData, history)
	}

	e.writeEndpointResponseOnCacheIfNeeded(ctx, executeData, history, response)

	return response
}

func (e endpointUseCase) readEndpointResponseOnCacheIfNeeded(ctx context.Context, executeData dto.ExecuteEndpoint,
) *vo.EndpointResponse {
	if !executeData.Endpoint.HasCache() || !executeData.Endpoint.AllowCache() {
		return nil
	}

	var endpointCacheEntry vo.EndpointCacheEntry
	err := e.cacheService.Read(ctx, executeData.Endpoint.Cache(), executeData.Request, nil, &endpointCacheEntry)
	if checker.NonNil(err) {
		e.endpointLog.PrintWarnf(executeData, "error to read endpoint response cache: %v", err)
		return nil
	}

	return endpointCacheEntry.Response()
}

func (e endpointUseCase) executeAllBackends(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backends []vo.BackendConfig,
) (*aggregate.History, bool) {
	seqCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	asyncDoneCh := make(chan backendExecResult, len(backends))
	abortCh := make(chan backendExecResult, 1)

	committed := make([]chan struct{}, len(backends))
	for i := range committed {
		committed[i] = make(chan struct{})
	}

	history := aggregate.NewHistoryWithSize(len(backends))

	waitDependencies := func(b *vo.BackendConfig) {
		if !b.HasDependencies() {
			return
		}
		for _, dependencyIndex := range b.Dependencies().Indexes() {
			select {
			case <-seqCtx.Done():
				return
			case <-committed[dependencyIndex]:
			}
		}
	}

	commit := func(r backendExecResult) {
		history.Add(r.i, r.backend, r.response)
		select {
		case <-committed[r.i]:
		default:
			close(committed[r.i])
		}
	}

	pollAbort := func() (backendExecResult, bool) {
		select {
		case r := <-abortCh:
			return r, true
		default:
			return backendExecResult{}, false
		}
	}

	pendingAsync := 0

	for i := range backends {
		if r, ok := pollAbort(); ok {
			commit(r)
			return history, true
		}

		backend := backends[i]

		if backend.Execution().Async() {
			pendingAsync++

			i := i
			backend := backend

			go func() {
				waitDependencies(&backend)

				response := e.executeBackend(seqCtx, executeData, &backend, history)

				r := backendExecResult{i: i, backend: &backend, response: response}

				if e.shouldBackendAbort(&backend, response) {
					select {
					case abortCh <- r:
						cancel()
					default:
					}
					return
				}

				asyncDoneCh <- r
			}()
			continue
		}

		waitDependencies(&backend)

		response := e.executeBackend(seqCtx, executeData, &backend, history)

		commit(backendExecResult{i: i, backend: &backend, response: response})

		if e.shouldBackendAbort(&backend, response) {
			cancel()
			return history, true
		}
	}

	for completed := 0; checker.IsLessThan(completed, pendingAsync); {
		select {
		case r := <-abortCh:
			commit(r)
			return history, true
		case r := <-asyncDoneCh:
			commit(r)
			completed++
		}
	}

	return history, false
}

func (e endpointUseCase) executeBackend(
	parentCtx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	history *aggregate.History,
) (backendResponse *vo.BackendResponse) {
	startTime := time.Now()

	defer func() {
		if r := recover(); checker.NonNil(r) {
			err := e.panicAsError(parentCtx, backend.ID(), r)
			backendResponse = e.backendResponseFactory.BuildResponseByError(executeData.Endpoint, backend, err,
				time.Since(startTime))
		}
		e.writeBackendResponseOnCacheIfNeeded(parentCtx, executeData, backend, history, backendResponse)
	}()

	timeout, ok := parentCtx.Deadline()

	if !ok {
		return e.backendResponseFactory.BuildResponseByError(executeData.Endpoint, backend, parentCtx.Err(),
			time.Since(startTime))
	}

	ctx, cancel := context.WithTimeout(parentCtx, time.Until(timeout))
	defer cancel()

	if err := e.checkIfCanBackendBeRun(executeData, backend, history); checker.NonNil(err) {
		backendResponse = e.backendResponseFactory.BuildResponseByError(executeData.Endpoint, backend, err,
			time.Since(startTime))
		return
	}

	cacheBackendResponse := e.readBackendResponseOnCacheIfNeeded(ctx, executeData, backend, startTime, history)
	if checker.NonNil(cacheBackendResponse) {
		backendResponse = cacheBackendResponse
		return
	}

	switch backend.Kind() {
	case enum.BackendKindHTTP:
		backendResponse = e.executeHTTPBackend(ctx, executeData, backend, startTime, history)
		return
	case enum.BackendKindPublisher:
		backendResponse = e.executePublisherBackend(ctx, executeData, backend, startTime, history)
		return
	default:
		panic(fmt.Sprintf("unknown backend kind: %v", backend.Kind()))
	}
}

func (e endpointUseCase) checkIfCanBackendBeRun(
	executeData dto.ExecuteEndpoint, backend *vo.BackendConfig, history *aggregate.History,
) error {
	if !history.DependenciesSatisfied(backend) {
		return app.NewErrBackendDependenciesNotExecuted(backend.Dependencies().IDs())
	}
	return e.evalBackendGuards(backend, executeData.Request, history)
}

func (e endpointUseCase) evalBackendGuards(
	backend *vo.BackendConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
) error {
	errs := e.dynamicValueService.EvalGuardsWithErr(backend.OnlyIf(), backend.IgnoreIf(), request, history)
	if errors.Only(errs, domain.ErrEvalGuards) {
		return errs[0]
	} else if checker.IsNotEmpty(errs) {
		return errors.JoinInheritf(errs, ", ", "failed to evaluate guard for backend id=%s kind=%s",
			backend.ID(), backend.Kind())
	} else {
		return nil
	}
}

func (e endpointUseCase) readBackendResponseOnCacheIfNeeded(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	startTime time.Time,
	history *aggregate.History,
) *vo.BackendResponse {
	if !backend.HasCache() || !backend.AllowCache() {
		return nil
	}

	var backendCacheEntry vo.BackendCacheEntry
	err := e.cacheService.Read(ctx, backend.Cache(), executeData.Request, history, &backendCacheEntry)
	if checker.NonNil(err) {
		e.backendLog.PrintWarnf(executeData, backend, "error to read backend response cache: %v", err)
		return nil
	}

	return backendCacheEntry.Response(time.Since(startTime))
}

func (e endpointUseCase) executeHTTPBackend(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	startTime time.Time,
	history *aggregate.History,
) *vo.BackendResponse {
	httpBackendRequest, err := e.buildHTTPBackendRequest(ctx, executeData, backend, history)
	if checker.NonNil(err) {
		return e.backendResponseFactory.BuildResponseByError(executeData.Endpoint, backend, err, time.Since(startTime))
	} else if backend.Execution().IsConcurrent() {
		return e.makeConcurrentBackendHTTPRequest(ctx, executeData, backend, startTime, httpBackendRequest)
	} else {
		return e.makeBackendHTTPRequest(ctx, executeData, backend, startTime, httpBackendRequest)
	}
}

func (e endpointUseCase) executePublisherBackend(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	startTime time.Time,
	history *aggregate.History,
) *vo.BackendResponse {
	publisherBackendRequest, err := e.buildPublisherRequest(ctx, executeData, backend, history)
	if checker.NonNil(err) {
		return e.backendResponseFactory.BuildResponseByError(executeData.Endpoint, backend, err, time.Since(startTime))
	} else {
		return e.makeBackendPublisherRequest(ctx, executeData, backend, startTime, publisherBackendRequest)
	}
}

func (e endpointUseCase) makeConcurrentBackendHTTPRequest(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	startTime time.Time,
	request *vo.HTTPBackendRequest,
) *vo.BackendResponse {
	concurrentCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	responseChan := make(chan *vo.BackendResponse)
	for i := 0; checker.IsLessThan(i, backend.Execution().Concurrent()); i++ {
		go func() {
			backendResponse := e.makeBackendHTTPRequest(concurrentCtx, executeData, backend, startTime, request)
			responseChan <- backendResponse
		}()
	}

	return <-responseChan
}

func (e endpointUseCase) makeBackendHTTPRequest(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	startTime time.Time,
	request *vo.HTTPBackendRequest,
) *vo.BackendResponse {
	e.backendLog.PrintHTTPRequest(executeData, backend, request)

	httpResponse, err := e.httpClient.MakeRequest(ctx, executeData.Endpoint, executeData.Request, request)

	var backendResponse *vo.BackendResponse
	if err = e.treatHTTPClientErr(err); checker.NonNil(err) {
		backendResponse = e.backendResponseFactory.BuildResponseByError(executeData.Endpoint, backend, err, time.Since(startTime))
	} else {
		backendResponse = e.backendResponseFactory.BuildResponseByHTTP(httpResponse, time.Since(startTime))
	}

	e.backendLog.PrintResponse(executeData, backend, backendResponse)

	return backendResponse
}

func (e endpointUseCase) makeBackendPublisherRequest(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	startTime time.Time,
	publisherBackendRequest *vo.PublisherBackendRequest,
) *vo.BackendResponse {
	e.backendLog.PrintPublisherRequest(executeData, backend, publisherBackendRequest)

	publisherResponse, err := e.publishClient.Publish(ctx, executeData.Request, publisherBackendRequest)

	var backendResponse *vo.BackendResponse
	if err = e.treatPublisherClientErr(err); checker.NonNil(err) {
		backendResponse = e.backendResponseFactory.BuildResponseByError(executeData.Endpoint, backend, err, time.Since(startTime))
	} else {
		backendResponse = e.backendResponseFactory.BuildResponseByPublisher(publisherResponse, time.Since(startTime))
	}

	e.backendLog.PrintResponse(executeData, backend, backendResponse)

	return backendResponse
}

func (e endpointUseCase) treatHTTPClientErr(err error) error {
	if checker.IsNil(err) {
		return nil
	}

	var urlErr *url.Error
	errors.As(err, &urlErr)

	if checker.NonNil(urlErr) && errors.Is(urlErr.Err, context.Canceled) {
		return app.NewErrBackendConcurrentCancelled()
	} else if checker.NonNil(urlErr) && urlErr.Timeout() {
		return app.NewErrBackendGatewayTimeout(err)
	} else {
		return app.NewErrBackendBadGateway(err)
	}
}

func (e endpointUseCase) treatPublisherClientErr(err error) error {
	if checker.IsNil(err) {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		return app.NewErrBackendConcurrentCancelled()
	} else if errors.Is(err, context.DeadlineExceeded) {
		return app.NewErrBackendGatewayTimeout(err)
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return app.NewErrBackendGatewayTimeout(err)
	}

	return app.NewErrBackendBadGateway(err)
}

func (e endpointUseCase) shouldBackendAbort(backend *vo.BackendConfig, response *vo.BackendResponse) bool {
	if backend.Execution().IsBestEffort() {
		return false
	}
	return backend.Execution().ShouldAbortOnResponseStatus(response.Status())
}

func (e endpointUseCase) buildHTTPBackendRequest(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	history *aggregate.History,
) (*vo.HTTPBackendRequest, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "factory.http.request",
		trace.WithAttributes(
			attribute.Int("http.transformations", backend.HTTP().CountAllDataTransforms()),
		),
	)
	defer span.End()

	httpBackendRequest, errs := e.backendRequestFactory.BuildHTTPRequest(backend, executeData.Request, history)
	if checker.IsEmpty(errs) {
		return httpBackendRequest, nil
	} else if backend.Execution().ContinueOn(enum.ExecutionOnBuild) {
		for _, err := range errs {
			e.backendLog.PrintWarnf(executeData, backend, "error build HTTP backend request: %v", err)
		}
		return httpBackendRequest, nil
	}

	err := errors.JoinInheritf(
		errs,
		", ",
		"failed to build http backend request (id=%s method=%s path=%s)",
		backend.ID(),
		backend.HTTP().Method(),
		backend.HTTP().Path(),
	)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	return nil, err
}

func (e endpointUseCase) buildPublisherRequest(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	history *aggregate.History,
) (*vo.PublisherBackendRequest, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "factory.publisher.request",
		trace.WithAttributes(
			attribute.Int("publisher.transformations", backend.Publisher().CountAllDataTransforms()),
		),
	)
	defer span.End()

	publisherRequest, errs := e.backendRequestFactory.BuildPublisherRequest(backend, executeData.Request, history)
	if checker.IsEmpty(errs) {
		return publisherRequest, nil
	} else if backend.Execution().ContinueOn(enum.ExecutionOnBuild) {
		for _, err := range errs {
			e.backendLog.PrintWarnf(executeData, backend, "error build PUBLISHER backend request: %s", err)
		}
		return publisherRequest, nil
	}

	err := errors.JoinInheritf(
		errs, ", ",
		"failed to build PUBLISHER backend request (id=%s broker=%s path=%s)",
		backend.ID(),
		backend.Publisher().Broker(),
		backend.Publisher().Path(),
	)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	return nil, err
}

func (e endpointUseCase) buildEndpointResponse(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	history *aggregate.History,
) *vo.EndpointResponse {
	ctx, span := telemetry.Tracer().Start(ctx, "factory.endpoint.response")
	defer span.End()

	err := e.buildFinalBackendResponses(executeData, history)
	if checker.NonNil(err) {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return e.endpointResponseFactory.BuildErrorResponse(executeData.Endpoint, err)
	}

	endpointResponse, errs := e.endpointResponseFactory.BuildResponse(executeData.Endpoint, executeData.Request, history)
	if checker.IsEmpty(errs) {
		return endpointResponse
	} else if executeData.Endpoint.Execution().ContinueOn(enum.ExecutionOnBuild) {
		for _, err = range errs {
			e.endpointLog.PrintWarnf(executeData, "error build endpoint response: %s", err)
		}
		return endpointResponse
	}

	buildErr := errors.JoinInheritf(
		errs, ", ",
		"failed to build endpoint response (method=%s path=%s)",
		executeData.Endpoint.Method(),
		executeData.Endpoint.Path(),
	)
	span.RecordError(buildErr)
	span.SetStatus(codes.Error, buildErr.Error())
	return e.endpointResponseFactory.BuildErrorResponse(executeData.Endpoint, buildErr)
}

func (e endpointUseCase) buildFinalBackendResponses(executeData dto.ExecuteEndpoint, history *aggregate.History) error {
	var allErrs []error
	for i := 0; checker.IsLessThan(i, history.Size()); i++ {
		backend, response := history.Get(i)

		var (
			finalResponse *vo.BackendResponse
			err           error
		)
		finalResponse, err = e.buildFinalBackendResponse(executeData, backend, response, history)
		if checker.NonNil(err) {
			allErrs = append(allErrs, err)
			continue
		}

		history.Add(i, backend, finalResponse)
	}
	if checker.IsNotEmpty(allErrs) {
		return errors.JoinInheritf(
			allErrs, ", ",
			"failed to build final backend responses (method=%s path=%s)",
			executeData.Endpoint.Method(),
			executeData.Endpoint.Path(),
		)
	}
	return nil
}

func (e endpointUseCase) buildFinalBackendResponse(
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	response *vo.BackendResponse,
	history *aggregate.History,
) (*vo.BackendResponse, error) {
	finalBackendResponse, errs := e.backendResponseFactory.BuildFinalResponse(backend, response, executeData.Request, history)

	if checker.IsEmpty(errs) {
		return finalBackendResponse, nil
	} else if backend.Execution().ContinueOn(enum.ExecutionOnBuild) {
		for _, err := range errs {
			e.backendLog.PrintWarnf(executeData, backend, "error build final backend response: %s", err)
		}
		return finalBackendResponse, nil
	}

	return nil, errors.JoinInheritf(
		errs, ", ",
		"failed to build final backend response (endpoint=%s id=%s)",
		executeData.Endpoint.Path(),
		backend.ID(),
	)
}

func (e endpointUseCase) panicAsError(ctx context.Context, where string, r any) error {
	err := errors.Newf("panic in backend goroutine (%s): %v", where, r)

	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return err
}

func (e endpointUseCase) writeBackendResponseOnCacheIfNeeded(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.BackendConfig,
	history *aggregate.History,
	backendResponse *vo.BackendResponse,
) {
	if !backend.HasCache() || !backend.AllowCache() || checker.IsNil(backendResponse) || backendResponse.ComesFromCache() {
		return
	}

	entry := vo.NewBackendCacheEntry(backend.Cache(), backendResponse)

	err := e.cacheService.Write(ctx, backend.Cache(), entry, executeData.Request, history)
	if checker.NonNil(err) {
		e.backendLog.PrintWarnf(executeData, backend, "error to write backend response cache: %v", err)
	}
}

func (e endpointUseCase) writeEndpointResponseOnCacheIfNeeded(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	history *aggregate.History,
	response *vo.EndpointResponse,
) {
	if !executeData.Endpoint.HasCache() || !executeData.Endpoint.AllowCache() {
		return
	}

	entry := vo.NewEndpointCacheEntry(executeData.Endpoint.Cache(), response)

	err := e.cacheService.Write(ctx, executeData.Endpoint.Cache(), entry, executeData.Request, history)
	if checker.NonNil(err) {
		e.endpointLog.PrintWarnf(executeData, "error to write endpoint response cache: %v", err)
	}
}
