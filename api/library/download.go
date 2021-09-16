package library

import (
	"context"
	"os"

	"github.com/go-resty/resty/v2"
)

type ResponseHandler func(c *Client, ctx context.Context, resp *resty.Response) error

type DownloadJob struct {
	Referer string
	Path    string
	Result  interface{}
	Handler ResponseHandler
}

func (j *DownloadJob) Do(c *Client, ctx context.Context) error {
	req := c.R(j.Referer)
	req.SetContext(ctx)
	if j.Result != nil {
		req.SetHeader("Accept", "application/json")
		req.SetResult(j.Result)
	}
	resp, err := req.Get(j.Path)
	if err != nil {
		os.Stdout.Write(resp.Body())
		return err
	}

	return j.Handler(c, ctx, resp)
}

func (wq *WorkQueue) Download(referer, path string, out interface{}, fn ResponseHandler) {
	if len(path) > 0 {
		req := &DownloadJob{
			Referer: referer,
			Path:    path,
			Result:  out,
			Handler: fn,
		}

		wq.d.Push(req)
		wq.Poke()
	}
}
