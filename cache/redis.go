package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"

	"time"
)

var (
	errFailedToSetCache = errors.New("failed to set cache")
)

type RedisCache struct {
	client redis.Cmdable
}

func (r *RedisCache) LoadAndDelete(ctx context.Context, key string) (any, error) {
	return r.client.GetDel(ctx, key).Result()
}

func NewRedisCache(client redis.Cmdable) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	res, err := r.client.Set(ctx, key, val, expiration).Result()
	if err != nil {
		return err
	}
	if res != "OK" {
		return errFailedToSetCache
	}
	return nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
