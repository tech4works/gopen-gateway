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
	"github.com/tech4works/checker"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type backendRequest struct {
	buildPipelineService service.BuildPipeline
}

type BackendRequest interface {
	BuildHTTPRequest(backend *vo.BackendConfig, request *vo.EndpointRequest, history *aggregate.History) (
		*vo.HTTPBackendRequest, []error)
	BuildPublisherRequest(backend *vo.BackendConfig, request *vo.EndpointRequest, history *aggregate.History) (
		*vo.PublisherBackendRequest, []error)
}

func NewBackendRequest(buildPipelineService service.BuildPipeline) BackendRequest {
	return backendRequest{
		buildPipelineService: buildPipelineService,
	}
}

func (f backendRequest) BuildHTTPRequest(
	backend *vo.BackendConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.HTTPBackendRequest, []error) {
	backendHTTP := backend.HTTP()
	useFallback := backend.Execution().UseFallback(enum.ExecutionOnBuild)

	var (
		body    *vo.Payload
		urlPath vo.URLPath
		query   vo.Query
		header  vo.Metadata
	)
	var bodyDegraded, urlPathDegraded, queryDegraded, headerDegraded bool
	var bodyErrs, urlPathErrs, queryErrs, headerErrs []error

	host := f.buildPipelineService.ApplyHost(backendHTTP)
	body, bodyDegraded, bodyErrs = f.buildHTTPRequestBody(backendHTTP.Request().Body(), request, history, useFallback)
	urlPath, urlPathDegraded, urlPathErrs = f.buildHTTPRequestURLPath(backendHTTP, request, history, useFallback)
	query, queryDegraded, queryErrs = f.buildHTTPRequestQuery(backendHTTP.Request().Query(), request, history, useFallback)
	header, headerDegraded, headerErrs = f.buildHTTPRequestHeader(backendHTTP.Request().Header(), request, history, useFallback)

	var degradationKinds []enum.DegradationKind
	if bodyDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindPayload)
	}
	if urlPathDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindURLPath)
	}
	if queryDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindQuery)
	}
	if headerDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindMetadata)
	}

	return vo.NewHTTPBackendRequest(
		vo.NewDegradation(degradationKinds...),
		host,
		backendHTTP.Method(),
		urlPath,
		header,
		query,
		body,
	), joinErrs(headerErrs, urlPathErrs, queryErrs, bodyErrs)
}

func (f backendRequest) BuildPublisherRequest(
	backend *vo.BackendConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.PublisherBackendRequest, []error) {
	backendPublisher := backend.Publisher()
	useFallback := backend.Execution().UseFallback(enum.ExecutionOnBuild)

	groupID, groupIDDegraded, groupErrs := f.buildPublisherRequestGroupID(backendPublisher, request, history, useFallback)
	dedupID, dedupDegraded, dedupErrs := f.buildPublisherRequestDedupID(backendPublisher, request, history, useFallback)
	attrs, attrDegraded, attrErrs := f.buildPublisherRequestAttributes(backendPublisher.Message(), request, history,
		useFallback)
	body, bodyDegraded, bodyErrs := f.buildPublisherRequestBody(backendPublisher.Message().Body(), request, history,
		useFallback)

	var degradationKinds []enum.DegradationKind
	if groupIDDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindGroupID)
	}
	if dedupDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindDeduplicationID)
	}
	if attrDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindAttributes)
	}
	if bodyDegraded {
		degradationKinds = append(degradationKinds, enum.DegradationKindPayload)
	}

	return vo.NewPublisherBackendRequest(
		vo.NewDegradation(degradationKinds...),
		backendPublisher.Broker(),
		backendPublisher.Path(),
		groupID,
		dedupID,
		backendPublisher.Delay(),
		attrs,
		body,
	), joinErrs(groupErrs, dedupErrs, attrErrs, bodyErrs)
}

func (f backendRequest) buildHTTPRequestBody(
	spec vo.PayloadPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
	useFallback bool,
) (*vo.Payload, bool, []error) {
	body, errs := f.buildPipelineService.ApplyPayload(spec, request.Payload(), request, history)
	return fallbackIf(useFallback, errs, body, request.Payload()), isDegraded(spec, errs), errs
}

func (f backendRequest) buildHTTPRequestURLPath(
	httpConfig *vo.BackendHTTPConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
	useFallback bool,
) (vo.URLPath, bool, []error) {
	spec := httpConfig.Request().URLPath()
	urlPathRaw := vo.NewURLPath(httpConfig.Path(), request.Params().Copy())

	urlPath, errs := f.buildPipelineService.ApplyURLPath(spec, urlPathRaw, request, history)
	return fallbackIf(useFallback, errs, urlPath, urlPathRaw), isDegraded(spec, errs), errs
}

func (f backendRequest) buildHTTPRequestQuery(
	spec vo.QueryPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
	useFallback bool,
) (vo.Query, bool, []error) {
	query, errs := f.buildPipelineService.ApplyQuery(spec, request.Query(), request, history)
	return fallbackIf(useFallback, errs, query, request.Query()), isDegraded(spec, errs), errs
}

func (f backendRequest) buildHTTPRequestHeader(
	spec vo.MetadataPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
	useFallback bool,
) (vo.Metadata, bool, []error) {
	baseMap := map[string][]string{}
	for k, v := range request.Metadata().Copy() {
		if checker.NotContains(app.TransportHTTPHeaderKeys(), k) {
			baseMap[k] = v
		}
	}

	base := vo.NewMetadata(baseMap)
	header, errs := f.buildPipelineService.ApplyMetadata(spec, base, app.TransportHTTPHeaderKeys(), request, history)
	return fallbackIf(useFallback, errs, header, request.Metadata()), isDegraded(spec, errs), errs
}

func (f backendRequest) buildPublisherRequestGroupID(
	spec vo.GroupIDPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
	_ bool,
) (string, bool, []error) {
	groupID, errs := f.buildPipelineService.ApplyGroupID(spec, request, history)
	return groupID, isDegraded(spec, errs), errs
}

func (f backendRequest) buildPublisherRequestDedupID(
	spec vo.DeduplicationIDPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
	_ bool,
) (string, bool, []error) {
	deduplicationID, errs := f.buildPipelineService.ApplyDeduplicationID(spec, request, history)
	return deduplicationID, isDegraded(spec, errs), errs
}

func (f backendRequest) buildPublisherRequestAttributes(
	spec vo.AttributesPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
	_ bool,
) (map[string]vo.AttributeValueConfig, bool, []error) {
	attributes, errs := f.buildPipelineService.ApplyAttributes(spec, request, history)
	return attributes, isDegraded(spec, errs), errs
}

func (f backendRequest) buildPublisherRequestBody(
	spec vo.PayloadPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
	useFallback bool,
) (*vo.Payload, bool, []error) {
	body, errs := f.buildPipelineService.ApplyPayload(spec, request.Payload(), request, history)
	return fallbackIf(useFallback, errs, body, request.Payload()), isDegraded(spec, errs), errs
}
