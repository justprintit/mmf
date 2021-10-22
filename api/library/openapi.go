package library

import (
	"context"
	"net/http"
)

func (w *Worker) openapiRequestEditor(ctx context.Context, req *http.Request) error {
	return nil
}

func (w *Worker) testOauth2(ctx context.Context) bool {
	return true
}
