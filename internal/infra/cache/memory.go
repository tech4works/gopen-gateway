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

package cache

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v2"
	"github.com/tech4works/checker"
	"github.com/tech4works/compressor"
	"github.com/tech4works/decompressor"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"go.elastic.co/apm/v2"
)

type memoryStore struct {
	ttlCache *ttlcache.Cache
}

func NewMemoryStore() domain.Store {
	ttlCache := ttlcache.NewCache()
	ttlCache.SkipTTLExtensionOnHit(true)
	return &memoryStore{
		ttlCache: ttlCache,
	}
}

func (m memoryStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	span, ctx := apm.StartSpan(ctx, "local.write", "cache")
	defer span.End()

	span.Context.SetLabel("key", key)

	b64, err := compressor.ToGzipBase64WithErr(value)
	if checker.NonNil(err) {
		return err
	}

	return m.ttlCache.SetWithTTL(key, b64, ttl)
}

func (m memoryStore) Del(_ context.Context, key string) error {
	return m.ttlCache.Remove(key)
}

func (m memoryStore) Get(ctx context.Context, key string) (string, error) {
	span, ctx := apm.StartSpan(ctx, "local.read", "cache")
	defer span.End()

	span.Context.SetLabel("key", key)

	cacheGzipBase64, err := m.ttlCache.Get(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return "", domain.NewErrCacheNotFound(key)
	} else if checker.NonNil(err) {
		return "", err
	}

	return decompressor.ToStringWithErr(decompressor.TypeGzipBase64, cacheGzipBase64)
}

func (m memoryStore) Close() error {
	return nil
}
