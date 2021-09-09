package library

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/html"
)

const (
	MaxLoginAttempts = 2
)

type LoginForm struct {
	Action  string
	EncType string
	Method  string
	Referer string
	Fields  map[string]string
}

func (form *LoginForm) Request(c *resty.Client) (*resty.Response, error) {
	if len(form.Action) == 0 {
		form.Action = form.Referer
	}
	if len(form.EncType) == 0 {
		form.EncType = "application/x-www-form-urlencoded"
	}
	if len(form.Method) == 0 {
		form.Method = "GET"
	} else {
		form.Method = strings.ToUpper(form.Method)
	}

	req := c.R()
	req.SetHeaders(map[string]string{
		"Cache-Control": "no-cache",
		"Content-Type":  form.EncType,
		"Referer":       form.Referer,
	})
	req.SetFormData(form.Fields)
	return req.Execute(form.Method, form.Action)
}

type LoginParseState func(z *html.Tokenizer, form *LoginForm) (LoginParseState, error)

type AutoLogin struct {
	attempts  int
	client    *Client
	Transport http.RoundTripper
}

func (m *AutoLogin) transport() http.RoundTripper {
	if m != nil && m.Transport != nil {
		return m.Transport
	} else {
		return http.DefaultTransport
	}
}

func (m *AutoLogin) roundTrip(req *http.Request) (resp *http.Response, err error) {
	return m.transport().RoundTrip(req)
}

func (m *AutoLogin) NewClient() *http.Client {
	hc := m.client.GetClient()

	return &http.Client{
		Transport:     m.transport(),
		CheckRedirect: hc.CheckRedirect,
		Jar:           hc.Jar,
		Timeout:       hc.Timeout,
	}
}

func (m *AutoLogin) IsLoginRedirect(resp *http.Response) bool {
	if resp.StatusCode == http.StatusFound {
		u, err := url.Parse(resp.Header.Get("Location"))

		if err == nil && u.Path == "/login" {
			return true
		}
	}
	return false
}

func (m *AutoLogin) processFormToken(t html.Token, form *LoginForm) (LoginParseState, error) {
	var actionURL *url.URL
	var tagId string
	var temp LoginForm

	// validate attributes
	for _, v := range t.Attr {
		switch v.Key {
		case "action":
			if u, err := url.Parse(v.Val); err == nil {
				actionURL = u
				temp.Action = v.Val
			}
		case "id":
			tagId = strings.ToLower(strings.TrimSpace(v.Val))
		case "method":
			temp.Method = strings.ToLower(strings.TrimSpace(v.Val))
		case "enctype":
			temp.EncType = strings.ToLower(strings.TrimSpace(v.Val))
		}
	}

	if actionURL == nil || strings.Index(tagId, "login") != -1 || strings.Index(actionURL.Path, "login") != -1 {
		// login form
		form.Action = temp.Action
		form.EncType = temp.EncType
		form.Method = temp.Method

		// find fields
		return m.parseForm, nil
	}

	return m.parseInitial, nil
}

func (m *AutoLogin) processFormField(t html.Token, form *LoginForm) {
	var id, name, value string

	for _, v := range t.Attr {
		switch v.Key {
		case "id":
			id = v.Val
		case "name":
			name = v.Val
		case "value":
			value = v.Val
		}
	}

	if len(name) > 0 {
		if len(value) == 0 {
			for _, k := range []string{"username"} {
				if strings.Index(name, k) != -1 || strings.Index(id, k) != -1 {
					value = m.client.Credentials.Username
					break
				}
			}
		}

		if len(value) == 0 {
			for _, k := range []string{"password"} {
				if strings.Index(name, k) != -1 || strings.Index(id, k) != -1 {
					value = m.client.Credentials.Password
					break
				}
			}
		}

		form.Fields[name] = value
	}
}

func (m *AutoLogin) parseForm(z *html.Tokenizer, form *LoginForm) (LoginParseState, error) {
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// EOF
			return nil, nil
		case html.EndTagToken:
			t := z.Token()
			if t.Data == "form" {
				// done
				return nil, nil
			}
		case html.SelfClosingTagToken:
			t := z.Token()
			switch t.Data {
			case "input":
				m.processFormField(t, form)
			default:
			}
		}
	}
}

func (m *AutoLogin) parseInitial(z *html.Tokenizer, form *LoginForm) (LoginParseState, error) {
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// EOF
			return nil, nil
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "form" {
				return m.processFormToken(t, form)
			}
		}
	}
}

func (m *AutoLogin) makeLoginRequest(c *resty.Client, referer string, resp *resty.Response) (*resty.Response, error) {

	form := &LoginForm{
		Referer: referer,
		Fields:  make(map[string]string),
	}

	// tokenizer
	z := html.NewTokenizer(bytes.NewReader(resp.Body()))
	// initial state
	s := LoginParseState(m.parseInitial)

	for s != nil {
		var err error

		// next state
		if s, err = s(z, form); err != nil {
			return nil, err
		}
	}

	return form.Request(c)
}

func (m *AutoLogin) tryLogin(prev *http.Request, redirect *http.Response) (*http.Response, error, bool) {
	// don't do this forever
	if m.attempts > 0 {
		m.attempts--
	} else {
		return nil, errors.New("too many login requests"), false
	}

	// new client
	c := resty.NewWithClient(m.NewClient())
	c.SetHostURL(m.client.HostURL)

	// get /login HTML page
	loginLocation := redirect.Header.Get("location")
	resp, err := c.R().SetHeader("Referer", prev.URL.String()).Get(loginLocation)
	if err == nil && resp.StatusCode() == http.StatusOK {
		// compose POST request
		resp, err = m.makeLoginRequest(c, loginLocation, resp)
		if err == nil && resp.StatusCode() == http.StatusOK {
			// success
			return nil, nil, true
		}
	}

	return resp.RawResponse, err, false
}

func (m *AutoLogin) RoundTrip(req *http.Request) (*http.Response, error) {
	for {
		resp, err := m.roundTrip(req)
		if err != nil || !m.IsLoginRedirect(resp) {
			// done
			return resp, err
		} else if resp, err, ok := m.tryLogin(req, resp); !ok {
			// login failed. bye.
			return resp, err
		} else {
			// success. apply new cookies and retry
			jar := m.client.GetClient().Jar
			for _, cookie := range jar.Cookies(req.URL) {
				req.AddCookie(cookie)
			}
		}
	}
}

func (c *Client) SetTransport(next http.RoundTripper) *Client {
	if next == nil {
		next = http.DefaultTransport
	}

	hc := c.GetClient()
	hc.Transport = &AutoLogin{
		attempts:  MaxLoginAttempts,
		client:    c,
		Transport: next,
	}
	return c
}
