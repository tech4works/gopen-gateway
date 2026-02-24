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
	"net/url"
	"runtime/debug"
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/factory"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
	"go.elastic.co/apm/v2"
)

type endpointUseCase struct {
	dynamicValueService service.DynamicValue
	backendFactory      factory.Backend
	httpResponseFactory factory.HTTPResponse
	httpClient          app.HTTPClient
	publishClient       app.PublisherClient
	endpointLog         app.EndpointLog
	backendLog          app.BackendLog
}

type backendExecResult struct {
	i       int
	backend vo.Backend

	httpReq  *vo.HTTPBackendRequest
	httpResp *vo.HTTPBackendResponse

	pubReq  *vo.PublisherBackendRequest
	pubResp *vo.PublisherBackendResponse

	err error
}

type Endpoint interface {
	Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.HTTPResponse
}

func NewEndpoint(
	dynamicValueService service.DynamicValue,
	backendFactory factory.Backend,
	responseFactory factory.HTTPResponse,
	httpClient app.HTTPClient,
	publishClient app.PublisherClient,
	endpointLog app.EndpointLog,
	backendLog app.BackendLog,
) Endpoint {
	return endpointUseCase{
		dynamicValueService: dynamicValueService,
		backendFactory:      backendFactory,
		httpResponseFactory: responseFactory,
		httpClient:          httpClient,
		publishClient:       publishClient,
		endpointLog:         endpointLog,
		backendLog:          backendLog,
	}
}

func (e endpointUseCase) Execute(ctx context.Context, executeData dto.ExecuteEndpoint) *vo.HTTPResponse {
	backends := executeData.Endpoint.Backends()

	history, aborted, err := e.executeAllBackends(ctx, executeData, backends)
	if checker.NonNil(err) {
		return e.httpResponseFactory.BuildErrorResponse(executeData.Endpoint, err)
	} else if aborted {
		return e.httpResponseFactory.BuildAbortedResponse(history)
	}

	return e.buildHTTPResponse(ctx, executeData, history)
}

func (e endpointUseCase) executeAllBackends(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backends []vo.Backend,
) (*aggregate.History, bool, error) {
	seqCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	asyncDoneCh := make(chan backendExecResult, len(backends))
	abortCh := make(chan backendExecResult, 1)

	committed := make([]chan struct{}, len(backends))
	for i := range committed {
		committed[i] = make(chan struct{})
	}

	history := aggregate.NewHistoryWithSize(len(backends))

	waitDependencies := func(b vo.Backend) {
		for _, dependencyIndex := range b.Dependencies().Indexes() {
			select {
			case <-seqCtx.Done():
				return
			case <-committed[dependencyIndex]:
			}
		}
	}

	commitResult := func(r backendExecResult) {
		history.AddBackend(r.i, &r.backend, r.httpReq, r.httpResp, r.pubReq, r.pubResp)
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
			commitResult(r)
			return history, true, nil
		}

		backend := backends[i]

		if backend.Async() {
			pendingAsync++

			i := i
			backend := backend

			safeSendBackendResult(seqCtx, "executeBackend.runAsync", asyncDoneCh, func() backendExecResult {
				waitDependencies(backend)

				httpReq, httpResp, pubReq, pubResp, err := e.executeBackend(seqCtx, executeData, &backend, history)

				if e.shouldBackendAbort(executeData.Endpoint, httpResp, pubResp, err) {
					select {
					case abortCh <- backendExecResult{
						i:        i,
						backend:  backend,
						httpReq:  httpReq,
						httpResp: httpResp,
						pubReq:   pubReq,
						pubResp:  pubResp,
						err:      nil,
					}:
						cancel()
					default:
					}
				}

				return backendExecResult{
					i:        i,
					backend:  backend,
					httpReq:  httpReq,
					httpResp: httpResp,
					pubReq:   pubReq,
					pubResp:  pubResp,
					err:      err,
				}
			})
			continue
		}

		waitDependencies(backend)

		httpReq, httpResp, pubReq, pubResp, err := e.executeBackend(seqCtx, executeData, &backend, history)
		if checker.NonNil(err) {
			cancel()
			return history, false, err
		}

		commitResult(backendExecResult{
			i:        i,
			backend:  backend,
			httpReq:  httpReq,
			httpResp: httpResp,
			pubReq:   pubReq,
			pubResp:  pubResp,
			err:      nil,
		})

		if e.shouldBackendAbort(executeData.Endpoint, httpResp, pubResp, err) {
			cancel()
			return history, true, nil
		}
	}

	for completed := 0; checker.IsLessThan(completed, pendingAsync); {
		select {
		case r := <-abortCh:
			commitResult(r)
			return history, true, nil
		case r := <-asyncDoneCh:
			if checker.NonNil(r.err) {
				cancel()
				return history, false, r.err
			}
			commitResult(r)
			completed++
		}
	}

	return history, false, nil
}

func (e endpointUseCase) executeBackend(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	history *aggregate.History,
) (*vo.HTTPBackendRequest, *vo.HTTPBackendResponse, *vo.PublisherBackendRequest, *vo.PublisherBackendResponse, error) {
	switch backend.Kind() {
	case enum.BackendKindHTTP:
		httpReq, httpResp, err := e.executeHTTPBackend(ctx, executeData, backend, history)
		return httpReq, httpResp, nil, nil, err
	case enum.BackendKindPublisher:
		pubReq, pubResp, err := e.executePublisherBackend(ctx, executeData, backend, history)
		return nil, nil, pubReq, pubResp, err
	default:
		return nil, nil, nil, nil, errors.Newf("invalid backend kind: %v", backend.Kind())
	}
}

func (e endpointUseCase) evalBackendGuards(
	backend *vo.Backend,
	request *vo.HTTPRequest,
	history *aggregate.History,
) error {
	errs := e.dynamicValueService.EvalGuardsWithErr(backend.OnlyIf(), backend.IgnoreIf(), request, history)
	if errors.Only(errs, mapper.ErrEvalGuards) {
		return errs[0]
	} else if checker.IsNotEmpty(errs) {
		return errors.JoinInheritf(errs, ", ", "failed to evaluate guard for backend kind=%v", backend.Kind())
	} else {
		return nil
	}
}

func (e endpointUseCase) executeHTTPBackend(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	history *aggregate.History,
) (*vo.HTTPBackendRequest, *vo.HTTPBackendResponse, error) {
	if err := e.evalBackendGuards(backend, executeData.Request, history); checker.NonNil(err) {
		return nil, e.backendFactory.BuildHTTPResponseByErr(executeData.Endpoint, backend, err), nil
	}

	httpBackendRequest, err := e.buildHTTPBackendRequest(ctx, executeData, backend, history)
	if checker.NonNil(err) {
		return httpBackendRequest, nil, err
	} else if backend.HTTP().HasRequest() && backend.HTTP().Request().IsConcurrent() {
		return httpBackendRequest, e.makeConcurrentBackendHTTPRequest(ctx, backend, executeData, httpBackendRequest), nil
	} else {
		return httpBackendRequest, e.makeBackendHTTPRequest(ctx, executeData, backend, httpBackendRequest), nil
	}
}

func (e endpointUseCase) executePublisherBackend(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	history *aggregate.History,
) (*vo.PublisherBackendRequest, *vo.PublisherBackendResponse, error) {
	if err := e.evalBackendGuards(backend, executeData.Request, history); checker.NonNil(err) {
		return nil, e.backendFactory.BuildPublisherResponseByErr(executeData.Endpoint, backend, err), nil
	}

	publisherBackendRequest, err := e.buildPublisherRequest(ctx, executeData, backend, history)
	if checker.IsNotEmpty(err) {
		return nil, nil, err
	}

	return publisherBackendRequest, e.makeBackendPublisherRequest(ctx, executeData, backend, publisherBackendRequest), nil
}

func (e endpointUseCase) makeConcurrentBackendHTTPRequest(
	ctx context.Context,
	backend *vo.Backend,
	executeData dto.ExecuteEndpoint,
	request *vo.HTTPBackendRequest,
) *vo.HTTPBackendResponse {
	concurrentCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	responseChan := make(chan *vo.HTTPBackendResponse)
	for i := 0; checker.IsLessThan(i, backend.HTTP().Request().Concurrent()); i++ {
		go func() {
			httpBackendResponse := e.makeBackendHTTPRequest(concurrentCtx, executeData, backend, request)
			responseChan <- httpBackendResponse
		}()
	}

	select {
	case httpBackendResponse := <-responseChan:
		return httpBackendResponse
	}
}

func (e endpointUseCase) makeBackendHTTPRequest(
	parent context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	request *vo.HTTPBackendRequest,
) *vo.HTTPBackendResponse {
	timeout, ok := parent.Deadline()
	if !ok {
		return e.backendFactory.BuildHTTPResponseByErr(executeData.Endpoint, backend, context.DeadlineExceeded)
	}

	ctx, cancel := context.WithTimeout(parent, time.Until(timeout))
	defer cancel()

	e.backendLog.PrintHTTPRequest(executeData, backend, request)

	startTime := time.Now()
	httpResponse, err := e.httpClient.MakeRequest(ctx, request)
	duration := time.Since(startTime)

	var httpBackendResponse *vo.HTTPBackendResponse
	if err = e.treatHTTPClientErr(err); errors.Is(err, mapper.ErrBackendConcurrentCancelled) {
		httpBackendResponse = e.backendFactory.BuildHTTPResponseByErr(executeData.Endpoint, backend, err)
	} else {
		httpBackendResponse = e.backendFactory.BuildHTTPResponse(httpResponse)
	}

	e.backendLog.PrintHTTPResponse(executeData, backend, httpBackendResponse, duration)

	return httpBackendResponse
}

func (e endpointUseCase) treatHTTPClientErr(err error) error {
	if checker.IsNil(err) {
		return nil
	}

	var urlErr *url.Error
	errors.As(err, &urlErr)
	if errors.Is(urlErr.Err, context.Canceled) {
		return mapper.NewErrBackendConcurrentCancelled()
	} else if urlErr.Timeout() {
		return mapper.NewErrBackendGatewayTimeout(err)
	} else {
		return mapper.NewErrBackendBadGateway(err)
	}
}

func (e endpointUseCase) makeBackendPublisherRequest(
	parent context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	publisherBackendRequest *vo.PublisherBackendRequest,
) *vo.PublisherBackendResponse {
	timeout, ok := parent.Deadline()
	if !ok {
		return e.backendFactory.BuildPublisherResponseByErr(executeData.Endpoint, backend, context.DeadlineExceeded)
	}

	ctx, cancel := context.WithTimeout(parent, time.Until(timeout))
	defer cancel()

	e.backendLog.PrintPublisherRequest(executeData, backend, publisherBackendRequest)

	startTime := time.Now()
	publisherResponse, err := e.publishClient.Publish(ctx, publisherBackendRequest)
	duration := time.Since(startTime)

	var publisherBackendResponse *vo.PublisherBackendResponse
	if err = e.treatPublisherClientErr(err); checker.NonNil(err) {
		publisherBackendResponse = e.backendFactory.BuildPublisherResponseByErr(executeData.Endpoint, backend, err)
	} else {
		publisherBackendResponse = e.backendFactory.BuildPublisherResponse(publisherResponse)
	}

	e.backendLog.PrintPublisherResponse(executeData, backend, publisherBackendResponse, duration)

	return publisherBackendResponse
}

func (e endpointUseCase) treatPublisherClientErr(err error) error {
	if errors.Is(err, context.Canceled) {
		return mapper.NewErrBackendConcurrentCancelled()
	} else {
		return err
	}
}

func (e endpointUseCase) shouldBackendAbort(
	endpoint *vo.Endpoint,
	httpResp *vo.HTTPBackendResponse,
	pubResp *vo.PublisherBackendResponse,
	err error,
) bool {
	if checker.NonNil(err) {
		return true
	} else if checker.IsNil(httpResp) && checker.IsNil(pubResp) {
		return false
	}

	var status vo.StatusCode
	if checker.NonNil(httpResp) {
		status = httpResp.StatusCode()
	} else {
		status = httpResp.StatusCode()
	}

	if endpoint.HasAbortIfStatusCodes() {
		return checker.Contains(endpoint.AbortIfStatusCodes(), status.Code())
	}

	return status.Failed()
}

func (e endpointUseCase) buildHTTPBackendRequest(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	history *aggregate.History,
) (*vo.HTTPBackendRequest, error) {
	span, ctx := apm.StartSpan(ctx, "http.request", "factory")
	defer span.End()

	span.Context.SetLabel("transformations", backend.HTTP().CountAllDataTransforms())

	httpBackendRequest, errs := e.backendFactory.BuildHTTPRequest(backend.HTTP(), executeData.Request, history)
	if checker.IsEmpty(errs) {
		return httpBackendRequest, nil
	} else if backend.HTTP().HasRequest() && backend.HTTP().Request().ContinueOnError() {
		for _, err := range errs {
			e.backendLog.PrintWarn(executeData, backend, err)
		}
		return httpBackendRequest, nil
	}

	return httpBackendRequest, errors.JoinInheritf(
		errs,
		", ",
		"failed to build backend request (endpoint=%s method=%s path=%s)",
		executeData.Endpoint.Path(),
		backend.HTTP().Method(),
		backend.HTTP().Path(),
	)
}

func (e endpointUseCase) buildFinalHTTPBackendResponse(
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	httpBackendResponse *vo.HTTPBackendResponse,
	history *aggregate.History,
) (*vo.HTTPBackendResponse, error) {
	if !backend.HasResponse() || !httpBackendResponse.Executed() {
		return httpBackendResponse, nil
	} else if backend.Response().Omit() {
		return nil, nil
	}

	finalHTTPBackendResponse, errs := e.backendFactory.BuildFinalHTTPResponse(
		backend,
		httpBackendResponse,
		executeData.Request,
		history,
	)
	if checker.IsEmpty(errs) {
		return finalHTTPBackendResponse, nil
	} else if backend.Response().ContinueOnError() {
		for _, err := range errs {
			e.backendLog.PrintWarn(executeData, backend, err)
		}
		return finalHTTPBackendResponse, nil
	}

	return finalHTTPBackendResponse, errors.JoinInheritf(
		errs, ", ",
		"failed to build final backend response (endpoint=%s method=%s path=%s)",
		executeData.Endpoint.Path(),
		backend.HTTP().Method(),
		backend.HTTP().Path(),
	)
}

func (e endpointUseCase) buildFinalPublisherBackendResponse(
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	publisherResponse *vo.PublisherBackendResponse,
	history *aggregate.History,
) (*vo.PublisherBackendResponse, error) {
	if !backend.HasResponse() || !publisherResponse.Executed() {
		return publisherResponse, nil
	}

	finalPublisherResponse, errs := e.backendFactory.BuildFinalPublisherResponse(
		backend,
		publisherResponse,
		executeData.Request,
		history,
	)
	if checker.IsEmpty(errs) {
		return finalPublisherResponse, nil
	} else if backend.HasResponse() && backend.Response().ContinueOnError() {
		for _, err := range errs {
			e.backendLog.PrintWarn(executeData, backend, err)
		}
		return finalPublisherResponse, nil
	}

	return finalPublisherResponse, errors.JoinInheritf(
		errs, ", ",
		"failed to build final backend publisher response (endpoint=%s broker=%s path=%s)",
		executeData.Endpoint.Path(),
		backend.Publisher().Broker(),
		backend.Publisher().Path(),
	)
}

func (e endpointUseCase) buildHTTPResponse(ctx context.Context, executeData dto.ExecuteEndpoint, history *aggregate.History,
) *vo.HTTPResponse {
	span, ctx := apm.StartSpan(ctx, "endpoint.response", "factory")
	defer span.End()

	filteredHistory, err := e.filterHistory(executeData, history)
	if checker.NonNil(err) {
		return e.httpResponseFactory.BuildErrorResponse(executeData.Endpoint, err)
	}

	httpResponse, errs := e.httpResponseFactory.BuildResponse(executeData.Endpoint, executeData.Request, filteredHistory)
	if checker.IsEmpty(errs) {
		return httpResponse
	} else if executeData.Endpoint.HasResponse() && executeData.Endpoint.Response().ContinueOnError() {
		for _, err := range errs {
			e.endpointLog.PrintWarn(executeData, err)
		}
		return httpResponse
	}

	return e.httpResponseFactory.BuildErrorResponse(executeData.Endpoint, errors.JoinInheritf(
		errs, ", ",
		"failed to build endpoint response (method=%s path=%s)",
		executeData.Endpoint.Method(),
		executeData.Endpoint.Path(),
	))
}

func (e endpointUseCase) filterHistory(executeData dto.ExecuteEndpoint, history *aggregate.History) (*aggregate.History, error) {
	for i := 0; checker.IsLessThan(i, history.Size()); i++ {
		backend, httpReq, tempHTTPRes, pubReq, tmpPubRes := history.GetBackend(i)
		if checker.IsNil(backend) {
			continue
		}

		var httpFinal *vo.HTTPBackendResponse
		var pubFinal *vo.PublisherBackendResponse
		var err error
		if backend.IsHTTP() && checker.NonNil(tempHTTPRes) {
			httpFinal, err = e.buildFinalHTTPBackendResponse(executeData, backend, tempHTTPRes, history)
		} else if backend.IsPublisher() && checker.NonNil(tmpPubRes) {
			pubFinal, err = e.buildFinalPublisherBackendResponse(executeData, backend, tmpPubRes, history)
		}
		if checker.NonNil(err) {
			return nil, err
		}

		history.AddBackend(i, backend, httpReq, httpFinal, pubReq, pubFinal)
	}

	return history, nil
}

func panicAsError(ctx context.Context, where string, r any) error {
	err := fmt.Errorf("panic in backend goroutine (%s): %v\n%s", where, r, string(debug.Stack()))

	apm.CaptureError(ctx, err).Send()

	return err
}

func safeSendBackendResult(
	ctx context.Context,
	where string,
	out chan<- backendExecResult,
	build func() backendExecResult,
) {
	go func() {
		var r backendExecResult

		defer func() {
			if rec := recover(); checker.NonNil(rec) {
				r.err = panicAsError(ctx, where, rec)
			}
			select {
			case <-ctx.Done():
				return
			case out <- r:
				return
			}
		}()

		r = build()
	}()
}

func (e endpointUseCase) buildPublisherRequest(
	ctx context.Context,
	executeData dto.ExecuteEndpoint,
	backend *vo.Backend,
	history *aggregate.History,
) (*vo.PublisherBackendRequest, error) {
	span, ctx := apm.StartSpan(ctx, "publisher.request", "factory")
	defer span.End()

	span.Context.SetLabel("publisher.transformations", backend.Publisher().CountAllDataTransforms())

	publisherRequest, errs := e.backendFactory.BuildPublisherRequest(executeData.Request, history, backend.Publisher())
	if checker.IsEmpty(errs) {
		return publisherRequest, nil
	} else if backend.Publisher().HasMessage() && backend.Publisher().Message().ContinueOnError() {
		for _, err := range errs {
			e.backendLog.PrintWarn(executeData, backend, err)
		}
		return publisherRequest, nil
	}

	return publisherRequest, errors.JoinInheritf(
		errs, ", ",
		"failed to build publisher request (endpoint=%s broker=%s path=%s)",
		executeData.Endpoint.Path(),
		backend.Publisher().Broker(),
		backend.Publisher().Path(),
	)
}
