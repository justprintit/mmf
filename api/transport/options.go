package transport

import (
	"github.com/justprintit/mmf/api/mmf"
	"net/http"
)

type ClientOptionFunc func(*Client) error

func (f ClientOptionFunc) Apply(c *Client) error {
	return f(c)
}

type ClientOption interface {
	Apply(*Client) error
}

func NewClientWithOptions(options ...ClientOption) (*Client, error) {
	c := &Client{}

	for _, opt := range options {
		if err := opt.Apply(c); err != nil {
			return nil, err
		}
	}

	if err := c.SetDefaults(); err != nil {
		return nil, err
	}
	return c, nil
}

func WithTransport(rt http.RoundTripper) ClientOptionFunc {
	return func(c *Client) error {
		c.Transport = rt
		return nil
	}
}

func WithCookieJar(jar http.CookieJar) ClientOptionFunc {
	return func(c *Client) error {
		c.Jar = jar
		return nil
	}
}

func WithUser(cred mmf.User) ClientOptionFunc {
	return func(c *Client) error {
		c.credentials = cred
		return nil
	}
}
