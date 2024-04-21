package mapper

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"os"
)

// BuildSettingViewDTO builds a `SettingView` DTO object using the provided `Gopen` object as input.
// It retrieves various properties from the `Gopen` object and sets them on the `SettingView` object.
func BuildSettingViewDTO(gopenVO *vo.Gopen) dto.SettingView {
	return dto.SettingView{
		Version:     os.Getenv("VERSION"),
		VersionDate: os.Getenv("VERSION_DATE"),
		Founder:     os.Getenv("FOUNDER"),
		CodeHelpers: os.Getenv("CODE_HELPERS"),
		Endpoints:   gopenVO.CountEndpoints(),
		Middlewares: gopenVO.CountMiddlewares(),
		Backends:    gopenVO.CountBackends(),
		Modifiers:   gopenVO.CountModifiers(),
		Setting:     BuildGopenDTO(gopenVO),
	}
}

// BuildGopenDTO builds a `Gopen` DTO object using the provided `Gopen` object as input.
// It retrieves various properties from the `Gopen` object and sets them on the `Gopen` object.
func BuildGopenDTO(gopenVO *vo.Gopen) dto.Gopen {
	return dto.Gopen{
		Version:      gopenVO.Version(),
		Port:         gopenVO.Port(),
		HotReload:    gopenVO.HotReload(),
		Timeout:      gopenVO.Timeout().String(),
		Limiter:      BuildLimiterDTOFromVO(gopenVO.Limiter()),
		Cache:        BuildCacheDTOFromVO(gopenVO.Cache()),
		SecurityCors: BuildSecurityCorsDTOFromVO(gopenVO.SecurityCors()),
		Middlewares:  BuildMiddlewaresDTOFromVO(gopenVO.Middlewares()),
		Endpoints:    BuildEndpointsDTOFromVOs(gopenVO.PureEndpoints()),
	}
}

// BuildGopenDTOFromCMD builds a `Gopen` DTO object using the provided `Gopen` and `Store` objects as input.
// It retrieves various properties from the `Gopen` object and sets them on the `Gopen` DTO object.
func BuildGopenDTOFromCMD(gopenVO *vo.Gopen, storeDTO *dto.Store) dto.Gopen {
	return dto.Gopen{
		Version:      gopenVO.Version(),
		Port:         gopenVO.Port(),
		Store:        storeDTO,
		HotReload:    gopenVO.HotReload(),
		Timeout:      gopenVO.Timeout().String(),
		Limiter:      BuildLimiterDTOFromVO(gopenVO.Limiter()),
		Cache:        BuildCacheDTOFromVO(gopenVO.Cache()),
		SecurityCors: BuildSecurityCorsDTOFromVO(gopenVO.SecurityCors()),
		Middlewares:  BuildMiddlewaresDTOFromVO(gopenVO.Middlewares()),
		Endpoints:    BuildEndpointsDTOFromVOs(gopenVO.PureEndpoints()),
	}
}

// BuildLimiterDTOFromVO builds a `Limiter` DTO object using the provided `Limiter` object as input.
// It retrieves various properties from the `Limiter` object and sets them on the `Limiter` DTO object.
func BuildLimiterDTOFromVO(limiterVO vo.Limiter) *dto.Limiter {
	maxHeaderSize := limiterVO.MaxHeaderSize()
	maxBodySize := limiterVO.MaxBodySize()
	maxMultipartMemorySize := limiterVO.MaxMultipartMemorySize()
	return &dto.Limiter{
		MaxHeaderSize:          maxHeaderSize.String(),
		MaxBodySize:            maxBodySize.String(),
		MaxMultipartMemorySize: maxMultipartMemorySize.String(),
		Rate:                   BuildRateDTOFromVO(limiterVO.Rate()),
	}
}

// BuildRateDTOFromVO builds a `Rate` DTO object using the provided `Rate` object as input.
// It retrieves various properties from the `Rate` object and sets them on the `Rate` DTO object.
// If the input `Rate` object is nil, it returns nil.
// Otherwise, it returns a pointer to a new `Rate` object with the following properties:
// - Capacity: the capacity field from the input `Rate` object
// - Every: the every field from the input `Rate` object as a string
// The returned `Rate` object is allocated on the heap.
func BuildRateDTOFromVO(rate vo.Rate) dto.Rate {
	return dto.Rate{
		Capacity: rate.Capacity(),
		Every:    rate.Every().String(),
	}
}

func BuildEndpointRateDTOFromVO(rateVO *vo.EndpointRate) *dto.Rate {
	return &dto.Rate{
		Capacity: rateVO.Capacity(),
		Every:    rateVO.EveryStr(),
	}
}

// BuildEndpointLimiterDTOFromEndpointVO builds a `Limiter` DTO object using the provided `EndpointLimiter` object as input.
// It retrieves various properties from the `EndpointLimiter` object and sets them on the `Limiter` object.
func BuildEndpointLimiterDTOFromEndpointVO(endpointLimiterVO *vo.EndpointLimiter) *dto.EndpointLimiter {
	if helper.IsNil(endpointLimiterVO) {
		return nil
	}
	maxHeaderSize := endpointLimiterVO.MaxHeaderSize()
	maxBodySize := endpointLimiterVO.MaxBodySize()
	maxMultipartMemorySize := endpointLimiterVO.MaxMultipartMemorySize()
	return &dto.EndpointLimiter{
		MaxHeaderSize:          maxHeaderSize.String(),
		MaxBodySize:            maxBodySize.String(),
		MaxMultipartMemorySize: maxMultipartMemorySize.String(),
		Rate:                   BuildEndpointRateDTOFromVO(endpointLimiterVO.Rate()),
	}
}

// BuildCacheDTOFromVO builds a `Cache` DTO object using the provided `Cache` object as input.
// It retrieves various properties from the `Cache` object and sets them on the `Cache` DTO object.
func BuildCacheDTOFromVO(cacheVO *vo.Cache) *dto.Cache {
	if helper.IsNil(cacheVO) {
		return nil
	}
	return &dto.Cache{
		Duration:          cacheVO.Duration().String(),
		StrategyHeaders:   cacheVO.StrategyHeaders(),
		OnlyIfStatusCodes: cacheVO.OnlyIfStatusCodes(),
		OnlyIfMethods:     cacheVO.OnlyIfMethods(),
		AllowCacheControl: cacheVO.AllowCacheControl(),
	}
}

// BuildEndpointCacheDTOFromVO builds an `EndpointCache` DTO object using the provided `EndpointCache` object as input.
// It retrieves various properties from the `EndpointCache` object and sets them on the `EndpointCache` object.
func BuildEndpointCacheDTOFromVO(endpointCacheVO *vo.EndpointCache) *dto.EndpointCache {
	if helper.IsNil(endpointCacheVO) || endpointCacheVO.Disabled() {
		return nil
	}
	return &dto.EndpointCache{
		Enabled:           endpointCacheVO.Enabled(),
		IgnoreQuery:       endpointCacheVO.IgnoreQuery(),
		Duration:          endpointCacheVO.DurationStr(),
		StrategyHeaders:   endpointCacheVO.StrategyHeaders(),
		OnlyIfStatusCodes: endpointCacheVO.OnlyIfStatusCodes(),
		AllowCacheControl: endpointCacheVO.AllowCacheControl(),
	}
}

// BuildSecurityCorsDTOFromVO builds a `SecurityCors` DTO object using the provided `SecurityCors` VO object as input.
// It checks if the `securityCorsVO` object is nil, and if so, returns nil. Otherwise, it creates a new `SecurityCors`
// object and sets its `AllowOrigins`, `AllowMethods`, and `AllowHeaders` fields by calling the corresponding methods
// on the `securityCorsVO` object. Finally, it returns the created `SecurityCors` object.
func BuildSecurityCorsDTOFromVO(securityCorsVO *vo.SecurityCors) *dto.SecurityCors {
	if helper.IsNil(securityCorsVO) {
		return nil
	}
	return &dto.SecurityCors{
		AllowOrigins: securityCorsVO.AllowOriginsData(),
		AllowMethods: securityCorsVO.AllowMethodsData(),
		AllowHeaders: securityCorsVO.AllowHeadersData(),
	}
}

// BuildMiddlewaresDTOFromVO builds a map of string keys to `Backend` DTO objects using the provided `Middlewares` object as input.
// It iterates over each key-value pair in the `Middlewares` map and calls the `BuildBackendDTO` function to create a `Backend` DTO object
// for each `Backend` value. The resulting `Backend` DTO objects are then stored in the `result` map with the corresponding key.
// Finally, the `result` map is returned.
func BuildMiddlewaresDTOFromVO(middlewaresVO vo.Middlewares) map[string]dto.Backend {
	result := map[string]dto.Backend{}
	for key, backendVO := range middlewaresVO {
		result[key] = BuildBackendDTOFromVO(backendVO)
	}
	return result
}

// BuildEndpointsDTOFromVOs builds a slice of `Endpoint` DTO objects using the provided `[]Endpoint` as input.
// It iterates over each `Endpoint` object and calls `BuildEndpointDTO` to build the individual `Endpoint` DTO object.
// The resulting DTO objects are then appended to the `result` slice and returned.
func BuildEndpointsDTOFromVOs(endpoints []vo.Endpoint) (result []dto.Endpoint) {
	for _, endpointVO := range endpoints {
		result = append(result, BuildEndpointDTOFromVO(endpointVO))
	}
	return result
}

// BuildEndpointDTOFromVO builds a `Endpoint` DTO object using the provided `Endpoint` object as input.
// It retrieves various properties from the `Endpoint` object and sets them on the `Endpoint` DTO object.
func BuildEndpointDTOFromVO(endpointVO vo.Endpoint) dto.Endpoint {
	return dto.Endpoint{
		Comment:            endpointVO.Comment(),
		Path:               endpointVO.Path(),
		Method:             endpointVO.Method(),
		Timeout:            endpointVO.TimeoutStr(),
		Limiter:            BuildEndpointLimiterDTOFromEndpointVO(endpointVO.Limiter()),
		Cache:              BuildEndpointCacheDTOFromVO(endpointVO.Cache()),
		ResponseEncode:     endpointVO.ResponseEncode(),
		AggregateResponses: endpointVO.AggregateResponses(),
		AbortIfStatusCodes: endpointVO.AbortIfStatusCodes(),
		Beforeware:         endpointVO.Beforeware(),
		Afterware:          endpointVO.Afterware(),
		Backends:           BuildBackendsDTOFromVO(endpointVO.Backends()),
	}
}

// BuildBackendsDTOFromVO builds an array of `Backend` DTO objects using the provided array of `Backend` objects as input.
// It iterates over each `Backend` in the input array and calls the `BuildBackendDTO` function to build the corresponding `Backend` DTO.
// The resulting DTOs are appended to the `result` array and returned.
func BuildBackendsDTOFromVO(backends []vo.Backend) (result []dto.Backend) {
	for _, backendVO := range backends {
		result = append(result, BuildBackendDTOFromVO(backendVO))
	}
	return result
}

// BuildBackendDTOFromVO builds a `Backend` DTO object using the provided `Backend` object as input.
// It retrieves various properties from the `Backend` object and sets them on the `Backend` DTO object.
func BuildBackendDTOFromVO(backendVO vo.Backend) dto.Backend {
	return dto.Backend{
		Name:           backendVO.Name(),
		Hosts:          backendVO.Hosts(),
		Path:           backendVO.Path(),
		Method:         backendVO.Method(),
		ForwardHeaders: backendVO.ForwardHeaders(),
		ForwardQueries: backendVO.ForwardQueries(),
		Modifiers:      BuildBackendModifiersDTOFromVO(backendVO.BackendModifiers()),
		ExtraConfig:    BuildBackendExtraConfigDTOFromVO(backendVO.ExtraConfig()),
	}
}

func BuildBackendExtraConfigDTOFromVO(backendExtraConfigVO *vo.BackendExtraConfig) *dto.BackendExtraConfig {
	if helper.IsNil(backendExtraConfigVO) || !backendExtraConfigVO.GroupResponse() ||
		!backendExtraConfigVO.OmitRequestBody() || !backendExtraConfigVO.OmitResponse() {
		return nil
	}
	return &dto.BackendExtraConfig{
		GroupResponse:   backendExtraConfigVO.GroupResponse(),
		OmitRequestBody: backendExtraConfigVO.OmitRequestBody(),
		OmitResponse:    backendExtraConfigVO.OmitResponse(),
	}
}

// BuildBackendModifiersDTOFromVO builds a `BackendModifiers` DTO object using the provided `BackendModifiers` object as input.
// It retrieves various properties from the `BackendModifiers` object and sets them on the `BackendModifiers` object.
func BuildBackendModifiersDTOFromVO(backendModifiersVO *vo.BackendModifiers) *dto.BackendModifiers {
	if helper.IsEmpty(backendModifiersVO) {
		return nil
	}
	return &dto.BackendModifiers{
		StatusCode: backendModifiersVO.StatusCode(),
		Header:     BuildModifiersDTOFromVO(backendModifiersVO.Header()),
		Param:      BuildModifiersDTOFromVO(backendModifiersVO.Param()),
		Query:      BuildModifiersDTOFromVO(backendModifiersVO.Query()),
		Body:       BuildModifiersDTOFromVO(backendModifiersVO.Body()),
	}
}

// BuildModifiersDTOFromVO builds a slice of `Modifier` DTO objects using the provided slice of `Modifier` objects as input.
// It iterates over each `Modifier` in the `modifiers` slice and calls `BuildModifierDTO` to create the DTO object.
// The DTO object is then appended to the `result` slice and finally returned.
func BuildModifiersDTOFromVO(modifiers []vo.Modifier) (result []dto.Modifier) {
	for _, modifierVO := range modifiers {
		result = append(result, *BuildModifierDTOFromVO(&modifierVO))
	}
	return result
}

// BuildModifierDTOFromVO builds a `Modifier` DTO object using the provided `Modifier` object as input.
// It checks if the input `Modifier` object is valid and returns `nil` if it is not.
// Otherwise, it creates a new `Modifier` object and sets its properties based on the input `Modifier` object.
// The created `Modifier` object is then returned as a pointer.
func BuildModifierDTOFromVO(modifierVO *vo.Modifier) *dto.Modifier {
	if helper.IsNil(modifierVO) {
		return nil
	}
	return &dto.Modifier{
		Context:   modifierVO.Context(),
		Scope:     modifierVO.Scope(),
		Action:    modifierVO.Action(),
		Propagate: modifierVO.Propagate(),
		Key:       modifierVO.Key(),
		Value:     modifierVO.Value(),
	}
}
