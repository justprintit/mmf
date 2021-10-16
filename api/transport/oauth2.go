package transport

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"time"

	"golang.org/x/oauth2"

	"go.sancus.dev/web/errors"

	"github.com/justprintit/mmf/api/mmf"
	"github.com/justprintit/mmf/util"
)

const (
	CodeExchangeTimeout = 2 * time.Second
	RandomStateLength   = 32
)

func WithOauth2(cred mmf.Client, base string, callback string) ClientOptionFunc {
	// callback = base + path
	if u, err := url.Parse(base); err != nil {
		panic(err)
	} else {
		u.Path = path.Clean(path.Join(u.Path, callback))
		callback = u.String()
	}

	return func(c *Client) error {
		c.callback = callback
		c.client = cred
		return c.setOauth2()
	}
}

func (c *Client) setOauth2() error {
	if len(c.callback) > 0 && c.client.Ok() {
		// *oauth2.Config
		scopes := []string{
			string(mmf.BasicScope),
			string(mmf.DownloadScope),
		}

		oc, err := c.client.NewOauth2(c.callback, scopes...)
		if err != nil {
			return err
		}

		c.oauth2 = oc

		// TokenSource
		return c.setToken(c.client.Token(), true)
	} else {
		c.oauth2 = nil
		return c.setToken(nil, false)
	}
}

func (c *Client) setToken(token *oauth2.Token, refresh bool) error {
	if token == nil {
		// no token
		c.ts = nil
		return nil
	}

	ctx := c.Context()

	if refresh {
		// token is incomplete, refresh it right away
		ctx, cancel := context.WithTimeout(ctx, CodeExchangeTimeout)
		defer cancel()

		ts := c.oauth2.TokenSource(ctx, token)
		if token, err := ts.Token(); err != nil {
			return err
		} else {
			// and allocate TokenSource
			return c.setToken(token, false)
		}
	} else {
		// allocate TokenSource
		c.ts = c.oauth2.TokenSource(ctx, token)

		// and remember new Token if needed
		if len(token.RefreshToken) == 0 {
			token.RefreshToken = c.client.RefreshToken
		}

		if token.AccessToken != c.client.AccessToken ||
			token.RefreshToken != c.client.RefreshToken {
			c.client.AccessToken = token.AccessToken
			c.client.RefreshToken = token.RefreshToken
		}

		return c.onNewToken(token)
	}
}

func (c *Client) RedirectHandler(rw http.ResponseWriter, req *http.Request) error {
	if b, _ := httputil.DumpRequest(req, true); len(b) > 0 {
		os.Stderr.Write(b)
	}

	if c.oauth2 == nil {
		return errors.ErrMissingField("%s.%s", "oauth2", "Config")
	}

	// generate random state
	// TODO: store state in cookie
	state, err := util.RandomString(RandomStateLength)
	if err != nil {
		return err
	}

	// and redirect to AuthCodeURL
	return errors.NewSeeOther(c.oauth2.AuthCodeURL(state))
}

func (c *Client) CallbackHandler(rw http.ResponseWriter, req *http.Request) error {
	if b, _ := httputil.DumpRequest(req, true); len(b) > 0 {
		os.Stderr.Write(b)
	}

	if c.oauth2 == nil {
		return errors.ErrMissingField("%s.%s", "oauth2", "Config")
	}

	q := req.URL.Query()
	state := q.Get("state")
	if len(state) != RandomStateLength {
		// wrong length
		return errors.ErrInvalidArgument("state")
	}

	// TODO: validate state against cookie

	// Exchange
	ctx, cancel := context.WithTimeout(req.Context(), CodeExchangeTimeout)
	defer cancel()

	token, err := c.oauth2.Exchange(ctx, q.Get("code"))
	if err != nil {
		return errors.AsInvalidArgumentError(err, "code")
	}

	return c.setToken(token, false)
}