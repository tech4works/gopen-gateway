package service

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
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
	if !securityCors.IsValidOrigin(request.ClientIP()) {
		return errors.New("Origin not mapped on security-cors.allow-origins")
	}
	return nil
}

func (s securityCorsService) ValidateMethod(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error {
	if !securityCors.IsValidMethod(request.Method()) {
		return errors.New("Method not mapped on security-cors.allow-methods")
	}
	return nil
}

func (s securityCorsService) ValidateHeaders(securityCors *vo.SecurityCors, request *vo.HTTPRequest) error {
	var headersNotAllowed []string
	for _, key := range request.Header().Keys() {
		if helper.IsNotEqualTo(key, mapper.XForwardedFor) && !securityCors.IsValidHeader(key) {
			headersNotAllowed = append(headersNotAllowed, key)
		}
	}
	if helper.IsNotEmpty(headersNotAllowed) {
		keys := strings.Join(headersNotAllowed, ", ")
		return errors.Newf("Headers contain not mapped fields on security-cors.allow-headers: %s", keys)
	}
	return nil
}
