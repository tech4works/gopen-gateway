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
	Comment      string        `json:"@comment,omitempty"`
	Version      string        `json:"version,omitempty"`
	HotReload    bool          `json:"hot-reload,omitempty"`
	Proxy        *Proxy        `json:"proxy,omitempty"`
	Store        *Store        `json:"store,omitempty"`
	Timeout      vo.Duration   `json:"timeout,omitempty"`
	Cache        *Cache        `json:"cache,omitempty"`
	Limiter      *Limiter      `json:"limiter,omitempty"`
	SecurityCors *SecurityCors `json:"security-cors,omitempty"`
	Templates    *Templates    `json:"templates,omitempty"`
	Endpoints    []Endpoint    `json:"endpoints,omitempty"`
}

type Proxy struct {
	Provider enum.ProxyProvider `json:"provider,omitempty"`
	Token    string             `json:"token,omitempty"`
	Domains  []string           `json:"domains,omitempty"`
}

type Store struct {
	Redis *Redis `json:"redis,omitempty"`
}

type Redis struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
}

type Cache struct {
	Duration          vo.Duration `json:"duration,omitempty"`
	StrategyHeaders   []string    `json:"strategy-headers,omitempty"`
	OnlyIfStatusCodes []int       `json:"only-if-status-codes,omitempty"`
	OnlyIfMethods     []string    `json:"only-if-methods,omitempty"`
	AllowCacheControl *bool       `json:"allow-cache-control,omitempty"`
}

type EndpointCache struct {
	Enabled           bool        `json:"enabled"`
	IgnoreQuery       bool        `json:"ignore-query,omitempty"`
	Duration          vo.Duration `json:"duration,omitempty"`
	StrategyHeaders   []string    `json:"strategy-headers,omitempty"`
	OnlyIfStatusCodes []int       `json:"only-if-status-codes,omitempty"`
	AllowCacheControl *bool       `json:"allow-cache-control,omitempty"`
}

type Limiter struct {
	MaxHeaderSize          *vo.Bytes `json:"max-header-size,omitempty"`
	MaxBodySize            *vo.Bytes `json:"max-body-size,omitempty"`
	MaxMultipartMemorySize *vo.Bytes `json:"max-multipart-memory-size,omitempty"`
	Rate                   *Rate     `json:"rate,omitempty"`
}

type Rate struct {
	Capacity *int         `json:"capacity,omitempty"`
	Every    *vo.Duration `json:"every,omitempty"`
}

type EndpointLimiter struct {
	MaxHeaderSize          vo.Bytes `json:"max-header-size,omitempty"`
	MaxBodySize            vo.Bytes `json:"max-body-size,omitempty"`
	MaxMultipartMemorySize vo.Bytes `json:"max-multipart-memory-size,omitempty"`
	Rate                   *Rate    `json:"rate,omitempty"`
}

type SecurityCors struct {
	AllowOrigins []string `json:"allow-origins"`
	AllowMethods []string `json:"allow-methods"`
	AllowHeaders []string `json:"allow-headers"`
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

type Endpoint struct {
	Comment                string            `json:"@comment,omitempty"`
	Path                   string            `json:"path,omitempty"`
	Method                 string            `json:"method,omitempty"`
	Timeout                vo.Duration       `json:"timeout,omitempty"`
	Limiter                *EndpointLimiter  `json:"limiter,omitempty"`
	Cache                  *EndpointCache    `json:"cache,omitempty"`
	AbortIfHTPPStatusCodes *[]int            `json:"abort-if-http-status-codes,omitempty"`
	Parallelism            bool              `json:"parallelism,omitempty"`
	Beforewares            []Backend         `json:"beforewares,omitempty"`
	Backends               []Backend         `json:"backends,omitempty"`
	Afterwares             []Backend         `json:"afterwares,omitempty"`
	Response               *EndpointResponse `json:"response,omitempty"`
}

type EndpointResponse struct {
	Comment         string                  `json:"@comment,omitempty"`
	ContinueOnError bool                    `json:"continue-on-error,omitempty"`
	Header          *EndpointResponseHeader `json:"header,omitempty"`
	Body            *EndpointResponseBody   `json:"body,omitempty"`
}

type EndpointResponseHeader struct {
	Comment   string     `json:"@comment,omitempty"`
	Mapper    *Mapper    `json:"mapper,omitempty"`
	Projector *Projector `json:"projector,omitempty"`
}

type EndpointResponseBody struct {
	Comment         string               `json:"@comment,omitempty"`
	Aggregate       bool                 `json:"aggregate,omitempty"`
	OmitEmpty       bool                 `json:"omit-empty,omitempty"`
	ContentType     enum.ContentType     `json:"content-type,omitempty"`
	ContentEncoding enum.ContentEncoding `json:"content-encoding,omitempty"`
	Nomenclature    enum.Nomenclature    `json:"nomenclature,omitempty"`
	Mapper          *Mapper              `json:"mapper,omitempty"`
	Projector       *Projector           `json:"projector,omitempty"`
}

type Backend struct {
	Comment string `json:"@comment,omitempty"`

	// id do backend dentro do endpoint (usado para dependencies, debug, etc).
	// - Em templates: será preenchido automaticamente com a chave do template (template.path).
	// - Em endpoint (backend bruto): se omitido, será preenchido automaticamente com o path (backend.path).
	ID string `json:"id,omitempty"`

	Dependencies []string `json:"dependencies,omitempty"`

	OnlyIf   []string `json:"only-if,omitempty"`
	IgnoreIf []string `json:"ignore-if,omitempty"`

	Kind     enum.BackendKind `json:"kind,omitempty"`
	Template *Template        `json:"template,omitempty"`

	Async *bool `json:"async,omitempty"`

	// ---- HTTP ----
	Hosts     []string          `json:"hosts,omitempty"`
	Path      string            `json:"path,omitempty"`
	Method    string            `json:"method,omitempty"`
	Request   *BackendRequest   `json:"request,omitempty"`
	Response  *BackendResponse  `json:"response,omitempty"`
	Propagate *BackendPropagate `json:"propagate,omitempty"`

	// ---- PUBLISHER ----
	Provider        enum.PublisherProvider `json:"provider,omitempty"`
	GroupID         string                 `json:"group-id,omitempty"`
	DeduplicationID string                 `json:"deduplication-id,omitempty"`
	Delay           vo.Duration            `json:"delay,omitempty"`
	Message         *PublisherMessage      `json:"message,omitempty"`
}

type BackendRequest struct {
	Comment         string                `json:"@comment,omitempty"`
	ContinueOnError *bool                 `json:"continue-on-error,omitempty"`
	Concurrent      int                   `json:"concurrent,omitempty"`
	Async           *bool                 `json:"async,omitempty"`
	Header          *BackendRequestHeader `json:"header,omitempty"`
	Param           *BackendRequestParam  `json:"param,omitempty"`
	Query           *BackendRequestQuery  `json:"query,omitempty"`
	Body            *BackendRequestBody   `json:"body,omitempty"`
}

type BackendPropagate struct {
	Comment         string                `json:"@comment,omitempty"`
	ContinueOnError *bool                 `json:"continue-on-error,omitempty"`
	Header          *BackendRequestHeader `json:"header,omitempty"`
	Param           *BackendRequestParam  `json:"param,omitempty"`
	Query           *BackendRequestQuery  `json:"query,omitempty"`
	Body            *BackendRequestBody   `json:"body,omitempty"`
}

type BackendPropagateRequest struct {
	Header *BackendRequestHeader `json:"header,omitempty"`
	Param  *BackendRequestParam  `json:"param,omitempty"`
	Query  *BackendRequestQuery  `json:"query,omitempty"`
	Body   *BackendRequestBody   `json:"body,omitempty"`
}

type BackendRequestHeader struct {
	Comment   string     `json:"@comment,omitempty"`
	Omit      bool       `json:"omit,omitempty"`
	Mapper    *Mapper    `json:"mapper,omitempty"`
	Projector *Projector `json:"projector,omitempty"`
	Modifiers []Modifier `json:"modifiers,omitempty"`
}

type BackendRequestParam struct {
	Comment   string     `json:"@comment,omitempty"`
	Modifiers []Modifier `json:"modifiers,omitempty"`
}

type BackendRequestQuery struct {
	Comment   string     `json:"@comment,omitempty"`
	Omit      bool       `json:"omit,omitempty"`
	Mapper    *Mapper    `json:"mapper,omitempty"`
	Projector *Projector `json:"projector,omitempty"`
	Modifiers []Modifier `json:"modifiers,omitempty"`
}

type BackendRequestBody struct {
	Comment         string               `json:"@comment,omitempty"`
	Omit            bool                 `json:"omit,omitempty"`
	OmitEmpty       bool                 `json:"omit-empty,omitempty"`
	ContentType     enum.ContentType     `json:"content-type,omitempty"`
	ContentEncoding enum.ContentEncoding `json:"content-encoding,omitempty"`
	Nomenclature    enum.Nomenclature    `json:"nomenclature,omitempty"`
	Mapper          *Mapper              `json:"mapper,omitempty"`
	Projector       *Projector           `json:"projector,omitempty"`
	Modifiers       []Modifier           `json:"modifiers,omitempty"`
}

type BackendResponse struct {
	Comment         string                 `json:"@comment,omitempty"`
	ContinueOnError bool                   `json:"continue-on-error,omitempty"`
	Omit            bool                   `json:"omit,omitempty"`
	Header          *BackendResponseHeader `json:"header,omitempty"`
	Body            *BackendResponseBody   `json:"body,omitempty"`
}

type BackendResponseHeader struct {
	Comment   string     `json:"@comment,omitempty"`
	Omit      bool       `json:"omit,omitempty"`
	Mapper    *Mapper    `json:"mapper,omitempty"`
	Projector *Projector `json:"projector,omitempty"`
	Modifiers []Modifier `json:"modifiers,omitempty"`
}

type BackendResponseBody struct {
	Comment   string     `json:"@comment,omitempty"`
	Omit      bool       `json:"omit,omitempty"`
	Group     string     `json:"group,omitempty"`
	Mapper    *Mapper    `json:"mapper,omitempty"`
	Projector *Projector `json:"projector,omitempty"`
	Modifiers []Modifier `json:"modifiers,omitempty"`
}

type Publisher struct {
	Comment         string                 `json:"@comment,omitempty"`
	OnlyIf          []string               `json:"only-if,omitempty"`
	IgnoreIf        []string               `json:"ignore-if,omitempty"`
	Provider        enum.PublisherProvider `json:"provider,omitempty"`
	Reference       string                 `json:"reference,omitempty"`
	GroupID         string                 `json:"group-id,omitempty"`
	DeduplicationID string                 `json:"deduplication-id,omitempty"`
	Delay           vo.Duration            `json:"delay,omitempty"`
	Async           *bool                  `json:"async,omitempty"`
	Message         *PublisherMessage      `json:"message,omitempty"`
}

type PublisherMessage struct {
	Comment         string                               `json:"@comment,omitempty"`
	ContinueOnError bool                                 `json:"continue-on-error,omitempty"`
	OnlyIf          []string                             `json:"only-if,omitempty"`
	IgnoreIf        []string                             `json:"ignore-if,omitempty"`
	Attributes      map[string]PublisherMessageAttribute `json:"attributes,omitempty"`
	Body            *PublisherMessageBody                `json:"body,omitempty"`
}

type PublisherMessageBody struct {
	OmitEmpty bool       `json:"omit-empty,omitempty"`
	Mapper    *Mapper    `json:"mapper,omitempty"`
	Projector *Projector `json:"projector,omitempty"`
	Modifiers []Modifier `json:"modifiers,omitempty"`
}

type PublisherMessageAttribute struct {
	DataType string `json:"data-type,omitempty"`
	Value    string `json:"value,omitempty"`
}

type Mapper struct {
	Comment  string   `json:"@comment,omitempty"`
	OnlyIf   []string `json:"only-if,omitempty"`
	IgnoreIf []string `json:"ignore-if,omitempty"`
	Map      vo.Map   `json:"map,omitempty"`
}

type Projector struct {
	Comment  string     `json:"@comment,omitempty"`
	OnlyIf   []string   `json:"only-if,omitempty"`
	IgnoreIf []string   `json:"ignore-if,omitempty"`
	Project  vo.Project `json:"project,omitempty"`
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

type ErrorBody struct {
	ID        string    `json:"id,omitempty"`
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Endpoint  string    `json:"endpoint"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
