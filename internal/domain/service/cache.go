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

package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/tech4works/checker"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
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

	key := c.buildKey(cache, request)

	cacheResponse, err := c.store.Get(ctx, key)
	if errors.Is(err, mapper.ErrCacheNotFound) {
		return nil, nil
	} else if checker.NonNil(err) {
		return nil, errors.Inheritf(err, "cache failed: unexpected error reading cache key=%s", key)
	}

	return cacheResponse, nil
}

func (c cacheService) Write(ctx context.Context, cache *vo.Cache, request *vo.HTTPRequest, response *vo.HTTPResponse) error {
	if !c.canWrite(cache, request, response) {
		return nil
	}

	key := c.buildKey(cache, request)

	err := c.store.Set(ctx, key, vo.NewCacheResponse(cache, response))
	if checker.NonNil(err) {
		return errors.Inheritf(err, "cache failed: unexpected error writing cache key=%s", key)
	}

	return nil
}

func (c cacheService) canRead(cache *vo.Cache, request *vo.HTTPRequest) bool {
	if cache.Disabled() {
		return false
	}

	return checker.NotEquals(enum.CacheControlNoCache, c.extractCacheControl(cache, request)) && c.allowMethod(cache, request)
}

func (c cacheService) canWrite(cache *vo.Cache, request *vo.HTTPRequest, response *vo.HTTPResponse) bool {
	if cache.Disabled() {
		return false
	}

	return checker.NotEquals(enum.CacheControlNoStore, c.extractCacheControl(cache, request)) &&
		c.allowMethod(cache, request) &&
		c.allowStatusCode(cache, response)
}

func (c cacheService) buildKey(cache *vo.Cache, request *vo.HTTPRequest) string {
	url := request.URL()
	if cache.IgnoreQuery() {
		url = request.Path().String()
	}

	strategyKey := fmt.Sprintf("%s:%s", request.Method(), url)

	var strategyHeaderValues []string
	for _, strategyHeaderKey := range cache.StrategyHeaders() {
		valueByStrategyKey := request.Header().Get(strategyHeaderKey)
		if checker.IsNotEmpty(valueByStrategyKey) {
			strategyHeaderValues = append(strategyHeaderValues, valueByStrategyKey)
		}
	}
	if checker.IsNotEmpty(strategyHeaderValues) {
		strategyKey = fmt.Sprintf("%s:%s", strategyKey, strings.Join(strategyHeaderValues, ":"))
	}

	return strategyKey
}

func (c cacheService) allowMethod(cache *vo.Cache, request *vo.HTTPRequest) bool {
	return !cache.HasOnlyIfMethods() || (!cache.HasAnyOnlyIfMethods() && checker.Equals(request.Method(), http.MethodGet)) ||
		checker.Contains(cache.OnlyIfMethods(), request.Method())
}

func (c cacheService) allowStatusCode(cache *vo.Cache, response *vo.HTTPResponse) bool {
	return (!cache.HasAnyOnlyIfStatusCodes() && response.StatusCode().OK()) ||
		(checker.NonNil(cache.OnlyIfStatusCodes()) && checker.Contains(cache.OnlyIfStatusCodes(), response.StatusCode().Code()))
}

func (c cacheService) extractCacheControl(cache *vo.Cache, request *vo.HTTPRequest) enum.CacheControl {
	var cacheControl enum.CacheControl
	if cache.AllowCacheControlNonNil() {
		cacheControl = enum.CacheControl(request.Header().GetFirst("Cache-Control"))
	}
	return cacheControl
}
