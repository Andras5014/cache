package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBuildInMapCache_Get(t *testing.T) {

	testCases := []struct {
		name    string
		key     string
		cache   func() *BuildInMapCache
		wantVal any
		wantErr error
	}{
		{
			name: "key not found",
			key:  "not exist key",
			cache: func() *BuildInMapCache {
				return NewBuildInMapCache(time.Second * 10)
			},
			wantVal: "not exist key",
			wantErr: errKeyNotFound,
		}, {
			name: "get val",
			key:  "key1",
			cache: func() *BuildInMapCache {
				res := NewBuildInMapCache(time.Second * 10)
				err := res.Set(nil, "key1", 123, time.Minute)
				require.NoError(t, err)
				return res
			},
			wantVal: 123,
		}, {
			name: "key expired",
			key:  "key1",
			cache: func() *BuildInMapCache {
				res := NewBuildInMapCache(time.Second * 10)
				err := res.Set(nil, "key1", 123, time.Nanosecond)
				require.NoError(t, err)
				time.Sleep(time.Second * 2)
				return res
			},
			wantErr: errKeyExpired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.cache()
			val, err := c.Get(nil, tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}
}
func TestBuildInMapCache_Loop(t *testing.T) {
	cnt := 0
	c := NewBuildInMapCache(time.Second, BuildInMapCacheWithEvictedCallback(func(key string, val any) {
		cnt++
	}))
	err := c.Set(nil, "key1", 123, time.Second)
	require.NoError(t, err)
	time.Sleep(time.Second * 3)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.data["key1"]
	require.False(t, ok)
	require.Equal(t, 1, cnt)
}
