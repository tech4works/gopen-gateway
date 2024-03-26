package dto

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"time"
)

type Gopen struct {
	Port         int                `json:"port,omitempty" validate:"required"`
	Timeout      time.Duration      `json:"timeout,omitempty" validate:"required"`
	Cache        time.Duration      `json:"cache,omitempty"`
	Limiter      Limiter            `json:"limiter,omitempty"`
	SecurityCors SecurityCors       `json:"security-cors,omitempty"`
	Middlewares  map[string]Backend `json:"middlewares,omitempty"`
	Endpoints    []Endpoint         `json:"endpoints,omitempty" validate:"required"`
}

type Limiter struct {
	MaxHeaderSize          vo.Bytes `json:"max-header-size,omitempty"`
	MaxBodySize            vo.Bytes `json:"max-body-size,omitempty"`
	MaxMultipartMemorySize vo.Bytes `json:"max-multipart-memory-size,omitempty"`
	Rate                   Rate     `json:"rate,omitempty"`
}

type Rate struct {
	Capacity int           `json:"capacity,omitempty"`
	Every    time.Duration `json:"every,omitempty"`
}

type SecurityCors struct {
	AllowCountries []string `json:"allow-countries,omitempty"`
	AllowOrigins   []string `json:"allow-origins,omitempty"`
	AllowMethods   []string `json:"allow-methods,omitempty"`
	AllowHeaders   []string `json:"allow-headers,omitempty"`
}

type Endpoint struct {
	Path               string        `json:"path,omitempty" validate:"required,url"`
	Method             string        `json:"method,omitempty" validate:"required"`
	Cache              time.Duration `json:"cache,omitempty"`
	Timeout            time.Duration `json:"timeout,omitempty"`
	Limiter            Limiter       `json:"limiter,omitempty"`
	AggregateResponses bool          `json:"aggregate-responses,omitempty"`
	AbortIfErrorOccurs bool          `json:"abort-if-error-occurs,omitempty"`
	Beforeware         []string      `json:"beforeware,omitempty"`
	Afterware          []string      `json:"afterware,omitempty"`
	Backends           []Backend     `json:"backends,omitempty" validate:"required"`
}

type Backend struct {
	Host           []string           `json:"host,omitempty" validate:"required"`
	Path           string             `json:"path,omitempty" validate:"required,url"`
	Method         string             `json:"method,omitempty" validate:"required"`
	ForwardHeaders []string           `json:"forward-headers,omitempty"`
	ForwardQueries []string           `json:"forward-queries,omitempty"`
	Modifier       BackendModifier    `json:"modifiers,omitempty"`
	ExtraConfig    BackendExtraConfig `json:"extra-config,omitempty"`
}

type BackendModifier struct {
	Headers []Modifier `json:"headers,omitempty"`
	Params  []Modifier `json:"params,omitempty"`
	Queries []Modifier `json:"queries,omitempty"`
	Body    []Modifier `json:"body,omitempty"`
}

type BackendExtraConfig struct {
	ResponseGroupName string `json:"response-group-name,omitempty"`
	OmitRequestBody   bool   `json:"omit-request-body,omitempty"`
	OmitResponse      bool   `json:"omit-response,omitempty"`
}

type Modifier struct {
	Context enum.ModifierContext `json:"context,omitempty" validate:"required,enum"`
	Scope   enum.ModifierScope   `json:"scope,omitempty" validate:"required,enum"`
	Action  enum.ModifierAction  `json:"action,omitempty" validate:"required,enum"`
	Key     string               `json:"key,omitempty" validate:"required"`
	Value   string               `json:"value,omitempty" validate:"required"`
}

type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (e Endpoint) IsCurrentRequest(ctx *gin.Context) bool {
	return (helper.Equals(e.Path, ctx.Request.URL.Path) || helper.Equals(e.Path, ctx.FullPath())) &&
		helper.Equals(e.Method, ctx.Request.Method)
}
