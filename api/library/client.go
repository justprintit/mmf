package library

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const (
	DefaultHost = "https://www.myminifactory.com/"
)

type Client struct {
	*resty.Client

	TraceEnabled bool
}

func newClient(rc *resty.Client) *Client {
	c := &Client{
		Client: rc,
	}
	c.SetHostURL(DefaultHost)
	return c
}

func New() *Client {
	rc := resty.New()
	return newClient(rc)
}

func NewWithClient(hc *http.Client) *Client {
	rc := resty.NewWithClient(hc)
	return newClient(rc)
}

func NewWithTransport(transport http.RoundTripper) *Client {
	rc := resty.New().SetTransport(transport)
	return newClient(rc)
}

func (c *Client) R(referer string, args ...interface{}) *resty.Request {
	req := c.Client.R()

	// TraceInfo
	if c.TraceEnabled {
		req.EnableTrace()
	}

	// Referer
	if len(args) > 0 {
		referer = fmt.Sprintf(referer, args...)
	}
	if len(referer) == 0 {
		referer = "/"
	} else if referer[0] != '/' {
		referer = "/" + referer
	}
	req.SetHeader("Referer", c.HostURL+referer)

	return req
}

func (c *Client) J(referer string, args ...interface{}) *resty.Request {
	req := c.R(referer, args...)
	req.SetHeader("Accept", "application/json")
	return req
}

func (c *Client) Head(library string) (*resty.Response, error) {
	path := fmt.Sprintf("/data-library/%s", library)
	return c.J("/library?v=%s", library).Head(path)
}

func (c *Client) Get(library string) (*resty.Response, error) {
	path := fmt.Sprintf("/data-library/%s", library)
	return c.J("/library?v=%s", library).Get(path)
}
