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
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/app/model/publisher"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type backendFactory struct {
	mapperService       service.Mapper
	projectorService    service.Projector
	dynamicValueService service.DynamicValue
	modifierService     service.Modifier
	enricherService     service.Enricher
	omitterService      service.Omitter
	nomenclatureService service.Nomenclature
	contentService      service.Content
	aggregatorService   service.Aggregator
}

type Backend interface {
	BuildHTTPRequest(http *vo.HTTP, request *vo.HTTPRequest, history *vo.History) (*vo.HTTPBackendRequest, []error)
	BuildPublisherRequest(request *vo.HTTPRequest, history *vo.History, publisher *vo.Publisher) (*vo.PublisherBackendRequest, []error)
	BuildHTTPResponse(httpResponse *http.Response) *vo.HTTPBackendResponse
	BuildHTTPResponseByErr(endpoint *vo.Endpoint, backend *vo.Backend, err error) *vo.HTTPBackendResponse
	BuildPublisherResponse(publisherResponse *publisher.Response) *vo.PublisherBackendResponse
	BuildPublisherResponseByErr(endpoint *vo.Endpoint, backend *vo.Backend, err error) *vo.PublisherBackendResponse
	BuildFinalHTTPResponse(backend *vo.Backend, response *vo.HTTPBackendResponse, request *vo.HTTPRequest, history *vo.History) (*vo.HTTPBackendResponse, []error)
	BuildFinalPublisherResponse(backend *vo.Backend, response *vo.PublisherBackendResponse, request *vo.HTTPRequest, history *vo.History) (*vo.PublisherBackendResponse, []error)
}

func NewBackend(mapperService service.Mapper, projectorService service.Projector,
	dynamicValueService service.DynamicValue, modifierService service.Modifier,
	enricherService service.Enricher, omitterService service.Omitter,
	nomenclatureService service.Nomenclature, contentService service.Content, aggregatorService service.Aggregator,
) Backend {
	return backendFactory{
		mapperService:       mapperService,
		projectorService:    projectorService,
		dynamicValueService: dynamicValueService,
		modifierService:     modifierService,
		enricherService:     enricherService,
		omitterService:      omitterService,
		nomenclatureService: nomenclatureService,
		contentService:      contentService,
		aggregatorService:   aggregatorService,
	}
}

func (f backendFactory) BuildHTTPRequest(http *vo.HTTP, request *vo.HTTPRequest, history *vo.History) (
	*vo.HTTPBackendRequest, []error) {
	host := f.buildHTTPRequestBalancedHost(http)
	body, bodyErrs := f.buildHTTPRequestBody(http, request, history)
	urlPath, urlPathErrs := f.buildHTTPRequestURLPath(http, request, history)
	header, headerErrs := f.buildHTTPRequestHeader(http, body, request, history)
	query, queryErrs := f.buildHTTPRequestQuery(http, request, history)

	var allErrs []error
	allErrs = append(allErrs, bodyErrs...)
	allErrs = append(allErrs, urlPathErrs...)
	allErrs = append(allErrs, headerErrs...)
	allErrs = append(allErrs, queryErrs...)

	return vo.NewHTTPBackendRequest(host, http.Method(), urlPath, header, query, body), allErrs
}

func (f backendFactory) BuildPublisherRequest(request *vo.HTTPRequest, history *vo.History, publisher *vo.Publisher) (
	*vo.PublisherBackendRequest, []error) {
	groupID, groupErrs := f.dynamicValueService.Get(publisher.GroupID(), request, history)
	if checker.IsNotEmpty(groupErrs) {
		return nil, groupErrs
	}

	deduplicateID, deduplicateErrs := f.dynamicValueService.Get(publisher.DeduplicationID(), request, history)
	if checker.IsNotEmpty(deduplicateErrs) {
		return nil, deduplicateErrs
	}

	var allErrors []error

	attributes := map[string]vo.PublisherMessageAttribute{}
	body := request.Body()
	if publisher.HasMessage() {
		message := publisher.Message()

		for key, attribute := range message.Attributes() {
			value, errs := f.dynamicValueService.Get(attribute.Value(), request, history)

			attributes[key] = vo.NewPublisherMessageAttribute(attribute.DataType(), value)

			allErrors = append(allErrors, errs...)
		}

		if publisher.Message().HasBody() {
			var modifyErrs, mapErrs, projectErrs []error

			body, modifyErrs = f.omitEmptyValuesFromBody(message.Body().OmitEmpty(), body)
			body, modifyErrs = f.modifierService.ExecuteBodyModifiers(message.Body().Modifiers(), body, request, history)
			body, mapErrs = f.mapperService.MapBody(message.Body().Mapper(), body, request, history)
			body, projectErrs = f.projectorService.ProjectBody(message.Body().Projector(), body, request, history)

			allErrors = append(allErrors, modifyErrs...)
			allErrors = append(allErrors, mapErrs...)
			allErrors = append(allErrors, projectErrs...)
		}
	}

	bodyString, err := body.CompactString()
	if checker.NonNil(err) {
		return nil, append(allErrors, err)
	}

	return vo.NewPublisherBackendRequest(
		publisher.Broker(),
		publisher.Path(),
		groupID,
		deduplicateID,
		publisher.Delay(),
		attributes,
		bodyString,
	), nil
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

	return vo.NewHTTPBackendResponse(statusCode, header, body)
}

func (f backendFactory) BuildHTTPResponseByErr(endpoint *vo.Endpoint, backend *vo.Backend, err error) *vo.HTTPBackendResponse {
	if errors.Is(err, mapper.ErrBackendConcurrentCancelled) {
		return nil
	}

	var code int
	if errors.Is(err, mapper.ErrBackendGatewayTimeout) {
		code = http.StatusGatewayTimeout
	} else if errors.Is(err, mapper.ErrBackendBadGateway) {
		code = http.StatusBadGateway
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

	return vo.NewHTTPBackendResponse(statusCode, header, body)
}

func (f backendFactory) BuildPublisherResponse(publisherResponse *publisher.Response) *vo.PublisherBackendResponse {
	return vo.NewPublisherBackendResponse(publisherResponse.OK, vo.NewBodyJson(converter.ToBuffer(publisherResponse.Body)))
}

func (f backendFactory) BuildPublisherResponseByErr(endpoint *vo.Endpoint, backend *vo.Backend, err error,
) *vo.PublisherBackendResponse {
	if errors.Is(err, mapper.ErrBackendConcurrentCancelled) {
		return nil
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

	return vo.NewPublisherBackendResponse(false, body)
}

func (f backendFactory) BuildFinalHTTPResponse(
	backend *vo.Backend,
	response *vo.HTTPBackendResponse,
	request *vo.HTTPRequest,
	history *vo.History,
) (*vo.HTTPBackendResponse, []error) {
	if !backend.HasResponse() {
		return response, nil
	} else if backend.Response().Omit() {
		return nil, nil
	}

	body, bodyErrs := f.buildFinalHTTPResponseBody(backend, response, request, history)
	header, headerErrs := f.buildFinalHTTPResponseHeader(backend, response, body, request, history)

	var allErrs []error
	allErrs = append(allErrs, bodyErrs...)
	allErrs = append(allErrs, headerErrs...)

	return vo.NewHTTPBackendResponse(response.StatusCode(), header, body), allErrs
}

func (f backendFactory) BuildFinalPublisherResponse(
	backend *vo.Backend,
	response *vo.PublisherBackendResponse,
	request *vo.HTTPRequest,
	history *vo.History,
) (*vo.PublisherBackendResponse, []error) {
	if !backend.HasResponse() {
		return response, nil
	} else if backend.Response().Omit() {
		return nil, nil
	}

	body := response.Body()
	body, modifyErrs := f.modifierService.ExecuteBodyModifiers(backend.Response().Body().Modifiers(), body, request, history)
	body, mapErrs := f.mapperService.MapBody(backend.Response().Body().Mapper(), body, request, history)
	body, projectErrs := f.projectorService.ProjectBody(backend.Response().Body().Projector(), body, request, history)
	body, aggregateErr := f.aggregatorService.AggregateBodyToKey(backend.Response().Body().Group(), body)

	var allErrs []error
	allErrs = append(allErrs, modifyErrs...)
	allErrs = append(allErrs, mapErrs...)
	allErrs = append(allErrs, projectErrs...)
	if checker.NonNil(aggregateErr) {
		allErrs = append(allErrs, aggregateErr)
	}

	return vo.NewPublisherBackendResponse(response.OK(), body), allErrs
}

func (f backendFactory) buildHTTPRequestBalancedHost(http *vo.HTTP) string {
	// todo: aqui obtemos o host (correto é criar um domínio chamado balancer aonde ele vai retornar o host
	//  disponível pegando como base, se ele esta de pé ou não, e sua config de porcentagem)
	if checker.IsLengthEquals(http.Hosts(), 1) {
		return http.Hosts()[0]
	}
	return http.Hosts()[rand.Intn(len(http.Hosts())-1)]
}

func (f backendFactory) buildHTTPRequestBody(http *vo.HTTP, request *vo.HTTPRequest, history *vo.History) (
	*vo.Body, []error) {
	if checker.IsNil(request.Body()) || !http.HasRequest() || !http.Request().HasBody() {
		return request.Body(), nil
	} else if http.Request().Body().Omit() {
		return nil, nil
	}

	body := request.Body()

	body, modifyErrs := f.modifierService.ExecuteBodyModifiers(http.Request().Body().Modifiers(), body, request, history)
	body, mapErrs := f.mapperService.MapBody(http.Request().Body().Mapper(), body, request, history)
	body, projectErrs := f.projectorService.ProjectBody(http.Request().Body().Projector(), body, request, history)

	body, omitErrs := f.omitEmptyValuesFromBody(http.Request().Body().OmitEmpty(), body)
	body, modifyCaseErrs := f.nomenclatureService.ToCase(http.Request().Body().Nomenclature(), body)
	body, modifyContentTypeErr := f.contentService.ModifyBodyContentType(http.Request().Body().ContentType(), body)
	body, modifyBodyContentEncodingErr := f.contentService.ModifyBodyContentEncoding(http.Request().Body().ContentEncoding(), body)

	var allErrors []error
	allErrors = append(allErrors, mapErrs...)
	allErrors = append(allErrors, projectErrs...)
	allErrors = append(allErrors, modifyErrs...)
	allErrors = append(allErrors, omitErrs...)
	allErrors = append(allErrors, modifyCaseErrs...)

	if checker.NonNil(modifyContentTypeErr) {
		allErrors = append(allErrors, modifyContentTypeErr)
	}
	if checker.NonNil(modifyBodyContentEncodingErr) {
		allErrors = append(allErrors, modifyBodyContentEncodingErr)
	}

	return body, allErrors
}

func (f backendFactory) buildHTTPRequestURLPath(http *vo.HTTP, request *vo.HTTPRequest, history *vo.History) (
	vo.URLPath, []error) {
	urlPath := vo.NewURLPath(http.Path(), request.Params().Copy())
	if !http.HasRequest() || !http.Request().HasParam() {
		return urlPath, nil
	}

	return f.modifierService.ExecuteURLPathModifiers(http.Request().Param().Modifiers(), urlPath, request, history)
}

func (f backendFactory) buildHTTPRequestHeader(http *vo.HTTP, body *vo.Body, request *vo.HTTPRequest,
	history *vo.History) (vo.Header, []error) {
	header := vo.NewHeaderByBody(body)

	if http.HasRequest() && http.Request().HasHeader() && http.Request().Header().Omit() {
		return header, nil
	}

	header = f.aggregatorService.AggregateHeaders(header, request.Header())

	if !http.HasRequest() || !http.Request().HasHeader() {
		return header, nil
	}

	header, modifierErrs := f.modifierService.ExecuteHeaderModifiers(http.Request().Header().Modifiers(), header, request, history)
	header, mapErr := f.mapperService.MapHeader(http.Request().Header().Mapper(), header, request, history)
	header, projectErr := f.projectorService.ProjectHeader(http.Request().Header().Projector(), header, request, history)

	var allErrors []error
	allErrors = append(allErrors, modifierErrs...)
	if checker.NonNil(mapErr) {
		allErrors = append(allErrors, mapErr)
	}
	if checker.NonNil(projectErr) {
		allErrors = append(allErrors, projectErr)
	}

	return header, allErrors
}

func (f backendFactory) buildHTTPRequestQuery(http *vo.HTTP, request *vo.HTTPRequest, history *vo.History) (
	vo.Query, []error) {
	if !http.HasRequest() || !http.Request().HasQuery() {
		return request.Query(), nil
	} else if http.Request().HasQuery() && http.Request().Query().Omit() {
		return vo.NewEmptyQuery(), nil
	}

	query := request.Query()

	query, modifierErrs := f.modifierService.ExecuteQueryModifiers(http.Request().Query().Modifiers(), query, request, history)
	query, mapErr := f.mapperService.MapQuery(http.Request().Query().Mapper(), query, request, history)
	query, projectErr := f.projectorService.ProjectQuery(http.Request().Query().Projector(), query, request, history)

	var allErrs []error
	allErrs = append(allErrs, modifierErrs...)
	if checker.NonNil(mapErr) {
		allErrs = append(allErrs, mapErr)
	}
	if checker.NonNil(projectErr) {
		allErrs = append(allErrs, projectErr)
	}

	return query, allErrs
}

func (f backendFactory) buildFinalHTTPResponseBody(
	backend *vo.Backend,
	response *vo.HTTPBackendResponse,
	request *vo.HTTPRequest,
	history *vo.History,
) (*vo.Body, []error) {
	if response.StatusCode().Failed() || !response.HasBody() ||
		(backend.Response().HasBody() && backend.Response().Body().Omit()) {
		return nil, nil
	}

	body := response.Body()
	body, enrichErrs := f.enricherService.ExecuteBodyEnrichers(backend.Response().Body().Enrichers(), body, request, history)
	body, modifyErrs := f.modifierService.ExecuteBodyModifiers(backend.Response().Body().Modifiers(), body, request, history)
	body, mapErrs := f.mapperService.MapBody(backend.Response().Body().Mapper(), body, request, history)
	body, projectErrs := f.projectorService.ProjectBody(backend.Response().Body().Projector(), body, request, history)
	body, aggregateErr := f.aggregatorService.AggregateBodyToKey(backend.Response().Body().Group(), body)

	var allErrors []error
	allErrors = append(allErrors, enrichErrs...)
	allErrors = append(allErrors, modifyErrs...)
	allErrors = append(allErrors, mapErrs...)
	allErrors = append(allErrors, projectErrs...)
	if checker.NonNil(aggregateErr) {
		allErrors = append(allErrors, aggregateErr)
	}

	return body, allErrors
}

func (f backendFactory) buildFinalHTTPResponseHeader(
	backend *vo.Backend,
	response *vo.HTTPBackendResponse,
	body *vo.Body,
	request *vo.HTTPRequest,
	history *vo.History,
) (vo.Header, []error) {
	header := vo.NewHeaderByBody(body)

	if backend.HasResponse() && backend.Response().HasHeader() && backend.Response().Header().Omit() {
		return header, nil
	}

	header = f.aggregatorService.AggregateHeaders(header, response.Header())

	if !backend.HasResponse() || !backend.Response().HasHeader() {
		return header, nil
	}

	header, modifyErrs := f.modifierService.ExecuteHeaderModifiers(backend.Response().Header().Modifiers(), header, request, history)
	header, mapErr := f.mapperService.MapHeader(backend.Response().Header().Mapper(), header, request, history)
	header, projectErr := f.projectorService.ProjectHeader(backend.Response().Header().Projector(), header, request, history)

	var allErrors []error
	allErrors = append(allErrors, modifyErrs...)
	if checker.NonNil(mapErr) {
		allErrors = append(allErrors, mapErr)
	}
	if checker.NonNil(projectErr) {
		allErrors = append(allErrors, projectErr)
	}
	return header, allErrors
}

func (f backendFactory) omitEmptyValuesFromBody(omitEmpty bool, body *vo.Body) (*vo.Body, []error) {
	if !omitEmpty {
		return body, nil
	}
	return f.omitterService.OmitEmptyValuesFromBody(body)
}
