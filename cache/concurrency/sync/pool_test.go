package sync

import (
	"sync"
	"testing"
)

func TestPool(t *testing.T) {

	p := sync.Pool{
		New: func() any {
			t.Log("创建资源")
			return "hello"

		},
	}
	obj := p.Get()
	p.Put(obj)
}
