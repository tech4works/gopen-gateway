package infra

import (
	"context"
	"time"
)

type CacheStore interface {
	// Set sets an item to the Cache, replacing any existing item.
	Set(ctx context.Context, key string, value any, expire time.Duration) error
	// Del removes an item from the Cache. Does nothing if the key is not in the Cache.
	Del(ctx context.Context, key string) error
	// Get retrieves an item from the Cache. if key does not exist in the store, return ErrCacheMiss
	Get(ctx context.Context, key string, dest any) error
	// Close client connection
	Close() error
}
