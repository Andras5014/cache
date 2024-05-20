package channel

import "context"

type Task func()
type TaskPool struct {
	tasks chan Task
	close chan struct{}
}

func NewTaskPool(numG int, capacity int) *TaskPool {
	res := &TaskPool{
		tasks: make(chan Task, capacity),
		close: make(chan struct{}),
	}

	for i := 0; i < numG; i++ {
		go func() {
			for {
				select {
				case <-res.close:
					return
				case t := <-res.tasks:
					t()
				}

			}

		}()
	}
	return res
}

func (p TaskPool) Submit(ctx context.Context, t Task) error {
	select {
	case p.tasks <- t:
	case <-ctx.Done():
		return ctx.Err()

	}
	return nil
}
func (p *TaskPool) Close() error {
	close(p.close)
	return nil
}
