package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	errKeyNotFound = errors.New("cache:key is not found")
	errKeyExpired  = errors.New("cache:key is expired")
)

type BuildInMapCacheOption func(*BuildInMapCache)
type BuildInMapCache struct {
	data      map[string]*Item
	mutex     sync.RWMutex
	close     chan struct{}
	onEvicted func(key string, val any)
}

func (b *BuildInMapCache) LoadAndDelete(ctx context.Context, key string) (any, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	val, ok := b.data[key]
	if !ok {
		return nil, errKeyNotFound
	}
	b.delete(key)
	return val.val, nil
}

func NewBuildInMapCache(interval time.Duration, opts ...BuildInMapCacheOption) *BuildInMapCache {
	res := &BuildInMapCache{
		data:  make(map[string]*Item, 100),
		close: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(res)
	}
	//maxCnt:=1000
	go func() {
		ticker := time.NewTicker(interval)
		for {

			select {
			case t := <-ticker.C:
				res.mutex.Lock()
				i := 0
				for k, v := range res.data {
					if i > 10000 {
						break
					}
					if v.deadlineBefore(t) {
						res.delete(k)
					}
					i++
				}
				res.mutex.Unlock()
			case <-res.close:
				return

			}
		}
	}()
	return res
}
func BuildInMapCacheWithEvictedCallback(fn func(key string, val any)) BuildInMapCacheOption {
	return func(cache *BuildInMapCache) {
		cache.onEvicted = fn
	}
}
func (b *BuildInMapCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.set(key, val, expiration)
}

func (b *BuildInMapCache) set(key string, val any, expiration time.Duration) error {
	var dl time.Time
	if expiration > 0 {
		dl = time.Now().Add(expiration)
	}
	b.data[key] = &Item{
		val:      val,
		deadline: dl,
	}
	return nil
}
func (b *BuildInMapCache) Get(ctx context.Context, key string) (any, error) {
	b.mutex.RLock()

	res, ok := b.data[key]
	b.mutex.RUnlock()

	if !ok {
		return nil, errKeyNotFound
	}
	now := time.Now()
	if res.deadlineBefore(now) {
		b.mutex.Lock()
		defer b.mutex.Unlock()
		res, ok = b.data[key]
		if !ok {
			return nil, errKeyNotFound
		}
		if res.deadlineBefore(now) {
			b.delete(key)
			return nil, errKeyExpired
		}

	}
	return res.val, nil
}

func (b *BuildInMapCache) Delete(ctx context.Context, key string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.delete(key)
	return nil
}
func (b *BuildInMapCache) delete(key string) {
	itm, ok := b.data[key]
	if ok {
		if b.onEvicted != nil {
			b.onEvicted(key, itm.val)
		}
		delete(b.data, key)
	}
}
func (b *BuildInMapCache) Close() error {
	select {
	case b.close <- struct{}{}:
	default:
		return errors.New("cache:close is already called")
	}
	return nil
}

type Item struct {
	val      any
	deadline time.Time
}

func (i *Item) deadlineBefore(t time.Time) bool {
	return !i.deadline.IsZero() && i.deadline.Before(t)
}
