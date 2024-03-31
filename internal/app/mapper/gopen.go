package mapper

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"os"
)

func BuildSettingViewDTO(gopenVO vo.GOpen) dto.SettingView {
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

func BuildGOpenViewDTO(gopenVO vo.GOpen) dto.GOpenView {
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

func BuildLimiterDTOFromGOpenVO(gopenVO vo.GOpen) *dto.Limiter {
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

func BuildCacheDTOFromGOpenVO(gopenVO vo.GOpen) *dto.Cache {
	return &dto.Cache{
		Duration:          gopenVO.CacheDuration().String(),
		StrategyHeaders:   gopenVO.CacheStrategyHeaders(),
		AllowCacheControl: helper.ConvertToPointer(gopenVO.AllowCacheControl()),
	}
}

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

func BuildSecurityCorsDTO(gopenVO vo.GOpen) *dto.SecurityCors {
	if helper.IsEmpty(gopenVO.SecurityCors()) {
		return nil
	}
	return &dto.SecurityCors{
		AllowCountries: gopenVO.SecurityCors().AllowCountriesData(),
		AllowOrigins:   gopenVO.SecurityCors().AllowOriginsData(),
		AllowMethods:   gopenVO.SecurityCors().AllowMethodsData(),
		AllowHeaders:   gopenVO.SecurityCors().AllowHeadersData(),
	}
}

func BuildMiddlewaresDTO(gopenVO vo.GOpen) map[string]dto.Backend {
	result := map[string]dto.Backend{}
	for key, backendVO := range gopenVO.Middlewares() {
		result[key] = BuildBackendDTO(backendVO)
	}
	return result
}

func BuildEndpointsDTO(endpoints []vo.Endpoint) (result []dto.Endpoint) {
	for _, endpointVO := range endpoints {
		result = append(result, BuildEndpointDTO(endpointVO))
	}
	return result
}

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

func BuildBackendsDTO(backends []vo.Backend) (result []dto.Backend) {
	for _, backendVO := range backends {
		result = append(result, BuildBackendDTO(backendVO))
	}
	return result
}

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

func BuildModifiersDTO(modifiers []vo.Modifier) (result []dto.Modifier) {
	for _, modifierVO := range modifiers {
		result = append(result, *BuildModifierDTO(modifierVO))
	}
	return result
}

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
