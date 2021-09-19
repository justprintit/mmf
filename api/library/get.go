package library

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

func (c *Client) GetWithContext(ctx context.Context, path string) (*resty.Response, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	req := c.R("/library")
	req.SetContext(ctx)

	if u.Host == "" && strings.IndexRune(u.Path, '/') == -1 {
		// data-library
		u.Path = "/data-library/" + u.Path
	}

	if strings.HasPrefix(u.Path, "/data-library/") {
		req.SetHeader("Accept", "application/json")
	}

	return req.Get(u.String())
}
