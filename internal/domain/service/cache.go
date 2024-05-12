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
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

// cacheService represents a cache service that interacts with a cache store.
type cacheService struct {
	// cacheStore represents an interface for interacting with a cache store.
	cacheStore interfaces.CacheStore
}

// Cache defines an interface for managing cache operations.
type Cache interface {
	// Read retrieves the cached response for a given HTTP request. It takes a context.Context object,
	// a vo.Cache pointer, and a vo.HttpRequest pointer as input parameters. It returns a pointer to
	// vo.CacheResponse and an error if an error occurred during the cache operation.
	Read(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest) (*vo.CacheResponse, error)
	// Write stores the response of an HTTP request in the cache. It takes a context.Context object,
	// a vo.Cache pointer, a vo.HttpRequest pointer, and a vo.HttpResponse pointer as input parameters.
	// It returns an error if an error occurred during the cache operation.
	Write(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest, httpResponse *vo.HttpResponse) error
}

// NewCache creates a new cache service instance. It takes a cacheStore object that implements the
// CacheStore interface as input parameter. It returns a Cache object.
func NewCache(cacheStore interfaces.CacheStore) Cache {
	return cacheService{
		cacheStore: cacheStore,
	}
}

func (c cacheService) Read(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest) (*vo.CacheResponse, error) {
	if cache.CantRead(httpRequest) {
		return nil, nil
	}

	cacheResponse, err := c.cacheStore.Get(ctx, cache.StrategyKey(httpRequest))
	if errors.Is(err, mapper.ErrCacheNotFound) {
		return nil, nil
	} else if helper.IsNotNil(err) {
		return nil, err
	}

	return cacheResponse, nil
}

func (c cacheService) Write(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse) error {
	if cache.CantWrite(httpRequest, httpResponse) {
		return nil
	}

	cacheResponse := vo.NewCacheResponse(httpResponse, cache.Duration())
	return c.cacheStore.Set(ctx, cache.StrategyKey(httpRequest), cacheResponse)
}
