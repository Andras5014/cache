package context

import (
	"context"
	"testing"
	"time"
)

type mykey struct {
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, mykey{}, "my-value")
	//ctx, cancel := context.WithCancel(ctx)
	//
	//cancel()
	val := ctx.Value(mykey{}).(string)
	t.Log(val)
	ctx.Value("my-key")
	//Testcases := []struct {
	//	Name string
	//}{
	//	{},
	//}
	//
	//for i, i2 := range collection {
	//
	//}
}
func TestContext_WithCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()
	<-ctx.Done()
	t.Log("hello,cancel:", ctx.Err())

}
func TestContext_WithDeadline(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	t.Log(time.Now())
	deadline, _ := ctx.Deadline()
	t.Log(deadline)
	defer cancel()
	<-ctx.Done()
	t.Log("hello,canceldeadline:", ctx.Err())
}
func TestContext_Parent(t *testing.T) {
	ctx := context.Background()
	parent := context.WithValue(ctx, "my-key", "my-value")
	child := context.WithValue(parent, "my-key", "my new val")

	t.Log("parent:", parent.Value("my-key"))
	t.Log("child:", child.Value("my-key"))
}
