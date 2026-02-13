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
	"net/http"
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type httpResponseFactory struct {
	aggregatorService   service.Aggregator
	omitterService      service.Omitter
	mapperService       service.Mapper
	projectorService    service.Projector
	nomenclatureService service.Nomenclature
	contentService      service.Content
	httpBackendFactory  Backend
}

type HTTPResponse interface {
	BuildErrorResponse(endpoint *vo.Endpoint, err error) *vo.HTTPResponse
	BuildAbortedResponse(history *vo.History) *vo.HTTPResponse
	BuildResponse(endpoint *vo.Endpoint, request *vo.HTTPRequest, history *vo.History) (*vo.HTTPResponse, []error)
}

func NewHTTPResponse(aggregatorService service.Aggregator, omitterService service.Omitter,
	mapperService service.Mapper, projectorService service.Projector, nomenclatureService service.Nomenclature,
	contentService service.Content, httpBackendFactory Backend) HTTPResponse {
	return httpResponseFactory{
		aggregatorService:   aggregatorService,
		omitterService:      omitterService,
		mapperService:       mapperService,
		projectorService:    projectorService,
		nomenclatureService: nomenclatureService,
		contentService:      contentService,
		httpBackendFactory:  httpBackendFactory,
	}
}

func (h httpResponseFactory) BuildErrorResponse(endpoint *vo.Endpoint, err error) *vo.HTTPResponse {
	statusCode := vo.NewStatusCode(http.StatusInternalServerError)
	header := vo.NewHeader(map[string][]string{
		mapper.XGopenCache:    {"false"},
		mapper.XGopenSuccess:  {converter.ToString(false)},
		mapper.XGopenComplete: {converter.ToString(false)},
	})
	details := errors.Wrap(err)
	buffer := converter.ToBuffer(dto.ErrorBody{
		File:      details.File(),
		Line:      details.Line(),
		Endpoint:  endpoint.Path(),
		Message:   details.Message(),
		Timestamp: time.Now(),
	})
	body := vo.NewBodyJson(buffer)

	return vo.NewHTTPResponse(statusCode, header, body)
}

func (h httpResponseFactory) BuildAbortedResponse(history *vo.History) *vo.HTTPResponse {
	lastestBackendResponse := history.BackendResponseLastest()
	lastestStatusCode := lastestBackendResponse.StatusCode()
	lastestHeader := lastestBackendResponse.Header()
	lastestBody := lastestBackendResponse.Body()

	header := vo.NewHeader(map[string][]string{
		mapper.XGopenCache:    {"false"},
		mapper.XGopenSuccess:  {converter.ToString(lastestBackendResponse.OK())},
		mapper.XGopenComplete: {converter.ToString(history.AllBackendsExecuted())},
	})
	header = h.aggregatorService.AggregateHeaders(header, lastestHeader)

	return vo.NewHTTPResponse(lastestStatusCode, header, lastestBody)
}

func (h httpResponseFactory) BuildResponse(endpoint *vo.Endpoint, request *vo.HTTPRequest, history *vo.History,
) (*vo.HTTPResponse, []error) {
	var allErrs []error

	statusCode := h.buildStatusCodeByHistory(history)
	body, bodyErrs := h.buildBodyByHistory(endpoint, request, history)
	header := h.buildHeaderByHistory(body, history)

	allErrs = append(allErrs, bodyErrs...)

	return vo.NewHTTPResponse(statusCode, header, body), allErrs
}

func (h httpResponseFactory) buildStatusCodeByHistory(history *vo.History) vo.StatusCode {
	if history.IsMultipleResponses() {
		return h.buildStatusCodeFromMultipleResponses(history)
	} else if history.IsSingleResponse() {
		return history.BackendResponseLastest().StatusCode()
	}
	return vo.NewStatusCode(http.StatusNoContent)
}

func (h httpResponseFactory) buildBodyByHistory(endpoint *vo.Endpoint, request *vo.HTTPRequest, history *vo.History) (
	*vo.Body, []error) {
	var body *vo.Body
	var errs []error

	if history.IsMultipleResponses() {
		body, errs = h.buildBodyFromMultipleResponses(endpoint, history)
	} else if history.IsSingleResponse() {
		body = history.BackendResponseLastest().Body()
	}

	if !endpoint.HasResponse() || !endpoint.Response().HasBody() {
		return body, nil
	}

	body, omitErrs := h.omitEmptyValuesFromBody(endpoint.Response().Body().OmitEmpty(), body)
	body, modifyCaseErrs := h.nomenclatureService.ToCase(endpoint.Response().Body().Nomenclature(), body)
	body, mapErrs := h.mapperService.MapBody(endpoint.Response().Body().Mapper(), body, request, history)
	body, projectErrs := h.projectorService.ProjectBody(endpoint.Response().Body().Projector(), body, request, history)
	body, modifyContentTypeErr := h.contentService.ModifyBodyContentType(endpoint.Response().Body().ContentType(), body)
	body, modifyBodyContentEncodingErr := h.contentService.ModifyBodyContentEncoding(endpoint.Response().Body().ContentEncoding(), body)

	errs = append(errs, omitErrs...)
	errs = append(errs, modifyCaseErrs...)
	errs = append(errs, mapErrs...)
	errs = append(errs, projectErrs...)
	if checker.NonNil(modifyContentTypeErr) {
		errs = append(errs, modifyContentTypeErr)
	}
	if checker.NonNil(modifyBodyContentEncodingErr) {
		errs = append(errs, modifyBodyContentEncodingErr)
	}

	return body, errs
}

func (h httpResponseFactory) buildHeaderByHistory(body *vo.Body, history *vo.History) vo.Header {
	mapHeader := map[string][]string{
		mapper.XGopenCache:    {"false"},
		mapper.XGopenSuccess:  {converter.ToString(history.AllOK())},
		mapper.XGopenComplete: {converter.ToString(history.AllBackendsExecuted())},
	}
	if checker.NonNil(body) {
		mapHeader[mapper.ContentType] = []string{body.ContentType().String()}
		mapHeader[mapper.ContentLength] = []string{body.SizeInString()}
		if body.HasContentEncoding() {
			mapHeader[mapper.ContentEncoding] = []string{body.ContentEncoding().String()}
		}
	}

	header := vo.NewHeader(mapHeader)

	for i := 0; i < history.Size(); i++ {
		backendResponse := history.GetBackendResponse(i)
		if checker.IsNil(backendResponse) {
			continue
		}
		header = h.aggregatorService.AggregateHeaders(header, backendResponse.Header())
	}

	return header
}

func (h httpResponseFactory) buildBodyFromMultipleResponses(endpoint *vo.Endpoint, history *vo.History) (*vo.Body, []error) {
	if endpoint.HasResponse() && endpoint.Response().HasBody() && endpoint.Response().Body().Aggregate() {
		return h.aggregatorService.AggregateBodies(history)
	}
	return h.aggregatorService.AggregateBodiesIntoSlice(history)
}

func (h httpResponseFactory) buildStatusCodeFromMultipleResponses(history *vo.History) vo.StatusCode {
	statusCodes := make(map[vo.StatusCode]int)
	for i := 0; i < history.Size(); i++ {
		backendResponse := history.GetBackendResponse(i)
		if checker.IsNil(backendResponse) {
			continue
		}
		statusCodes[backendResponse.StatusCode()]++
	}

	mostFrequentCode := vo.NewStatusCode(http.StatusNoContent)
	maxCount := 0
	for statusCode, count := range statusCodes {
		if checker.IsGreaterThanOrEqual(count, maxCount) {
			mostFrequentCode = statusCode
			maxCount = count
		}
	}

	return mostFrequentCode
}

func (h httpResponseFactory) omitEmptyValuesFromBody(omitEmpty bool, body *vo.Body) (*vo.Body, []error) {
	if !omitEmpty {
		return body, nil
	}
	return h.omitterService.OmitEmptyValuesFromBody(body)
}
