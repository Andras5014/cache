package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 有问题！！！
func TestClient_Lock(t *testing.T) {
	testCases := []struct {
		name string
		mock func() redis.Cmdable
		key  string

		wantErr  error
		wantLock *Lock
	}{
		{
			name: "set nx error",
			//mock: func() redis.Cmdable {
			//},
			key:     "test",
			wantErr: nil,
			wantLock: &Lock{
				key:    "test",
				value:  "test",
				client: nil,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(tc.mock())
			l, err := client.tryLock(context.Background(), tc.key, time.Second)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLock, l)
		})
	}

}
