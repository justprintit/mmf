package library

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf"
	"github.com/justprintit/mmf/api/client"
)

type Client struct {
	*client.Client
}

func New(cred mmf.Credentials) *Client {
	return &Client{
		Client: client.New(cred),
	}
}

func NewWithClient(cred mmf.Credentials, hc *http.Client) *Client {
	return &Client{
		Client: client.NewWithClient(cred, hc),
	}
}

func NewWithTransport(cred mmf.Credentials, transport http.RoundTripper) *Client {
	return &Client{
		Client: client.NewWithTransport(cred, transport),
	}
}

func (c *Client) GetLibrary(library string) (*resty.Response, error) {
	path := fmt.Sprintf("/data-library/%s", library)
	return c.J("/library?v=%s", library).Get(path)
}
