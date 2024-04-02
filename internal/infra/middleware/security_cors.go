package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

type securityCors struct {
	securityCorsVO vo.SecurityCors
}

type SecurityCors interface {
	Do(req *api.Request)
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
func (c securityCors) Do(req *api.Request) {
	// chamamos o objeto de valor para validar se o ip de origem é permitida a partir do objeto de valor fornecido
	if err := c.securityCorsVO.AllowOrigins(req.HeaderValue(consts.XForwardedFor)); helper.IsNotNil(err) {
		req.WriteError(http.StatusForbidden, err)
		return
	}
	// chamamos o objeto de valor para validar se o method é permitida a partir do objeto de valor fornecido
	if err := c.securityCorsVO.AllowMethods(req.Method()); helper.IsNotNil(err) {
		req.WriteError(http.StatusForbidden, err)
		return
	}
	// chamamos o domínio para validar se o headers fornecido estão permitidas a partir do objeto de valor fornecido
	if err := c.securityCorsVO.AllowHeaders(req.Header()); helper.IsNotNil(err) {
		req.WriteError(http.StatusForbidden, err)
		return
	}

	// se tudo ocorreu bem seguimos para o próximo
	req.Next()
}
