package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/util"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type securityCors struct {
	securityCorsVO vo.SecurityCors
}

type SecurityCors interface {
	Do(endpointVO vo.Endpoint) gin.HandlerFunc
}

func NewSecurityCors(securityCorsVO vo.SecurityCors) SecurityCors {
	return securityCors{
		securityCorsVO: securityCorsVO,
	}
}

func (c securityCors) Do(endpointVO vo.Endpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// chamamos o objeto de valor para validar se o ip de origem é permitida a partir do objeto de valor fornecido
		if err := c.securityCorsVO.AllowOrigins(ctx.GetHeader(consts.XForwardedFor)); helper.IsNotNil(err) {
			util.RespondCodeWithError(ctx, endpointVO.ResponseEncode(), http.StatusForbidden, err)
			return
		}
		// chamamos o objeto de valor para validar se o method é permitida a partir do objeto de valor fornecido
		if err := c.securityCorsVO.AllowMethods(ctx.Request.Method); helper.IsNotNil(err) {
			util.RespondCodeWithError(ctx, endpointVO.ResponseEncode(), http.StatusForbidden, err)
			return
		}
		// chamamos o domínio para validar se o headers fornecido estão permitidas a partir do objeto de valor fornecido
		if err := c.securityCorsVO.AllowHeaders(ctx.Request.Header); helper.IsNotNil(err) {
			util.RespondCodeWithError(ctx, endpointVO.ResponseEncode(), http.StatusForbidden, err)
			return
		}

		// se tudo ocorreu bem seguimos para o próximo
		ctx.Next()
	}
}
