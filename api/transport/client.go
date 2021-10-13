package transport

import (
	"context"
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

	// cancel
	ctx    context.Context
	cancel context.CancelFunc

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

	// cancelation
	if c.cancel == nil {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)

		// preserve context if present
		if ctx = c.ctx; ctx == nil {
			ctx = context.Background()
		}

		ctx, cancel = context.WithCancel(ctx)
		c.ctx = ctx
		c.cancel = cancel
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
