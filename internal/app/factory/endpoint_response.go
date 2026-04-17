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
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type endpointResponse struct {
	aggregatorService    service.Aggregator
	buildPipelineService service.BuildPipeline
}

type EndpointResponse interface {
	BuildErrorResponse(endpoint *vo.EndpointConfig, err error) *vo.EndpointResponse
	BuildAbortedResponse(history *aggregate.History) *vo.EndpointResponse
	BuildResponse(endpoint *vo.EndpointConfig, request *vo.EndpointRequest, history *aggregate.History) (
		*vo.EndpointResponse, []error)
}

func NewEndpointResponse(aggregatorService service.Aggregator, buildPipelineService service.BuildPipeline) EndpointResponse {
	return endpointResponse{
		aggregatorService:    aggregatorService,
		buildPipelineService: buildPipelineService,
	}
}

func (f endpointResponse) BuildErrorResponse(endpoint *vo.EndpointConfig, err error) *vo.EndpointResponse {
	status := vo.NewResponseStatusByValue(enum.ResponseStatusInternalError)

	wrapped := errors.Wrap(err)
	payload := vo.NewPayloadJSON(converter.ToBuffer(dto.ErrorPayload{
		File:      wrapped.File(),
		Line:      wrapped.Line(),
		Endpoint:  endpoint.Path(),
		Message:   wrapped.Message(),
		Stack:     wrapped.Stack(),
		Timestamp: time.Now(),
	}))
	metadata := vo.NewEmptyMetadata()

	return vo.NewEndpointResponse(status, metadata, payload)
}

func (f endpointResponse) BuildAbortedResponse(history *aggregate.History) *vo.EndpointResponse {
	lastestResponse := history.GetResponseLastest()

	lastestStatusCode := lastestResponse.Status()
	lastestMetadata := lastestResponse.Metadata()
	lastestPayload := lastestResponse.Payload()

	degradation := lastestResponse.Degradation()

	return vo.NewEndpointResponseWithExecution(
		degradation,
		f.buildEndpointExecution(history),
		lastestStatusCode,
		lastestMetadata,
		lastestPayload,
	)
}

func (f endpointResponse) BuildResponse(
	endpoint *vo.EndpointConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.EndpointResponse, []error) {
	status := f.buildStatusByHistory(history)
	payload, payloadDegraded, payloadErrs := f.buildPayloadByHistory(endpoint, request, history)
	metadata, metadataDegraded, metadataErrs := f.buildMetadataByHistory(endpoint, request, history)

	var degradationKinds []enum.DegradationKind
	if payloadDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindPayload)
	}
	if metadataDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindMetadata)
	}

	return vo.NewEndpointResponseWithBackendCache(
		vo.NewEmptyCacheInfo(),
		vo.NewDegradation(degradationKinds...),
		f.buildEndpointExecution(history),
		history.BackendsCachedIDs(),
		history.AllBackendsFromCache(),
		history.NewestBackendCacheTTLMillis(),
		status,
		metadata,
		payload,
	), joinErrs(payloadErrs, metadataErrs)
}

func (f endpointResponse) buildStatusByHistory(history *aggregate.History) vo.ResponseStatus {
	if history.IsMultipleFinalResponse() {
		return f.buildStatusFromMultipleResponses(history)
	} else if history.IsSingleFinalResponse() {
		return history.GetResponseLastest().Status()
	} else {
		return vo.NewResponseStatusByValue(enum.ResponseStatusOK)
	}
}

func (f endpointResponse) buildPayloadByHistory(
	endpoint *vo.EndpointConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.Payload, bool, []error) {
	var payload *vo.Payload
	var buildDegraded bool
	var errs []error

	if history.IsMultipleFinalResponse() {
		payload, errs = f.buildPayloadFromMultipleResponses(endpoint, history)
		buildDegraded = checker.NonNil(errs)
	} else if history.IsSingleFinalResponse() {
		responseLastest := history.GetResponseLastest()

		payload = responseLastest.Payload()
		buildDegraded = responseLastest.PayloadDegraded()
	}

	if !endpoint.Response().HasPayload() || checker.IsNil(payload) {
		return payload, buildDegraded, errs
	}

	payload, payloadErrs := f.buildPipelineService.ApplyPayload(endpoint.Response().Payload(), payload, request, history)

	return payload, buildDegraded && isDegraded(endpoint.Response().Payload(), payloadErrs), joinErrs(errs, payloadErrs)
}

func (f endpointResponse) buildMetadataByHistory(
	endpoint *vo.EndpointConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Metadata, bool, []error) {
	var buildDegraded bool

	metadata := vo.NewEmptyMetadata()
	for i := 0; checker.IsLessThan(i, history.Size()); i++ {
		response := history.GetResponse(i)
		if history.ShouldContributeMetadata(i) {
			metadata = f.aggregatorService.AggregateMetadata(metadata, response.Metadata(), app.TransportHTTPHeaderKeys())
			buildDegraded = buildDegraded || response.MetadataDegraded()
		}
	}
	metadata, errs := f.buildPipelineService.ApplyMetadata(
		endpoint.Response().Metadata(),
		metadata,
		app.TransportHTTPHeaderKeys(),
		request,
		history,
	)

	return metadata, buildDegraded && isDegraded(endpoint.Response().Metadata(), errs), errs
}

func (f endpointResponse) buildStatusFromMultipleResponses(history *aggregate.History) vo.ResponseStatus {
	statuses := make(map[vo.ResponseStatus]int)
	for i := 0; checker.IsLessThan(i, history.Size()); i++ {
		response := history.GetResponse(i)
		if history.ShouldBeInFinalResponse(i) {
			statuses[response.Status()]++
		}
	}

	mostFrequentStatus := vo.NewResponseStatusByValue(enum.ResponseStatusOK)
	maxCount := 0
	for status, count := range statuses {
		if checker.IsGreaterThanOrEqual(count, maxCount) {
			mostFrequentStatus = status
			maxCount = count
		}
	}

	return mostFrequentStatus
}

func (f endpointResponse) buildPayloadFromMultipleResponses(endpoint *vo.EndpointConfig, history *aggregate.History) (
	*vo.Payload, []error) {
	if endpoint.Response().HasPayload() && endpoint.Response().Payload().Aggregate() {
		return f.aggregatorService.AggregatePayloads(history)
	} else {
		return f.aggregatorService.AggregatePayloadsIntoSlice(history)
	}
}

func (f endpointResponse) buildEndpointExecution(history *aggregate.History) vo.EndpointExecution {
	return vo.NewEndpointExecution(history.AllExecuted(), history.AllOK(), history.Degradations())
}
