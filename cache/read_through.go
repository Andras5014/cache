package cache

import (
	"context"
	"errors"
	"golang.org/x/sync/singleflight"
	"log"
	"time"
)

var (
	ErrFailedToRefreshCache = errors.New("failed to refresh cache")
)

// ReadThroughCache 一定要赋值LoadFunc和Expiration
// 否则会panic
type ReadThroughCache struct {
	Cache
	LoadFunc   func(ctx context.Context, key string) (any, error)
	Expiration time.Duration
	g          singleflight.Group
}

func (r *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if errors.Is(err, errKeyNotFound) {
		val, err = r.LoadFunc(ctx, key)
		if err != nil {
			er := r.Cache.Set(ctx, key, val, r.Expiration)
			if er != nil {
				return nil, ErrFailedToRefreshCache
			}
		}

	}
	return val, err
}
func (r *ReadThroughCache) GetV1(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if errors.Is(err, errKeyNotFound) {
		go func() {
			val, err = r.LoadFunc(ctx, key)
			if err != nil {
				er := r.Cache.Set(ctx, key, val, r.Expiration)
				if er != nil {
					log.Fatalln(er)
				}
			}
		}()

	}
	return val, err
}
func (r *ReadThroughCache) GetV2(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if errors.Is(err, errKeyNotFound) {
		val, err = r.LoadFunc(ctx, key)
		if err != nil {
			go func() {
				er := r.Cache.Set(ctx, key, val, r.Expiration)
				if er != nil {
					log.Fatalln(er)
				}
			}()
		}

	}
	return val, err
}

//func (r *ReadThroughCache) GetV3(ctx context.Context, key string) (any, error) {
//	val, err := r.Cache.Get(ctx, key)
//	if errors.Is(err, errKeyNotFound) {
//		val, err, _ = r.g.Do(key, func() (interface{}, error) {
//			v, er := r.LoadFunc(ctx, key)
//			if er == nil {
//				er = r.Cache.Set(ctx, key, v, r.Expiration)
//				if er != nil {
//					return v, ErrFailedToRefreshCache
//				}
//			}
//			return v, er
//		})
//
//	}
//	return val, err
//}
