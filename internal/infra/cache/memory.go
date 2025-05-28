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
	"github.com/jellydator/ttlcache/v2"
	"github.com/tech4works/checker"
	"github.com/tech4works/compressor"
	"github.com/tech4works/converter"
	"github.com/tech4works/decompressor"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
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

func (m memoryStore) Set(ctx context.Context, key string, cacheResponse *vo.CacheResponse) error {
	span, ctx := apm.StartSpan(ctx, "local.write", "cache")
	defer span.End()

	span.Context.SetLabel("key", key)

	b64, err := compressor.ToGzipBase64WithErr(cacheResponse)
	if checker.NonNil(err) {
		return err
	}

	return m.ttlCache.SetWithTTL(key, b64, cacheResponse.Duration.Time())
}

func (m memoryStore) Del(_ context.Context, key string) error {
	return m.ttlCache.Remove(key)
}

func (m memoryStore) Get(ctx context.Context, key string) (*vo.CacheResponse, error) {
	span, ctx := apm.StartSpan(ctx, "local.read", "cache")
	defer span.End()

	span.Context.SetLabel("key", key)

	value, err := m.ttlCache.Get(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return nil, mapper.NewErrCacheNotFound()
	} else if checker.NonNil(err) {
		return nil, err
	}

	bs, err := decompressor.ToBytesWithErr(decompressor.TypeGzipBase64, value)
	if checker.NonNil(err) {
		return nil, err
	}

	var cacheResponse vo.CacheResponse
	err = converter.ToDestWithErr(bs, &cacheResponse)
	if checker.NonNil(err) {
		return nil, err
	}

	return &cacheResponse, nil
}

func (m memoryStore) Close() error {
	return nil
}
