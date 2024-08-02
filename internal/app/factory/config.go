package factory

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"strings"
)

func BuildGopen(gopen *dto.Gopen) *vo.Gopen {
	return vo.NewGopen(
		gopen.Port,
		buildSecurityCors(gopen.SecurityCors),
		buildMiddlewares(gopen.Middlewares),
		buildEndpoints(gopen),
	)
}

func buildSecurityCors(securityCors *dto.SecurityCors) *vo.SecurityCors {
	if helper.IsNil(securityCors) {
		return nil
	}
	return vo.NewSecurityCors(securityCors.AllowOrigins, securityCors.AllowMethods, securityCors.AllowHeaders)
}

func buildMiddlewares(middlewares map[string]dto.Backend) vo.Middlewares {
	var mapp map[string]vo.Backend
	for k, v := range middlewares {
		mapp[k] = buildMiddlewareBackend(v)
	}
	return vo.NewMiddlewares(mapp)
}

func buildEndpoints(gopen *dto.Gopen) []vo.Endpoint {
	var endpoints []vo.Endpoint
	var errs []string

	for _, endpoint := range gopen.Endpoints {
		var err string
		for _, registeredEndpoint := range endpoints {
			if helper.IsNotEqualTo(endpoint.Path, registeredEndpoint.Path()) ||
				helper.IsNotEqualTo(endpoint.Method, registeredEndpoint.Method()) {
				continue
			}
			err = fmt.Sprintf("- Duplicate endpoint path: %s method: %s", endpoint.Path, endpoint.Method)
			errs = append(errs, err)
		}
		if helper.IsEmpty(err) {
			endpoints = append(endpoints, buildEndpoint(gopen, endpoint))
		}
	}

	if helper.IsNotEmpty(errs) {
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
		endpoint.AbortIfStatusCodes,
		buildEndpointResponse(endpoint.Response),
		endpoint.Beforewares,
		endpoint.Afterwares,
		buildBackends(endpoint.Backends),
	)
}

func buildTimeout(timeout, endpointTimeout vo.Duration) vo.Duration {
	if helper.IsGreaterThan(endpointTimeout, 0) {
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

	if helper.IsNotNil(limiter) {
		if helper.IsNotNil(limiter.MaxHeaderSize) {
			maxHeaderSize = *limiter.MaxHeaderSize
		}
		if helper.IsNotNil(limiter.MaxBodySize) {
			maxBodySize = *limiter.MaxBodySize
		}
		if helper.IsNotNil(limiter.MaxMultipartMemorySize) {
			maxMultipartForm = *limiter.MaxMultipartMemorySize
		}
		rate = limiter.Rate
	}

	if helper.IsNotNil(endpointLimiter) {
		if helper.IsNotNil(endpointLimiter.MaxHeaderSize) {
			maxHeaderSize = endpointLimiter.MaxHeaderSize
		}
		if helper.IsNotNil(endpointLimiter.MaxBodySize) {
			maxBodySize = endpointLimiter.MaxBodySize
		}
		if helper.IsNotNil(endpointLimiter.MaxMultipartMemorySize) {
			maxMultipartForm = endpointLimiter.MaxMultipartMemorySize
		}
		endpointRate = endpointLimiter.Rate
	}

	return vo.NewLimiter(maxHeaderSize, maxBodySize, maxMultipartForm, buildLimiterRate(rate, endpointRate))
}

func buildLimiterRate(rate, endpointRate *dto.Rate) vo.Rate {
	var every vo.Duration
	var capacity int

	if helper.IsNotNil(rate) {
		if helper.IsNotNil(rate.Every) {
			every = *rate.Every
		}
		if helper.IsNotNil(rate.Capacity) {
			capacity = *rate.Capacity
		}
	}

	if helper.IsNotNil(endpointRate) {
		if helper.IsNotNil(endpointRate.Every) {
			every = *endpointRate.Every
		}
		if helper.IsNotNil(endpointRate.Capacity) {
			capacity = *endpointRate.Capacity
		}
	}

	return vo.NewRate(every, capacity)
}

func buildCache(cache *dto.Cache, endpointCache *dto.EndpointCache) *vo.Cache {
	if helper.IsNil(cache) && helper.IsNil(endpointCache) {
		return nil
	}

	var enabled bool
	var ignoreQuery bool
	var duration vo.Duration
	var strategyHeaders []string
	var onlyIfStatusCodes []int
	var onlyIfMethods []string
	var allowCacheControl *bool

	if helper.IsNotNil(cache) {
		duration = cache.Duration
		strategyHeaders = cache.StrategyHeaders
		onlyIfStatusCodes = cache.OnlyIfStatusCodes
		onlyIfMethods = cache.OnlyIfMethods
		allowCacheControl = cache.AllowCacheControl
	}

	if helper.IsNotNil(endpointCache) {
		enabled = endpointCache.Enabled
		ignoreQuery = endpointCache.IgnoreQuery
		if helper.IsGreaterThan(endpointCache.Duration, 0) {
			duration = endpointCache.Duration
		}
		if helper.IsNotNil(endpointCache.StrategyHeaders) {
			strategyHeaders = endpointCache.StrategyHeaders
		}
		if helper.IsNotNil(endpointCache.AllowCacheControl) {
			allowCacheControl = endpointCache.AllowCacheControl
		}
		if helper.IsNotNil(endpointCache.OnlyIfStatusCodes) {
			onlyIfStatusCodes = endpointCache.OnlyIfStatusCodes
		}
	}

	return vo.NewCache(enabled, ignoreQuery, duration, strategyHeaders, onlyIfStatusCodes, onlyIfMethods, allowCacheControl)
}

func buildEndpointResponse(endpointResponse *dto.EndpointResponse) *vo.EndpointResponse {
	if helper.IsNil(endpointResponse) {
		return nil
	}
	return vo.NewEndpointResponse(
		endpointResponse.Aggregate,
		endpointResponse.ContentType,
		endpointResponse.ContentEncoding,
		endpointResponse.Nomenclature,
		endpointResponse.OmitEmpty,
	)
}

func buildBackends(backends []dto.Backend) []vo.Backend {
	var result []vo.Backend
	for _, backend := range backends {
		result = append(result, buildBackend(backend))
	}
	return result
}

func buildBackend(backend dto.Backend) vo.Backend {
	return vo.NewBackend(
		backend.Hosts,
		backend.Path,
		backend.Method,
		buildBackendRequest(backend.Request),
		buildBackendResponse(backend.Response),
	)
}

func buildMiddlewareBackend(backend dto.Backend) vo.Backend {
	return vo.NewBackend(
		backend.Hosts,
		backend.Path,
		backend.Method,
		buildBackendRequest(backend.Request),
		vo.NewBackendResponseForMiddleware(),
	)
}

func buildBackendRequest(backendRequest *dto.BackendRequest) *vo.BackendRequest {
	if helper.IsNil(backendRequest) {
		return nil
	}

	return vo.NewBackendRequest(
		backendRequest.OmitHeader,
		backendRequest.OmitQuery,
		backendRequest.OmitBody,
		backendRequest.ContentType,
		backendRequest.ContentEncoding,
		backendRequest.Nomenclature,
		backendRequest.OmitEmpty,
		backendRequest.HeaderMapper,
		backendRequest.QueryMapper,
		backendRequest.BodyMapper,
		backendRequest.HeaderProjection,
		backendRequest.QueryProjection,
		backendRequest.BodyProjection,
		buildModifiers(backendRequest.HeaderModifiers),
		buildModifiers(backendRequest.ParamModifiers),
		buildModifiers(backendRequest.QueryModifiers),
		buildModifiers(backendRequest.BodyModifiers),
	)
}

func buildBackendResponse(backendResponse *dto.BackendResponse) *vo.BackendResponse {
	if helper.IsNil(backendResponse) {
		return nil
	}

	return vo.NewBackendResponse(
		backendResponse.Omit,
		backendResponse.OmitHeader,
		backendResponse.OmitBody,
		backendResponse.Group,
		backendResponse.HeaderMapper,
		backendResponse.BodyMapper,
		backendResponse.HeaderProjection,
		backendResponse.BodyProjection,
		buildModifiers(backendResponse.HeaderModifiers),
		buildModifiers(backendResponse.BodyModifiers),
	)
}

func buildModifiers(modifiers []dto.Modifier) []vo.Modifier {
	var result []vo.Modifier
	for _, modifier := range modifiers {
		result = append(result, buildModifier(modifier))
	}
	return result
}

func buildModifier(modifier dto.Modifier) vo.Modifier {
	return vo.NewModifier(modifier.Action, modifier.Propagate, modifier.Key, modifier.Value)
}
