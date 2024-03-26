package service

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"net/http"
	"strings"
)

type securityCors struct {
	securityCorsVO vo.SecurityCors
}

type SecurityCors interface {
	ValidateOrigins(ip string) error
	ValidateMethods(method string) error
	ValidateHeaders(header http.Header) error
}

func NewSecurityCors(securityCorsVO vo.SecurityCors) SecurityCors {
	return securityCors{
		securityCorsVO: securityCorsVO,
	}
}

func (s securityCors) ValidateOrigins(ip string) (err error) {
	// verificamos se na configuração security-cors.allow-origins ta vazia, tem * ou tá informado o ip da requisição
	if helper.IsNotEmpty(s.securityCorsVO.AllowOrigins) && helper.NotContains(s.securityCorsVO.AllowOrigins, "*") &&
		helper.NotContains(s.securityCorsVO.AllowOrigins, ip) {
		err = errors.New("Origin not mapped on security-cors.allow-origins")
	}
	return err
}

func (s securityCors) ValidateMethods(method string) (err error) {
	// verificamos se na configuração security-cors.allow-methods ta vazia, tem * ou tá informado com oq vem na requisição
	if helper.IsNotEmpty(s.securityCorsVO.AllowMethods) && helper.NotContains(s.securityCorsVO.AllowMethods, "*") &&
		helper.NotContains(s.securityCorsVO.AllowMethods, method) {
		err = errors.New("method not mapped on security-cors.allow-methods")
	}
	return err
}

func (s securityCors) ValidateHeaders(header http.Header) (err error) {
	// verificamos se na configuração security-cors.allow-headers ta vazia, tem * para retornar ok
	if helper.IsEmpty(s.securityCorsVO.AllowHeaders) || helper.Contains(s.securityCorsVO.AllowHeaders, "*") {
		return nil
	}
	// inicializamos os headers não permitidos
	var headersNotAllowed []string
	// iteramos o header da requisição para verificar os headers que contain
	for key := range header {
		// caso o campo do header não esteja mapeado na lista security-cors.allow-headers e nao seja X-Forwarded-For
		// e nem X-Trace-Id adicionamos na lista
		if helper.NotContains(s.securityCorsVO.AllowHeaders, key) &&
			helper.IsNotEqualToIgnoreCase(key, consts.XForwardedFor) &&
			helper.IsNotEqualToIgnoreCase(key, consts.XTraceId) {
			headersNotAllowed = append(headersNotAllowed, key)
		}
	}
	// caso a lista não esteja vazia, quer dizer que tem headers não permitidos
	if helper.IsNotEmpty(headersNotAllowed) {
		headersFields := strings.Join(headersNotAllowed, ", ")
		return errors.New("Headers contains not mapped fields on security-cors.allow-headers:", headersFields)
	}

	// se tudo ocorreu bem retornamos nil
	return nil
}
