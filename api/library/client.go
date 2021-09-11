package library

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf"
	"github.com/justprintit/mmf/api/client"
)

type Client struct {
	client.Client
}

func (c *Client) Init(cred mmf.Credentials, rc *resty.Client) *Client {
	c.Client.Init(cred, rc)
	return c
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

func NewWithOptions(options ...client.ClientOption) (*Client, error) {
	c := new(Client)
	for _, opt := range options {
		if err := opt.Apply(&c.Client); err != nil {
			return nil, err
		}
	}
	c.Client.Init(mmf.Credentials{}, nil)
	return c, nil
}

func (c *Client) GetLibrary(library string) (*resty.Response, error) {
	path := fmt.Sprintf("/data-library/%s", library)
	return c.J("/library?v=%s", library).Get(path)
}

func WithCredentials(cred mmf.Credentials) client.ClientOption {
	return client.WithCredentials(cred)
}

func WithTransport(transport http.RoundTripper) client.ClientOption {
	return client.WithTransport(transport)
}

func WithCookieJar(jar http.CookieJar) client.ClientOption {
	return client.WithCookieJar(jar)
}
