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
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/app/model/publisher"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type backendFactory struct {
	buildPipelineService service.BuildPipeline
}

type Backend interface {
	BuildHTTPRequest(http *vo.HTTP, request *vo.HTTPRequest, history *aggregate.History) (*vo.HTTPBackendRequest, []error)
	BuildPublisherRequest(request *vo.HTTPRequest, history *aggregate.History, publisher *vo.Publisher) (*vo.PublisherBackendRequest, []error)
	BuildHTTPResponse(httpResponse *http.Response) *vo.HTTPBackendResponse
	BuildHTTPResponseByErr(endpoint *vo.Endpoint, backend *vo.Backend, err error) *vo.HTTPBackendResponse
	BuildPublisherResponse(publisherResponse *publisher.Response) *vo.PublisherBackendResponse
	BuildPublisherResponseByErr(endpoint *vo.Endpoint, backend *vo.Backend, err error) *vo.PublisherBackendResponse
	BuildFinalHTTPResponse(backend *vo.Backend, response *vo.HTTPBackendResponse, request *vo.HTTPRequest, history *aggregate.History) (*vo.HTTPBackendResponse, []error)
	BuildFinalPublisherResponse(backend *vo.Backend, response *vo.PublisherBackendResponse, request *vo.HTTPRequest, history *aggregate.History) (*vo.PublisherBackendResponse, []error)
}

func NewBackend(buildPipelineService service.BuildPipeline) Backend {
	return backendFactory{
		buildPipelineService: buildPipelineService,
	}
}

func (f backendFactory) BuildHTTPRequest(http *vo.HTTP, request *vo.HTTPRequest, history *aggregate.History) (
	*vo.HTTPBackendRequest, []error) {
	hostPipelineSpec := vo.NewHostPipelineSpec(http.Hosts())
	bodyPipelineSpec := vo.NewBodyPipelineSpecFromBackendRequest(http.Request())
	headerPipelineSpec := vo.NewHeaderPipelineSpecFromBackendRequest(http.Request())
	urlPathPipelineSpec := vo.NewURLPathPipelineSpecFromBackendRequest(http.Request())
	queryPipelineSpec := vo.NewQueryPipelineSpecFromBackendRequest(http.Request())

	host, hostErrs := f.buildPipelineService.ApplyHost(hostPipelineSpec)
	body, bodyErrs := f.buildPipelineService.ApplyBody(bodyPipelineSpec, request.Body(), request, history)
	urlPath, urlPathErrs := f.buildPipelineService.ApplyURLPath(urlPathPipelineSpec, request.Path(), request, history)
	header, headerErrs := f.buildPipelineService.ApplyHeader(headerPipelineSpec, request.Header(), request, history)
	query, queryErrs := f.buildPipelineService.ApplyQuery(queryPipelineSpec, request.Query(), request, history)

	var allErrs []error
	allErrs = append(allErrs, hostErrs...)
	allErrs = append(allErrs, bodyErrs...)
	allErrs = append(allErrs, urlPathErrs...)
	allErrs = append(allErrs, headerErrs...)
	allErrs = append(allErrs, queryErrs...)

	return vo.NewHTTPBackendRequest(host, http.Method(), urlPath, header, query, body), allErrs
}

func (f backendFactory) BuildPublisherRequest(
	request *vo.HTTPRequest,
	history *aggregate.History,
	publisher *vo.Publisher,
) (*vo.PublisherBackendRequest, []error) {
	groupIDSpec := vo.NewGroupIDPipelineSpecFromPublisher(publisher)
	dedupIDSpec := vo.NewDeduplicationIDPipelineSpecFromPublisher(publisher)
	attrsSpec := vo.NewPublisherAttributesPipelineSpecFromPublisher(publisher)
	bodySpec := vo.NewBodyPipelineSpecFromPublisherMessage(publisher.Message())

	groupID, groupErrs := f.buildPipelineService.ApplyGroupID(groupIDSpec, request, history)
	dedupID, dedupErrs := f.buildPipelineService.ApplyDeduplicationID(dedupIDSpec, request, history)
	attributes, attrErrs := f.buildPipelineService.ApplyPublisherAttributes(attrsSpec, request, history)
	body, bodyErrs := f.buildPipelineService.ApplyBody(bodySpec, request.Body(), request, history)

	var allErrs []error

	allErrs = append(allErrs, groupErrs...)
	allErrs = append(allErrs, dedupErrs...)
	allErrs = append(allErrs, attrErrs...)
	allErrs = append(allErrs, bodyErrs...)

	return vo.NewPublisherBackendRequest(
		publisher.Broker(),
		publisher.Path(),
		groupID,
		dedupID,
		publisher.Delay(),
		attributes,
		body,
	), allErrs
}

func (f backendFactory) BuildHTTPResponse(httpResponse *http.Response) *vo.HTTPBackendResponse {
	statusCode := vo.NewStatusCode(httpResponse.StatusCode)
	header := vo.NewHeader(httpResponse.Header)

	var body *vo.Body
	if checker.NonNil(httpResponse.Body) {
		contentType := httpResponse.Header.Get(mapper.ContentType)
		contentEncoding := httpResponse.Header.Get(mapper.ContentEncoding)

		bodyBytes, err := io.ReadAll(httpResponse.Body)
		if checker.NonNil(err) {
			panic(err)
		}

		body = vo.NewBody(contentType, contentEncoding, bytes.NewBuffer(bodyBytes))
	}

	return vo.NewHTTPBackendResponse(enum.BackendOutcomeExecuted, statusCode, header, body)
}

func (f backendFactory) BuildHTTPResponseByErr(endpoint *vo.Endpoint, backend *vo.Backend, err error) *vo.HTTPBackendResponse {
	outcome := enum.BackendOutcomeExecuted

	var code int
	if errors.Is(err, mapper.ErrBackendConcurrentCancelled) {
		code = 499
		outcome = enum.BackendOutcomeCancelled
	} else if errors.Is(err, mapper.ErrBackendGatewayTimeout) || errors.Is(err, context.DeadlineExceeded) {
		code = http.StatusGatewayTimeout
	} else if errors.Is(err, mapper.ErrBackendBadGateway) {
		code = http.StatusBadGateway
	} else if errors.Is(err, mapper.ErrEvalGuards) {
		code = 499
		outcome = enum.BackendOutcomeIgnored
	}
	statusCode := vo.NewStatusCode(code)

	wrapped := errors.Wrap(err)
	buffer := converter.ToBuffer(dto.ErrorBody{
		ID:        backend.ID(),
		File:      wrapped.File(),
		Line:      wrapped.Line(),
		Endpoint:  endpoint.Path(),
		Message:   wrapped.Message(),
		Timestamp: time.Now(),
	})
	body := vo.NewBodyJson(buffer)
	header := vo.NewHeaderByBody(body)

	return vo.NewHTTPBackendResponse(outcome, statusCode, header, body)
}

func (f backendFactory) BuildPublisherResponse(publisherResponse *publisher.Response) *vo.PublisherBackendResponse {
	return vo.NewPublisherBackendResponse(
		enum.BackendOutcomeExecuted,
		publisherResponse.OK,
		vo.NewBodyJson(converter.ToBuffer(publisherResponse.Body)),
	)
}

func (f backendFactory) BuildPublisherResponseByErr(endpoint *vo.Endpoint, backend *vo.Backend, err error,
) *vo.PublisherBackendResponse {
	outcome := enum.BackendOutcomeExecuted
	if errors.Is(err, mapper.ErrBackendConcurrentCancelled) {
		outcome = enum.BackendOutcomeCancelled
	} else if errors.Is(err, mapper.ErrEvalGuards) {
		outcome = enum.BackendOutcomeIgnored
	}

	wrapped := errors.Wrap(err)
	buffer := converter.ToBuffer(dto.ErrorBody{
		ID:        backend.ID(),
		File:      wrapped.File(),
		Line:      wrapped.Line(),
		Endpoint:  endpoint.Path(),
		Message:   wrapped.Message(),
		Timestamp: time.Now(),
	})
	body := vo.NewBodyJson(buffer)

	return vo.NewPublisherBackendResponse(outcome, false, body)
}

func (f backendFactory) BuildFinalHTTPResponse(
	backend *vo.Backend,
	response *vo.HTTPBackendResponse,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (*vo.HTTPBackendResponse, []error) {
	if !backend.HasResponse() || !response.Executed() {
		return response, nil
	} else if backend.Response().Omit() {
		return nil, nil
	}

	bodyPipelineSpec := vo.NewBodyPipelineSpecFromBackendResponse(backend.Response())
	headerPipelineSpec := vo.NewHeaderPipelineSpecFromBackendResponse(backend.Response())

	body, bodyErrs := f.buildPipelineService.ApplyBody(bodyPipelineSpec, response.Body(), request, history)
	header, headerErrs := f.buildPipelineService.ApplyHeader(headerPipelineSpec, vo.NewHeaderByBody(body), request, history)

	return vo.NewHTTPBackendResponse(
		enum.BackendOutcomeExecuted,
		response.StatusCode(),
		header,
		body,
	), append(bodyErrs, headerErrs...)
}

func (f backendFactory) BuildFinalPublisherResponse(
	backend *vo.Backend,
	response *vo.PublisherBackendResponse,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (*vo.PublisherBackendResponse, []error) {
	if !backend.HasResponse() || !response.Executed() {
		return response, nil
	} else if backend.Response().Omit() {
		return nil, nil
	}

	bodyPipelineSpec := vo.NewBodyPipelineSpecFromBackendResponse(backend.Response())

	body, bodyErrs := f.buildPipelineService.ApplyBody(bodyPipelineSpec, response.Body(), request, history)

	return vo.NewPublisherBackendResponse(enum.BackendOutcomeExecuted, response.OK(), body), bodyErrs
}
