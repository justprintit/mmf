package login

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

type LoginRedirectError struct {
	req *http.Request
}

func (e *LoginRedirectError) Error() string {
	return fmt.Sprintf("Redirected to %s", e.req.URL.String())
}

func (e *LoginRedirectError) Unwrap() *http.Request {
	return e.req
}

func NewLoginRedirectError(req *http.Request, via []*http.Request) error {

	if b, _ := httputil.DumpRequest(req, false); len(b) > 0 {
		os.Stderr.Write(b)
	}

	if req.URL.Path == "/login" {
		return &LoginRedirectError{
			req: req,
		}
	}
	return nil
}
