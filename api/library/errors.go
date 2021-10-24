package library

import (
	"net/http"
)

type ResponseError struct {
	*http.Response

	Body []byte
}

func (resp ResponseError) Error() string {
	return http.StatusText(resp.Status())
}

func (resp ResponseError) Status() int {
	if p := resp.Response; p != nil {
		return resp.StatusCode
	} else {
		return 0
	}
}

func ErrResponse(resp *http.Response, body []byte) error {
	return ResponseError{resp, body}
}
