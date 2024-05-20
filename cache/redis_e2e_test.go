//go:build e2e

package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedis_e2e_Set(t *testing.T) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	c := NewRedisCache(rdb)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := c.Set(ctx, "key", "value", time.Second*10)
	require.NoError(t, err)
	val, err := c.Get(ctx, "key")
	require.NoError(t, err)
	assert.Equal(t, "value", val)
}
