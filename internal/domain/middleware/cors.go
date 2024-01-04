package middleware

import (
	"github.com/GabrielHCataldo/go-error-detail/errors"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/dto"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/handler"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

type cors struct {
	securityCors  dto.SecurityCors
	localeUseCase service.Locale
}

type Cors interface {
	PreHandlerRequest(ctx *gin.Context)
}

func NewCors(securityCors dto.SecurityCors, localeUseCase service.Locale) Cors {
	return cors{
		securityCors:  securityCors,
		localeUseCase: localeUseCase,
	}
}

func (c cors) PreHandlerRequest(ctx *gin.Context) {
	logger.Info("Start request on URI:", ctx.Request.RequestURI)

	//todo -> depois pensar se vale a pena mesmo colocar esse cara, obs: deixa a request lenta e gera custo
	//ipLocale, err := c.localeUseCase.GetLocaleByIP(ctx, ctx.ClientIP())
	//if err != nil {
	//	handler.RespondCodeWithError(ctx, http.StatusInternalServerError, errors.NewByError(err))
	//	return
	//}
	//ctx.Set("Locale", ipLocale)
	//ctx.Request.Header.Set("X-Forwarded-For", ipLocale.IpAddress)

	////check allow-countries
	//if len(c.securityCors.AllowCountries) > 0 &&
	//	!slices.Contains(c.securityCors.AllowCountries, "*") &&
	//	!slices.Contains(c.securityCors.AllowCountries, ipLocale.CountryCode) {
	//	handler.RespondCodeWithError(ctx, http.StatusForbidden, errors.New(
	//		"Client country not mapped on allow-countries",
	//	))
	//	return
	//}
	//check allow-origins
	if len(c.securityCors.AllowOrigins) > 0 &&
		!slices.Contains(c.securityCors.AllowOrigins, "*") &&
		!slices.Contains(c.securityCors.AllowOrigins, ctx.ClientIP()) {
		handler.RespondCodeWithError(ctx, http.StatusForbidden, errors.New(
			"Client IP origin not mapped on allow-origins",
		))
		return
	}
	//check allow-methods
	if len(c.securityCors.AllowMethods) > 0 &&
		!slices.Contains(c.securityCors.AllowMethods, "*") &&
		!slices.Contains(c.securityCors.AllowMethods, ctx.Request.Method) {
		handler.RespondCodeWithError(ctx, http.StatusForbidden, errors.New(
			"Request method not mapped on allow-methods",
		))
		return
	}
	//check allow-headers
	if len(c.securityCors.AllowHeaders) > 0 && !slices.Contains(c.securityCors.AllowHeaders, "*") {
		containsNotAllowHeader := false
		for key := range ctx.Request.Header {
			if !slices.Contains(c.securityCors.AllowHeaders, key) {
				logger.Info("Allow-headers not contains:", key)
				containsNotAllowHeader = true
				break
			}
		}
		if containsNotAllowHeader {
			handler.RespondCodeWithError(ctx, http.StatusForbidden, errors.New(
				"Request headers contains not mapped field on allow-headers",
			))
			return
		}
	}
	logger.Info("Finish cors validate on URI:", ctx.Request.RequestURI)
}
