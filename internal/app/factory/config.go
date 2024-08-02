package factory

import (
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"strings"
)

func BuildGopen(gopen *dto.Gopen) *vo.Gopen {
	return vo.NewGopen(
		gopen.Port,
		buildSecurityCors(gopen.SecurityCors),
		buildEndpoints(gopen),
	)
}

func buildSecurityCors(securityCors *dto.SecurityCors) *vo.SecurityCors {
	if helper.IsNil(securityCors) {
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
		buildBackends(gopen.Middlewares, endpoint),
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

func buildBackends(middlewares map[string]dto.Backend, endpoint dto.Endpoint) []vo.Backend {
	var result []vo.Backend

	propagateHeaderModifiers := helper.ConvertToPointer([]vo.Modifier{})
	propagateParamModifiers := helper.ConvertToPointer([]vo.Modifier{})
	propagateQueryModifiers := helper.ConvertToPointer([]vo.Modifier{})
	propagateBodyModifiers := helper.ConvertToPointer([]vo.Modifier{})

	result = append(result, buildMiddlewareBackend(endpoint.Beforewares, middlewares, enum.BackendTypeBeforeware,
		propagateHeaderModifiers, propagateParamModifiers, propagateBodyModifiers, propagateQueryModifiers)...)

	result = append(result, buildNormalBackend(endpoint.Backends, propagateHeaderModifiers, propagateParamModifiers,
		propagateBodyModifiers, propagateQueryModifiers)...)

	result = append(result, buildMiddlewareBackend(endpoint.Afterwares, middlewares, enum.BackendTypeAfterware,
		propagateHeaderModifiers, propagateParamModifiers, propagateBodyModifiers, propagateQueryModifiers)...)

	return result
}

func buildNormalBackend(backends []dto.Backend, propagateHeaderModifiers, propagateParamModifiers,
	propagateBodyModifiers, propagateQueryModifiers *[]vo.Modifier) []vo.Backend {
	var result []vo.Backend
	for _, backend := range backends {
		result = append(result, buildBackend(backend, enum.BackendTypeNormal, propagateHeaderModifiers,
			propagateParamModifiers, propagateBodyModifiers, propagateQueryModifiers))
	}
	return result
}

func buildMiddlewareBackend(middlewareKeys []string, middlewares map[string]dto.Backend, backendType enum.BackendType,
	propagateHeaderModifiers, propagateParamModifiers, propagateBodyModifiers, propagateQueryModifiers *[]vo.Modifier,
) []vo.Backend {
	var result []vo.Backend
	for _, middlewareKey := range middlewareKeys {
		middleware, ok := middlewares[middlewareKey]
		if !ok {
			panic(errors.Newf("Middleware \"%s\" not configured on middlewares field!", middlewareKey))
		}
		result = append(result, buildBackend(middleware, backendType, propagateHeaderModifiers, propagateParamModifiers,
			propagateBodyModifiers, propagateQueryModifiers))
	}
	return result
}

func buildBackend(
	backend dto.Backend,
	backendType enum.BackendType,
	propagateHeaderModifiers,
	propagateParamModifiers,
	propagateQueryModifiers,
	propagateBodyModifiers *[]vo.Modifier,
) vo.Backend {
	return vo.NewBackend(
		backendType,
		backend.Hosts,
		backend.Path,
		backend.Method,
		buildBackendRequest(backend, propagateHeaderModifiers, propagateParamModifiers, propagateQueryModifiers, propagateBodyModifiers),
		buildBackendResponse(backend, backendType),
	)
}

func buildBackendRequest(
	backend dto.Backend,
	propagateHeaderModifiers,
	propagateParamModifiers,
	propagateQueryModifiers,
	propagateBodyModifiers *[]vo.Modifier,
) *vo.BackendRequest {
	if helper.IsNotNil(backend.Request) {
		buildAndPropagateModifiers(backend.Request.HeaderModifiers, propagateHeaderModifiers)
		buildAndPropagateModifiers(backend.Request.ParamModifiers, propagateParamModifiers)
		buildAndPropagateModifiers(backend.Request.QueryModifiers, propagateQueryModifiers)
		buildAndPropagateModifiers(backend.Request.BodyModifiers, propagateBodyModifiers)

		return vo.NewBackendRequest(
			backend.Request.OmitHeader,
			backend.Request.OmitQuery,
			backend.Request.OmitBody,
			backend.Request.ContentType,
			backend.Request.ContentEncoding,
			backend.Request.Nomenclature,
			backend.Request.OmitEmpty,
			backend.Request.HeaderMapper,
			backend.Request.QueryMapper,
			backend.Request.BodyMapper,
			backend.Request.HeaderProjection,
			backend.Request.QueryProjection,
			backend.Request.BodyProjection,
			*propagateHeaderModifiers,
			*propagateParamModifiers,
			*propagateQueryModifiers,
			*propagateBodyModifiers,
		)
	} else if helper.IsNotEmpty(propagateHeaderModifiers) || helper.IsNotEmpty(propagateParamModifiers) ||
		helper.IsNotEmpty(propagateQueryModifiers) || helper.IsNotEmpty(propagateBodyModifiers) {
		return vo.NewBackendRequestOnlyModifiers(
			*propagateHeaderModifiers,
			*propagateParamModifiers,
			*propagateQueryModifiers,
			*propagateBodyModifiers,
		)
	} else {
		return nil
	}
}

func buildAndPropagateModifiers(modifiers []dto.Modifier, propagateModifiers *[]vo.Modifier) {
	newModifiers := buildModifiers(modifiers)
	*propagateModifiers = append(*propagateModifiers, newModifiers...)
	for _, newModifier := range newModifiers {
		if newModifier.Propagate() {
			*propagateModifiers = append(*propagateModifiers, newModifier)
		}
	}
}

func buildBackendResponse(backend dto.Backend, backendType enum.BackendType) *vo.BackendResponse {
	if helper.IsNil(backend.Response) {
		return nil
	} else if helper.Equals(backendType, enum.BackendTypeBeforeware) || helper.Equals(backendType, enum.BackendTypeAfterware) {
		return vo.NewBackendResponseForMiddleware()
	}

	return vo.NewBackendResponse(
		backend.Response.Omit,
		backend.Response.OmitHeader,
		backend.Response.OmitBody,
		backend.Response.Group,
		backend.Response.HeaderMapper,
		backend.Response.BodyMapper,
		backend.Response.HeaderProjection,
		backend.Response.BodyProjection,
		buildModifiers(backend.Response.HeaderModifiers),
		buildModifiers(backend.Response.BodyModifiers),
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
