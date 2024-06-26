package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Delete(ctx context.Context, key string) error

	LoadAndDelete(ctx context.Context, key string) (any, error)
}

//type CacheV2[T any] interface {
//	Set(ctx context.Context, key string, val T, expiration time.Duration) error
//	Get(ctx context.Context, key string) (T, error)
//	Delete(ctx context.Context, key string) error
//}
