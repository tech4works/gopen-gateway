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

package vo

import (
	"net/http"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type EndpointRequest struct {
	id       string
	traceID  string
	clientIP string

	protocol  enum.Protocol
	route     string
	path      URLPath
	query     Query
	operation string
	metadata  Metadata
	payload   *Payload
}

func NewHTTPEndpointRequest(
	id,
	traceID,
	clientIP,
	url string,
	path URLPath,
	query Query,
	method string,
	header Metadata,
	body *Payload,
) *EndpointRequest {
	return &EndpointRequest{
		id:        id,
		traceID:   traceID,
		clientIP:  clientIP,
		protocol:  enum.ProtocolHTTP,
		route:     url,
		path:      path,
		operation: method,
		metadata:  header,
		query:     query,
		payload:   body,
	}
}

func NewGRPCEndpointRequest(
	id,
	traceID,
	clientIP,
	fullMethod string,
	operation string,
	metadata Metadata,
	payload *Payload,
) *EndpointRequest {
	return &EndpointRequest{
		id:        id,
		traceID:   traceID,
		clientIP:  clientIP,
		protocol:  enum.ProtocolGRPC,
		route:     fullMethod,
		operation: operation,
		metadata:  metadata,
		payload:   payload,
	}
}

func (r *EndpointRequest) ID() string {
	return r.id
}

func (r *EndpointRequest) TraceID() string {
	return r.traceID
}

func (r *EndpointRequest) ClientIP() string {
	return r.clientIP
}

func (r *EndpointRequest) Protocol() enum.Protocol {
	return r.protocol
}

func (r *EndpointRequest) Route() string {
	return r.route
}

func (r *EndpointRequest) Path() URLPath {
	return r.path
}

func (r *EndpointRequest) Operation() string {
	return r.operation
}

func (r *EndpointRequest) Metadata() Metadata {
	return r.metadata
}

func (r *EndpointRequest) Params() Params {
	return r.Path().Params()
}

func (r *EndpointRequest) Query() Query {
	return r.query
}

func (r *EndpointRequest) HasPayload() bool {
	return checker.NonNil(r.payload)
}

func (r *EndpointRequest) Payload() *Payload {
	return r.payload
}

func (r *EndpointRequest) IsHTTP() bool {
	return checker.Equals(r.protocol, enum.ProtocolHTTP)
}

func (r *EndpointRequest) IsWebSocket() bool {
	return checker.Equals(r.protocol, enum.ProtocolWebSocket)
}

func (r *EndpointRequest) IsGRPC() bool {
	return checker.Equals(r.protocol, enum.ProtocolGRPC)
}

func (r *EndpointRequest) IsCORS() bool {
	return r.Metadata().Exists("Origin")
}

func (r *EndpointRequest) IsPreflight() bool {
	return checker.EqualsIgnoreCase(r.Operation(), http.MethodOptions) &&
		r.Metadata().Exists("Access-Control-Request-Method")
}

func (r *EndpointRequest) Map() (string, error) {
	var payload any
	if checker.NonNil(r.Payload()) {
		bodyMap, err := r.Payload().Map()
		if checker.NonNil(err) {
			return "", err
		}
		payload = bodyMap
	}
	switch r.protocol {
	case enum.ProtocolHTTP, enum.ProtocolWebSocket:
		return converter.ToStringWithErr(map[string]any{
			"id":        r.ID(),
			"client-ip": r.ClientIP(),
			"url":       r.Route(),
			"path":      r.Path().String(),
			"header":    r.Metadata().Map(),
			"method":    r.Operation(),
			"params":    r.Params().Map(),
			"query":     r.Query().Map(),
			"body":      payload,
		})
	case enum.ProtocolGRPC:
		return converter.ToStringWithErr(map[string]any{
			"id":          r.ID(),
			"client-ip":   r.ClientIP(),
			"full-method": r.Route(),
			"metadata":    r.Metadata().Map(),
			"operation":   r.Operation(),
			"payload":     payload,
		})
	default:
		return "", errors.New("unknown protocol to map request")
	}
}
