package login

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"github.com/justprintit/mmf/api/mmf"
)

func ParseLoginResponse(resp *http.Response) (*Form, error) {
	// read all the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// and allow rereading, just in case
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return ParseLoginDocument(body)
}

func ParseLoginDocument(body []byte) (*Form, error) {
	var form *Form

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// for each <form />
	doc.Find("form").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if f := ParseLoginForm(i, s); f != nil {
			form = f
			return false
		}
		return true
	})

	if form != nil {
		return form, nil
	}

	return nil, errors.New("login form not found")
}

func ParseLoginForm(i int, s *goquery.Selection) *Form {
	f := ParseForm(i, s)
	if f != nil && f.Action != "" {
		if u, err := url.Parse(f.Action); err == nil {
			if u.Path == "/login_check" {
				return f
			}
		}
	}
	return nil
}

func NewLoginRequest(prev *http.Request, form *Form, cred mmf.User) (*http.Request, error) {
	form.EachInput(func(i int, d *Input) {
		switch d.Type {
		case "text":
			d.Value = cred.Username
		case "password":
			d.Value = cred.Password
		}
	})

	return form.NewRequest(prev.URL.String())
}

func LoginSuccessResponse(req *http.Request, resp *http.Response) bool {
	if resp.StatusCode == http.StatusFound {
		return true
	}
	return false
}
