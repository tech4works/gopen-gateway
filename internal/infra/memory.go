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

package infra

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/jellydator/ttlcache/v2"
)

// memoryStore represents an in-memory cache store that implements the CacheStore interface.
type memoryStore struct {
	ttlCache *ttlcache.Cache
}

// NewMemoryStore returns a new instance of the MemoryStore structure that implements the CacheStore interface.
// This implementation uses an in-memory cache with a time-to-live (TTL)
func NewMemoryStore() interfaces.CacheStore {
	ttlCache := ttlcache.NewCache()
	ttlCache.SkipTTLExtensionOnHit(true)
	return &memoryStore{
		ttlCache: ttlCache,
	}
}

func (r memoryStore) Set(_ context.Context, key string, cacheResponse *vo.CacheResponse) error {
	gzipBase64, err := helper.CompressWithGzipToBase64(cacheResponse)
	if helper.IsNotNil(err) {
		return err
	}
	return r.ttlCache.SetWithTTL(key, gzipBase64, cacheResponse.Duration.Time())
}

// Del removes a key-value pair from the memory cache with the specified key.
// The key is a string that serves as the identifier for the key-value pair to be removed.
// The error returned indicates any issues encountered while removing the key-value pair.
// Implementing the CacheStore interface, this method uses the underlying ttlCache to remove the data.
// The ttlCache.Remove function is used to remove the key-value pair from the cache.
func (r memoryStore) Del(_ context.Context, key string) error {
	return r.ttlCache.Remove(key)
}

func (r memoryStore) Get(_ context.Context, key string) (*vo.CacheResponse, error) {
	value, err := r.ttlCache.Get(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return nil, domainmapper.NewErrCacheNotFound()
	} else if helper.IsNotNil(err) {
		return nil, err
	}
	var cacheResponse vo.CacheResponse
	err = helper.DecompressFromBase64WithGzipToDest(value, &cacheResponse)
	if helper.IsNotNil(err) {
		return nil, err
	}
	return &cacheResponse, nil
}

func (r memoryStore) Close() error {
	return nil
}
