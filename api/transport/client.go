package transport

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"

	"github.com/justprintit/mmf/api/mmf"
)

const (
	DefaultServer = "https://www.myminifactory.com/"
)

type Client struct {
	Server    string
	Jar       http.CookieJar
	Transport http.RoundTripper

	// scrap
	credentials mmf.User

	done1 chan struct{}
	done2 chan struct{}
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

func (c *Client) Run() {
	c.done1 = make(chan struct{})
	c.done2 = make(chan struct{})

	defer close(c.done2)
	<-c.done1
}

func (c *Client) Abort() {
	close(c.done1)
}

func (c *Client) Wait() {
	<-c.done2
}
