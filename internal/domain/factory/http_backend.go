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
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type httpBackendFactory struct {
	mapperService       service.Mapper
	projectorService    service.Projector
	dynamicValueService service.DynamicValue
	modifierService     service.Modifier
	omitterService      service.Omitter
	nomenclatureService service.Nomenclature
	contentService      service.Content
	aggregatorService   service.Aggregator
}

type HTTPBackend interface {
	BuildRequest(backend *vo.Backend, request *vo.HTTPRequest, history *vo.History) (*vo.HTTPBackendRequest, []error)
	BuildTemporaryResponse(httpResponse *http.Response) *vo.HTTPBackendResponse
	BuildTemporaryResponseByErr(endpoint *vo.Endpoint, err error) *vo.HTTPBackendResponse
	BuildResponse(backend *vo.Backend, temporaryResponse *vo.HTTPBackendResponse, request *vo.HTTPRequest, history *vo.History) (*vo.HTTPBackendResponse, []error)
}

func NewHTTPBackend(mapperService service.Mapper, projectorService service.Projector,
	dynamicValueService service.DynamicValue, modifierService service.Modifier, omitterService service.Omitter,
	nomenclatureService service.Nomenclature, contentService service.Content, aggregatorService service.Aggregator,
) HTTPBackend {
	return httpBackendFactory{
		mapperService:       mapperService,
		projectorService:    projectorService,
		dynamicValueService: dynamicValueService,
		modifierService:     modifierService,
		omitterService:      omitterService,
		nomenclatureService: nomenclatureService,
		contentService:      contentService,
		aggregatorService:   aggregatorService,
	}
}

func (f httpBackendFactory) BuildRequest(backend *vo.Backend, request *vo.HTTPRequest, history *vo.History) (
	*vo.HTTPBackendRequest, []error) {
	host := f.buildRequestBalancedHost(backend)
	body, bodyErrs := f.buildRequestBody(backend, request, history)
	urlPath, urlPathErrs := f.buildRequestURLPath(backend, request, history)
	header, headerErrs := f.buildRequestHeader(backend, body, request, history)
	query, queryErrs := f.buildRequestQuery(backend, request, history)

	var allErrs []error
	allErrs = append(allErrs, bodyErrs...)
	allErrs = append(allErrs, urlPathErrs...)
	allErrs = append(allErrs, headerErrs...)
	allErrs = append(allErrs, queryErrs...)

	return vo.NewHTTPBackendRequest(host, backend.Method(), urlPath, header, query, body), allErrs
}

func (f httpBackendFactory) BuildTemporaryResponse(httpResponse *http.Response) *vo.HTTPBackendResponse {
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

func (f httpBackendFactory) BuildTemporaryResponseByErr(endpoint *vo.Endpoint, err error) *vo.HTTPBackendResponse {
	if errors.Is(err, mapper.ErrConcurrentCanceled) {
		return nil
	}

	var code int
	if errors.Is(err, mapper.ErrGatewayTimeout) {
		code = http.StatusGatewayTimeout
	} else if errors.Is(err, mapper.ErrBadGateway) {
		code = http.StatusBadGateway
	}
	statusCode := vo.NewStatusCode(code)

	details := errors.Details(err)
	buffer := converter.ToBuffer(dto.ErrorBody{
		File:      details.File(),
		Line:      details.Line(),
		Endpoint:  endpoint.Path(),
		Message:   details.Message(),
		Timestamp: time.Now(),
	})
	body := vo.NewBodyJson(buffer)
	header := vo.NewHeaderByBody(body)

	return vo.NewHTTPBackendResponse(statusCode, header, body)
}

func (f httpBackendFactory) BuildResponse(backend *vo.Backend, temporaryResponse *vo.HTTPBackendResponse,
	request *vo.HTTPRequest, history *vo.History) (*vo.HTTPBackendResponse, []error) {
	if !backend.HasResponse() {
		return temporaryResponse, nil
	} else if backend.Response().Omit() {
		return nil, nil
	}

	body, bodyErrs := f.buildResponseBody(backend, temporaryResponse, request, history)
	header, headerErrs := f.buildResponseHeader(backend, temporaryResponse, body, request, history)

	var allErrs []error
	allErrs = append(allErrs, bodyErrs...)
	allErrs = append(allErrs, headerErrs...)

	return vo.NewHTTPBackendResponse(temporaryResponse.StatusCode(), header, body), allErrs
}

func (f httpBackendFactory) buildRequestBalancedHost(backend *vo.Backend) string {
	// todo: aqui obtemos o host (correto é criar um domínio chamado balancer aonde ele vai retornar o host
	//  disponível pegando como base, se ele esta de pé ou não, e sua config de porcentagem)
	if checker.IsLengthEquals(backend.Hosts(), 1) {
		return backend.Hosts()[0]
	}
	return backend.Hosts()[rand.Intn(len(backend.Hosts())-1)]
}

func (f httpBackendFactory) buildRequestBody(backend *vo.Backend, request *vo.HTTPRequest, history *vo.History) (
	*vo.Body, []error) {
	if !backend.HasRequest() {
		return request.Body(), nil
	} else if checker.IsNil(request.Body()) || backend.Request().OmitBody() {
		return nil, nil
	}

	body := request.Body()
	body, modifyErrs := f.modifyBody(backend.Request().BodyModifiers(), body, request, history)
	body, mapErrs := f.mapperService.MapBody(body, backend.Request().BodyMapper())
	body, projectErrs := f.projectorService.ProjectBody(body, backend.Request().BodyProjection())
	body, omitErrs := f.omitEmptyValuesFromBody(backend.Request().OmitEmpty(), body)
	body, modifyCaseErrs := f.modifyBodyCase(backend.Request().Nomenclature(), body)
	body, modifyContentTypeErrs := f.modifyBodyContentType(backend.Request().ContentType(), body)
	body, modifyBodyContentEncodingErrs := f.modifyBodyContentEncoding(backend.Request().ContentEncoding(), body)

	var allErrors []error
	allErrors = append(allErrors, mapErrs...)
	allErrors = append(allErrors, projectErrs...)
	allErrors = append(allErrors, modifyErrs...)
	allErrors = append(allErrors, omitErrs...)
	allErrors = append(allErrors, modifyCaseErrs...)
	allErrors = append(allErrors, modifyContentTypeErrs...)
	allErrors = append(allErrors, modifyBodyContentEncodingErrs...)

	return body, allErrors
}

func (f httpBackendFactory) buildRequestURLPath(backend *vo.Backend, request *vo.HTTPRequest, history *vo.History) (
	vo.URLPath, []error) {
	urlPath := vo.NewURLPath(backend.Path(), request.Params().Copy())
	if !backend.HasRequest() {
		return urlPath, nil
	}

	return f.modifyURLPath(backend, urlPath, request, history)
}

func (f httpBackendFactory) buildRequestHeader(backend *vo.Backend, body *vo.Body, request *vo.HTTPRequest,
	history *vo.History) (vo.Header, []error) {
	header := vo.NewHeaderByBody(body)
	if backend.HasRequest() && backend.Request().OmitHeader() {
		return header, nil
	}

	header = f.aggregatorService.AggregateHeaders(header, request.Header())
	if !backend.HasRequest() {
		return header, nil
	}

	header = f.mapperService.MapHeader(header, backend.Request().HeaderMapper())
	header = f.projectorService.ProjectHeader(header, backend.Request().HeaderProjection())

	return f.modifyHeader(backend.Request().HeaderModifiers(), header, request, history)
}

func (f httpBackendFactory) buildRequestQuery(backend *vo.Backend, request *vo.HTTPRequest, history *vo.History) (
	vo.Query, []error) {
	if !backend.HasRequest() {
		return request.Query(), nil
	} else if backend.Request().OmitQuery() {
		return vo.NewEmptyQuery(), nil
	}

	query := request.Query()
	query = f.mapperService.MapQuery(query, backend.Request().QueryMapper())
	query = f.projectorService.ProjectQuery(query, backend.Request().QueryProjection())

	return f.modifyQuery(backend, query, request, history)
}

func (f httpBackendFactory) buildResponseBody(backend *vo.Backend, temporaryResponse *vo.HTTPBackendResponse,
	request *vo.HTTPRequest, history *vo.History) (*vo.Body, []error) {
	if temporaryResponse.StatusCode().Failed() || !temporaryResponse.HasBody() || backend.Response().OmitBody() {
		return nil, nil
	}

	body := temporaryResponse.Body()
	body, modifyErrs := f.modifyBody(backend.Response().BodyModifiers(), body, request, history)
	body, mapErrs := f.mapperService.MapBody(body, backend.Response().BodyMapper())
	body, projectErrs := f.projectorService.ProjectBody(body, backend.Response().BodyProjection())
	body, aggregateErr := f.aggregatorService.AggregateBodyToKey(backend.Response().Group(), body)

	var allErrors []error
	allErrors = append(allErrors, mapErrs...)
	allErrors = append(allErrors, projectErrs...)
	allErrors = append(allErrors, modifyErrs...)
	if checker.NonNil(aggregateErr) {
		allErrors = append(allErrors, aggregateErr)
	}

	return body, allErrors
}

func (f httpBackendFactory) buildResponseHeader(backend *vo.Backend, temporaryResponse *vo.HTTPBackendResponse,
	body *vo.Body, request *vo.HTTPRequest, history *vo.History) (vo.Header, []error) {
	header := vo.NewHeaderByBody(body)

	if backend.Response().OmitHeader() {
		return header, nil
	}

	header = f.aggregatorService.AggregateHeaders(header, temporaryResponse.Header())
	header = f.mapperService.MapHeader(header, backend.Response().HeaderMapper())
	header = f.projectorService.ProjectHeader(header, backend.Response().HeaderProjection())

	return f.modifyHeader(backend.Response().HeaderModifiers(), header, request, history)
}

func (f httpBackendFactory) modifyBody(modifiers []vo.Modifier, body *vo.Body, request *vo.HTTPRequest,
	history *vo.History) (*vo.Body, []error) {
	var errs []error

	for _, bodyModifier := range modifiers {
		modifierValue, dynamicValueErrs := f.dynamicValueService.Get(bodyModifier.Value(), request, history)
		if checker.IsNotEmpty(dynamicValueErrs) {
			errs = append(errs, dynamicValueErrs...)
		}
		modifiedBody, err := f.modifierService.ModifyBody(body, bodyModifier.Action(), bodyModifier.Key(), modifierValue)
		if checker.NonNil(err) {
			errs = append(errs, err)
		}
		body = modifiedBody
	}

	return body, errs
}

func (f httpBackendFactory) omitEmptyValuesFromBody(omitEmpty bool, body *vo.Body) (*vo.Body, []error) {
	if !omitEmpty {
		return body, nil
	}
	return f.omitterService.OmitEmptyValuesFromBody(body)
}

func (f httpBackendFactory) modifyBodyCase(nomenclature enum.Nomenclature, body *vo.Body) (*vo.Body, []error) {
	if !nomenclature.IsEnumValid() {
		return body, nil
	}
	return f.nomenclatureService.ToCase(body, nomenclature)
}

func (f httpBackendFactory) modifyBodyContentType(contentTypeConfig enum.ContentType, body *vo.Body) (*vo.Body, []error) {
	var contentType enum.ContentType
	if contentTypeConfig.IsEnumValid() {
		contentType = contentTypeConfig
	} else {
		contentType = body.ContentType().ToEnum()
	}

	newBody, err := f.contentService.ModifyBodyContentType(body, contentType)
	if checker.NonNil(err) {
		return body, []error{err}
	}

	return newBody, nil
}

func (f httpBackendFactory) modifyBodyContentEncoding(contentEncodingConfig enum.ContentEncoding, body *vo.Body) (
	*vo.Body, []error) {
	var contentEncoding enum.ContentEncoding
	if contentEncodingConfig.IsEnumValid() {
		contentEncoding = contentEncodingConfig
	} else if body.ContentEncoding().IsGzip() {
		contentEncoding = enum.ContentEncodingGzip
	} else if body.ContentEncoding().IsDeflate() {
		contentEncoding = enum.ContentEncodingDeflate
	} else {
		contentEncoding = enum.ContentEncodingNone
	}

	newBody, err := f.contentService.ModifyBodyContentEncoding(body, contentEncoding)
	if checker.NonNil(err) {
		return body, []error{err}
	}

	return newBody, nil
}

func (f httpBackendFactory) modifyURLPath(backend *vo.Backend, urlPath vo.URLPath, request *vo.HTTPRequest,
	history *vo.History) (vo.URLPath, []error) {
	var errs []error

	for _, paramModifier := range backend.Request().ParamModifiers() {
		modifierValue, dynamicValueErrs := f.dynamicValueService.Get(paramModifier.Value(), request, history)
		if checker.IsNotEmpty(dynamicValueErrs) {
			errs = append(errs, dynamicValueErrs...)
		}
		modifiedUrlPath, err := f.modifierService.ModifyUrlPath(urlPath, paramModifier.Action(), paramModifier.Key(), modifierValue)
		if checker.NonNil(err) {
			errs = append(errs, err)
		}
		urlPath = modifiedUrlPath
	}

	return urlPath, errs
}

func (f httpBackendFactory) modifyHeader(modifiers []vo.Modifier, header vo.Header, request *vo.HTTPRequest,
	history *vo.History) (vo.Header, []error) {
	var errs []error

	for _, headerModifier := range modifiers {
		modifierValue, dynamicValueErrs := f.dynamicValueService.GetAsSliceOfString(headerModifier.Value(), request, history)
		if checker.IsNotEmpty(dynamicValueErrs) {
			errs = append(errs, dynamicValueErrs...)
		}
		modifiedHeader, err := f.modifierService.ModifyHeader(header, headerModifier.Action(), headerModifier.Key(), modifierValue)
		if checker.NonNil(err) {
			errs = append(errs, err)
		}
		header = modifiedHeader
	}

	return header, errs
}

func (f httpBackendFactory) modifyQuery(backend *vo.Backend, query vo.Query, request *vo.HTTPRequest, history *vo.History,
) (vo.Query, []error) {
	var errs []error

	for _, queryModifier := range backend.Request().QueryModifiers() {
		modifierValue, dynamicValueErrs := f.dynamicValueService.GetAsSliceOfString(queryModifier.Value(), request, history)
		if checker.IsNotEmpty(dynamicValueErrs) {
			errs = append(errs, dynamicValueErrs...)
		}
		modifiedQuery, err := f.modifierService.ModifyQuery(query, queryModifier.Action(), queryModifier.Key(), modifierValue)
		if checker.NonNil(err) {
			errs = append(errs, err)
		}
		query = modifiedQuery
	}

	return query, errs
}
