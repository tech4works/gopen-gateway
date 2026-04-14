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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/tech4works/checker"
	"github.com/tech4works/compressor"
	"github.com/tech4works/decompressor"
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/infra/telemetry"
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
	ctx, span := telemetry.Tracer().Start(ctx, "cache.global.write")
	defer span.End()

	span.SetAttributes(attribute.String("cache.key", key))

	b64, err := compressor.ToGzipBase64WithErr(value)
	if checker.NonNil(err) {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return r.client.Set(ctx, key, b64, ttl).Err()
}

func (r redisStore) Del(ctx context.Context, key string) error {
	ctx, span := telemetry.Tracer().Start(ctx, "cache.global.delete")
	defer span.End()

	span.SetAttributes(attribute.String("cache.key", key))

	return r.client.Del(ctx, key).Err()
}

func (r redisStore) Get(ctx context.Context, key string) (string, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "cache.global.read")
	defer span.End()

	span.SetAttributes(attribute.String("cache.key", key))

	cacheGzipBase64, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.NewErrCacheNotFound(key)
	} else if checker.NonNil(err) {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}

	return decompressor.ToStringWithErr(decompressor.TypeGzipBase64, cacheGzipBase64)
}

func (r redisStore) Close() error {
	return r.client.Close()
}
