/*
 * Copyright 2024 Tech4Works
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dto

import (
	"time"

	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type Gopen struct {
	Comment      string          `json:"@comment,omitempty"`
	Version      string          `json:"version,omitempty"`
	HotReload    bool            `json:"hot-reload,omitempty"`
	Store        *Store          `json:"store,omitempty"`
	Timeout      vo.Duration     `json:"timeout,omitempty"`
	Execution    *GopenExecution `json:"execution,omitempty"`
	SecurityCors *SecurityCors   `json:"security-cors,omitempty"`
	Limiter      *Limiter        `json:"limiter,omitempty"`
	Cache        *Cache          `json:"cache,omitempty"`
	Request      *Request        `json:"request,omitempty"`
	Components   *Components     `json:"components,omitempty"`
	Templates    *Templates      `json:"templates,omitempty"`
	Endpoints    []Endpoint      `json:"endpoints,omitempty"`
}

type GopenExecution struct {
	Comment string             `json:"@comment,omitempty"`
	Mode    enum.ExecutionMode `json:"mode,omitempty"`
	On      []enum.ExecutionOn `json:"on,omitempty"`
}

type Proxy struct {
	Token   string   `json:"token,omitempty"`
	Domains []string `json:"domains,omitempty"`
}

type Request struct {
	Client *RequestClient `json:"client,omitempty"`
}

type RequestClient struct {
	RequestID        *RequestClientValue            `json:"request-id,omitempty"`
	Trace            *RequestClientValue            `json:"trace,omitempty"`
	IP               *RequestClientIP               `json:"ip,omitempty"`
	TransportHeaders *RequestClientTransportHeaders `json:"transport-headers,omitempty"`
}

// RequestClientTransportHeaders controls which groups of transport headers are injected.
type RequestClientTransportHeaders struct {
	Request  *RequestClientTransportHeadersRequest  `json:"request,omitempty"`
	Response *RequestClientTransportHeadersResponse `json:"response,omitempty"`
}

// RequestClientTransportHeadersRequest controls headers injected in the request to the backend.
type RequestClientTransportHeadersRequest struct {
	Degradation *bool `json:"degradation,omitempty"`
	Timeout     *bool `json:"timeout,omitempty"`
}

// RequestClientTransportHeadersResponse controls headers injected in the response to the client.
type RequestClientTransportHeadersResponse struct {
	Cache           *bool `json:"cache,omitempty"`
	ExecutionStatus *bool `json:"execution-status,omitempty"`
	Degradation     *bool `json:"degradation,omitempty"`
}

type RequestClientValue struct {
	Headers   []string                `json:"headers,omitempty"`
	Fallback  *bool                   `json:"fallback,omitempty"`
	Propagate *RequestClientPropagate `json:"propagate,omitempty"`
}

type RequestClientIP struct {
	Headers        []string                `json:"headers,omitempty"`
	TrustedProxies []string                `json:"trusted-proxies,omitempty"`
	Propagate      *RequestClientPropagate `json:"propagate,omitempty"`
}

type RequestClientPropagate struct {
	Request  string `json:"request,omitempty"`
	Response string `json:"response,omitempty"`
}

type Store struct {
	Redis *Redis `json:"redis,omitempty"`
}

type Redis struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
}

type Cache struct {
	OnlyIf   []string      `json:"only-if,omitempty"`
	IgnoreIf []string      `json:"ignore-if,omitempty"`
	Read     CacheDecision `json:"read,omitempty"`
	Write    CacheDecision `json:"write,omitempty"`
	Key      string        `json:"key,omitempty"`
	TTL      vo.Duration   `json:"ttl,omitempty"`
}

type CacheDecision struct {
	OnlyIf   []string `json:"only-if,omitempty"`
	IgnoreIf []string `json:"ignore-if,omitempty"`
}

type Limiter struct {
	Size *LimiterSize `json:"size,omitempty"`
	Rate *LimiterRate `json:"rate,omitempty"`
}

type LimiterSize struct {
	OnlyIf   []string `json:"only-if,omitempty"`
	IgnoreIf []string `json:"ignore-if,omitempty"`

	// ---- HTTP ----
	MaxHeader *vo.Bytes `json:"max-header,omitempty"`
	MaxBody   *vo.Bytes `json:"max-body,omitempty"`
}

type LimiterRate struct {
	OnlyIf   []string     `json:"only-if,omitempty"`
	IgnoreIf []string     `json:"ignore-if,omitempty"`
	Capacity *int         `json:"capacity,omitempty"`
	Every    *vo.Duration `json:"every,omitempty"`
}

type SecurityCors struct {
	OnlyIf           []string `json:"only-if,omitempty"`
	IgnoreIf         []string `json:"ignore-if,omitempty"`
	AllowOrigins     []string `json:"allow-origins"`
	AllowMethods     []string `json:"allow-methods"`
	AllowHeaders     []string `json:"allow-headers"`
	AllowCredentials bool     `json:"allow-credentials"`
}

type EndpointExecution struct {
	Comment     string             `json:"@comment,omitempty"`
	Parallelism bool               `json:"parallelism,omitempty"`
	Mode        enum.ExecutionMode `json:"mode,omitempty"`
	On          []enum.ExecutionOn `json:"on,omitempty"`
}

type Templates struct {
	Comment     string             `json:"@comment,omitempty"`
	Beforewares map[string]Backend `json:"beforewares,omitempty"`
	Backends    map[string]Backend `json:"backends,omitempty"`
	Afterwares  map[string]Backend `json:"afterwares,omitempty"`
}

type Template struct {
	Path  string             `json:"path,omitempty"`
	Merge enum.TemplateMerge `json:"merge,omitempty"`
}

type ComponentReference struct {
	Path string `json:"path"`
}

type Components struct {
	Backends *ComponentsBackends `json:"backends,omitempty"`
}

type ComponentsBackends struct {
	Request   map[string]BackendRequest   `json:"request,omitempty"`
	Propagate map[string]BackendPropagate `json:"propagate,omitempty"`
}

type Endpoint struct {
	Comment      string             `json:"@comment,omitempty"`
	Execution    *EndpointExecution `json:"execution,omitempty"`
	Path         string             `json:"path,omitempty"`
	Method       string             `json:"method,omitempty"`
	Timeout      vo.Duration        `json:"timeout,omitempty"`
	SecurityCors *SecurityCors      `json:"security-cors,omitempty"`
	Limiter      *Limiter           `json:"limiter,omitempty"`
	Cache        *Cache             `json:"cache,omitempty"`
	Beforewares  []Backend          `json:"beforewares,omitempty"`
	Backends     []Backend          `json:"backends,omitempty"`
	Afterwares   []Backend          `json:"afterwares,omitempty"`
	Response     *EndpointResponse  `json:"response,omitempty"`
}

type EndpointResponse struct {
	Comment string `json:"@comment,omitempty"`

	// ---- HTTP ----
	Header *MetadataTransformation `json:"header,omitempty"`
	Body   *PayloadTransformation  `json:"body,omitempty"`
}

type Backend struct {
	Comment string `json:"@comment,omitempty"`

	// id do backend dentro do endpoint (usado para dependencies, debug, etc).
	// - Em templates: será preenchido automaticamente com a chave do template (template.path).
	// - Em endpoint (backend bruto): se omitido, será preenchido automaticamente com o path (backend.path).
	ID string `json:"id,omitempty"`

	Execution    *BackendExecution `json:"execution,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`

	OnlyIf   []string `json:"only-if,omitempty"`
	IgnoreIf []string `json:"ignore-if,omitempty"`

	Kind     enum.BackendKind `json:"kind,omitempty"`
	Template *Template        `json:"template,omitempty"`

	// ---- HTTP ----
	Cache *Cache `json:"cache,omitempty"`

	Hosts     []string         `json:"hosts,omitempty"`
	Path      string           `json:"path,omitempty"`
	Method    string           `json:"method,omitempty"`
	Request   BackendRequest   `json:"request,omitempty"`
	Propagate BackendPropagate `json:"propagate,omitempty"`

	// ---- PUBLISHER ----
	Broker enum.BackendBroker `json:"broker,omitempty"`

	GroupID         string           `json:"group-id,omitempty"`
	DeduplicationID string           `json:"deduplication-id,omitempty"`
	Delay           vo.Duration      `json:"delay,omitempty"`
	Message         PublisherMessage `json:"message,omitempty"`

	Response BackendResponse `json:"response,omitempty"`
}

type BackendRequest struct {
	Comment    string                    `json:"@comment,omitempty"`
	Components *BackendRequestComponents `json:"components,omitempty"`

	// ---- HTTP ----
	Header *MetadataTransformation `json:"header,omitempty"`
	Param  *URLPathTransformation  `json:"param,omitempty"`
	Query  *QueryTransformation    `json:"query,omitempty"`
	Body   *PayloadTransformation  `json:"body,omitempty"`
}

type BackendRequestComponents struct {
	Path  []ComponentReference `json:"path,omitempty"`
	Merge enum.ComponentMerge  `json:"merge,omitempty"`
}

type BackendPropagate struct {
	Comment    string                      `json:"@comment,omitempty"`
	Components *BackendPropagateComponents `json:"components,omitempty"`

	// ---- HTTP ----
	Header  *MetadataTransformation `json:"header,omitempty"`
	URLPath *URLPathTransformation  `json:"url-path,omitempty"`
	Query   *QueryTransformation    `json:"query,omitempty"`
	Body    *PayloadTransformation  `json:"body,omitempty"`
}

type BackendPropagateComponents struct {
	Path  []ComponentReference `json:"path,omitempty"`
	Merge enum.ComponentMerge  `json:"merge,omitempty"`
}

type MetadataTransformation struct {
	Comment   string     `json:"@comment,omitempty"`
	Omit      bool       `json:"omit,omitempty"`
	Mapper    *Mapper    `json:"mapper,omitempty"`
	Projector *Projector `json:"projector,omitempty"`
	Modifiers []Modifier `json:"modifiers,omitempty"`
}

type URLPathTransformation struct {
	Comment   string     `json:"@comment,omitempty"`
	Modifiers []Modifier `json:"modifiers,omitempty"`
}

type QueryTransformation struct {
	Comment   string     `json:"@comment,omitempty"`
	Omit      bool       `json:"omit,omitempty"`
	Mapper    *Mapper    `json:"mapper,omitempty"`
	Projector *Projector `json:"projector,omitempty"`
	Modifiers []Modifier `json:"modifiers,omitempty"`
}

type PayloadTransformation struct {
	Comment         string            `json:"@comment,omitempty"`
	Aggregate       bool              `json:"aggregate,omitempty"`
	Group           string            `json:"group,omitempty"`
	Omit            bool              `json:"omit,omitempty"`
	OmitEmpty       bool              `json:"omit-empty,omitempty"`
	ContentType     string            `json:"content-type,omitempty"`
	ContentEncoding string            `json:"content-encoding,omitempty"`
	Nomenclature    enum.Nomenclature `json:"nomenclature,omitempty"`
	Mapper          *Mapper           `json:"mapper,omitempty"`
	Projector       *Projector        `json:"projector,omitempty"`
	Modifiers       []Modifier        `json:"modifiers,omitempty"`
	Joins           []Join            `json:"joins,omitempty"`
}

type BackendResponse struct {
	Comment string `json:"@comment,omitempty"`
	Omit    bool   `json:"omit,omitempty"`

	// ---- HTTP ----
	Header *MetadataTransformation `json:"header,omitempty"`
	Body   *PayloadTransformation  `json:"body,omitempty"`
}

type PublisherMessage struct {
	Comment    string                    `json:"@comment,omitempty"`
	OnlyIf     []string                  `json:"only-if,omitempty"`
	IgnoreIf   []string                  `json:"ignore-if,omitempty"`
	Attributes map[string]AttributeValue `json:"attributes,omitempty"`
	Body       *PayloadTransformation    `json:"body,omitempty"`
}

type AttributeValue struct {
	Type  enum.AttributeValueType `json:"type,omitempty"`
	Value string                  `json:"value,omitempty"`
}

type Mapper struct {
	Comment  string            `json:"@comment,omitempty"`
	OnlyIf   []string          `json:"only-if,omitempty"`
	IgnoreIf []string          `json:"ignore-if,omitempty"`
	Policy   enum.MapperPolicy `json:"policy,omitempty"`
	Map      vo.MapConfig      `json:"map,omitempty"`
}

type Projector struct {
	Comment  string           `json:"@comment,omitempty"`
	OnlyIf   []string         `json:"only-if,omitempty"`
	IgnoreIf []string         `json:"ignore-if,omitempty"`
	Project  vo.ProjectConfig `json:"project,omitempty"`
}

type Modifier struct {
	Comment   string              `json:"@comment,omitempty"`
	OnlyIf    []string            `json:"only-if,omitempty"`
	IgnoreIf  []string            `json:"ignore-if,omitempty"`
	Action    enum.ModifierAction `json:"action,omitempty"`
	Propagate bool                `json:"propagate,omitempty"`
	Key       string              `json:"key,omitempty"`
	Value     string              `json:"value,omitempty"`
}

type Join struct {
	Comment  string     `json:"@comment,omitempty"`
	OnlyIf   []string   `json:"only-if,omitempty"`
	IgnoreIf []string   `json:"ignore-if,omitempty"`
	Source   JoinSource `json:"source,omitempty"`
	Target   JoinTarget `json:"target,omitempty"`
}

type JoinSource struct {
	Path string `json:"path,omitempty"`
	Key  string `json:"key,omitempty"`
}

type JoinTarget struct {
	Policy    enum.JoinTargetPolicy    `json:"policy,omitempty"`
	Path      string                   `json:"path,omitempty"`
	Key       string                   `json:"key,omitempty"`
	As        string                   `json:"as,omitempty"`
	OnMissing enum.JoinTargetOnMissing `json:"on-missing,omitempty"`
}

type BackendExecution struct {
	Comment    string             `json:"@comment,omitempty"`
	Concurrent int                `json:"concurrent,omitempty"`
	Async      *bool              `json:"async,omitempty"`
	Mode       enum.ExecutionMode `json:"mode,omitempty"`
	On         []enum.ExecutionOn `json:"on,omitempty"`
}

type ErrorPayload struct {
	ID        string    `json:"id,omitempty"`
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Endpoint  string    `json:"endpoint"`
	Message   string    `json:"message"`
	Stack     string    `json:"stack,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
