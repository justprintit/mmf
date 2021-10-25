package types

import (
	"net/http"
)

type HttpRequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type HttpRequestDoerFunc func(*http.Request) (*http.Response, error)

func (f HttpRequestDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}
