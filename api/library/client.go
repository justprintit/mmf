package library

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf"
)

const (
	DefaultHost = "https://www.myminifactory.com/"
)

type Client struct {
	*resty.Client

	Credentials  mmf.Credentials
	TraceEnabled bool
}

func newClient(cred mmf.Credentials, rc *resty.Client) *Client {
	c := &Client{
		Client:      rc,
		Credentials: cred,
	}
	c.SetHostURL(DefaultHost)

	// inject auto-login middleware
	hc := rc.GetClient()
	return c.SetTransport(hc.Transport)
}

func New(cred mmf.Credentials) *Client {
	rc := resty.New()
	return newClient(cred, rc)
}

func NewWithClient(cred mmf.Credentials, hc *http.Client) *Client {
	rc := resty.NewWithClient(hc)
	return newClient(cred, rc)
}

func NewWithTransport(cred mmf.Credentials, transport http.RoundTripper) *Client {
	rc := resty.New().SetTransport(transport)
	return newClient(cred, rc)
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
