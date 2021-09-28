package json

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/justprintit/mmf/api/client"
	"github.com/justprintit/mmf/api/library/types"
)

type UserSharedLibrary struct {
	Objects Objects `json:"objects,omitempty"`
	Groups  Groups  `json:"groups,omitempty"`
}

func (p *UserSharedLibrary) Apply(d *types.Library, u *types.User) error {
	if err := p.Groups.Apply(d, u); err != nil {
		return err
	}
	if err := p.Objects.Apply(d, u, nil); err != nil {
		return err
	}
	return nil
}

// /data-library/shared/{username}
func NewUserSharedLibraryRequest(u *types.User) client.RequestOptions {
	referer := SharedLibraryRequest.Referer
	referer = fmt.Sprintf("%s&s=all/%s", referer, url.QueryEscape(u.Username))

	return client.RequestOptions{
		Accept:  "application/json",
		Referer: referer,
		Path:    u.GetSharedGroupsURL(),
		Method:  http.MethodGet,
		Result:  UserSharedLibrary{},
	}
}

func NewUserSharedLibraryFromRequest(req *http.Request) client.RequestOptions {
	return client.RequestOptions{
		Accept:  "application/json",
		Referer: req.Header.Get("Referer"),
		Path:    req.URL.Path,
		Method:  req.Method,
		Result:  UserSharedLibrary{},
	}
}

func UserSharedLibraryPages(p *UserSharedLibrary) *client.Pagination {
	p1 := ObjectsPages(&p.Objects)
	p2 := GroupsPages(&p.Groups)

	d1 := p1.Total - p1.Size
	d2 := p2.Total - p2.Size

	if d1 > d2 {
		return p1
	} else {
		return p2
	}
}
