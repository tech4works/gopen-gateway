package infra

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-redis-template/redis"
	"github.com/GabrielHCataldo/go-redis-template/redis/option"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	appmapper "github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"time"
)

type redisStore struct {
	redisTemplate *redis.Template
}

func NewRedisStore(address, password string) interfaces.CacheStore {
	return &redisStore{
		redisTemplate: redis.NewTemplate(option.Client{
			Addr:     address,
			Password: password,
		}),
	}
}

func (r redisStore) Set(ctx context.Context, key string, value any, expire time.Duration) error {
	return r.redisTemplate.Set(ctx, key, value, option.NewSet().SetTTL(expire))
}

func (r redisStore) Del(ctx context.Context, key string) error {
	return r.redisTemplate.Del(ctx, key)
}

func (r redisStore) Get(ctx context.Context, key string, dest any) error {
	err := r.redisTemplate.Get(ctx, key, dest)
	if errors.Is(err, redis.ErrKeyNotFound) {
		return appmapper.NewErrCacheNotFound()
	} else if helper.IsNotNil(err) {
		return err
	}
	return nil
}

func (r redisStore) Close() error {
	return r.redisTemplate.Disconnect()
}
