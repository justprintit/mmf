package client

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

func (c *Client) Init(cred mmf.Credentials, rc *resty.Client) *Client {
	if rc != nil {
		c.Client = rc
	} else if c.Client == nil {
		c.Client = resty.New()
	}

	if len(cred.Username) > 0 {
		c.Credentials = cred
	}

	c.SetHostURL(DefaultHost)

	// inject auto-login middleware
	hc := rc.GetClient()
	return c.SetTransport(hc.Transport)
}

func New(cred mmf.Credentials) *Client {
	return new(Client).Init(cred, nil)
}

func NewWithClient(cred mmf.Credentials, hc *http.Client) *Client {
	rc := resty.NewWithClient(hc)
	return new(Client).Init(cred, rc)
}

func NewWithTransport(cred mmf.Credentials, transport http.RoundTripper) *Client {
	rc := resty.New().SetTransport(transport)
	return new(Client).Init(cred, rc)
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
