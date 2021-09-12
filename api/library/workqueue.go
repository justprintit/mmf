package library

import (
	"context"
	"sync"

	"go.sancus.dev/core/queue"
)

type WorkJob func(c *Client, ctx context.Context) error

type WorkQueue struct {
	c      *Client
	mu     sync.Mutex
	q      queue.Queue
	done   chan error
	cancel context.CancelFunc
	ctx    context.Context
	wg     sync.WaitGroup
}

func (wq *WorkQueue) Init(c *Client) {
	wq.c = c
}

func (wq *WorkQueue) Start() {

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	// init
	*wq = WorkQueue{
		c:      wq.c,
		done:   done,
		ctx:    ctx,
		cancel: cancel,
	}

	go func() {
		defer close(done)
		for {
			if v, ok := wq.q.Pop(); !ok {
				// empty
				break
			} else if f, ok := v.(WorkJob); !ok {
				// wtf? ignore
			} else if err := f(wq.c, wq.ctx); err != nil {
				// abort
				wq.Cancel()
				break
			}
		}

		wq.wg.Wait()
	}()
}

func (wq *WorkQueue) Add(f WorkJob) {
	if f != nil {
		wq.q.Push(f)
	}
}

func (wq *WorkQueue) Spawn(f WorkJob) {
	if f != nil {
		wq.wg.Add(1)

		go func() {
			defer wq.wg.Done()
			f(wq.c, wq.ctx)
		}()
	}
}

func (wq *WorkQueue) Context() context.Context {
	return wq.ctx
}

func (wq *WorkQueue) Cancel() {
	wq.cancel()
}

func (wq *WorkQueue) Done() <-chan error {
	return wq.done
}
