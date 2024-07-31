/*
 * Copyright 2024 Gabriel Cataldo
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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"time"
)

type Gopen struct {
	Comment      string             `json:"@comment,omitempty"`
	Version      string             `json:"version,omitempty"`
	Port         int                `json:"port,omitempty"`
	HotReload    bool               `json:"hot-reload,omitempty"`
	Store        *Store             `json:"store,omitempty"`
	Timeout      vo.Duration        `json:"timeout,omitempty"`
	Cache        *Cache             `json:"cache,omitempty"`
	Limiter      *Limiter           `json:"limiter,omitempty"`
	SecurityCors *SecurityCors      `json:"security-cors,omitempty"`
	Middlewares  map[string]Backend `json:"middlewares,omitempty"`
	Endpoints    []Endpoint         `json:"endpoints,omitempty"`
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

type Middlewares map[string]Backend

type Endpoint struct {
	Comment            string            `json:"@comment,omitempty"`
	Path               string            `json:"path,omitempty"`
	Method             string            `json:"method,omitempty"`
	Timeout            vo.Duration       `json:"timeout,omitempty"`
	Limiter            *EndpointLimiter  `json:"limiter,omitempty"`
	Cache              *EndpointCache    `json:"cache,omitempty"`
	AbortIfStatusCodes *[]int            `json:"abort-if-status-codes,omitempty"`
	Response           *EndpointResponse `json:"response,omitempty"`
	Beforewares        []string          `json:"beforewares,omitempty"`
	Afterwares         []string          `json:"afterwares,omitempty"`
	Backends           []Backend         `json:"backends,omitempty"`
}

type EndpointResponse struct {
	Aggregate       bool                 `json:"aggregate,omitempty"`
	ContentType     enum.ContentType     `json:"content-type,omitempty"`
	ContentEncoding enum.ContentEncoding `json:"content-encoding,omitempty"`
	Nomenclature    enum.Nomenclature    `json:"nomenclature,omitempty"`
	OmitEmpty       bool                 `json:"omit-empty,omitempty"`
}

type Backend struct {
	Comment  string           `json:"@comment,omitempty"`
	Hosts    []string         `json:"hosts,omitempty"`
	Path     string           `json:"path,omitempty"`
	Method   string           `json:"method,omitempty"`
	Request  *BackendRequest  `json:"request,omitempty"`
	Response *BackendResponse `json:"response,omitempty"`
}

type BackendRequest struct {
	OmitHeader       bool                 `json:"omit-header,omitempty"`
	OmitQuery        bool                 `json:"omit-query,omitempty"`
	OmitBody         bool                 `json:"omit-body,omitempty"`
	ContentType      enum.ContentType     `json:"content-type,omitempty"`
	ContentEncoding  enum.ContentEncoding `json:"content-encoding,omitempty"`
	Nomenclature     enum.Nomenclature    `json:"nomenclature,omitempty"`
	OmitEmpty        bool                 `json:"omit-empty,omitempty"`
	HeaderMapper     *vo.Mapper           `json:"header-mapper,omitempty"`
	QueryMapper      *vo.Mapper           `json:"query-mapper,omitempty"`
	BodyMapper       *vo.Mapper           `json:"body-mapper,omitempty"`
	HeaderProjection *vo.Projection       `json:"header-projection,omitempty"`
	QueryProjection  *vo.Projection       `json:"query-projection,omitempty"`
	BodyProjection   *vo.Projection       `json:"body-projection,omitempty"`
	HeaderModifiers  []Modifier           `json:"header-modifiers,omitempty"`
	ParamModifiers   []Modifier           `json:"param-modifiers,omitempty"`
	QueryModifiers   []Modifier           `json:"query-modifiers,omitempty"`
	BodyModifiers    []Modifier           `json:"body-modifiers,omitempty"`
}

type BackendResponse struct {
	Omit             bool           `json:"omit,omitempty"`
	OmitHeader       bool           `json:"omit-header,omitempty"`
	OmitBody         bool           `json:"omit-body,omitempty"`
	Group            string         `json:"group,omitempty"`
	HeaderMapper     *vo.Mapper     `json:"header-mapper,omitempty"`
	BodyMapper       *vo.Mapper     `json:"body-mapper,omitempty"`
	HeaderProjection *vo.Projection `json:"header-projection,omitempty"`
	BodyProjection   *vo.Projection `json:"body-projection,omitempty"`
	HeaderModifiers  []Modifier     `json:"header-modifiers,omitempty"`
	BodyModifiers    []Modifier     `json:"body-modifiers,omitempty"`
}

type Modifier struct {
	Comment   string              `json:"@comment,omitempty"`
	Action    enum.ModifierAction `json:"action,omitempty"`
	Propagate bool                `json:"propagate,omitempty"`
	Key       string              `json:"key,omitempty"`
	Value     string              `json:"value,omitempty"`
}

type ErrorBody struct {
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Endpoint  string    `json:"endpoint"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
