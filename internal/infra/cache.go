package infra

import (
	"context"
	"time"
)

type CacheStore interface {
	Set(ctx context.Context, key string, value any, expire time.Duration) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string, dest any) error
	Close() error
}
