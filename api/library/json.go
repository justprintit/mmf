package library

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
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
