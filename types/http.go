package types

import (
	"net/http"

	"github.com/motemen/go-loghttp"
	"golang.org/x/oauth2"
)

type HttpRequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type HttpRequestDoerFunc func(*http.Request) (*http.Response, error)

func (f HttpRequestDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type TransportUnwrapper interface {
	Unwrap() *http.Transport
}

func TransportUnwrap(v interface{}) *http.Transport {
	for v != nil {

		switch t := v.(type) {
		case *http.Transport:
			// direct
			return t
		case TransportUnwrapper:
			// semi-direct
			return t.Unwrap()

		case *http.Client:
			// next http.RoundTripper
			v = t.Transport
		case *oauth2.Transport:
			// next http.RoundTripper
			v = t.Base
		case *loghttp.Transport:
			// next http.RoundTripper
			v = t.Transport

		default:
			// give up
			v = nil
		}
	}
	return nil
}
