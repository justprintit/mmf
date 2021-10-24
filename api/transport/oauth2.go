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

	ApiPath = "/api/v2"
)

func (c *Client) OpenAPIServer() string {
	return c.ServerJoinPath(ApiPath)
}

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

func (c *Client) resetOauth2() {
	c.setTokenSource(nil, false)
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
		return c.setTokenSource(c.client.Token(), true)
	} else {
		c.oauth2 = nil
		return c.setTokenSource(nil, false)
	}
}

func (c *Client) setTokenSource(token *oauth2.Token, refresh bool) error {
	if token == nil {
		// no token
		c.ts = nil
		return c.rememberToken(nil)
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
			return c.setTokenSource(token, false)
		}
	} else {
		// allocate TokenSource
		c.ts = c.oauth2.TokenSource(ctx, token)
		return c.rememberToken(token)
	}
}

func (c *Client) rememberToken(token *oauth2.Token) error {

	if token == nil {
		// forget
		c.client.AccessToken = ""
		c.client.RefreshToken = ""
		return c.onNewToken(nil)
	}

	if len(token.RefreshToken) == 0 {
		token.RefreshToken = c.client.RefreshToken
	}

	if token.AccessToken != c.client.AccessToken ||
		token.RefreshToken != c.client.RefreshToken {
		// and remember new Token if needed

		c.client.AccessToken = token.AccessToken
		c.client.RefreshToken = token.RefreshToken
		return c.onNewToken(token)
	}

	return nil
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

	return c.setTokenSource(token, false)
}

// Token() lets us implement oauth2.TokenSource and capture token updates
func (c *Client) Token() (*oauth2.Token, error) {
	if c.ts != nil {
		t, err := c.ts.Token()
		if err == nil {
			err = c.rememberToken(t)
		}
		return t, err
	}
	return nil, ErrTokenNotAvailable
}

// Do makes a request using the oauth2 token
func (c *Client) NewOauth2Doer() HttpRequestDoerFunc {

	rt := &http.Client{
		Transport: &oauth2.Transport{
			Source: c,
			Base:   c.Transport,
		},
		Jar: c.Jar,
	}

	fn := func(req *http.Request) (*http.Response, error) {
		resp, err := rt.Do(req)
		if err == nil {
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				c.resetOauth2()
			}
		}
		return resp, err
	}

	return fn
}
