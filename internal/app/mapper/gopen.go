package mapper

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"os"
)

// BuildSettingViewDTO builds a `SettingView` DTO object using the provided `Gopen` object as input.
// It retrieves various properties from the `Gopen` object and sets them on the `SettingView` object.
func BuildSettingViewDTO(gopenVO vo.Gopen) dto.SettingView {
	return dto.SettingView{
		Version:     os.Getenv("VERSION"),
		VersionDate: os.Getenv("VERSION_DATE"),
		Founder:     os.Getenv("FOUNDER"),
		CodeHelpers: os.Getenv("CODE_HELPERS"),
		Endpoints:   gopenVO.CountEndpoints(),
		Middlewares: gopenVO.CountMiddlewares(),
		Backends:    gopenVO.CountBackends(),
		Modifiers:   gopenVO.CountModifiers(),
		Setting:     BuildGOpenViewDTO(gopenVO),
	}
}

// BuildGOpenViewDTO builds a `GOpenView` DTO object using the provided `Gopen` object as input.
// It retrieves various properties from the `Gopen` object and sets them on the `GOpenView` object.
func BuildGOpenViewDTO(gopenVO vo.Gopen) dto.GOpenView {
	return dto.GOpenView{
		Version:      gopenVO.Version(),
		Port:         gopenVO.Port(),
		HotReload:    gopenVO.HotReload(),
		Timeout:      gopenVO.Timeout().String(),
		Limiter:      BuildLimiterDTOFromGOpenVO(gopenVO),
		Cache:        BuildCacheDTOFromGOpenVO(gopenVO),
		SecurityCors: BuildSecurityCorsDTO(gopenVO),
		Middlewares:  BuildMiddlewaresDTO(gopenVO),
		Endpoints:    BuildEndpointsDTO(gopenVO.Endpoints()),
	}
}

// BuildLimiterDTOFromGOpenVO builds a `Limiter` DTO object using the provided `Gopen` object as input.
// It retrieves various properties from the `Gopen` object and sets them on the `Limiter` object.
func BuildLimiterDTOFromGOpenVO(gopenVO vo.Gopen) *dto.Limiter {
	maxHeaderSize := gopenVO.LimiterMaxHeaderSize()
	maxBodySize := gopenVO.LimiterMaxBodySize()
	maxMultipartMemorySize := gopenVO.LimiterMaxMultipartMemorySize()

	return &dto.Limiter{
		MaxHeaderSize:          maxHeaderSize.String(),
		MaxBodySize:            maxBodySize.String(),
		MaxMultipartMemorySize: maxMultipartMemorySize.String(),
		Rate: &dto.Rate{
			Capacity: gopenVO.LimiterRateCapacity(),
			Every:    gopenVO.LimiterRateEvery().String(),
		},
	}
}

// BuildLimiterDTOFromEndpointVO builds a `Limiter` DTO object using the provided `Endpoint` object as input.
// If the `Endpoint` object does not have a limiter, it returns nil.
// It retrieves various properties from the `Endpoint` object and sets them on the `Limiter` object.
// It also sets the `Rate` property by creating a new `Rate` object and setting its properties using the `Endpoint` object.
// The `MaxHeaderSize`, `MaxBodySize`, and `MaxMultipartMemorySize` properties are converted to string format using the `String` method.
func BuildLimiterDTOFromEndpointVO(endpointVO vo.Endpoint) *dto.Limiter {
	if !endpointVO.HasLimiter() {
		return nil
	}

	maxHeaderSize := endpointVO.LimiterMaxHeaderSize()
	maxBodySize := endpointVO.LimiterMaxBodySize()
	maxMultipartMemorySize := endpointVO.LimiterMaxMultipartMemorySize()

	return &dto.Limiter{
		MaxHeaderSize:          maxHeaderSize.String(),
		MaxBodySize:            maxBodySize.String(),
		MaxMultipartMemorySize: maxMultipartMemorySize.String(),
		Rate: &dto.Rate{
			Capacity: endpointVO.LimiterRateCapacity(),
			Every:    endpointVO.LimiterRateEvery().String(),
		},
	}
}

// BuildCacheDTOFromGOpenVO builds a `Cache` DTO object using the provided `Gopen` object as input.
// It retrieves the cache duration, cache strategy headers, and allow cache control properties from the `Gopen` object
// and sets them on the `Cache` object.
func BuildCacheDTOFromGOpenVO(gopenVO vo.Gopen) *dto.Cache {
	return &dto.Cache{
		Duration:          gopenVO.CacheDuration().String(),
		StrategyHeaders:   gopenVO.CacheStrategyHeaders(),
		AllowCacheControl: helper.ConvertToPointer(gopenVO.AllowCacheControl()),
	}
}

// BuildCacheDTOFromEndpointVO builds a `Cache` DTO object using the provided `Endpoint` object as input.
// It retrieves various properties from the `Endpoint` object and sets them on the `Cache` object.
// If the `Endpoint` object does not have any cache, it returns nil.
// It returns a pointer to the `Cache` object.
func BuildCacheDTOFromEndpointVO(endpointVO vo.Endpoint) *dto.Cache {
	if !endpointVO.HasCache() {
		return nil
	}
	return &dto.Cache{
		Duration:          endpointVO.CacheDuration().String(),
		StrategyHeaders:   endpointVO.CacheStrategyHeaders(),
		AllowCacheControl: helper.ConvertToPointer(endpointVO.AllowCacheControl()),
	}
}

// BuildSecurityCorsDTO builds a `SecurityCors` DTO object using the provided `Gopen` object as input.
// It checks if the `SecurityCors` object from `Gopen` is empty. If not, it retrieves the `allowOriginsData`,
// `allowMethodsData`, and `allowHeadersData` properties from the `SecurityCors` object and sets them on the
// `SecurityCors` DTO object.
func BuildSecurityCorsDTO(gopenVO vo.Gopen) *dto.SecurityCors {
	if helper.IsEmpty(gopenVO.SecurityCors()) {
		return nil
	}
	return &dto.SecurityCors{
		AllowOrigins: gopenVO.SecurityCors().AllowOriginsData(),
		AllowMethods: gopenVO.SecurityCors().AllowMethodsData(),
		AllowHeaders: gopenVO.SecurityCors().AllowHeadersData(),
	}
}

// BuildMiddlewaresDTO builds a map of `Backend` DTO objects using the provided `Gopen` object as input.
// It iterates through the `Middlewares` in the `Gopen` object and calls the `BuildBackendDTO` function to build the `Backend` DTO for each one.
// The resulting `Backend` DTO objects are then added to the `result` map with the key being the key of the `Middleware` in the `Gopen` object.
// The `result` map is returned as the output of the function.
func BuildMiddlewaresDTO(gopenVO vo.Gopen) map[string]dto.Backend {
	result := map[string]dto.Backend{}
	for key, backendVO := range gopenVO.Middlewares() {
		result[key] = BuildBackendDTO(backendVO)
	}
	return result
}

// BuildEndpointsDTO builds a slice of `Endpoint` DTO objects using the provided `[]Endpoint` as input.
// It iterates over each `Endpoint` object and calls `BuildEndpointDTO` to build the individual `Endpoint` DTO object.
// The resulting DTO objects are then appended to the `result` slice and returned.
func BuildEndpointsDTO(endpoints []vo.Endpoint) (result []dto.Endpoint) {
	for _, endpointVO := range endpoints {
		result = append(result, BuildEndpointDTO(endpointVO))
	}
	return result
}

// BuildEndpointDTO builds a `Endpoint` DTO object using the provided `Endpoint` object as input.
// It retrieves various properties from the `Endpoint` object and sets them on the `Endpoint` DTO object.
func BuildEndpointDTO(endpointVO vo.Endpoint) dto.Endpoint {
	timeout := ""
	if endpointVO.HasTimeout() {
		timeout = endpointVO.Timeout().String()
	}

	return dto.Endpoint{
		Path:               endpointVO.Path(),
		Method:             endpointVO.Method(),
		Timeout:            timeout,
		Limiter:            BuildLimiterDTOFromEndpointVO(endpointVO),
		Cache:              BuildCacheDTOFromEndpointVO(endpointVO),
		ResponseEncode:     endpointVO.ResponseEncode(),
		AggregateResponses: endpointVO.AggregateResponses(),
		AbortIfStatusCodes: endpointVO.AbortIfStatusCodes(),
		Beforeware:         endpointVO.Beforeware(),
		Afterware:          endpointVO.Afterware(),
		Backends:           BuildBackendsDTO(endpointVO.Backends()),
	}
}

// BuildBackendsDTO builds an array of `Backend` DTO objects using the provided array of `Backend` objects as input.
// It iterates over each `Backend` in the input array and calls the `BuildBackendDTO` function to build the corresponding `Backend` DTO.
// The resulting DTOs are appended to the `result` array and returned.
func BuildBackendsDTO(backends []vo.Backend) (result []dto.Backend) {
	for _, backendVO := range backends {
		result = append(result, BuildBackendDTO(backendVO))
	}
	return result
}

// BuildBackendDTO builds a `Backend` DTO object using the provided `Backend` object as input.
// It retrieves various properties from the `Backend` object and sets them on the `Backend` DTO object.
func BuildBackendDTO(backendVO vo.Backend) dto.Backend {
	return dto.Backend{
		Name:           backendVO.Name(),
		Host:           backendVO.Host(),
		Path:           backendVO.Path(),
		Method:         backendVO.Method(),
		ForwardHeaders: backendVO.ForwardHeaders(),
		ForwardQueries: backendVO.ForwardQueries(),
		Modifiers:      BuildBackendModifiersDTO(backendVO.BackendModifiers()),
		ExtraConfig:    nil,
	}
}

// BuildBackendModifiersDTO builds a `BackendModifiers` DTO object using the provided `BackendModifiers` object as input.
// It retrieves various properties from the `BackendModifiers` object and sets them on the `BackendModifiers` object.
func BuildBackendModifiersDTO(backendModifiersVO vo.BackendModifiers) *dto.BackendModifiers {
	if helper.IsEmpty(backendModifiersVO) {
		return nil
	}
	return &dto.BackendModifiers{
		StatusCode: BuildModifierDTO(backendModifiersVO.StatusCode()),
		Header:     BuildModifiersDTO(backendModifiersVO.Header()),
		Params:     BuildModifiersDTO(backendModifiersVO.Params()),
		Query:      BuildModifiersDTO(backendModifiersVO.Query()),
		Body:       BuildModifiersDTO(backendModifiersVO.Body()),
	}
}

// BuildModifiersDTO builds a slice of `Modifier` DTO objects using the provided slice of `Modifier` objects as input.
// It iterates over each `Modifier` in the `modifiers` slice and calls `BuildModifierDTO` to create the DTO object.
// The DTO object is then appended to the `result` slice and finally returned.
func BuildModifiersDTO(modifiers []vo.Modifier) (result []dto.Modifier) {
	for _, modifierVO := range modifiers {
		result = append(result, *BuildModifierDTO(modifierVO))
	}
	return result
}

// BuildModifierDTO builds a `Modifier` DTO object using the provided `Modifier` object as input.
// It checks if the input `Modifier` object is valid and returns `nil` if it is not.
// Otherwise, it creates a new `Modifier` object and sets its properties based on the input `Modifier` object.
// The created `Modifier` object is then returned as a pointer.
func BuildModifierDTO(modifierVO vo.Modifier) *dto.Modifier {
	if !modifierVO.Valid() {
		return nil
	}
	return &dto.Modifier{
		Context: modifierVO.Context(),
		Scope:   modifierVO.Scope(),
		Action:  modifierVO.Action(),
		Global:  modifierVO.Global(),
		Key:     modifierVO.Key(),
		Value:   modifierVO.Value(),
	}
}
