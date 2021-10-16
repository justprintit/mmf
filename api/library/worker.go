package library

import (
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
