package login

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"

	"go.sancus.dev/web/errors"

	"github.com/justprintit/mmf/api/mmf"
	"github.com/justprintit/mmf/types"
)

const (
	MaxRedirect = 10
)

type AutoLogin struct {
	Base *http.Client
	User mmf.User
}

func (m *AutoLogin) CloneBase(redirect func(*http.Request, []*http.Request) error) (*http.Client, error) {
	// make sure we have a Base client
	if m.Base == nil {
		m.Base = &http.Client{}
	}

	// with a cookiejar
	if m.Base.Jar == nil {
		jar, err := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		})
		if err != nil {
			return nil, err
		}
		m.Base.Jar = jar
	}

	// redirect policy
	if redirect == nil {
		// don't follow redirects
		redirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// shallow copy with custom redirect policy
	c := *m.Base
	c.CheckRedirect = redirect
	return &c, nil
}

func (m *AutoLogin) newLoginRequest(resp *http.Response) (*http.Request, error) {
	form, err := ParseLoginResponse(resp)
	if err != nil {
		return nil, err
	}

	return NewLoginRequest(resp.Request, form, m.User)
}

func (m *AutoLogin) tryLogin(redir *LoginRedirectError) (*http.Response, error, bool) {
	base, err := m.CloneBase(nil)
	if err != nil {
		return nil, err, false
	}

	// get login page
	req := redir.Unwrap()
	resp, err := base.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		// parse form and prepare request
		req, err = m.newLoginRequest(resp)
		if err != nil {
			goto fail
		}

		// attempt login
		resp, err = base.Do(req)
		if err == nil && LoginSuccessResponse(req, resp) {
			// success
			return nil, nil, true
		}
	}

fail:
	return resp, err, false
}

func (m *AutoLogin) redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	if err := NewLoginRedirectError(req, via); err != nil {
		// login redirect detected
		return err
	} else if m.Base != nil && m.Base.CheckRedirect != nil {
		// given base had its own redirect policy
		return m.Base.CheckRedirect(req, via)
	} else if len(via) >= MaxRedirect {
		// limit number of redirects as the default standard client does
		return errors.New("stopped after %v redirects", len(via))
	} else {
		// accept
		return nil
	}
}

func (m *AutoLogin) roundTrip(req *http.Request) (*http.Response, error) {

	// flag to tryLogin only once
	loop := false

	// prepare base client with custom redirect policy
	base, err := m.CloneBase(m.redirectPolicyFunc)
	if err != nil {
		return nil, err
	}

	// and try the request
	for {
		resp, err := base.Do(req)
		if redir, ok := errors.Unwrap(err).(*LoginRedirectError); !ok {
			// done
			return resp, err
		} else if loop {
			// we already tried to login, sorry
			return nil, errors.New("login loop")
		} else {
			// try to login only once
			loop = true

			if resp, err, ok := m.tryLogin(redir); !ok {
				// login failed. bye
				return resp, err
			} else {
				// success. apply new cookies and retry
				jar := base.Jar
				for _, cookie := range jar.Cookies(req.URL) {
					req.AddCookie(cookie)
				}
			}
		}
	}
}

// http.RoundTripper
func (m *AutoLogin) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

// types.HttpRequestDoer
func (m *AutoLogin) Do(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

// types.TransportUnwrapper
func (m *AutoLogin) Unwrap() *http.Transport {
	return types.TransportUnwrap(m.Base)
}
