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
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
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
}

// NewRestTemplate returns a new instance of a restTemplate object.
// It implements the interfaces.RestTemplate interface.
func NewRestTemplate() interfaces.RestTemplate {
	return restTemplate{}
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
	latency := time.Since(startTime).String()

	err = r.treatHttpClientErr(err)
	if helper.IsNotNil(err) {
		r.traceHttpBackendResponseError(httpBackendRequest, span, latency, err)
		return nil, err
	}
	defer netHttpResponse.Body.Close()

	httpBackendResponse := vo.NewHttpBackendResponse(backend, netHttpResponse)
	r.traceHttpBackendResponse(httpBackendRequest, httpBackendResponse, span, latency)

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
	span := r.buildRequestSubSpan(ctx, httpBackendRequest)
	logger.Debugf("Backend HTTP request: %s --> %s", httpBackendRequest.Method(), httpBackendRequest.Url())
	return span
}

func (r restTemplate) traceHttpBackendResponseError(
	httpBackendRequest *vo.HttpBackendRequest,
	span opentracing.Span,
	latency string,
	err error,
) {
	reqMethod := httpBackendRequest.Method()
	reqUrl := httpBackendRequest.Url()

	errorDetails := errors.Details(err)
	errCause := errorDetails.GetCause()

	span.SetTag("response.error", errCause)
	logger.Errorf("Backend HTTP response: %s --> %s latency: %s err: %s", reqMethod, reqUrl, latency, errCause)
}

func (r restTemplate) traceHttpBackendResponse(
	httpBackendRequest *vo.HttpBackendRequest,
	httpBackendResponse *vo.HttpBackendResponse,
	span opentracing.Span,
	latency string,
) {
	reqMethod := httpBackendRequest.Method()
	reqUrl := httpBackendRequest.Url()
	resStatus := httpBackendResponse.Status()

	span.SetTag("response.status", resStatus)
	span.SetTag("response.header", httpBackendResponse.Header().String())
	if helper.IsNotNil(httpBackendResponse.Body()) {
		span.SetTag("response.body", httpBackendResponse.Body().CompactString())
	} else {
		span.SetTag("response.body", "")
	}

	logger.Debugf("Backend HTTP response: %s --> %s latency: %s status: %v", reqMethod, reqUrl, latency, resStatus)
}

func (r restTemplate) buildRequestSubSpan(ctx context.Context, httpBackendRequest *vo.HttpBackendRequest) opentracing.Span {
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
