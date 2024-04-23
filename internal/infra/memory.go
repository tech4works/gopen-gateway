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
	appmapper "github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/jellydator/ttlcache/v2"
	"time"
)

// memoryStore represents an in-memory cache store that implements the CacheStore interface.
type memoryStore struct {
	ttlCache *ttlcache.Cache
}

// NewMemoryStore returns a new instance of the MemoryStore structure that implements the CacheStore interface.
// This implementation uses an in-memory cache with a time-to-live (TTL)
func NewMemoryStore() CacheStore {
	ttlCache := ttlcache.NewCache()
	ttlCache.SkipTTLExtensionOnHit(true)
	return &memoryStore{
		ttlCache: ttlCache,
	}
}

// Set sets a key-value pair in the memory cache with the specified expiration duration.
// The key is a string that serves as the identifier, the value is any data that can be stored in the cache,
// and the expiry duration determines how long the key-value pair will remain in the cache.
// The error returned indicates any issues encountered while setting the key-value pair.
// Implementing the CacheStore interface, this method uses the underlying ttlCache to store the data.
// The ttlCache.SetWithTTL function is used to set the key-value pair with the specified expiration.
func (r memoryStore) Set(_ context.Context, key string, value any, expire time.Duration) error {
	return r.ttlCache.SetWithTTL(key, value, expire)
}

// Del removes a key-value pair from the memory cache with the specified key.
// The key is a string that serves as the identifier for the key-value pair to be removed.
// The error returned indicates any issues encountered while removing the key-value pair.
// Implementing the CacheStore interface, this method uses the underlying ttlCache to remove the data.
// The ttlCache.Remove function is used to remove the key-value pair from the cache.
func (r memoryStore) Del(_ context.Context, key string) error {
	return r.ttlCache.Remove(key)
}

// Get retrieves the value associated with the specified key from the memory cache.
// The key is a string that serves as the identifier.
// The dest parameter represents the destination where the retrieved value will be stored.
// The error returned indicates any issues encountered while retrieving the value.
// If the specified key is not found in the memory cache, it returns a cache not found error.
// If an error other than cache not found occurs, it returns that error.
// The helper.ConvertToDest function is used to convert the retrieved value to the destination type.
func (r memoryStore) Get(_ context.Context, key string, dest any) error {
	value, err := r.ttlCache.Get(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return appmapper.NewErrCacheNotFound()
	} else if helper.IsNotNil(err) {
		return err
	}
	return helper.ConvertToDest(value, dest)
}

func (r memoryStore) Close() error {
	return nil
}
