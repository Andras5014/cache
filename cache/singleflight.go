package cache

import (
	"context"
	"errors"
	"golang.org/x/sync/singleflight"
	"time"
)

type SingleflightCacheV1 struct {
	ReadThroughCache
}

func NewSingleflightCacheV1(cache Cache, loadFunc func(ctx context.Context, key string) (any, error), expiration time.Duration) *SingleflightCacheV1 {
	g := singleflight.Group{}
	return &SingleflightCacheV1{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				val, err, _ := g.Do(key, func() (interface{}, error) {
					return loadFunc(ctx, key)
				})
				return val, err
			},
			Expiration: 0,
		},
	}
}

type SingleflightCacheV2 struct {
	ReadThroughCache
}

func (r *SingleflightCacheV2) GetV3(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if errors.Is(err, errKeyNotFound) {
		val, err, _ = r.g.Do(key, func() (interface{}, error) {
			v, er := r.LoadFunc(ctx, key)
			if er == nil {
				er = r.Cache.Set(ctx, key, v, r.Expiration)
				if er != nil {
					return v, ErrFailedToRefreshCache
				}
			}
			return v, er
		})

	}
	return val, err
}
