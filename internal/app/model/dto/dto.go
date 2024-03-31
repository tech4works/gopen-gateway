package dto

import (
	"bytes"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/gin-gonic/gin"
)

type GOpen struct {
	Version      string             `json:"version,omitempty"`
	Port         int                `json:"port,omitempty" validate:"required,gte=1"`
	HotReload    bool               `json:"hot-reload,omitempty"`
	Timeout      string             `json:"timeout,omitempty" validate:"omitempty,duration"`
	Store        *Store             `json:"store,omitempty"`
	Limiter      *Limiter           `json:"limiter,omitempty"`
	Cache        *Cache             `json:"cache,omitempty"`
	SecurityCors *SecurityCors      `json:"security-cors,omitempty"`
	Middlewares  map[string]Backend `json:"middlewares,omitempty"`
	Endpoints    []Endpoint         `json:"endpoints,omitempty" validate:"required"`
}

type GOpenView struct {
	Version      string             `json:"version,omitempty"`
	Port         int                `json:"port,omitempty"`
	HotReload    bool               `json:"hot-reload,omitempty"`
	Timeout      string             `json:"timeout,omitempty"`
	Limiter      *Limiter           `json:"limiter,omitempty"`
	Cache        *Cache             `json:"cache,omitempty"`
	SecurityCors *SecurityCors      `json:"security-cors,omitempty"`
	Middlewares  map[string]Backend `json:"middlewares,omitempty"`
	Endpoints    []Endpoint         `json:"endpoints,omitempty"`
}

type Store struct {
	Redis Redis `json:"redis,omitempty" validate:"required"`
}

type Redis struct {
	Address  string `json:"address,omitempty" validate:"required,url"`
	Password string `json:"password,omitempty" validate:"required"`
}

type Cache struct {
	Duration          string   `json:"duration,omitempty" validate:"required,duration"`
	StrategyHeaders   []string `json:"strategyHeaders,omitempty"`
	AllowCacheControl *bool    `json:"allowCacheControl,omitempty"`
}

type Limiter struct {
	MaxHeaderSize          string `json:"max-header-size,omitempty" validate:"omitempty,byte_unit"`
	MaxBodySize            string `json:"max-body-size,omitempty" validate:"omitempty,byte_unit"`
	MaxMultipartMemorySize string `json:"max-multipart-memory-size,omitempty" validate:"omitempty,byte_unit"`
	Rate                   *Rate  `json:"rate,omitempty"`
}

type Rate struct {
	Capacity int    `json:"capacity,omitempty" validate:"omitempty,gte=1"`
	Every    string `json:"every,omitempty" validate:"omitempty,duration"`
}

type SecurityCors struct {
	AllowOrigins []string `json:"allow-origins,omitempty"`
	AllowMethods []string `json:"allow-methods,omitempty"`
	AllowHeaders []string `json:"allow-headers,omitempty"`
}

type Endpoint struct {
	Path               string              `json:"path,omitempty" validate:"required,url_path"`
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
	Host           []string            `json:"host,omitempty" validate:"required,url"`
	Path           string              `json:"path,omitempty" validate:"required,url_path"`
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
	Context enum.ModifierContext `json:"context,omitempty" validate:"required,enum"`
	Scope   enum.ModifierScope   `json:"scope,omitempty" validate:"omitempty,enum"`
	Action  enum.ModifierAction  `json:"action,omitempty" validate:"required,enum"`
	Global  bool                 `json:"global,omitempty"`
	Key     string               `json:"key,omitempty" validate:"required"`
	Value   string               `json:"value,omitempty" validate:"required"`
}

type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (r ResponseWriter) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r ResponseWriter) WriteString(s string) (n int, err error) {
	r.Body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}
