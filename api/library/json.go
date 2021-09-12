package library

import (
	"context"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf/api/library/json"
)

func (c *Client) GetLibrary(ctx context.Context, library string, out interface{}) (*resty.Response, error) {
	path := fmt.Sprintf("/data-library/%s", library)
	req := c.J("/library?v=%s", library)

	if ctx != nil {
		req.SetContext(ctx)
	}

	if out != nil {
		req.SetResult(out)
	}

	return req.Get(path)
}

func (c *Client) GetSharedLibrary(ctx context.Context) (*json.Users, error) {
	out := &json.Users{}

	resp, err := c.GetLibrary(ctx, "shared", out)
	if err != nil {
		os.Stdout.Write(resp.Body())
	}
	return out, err
}

func (c *Client) GetPurchasesLibrary(ctx context.Context) (*json.Objects, error) {
	out := &json.Objects{}

	resp, err := c.GetLibrary(ctx, "purchases", out)
	if err != nil {
		os.Stdout.Write(resp.Body())
	}
	return out, err
}

func (c *Client) GetPledgesLibrary(ctx context.Context) (*json.Objects, error) {
	out := &json.Objects{}

	resp, err := c.GetLibrary(ctx, "campaigns", out)
	if err != nil {
		os.Stdout.Write(resp.Body())
	}
	return out, err
}
