package service

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"strings"
)

type securityCorsService struct {
}

type SecurityCors interface {
	ValidateOrigin(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error
	ValidateMethod(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error
	ValidateHeaders(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error
}

func NewSecurityCors() SecurityCors {
	return securityCorsService{}
}

func (s securityCorsService) ValidateOrigin(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error {
	if !securityCors.AllowOrigin(request.ClientIP()) {
		return errors.New("Origin not mapped on security-cors.allow-origins")
	}
	return nil
}

func (s securityCorsService) ValidateMethod(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error {
	if !securityCors.AllowMethod(request.Method()) {
		return errors.New("Method not mapped on security-cors.allow-methods")
	}
	return nil
}

func (s securityCorsService) ValidateHeaders(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error {
	var headersNotAllowed []string
	for _, key := range request.Header().Keys() {
		if checker.NotEquals(key, mapper.XForwardedFor) && !securityCors.AllowHeader(key) {
			headersNotAllowed = append(headersNotAllowed, key)
		}
	}

	if checker.IsNotEmpty(headersNotAllowed) {
		keys := strings.Join(headersNotAllowed, ", ")
		return errors.Newf("Headers contain not mapped fields on security-cors.allow-headers: %s", keys)
	}

	return nil
}
