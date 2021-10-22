package library

import (
	"context"

	"github.com/justprintit/mmf/api/library/store"
	"github.com/justprintit/mmf/api/transport"
	"github.com/justprintit/mmf/types"
)

type Worker struct {
	*transport.Client

	data types.Store
}

func NewWorker(c *transport.Client, data types.Store) *Worker {
	if c != nil {

		if data == nil {
			data, _ = store.NewDummy()
		}

		return &Worker{
			Client: c,
			data:   data,
		}
	}
	return nil
}

func (c *Worker) Refresh() error {
	return nil
}

func (c *Worker) Start(ctx context.Context, downloaders int32) {
	c.Client.Spawn(ctx, downloaders)
	c.Client.Start()
}

func (c *Worker) Run(ctx context.Context, downloaders int32) {
	c.Start(ctx, downloaders)
	c.Wait()
}
