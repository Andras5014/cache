package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroker_Send(t *testing.T) {
	b := &Broker{}
	go func() {
		for {
			err := b.Send(Msg{Content: time.Now().String()})
			if err != nil {
				t.Log(err)
				return
			}
			time.Sleep(100 * time.Millisecond)

		}
	}()
	wg := sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("消费者%d", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			msgs, err := b.Subscribe(100)
			if err != nil {
				t.Log(err)
				return
			}
			for msg := range msgs {
				fmt.Println(name, ":", msg.Content)
			}
		}()
	}
	wg.Wait()
}
