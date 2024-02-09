package usecase

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/dto"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/enum"
	"net/http"
	"strings"
)

type cors struct {
	securityCors dto.SecurityCors
}

type Cors interface {
	ValidateOrigins(ip string) error
	ValidateMethods(method string) error
	ValidateHeaders(header http.Header) error
}

func NewCors(securityCors dto.SecurityCors) Cors {
	return cors{
		securityCors: securityCors,
	}
}

func (c cors) ValidateOrigins(ip string) (err error) {
	if helper.IsNotEmpty(c.securityCors.AllowOrigins) &&
		helper.NotContains(c.securityCors.AllowOrigins, "*") &&
		helper.NotContains(c.securityCors.AllowOrigins, ip) {
		err = errors.New("Client IP origin not mapped on allow-origins")
	}
	return err
}

func (c cors) ValidateMethods(method string) (err error) {
	if helper.IsNotEmpty(c.securityCors.AllowMethods) &&
		helper.NotContains(c.securityCors.AllowMethods, "*") &&
		helper.NotContains(c.securityCors.AllowMethods, method) {
		err = errors.New("Request method not mapped on allow-methods")
	}
	return err
}

func (c cors) ValidateHeaders(header http.Header) (err error) {
	if helper.IsEmpty(c.securityCors.AllowHeaders) || helper.Contains(c.securityCors.AllowHeaders, "*") {
		return
	}
	var headerNotContains []string
	for key := range header {
		if helper.NotContains(c.securityCors.AllowHeaders, key) &&
			helper.IsNotEqualToIgnoreCase(key, enum.XForwardedFor) &&
			helper.IsNotEqualToIgnoreCase(key, enum.XTraceId) {
			headerNotContains = append(headerNotContains, key)
		}
	}
	if helper.IsNotEmpty(headerNotContains) {
		err = errors.New("Request headers contains not mapped fields on allow-headers:",
			strings.Join(headerNotContains, ", "))
	}
	return err
}
