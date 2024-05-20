package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

var (
	ErrFailedPreemptLock = errors.New("failed to preempt lock")
	ErrLockNotHold       = errors.New("lock not hold")
	//go:embed lua/unlock.lua
	luaUnlock string
	//go:embed lua/refresh.lua
	luaRefresh string
	//go:embed lua/lock.lua
	luaLock string
)

// Client 对redis.Cmdable的二次封装
type Client struct {
	client redis.Cmdable
}

func NewClient(client redis.Cmdable) *Client {
	return &Client{
		client: client,
	}
}
func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration, timeout time.Duration, retry RetryStrategy) (*Lock, error) {
	var timer *time.Timer
	val := uuid.New().String()
	for {
		//retry
		lctx, cancel := context.WithTimeout(ctx, timeout)
		res, err := c.client.Eval(lctx, luaLock, []string{key}, val, expiration.Seconds()).Result()
		cancel()
		if err != nil && errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		if res == "OK" {
			return &Lock{
				client:     c.client,
				key:        key,
				value:      val,
				expiration: expiration,
				unlockChan: make(chan struct{}, 1),
			}, nil
		}
		interval, ok := retry.Next()
		if !ok {
			return nil, fmt.Errorf("reids-lock:超出重试限制,%w", ErrFailedPreemptLock)
		}
		if timer == nil {
			timer = time.NewTimer(interval)
		} else {
			timer.Reset(interval)
		}
		select {
		case <-timer.C:
		case <-ctx.Done():
			return nil, ctx.Err()

		}

	}
}

func (c *Client) tryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	val := uuid.New().String()
	ok, err := c.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		//别人抢到锁
		return nil, ErrFailedPreemptLock
	}
	//defer func() {
	//	if err = c.client.Del(ctx,key).Err();err!=nil {
	//		panic(err)
	//	}
	//}()
	return &Lock{
		client:     c.client,
		key:        key,
		value:      val,
		expiration: expiration,
		unlockChan: make(chan struct{}, 1),
	}, nil
}

type Lock struct {
	client     redis.Cmdable
	key        string
	value      string
	expiration time.Duration
	unlockChan chan struct{}
}

func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.value).Int64()
	if errors.Is(err, redis.Nil) {
		return ErrLockNotHold
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}
func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.expiration.Seconds()).Int64()
	//defer func() {
	//	select {
	//	case l.unlockChan <- struct{}{}:
	//	default:
	//
	//	}
	//}()
	if errors.Is(err, redis.Nil) {
		return ErrLockNotHold
	}
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}

//func (l *Lock) Unlock(ctx context.Context) error {
//	//判断是不是我的锁
//	if l.value == l.client.Get(ctx, l.key).Val() {
//	}
//	cnt, err := l.client.Del(ctx, l.key).Result()
//	if err != nil {
//		return err
//	}
//	if cnt != 1 {
//		//加锁过期
//		return ErrLockNotExist
//	}
//	return nil
//
//}
