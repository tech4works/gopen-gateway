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

package http

import (
	"bytes"
	"context"
	goerrors "errors"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/opentracing/opentracing-go"
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"io"
	net "net/http"
	"net/url"
	"time"
)

type client struct {
}

func NewClient() app.HTTPClient {
	return client{}
}

func (c client) MakeRequest(ctx context.Context, endpoint *vo.Endpoint, request *vo.HTTPBackendRequest,
) *vo.HTTPBackendResponse {
	httpRequest, err := c.buildNetHTTPRequest(ctx, request)
	if checker.NonNil(err) {
		return c.buildHTTPBackendResponseByErr(endpoint, err)
	}

	span := c.startSpan(ctx, request)

	httpResponse, err := net.DefaultClient.Do(httpRequest)

	var httpBackendResponse *vo.HTTPBackendResponse
	if err = c.treatHTTPClientErr(err); checker.NonNil(err) {
		httpBackendResponse = c.buildHTTPBackendResponseByErr(endpoint, err)
	} else {
		httpBackendResponse, err = c.buildHTTPBackendResponse(httpResponse)
		if checker.NonNil(err) {
			httpBackendResponse = c.buildHTTPBackendResponseByErr(endpoint, err)
		}
	}

	c.finishSpan(span, httpBackendResponse)
	return httpBackendResponse
}

func (c client) buildNetHTTPRequest(ctx context.Context, request *vo.HTTPBackendRequest) (*net.Request, error) {
	var body io.ReadCloser
	if request.HasBody() {
		body = io.NopCloser(request.Body().Buffer())
	}
	netReq, err := net.NewRequestWithContext(ctx, request.Method(), request.Url(), body)
	if checker.NonNil(err) {
		return nil, err
	}

	header := *request.Header()
	query := request.Query()

	netReq.Header = header.Http()
	netReq.URL.RawQuery = query.Encode()

	span := opentracing.SpanFromContext(ctx)
	if checker.NonNil(span) {
		err = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(netReq.Header))
		if checker.NonNil(err) {
			return nil, err
		}
	}

	return netReq, nil
}

func (c client) startSpan(ctx context.Context, request *vo.HTTPBackendRequest) opentracing.Span {
	span := opentracing.SpanFromContext(ctx)
	if checker.IsNil(span) {
		return nil
	}

	urlTag := opentracing.Tag{Key: "request.url", Value: request.Url()}
	methodTag := opentracing.Tag{Key: "request.method", Value: request.Method()}
	headerTag := opentracing.Tag{Key: "request.header", Value: request.Header().String()}
	bodyTag := opentracing.Tag{Key: "request.body", Value: ""}
	if request.HasBody() {
		bodyTag.Value = request.Body().Resume()
	}
	childOf := opentracing.ChildOf(span.Context())

	return span.Tracer().StartSpan(request.Path().Raw(), childOf, urlTag, methodTag, headerTag, bodyTag)
}

func (c client) finishSpan(span opentracing.Span, httpBackendResponse *vo.HTTPBackendResponse) {
	if checker.IsNil(span) {
		return
	}

	statusCode := httpBackendResponse.StatusCode()
	header := httpBackendResponse.Header()

	span.SetTag("response.status", statusCode.String())
	span.SetTag("response.header", header.String())
	if checker.NonNil(httpBackendResponse.Body()) {
		span.SetTag("response.body", httpBackendResponse.Body().Resume())
	} else {
		span.SetTag("response.body", "")
	}

	span.Finish()
}

func (c client) treatHTTPClientErr(err error) error {
	if checker.IsNil(err) {
		return nil
	}

	var urlErr *url.Error
	goerrors.As(err, &urlErr)
	if urlErr.Timeout() {
		return mapper.NewErrGatewayTimeoutByErr(err)
	}

	return mapper.NewErrBadGateway(err)
}

func (c client) buildHTTPBackendResponseByErr(endpoint *vo.Endpoint, err error) *vo.HTTPBackendResponse {
	code := net.StatusInternalServerError
	if errors.Is(err, mapper.ErrGatewayTimeout) {
		code = net.StatusGatewayTimeout
	} else if errors.Is(err, mapper.ErrBadGateway) {
		code = net.StatusBadGateway
	}
	statusCode := vo.NewStatusCode(code)

	details := errors.Details(err)
	buffer := converter.ToBuffer(dto.ErrorBody{
		File:      details.GetFile(),
		Line:      details.GetLine(),
		Endpoint:  endpoint.Path(),
		Message:   details.GetMessage(),
		Timestamp: time.Now(),
	})
	body := vo.NewBodyJson(buffer)
	header := vo.NewHeaderByBody(body)

	return vo.NewHTTPBackendResponse(statusCode, header, body)
}

func (c client) buildHTTPBackendResponse(httpResponse *net.Response) (*vo.HTTPBackendResponse, error) {
	statusCode := vo.NewStatusCode(httpResponse.StatusCode)
	header := vo.NewHeader(httpResponse.Header)

	var body *vo.Body
	if checker.NonNil(httpResponse.Body) {
		contentType := httpResponse.Header.Get(mapper.ContentType)
		contentEncoding := httpResponse.Header.Get(mapper.ContentEncoding)

		bodyBytes, err := io.ReadAll(httpResponse.Body)
		if checker.NonNil(err) {
			return nil, err
		}
		body = vo.NewBody(contentType, contentEncoding, bytes.NewBuffer(bodyBytes))
	}

	return vo.NewHTTPBackendResponse(statusCode, header, body), nil
}
