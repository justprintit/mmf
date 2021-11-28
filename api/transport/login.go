package transport

import (
	"net/http"

	"github.com/justprintit/mmf/api/mmf"
	"github.com/justprintit/mmf/api/transport/login"
	"github.com/justprintit/mmf/types"
)

const (
	DataLibraryPath = "/data-library"
)

func WithUser(cred mmf.User) ClientOptionFunc {
	return func(c *Client) error {
		c.credentials = cred
		return nil
	}
}

func (c *Client) DataLibraryServer() string {
	return c.ServerJoinPath(DataLibraryPath)
}

// Do makes a request using the login authentication
func (c *Client) NewLoginDoer() types.HttpRequestDoer {

	return &login.AutoLogin{
		Base: &http.Client{
			Transport: c.Transport,
			Jar:       c.Jar,
		},

		User: c.credentials,
	}
}
