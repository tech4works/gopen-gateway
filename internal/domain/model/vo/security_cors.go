package vo

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"strings"
)

type SecurityCors struct {
	allowOrigins []string
	allowMethods []string
	allowHeaders []string
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

// AllowOrigins checks if the given IP is allowed by the security-cors.allow-origins configuration
// It returns an error if the configuration is not empty, does not contain "*", and does not contain the IP
func (s SecurityCors) AllowOrigins(ip string) (err error) {
	// verificamos se na configuração security-cors.allow-origins ta vazia, tem * ou tá informado o ip da requisição
	if helper.IsNotEmpty(s.allowOrigins) && helper.NotContains(s.allowOrigins, "*") &&
		helper.NotContains(s.allowOrigins, ip) {
		err = errors.New("Origin not mapped on security-cors.allow-origins")
	}
	return err
}

// AllowMethods verifies if the specified method is allowed based on the configuration in SecurityCors.
// The method checks if the allow-methods configuration is not empty, does not contain "*", and does not have the specified method.
// If the method is not allowed, it returns an error with the message "method not mapped on security-cors.allow-methods".
// It returns the error, if any, otherwise it returns nil.
func (s SecurityCors) AllowMethods(method string) (err error) {
	// verificamos se na configuração security-cors.allow-methods ta vazia, tem * ou tá informado com oq vem na requisição
	if helper.IsNotEmpty(s.allowMethods) && helper.NotContains(s.allowMethods, "*") &&
		helper.NotContains(s.allowMethods, method) {
		err = errors.New("Method not mapped on security-cors.allow-methods")
	}
	return err
}

// AllowHeaders checks if the headers in the provided Header map are allowed according to the configuration in the SecurityCors struct.
// It returns an error if any of the headers are not allowed, based on the SecurityCors.allowHeaders field.
// If the SecurityCors.allowHeaders field is empty or contains "*", all headers are considered allowed.
// Headers listed in SecurityCors.allowHeaders, as well as the X-Forwarded-For and X-Trace-Id headers, are always allowed.
// If there are headers that are not allowed, the function returns an error with a message indicating which headers are not allowed.
// If there are no headers that are not allowed, it returns nil.
func (s SecurityCors) AllowHeaders(header Header) (err error) {
	// verificamos se na configuração security-cors.allow-headers ta vazia, tem * para retornar ok
	if helper.IsEmpty(s.allowHeaders) || helper.Contains(s.allowHeaders, "*") {
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
