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

// Transport implements HttpRequestDoer, http.RoundTripper, and TransportUnwrapper
type Transport struct {
	doer HttpRequestDoer
	do   HttpRequestDoerFunc
}

func NewTransport(doer HttpRequestDoer, do HttpRequestDoerFunc) *Transport {
	if doer != nil || do != nil {
		return &Transport{doer, do}
	} else {
		return nil
	}
}

func (rt *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.do != nil {
		return rt.do(req)
	} else if rt.doer != nil {
		return rt.doer.Do(req)
	} else {
		return nil, nil
	}
}

func (rt *Transport) Do(req *http.Request) (*http.Response, error) {
	return rt.RoundTrip(req)
}

func (rt *Transport) Unwrap() *http.Transport {
	return TransportUnwrap(rt.doer)
}
