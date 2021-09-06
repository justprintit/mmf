package library

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

const (
	DefaultHost = "https://www.myminifactory.com/"
)

type Client struct {
	*resty.Client
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
