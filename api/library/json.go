package library

import (
	"context"

	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf/api/client"
)

func (c *Client) GetPage(ctx context.Context, opt client.RequestOptions, page int) (*resty.Response, error) {
	if page > 0 {
		opt = opt.Clone()
		opt.Setf("page", "%v", page)
	}

	return opt.New(&c.Client, ctx).Get()
}

func (c *Client) SchedulePageRequest(opt client.RequestOptions, page int, fn ResponseHandler) {
	if page > 0 {
		opt = opt.Clone()
		opt.Setf("page", "%v", page)
	}

	c.Download(opt, fn)
}
