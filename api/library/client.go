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

func (c *Client) GetLibrary(library string) (*resty.Response, error) {
	path := fmt.Sprintf("/data-library/%s", library)
	return c.J("/library?v=%s", library).Get(path)
}
