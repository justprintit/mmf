package library

import (
	"context"
	"sync"
	"time"

	"github.com/justprintit/mmf/api/library/store"
	"github.com/justprintit/mmf/api/library/types"
	"github.com/justprintit/mmf/api/openapi"
	"github.com/justprintit/mmf/api/transport"
)

type OpenAPIClient struct {
	openapi.ClientWithResponsesInterface
	openapi.ClientInterface
}

type Worker struct {
	*transport.Client

	mu      sync.Mutex
	refresh map[string]time.Time
	oac     OpenAPIClient
	data    types.Store
}

func NewWorker(c *transport.Client, data types.Store) *Worker {
	var w *Worker

	if c != nil {

		// Data Store
		if data == nil {
			data, _ = store.NewDummy()
		}

		// Worker
		w = &Worker{
			Client:  c,
			refresh: make(map[string]time.Time),
			data:    data,
		}

		// OpenAPI
		oac, err := openapi.NewClient(c.OpenAPIServer(),
			openapi.WithHTTPClient(c.NewOauth2Doer()),
			openapi.WithRequestEditorFn(w.openapiRequestEditor))

		if err != nil {
			panic(err)
		}

		w.oac = OpenAPIClient{
			ClientInterface:              oac,
			ClientWithResponsesInterface: &openapi.ClientWithResponses{oac},
		}

	}

	return w
}

func (c *Worker) Refresh() error {
	return nil
}

func (c *Worker) Start(ctx context.Context, downloaders int32) {
	c.Client.Spawn(ctx, downloaders)
	c.Client.Go(c.initiate, nil)
}

func (c *Worker) Run(ctx context.Context, downloaders int32) {
	c.Start(ctx, downloaders)
	c.Wait()
}

func (w *Worker) initiate(_ *transport.Client, ctx context.Context, _ interface{}) error {
	var ok = true

	if !w.testCredentials(ctx) {
		ok = false
	}

	if !w.testOauth2(ctx) {
		ok = false
	}

	if ok {
		w.Client.Start()
	} else {
		w.Client.Pause()
	}
	return nil
}

func (w *Worker) testCredentials(ctx context.Context) bool {
	return true
}
