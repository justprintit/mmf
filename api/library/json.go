package library

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

func (c *Client) composeLibraryPageRequest(library string, page int) (string, string) {
	path := fmt.Sprintf("/data-library/%s", url.PathEscape(library))
	if page > 0 {
		path = fmt.Sprintf("%s?page=%v", path, page)
	}
	referer := fmt.Sprintf("/library?v=%s", url.QueryEscape(library))
	return referer, path
}

func (c *Client) GetLibrary(ctx context.Context, library string, out interface{}) (*resty.Response, error) {
	return c.GetLibraryPage(ctx, library, 0, out)
}

func (c *Client) GetLibraryPage(ctx context.Context, library string, page int, out interface{}) (*resty.Response, error) {
	referer, path := c.composeLibraryPageRequest(library, page)
	req := c.J(referer)

	if ctx != nil {
		req.SetContext(ctx)
	}

	if out != nil {
		req.SetResult(out)
	}

	return req.Get(path)
}

// lazy shortcut for getting data without page= parameter
func (c *Client) ScheduleLibraryRequest(library string, out interface{}, fn ResponseHandler) {
	c.ScheduleLibraryPageRequest(library, 0, out, fn)
}

func (c *Client) ScheduleLibraryPageRequest(library string, page int, out interface{}, fn ResponseHandler) {
	referer, path := c.composeLibraryPageRequest(library, page)

	c.Download(referer, path, out, fn)
}
