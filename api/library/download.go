package library

import (
	"context"
	"os"

	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf/api/client"
)

type ResponseHandler func(c *Client, ctx context.Context, resp *resty.Response) error

type DownloadJob struct {
	client.RequestOptions

	h ResponseHandler
}

func (j *DownloadJob) Do(c *Client, ctx context.Context) error {

	req := j.New(&c.Client, ctx)

	if req.Method == "HEAD" {
		req.SetDoNotParseResponse(true)
	}

	resp, err := req.Get()
	if err != nil {
		os.Stdout.Write(resp.Body())
		return err
	}

	return j.h(c, ctx, resp)
}

func (wq *WorkQueue) Download(opt client.RequestOptions, fn ResponseHandler) {
	if len(opt.Path) > 0 {
		req := &DownloadJob{
			RequestOptions: opt,
			h:              fn,
		}

		wq.Request(req)
	}
}

func (wq *WorkQueue) Peek(opt client.RequestOptions, fn ResponseHandler) {
	if len(opt.Path) > 0 {
		opt.Method = resty.MethodHead

		req := &DownloadJob{
			RequestOptions: opt,
			h:              fn,
		}

		wq.Request(req)
	}
}

func (wq *WorkQueue) Request(req *DownloadJob) {
	if req != nil {
		wq.d.Push(req)
		wq.Poke()
	}
}
