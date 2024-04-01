package infra

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	appmapper "github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/jellydator/ttlcache/v2"
	"time"
)

type memoryStore struct {
	ttlCache *ttlcache.Cache
}

func NewMemoryStore() CacheStore {
	ttlCache := ttlcache.NewCache()
	ttlCache.SkipTTLExtensionOnHit(true)

	return &memoryStore{
		ttlCache: ttlCache,
	}
}

func (r memoryStore) Set(_ context.Context, key string, value any, expire time.Duration) error {
	return r.ttlCache.SetWithTTL(key, value, expire)
}

func (r memoryStore) Del(_ context.Context, key string) error {
	return r.ttlCache.Remove(key)
}

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
