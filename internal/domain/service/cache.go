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

package service

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"net/http"
	"strings"
)

type cacheService struct {
	store domain.Store
}

type Cache interface {
	Read(ctx context.Context, cache *vo.Cache, request *vo.HTTPRequest) (*vo.CacheResponse, error)
	Write(ctx context.Context, cache *vo.Cache, request *vo.HTTPRequest, response *vo.HTTPResponse) error
}

func NewCache(store domain.Store) Cache {
	return cacheService{
		store: store,
	}
}

func (c cacheService) Read(ctx context.Context, cache *vo.Cache, request *vo.HTTPRequest) (*vo.CacheResponse, error) {
	if !c.canRead(cache, request) {
		return nil, nil
	}

	cacheResponse, err := c.store.Get(ctx, c.buildKey(cache, request))
	if errors.Is(err, mapper.ErrCacheNotFound) {
		return nil, nil
	} else if helper.IsNotNil(err) {
		return nil, err
	}

	return cacheResponse, nil
}

func (c cacheService) Write(ctx context.Context, cache *vo.Cache, request *vo.HTTPRequest, response *vo.HTTPResponse) error {
	if !c.canWrite(cache, request, response) {
		return nil
	}

	return c.store.Set(ctx, c.buildKey(cache, request), vo.NewCacheResponse(cache, response))
}

func (c cacheService) canRead(cache *vo.Cache, request *vo.HTTPRequest) bool {
	if cache.Disabled() {
		return false
	}

	return helper.IsNotEqualTo(enum.CacheControlNoCache, c.extractCacheControl(cache, request)) &&
		c.allowMethod(cache, request)
}

func (c cacheService) canWrite(cache *vo.Cache, request *vo.HTTPRequest, response *vo.HTTPResponse) bool {
	if cache.Disabled() {
		return false
	}

	return helper.IsNotEqualTo(enum.CacheControlNoStore, c.extractCacheControl(cache, request)) &&
		c.allowMethod(cache, request) && c.allowStatusCode(cache, response)
}

func (c cacheService) buildKey(cache *vo.Cache, request *vo.HTTPRequest) string {
	url := request.Url()
	if cache.IgnoreQuery() {
		url = request.Path().String()
	}
	strategyKey := fmt.Sprintf("%s:%s", request.Method(), url)

	var strategyHeaderValues []string
	for _, strategyHeaderKey := range cache.StrategyHeaders() {
		valueByStrategyKey := request.Header().Get(strategyHeaderKey)
		if helper.IsNotEmpty(valueByStrategyKey) {
			strategyHeaderValues = append(strategyHeaderValues, valueByStrategyKey)
		}
	}
	if helper.IsNotEmpty(strategyHeaderValues) {
		strategyKey = fmt.Sprintf("%s:%s", strategyKey, strings.Join(strategyHeaderValues, ":"))
	}

	return strategyKey
}

func (c cacheService) allowMethod(cache *vo.Cache, request *vo.HTTPRequest) bool {
	return !cache.HasOnlyIfMethods() || (!cache.HasAnyOnlyIfMethods() && helper.Equals(request.Method(), http.MethodGet)) ||
		helper.Contains(cache.OnlyIfMethods(), request.Method())
}

func (c cacheService) allowStatusCode(cache *vo.Cache, response *vo.HTTPResponse) bool {
	statusCode := response.StatusCode()
	return !cache.HasOnlyIfStatusCodes() || (!cache.HasAnyOnlyIfStatusCodes() && statusCode.OK()) ||
		helper.Contains(cache.OnlyIfStatusCodes(), statusCode.Code())
}

func (c cacheService) extractCacheControl(cache *vo.Cache, request *vo.HTTPRequest) enum.CacheControl {
	var cacheControl enum.CacheControl
	if cache.AllowCacheControlNonNil() {
		cacheControl = enum.CacheControl(request.Header().GetFirst("Cache-Control"))
	}
	return cacheControl
}
