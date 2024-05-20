package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	errOverCapacity = errors.New("cache:over capacity")
)

// MaxCntCache 控制住缓存住的键值对数量
type MaxCntCache struct {
	*BuildInMapCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(c *BuildInMapCache, maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		BuildInMapCache: c,
		maxCnt:          maxCnt,
	}
	origin := c.onEvicted
	res.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}
	return res
}
func (c *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	//如果key存在，技术不准确
	//cnt := atomic.AddInt32(&c.cnt, 1)
	//if cnt > c.maxCnt {
	//	atomic.AddInt32(&c.cnt, -1)
	//	return errOverCapacity
	//}
	//return c.Set(ctx, key, val, expiration)

	//并发问题，同时访问lock多次
	//c.mutex.Lock()
	//_, ok := c.data[key]
	//if !ok {
	//	c.cnt++
	//}
	//
	//if c.cnt > c.maxCnt {
	//	c.mutex.Unlock()
	//	return errOverCapacity
	//}
	//defer c.mutex.Unlock()
	//return c.Set(ctx, key, val, expiration)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.data[key]
	if !ok {
		if c.cnt+1 > c.maxCnt {
			return errOverCapacity
		}
		c.cnt++
	}
	return c.set(key, val, expiration)
}
