package library

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf/api/library/json"
)

func (c *Client) GetLibrary(ctx context.Context, library string, out interface{}) (*resty.Response, error) {
	return c.GetLibraryPage(ctx, library, 0, out)
}

func (c *Client) GetLibraryPage(ctx context.Context, library string, page int, out interface{}) (*resty.Response, error) {
	path := fmt.Sprintf("/data-library/%s", url.PathEscape(library))
	if page > 0 {
		path = fmt.Sprintf("%s?page=%v", path, page)
	}
	req := c.J("/library?v=%s", url.QueryEscape(library))

	if ctx != nil {
		req.SetContext(ctx)
	}

	if out != nil {
		req.SetResult(out)
	}

	return req.Get(path)
}

func (c *Client) GetSharedLibrary(ctx context.Context) (*json.Users, error) {
	return c.GetSharedLibraryPage(ctx, 0)
}

func (c *Client) GetSharedLibraryPage(ctx context.Context, page int) (*json.Users, error) {
	out := &json.Users{}

	resp, err := c.GetLibraryPage(ctx, "shared", page, out)
	if err != nil {
		os.Stdout.Write(resp.Body())
	}
	return out, err
}

func (c *Client) GetPurchasesLibrary(ctx context.Context) (*json.Objects, error) {
	return c.GetPurchasesLibraryPage(ctx, 0)
}

func (c *Client) GetPurchasesLibraryPage(ctx context.Context, page int) (*json.Objects, error) {
	out := &json.Objects{}

	resp, err := c.GetLibraryPage(ctx, "purchases", page, out)
	if err != nil {
		os.Stdout.Write(resp.Body())
	}
	return out, err
}

func (c *Client) GetPledgesLibrary(ctx context.Context) (*json.Objects, error) {
	return c.GetPledgesLibraryPage(ctx, 0)
}

func (c *Client) GetPledgesLibraryPage(ctx context.Context, page int) (*json.Objects, error) {
	out := &json.Objects{}

	resp, err := c.GetLibraryPage(ctx, "campaigns", page, out)
	if err != nil {
		os.Stdout.Write(resp.Body())
	}
	return out, err
}
