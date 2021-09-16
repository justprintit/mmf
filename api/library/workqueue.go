package library

import (
	"context"
	"log"
	"sync"
	"sync/atomic"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/core/queue"
)

type Queue struct {
	queue.Queue

	wg    *sync.WaitGroup
	Count int32
	Max   int32
}

func (p *Queue) MayAdd() bool {
	if n0 := p.Count; n0 < p.Max {
		// try
		if atomic.CompareAndSwapInt32(&p.Count, n0, n0+1) {
			p.wg.Add(1)
			return true
		}
	}
	return false
}

func (p *Queue) Add() {
	atomic.AddInt32(&p.Count, 1)
	p.wg.Add(1)
}

func (p *Queue) Done() {
	atomic.AddInt32(&p.Count, -1)
	p.wg.Done()
}

type WorkJob func(c *Client, ctx context.Context) error

type WorkQueue struct {
	c      *Client
	mu     sync.Mutex
	q      Queue // WorkJob
	d      Queue // DownloadJob
	done   chan error
	once   sync.Once // cancel only once
	cancel context.CancelFunc
	ctx    context.Context
	wg     sync.WaitGroup
}

func (wq *WorkQueue) Init(c *Client) {
	ctx, cancel := context.WithCancel(context.Background())

	// client
	wq.c = c
	// cancel
	wq.ctx = ctx
	wq.cancel = cancel
	// waitgroup
	wq.q.wg = &wq.wg
	wq.d.wg = &wq.wg
}

func (wq *WorkQueue) Start(n int32) {

	if wq.done != nil {
		panic(errors.New("%T can only be Started once", wq))
	}

	// init
	wq.done = make(chan error)
	wq.q.Max = 1 // WorkJobs run sequentially
	wq.d.Max = n // DownloadJobs in parallel

	go func() {
		// launch workers
		wq.Poke()

		// and wait until they are all done
		defer close(wq.done)
		wq.wg.Wait()
	}()
}

func (wq *WorkQueue) Poke() {
	wq.mu.Lock()
	defer wq.mu.Unlock()

	// WorkJobs
	if !wq.q.Empty() {
		if wq.q.MayAdd() {
			go func() {
				defer wq.q.Done()
				wq.runQueueWorker()
			}()
		}
	}

	// DownloadJob
	if !wq.d.Empty() {
		if wq.d.MayAdd() {
			go func() {
				defer wq.d.Done()
				wq.runDownloadWorker()
			}()
		}
	}
}

func (wq *WorkQueue) runQueueWorker() {
	for {
		if v, ok := wq.q.Pop(); !ok {
			// empty
			break
		} else if f, ok := v.(WorkJob); !ok {
			// wtf? ignore
		} else if err := f(wq.c, wq.ctx); err != nil {
			// abort
			log.Printf("Fatal: %v: %s", f, err)
			wq.Cancel()
			break
		}
	}
}

func (wq *WorkQueue) runDownloadWorker() {
	for {
		if v, ok := wq.d.Pop(); !ok {
			// empty
			break
		} else if req, ok := v.(*DownloadJob); !ok {
			// wtf? ignore
		} else if err := req.Do(wq.c, wq.ctx); err != nil {
			// abort
			log.Printf("Fatal: %s: %s", "Download", err)
			wq.Cancel()
			break
		}
	}
}

func (wq *WorkQueue) Add(f WorkJob) {
	if f != nil {
		wq.q.Push(f)
		wq.Poke()
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
	wq.once.Do(wq.cancel)
}

func (wq *WorkQueue) Done() <-chan error {
	return wq.done
}
