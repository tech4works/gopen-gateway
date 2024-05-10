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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
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

// Read retrieves the cache response for the given cache and HTTP request.
// If the cache is unable to read or if the cache response is not found, it returns nil.
// It initializes the cache response value before obtaining it from the cache store.
// The cache store is responsible for getting the cache response based on the strategy key derived from the HTTP request.
// If there is an error during the cache retrieval process, it returns the error.
// Otherwise, it returns the obtained cache response.
func (c cacheService) Read(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest) (*vo.CacheResponse, error) {
	if cache.CantRead(httpRequest) {
		return nil, nil
	}

	var cacheGzipBase64 string
	err := c.cacheStore.Get(ctx, cache.StrategyKey(httpRequest), &cacheGzipBase64)
	if errors.Is(err, mapper.ErrCacheNotFound) {
		return nil, nil
	} else if helper.IsNotNil(err) {
		return nil, err
	}

	var cacheResponse vo.CacheResponse
	err = helper.ConvertGzipBase64ToDest(cacheGzipBase64, &cacheResponse)
	if helper.IsNotNil(err) {
		return nil, err
	}

	return &cacheResponse, nil
}

// Write stores the cache response obtained from the HTTP response
// into the cache store. If the cache is unable to write, it returns nil.
// It retrieves the cache duration from the cache and creates a new
// cache response using the HTTP response and duration. Then, it calls
// the cache store's Set method to store the cache response based on the
// strategy key derived from the HTTP request. If there is an error during
// the cache storage process, it returns the error.
func (c cacheService) Write(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse) error {
	if cache.CantWrite(httpRequest, httpResponse) {
		return nil
	}

	duration := cache.Duration()
	cacheResponse := vo.NewCacheResponse(httpResponse, duration)

	cacheGzipBase64, err := helper.ConvertToGzipBase64(cacheResponse)
	if helper.IsNotNil(err) {
		return err
	}

	return c.cacheStore.Set(ctx, cache.StrategyKey(httpRequest), cacheGzipBase64, duration.Time())
}
