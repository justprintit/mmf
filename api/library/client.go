package library

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf"
	"github.com/justprintit/mmf/api/client"
)

type Client struct {
	client.Client

	DataDir string
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
		if lopt, ok := opt.(LibraryClientOption); ok {
			if err := lopt.ApplyLibrary(c); err != nil {
				return nil, err
			}
		} else if err := opt.Apply(&c.Client); err != nil {
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

type LibraryClientOption interface {
	client.ClientOption

	ApplyLibrary(*Client) error
}

type LibraryClientOptionFunc func(*Client) error

func (f LibraryClientOptionFunc) ApplyLibrary(c *Client) error {
	return f(c)
}

func (f LibraryClientOptionFunc) Apply(c *client.Client) error {
	return errors.New("Invalid usage of %T", f)
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

func WithDataDir(datadir string) client.ClientOption {
	return LibraryClientOptionFunc(func(c *Client) error {
		if datadir == "" {
			datadir = "."
		}
		c.DataDir = datadir
		return nil
	})
}
