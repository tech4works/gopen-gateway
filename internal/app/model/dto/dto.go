package dto

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/gin-gonic/gin"
)

type GOpen struct {
	Version      string             `json:"version,omitempty"`
	Port         int                `json:"port,omitempty" validate:"required"`
	Timeout      string             `json:"timeout,omitempty"`
	Limiter      *Limiter           `json:"limiter,omitempty"`
	Cache        *Cache             `json:"cache,omitempty"`
	SecurityCors *SecurityCors      `json:"security-cors,omitempty"`
	Middlewares  map[string]Backend `json:"middlewares,omitempty"`
	Endpoints    []Endpoint         `json:"endpoints,omitempty" validate:"required"`
}

type Cache struct {
	Duration          string   `json:"duration,omitempty"`
	StrategyHeaders   []string `json:"strategyHeaders,omitempty"`
	AllowCacheControl *bool    `json:"allowCacheControl,omitempty"`
}

type Limiter struct {
	MaxHeaderSize          string `json:"max-header-size,omitempty"`
	MaxBodySize            string `json:"max-body-size,omitempty"`
	MaxMultipartMemorySize string `json:"max-multipart-memory-size,omitempty"`
	Rate                   *Rate  `json:"rate,omitempty"`
}

type Rate struct {
	Capacity int    `json:"capacity,omitempty"`
	Every    string `json:"every,omitempty"`
}

type SecurityCors struct {
	AllowCountries []string `json:"allow-countries,omitempty"`
	AllowOrigins   []string `json:"allow-origins,omitempty"`
	AllowMethods   []string `json:"allow-methods,omitempty"`
	AllowHeaders   []string `json:"allow-headers,omitempty"`
}

type Endpoint struct {
	Path               string              `json:"path,omitempty" validate:"required,url"`
	Method             string              `json:"method,omitempty" validate:"required"`
	Timeout            string              `json:"timeout,omitempty"`
	Limiter            *Limiter            `json:"limiter,omitempty"`
	Cache              *Cache              `json:"cache,omitempty"`
	ResponseEncode     enum.ResponseEncode `json:"response-encode,omitempty"`
	AggregateResponses bool                `json:"aggregate-responses,omitempty"`
	AbortIfStatusCodes []int               `json:"abort-if-status-codes,omitempty"`
	Beforeware         []string            `json:"beforeware,omitempty"`
	Afterware          []string            `json:"afterware,omitempty"`
	Backends           []Backend           `json:"backends,omitempty" validate:"required"`
}

type Backend struct {
	Name           string              `json:"name,omitempty"`
	Host           []string            `json:"host,omitempty" validate:"required"`
	Path           string              `json:"path,omitempty" validate:"required,url"`
	Method         string              `json:"method,omitempty" validate:"required"`
	ForwardHeaders []string            `json:"forward-headers,omitempty"`
	ForwardQueries []string            `json:"forward-queries,omitempty"`
	Modifiers      *BackendModifiers   `json:"modifiers,omitempty"`
	ExtraConfig    *BackendExtraConfig `json:"extra-config,omitempty"`
}

type BackendModifiers struct {
	StatusCode *Modifier  `json:"status-code,omitempty"`
	Header     []Modifier `json:"header,omitempty"`
	Params     []Modifier `json:"params,omitempty"`
	Query      []Modifier `json:"query,omitempty"`
	Body       []Modifier `json:"body,omitempty"`
}

type BackendExtraConfig struct {
	GroupResponse   bool `json:"group-response,omitempty"`
	OmitRequestBody bool `json:"omit-request-body,omitempty"`
	OmitResponse    bool `json:"omit-response,omitempty"`
}

type Modifier struct {
	Context enum.ModifierContext `json:"context,omitempty" validate:"omitempty,enum"`
	Scope   enum.ModifierScope   `json:"scope,omitempty" validate:"omitempty,enum"`
	Action  enum.ModifierAction  `json:"action,omitempty" validate:"omitempty,enum"`
	Global  bool                 `json:"global,omitempty"`
	Key     string               `json:"key,omitempty"`
	Value   string               `json:"value,omitempty" validate:"required"`
}

type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (g GOpen) CountMiddlewares() int {
	return len(g.Middlewares)
}

func (g GOpen) CountEndpoints() int {
	return len(g.Endpoints)
}

func (g GOpen) CountBackends() (count int) {
	count += g.CountMiddlewares()
	for _, endpointDTO := range g.Endpoints {
		count += endpointDTO.CountBackends()
	}
	return count
}

func (g GOpen) CountModifiers() (count int) {
	for _, middlewareBackend := range g.Middlewares {
		count += middlewareBackend.CountModifiers()
	}
	for _, endpointDTO := range g.Endpoints {
		count += endpointDTO.CountModifiers()
	}
	return count
}

func (e Endpoint) CountBackends() int {
	return len(e.Backends)
}

func (e Endpoint) CountModifiers() (count int) {
	for _, backendDTO := range e.Backends {
		count += backendDTO.CountModifiers()
	}
	return count
}

func (b Backend) CountModifiers() int {
	if helper.IsNotNil(b.Modifiers) {
		return b.Modifiers.CountAll()
	}
	return 0
}

func (m BackendModifiers) CountAll() (count int) {
	if helper.IsNotNil(m.StatusCode) {
		count++
	}
	count += len(m.Header) + len(m.Params) + len(m.Query) + len(m.Body)
	return count
}

func (r ResponseWriter) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r ResponseWriter) WriteString(s string) (n int, err error) {
	r.Body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}
