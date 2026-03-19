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

	"github.com/redis/go-redis/v9"
	"github.com/tech4works/checker"
	"github.com/tech4works/compressor"
	"github.com/tech4works/decompressor"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"go.elastic.co/apm/v2"
)

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(address, password string) domain.Store {
	return &redisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
		}),
	}
}

func (r redisStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	span, ctx := apm.StartSpan(ctx, "global.write", "cache")
	defer span.End()

	span.Context.SetLabel("key", key)

	b64, err := compressor.ToGzipBase64WithErr(value)
	if checker.NonNil(err) {
		return err
	}

	return r.client.Set(ctx, key, b64, ttl).Err()
}

func (r redisStore) Del(ctx context.Context, key string) error {
	span, ctx := apm.StartSpan(ctx, "global.read", "cache")
	defer span.End()

	span.Context.SetLabel("key", key)

	return r.client.Del(ctx, key).Err()
}

func (r redisStore) Get(ctx context.Context, key string) (string, error) {
	span, ctx := apm.StartSpan(ctx, "global.read", "cache")
	defer span.End()

	span.Context.SetLabel("key", key)

	cacheGzipBase64, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.NewErrCacheNotFound(key)
	} else if checker.NonNil(err) {
		return "", err
	}

	return decompressor.ToStringWithErr(decompressor.TypeGzipBase64, cacheGzipBase64)
}

func (r redisStore) Close() error {
	return r.client.Close()
}
