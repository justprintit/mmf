package transport

import (
	"context"
	"log"
	"sync"
	"sync/atomic"

	"go.sancus.dev/core/queue"
)

type QueueIndex uint

type QueueFunc func(c *Client, ctx context.Context, data interface{}) error

type QueueEntry struct {
	fn   QueueFunc
	data interface{}
}

type WorkQueueState int

const (
	WorkQueueStarting WorkQueueState = iota
	WorkQueuePaused
	WorkQueueRunning
	WorkQueueTerminating
)

type WorkQueue struct {
	sync.Mutex

	state atomic.Value

	// work
	c  *Client
	q  []Queue
	wg sync.WaitGroup

	// cancel
	ctx    context.Context
	cancel context.CancelFunc
	once   sync.Once
	done   chan struct{}
}

func (wq *WorkQueue) setState(s WorkQueueState) {
	wq.Lock()
	defer wq.Unlock()

	wq.state.Store(s)
}

func (wq *WorkQueue) State() WorkQueueState {
	if v := wq.state.Load(); v != nil {
		return v.(WorkQueueState)
	}
	return WorkQueueStarting
}

func (wq *WorkQueue) running() bool {
	return wq.State() == WorkQueueRunning
}

func (wq *WorkQueue) Len() int {
	return len(wq.q)
}

func (wq *WorkQueue) Init(c *Client, count int) {
	wq.c = c
	wq.q = make([]Queue, count)
}

func (wq *WorkQueue) Start(ctx context.Context, limits ...int32) {
	wq.Lock()
	defer wq.Unlock()

	// prepare queues
	//
	if len(limits) > len(wq.q) {
		// align sizes, extend q
		extra := make([]Queue, len(limits)-len(wq.q))
		wq.q = append(wq.q, extra...)
	}

	if len(limits) < len(wq.q) {
		// align sizes, extend count
		extra := make([]int32, len(wq.q)-len(limits))
		limits = append(limits, extra...)
	}

	for i := range wq.q {
		q := &wq.q[i]
		q.wg = &wq.wg

		if limits[i] > 0 {
			q.Max = limits[i]
		} else if q.Max == 0 {
			q.Max = 1
		}
	}

	// cancel
	//
	if ctx == nil {
		ctx = context.Background()
	}

	wq.ctx, wq.cancel = context.WithCancel(ctx)
	wq.done = make(chan struct{})
	wq.state.Store(WorkQueueStarting)

	// watcher
	go func() {
		defer close(wq.done)

		// wait for cancellation
		select {
		case <-wq.ctx.Done():
			wq.setState(WorkQueueTerminating)
		}

		// and wait for workers
		wq.wg.Wait()
	}()
}

func (wq *WorkQueue) run(q *Queue) {
	for wq.running() {
		if v, ok := q.Pop(); !ok {
			// empty
			break
		} else if qe, ok := v.(QueueEntry); !ok {
			// wtf? ignore
		} else if err := qe.fn(wq.c, wq.ctx, qe.data); err != nil {
			// abort
			log.Printf("Fatal: %v: %s", qe.fn, err)
			wq.Cancel()
			break
		}
	}
}

// Poke() wakes workers if they are needed
func (wq *WorkQueue) Poke() {
	if wq.running() {
		wq.Lock()
		defer wq.Unlock()

		wq.poke()
	}
}

func (wq *WorkQueue) poke() {
	for i := range wq.q {
		q := &wq.q[i]
		if !q.Empty() {
			if q.Add() {
				// spawn worker
				go func() {
					defer q.Done()
					wq.run(q)
				}()
			}
		}
	}
}

// Adds a Task to a queue
func (wq *WorkQueue) Add(i QueueIndex, f QueueFunc, data interface{}) {
	if f != nil {
		wq.q[i].Push(QueueEntry{f, data})
		wq.Poke()
	}
}

// Runs a Task on its own goroutine
func (wq *WorkQueue) Go(f QueueFunc, data interface{}) {
	if f != nil {
		wq.wg.Add(1)

		go func() {
			defer wq.wg.Done()
			f(wq.c, wq.ctx, data)
		}()
	}
}

func (wq *WorkQueue) Context() context.Context {
	return wq.ctx
}

func (wq *WorkQueue) Cancel() {
	wq.once.Do(wq.cancel)
}

func (wq *WorkQueue) Done() <-chan struct{} {
	return wq.done
}

func (wq *WorkQueue) Wait() {
	<-wq.done
}

type Queue struct {
	queue.Queue

	wg    *sync.WaitGroup
	Count int32
	Max   int32
}

func (q *Queue) Add() bool {
	if n0 := q.Count; n0 < q.Max {
		// try
		if atomic.CompareAndSwapInt32(&q.Count, n0, n0+1) {
			q.wg.Add(1)
			return true
		}
	}
	return false
}

func (q *Queue) Done() {
	atomic.AddInt32(&q.Count, -1)
	q.wg.Done()
}
