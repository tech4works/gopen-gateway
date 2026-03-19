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

package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"go.elastic.co/apm/module/apmhttp/v2"
)

type client struct {
	engine *http.Client
}

func NewClient() app.HTTPClient {
	return client{
		engine: apmhttp.WrapClient(&http.Client{}),
	}
}

func (c client) MakeRequest(ctx context.Context, parent *vo.EndpointRequest, request *vo.HTTPBackendRequest) (
	*http.Response, error) {
	httpRequest, err := c.buildNetHTTPRequest(ctx, parent, request)
	if checker.NonNil(err) {
		return nil, err
	}
	return c.engine.Do(httpRequest)
}

func (c client) buildNetHTTPRequest(ctx context.Context, parent *vo.EndpointRequest, request *vo.HTTPBackendRequest) (
	*http.Request, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, request.Method(), request.URL(), c.buildNetHTTPRequestBody(request))
	if checker.NonNil(err) {
		return nil, err
	}

	httpRequest.Header = c.buildNetHTTPRequestHeader(ctx, parent, request)
	httpRequest.URL.RawQuery = request.Query().Encode()

	return httpRequest, nil
}

func (c client) buildNetHTTPRequestHeader(ctx context.Context, parent *vo.EndpointRequest, request *vo.HTTPBackendRequest,
) http.Header {
	httpHeader := http.Header(request.Header().Copy())

	httpHeader.Set(app.XForwardedFor, parent.ClientIP())
	httpHeader.Set(app.XGopenRequestID, parent.ID())
	httpHeader.Set(app.XGopenDegraded, converter.ToString(request.Degraded()))
	httpHeader.Set(app.XGopenURLPathDegraded, converter.ToString(request.URLPathDegraded()))
	httpHeader.Set(app.XGopenHeaderDegraded, converter.ToString(request.HeaderDegraded()))
	httpHeader.Set(app.XGopenQueryDegraded, converter.ToString(request.QueryDegraded()))
	httpHeader.Set(app.XGopenBodyDegraded, converter.ToString(request.BodyDegraded()))

	if request.HasBody() {
		httpHeader.Set(app.ContentType, request.Body().ContentType().String())
		httpHeader.Set(app.ContentLength, request.Body().SizeInString())

		if request.Body().HasContentEncoding() {
			httpHeader.Set(app.ContentEncoding, request.Body().ContentEncoding().String())
		}
	}

	timeout, ok := ctx.Deadline()
	if ok {
		remaining := time.Until(timeout)
		httpHeader.Set(app.XGopenTimeout, converter.ToString(remaining.Milliseconds()))
	}

	return httpHeader
}

func (c client) buildNetHTTPRequestBody(request *vo.HTTPBackendRequest) io.Reader {
	var body io.ReadCloser
	if request.HasBody() {
		body = io.NopCloser(request.Body().Buffer())
	}
	return body
}
