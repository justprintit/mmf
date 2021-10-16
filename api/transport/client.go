package transport

import (
	"context"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
	"golang.org/x/oauth2"

	"github.com/justprintit/mmf/api/mmf"
)

const (
	SequentialQueue QueueIndex = iota
	DownloadQueue

	QueuesCount int = iota

	DefaultServer = "https://www.myminifactory.com/"
)

type Client struct {
	wq WorkQueue

	Server    string
	Jar       http.CookieJar
	Transport http.RoundTripper

	// scrap
	credentials mmf.User
	// oauth2
	callback string
	client   mmf.Client
	oauth2   *oauth2.Config
	ts       oauth2.TokenSource

	events ClientEvents
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

	// queue
	if c.wq.Len() == 0 {
		c.wq.Init(c, QueuesCount)
	}

	// oauth2
	if c.oauth2 == nil {
		if err := c.setOauth2(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Context() context.Context {
	if ctx := c.wq.Context(); ctx != nil {
		return ctx
	}
	return context.Background()
}

func (c *Client) Run(ctx context.Context, n int32) {
	c.wq.Start(ctx, 1, n)
	c.wq.Wait()
}

func (c *Client) Cancel() {
	c.wq.Cancel()
}

func (c *Client) Wait() {
	c.wq.Wait()
}

func (c *Client) Schedule(f QueueFunc, v interface{}) {
	c.wq.Add(SequentialQueue, f, v)
}

func (c *Client) ScheduleDownloader(f QueueFunc, v interface{}) {
	c.wq.Add(DownloadQueue, f, v)
}
