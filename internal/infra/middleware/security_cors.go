package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

// securityCors implements the SecurityCors interface.
// It represents the configuration options for Cross-Origin Resource Sharing (CORS) settings in Gopen.
type securityCors struct {
	securityCorsVO vo.SecurityCors
}

// SecurityCors is an interface that defines the behavior of handling Cross-Origin Resource Sharing (CORS) settings in Gopen.
// Implementing types must define the Do method, which takes a *api.Context as an argument.
// The Do method is responsible for handling CORS-related tasks, such as validating and processing CORS requests.
type SecurityCors interface {
	// Do handles Cross-Origin Resource Sharing (CORS) related tasks.
	// It takes a *api.Context as an argument and is responsible for validating and processing CORS requests.
	Do(ctx *api.Context)
}

// NewSecurityCors is a function that creates a new instance of SecurityCors with the given securityCorsVO configuration.
// It returns the new SecurityCors object.
func NewSecurityCors(securityCorsVO vo.SecurityCors) SecurityCors {
	return securityCors{
		securityCorsVO: securityCorsVO,
	}
}

// Do perform CORS validation for the incoming request.
// It checks if the origin IP, HTTP method, and headers are allowed based on the provided securityCorsVO configuration.
//
// If the origin IP is not allowed, it writes a forbidden error response and exits the method.
// If the HTTP method is not allowed, it writes a forbidden error response and exits the method.
// If the headers are not allowed, it writes a forbidden error response and exits the method.
//
// If all validations passed, it proceeds to the next middleware or endpoint handler.
func (c securityCors) Do(ctx *api.Context) {
	// chamamos o objeto de valor para validar se o ip de origem é permitida a partir do objeto de valor fornecido
	if err := c.securityCorsVO.AllowOrigins(ctx.HeaderValue(consts.XForwardedFor)); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}
	// chamamos o objeto de valor para validar se o method é permitida a partir do objeto de valor fornecido
	if err := c.securityCorsVO.AllowMethods(ctx.Method()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}
	// chamamos o domínio para validar se o headers fornecido estão permitidas a partir do objeto de valor fornecido
	if err := c.securityCorsVO.AllowHeaders(ctx.Header()); helper.IsNotNil(err) {
		ctx.WriteError(http.StatusForbidden, err)
		return
	}

	// se tudo ocorreu bem seguimos para o próximo
	ctx.Next()
}
