package infra

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-redis-template/redis"
	"github.com/GabrielHCataldo/go-redis-template/redis/option"
	appmapper "github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"time"
)

// redisStore represents a Redis cache store that implements the CacheStore interface.
type redisStore struct {
	redisTemplate *redis.Template
}

// NewRedisStore creates a new Redis cache store with the given address and password.
// It returns a CacheStore interface that can be used to interact with the Redis cache.
func NewRedisStore(address, password string) CacheStore {
	return &redisStore{
		redisTemplate: redis.NewTemplate(option.Client{
			Addr:     address,
			Password: password,
		}),
	}
}

// Set sets the value of the given key in the Redis cache.
// It takes the context, key, value, and expiration duration as parameters.
// The value can be of any type.
// The expiry parameter specifies the time after which the key will expire in the cache.
// It returns an error if there was a problem setting the value in the cache.
// If the key does not exist in the cache, it will be created.
// If the value is not already of type "any", it will be converted to "any".
func (r redisStore) Set(ctx context.Context, key string, value any, expire time.Duration) error {
	return r.redisTemplate.Set(ctx, key, value, option.NewSet().SetTTL(expire))
}

// Del deletes the value associated with the given key from the Redis cache.
// It takes the context and key as parameters.
// If the key does not exist in the cache, Del returns nil (no error is returned).
// If there is any error during the deletion, that error is returned.
// If everything goes well, Del returns nil.
func (r redisStore) Del(ctx context.Context, key string) error {
	return r.redisTemplate.Del(ctx, key)
}

// Get retrieves the value associated with the given key from the Redis cache.
// It takes the context, key, and a destination variable as parameters.
// The destination variable represents the result of the retrieval operation.
// If the key does not exist in the cache, Get returns an error of type ErrCacheNotFound.
// If there is any other error during the retrieval, that error is returned.
// If everything goes well, Get returns nil.
func (r redisStore) Get(ctx context.Context, key string, dest any) error {
	err := r.redisTemplate.Get(ctx, key, dest)
	if errors.Is(err, redis.ErrKeyNotFound) {
		return appmapper.NewErrCacheNotFound()
	} else if helper.IsNotNil(err) {
		return err
	}
	return nil
}

// Close closes the connection to the Redis server.
// It returns an error if there was a problem disconnecting from the server.
func (r redisStore) Close() error {
	return r.redisTemplate.Disconnect()
}
