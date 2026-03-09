package cache

import (
	"context"
	"errors"
	"time"
)

var ErrCacheKeyNotFound = errors.New("cache key not found")

type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
}
