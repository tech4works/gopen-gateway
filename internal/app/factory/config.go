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
	"fmt"
	"strings"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type propagateState struct {
	continueOnError *bool
	header          *dto.BackendRequestHeader
	param           *dto.BackendRequestParam
	query           *dto.BackendRequestQuery
	body            *dto.BackendRequestBody
}

func BuildGopen(gopen *dto.Gopen) *vo.Gopen {
	return vo.NewGopen(
		buildProxy(gopen.Proxy),
		buildSecurityCors(gopen.SecurityCors),
		buildEndpoints(gopen),
	)
}

func buildProxy(proxy *dto.Proxy) *vo.Proxy {
	if checker.IsNil(proxy) {
		return nil
	}
	return vo.NewProxy(proxy.Provider, proxy.Token, proxy.Domains)
}

func buildSecurityCors(securityCors *dto.SecurityCors) *vo.SecurityCors {
	if checker.IsNil(securityCors) {
		return nil
	}
	return vo.NewSecurityCors(securityCors.AllowOrigins, securityCors.AllowMethods, securityCors.AllowHeaders)
}

func buildEndpoints(gopen *dto.Gopen) []vo.Endpoint {
	var endpoints []vo.Endpoint
	var errs []string

	for _, endpoint := range gopen.Endpoints {
		var err string
		for _, registeredEndpoint := range endpoints {
			if checker.NotEquals(endpoint.Path, registeredEndpoint.Path()) ||
				checker.NotEquals(endpoint.Method, registeredEndpoint.Method()) {
				continue
			}
			err = fmt.Sprintf("- Duplicate endpoint path: %s method: %s", endpoint.Path, endpoint.Method)
			errs = append(errs, err)
		}
		if checker.IsEmpty(err) {
			endpoints = append(endpoints, buildEndpoint(gopen, endpoint))
		}
	}

	if checker.IsNotEmpty(errs) {
		panic(strings.Join(errs, "\n"))
	}

	return endpoints
}

func buildEndpoint(gopen *dto.Gopen, endpoint dto.Endpoint) vo.Endpoint {
	return vo.NewEndpoint(
		endpoint.Path,
		endpoint.Method,
		buildTimeout(gopen.Timeout, endpoint.Timeout),
		buildLimiter(gopen.Limiter, endpoint.Limiter),
		buildCache(gopen.Cache, endpoint.Cache),
		endpoint.AbortIfHTPPStatusCodes,
		endpoint.Parallelism,
		buildEndpointResponse(endpoint.Response),
		buildBackends(gopen.Templates, endpoint),
	)
}

func buildTimeout(timeout, endpointTimeout vo.Duration) vo.Duration {
	if checker.IsGreaterThan(endpointTimeout, 0) {
		return endpointTimeout
	} else {
		return timeout
	}
}

func buildLimiter(limiter *dto.Limiter, endpointLimiter *dto.EndpointLimiter) vo.Limiter {
	var maxHeaderSize vo.Bytes
	var maxBodySize vo.Bytes
	var maxMultipartForm vo.Bytes
	var endpointRate, rate *dto.Rate

	if checker.NonNil(limiter) {
		if checker.NonNil(limiter.MaxHeaderSize) {
			maxHeaderSize = *limiter.MaxHeaderSize
		}
		if checker.NonNil(limiter.MaxBodySize) {
			maxBodySize = *limiter.MaxBodySize
		}
		if checker.NonNil(limiter.MaxMultipartMemorySize) {
			maxMultipartForm = *limiter.MaxMultipartMemorySize
		}
		rate = limiter.Rate
	}

	if checker.NonNil(endpointLimiter) {
		if checker.NonNil(endpointLimiter.MaxHeaderSize) {
			maxHeaderSize = endpointLimiter.MaxHeaderSize
		}
		if checker.NonNil(endpointLimiter.MaxBodySize) {
			maxBodySize = endpointLimiter.MaxBodySize
		}
		if checker.NonNil(endpointLimiter.MaxMultipartMemorySize) {
			maxMultipartForm = endpointLimiter.MaxMultipartMemorySize
		}
		endpointRate = endpointLimiter.Rate
	}

	return vo.NewLimiter(maxHeaderSize, maxBodySize, maxMultipartForm, buildLimiterRate(rate, endpointRate))
}

func buildLimiterRate(rate, endpointRate *dto.Rate) vo.Rate {
	var every vo.Duration
	var capacity int

	if checker.NonNil(rate) {
		if checker.NonNil(rate.Every) {
			every = *rate.Every
		}
		if checker.NonNil(rate.Capacity) {
			capacity = *rate.Capacity
		}
	}

	if checker.NonNil(endpointRate) {
		if checker.NonNil(endpointRate.Every) {
			every = *endpointRate.Every
		}
		if checker.NonNil(endpointRate.Capacity) {
			capacity = *endpointRate.Capacity
		}
	}

	return vo.NewRate(every, capacity)
}

func buildCache(cache *dto.Cache, endpointCache *dto.EndpointCache) *vo.Cache {
	if checker.IsNil(cache) && checker.IsNil(endpointCache) {
		return nil
	}

	var enabled bool
	var ignoreQuery bool
	var duration vo.Duration
	var strategyHeaders []string
	var onlyIfStatusCodes []int
	var onlyIfMethods []string
	var allowCacheControl *bool

	if checker.NonNil(cache) {
		duration = cache.Duration
		strategyHeaders = cache.StrategyHeaders
		onlyIfStatusCodes = cache.OnlyIfStatusCodes
		onlyIfMethods = cache.OnlyIfMethods
		allowCacheControl = cache.AllowCacheControl
	}

	if checker.NonNil(endpointCache) {
		enabled = endpointCache.Enabled
		ignoreQuery = endpointCache.IgnoreQuery
		if checker.IsGreaterThan(endpointCache.Duration, 0) {
			duration = endpointCache.Duration
		}
		if checker.NonNil(endpointCache.StrategyHeaders) {
			strategyHeaders = endpointCache.StrategyHeaders
		}
		if checker.NonNil(endpointCache.AllowCacheControl) {
			allowCacheControl = endpointCache.AllowCacheControl
		}
		if checker.NonNil(endpointCache.OnlyIfStatusCodes) {
			onlyIfStatusCodes = endpointCache.OnlyIfStatusCodes
		}
	}

	return vo.NewCache(enabled, ignoreQuery, duration, strategyHeaders, onlyIfStatusCodes, onlyIfMethods, allowCacheControl)
}

func buildEndpointResponse(endpointResponse *dto.EndpointResponse) *vo.EndpointResponse {
	if checker.IsNil(endpointResponse) {
		return nil
	}
	return vo.NewEndpointResponse(
		endpointResponse.ContinueOnError,
		buildEndpointResponseHeader(endpointResponse.Header),
		buildEndpointResponseBody(endpointResponse.Body),
	)
}

func buildEndpointResponseHeader(endpointResponseHeader *dto.EndpointResponseHeader) *vo.EndpointResponseHeader {
	if checker.IsNil(endpointResponseHeader) {
		return nil
	}
	return vo.NewEndpointResponseHeader(
		buildMapper(endpointResponseHeader.Mapper),
		buildProjector(endpointResponseHeader.Projector),
	)
}

func buildEndpointResponseBody(endpointResponseBody *dto.EndpointResponseBody) *vo.EndpointResponseBody {
	if checker.IsNil(endpointResponseBody) {
		return nil
	}
	return vo.NewEndpointResponseBody(
		endpointResponseBody.Aggregate,
		endpointResponseBody.OmitEmpty,
		endpointResponseBody.ContentType,
		endpointResponseBody.ContentEncoding,
		endpointResponseBody.Nomenclature,
		buildMapper(endpointResponseBody.Mapper),
		buildProjector(endpointResponseBody.Projector),
	)
}

func buildBackends(templates *dto.Templates, endpoint dto.Endpoint) []vo.Backend {
	var result []vo.Backend
	var ps propagateState
	var backendIndex int

	seen := map[string]struct{}{}

	consume := func(flow enum.BackendFlow, items []dto.Backend, allowParallelism bool) {
		for _, b := range items {
			effective := resolveBackendTemplate(b, templates)

			validateBackendID(endpoint, flow, backendIndex, effective)
			validateBackendDependencies(endpoint, flow, backendIndex, effective, seen)
			markBackendSeen(endpoint, flow, backendIndex, effective, seen)

			result = append(result, buildBackendUnified(
				effective,
				flow,
				&ps,
				backendIndex,
				allowParallelism && endpoint.Parallelism,
			))
			backendIndex++
		}
	}

	consume(enum.BackendFlowBeforeware, endpoint.Beforewares, false)
	consume(enum.BackendFlowNormal, endpoint.Backends, true)
	consume(enum.BackendFlowAfterware, endpoint.Afterwares, false)

	return result
}

func validateBackendID(endpoint dto.Endpoint, flow enum.BackendFlow, backendIndex int, b dto.Backend) {
	if checker.IsEmpty(b.ID) {
		panic(errors.Newf(
			"backend has empty id after resolution (endpoint=%s %s, flow=%s, index=%d). "+
				"For templates, id must be template.path. For raw backends, default id is backend.path (must be non-empty).",
			endpoint.Method, endpoint.Path, flow, backendIndex,
		))
	}
}

func validateBackendDependencies(
	endpoint dto.Endpoint,
	flow enum.BackendFlow,
	backendIndex int,
	b dto.Backend,
	seen map[string]struct{},
) {
	if checker.IsEmpty(b.Dependencies) {
		return
	}

	var missing []string
	for _, dep := range b.Dependencies {
		dep = strings.TrimSpace(dep)
		if checker.IsEmpty(dep) {
			continue
		}
		if _, ok := seen[dep]; !ok {
			missing = append(missing, dep)
		}
	}

	if len(missing) > 0 {
		panic(errors.Newf(
			"dependencies not satisfied (endpoint=%s %s, flow=%s, index=%d, backend.id=%q). "+
				"Missing (must appear BEFORE in the same endpoint): %s",
			endpoint.Method,
			endpoint.Path,
			flow,
			backendIndex,
			b.ID,
			strings.Join(missing, ", "),
		))
	}
}

func markBackendSeen(endpoint dto.Endpoint, flow enum.BackendFlow, backendIndex int, b dto.Backend, seen map[string]struct{}) {
	if _, exists := seen[b.ID]; exists {
		panic(errors.Newf(
			"duplicate backend.id within the same endpoint (endpoint=%s %s, flow=%s, index=%d, backend.id=%q). "+
				"The id must be unique per endpoint for dependencies to work correctly.",
			endpoint.Method, endpoint.Path, flow, backendIndex, b.ID,
		))
	}
	seen[b.ID] = struct{}{}
}

func buildBackendUnified(
	backend dto.Backend,
	flow enum.BackendFlow,
	ps *propagateState,
	backendIndex int,
	parallelism bool,
) vo.Backend {
	resp := buildBackendResponse(backend, flow)

	switch backend.Kind {
	case enum.BackendKindPublisher:
		return vo.NewBackendPublisher(
			backend.ID,
			flow,
			backend.OnlyIf,
			backend.IgnoreIf,
			buildUnifiedPublisher(backend, parallelism),
			resp,
		)
	case enum.BackendKindHTTP:
		return vo.NewBackendHTTP(
			backend.ID,
			flow,
			backend.OnlyIf,
			backend.IgnoreIf,
			backend.Hosts,
			backend.Path,
			backend.Method,
			buildHTTPBackendRequest(backend, ps, backendIndex, parallelism),
			resp,
		)
	default:
		panic(errors.Newf("invalid backend.kind=%v (endpoint=%s %s)", backend.Kind, flow, backend.Path))
	}
}

func buildHTTPBackendRequest(
	backend dto.Backend,
	ps *propagateState,
	backendIndex int,
	parallelism bool,
) *vo.BackendRequest {
	effective := mergeBackendRequestWithPropagation(backend.Request, ps)

	applyBackendPropagateIntoState(backend.Propagate, ps, backendIndex)
	collectPropagatingModifiersFromRequestIntoState(backend.Request, ps, backendIndex)

	if checker.IsNil(effective) {
		return nil
	}

	async := effective.Async
	if checker.IsNil(async) && checker.NonNil(backend.Async) {
		async = backend.Async
	}

	return vo.NewBackendRequest(
		checker.IfNilReturns(effective.ContinueOnError, false),
		effective.Concurrent,
		resolveAsync(async, parallelism),
		buildBackendRequestHeader(effective.Header),
		buildBackendRequestParam(effective.Param),
		buildBackendRequestQuery(effective.Query),
		buildBackendRequestBody(effective.Body),
	)
}

func buildBackendResponse(backend dto.Backend, flow enum.BackendFlow) *vo.BackendResponse {
	if checker.Equals(flow, enum.BackendFlowBeforeware) || checker.Equals(flow, enum.BackendFlowAfterware) {
		// middleware: força o retorno de um BackendResponseForMiddleware
		if checker.IsNil(backend.Response) {
			return vo.NewBackendResponseForMiddleware()
		}
	}

	if checker.IsNil(backend.Response) {
		return nil
	}

	return vo.NewBackendResponse(
		backend.Response.ContinueOnError,
		backend.Response.Omit,
		buildBackendResponseHeader(backend.Response.Header),
		buildBackendResponseBody(backend.Response.Body),
	)
}

func buildUnifiedPublisher(backend dto.Backend, parallelism bool) vo.Publisher {
	return vo.NewPublisher(
		backend.OnlyIf,
		backend.IgnoreIf,
		backend.Provider,
		backend.Path,
		backend.GroupID,
		backend.DeduplicationID,
		backend.Delay,
		resolveAsync(backend.Async, parallelism),
		buildPublisherMessage(backend.Message),
	)
}

func buildPublisherMessage(publisherMessage *dto.PublisherMessage) *vo.PublisherMessage {
	if checker.IsNil(publisherMessage) {
		return nil
	}
	return vo.NewPublisherMessage(
		publisherMessage.ContinueOnError,
		publisherMessage.OnlyIf,
		publisherMessage.IgnoreIf,
		buildPublisherMessageAttributes(publisherMessage.Attributes),
		buildPublisherMessageBody(publisherMessage.Body),
	)
}

func buildPublisherMessageAttributes(publisherMessageAttributes map[string]dto.PublisherMessageAttribute,
) map[string]vo.PublisherMessageAttribute {
	if checker.IsNil(publisherMessageAttributes) {
		return nil
	}

	var result = make(map[string]vo.PublisherMessageAttribute)
	for key, value := range publisherMessageAttributes {
		result[key] = vo.NewPublisherMessageAttribute(value.DataType, value.Value)
	}
	return result
}

func buildPublisherMessageBody(publisherMessageBody *dto.PublisherMessageBody) *vo.PublisherMessageBody {
	if checker.IsNil(publisherMessageBody) {
		return nil
	}
	return vo.NewPublisherMessageBody(
		publisherMessageBody.OmitEmpty,
		buildMapper(publisherMessageBody.Mapper),
		buildProjector(publisherMessageBody.Projector),
		buildModifiers(publisherMessageBody.Modifiers),
	)
}

func buildBackendRequestHeader(backendRequestHeader *dto.BackendRequestHeader) *vo.BackendRequestHeader {
	if checker.IsNil(backendRequestHeader) {
		return nil
	}
	return vo.NewBackendRequestHeader(
		backendRequestHeader.Omit,
		buildMapper(backendRequestHeader.Mapper),
		buildProjector(backendRequestHeader.Projector),
		buildModifiers(backendRequestHeader.Modifiers),
	)
}

func buildBackendRequestParam(backendRequestParam *dto.BackendRequestParam) *vo.BackendRequestParam {
	if checker.IsNil(backendRequestParam) {
		return nil
	}
	return vo.NewBackendRequestParam(buildModifiers(backendRequestParam.Modifiers))
}

func buildBackendRequestQuery(backendRequestQuery *dto.BackendRequestQuery) *vo.BackendRequestQuery {
	if checker.IsNil(backendRequestQuery) {
		return nil
	}
	return vo.NewBackendRequestQuery(
		backendRequestQuery.Omit,
		buildMapper(backendRequestQuery.Mapper),
		buildProjector(backendRequestQuery.Projector),
		buildModifiers(backendRequestQuery.Modifiers),
	)
}

func buildBackendRequestBody(backendRequestBody *dto.BackendRequestBody) *vo.BackendRequestBody {
	if checker.IsNil(backendRequestBody) {
		return nil
	}
	return vo.NewBackendRequestBody(
		backendRequestBody.Omit,
		backendRequestBody.OmitEmpty,
		backendRequestBody.ContentType,
		backendRequestBody.ContentEncoding,
		backendRequestBody.Nomenclature,
		buildMapper(backendRequestBody.Mapper),
		buildProjector(backendRequestBody.Projector),
		buildModifiers(backendRequestBody.Modifiers),
	)
}

func buildBackendResponseHeader(backendResponseHeader *dto.BackendResponseHeader) *vo.BackendResponseHeader {
	if checker.IsNil(backendResponseHeader) {
		return nil
	}
	return vo.NewBackendResponseHeader(
		backendResponseHeader.Omit,
		buildMapper(backendResponseHeader.Mapper),
		buildProjector(backendResponseHeader.Projector),
		buildModifiers(backendResponseHeader.Modifiers),
	)
}

func buildBackendResponseBody(backendResponseBody *dto.BackendResponseBody) *vo.BackendResponseBody {
	if checker.IsNil(backendResponseBody) {
		return nil
	}
	return vo.NewBackendResponseBody(
		backendResponseBody.Omit,
		backendResponseBody.Group,
		buildMapper(backendResponseBody.Mapper),
		buildProjector(backendResponseBody.Projector),
		buildModifiers(backendResponseBody.Modifiers),
	)
}

func buildMapper(mapper *dto.Mapper) *vo.Mapper {
	if checker.IsNil(mapper) {
		return nil
	}
	return vo.NewMapper(mapper.OnlyIf, mapper.IgnoreIf, mapper.Map)
}

func buildProjector(projector *dto.Projector) *vo.Projector {
	if checker.IsNil(projector) {
		return nil
	}
	return vo.NewProjector(projector.OnlyIf, projector.IgnoreIf, projector.Project)
}

func buildModifiers(modifiers []dto.Modifier) []vo.Modifier {
	var result []vo.Modifier
	for _, modifier := range modifiers {
		result = append(result, buildModifier(modifier))
	}
	return result
}

func buildModifier(modifier dto.Modifier) vo.Modifier {
	return vo.NewModifier(modifier.OnlyIf, modifier.IgnoreIf, modifier.Action, modifier.Propagate, modifier.Key, modifier.Value)
}

func resolveBackendTemplate(cur dto.Backend, templates *dto.Templates) dto.Backend {
	if checker.IsNil(cur.Template) || checker.IsEmpty(cur.Template.Path) {
		if checker.IsEmpty(cur.ID) && checker.IsNotEmpty(cur.Path) {
			cur.ID = cur.Path
		}
		return cur
	}

	tpl := findTemplateBackend(cur.Template.Path, templates)
	mergeMode := normalizeMerge(cur.Template.Merge)

	base := tpl
	base.ID = cur.Template.Path
	if mergeMode == "BASE" {
		base = keepTemplateBase(tpl)
	}

	merged := mergeBackendFULL(base, cur)
	locked := forceLockedFieldsFromTemplate(tpl, merged)
	return locked
}

func findTemplateBackend(path string, templates *dto.Templates) dto.Backend {
	if checker.IsNil(templates) {
		panic(errors.Newf("templates not configured, it was not possible to resolve template.path=%q", path))
	}

	if checker.NonNil(templates.Beforewares) {
		if b, ok := templates.Beforewares[path]; ok {
			return b
		}
	}
	if checker.NonNil(templates.Backends) {
		if b, ok := templates.Backends[path]; ok {
			return b
		}
	}
	if checker.NonNil(templates.Afterwares) {
		if b, ok := templates.Afterwares[path]; ok {
			return b
		}
	}

	panic(errors.Newf("template %q not found in templates.beforewares/backends/publishers/afterwares", path))
}

func normalizeMerge(in enum.TemplateMerge) enum.TemplateMerge {
	if !in.IsEnumValid() {
		return enum.TemplateMergeFull
	}
	return in
}

func keepTemplateBase(tpl dto.Backend) dto.Backend {
	out := dto.Backend{
		Comment:      tpl.Comment,
		ID:           tpl.ID,
		Dependencies: tpl.Dependencies,
		OnlyIf:       tpl.OnlyIf,
		IgnoreIf:     tpl.IgnoreIf,
		Kind:         tpl.Kind,
	}

	switch tpl.Kind {
	case enum.BackendKindPublisher:
		out.Provider = tpl.Provider
		out.Path = tpl.Path
		return out
	case enum.BackendKindHTTP:
		out.Hosts = tpl.Hosts
		out.Path = tpl.Path
		out.Method = tpl.Method
		return out
	default:
		return out
	}
}

func forceLockedFieldsFromTemplate(tpl dto.Backend, merged dto.Backend) dto.Backend {
	// kind sempre travado
	merged.Kind = tpl.Kind

	switch tpl.Kind {
	case enum.BackendKindPublisher:
		merged.Provider = tpl.Provider
		merged.Path = tpl.Path
		return merged

	case enum.BackendKindHTTP:
		merged.Hosts = tpl.Hosts
		merged.Path = tpl.Path
		merged.Method = tpl.Method
		return merged

	default:
		return merged
	}
}

func mergeBackendFULL(tpl dto.Backend, cur dto.Backend) dto.Backend {
	out := tpl

	if checker.IsNotEmpty(cur.Comment) {
		out.Comment = cur.Comment
	}
	if checker.IsNotEmpty(cur.Dependencies) {
		out.Dependencies = append(append([]string{}, out.Dependencies...), cur.Dependencies...)
	}
	if checker.IsNotEmpty(cur.OnlyIf) {
		out.OnlyIf = append(append([]string{}, out.OnlyIf...), cur.OnlyIf...)
	}
	if checker.IsNotEmpty(cur.IgnoreIf) {
		out.IgnoreIf = append(append([]string{}, out.IgnoreIf...), cur.IgnoreIf...)
	}

	if cur.Kind.IsEnumValid() {
		out.Kind = cur.Kind
	}
	if checker.NonNil(cur.Async) {
		out.Async = cur.Async
	}

	// HTTP
	if checker.IsNotEmpty(cur.Hosts) {
		out.Hosts = cur.Hosts
	}
	if checker.IsNotEmpty(cur.Path) {
		out.Path = cur.Path
	}
	if checker.IsNotEmpty(cur.Method) {
		out.Method = cur.Method
	}

	out.Request = mergeBackendRequest(out.Request, cur.Request)
	out.Response = mergeBackendResponse(out.Response, cur.Response)
	out.Propagate = mergeBackendPropagate(out.Propagate, cur.Propagate)

	// PUBLISHER
	if cur.Provider.IsEnumValid() {
		out.Provider = cur.Provider
	}
	if checker.IsNotEmpty(cur.GroupID) {
		out.GroupID = cur.GroupID
	}
	if checker.IsNotEmpty(cur.DeduplicationID) {
		out.DeduplicationID = cur.DeduplicationID
	}
	if checker.IsGreaterThan(cur.Delay, 0) {
		out.Delay = cur.Delay
	}
	if checker.NonNil(cur.Message) {
		out.Message = mergePublisherMessage(out.Message, cur.Message)
	}

	return out
}

func mergePublisherMessage(tpl, cur *dto.PublisherMessage) *dto.PublisherMessage {
	if checker.IsNil(tpl) && checker.IsNil(cur) {
		return nil
	}
	var out dto.PublisherMessage
	if checker.NonNil(tpl) {
		out = *tpl
	}
	if checker.NonNil(cur) {
		out.ContinueOnError = cur.ContinueOnError
		if checker.IsNotEmpty(cur.OnlyIf) {
			out.OnlyIf = append(append([]string{}, out.OnlyIf...), cur.OnlyIf...)
		}
		if checker.IsNotEmpty(cur.IgnoreIf) {
			out.IgnoreIf = append(append([]string{}, out.IgnoreIf...), cur.IgnoreIf...)
		}
		if checker.NonNil(cur.Attributes) {
			if out.Attributes == nil {
				out.Attributes = map[string]dto.PublisherMessageAttribute{}
			}
			for k, v := range cur.Attributes {
				out.Attributes[k] = v
			}
		}
		out.Body = mergePublisherMessageBody(out.Body, cur.Body)
	}
	return &out
}

func mergePublisherMessageBody(tpl, cur *dto.PublisherMessageBody) *dto.PublisherMessageBody {
	if checker.IsNil(tpl) && checker.IsNil(cur) {
		return nil
	}
	var out dto.PublisherMessageBody
	if checker.NonNil(tpl) {
		out = *tpl
	}
	if checker.NonNil(cur) {
		out.OmitEmpty = out.OmitEmpty || cur.OmitEmpty
		if checker.NonNil(cur.Mapper) {
			out.Mapper = cur.Mapper
		}
		if checker.NonNil(cur.Projector) {
			out.Projector = cur.Projector
		}
		out.Modifiers = append(append([]dto.Modifier{}, out.Modifiers...), cur.Modifiers...)
	}
	return &out
}

func mergeBackendRequest(tpl, cur *dto.BackendRequest) *dto.BackendRequest {
	if checker.IsNil(tpl) && checker.IsNil(cur) {
		return nil
	}

	var out dto.BackendRequest
	if checker.NonNil(tpl) {
		out = *tpl
	}

	if checker.IsNil(cur) {
		return &out
	}

	if checker.NonNil(cur.ContinueOnError) {
		out.ContinueOnError = cur.ContinueOnError
	}
	if checker.IsGreaterThan(cur.Concurrent, 0) {
		out.Concurrent = cur.Concurrent
	}
	if checker.NonNil(cur.Async) {
		out.Async = cur.Async
	}

	out.Header = mergeBackendRequestHeader(out.Header, cur.Header)
	out.Param = mergeBackendRequestParam(out.Param, cur.Param)
	out.Query = mergeBackendRequestQuery(out.Query, cur.Query)
	out.Body = mergeBackendRequestBody(out.Body, cur.Body)

	return &out
}

func mergeBackendRequestWithPropagation(req *dto.BackendRequest, ps *propagateState) *dto.BackendRequest {
	if checker.IsNil(req) &&
		checker.IsNil(ps.header) &&
		checker.IsNil(ps.param) &&
		checker.IsNil(ps.query) &&
		checker.IsNil(ps.body) {
		return nil
	}

	var out dto.BackendRequest
	if checker.NonNil(req) {
		out = *req
	}

	// ContinueOnError (*bool):
	// - se o backend atual NÃO configurou (nil), herda o valor propagado (pode ser nil também)
	// - se o backend atual configurou false, ele bloqueia a herança (mantém false)
	// - se o backend atual configurou true, ele ganha (mantém true)
	if checker.IsNil(out.ContinueOnError) {
		out.ContinueOnError = ps.continueOnError
	}

	out.Header = mergeBackendRequestHeader(ps.header, out.Header)
	out.Param = mergeBackendRequestParam(ps.param, out.Param)
	out.Query = mergeBackendRequestQuery(ps.query, out.Query)
	out.Body = mergeBackendRequestBody(ps.body, out.Body)

	return &out
}

func mergeBackendRequestHeader(inh, cur *dto.BackendRequestHeader) *dto.BackendRequestHeader {
	if checker.IsNil(inh) && checker.IsNil(cur) {
		return nil
	}
	var out dto.BackendRequestHeader
	if checker.NonNil(inh) {
		out = *inh
	}
	if checker.NonNil(cur) {
		// bool: herda por OR (se algum upstream pediu omit, mantém omit)
		out.Omit = out.Omit || cur.Omit

		// ponteiros: o "cur" tem precedência quando vier preenchido
		if checker.NonNil(cur.Mapper) {
			out.Mapper = cur.Mapper
		}
		if checker.NonNil(cur.Projector) {
			out.Projector = cur.Projector
		}

		// modifiers: inherited primeiro, depois current
		out.Modifiers = append(append([]dto.Modifier{}, out.Modifiers...), cur.Modifiers...)
	}
	return &out
}

func mergeBackendRequestParam(inh, cur *dto.BackendRequestParam) *dto.BackendRequestParam {
	if checker.IsNil(inh) && checker.IsNil(cur) {
		return nil
	}
	var out dto.BackendRequestParam
	if checker.NonNil(inh) {
		out = *inh
	}
	if checker.NonNil(cur) {
		out.Modifiers = append(append([]dto.Modifier{}, out.Modifiers...), cur.Modifiers...)
	}
	return &out
}

func mergeBackendRequestQuery(inh, cur *dto.BackendRequestQuery) *dto.BackendRequestQuery {
	if checker.IsNil(inh) && checker.IsNil(cur) {
		return nil
	}
	var out dto.BackendRequestQuery
	if checker.NonNil(inh) {
		out = *inh
	}
	if checker.NonNil(cur) {
		out.Omit = out.Omit || cur.Omit

		if checker.NonNil(cur.Mapper) {
			out.Mapper = cur.Mapper
		}
		if checker.NonNil(cur.Projector) {
			out.Projector = cur.Projector
		}

		out.Modifiers = append(append([]dto.Modifier{}, out.Modifiers...), cur.Modifiers...)
	}
	return &out
}

func mergeBackendRequestBody(inh, cur *dto.BackendRequestBody) *dto.BackendRequestBody {
	if checker.IsNil(inh) && checker.IsNil(cur) {
		return nil
	}
	var out dto.BackendRequestBody
	if checker.NonNil(inh) {
		out = *inh
	}
	if checker.NonNil(cur) {
		out.Omit = out.Omit || cur.Omit
		out.OmitEmpty = out.OmitEmpty || cur.OmitEmpty

		// enums: "cur" ganha quando for válido, senão herda
		if cur.ContentType.IsEnumValid() {
			out.ContentType = cur.ContentType
		}
		if cur.ContentEncoding.IsEnumValid() {
			out.ContentEncoding = cur.ContentEncoding
		}
		if cur.Nomenclature.IsEnumValid() {
			out.Nomenclature = cur.Nomenclature
		}

		if checker.NonNil(cur.Mapper) {
			out.Mapper = cur.Mapper
		}
		if checker.NonNil(cur.Projector) {
			out.Projector = cur.Projector
		}

		out.Modifiers = append(append([]dto.Modifier{}, out.Modifiers...), cur.Modifiers...)
	}
	return &out
}

func mergeBackendResponse(tpl, cur *dto.BackendResponse) *dto.BackendResponse {
	if checker.IsNil(tpl) && checker.IsNil(cur) {
		return nil
	}
	var out dto.BackendResponse
	if checker.NonNil(tpl) {
		out = *tpl
	}
	if checker.NonNil(cur) {
		// bool aqui: caller decide (não fazemos OR), porque omit=true é uma escolha forte.
		out.ContinueOnError = cur.ContinueOnError
		out.Omit = cur.Omit

		out.Header = mergeBackendResponseHeader(out.Header, cur.Header)
		out.Body = mergeBackendResponseBody(out.Body, cur.Body)
	}
	return &out
}

func mergeBackendResponseHeader(tpl, cur *dto.BackendResponseHeader) *dto.BackendResponseHeader {
	if checker.IsNil(tpl) && checker.IsNil(cur) {
		return nil
	}

	var out dto.BackendResponseHeader
	if checker.NonNil(tpl) {
		out = *tpl
	}

	if checker.IsNil(cur) {
		return &out
	}

	out.Omit = out.Omit || cur.Omit
	if checker.NonNil(cur.Mapper) {
		out.Mapper = cur.Mapper
	}
	if checker.NonNil(cur.Projector) {
		out.Projector = cur.Projector
	}
	out.Modifiers = append(append([]dto.Modifier{}, out.Modifiers...), cur.Modifiers...)

	return &out
}

func mergeBackendResponseBody(tpl, cur *dto.BackendResponseBody) *dto.BackendResponseBody {
	if checker.IsNil(tpl) && checker.IsNil(cur) {
		return nil
	}

	var out dto.BackendResponseBody
	if checker.NonNil(tpl) {
		out = *tpl
	}
	if checker.NonNil(cur) {
		out.Omit = out.Omit || cur.Omit
		if checker.IsNotEmpty(cur.Group) {
			out.Group = cur.Group
		}
		if checker.NonNil(cur.Mapper) {
			out.Mapper = cur.Mapper
		}
		if checker.NonNil(cur.Projector) {
			out.Projector = cur.Projector
		}
		out.Modifiers = append(append([]dto.Modifier{}, out.Modifiers...), cur.Modifiers...)
	}
	return &out
}

func mergeBackendPropagate(tpl, cur *dto.BackendPropagate) *dto.BackendPropagate {
	if checker.IsNil(tpl) && checker.IsNil(cur) {
		return nil
	}

	var out dto.BackendPropagate
	if checker.NonNil(tpl) {
		out = *tpl
	}

	if checker.IsNil(cur) {
		return &out
	}

	if checker.NonNil(cur.ContinueOnError) {
		out.ContinueOnError = cur.ContinueOnError
	}

	out.Header = mergeBackendRequestHeader(out.Header, cur.Header)
	out.Param = mergeBackendRequestParam(out.Param, cur.Param)
	out.Query = mergeBackendRequestQuery(out.Query, cur.Query)
	out.Body = mergeBackendRequestBody(out.Body, cur.Body)

	return &out
}

func applyBackendPropagateIntoState(p *dto.BackendPropagate, ps *propagateState, backendIndex int) {
	if checker.IsNil(p) {
		return
	}

	if checker.NonNil(p.ContinueOnError) {
		ps.continueOnError = converter.ToPointer(*p.ContinueOnError)
	}

	var rewritten = *p

	if checker.NonNil(rewritten.Header) {
		rewritten.Header = rewriteBackendRequestHeaderResponseRefs(rewritten.Header, backendIndex)
	}
	if checker.NonNil(rewritten.Param) {
		rewritten.Param = rewriteBackendRequestParamResponseRefs(rewritten.Param, backendIndex)
	}
	if checker.NonNil(rewritten.Query) {
		rewritten.Query = rewriteBackendRequestQueryResponseRefs(rewritten.Query, backendIndex)
	}
	if checker.NonNil(rewritten.Body) {
		rewritten.Body = rewriteBackendRequestBodyResponseRefs(rewritten.Body, backendIndex)
	}

	ps.header = mergeBackendRequestHeader(ps.header, rewritten.Header)
	ps.param = mergeBackendRequestParam(ps.param, rewritten.Param)
	ps.query = mergeBackendRequestQuery(ps.query, rewritten.Query)
	ps.body = mergeBackendRequestBody(ps.body, rewritten.Body)
}

func collectPropagatingModifiersFromRequestIntoState(req *dto.BackendRequest, ps *propagateState, backendIndex int) {
	if checker.IsNil(req) {
		return
	}

	if checker.NonNil(req.Header) && checker.IsNotEmpty(req.Header.Modifiers) {
		ps.header = mergeBackendRequestHeader(ps.header, &dto.BackendRequestHeader{
			Modifiers: rewriteModifiersResponseRefs(onlyPropagate(req.Header.Modifiers), backendIndex),
		})
	}
	if checker.NonNil(req.Param) && checker.IsNotEmpty(req.Param.Modifiers) {
		ps.param = mergeBackendRequestParam(ps.param, &dto.BackendRequestParam{
			Modifiers: rewriteModifiersResponseRefs(onlyPropagate(req.Param.Modifiers), backendIndex),
		})
	}
	if checker.NonNil(req.Query) && checker.IsNotEmpty(req.Query.Modifiers) {
		ps.query = mergeBackendRequestQuery(ps.query, &dto.BackendRequestQuery{
			Modifiers: rewriteModifiersResponseRefs(onlyPropagate(req.Query.Modifiers), backendIndex),
		})
	}
	if checker.NonNil(req.Body) && checker.IsNotEmpty(req.Body.Modifiers) {
		ps.body = mergeBackendRequestBody(ps.body, &dto.BackendRequestBody{
			Modifiers: rewriteModifiersResponseRefs(onlyPropagate(req.Body.Modifiers), backendIndex),
		})
	}
}

func rewriteBackendRequestHeaderResponseRefs(h *dto.BackendRequestHeader, backendIndex int) *dto.BackendRequestHeader {
	if checker.IsNil(h) {
		return nil
	}
	out := *h
	out.Modifiers = rewriteModifiersResponseRefs(out.Modifiers, backendIndex)
	return &out
}

func rewriteBackendRequestParamResponseRefs(p *dto.BackendRequestParam, backendIndex int) *dto.BackendRequestParam {
	if checker.IsNil(p) {
		return nil
	}
	out := *p
	out.Modifiers = rewriteModifiersResponseRefs(out.Modifiers, backendIndex)
	return &out
}

func rewriteBackendRequestQueryResponseRefs(q *dto.BackendRequestQuery, backendIndex int) *dto.BackendRequestQuery {
	if checker.IsNil(q) {
		return nil
	}
	out := *q
	out.Modifiers = rewriteModifiersResponseRefs(out.Modifiers, backendIndex)
	return &out
}

func rewriteBackendRequestBodyResponseRefs(b *dto.BackendRequestBody, backendIndex int) *dto.BackendRequestBody {
	if checker.IsNil(b) {
		return nil
	}
	out := *b
	out.Modifiers = rewriteModifiersResponseRefs(out.Modifiers, backendIndex)
	return &out
}

func rewriteModifiersResponseRefs(in []dto.Modifier, backendIndex int) []dto.Modifier {
	if len(in) == 0 {
		return nil
	}
	out := make([]dto.Modifier, 0, len(in))
	for _, m := range in {
		m.Value = rewriteResponseRef(m.Value, backendIndex)
		out = append(out, m)
	}
	return out
}

func rewriteResponseRef(value string, backendIndex int) string {
	const prefix = "#response."
	if !strings.Contains(value, prefix) {
		return value
	}
	return strings.ReplaceAll(value, prefix, fmt.Sprintf("#responses[%d].", backendIndex))
}

func onlyPropagate(in []dto.Modifier) []dto.Modifier {
	if len(in) == 0 {
		return nil
	}
	out := make([]dto.Modifier, 0, len(in))
	for _, m := range in {
		if m.Propagate {
			out = append(out, m)
		}
	}
	return out
}

func resolveAsync(cur *bool, endpointParallelism bool) bool {
	if checker.NonNil(cur) {
		return *cur
	}
	return endpointParallelism
}
