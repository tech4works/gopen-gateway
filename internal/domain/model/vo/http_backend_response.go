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

package vo

import (
	"bytes"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"io"
	"net/http"
	"time"
)

// HttpBackendResponse represents the structure of a backend HTTP response.
type HttpBackendResponse struct {
	// written is a boolean field that represents whether the response has been written or not.
	written bool
	// config is a pointer to the BackendResponse struct that holds configuration settings for the
	// HttpBackendResponse instance.
	config  *BackendResponse
	latency time.Duration
	// statusCode represents HTTP statusCode of a backend httpResponse.
	statusCode StatusCode
	// header represents the body fields of a backend httpResponse.
	header Header
	// body represents the body of a backend httpResponse.
	body *Body
}

// NewHttpBackendResponse creates a new HttpBackendResponse object based on the provided parameters.
// It constructs the header from the netHttpResponse, parses the response bytes into a body interface,
// and builds the backend httpResponse value object.
// It then calls the ApplyConfig method with the enum.BackendResponseApplyEarly flag and returns the result.
func NewHttpBackendResponse(backend *Backend, netHttpResponse *http.Response, latency time.Duration) *HttpBackendResponse {
	contentType := netHttpResponse.Header.Get("Content-Type")
	contentEncoding := netHttpResponse.Header.Get("Content-Encoding")

	header := NewHeader(netHttpResponse.Header)
	bodyBytes, _ := io.ReadAll(netHttpResponse.Body)
	body := NewBody(contentType, contentEncoding, bytes.NewBuffer(bodyBytes))

	return &HttpBackendResponse{
		config:     backend.Response(),
		latency:    latency,
		statusCode: NewStatusCode(netHttpResponse.StatusCode),
		header:     header,
		body:       body,
	}
}

// Written returns a boolean indicating whether the response has been written or not.
func (h *HttpBackendResponse) Written() bool {
	return h.written
}

// Ok returns a boolean value indicating whether the HttpBackendResponse's StatusCode is within the valid range of 200 to 299.
func (h *HttpBackendResponse) Ok() bool {
	return h.StatusCode().Ok()
}

// Key returns a string representing the key for the HttpBackendResponse instance based on the given index.
// If a Config exists, and it has a group, the key is set to the group value.
// If no group is configured, the key is generated as "backend-{index}".
func (h *HttpBackendResponse) Key(index int) (key string) {
	key = fmt.Sprintf("backend-%v", index)
	if helper.IsNotNil(h.Config()) && h.Config().HasGroup() {
		key = h.Config().Group()
	}
	return key
}

// StatusCode returns the statusCode of the HttpBackendResponse instance.
func (h *HttpBackendResponse) StatusCode() StatusCode {
	return h.statusCode
}

func (h *HttpBackendResponse) Status() string {
	code := h.StatusCode()
	return fmt.Sprintf("%v (%s)", code, http.StatusText(code.AsInt()))
}

// Header returns the header of the HttpBackendResponse instance.
func (h *HttpBackendResponse) Header() Header {
	return h.header
}

// Body returns the body of the HttpBackendResponse instance.
func (h *HttpBackendResponse) Body() *Body {
	return h.body
}

// Config returns a pointer to the BackendResponse struct that holds configuration settings for the
// HttpBackendResponse instance.
func (h *HttpBackendResponse) Config() *BackendResponse {
	return h.config
}

func (h *HttpBackendResponse) Latency() time.Duration {
	return h.latency
}

// GroupByType returns a boolean indicating if the body of the HttpBackendResponse instance
// is either a text or a slice of bytes. It checks whether the body is not nil and it's either
// a text or a slice of bytes.
func (h *HttpBackendResponse) GroupByType() bool {
	body := h.Body()
	return helper.IsNotNil(body) && body.ContentType().IsText() || helper.IsSlice(body.Bytes())
}

// Map returns a map[string]any containing the statusCode, header, and body of the HttpBackendResponse instance.
// The body is only included if it is not nil, and it is converted to its underlying interface value.
func (h *HttpBackendResponse) Map() any {
	var body any
	if helper.IsNotNil(h.Body()) {
		body = h.Body().Interface()
	}
	return map[string]any{
		"statusCode": h.StatusCode(),
		"header":     h.Header(),
		"body":       body,
	}
}

// ApplyConfig applies the given configuration settings to the HttpBackendResponse instance.
// It checks if the response has already been written or if the momentToApply value is not equal
// to the moment to apply specified in the backendResponse configuration. If any of these conditions
// are true, it returns the original HttpBackendResponse instance. If the backendResponse is configured to be omitted,
// it returns nil. Otherwise, it builds the body based on the backendResponse configuration and the given httpRequest
// and httpResponse. It creates a new HttpBackendResponse instance with the updated configuration settings and returns it.
func (h *HttpBackendResponse) ApplyConfig(momentToApply enum.BackendResponseApply, httpRequest *HttpRequest,
	httpResponse *HttpResponse) *HttpBackendResponse {
	backendResponse := h.Config()
	if h.Written() || helper.IsNil(backendResponse) || helper.IsNotEqualTo(momentToApply, backendResponse.Apply()) {
		return h
	} else if backendResponse.Omit() {
		return nil
	}

	body := h.buildBodyByConfig(backendResponse, httpRequest, httpResponse)
	return &HttpBackendResponse{
		latency:    h.latency,
		statusCode: h.statusCode,
		header:     h.buildHeaderByConfig(backendResponse, body, httpRequest, httpResponse),
		body:       body,
		written:    true,
	}
}

func (h *HttpBackendResponse) buildHeaderByConfig(backendResponse *BackendResponse, body *Body, httpRequest *HttpRequest,
	httpResponse *HttpResponse) Header {
	header := h.Header()

	if backendResponse.OmitHeader() {
		header = header.OnlyMandatoryKeys()
	} else {
		header = header.Map(backendResponse.HeaderMapper())
		header = header.Projection(backendResponse.HeaderProjection())
		for _, modifier := range backendResponse.HeaderModifiers() {
			header = header.Modify(&modifier, httpRequest, httpResponse)
		}
	}

	return header.Write(body)
}

// buildBodyByConfig builds the body for the HttpBackendResponse instance based on the configuration settings
// specified in the BackendResponse. If the BackendResponse is configured to omit the body or the existing body is nil,
// it returns nil. Otherwise, it applies the body mappings, projections, modifiers, and grouping specified in the BackendResponse
// to the existing body of the HttpBackendResponse. Finally, it returns the modified body.
func (h *HttpBackendResponse) buildBodyByConfig(backendResponse *BackendResponse, httpRequest *HttpRequest,
	httpResponse *HttpResponse) *Body {
	if backendResponse.OmitBody() || helper.IsNil(h.Body()) {
		return nil
	}

	body := h.Body()
	body = body.Map(backendResponse.BodyMapper())
	body = body.Projection(backendResponse.BodyProjection())
	for _, modifier := range backendResponse.BodyModifiers() {
		body = body.Modify(&modifier, httpRequest, httpResponse)
	}
	if backendResponse.HasGroup() {
		body = NewBodyAggregateByKey(backendResponse.Group(), body)
	}

	return body
}
