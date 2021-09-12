package client

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/go-resty/resty/v2"
	"github.com/json-iterator/go"
	"golang.org/x/net/publicsuffix"

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

	c.JSONMarshal = jsoniter.Marshal
	c.JSONUnmarshal = jsoniter.Unmarshal

	c.SetHostURL(DefaultHost)

	// inject auto-login middleware
	hc := c.GetClient()
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

func NewWithOptions(options ...ClientOption) (*Client, error) {
	c := new(Client)
	for _, opt := range options {
		if err := opt.Apply(c); err != nil {
			return nil, err
		}
	}
	c.Init(mmf.Credentials{}, nil)
	return c, nil
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

type ClientOption interface {
	Apply(c *Client) error
}

type ClientOptionFunc func(*Client) error

func (f ClientOptionFunc) Apply(c *Client) error {
	return f(c)
}

func WithCredentials(cred mmf.Credentials) ClientOption {
	return ClientOptionFunc(func(c *Client) error {
		c.Credentials = cred
		return nil
	})
}

func WithTransport(transport http.RoundTripper) ClientOption {
	return ClientOptionFunc(func(c *Client) error {
		if c.Client == nil {
			c.Client = resty.New()
		}
		c.Client.SetTransport(transport)
		return nil
	})
}

func WithCookieJar(jar http.CookieJar) ClientOption {
	return ClientOptionFunc(func(c *Client) error {
		if jar == nil {
			var err error

			jar, err = cookiejar.New(&cookiejar.Options{
				PublicSuffixList: publicsuffix.List,
			})

			if err != nil {
				return err
			}

		}

		if c.Client == nil {
			c.Client = resty.New()
		}

		c.SetCookieJar(jar)
		return nil
	})
}
