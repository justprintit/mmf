package json

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/justprintit/mmf/api/client"
	"github.com/justprintit/mmf/api/library/types"
)

// /data-library/group/{id}
func NewUserSharedGroupRequest(g *types.Group) client.RequestOptions {
	opt := SharedLibraryRequest
	if u := g.User(); u != nil {
		opt.Referer += fmt.Sprintf("&s=all/%s", url.QueryEscape(u.Id()))
	}
	opt.Path = g.GetObjectsURL()
	opt.Result = Objects{}
	return opt
}

func NewUserSharedGroupFromRequest(req *http.Request) client.RequestOptions {
	return client.RequestOptions{
		Accept:  "application/json",
		Referer: req.Header.Get("Referer"),
		Path:    req.URL.Path,
		Method:  req.Method,
		Result:  Objects{},
	}
}
