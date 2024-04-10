package vo

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"strings"
)

type SecurityCors struct {
	allowOrigins []string
	allowMethods []string
	allowHeaders []string
}

// newSecurityCors creates a new instance of SecurityCors based on the provided securityCorsDTO.
// It sets the allowOrigins, allowMethods, and allowHeaders fields of SecurityCors based on the values from securityCorsDTO.
func newSecurityCors(securityCorsDTO dto.SecurityCors) SecurityCors {
	return SecurityCors{
		allowOrigins: securityCorsDTO.AllowOrigins,
		allowMethods: securityCorsDTO.AllowMethods,
		allowHeaders: securityCorsDTO.AllowHeaders,
	}
}

// AllowOriginsData returns the allowOrigins field in the SecurityCors struct.
func (s SecurityCors) AllowOriginsData() []string {
	return s.allowOrigins
}

// AllowMethodsData returns the allowMethods field in the SecurityCors struct.
func (s SecurityCors) AllowMethodsData() []string {
	return s.allowMethods
}

// AllowHeadersData returns the allowHeaders field in the SecurityCors struct.
func (s SecurityCors) AllowHeadersData() []string {
	return s.allowHeaders
}

func (s SecurityCors) AllowOrigins(ip string) (err error) {
	// verificamos se na configuração security-cors.allow-origins ta vazia, ou tá informado o ip da requisição
	if helper.IsNotEmpty(s.allowOrigins) && helper.NotContains(s.allowOrigins, ip) {
		err = errors.New("Origin not mapped on security-cors.allow-origins")
	}
	return err
}

func (s SecurityCors) AllowMethods(method string) (err error) {
	// verificamos se na configuração security-cors.allow-methods ta vazia ou tá informado com oq vem na requisição
	if helper.IsNotEmpty(s.allowMethods) && helper.NotContains(s.allowMethods, method) {
		err = errors.New("Method not mapped on security-cors.allow-methods")
	}
	return err
}

func (s SecurityCors) AllowHeaders(header Header) (err error) {
	// verificamos se na configuração security-cors.allow-headers ta vazia
	if helper.IsEmpty(s.allowHeaders) {
		return nil
	}
	// inicializamos os headers não permitidos
	var headersNotAllowed []string
	// iteramos o header da requisição para verificar os headers que contain
	for key := range header {
		// caso o campo do header não esteja mapeado na lista security-cors.allow-headers e nao seja X-Forwarded-For
		// e nem X-Trace-Id adicionamos na lista
		if helper.NotContains(s.allowHeaders, key) && helper.IsNotEqualToIgnoreCase(key, consts.XForwardedFor) &&
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
