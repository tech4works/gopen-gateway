/*
 * Copyright 2024 Gabriel Cataldo
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

package infra

import (
	"context"
	berrors "errors"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"net/url"
	"time"
)

// restTemplate is a struct that represents a template for making HTTP requests.
// It implements the interfaces.RestTemplate interface, which provides a method MakeRequest for sending
// an HTTP request and returning the corresponding HTTP response and an error, if any.
type restTemplate struct {
	logger interfaces.LoggerProvider
}

// NewRestTemplate returns a new instance of a restTemplate object.
// It implements the interfaces.RestTemplate interface.
func NewRestTemplate(logger interfaces.LoggerProvider) domain.RestTemplate {
	return restTemplate{
		logger: logger,
	}
}

func (r restTemplate) MakeRequest(ctx context.Context, backend *vo.Backend, httpBackendRequest *vo.HttpBackendRequest) (
	*vo.HttpBackendResponse, error) {
	netHttp, err := httpBackendRequest.NetHttp(ctx)
	if helper.IsNotNil(err) {
		return nil, err
	}

	span := r.traceHttpBackendRequest(ctx, httpBackendRequest)
	defer span.Finish()

	startTime := time.Now()
	netHttpResponse, err := http.DefaultClient.Do(netHttp)
	latency := time.Since(startTime)

	err = r.treatHttpClientErr(err)
	if helper.IsNotNil(err) {
		r.traceHttpBackendResponseError(backend, span, latency, err)
		return nil, err
	}
	defer netHttpResponse.Body.Close()

	httpBackendResponse := vo.NewHttpBackendResponse(backend, netHttpResponse, latency)
	r.traceHttpBackendResponse(backend, httpBackendResponse, span)

	return httpBackendResponse, nil
}

// treatHttpClientErr handles and transforms HTTP client errors.
// If the error is nil, it returns nil.
// If the error is an *url.Error and has a timeout, it returns a new domainmapper.ErrGatewayTimeout error.
// For any other error, it returns a new domainmapper.ErrBadGateway error.
//
// Inputs:
//   - err: The HTTP client error to be treated.
//
// Returns:
//   - error: The transformed error, if any.
//
// Example:
//
//	err := r.treatHttpClientErr(http.ErrTimeout)
//
// Note: For more details, see mapper.NewErrGatewayTimeoutByErr and mapper.NewErrBadGateway declarations.
func (r restTemplate) treatHttpClientErr(err error) error {
	if helper.IsNil(err) {
		return nil
	}

	var urlErr *url.Error
	berrors.As(err, &urlErr)
	if urlErr.Timeout() {
		err = mapper.NewErrGatewayTimeoutByErr(err)
	} else {
		err = mapper.NewErrBadGateway(err)
	}
	return err
}

func (r restTemplate) traceHttpBackendRequest(ctx context.Context, httpBackendRequest *vo.HttpBackendRequest,
) opentracing.Span {
	span := opentracing.SpanFromContext(ctx)
	if helper.IsNil(span) {
		return nil
	}

	urlTag := opentracing.Tag{Key: "request.url", Value: httpBackendRequest.Url()}
	methodTag := opentracing.Tag{Key: "request.method", Value: httpBackendRequest.Method()}
	headerTag := opentracing.Tag{Key: "request.header", Value: httpBackendRequest.Header().String()}
	bodyTag := opentracing.Tag{Key: "request.body", Value: ""}
	if helper.IsNotNil(httpBackendRequest.Body()) {
		bodyTag.Value = httpBackendRequest.Body().CompactString()
	}
	return span.Tracer().StartSpan(httpBackendRequest.Path().RawString(), opentracing.ChildOf(span.Context()),
		urlTag, methodTag, headerTag, bodyTag)
}

func (r restTemplate) traceHttpBackendResponseError(backend *vo.Backend, span opentracing.Span, latency time.Duration,
	err error) {
	errorDetails := errors.Details(err)
	span.SetTag("response.error", errorDetails.GetCause())

	format := "Error when trying to communicate with the backend service! latency: %s detail: %s"
	r.logger.PrintBackendErrorf(backend, format, latency.String(), errorDetails.GetMessage())
}

func (r restTemplate) traceHttpBackendResponse(backend *vo.Backend, httpBackendResponse *vo.HttpBackendResponse,
	span opentracing.Span) {
	span.SetTag("response.status", httpBackendResponse.Status())
	span.SetTag("response.header", httpBackendResponse.Header().String())
	if helper.IsNotNil(httpBackendResponse.Body()) {
		span.SetTag("response.body", httpBackendResponse.Body().CompactString())
	} else {
		span.SetTag("response.body", "")
	}

	r.logger.PrintBackendResponseInfo(backend, httpBackendResponse)
}
