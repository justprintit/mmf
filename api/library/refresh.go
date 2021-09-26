package library

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/core/typeconv"

	"github.com/justprintit/mmf/api/client"
	"github.com/justprintit/mmf/api/library/json"
	"github.com/justprintit/mmf/api/library/types"
)

const (
	NextUserSharedLibraryUpdate = 2 * time.Minute
	NextUserSharedGroupsUpdate  = 2 * time.Minute
	NextGroupObjectsUpdate      = 2 * time.Minute
)

// Reload persistent Library data
func (c *Client) Reload() error {
	l, err := c.store.Load()
	if err == nil {
		c.library = l
	}
	return err
}

// Stores library data persistently
func (c *Client) Commit() error {
	return c.store.Store(c.library)
}

func (c *Client) scheduleUserSharedLibrary(ctx context.Context, u *types.User) error {
	select {
	case <-ctx.Done():
		// cancelled
		return ctx.Err()
	default:
		t := time.Now()
		if t.After(u.NextUserSharedLibraryUpdate) {
			u.NextUserSharedLibraryUpdate = t.Add(NextUserSharedLibraryUpdate)

			c.Download(NewUserSharedLibraryRequest(u), refreshUserSharedLibraryCallback)
		}
		return nil
	}
}

func NewUserSharedLibraryRequest(u *types.User) client.RequestOptions {
	opt := SharedLibraryRequest
	opt.Referer += fmt.Sprintf("&s=all/%s", url.QueryEscape(u.Username))
	opt.Path += fmt.Sprintf("/%s", url.PathEscape(u.Username))
	opt.Result = json.UserSharedLibrary{}
	return opt
}

func UserSharedLibraryResult(resp *resty.Response) *json.UserSharedLibrary {
	if out := resp.Result(); out != nil {
		return out.(*json.UserSharedLibrary)
	}
	return nil
}

func refreshUserSharedLibraryCallback(c *Client, ctx context.Context, resp *resty.Response) error {
	if p := UserSharedLibraryResult(resp); p != nil {

		// grab Username from Path
		path := resp.RawResponse.Request.URL.Path
		username := strings.TrimPrefix(path, SharedLibraryRequest.Path)
		if username != path && len(username) > 1 && username[0] == '/' {
			if username, err := url.PathUnescape(username[1:]); err == nil {

				// find User, and Apply data
				u, err := c.library.GetUser(username)
				if err == nil {
					err = p.Apply(c.library, u)
				}
				return err

			}
		}

		return errors.New("Invalid Path %q", path)
	}
	return nil
}

func (c *Client) scheduleUserSharedGroups(ctx context.Context, u *types.User) error {
	select {
	case <-ctx.Done():
		// cancelled
		return ctx.Err()
	default:
		t := time.Now()
		if t.After(u.NextUserSharedGroupsUpdate) {
			u.NextUserSharedGroupsUpdate = t.Add(NextUserSharedGroupsUpdate)

			for _, g := range u.GroupsAll() {
				if t.After(g.NextGroupObjectsUpdate) {
					g.NextGroupObjectsUpdate = t.Add(NextGroupObjectsUpdate)

					c.Download(NewGroupObjectsRequest(g), refreshUserSharedGroupsCallback)
				}
			}
		}
		return nil
	}
}

func NewGroupObjectsRequest(g *types.Group) client.RequestOptions {
	opt := SharedLibraryRequest
	opt.Referer += fmt.Sprintf("&s=all/%s", url.QueryEscape(g.User().Username))
	opt.Path = g.GetObjectsURL()
	opt.Result = json.Objects{}
	return opt
}

func GroupObjectsResult(resp *resty.Response) *json.Objects {
	if out := resp.Result(); out != nil {
		return out.(*json.Objects)
	}
	return nil
}

func refreshUserSharedGroupsCallback(c *Client, ctx context.Context, resp *resty.Response) error {
	if p := GroupObjectsResult(resp); p != nil {

		// grab GroupId from Path
		path := resp.RawResponse.Request.URL.Path
		s := strings.TrimPrefix(path, "/data-library/group/")
		if id, err := typeconv.AsInt(s); err == nil {

			// find Group, and Appy data
			g, err := c.library.GetGroup(id)
			if err == nil {
				err = p.Apply(c.library, g.User(), g)
			}
			return err
		}

		return errors.New("Invalid Path %q", path)
	}
	return nil
}

func (c *Client) refreshSharedLibrary(ctx context.Context, offset int, users ...json.User) error {
	for _, w := range users {
		if u, err := w.Apply(c.library); err != nil {
			log.Println(err)
		} else if err := c.scheduleUserSharedLibrary(ctx, u); err != nil {
			log.Println(err)
		} else if err := c.scheduleUserSharedGroups(ctx, u); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (c *Client) refreshPurchasesLibrary(ctx context.Context, offset int, objects ...json.Object) error {
	for _, obj := range objects {
		if err := obj.Apply(c.library, nil, nil); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (c *Client) refreshPledgesLibrary(ctx context.Context, offset int, objects ...json.Object) error {
	for _, obj := range objects {
		if err := obj.Apply(c.library, nil, nil); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (c *Client) refreshTribesLibrary(ctx context.Context, offset int, tribes ...json.Tribe) error {
	i := offset
	for _, u := range tribes {
		i++

		log.Printf("Tribe.%v: %s (%v)", i, u.Name, u.Id)
	}
	return nil
}
