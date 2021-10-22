package transport

import (
	"net/http"
)

type HttpRequestDoerFunc func(*http.Request) (*http.Response, error)

func (f HttpRequestDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}
