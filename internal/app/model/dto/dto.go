package dto

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

type Gopen struct {
	Version      string             `json:"version,omitempty"`
	Port         int                `json:"port,omitempty"`
	HotReload    bool               `json:"hot-reload,omitempty"`
	Timeout      string             `json:"timeout,omitempty"`
	Store        *Store             `json:"store,omitempty"`
	Limiter      *Limiter           `json:"limiter,omitempty"`
	Cache        *Cache             `json:"cache,omitempty"`
	SecurityCors *SecurityCors      `json:"security-cors,omitempty"`
	Middlewares  map[string]Backend `json:"middlewares,omitempty"`
	Endpoints    []Endpoint         `json:"endpoints,omitempty"`
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
	Redis Redis `json:"redis,omitempty"`
}

type Redis struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
}

type Cache struct {
	Duration          string   `json:"duration,omitempty"`
	StrategyHeaders   []string `json:"strategy-headers,omitempty"`
	OnlyIfStatusCodes []int    `json:"only-if-status-codes,omitempty"`
	OnlyIfMethods     []string `json:"only-if-methods,omitempty"`
	AllowCacheControl *bool    `json:"allow-cache-control,omitempty"`
}

type EndpointCache struct {
	Enabled           bool     `json:"enabled,omitempty"`
	IgnoreQuery       bool     `json:"ignore-query,omitempty"`
	Duration          string   `json:"duration,omitempty"`
	StrategyHeaders   []string `json:"strategy-headers,omitempty"`
	OnlyIfStatusCodes []int    `json:"only-if-status-codes,omitempty"`
	AllowCacheControl *bool    `json:"allow-cache-control,omitempty"`
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
	AllowOrigins []string `json:"allow-origins,omitempty"`
	AllowMethods []string `json:"allow-methods,omitempty"`
	AllowHeaders []string `json:"allow-headers,omitempty"`
}

type Endpoint struct {
	Path               string              `json:"path,omitempty"`
	Method             string              `json:"method,omitempty"`
	Timeout            string              `json:"timeout,omitempty"`
	Limiter            *Limiter            `json:"limiter,omitempty"`
	Cache              *EndpointCache      `json:"cache,omitempty"`
	ResponseEncode     enum.ResponseEncode `json:"response-encode,omitempty"`
	AggregateResponses bool                `json:"aggregate-responses,omitempty"`
	AbortIfStatusCodes []int               `json:"abort-if-status-codes,omitempty"`
	Beforeware         []string            `json:"beforeware,omitempty"`
	Afterware          []string            `json:"afterware,omitempty"`
	Backends           []Backend           `json:"backends,omitempty"`
}

type Backend struct {
	Name           string              `json:"name,omitempty"`
	Hosts          []string            `json:"hosts,omitempty"`
	Path           string              `json:"path,omitempty"`
	Method         string              `json:"method,omitempty"`
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
	Context   enum.ModifierContext `json:"context,omitempty"`
	Scope     enum.ModifierScope   `json:"scope,omitempty"`
	Action    enum.ModifierAction  `json:"action,omitempty"`
	Propagate bool                 `json:"propagate,omitempty"`
	Key       string               `json:"key,omitempty"`
	Value     string               `json:"value,omitempty"`
}
