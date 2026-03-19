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

package factory

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/app/model/publisher"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type backendResponse struct {
	buildPipelineService service.BuildPipeline
}

type BackendResponse interface {
	BuildResponseByError(
		endpoint *vo.EndpointConfig,
		backend *vo.BackendConfig,
		err error,
		duration time.Duration,
	) *vo.BackendResponse
	BuildResponseByHTTP(httpResponse *http.Response, duration time.Duration) *vo.BackendResponse
	BuildResponseByPublisher(publisherResponse *publisher.Response, duration time.Duration) *vo.BackendResponse
	BuildFinalResponse(
		backend *vo.BackendConfig,
		response *vo.BackendResponse,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) (*vo.BackendResponse, []error)
}

func NewBackendResponse(buildPipelineService service.BuildPipeline) BackendResponse {
	return backendResponse{
		buildPipelineService: buildPipelineService,
	}
}

func (f backendResponse) BuildResponseByError(
	endpoint *vo.EndpointConfig,
	backend *vo.BackendConfig,
	err error,
	duration time.Duration,
) *vo.BackendResponse {
	outcome := f.buildBackendOutcomeByError(err)
	status := f.buildBackendResponseStatusByError(err)

	wrapped := errors.Wrap(err)
	payload := vo.NewPayloadJSON(converter.ToBuffer(dto.ErrorPayload{
		ID:        backend.ID(),
		File:      wrapped.File(),
		Line:      wrapped.Line(),
		Endpoint:  endpoint.Path(),
		Message:   wrapped.Message(),
		Timestamp: time.Now(),
	}))
	metadata := vo.NewEmptyMetadata()

	return vo.NewBackendResponse(backend.Kind(), outcome, duration, status, metadata, payload)
}

func (f backendResponse) BuildResponseByHTTP(httpResponse *http.Response, duration time.Duration) *vo.BackendResponse {
	status := f.buildResponseStatusFromHTTP(httpResponse.StatusCode)

	var body *vo.Payload
	if checker.NonNil(httpResponse.Body) {
		contentType := httpResponse.Header.Get(app.ContentType)
		contentEncoding := httpResponse.Header.Get(app.ContentEncoding)

		bodyBytes, err := io.ReadAll(httpResponse.Body)
		if checker.NonNil(err) {
			panic(err)
		}

		body = vo.NewPayload(contentType, contentEncoding, bytes.NewBuffer(bodyBytes))
	}

	mapHeader := map[string][]string{}
	for k, v := range httpResponse.Header {
		if checker.NotContains(app.TransportHTTPHeaderKeys(), k) {
			mapHeader[k] = v
		}
	}
	header := vo.NewMetadata(mapHeader)

	return vo.NewBackendResponse(
		enum.BackendKindHTTP,
		enum.BackendOutcomeExecuted,
		duration,
		status,
		header,
		body,
	)
}

func (f backendResponse) BuildResponseByPublisher(publisherResponse *publisher.Response, duration time.Duration,
) *vo.BackendResponse {
	return vo.NewBackendResponse(
		enum.BackendKindPublisher,
		enum.BackendOutcomeExecuted,
		duration,
		vo.NewResponseStatusByValue(enum.ResponseStatusOK),
		vo.NewEmptyMetadata(),
		vo.NewPayloadJSON(converter.ToBuffer(publisherResponse.Body)),
	)
}

func (f backendResponse) BuildFinalResponse(
	backend *vo.BackendConfig,
	response *vo.BackendResponse,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.BackendResponse, []error) {
	if !backend.HasResponse() || checker.IsNil(response) || response.ShouldIgnoreFinalResponseBuild() {
		return response, nil
	} else if backend.Response().Omit() {
		return nil, nil
	}

	useFallback := backend.Execution().UseFallback(enum.ExecutionOnBuild)

	payload, payloadDegraded, payloadErrs := f.buildFinalResponsePayload(
		backend.Response().Payload(), response, request, history, useFallback)
	metadata, metadataDegraded, metadataErrs := f.buildFinalResponseMetadata(
		backend.Response().Metadata(), response, request, history, useFallback)

	var degradationKinds []enum.DegradationKind
	if payloadDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindPayload)
	}
	if metadataDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindMetadata)
	}

	return vo.NewBackendResponseWithDegradation(
		backend.Kind(),
		response.Outcome(),
		vo.NewDegradation(degradationKinds...),
		response.Duration(),
		response.Status(),
		metadata,
		payload,
	), joinErrs(payloadErrs, metadataErrs)
}

func (f backendResponse) buildBackendOutcomeByError(err error) enum.BackendOutcome {
	if errors.Is(err, app.ErrBackendDependenciesNotExecuted) {
		return enum.BackendOutcomeCancelled
	} else if errors.Is(err, domain.ErrEvalGuards) {
		return enum.BackendOutcomeIgnored
	} else {
		return enum.BackendOutcomeError
	}
}

func (f backendResponse) buildBackendResponseStatusByError(err error) vo.ResponseStatus {
	var status enum.ResponseStatus
	if errors.Is(err, app.ErrBackendDependenciesNotExecuted) || errors.Is(err, domain.ErrEvalGuards) {
		status = enum.ResponseStatusCancelled
	} else if errors.Is(err, app.ErrBackendBrokerNotImplemented) {
		status = enum.ResponseStatusUnimplemented
	} else if errors.Is(err, app.ErrBackendBadGateway) {
		status = enum.ResponseStatusBadGateway
	} else if errors.Is(err, app.ErrBackendGatewayTimeout) || errors.Is(err, context.DeadlineExceeded) {
		status = enum.ResponseStatusDeadlineExceeded
	} else {
		status = enum.ResponseStatusInternalError
	}
	return vo.NewResponseStatus(status, err, err.Error())
}

func (f backendResponse) buildFinalResponsePayload(
	spec vo.PayloadPipelineSpec,
	response *vo.BackendResponse,
	request *vo.EndpointRequest,
	history *aggregate.History,
	useFallback bool,
) (*vo.Payload, bool, []error) {
	payload, errs := f.buildPipelineService.ApplyPayload(spec, response.Payload(), request, history)
	return fallbackIf(useFallback, errs, payload, response.Payload()), isDegraded(spec, errs), errs
}

func (f backendResponse) buildFinalResponseMetadata(
	spec vo.MetadataPipelineSpec,
	response *vo.BackendResponse,
	request *vo.EndpointRequest,
	history *aggregate.History,
	useFallback bool,
) (vo.Metadata, bool, []error) {
	metadata, errs := f.buildPipelineService.ApplyMetadata(spec, response.Metadata(), app.TransportHTTPHeaderKeys(),
		request, history)
	return fallbackIf(useFallback, errs, metadata, response.Metadata()), isDegraded(spec, errs), errs
}

func (f backendResponse) buildResponseStatusFromHTTP(code int) vo.ResponseStatus {
	var responseStatusEnum enum.ResponseStatus

	switch code {
	case 200, 201, 202, 204:
		responseStatusEnum = enum.ResponseStatusOK
	case 400:
		responseStatusEnum = enum.ResponseStatusInvalidArgument
	case 401:
		responseStatusEnum = enum.ResponseStatusUnauthenticated
	case 403:
		responseStatusEnum = enum.ResponseStatusPermissionDenied
	case 404:
		responseStatusEnum = enum.ResponseStatusNotFound
	case 408, 504:
		responseStatusEnum = enum.ResponseStatusDeadlineExceeded
	case 409:
		responseStatusEnum = enum.ResponseStatusConflict
	case 413:
		responseStatusEnum = enum.ResponseStatusPayloadTooLarge
	case 429:
		responseStatusEnum = enum.ResponseStatusResourceExhausted
	case 431:
		responseStatusEnum = enum.ResponseStatusMetadataTooLarge
	case 500:
		responseStatusEnum = enum.ResponseStatusInternalError
	case 501:
		responseStatusEnum = enum.ResponseStatusUnimplemented
	case 502:
		responseStatusEnum = enum.ResponseStatusBadGateway
	case 503:
		responseStatusEnum = enum.ResponseStatusUnavailable
	default:
		if code >= 200 && code <= 299 {
			responseStatusEnum = enum.ResponseStatusOK
		}
		if code >= 500 && code <= 599 {
			responseStatusEnum = enum.ResponseStatusInternalError
		}
		responseStatusEnum = enum.ResponseStatusUnknown
	}
	return vo.NewResponseStatus(responseStatusEnum, code, http.StatusText(code))
}
