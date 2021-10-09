package transport

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

const (
	DefaultServer = "https://www.myminifactory.com/"
)

type Client struct {
	Server string
	Jar    http.CookieJar
}

func (c *Client) SetDefaults() error {

	// Server
	if c.Server == "" {
		c.Server = DefaultServer
	}

	// CookieJar
	if c.Jar == nil {
		jar, err := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		})
		if err != nil {
			return err
		}
		c.Jar = jar
	}

	return nil
}
