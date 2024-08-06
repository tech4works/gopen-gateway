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

package cache

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-redis-template/redis"
	"github.com/GabrielHCataldo/go-redis-template/redis/option"
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

// redisStore represents a Redis cache store that implements the CacheStore interface.
type redisStore struct {
	redisTemplate *redis.Template
}

// NewRedisStore creates a new Redis cache store with the given address and password.
// It returns a CacheStore interface that can be used to interact with the Redis cache.
func NewRedisStore(address, password string) domain.Store {
	return &redisStore{
		redisTemplate: redis.NewTemplate(option.Client{
			Addr:     address,
			Password: password,
		}),
	}
}

func (r redisStore) Set(ctx context.Context, key string, cacheResponse *vo.CacheResponse) error {
	gzipBase64, err := helper.CompressWithGzipToBase64(cacheResponse)
	if checker.NonNil(err) {
		return err
	}

	return r.redisTemplate.Set(ctx, key, gzipBase64, option.NewSet().SetTTL(cacheResponse.Duration.Time()))
}

// Del deletes the value associated with the given key from the Redis cache.
// It takes the context and key as parameters.
// If the key does not exist in the cache, Del returns nil (no error is returned).
// If there is any error during the deletion, that error is returned.
// If everything goes well, Del returns nil.
func (r redisStore) Del(ctx context.Context, key string) error {
	return r.redisTemplate.Del(ctx, key)
}

func (r redisStore) Get(ctx context.Context, key string) (*vo.CacheResponse, error) {
	var cacheGzipBase64 string
	err := r.redisTemplate.Get(ctx, key, &cacheGzipBase64)
	if errors.Is(err, redis.ErrKeyNotFound) {
		return nil, mapper.NewErrCacheNotFound()
	} else if checker.NonNil(err) {
		return nil, err
	}

	var cacheResponse vo.CacheResponse
	err = helper.DecompressFromBase64WithGzipToDest(cacheGzipBase64, &cacheResponse)
	if checker.NonNil(err) {
		return nil, err
	}

	return &cacheResponse, nil
}

// Close closes the connection to the Redis server.
// It returns an error if there was a problem disconnecting from the server.
func (r redisStore) Close() error {
	return r.redisTemplate.Disconnect()
}
