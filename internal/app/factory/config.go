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
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type propagateState struct {
	header  *dto.MetadataTransformation
	urlPath *dto.URLPathTransformation
	query   *dto.QueryTransformation
	body    *dto.PayloadTransformation
}

func BuildGopen(gopen *dto.Gopen) *vo.GopenConfig {
	return vo.NewGopenConfig(buildEndpoints(gopen))
}

func buildEndpoints(gopen *dto.Gopen) []vo.EndpointConfig {
	var endpoints []vo.EndpointConfig
	var errs []string

	for _, endpoint := range gopen.Endpoints {
		var err string
		for _, registeredEndpoint := range endpoints {
			if checker.NotEquals(endpoint.Path, registeredEndpoint.Path()) ||
				checker.NotEquals(endpoint.Method, registeredEndpoint.Method()) {
				continue
			}
			errs = append(errs, fmt.Sprintf("- Duplicate endpoint path: %s method: %s", endpoint.Path, endpoint.Method))
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

func buildEndpoint(gopen *dto.Gopen, endpoint dto.Endpoint) vo.EndpointConfig {
	return vo.NewEndpointConfig(
		resolveEndpointExecution(endpoint.Execution),
		endpoint.Path,
		endpoint.Method,
		buildTimeout(gopen.Timeout, endpoint.Timeout),
		buildSecurityCors(gopen.SecurityCors, endpoint.SecurityCors),
		buildLimiter(gopen.Limiter, endpoint.Limiter),
		buildEndpointCache(gopen.Cache, endpoint.Cache),
		buildBackends(gopen.Templates, endpoint),
		buildEndpointResponse(endpoint.Response),
	)
}

func buildTimeout(base, timeout vo.Duration) vo.Duration {
	if checker.IsGreaterThan(timeout, 0) {
		return timeout
	} else {
		return base
	}
}

func buildSecurityCors(securityCors *dto.SecurityCors, endpointSecurityCors *dto.SecurityCors,
) *vo.SecurityCorsConfig {
	var onlyIf, ignoreIf, allowOrigins, allowMethods, allowHeaders []string
	var allowCredentials bool

	if checker.NonNil(securityCors) {
		if checker.NonNil(securityCors.OnlyIf) {
			onlyIf = securityCors.OnlyIf
		}
		if checker.NonNil(securityCors.IgnoreIf) {
			ignoreIf = securityCors.IgnoreIf
		}
		if checker.NonNil(securityCors.AllowOrigins) {
			allowOrigins = securityCors.AllowOrigins
		}
		if checker.NonNil(securityCors.AllowMethods) {
			allowMethods = securityCors.AllowMethods
		}
		if checker.NonNil(securityCors.AllowHeaders) {
			allowHeaders = securityCors.AllowHeaders
		}
		allowCredentials = securityCors.AllowCredentials
	}

	if checker.NonNil(endpointSecurityCors) {
		if checker.NonNil(endpointSecurityCors.OnlyIf) {
			onlyIf = endpointSecurityCors.OnlyIf
		}
		if checker.NonNil(endpointSecurityCors.IgnoreIf) {
			ignoreIf = endpointSecurityCors.IgnoreIf
		}
		if checker.NonNil(endpointSecurityCors.AllowOrigins) {
			allowOrigins = endpointSecurityCors.AllowOrigins
		}
		if checker.NonNil(endpointSecurityCors.AllowMethods) {
			allowMethods = endpointSecurityCors.AllowMethods
		}
		if checker.NonNil(endpointSecurityCors.AllowHeaders) {
			allowHeaders = endpointSecurityCors.AllowHeaders
		}
		allowCredentials = endpointSecurityCors.AllowCredentials
	}

	return vo.NewSecurityCorsConfig(onlyIf, ignoreIf, allowOrigins, allowMethods, allowHeaders, allowCredentials)
}

func buildLimiter(limiter *dto.Limiter, endpointLimiter *dto.Limiter) *vo.LimiterConfig {
	if checker.IsNil(limiter) && checker.IsNil(endpointLimiter) {
		return nil
	}

	var endpointSize, size *dto.LimiterSize
	var endpointRate, rate *dto.LimiterRate

	if checker.NonNil(limiter) {
		size = limiter.Size
		rate = limiter.Rate
	}

	if checker.NonNil(endpointLimiter) {
		endpointSize = endpointLimiter.Size
		endpointRate = endpointLimiter.Rate
	}

	return vo.NewLimiterConfig(buildLimiterSize(size, endpointSize), buildLimiterRate(rate, endpointRate))
}

func buildLimiterSize(size, endpointSize *dto.LimiterSize) *vo.LimiterSizeConfig {
	if checker.IsNil(size) && checker.IsNil(endpointSize) {
		return nil
	}

	var maxMetadata, maxPayload vo.Bytes

	if checker.NonNil(size) {
		if checker.NonNil(size.MaxHeader) {
			maxMetadata = *size.MaxHeader
		}
		if checker.NonNil(size.MaxBody) {
			maxPayload = *size.MaxBody
		}
	}

	if checker.NonNil(endpointSize) {
		if checker.NonNil(endpointSize.MaxHeader) {
			maxMetadata = *endpointSize.MaxHeader
		}
		if checker.NonNil(endpointSize.MaxBody) {
			maxPayload = *endpointSize.MaxBody
		}
	}

	return vo.NewLimiterSizeConfig(maxMetadata, maxPayload)
}

func buildLimiterRate(rate, endpointRate *dto.LimiterRate) *vo.LimiterRateConfig {
	if checker.IsNil(rate) && checker.IsNil(endpointRate) {
		return nil
	}

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

	return vo.NewLimiterRateConfig(every, capacity)
}

func buildEndpointCache(cache *dto.Cache, endpointCache *dto.Cache) *vo.CacheConfig {
	if checker.IsNil(cache) && checker.IsNil(endpointCache) {
		return nil
	}

	effective := mergeCacheDTO(cache, endpointCache)

	return vo.NewCacheConfig(
		enum.CacheKindEndpoint,
		vo.NewCacheDecisionConfig(
			mergeCacheConditions(effective.OnlyIf, effective.Read.OnlyIf),
			mergeCacheConditions(effective.IgnoreIf, effective.Read.IgnoreIf),
		),
		vo.NewCacheDecisionConfig(
			mergeCacheConditions(effective.OnlyIf, effective.Write.OnlyIf),
			mergeCacheConditions(effective.IgnoreIf, effective.Write.IgnoreIf),
		),
		effective.Key,
		effective.TTL,
	)
}

func mergeCacheDTO(base, override *dto.Cache) dto.Cache {
	if checker.IsNil(base) {
		return *override
	}
	if checker.IsNil(override) {
		return *base
	}
	result := *base
	if checker.IsNotEmpty(override.OnlyIf) {
		result.OnlyIf = override.OnlyIf
	}
	if checker.IsNotEmpty(override.IgnoreIf) {
		result.IgnoreIf = override.IgnoreIf
	}
	if checker.IsNotEmpty(override.Read.OnlyIf) {
		result.Read.OnlyIf = override.Read.OnlyIf
	}
	if checker.IsNotEmpty(override.Read.IgnoreIf) {
		result.Read.IgnoreIf = override.Read.IgnoreIf
	}
	if checker.IsNotEmpty(override.Write.OnlyIf) {
		result.Write.OnlyIf = override.Write.OnlyIf
	}
	if checker.IsNotEmpty(override.Write.IgnoreIf) {
		result.Write.IgnoreIf = override.Write.IgnoreIf
	}
	if checker.IsNotEmpty(override.Key) {
		result.Key = override.Key
	}
	if override.TTL > 0 {
		result.TTL = override.TTL
	}
	return result
}

func mergeCacheConditions(base, specific []string) []string {
	return append(base, specific...)
}

func buildEndpointResponse(endpointResponse *dto.EndpointResponse) vo.EndpointResponseConfig {
	return vo.NewEndpointResponseConfig(
		buildMetadata(endpointResponse.Header),
		buildPayload(endpointResponse.Body),
	)
}

func buildBackends(templates *dto.Templates, endpoint dto.Endpoint) []vo.BackendConfig {
	var result []vo.BackendConfig
	var ps propagateState
	var backendIndex int

	seen := map[string]struct{}{}
	idToIndex := map[string]int{}

	consume := func(flow enum.BackendFlow, items []dto.Backend, allowParallelism bool) {
		for _, b := range items {
			effective := resolveBackendTemplate(b, templates)

			validateBackendID(endpoint, flow, backendIndex, effective)
			validateBackendDependencies(endpoint, flow, backendIndex, effective, seen)
			markBackendSeen(endpoint, flow, backendIndex, effective, seen)

			result = append(result, buildBackendUnified(
				endpoint,
				effective,
				flow,
				&ps,
				backendIndex,
				idToIndex,
			))

			idToIndex[effective.ID] = backendIndex

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

	if checker.IsNotEmpty(missing) {
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
	endpoint dto.Endpoint,
	backend dto.Backend,
	flow enum.BackendFlow,
	ps *propagateState,
	backendIndex int,
	idToIndex map[string]int,
) vo.BackendConfig {
	var http *vo.BackendHTTPConfig
	var publisher *vo.BackendPublisherConfig

	switch backend.Kind {
	case enum.BackendKindPublisher:
		publisher = vo.NewBackendPublisherConfig(
			backend.Broker,
			backend.Path,
			backend.GroupID,
			backend.DeduplicationID,
			backend.Delay,
			buildPublisherMessage(backend.Message),
		)
	case enum.BackendKindHTTP:
		http = vo.NewBackendHTTPConfig(
			backend.Hosts,
			backend.Path,
			backend.Method,
			buildHTTPBackendRequest(backend, ps, backendIndex),
		)
	default:
		panic(errors.Newf("invalid backend.kind=%v (endpoint=%s %s)", backend.Kind, flow, backend.Path))
	}

	return vo.NewBackendConfig(
		flow,
		backend.OnlyIf,
		backend.IgnoreIf,
		backend.ID,
		resolveBackendExecution(endpoint.Execution, backend.Execution),
		buildBackendDependencies(backend.Dependencies, idToIndex),
		backend.Kind,
		buildBackendCache(backend.Cache),
		http,
		publisher,
		buildBackendResponse(backend, flow),
	)
}

func buildBackendDependencies(deps []string, idToIndex map[string]int) *vo.BackendDependenciesConfig {
	if checker.IsEmpty(deps) {
		return nil
	}

	var idxs []int
	for _, d := range deps {
		d = strings.TrimSpace(d)
		if idx, ok := idToIndex[d]; ok {
			idxs = append(idxs, idx)
		}
	}
	return vo.NewBackendDependenciesConfig(deps, idxs)
}

func buildBackendCache(cache *dto.Cache) *vo.CacheConfig {
	if checker.IsNil(cache) {
		return nil
	}
	return vo.NewCacheConfig(
		enum.CacheKindBackend,
		vo.NewCacheDecisionConfig(
			mergeCacheConditions(cache.OnlyIf, cache.Read.OnlyIf),
			mergeCacheConditions(cache.IgnoreIf, cache.Read.IgnoreIf),
		),
		vo.NewCacheDecisionConfig(
			mergeCacheConditions(cache.OnlyIf, cache.Write.OnlyIf),
			mergeCacheConditions(cache.IgnoreIf, cache.Write.IgnoreIf),
		),
		cache.Key,
		cache.TTL,
	)
}

func buildHTTPBackendRequest(backend dto.Backend, ps *propagateState, backendIndex int) vo.BackendHTTPRequestConfig {
	effective := mergeBackendRequestWithPropagation(backend.Request, ps)

	applyBackendPropagateIntoState(backend.Propagate, ps, backendIndex)
	collectPropagatingModifiersFromRequestIntoState(backend.Request, ps, backendIndex)

	return vo.NewBackendHTTPRequestConfig(
		buildMetadata(effective.Header),
		buildURLPath(effective.Param),
		buildQuery(effective.Query),
		buildPayload(effective.Body),
	)
}

func buildBackendResponse(backend dto.Backend, flow enum.BackendFlow) *vo.BackendResponseConfig {
	if checker.Equals(flow, enum.BackendFlowBeforeware) || checker.Equals(flow, enum.BackendFlowAfterware) {
		return buildMiddlewareBackendResponse(backend)
	} else if checker.IsNil(backend.Response) {
		return nil
	} else {
		return vo.NewBackendResponseConfig(
			backend.Response.Omit,
			buildMetadata(backend.Response.Header),
			buildPayload(backend.Response.Body),
		)
	}
}

func buildMiddlewareBackendResponse(backend dto.Backend) *vo.BackendResponseConfig {
	if checker.IsNil(backend.Response) {
		return vo.NewBackendResponseConfigForMiddleware(false, nil)
	} else {
		return vo.NewBackendResponseConfigForMiddleware(backend.Response.Omit, buildMetadata(backend.Response.Header))
	}
}

func buildPublisherMessage(publisherMessage dto.PublisherMessage) vo.BackendPublisherMessageConfig {
	return vo.NewBackendPublisherMessageConfig(
		publisherMessage.OnlyIf,
		publisherMessage.IgnoreIf,
		buildAttributeValues(publisherMessage.Attributes),
		buildPayload(publisherMessage.Body),
	)
}

func buildAttributeValues(attributes map[string]dto.AttributeValue) map[string]vo.AttributeValueConfig {
	if checker.IsNil(attributes) {
		return nil
	}

	var result = make(map[string]vo.AttributeValueConfig)
	for key, value := range attributes {
		result[key] = vo.NewAttributeValueConfig(value.Type, value.Value)
	}
	return result
}

func buildURLPath(backendRequestParam *dto.URLPathTransformation) *vo.URLPathConfig {
	if checker.IsNil(backendRequestParam) {
		return nil
	}
	return vo.NewURLPathConfig(buildModifiers(backendRequestParam.Modifiers))
}

func buildQuery(backendRequestQuery *dto.QueryTransformation) *vo.QueryConfig {
	if checker.IsNil(backendRequestQuery) {
		return nil
	}
	return vo.NewQueryConfig(
		backendRequestQuery.Omit,
		buildMapper(backendRequestQuery.Mapper),
		buildProjector(backendRequestQuery.Projector),
		buildModifiers(backendRequestQuery.Modifiers),
	)
}

func buildMetadata(metadata *dto.MetadataTransformation) *vo.MetadataConfig {
	if checker.IsNil(metadata) {
		return nil
	}
	return vo.NewMetadataConfig(
		metadata.Omit,
		buildMapper(metadata.Mapper),
		buildProjector(metadata.Projector),
		buildModifiers(metadata.Modifiers),
	)
}

func buildPayload(payload *dto.PayloadTransformation) *vo.PayloadConfig {
	if checker.IsNil(payload) {
		return nil
	}
	return vo.NewPayloadConfig(
		payload.Aggregate,
		payload.Omit,
		payload.OmitEmpty,
		payload.Group,
		payload.ContentType,
		payload.ContentEncoding,
		payload.Nomenclature,
		buildMapper(payload.Mapper),
		buildProjector(payload.Projector),
		buildModifiers(payload.Modifiers),
		buildJoins(payload.Joins),
	)
}

func buildMapper(mapper *dto.Mapper) *vo.MapperConfig {
	if checker.IsNil(mapper) {
		return nil
	}
	return vo.NewMapperConfig(mapper.OnlyIf, mapper.IgnoreIf, mapper.Policy, mapper.Map)
}

func buildProjector(projector *dto.Projector) *vo.ProjectorConfig {
	if checker.IsNil(projector) {
		return nil
	}
	return vo.NewProjectorConfig(projector.OnlyIf, projector.IgnoreIf, projector.Project)
}

func buildModifiers(modifiers []dto.Modifier) []vo.ModifierConfig {
	var result []vo.ModifierConfig
	for _, modifier := range modifiers {
		result = append(result, buildModifier(modifier))
	}
	return result
}

func buildModifier(modifier dto.Modifier) vo.ModifierConfig {
	return vo.NewModifierConfig(modifier.OnlyIf, modifier.IgnoreIf, modifier.Action, modifier.Propagate, modifier.Key, modifier.Value)
}

func buildJoins(joins []dto.Join) []vo.JoinConfig {
	var result []vo.JoinConfig
	for _, join := range joins {
		result = append(result, buildJoin(join))
	}
	return result
}

func buildJoin(join dto.Join) vo.JoinConfig {
	return vo.NewJoin(
		join.OnlyIf,
		join.IgnoreIf,
		vo.NewJoinSource(join.Source.Path, join.Source.Key),
		vo.NewJoinTarget(
			join.Target.Policy,
			join.Target.Path,
			join.Target.Key,
			join.Target.As,
			join.Target.OnMissing,
		),
	)
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
	if checker.Equals(mergeMode, enum.TemplateMergeBase) {
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
		ID:           tpl.ID,
		Dependencies: tpl.Dependencies,
		Kind:         tpl.Kind,
	}

	switch tpl.Kind {
	case enum.BackendKindPublisher:
		out.Broker = tpl.Broker
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
	merged.Kind = tpl.Kind

	switch tpl.Kind {
	case enum.BackendKindPublisher:
		merged.Broker = tpl.Broker
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
	if checker.NonNil(cur.Execution) {
		out.Execution = cur.Execution
	}
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

	if cur.Broker.IsEnumValid() {
		out.Broker = cur.Broker
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

func mergePublisherMessage(tpl, cur dto.PublisherMessage) dto.PublisherMessage {
	out := tpl

	if checker.IsNotEmpty(cur.OnlyIf) {
		out.OnlyIf = append(append([]string{}, out.OnlyIf...), cur.OnlyIf...)
	}
	if checker.IsNotEmpty(cur.IgnoreIf) {
		out.IgnoreIf = append(append([]string{}, out.IgnoreIf...), cur.IgnoreIf...)
	}
	if checker.NonNil(cur.Attributes) {
		if checker.IsNil(out.Attributes) {
			out.Attributes = map[string]dto.AttributeValue{}
		}
		for k, v := range cur.Attributes {
			out.Attributes[k] = v
		}
	}
	out.Body = mergePayloadTransformation(out.Body, cur.Body)

	return out
}

func mergeBackendRequest(tpl, cur dto.BackendRequest) dto.BackendRequest {
	out := tpl

	out.Header = mergeMetadataTransformation(out.Header, cur.Header)
	out.Param = mergeURLPathTransformation(out.Param, cur.Param)
	out.Query = mergeQueryTransformation(out.Query, cur.Query)
	out.Body = mergePayloadTransformation(out.Body, cur.Body)

	return out
}

func mergeBackendRequestWithPropagation(req dto.BackendRequest, ps *propagateState) dto.BackendRequest {
	if checker.IsNil(ps.header) &&
		checker.IsNil(ps.urlPath) &&
		checker.IsNil(ps.query) &&
		checker.IsNil(ps.body) {
		return req
	}

	out := req

	out.Header = mergeMetadataTransformation(ps.header, out.Header)
	out.Param = mergeURLPathTransformation(ps.urlPath, out.Param)
	out.Query = mergeQueryTransformation(ps.query, out.Query)
	out.Body = mergePayloadTransformation(ps.body, out.Body)

	return out
}

func mergeMetadataTransformation(inh, cur *dto.MetadataTransformation) *dto.MetadataTransformation {
	if checker.IsNil(inh) && checker.IsNil(cur) {
		return nil
	}
	var out dto.MetadataTransformation
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

func mergeURLPathTransformation(inh, cur *dto.URLPathTransformation) *dto.URLPathTransformation {
	if checker.IsNil(inh) && checker.IsNil(cur) {
		return nil
	}
	var out dto.URLPathTransformation
	if checker.NonNil(inh) {
		out = *inh
	}
	if checker.NonNil(cur) {
		out.Modifiers = append(append([]dto.Modifier{}, out.Modifiers...), cur.Modifiers...)
	}
	return &out
}

func mergeQueryTransformation(inh, cur *dto.QueryTransformation) *dto.QueryTransformation {
	if checker.IsNil(inh) && checker.IsNil(cur) {
		return nil
	}
	var out dto.QueryTransformation
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

func mergePayloadTransformation(inh, cur *dto.PayloadTransformation) *dto.PayloadTransformation {
	if checker.IsNil(inh) && checker.IsNil(cur) {
		return nil
	}
	var out dto.PayloadTransformation
	if checker.NonNil(inh) {
		out = *inh
	}
	if checker.NonNil(cur) {
		out.Omit = out.Omit || cur.Omit
		out.OmitEmpty = out.OmitEmpty || cur.OmitEmpty
		out.ContentType = cur.ContentType
		out.ContentEncoding = cur.ContentEncoding
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
		out.Joins = append(append([]dto.Join{}, out.Joins...), cur.Joins...)
	}
	return &out
}

func mergeBackendResponse(tpl, cur dto.BackendResponse) dto.BackendResponse {
	out := tpl

	out.Omit = out.Omit || cur.Omit
	out.Header = mergeMetadataTransformation(out.Header, cur.Header)
	out.Body = mergePayloadTransformation(out.Body, cur.Body)

	return out
}

func mergeBackendPropagate(tpl, cur dto.BackendPropagate) dto.BackendPropagate {
	out := tpl

	out.Header = mergeMetadataTransformation(out.Header, cur.Header)
	out.URLPath = mergeURLPathTransformation(out.URLPath, cur.URLPath)
	out.Query = mergeQueryTransformation(out.Query, cur.Query)
	out.Body = mergePayloadTransformation(out.Body, cur.Body)

	return out
}

func applyBackendPropagateIntoState(p dto.BackendPropagate, ps *propagateState, backendIndex int) {
	var rewritten = p

	if checker.NonNil(rewritten.Header) {
		rewritten.Header = rewriteBackendRequestHeaderResponseRefs(rewritten.Header, backendIndex)
	}
	if checker.NonNil(rewritten.URLPath) {
		rewritten.URLPath = rewriteBackendRequestParamResponseRefs(rewritten.URLPath, backendIndex)
	}
	if checker.NonNil(rewritten.Query) {
		rewritten.Query = rewriteBackendRequestQueryResponseRefs(rewritten.Query, backendIndex)
	}
	if checker.NonNil(rewritten.Body) {
		rewritten.Body = rewriteBackendRequestBodyResponseRefs(rewritten.Body, backendIndex)
	}

	ps.header = mergeMetadataTransformation(ps.header, rewritten.Header)
	ps.urlPath = mergeURLPathTransformation(ps.urlPath, rewritten.URLPath)
	ps.query = mergeQueryTransformation(ps.query, rewritten.Query)
	ps.body = mergePayloadTransformation(ps.body, rewritten.Body)
}

func collectPropagatingModifiersFromRequestIntoState(req dto.BackendRequest, ps *propagateState, backendIndex int) {
	if checker.NonNil(req.Header) && checker.IsNotEmpty(req.Header.Modifiers) {
		ps.header = mergeMetadataTransformation(ps.header, &dto.MetadataTransformation{
			Modifiers: rewriteModifiersResponseRefs(onlyPropagate(req.Header.Modifiers), backendIndex),
		})
	}
	if checker.NonNil(req.Param) && checker.IsNotEmpty(req.Param.Modifiers) {
		ps.urlPath = mergeURLPathTransformation(ps.urlPath, &dto.URLPathTransformation{
			Modifiers: rewriteModifiersResponseRefs(onlyPropagate(req.Param.Modifiers), backendIndex),
		})
	}
	if checker.NonNil(req.Query) && checker.IsNotEmpty(req.Query.Modifiers) {
		ps.query = mergeQueryTransformation(ps.query, &dto.QueryTransformation{
			Modifiers: rewriteModifiersResponseRefs(onlyPropagate(req.Query.Modifiers), backendIndex),
		})
	}
	if checker.NonNil(req.Body) && checker.IsNotEmpty(req.Body.Modifiers) {
		ps.body = mergePayloadTransformation(ps.body, &dto.PayloadTransformation{
			Modifiers: rewriteModifiersResponseRefs(onlyPropagate(req.Body.Modifiers), backendIndex),
		})
	}
}

func rewriteBackendRequestHeaderResponseRefs(h *dto.MetadataTransformation, backendIndex int) *dto.MetadataTransformation {
	if checker.IsNil(h) {
		return nil
	}
	out := *h
	out.Modifiers = rewriteModifiersResponseRefs(out.Modifiers, backendIndex)
	return &out
}

func rewriteBackendRequestParamResponseRefs(p *dto.URLPathTransformation, backendIndex int) *dto.URLPathTransformation {
	if checker.IsNil(p) {
		return nil
	}
	out := *p
	out.Modifiers = rewriteModifiersResponseRefs(out.Modifiers, backendIndex)
	return &out
}

func rewriteBackendRequestQueryResponseRefs(q *dto.QueryTransformation, backendIndex int) *dto.QueryTransformation {
	if checker.IsNil(q) {
		return nil
	}
	out := *q
	out.Modifiers = rewriteModifiersResponseRefs(out.Modifiers, backendIndex)
	return &out
}

func rewriteBackendRequestBodyResponseRefs(b *dto.PayloadTransformation, backendIndex int) *dto.PayloadTransformation {
	if checker.IsNil(b) {
		return nil
	}
	out := *b
	out.Modifiers = rewriteModifiersResponseRefs(out.Modifiers, backendIndex)
	return &out
}

func rewriteModifiersResponseRefs(in []dto.Modifier, backendIndex int) []dto.Modifier {
	if checker.IsEmpty(in) {
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
	if checker.IsEmpty(in) {
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

func resolveEndpointExecution(endpointExec *dto.EndpointExecution) vo.EndpointExecutionConfig {
	if checker.IsNil(endpointExec) {
		return vo.NewEndpointExecutionConfigDefault()
	}
	return vo.NewEndpointExecutionConfig(endpointExec.Mode, endpointExec.On)
}

func resolveBackendExecution(endpointExec *dto.EndpointExecution, backendExec *dto.BackendExecution) vo.BackendExecutionConfig {
	if checker.IsNil(endpointExec) && checker.IsNil(backendExec) {
		return vo.NewBackendExecutionConfigDefault()
	}

	var concurrent int
	var async bool
	var mode enum.ExecutionMode
	var on []enum.ExecutionOn

	if checker.NonNil(endpointExec) {
		async = endpointExec.Parallelism
		if endpointExec.Mode.IsEnumValid() {
			mode = endpointExec.Mode
		}
		if checker.IsNotEmpty(endpointExec.On) {
			on = endpointExec.On
		}
	}

	if checker.NonNil(backendExec) {
		concurrent = backendExec.Concurrent
		async = resolveAsync(backendExec.Async, async)
		if backendExec.Mode.IsEnumValid() {
			mode = backendExec.Mode
		}
		if checker.IsNotEmpty(backendExec.On) {
			on = backendExec.On
		}
	}

	return vo.NewBackendExecutionConfig(concurrent, async, mode, on)
}

func resolveAsync(cur *bool, endpointParallelism bool) bool {
	if checker.NonNil(cur) {
		return *cur
	}
	return endpointParallelism
}
