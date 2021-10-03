package library

import (
	"context"
	"net/http"
	"time"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/library/json"
	"github.com/justprintit/mmf/api/library/types"
)

const (
	NextUserSharedLibraryUpdate = 2 * time.Minute
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

func (c *Client) refreshUserSharedLibraryFromRequest(ctx context.Context, req *http.Request, d *json.UserSharedLibrary) error {
	// grab Username from Path
	path := req.URL.Path
	if username, err := types.UsernameFromPath(path); err == nil {

		// find User, and Apply data
		u, err := c.library.GetUser(username)
		if err == nil {
			return c.refreshUserSharedLibrary(ctx, u, d)
		}
		return err
	}

	return errors.New("Invalid Path %q", path)
}

func (c *Client) refreshUserSharedLibrary(ctx context.Context, u *types.User, d *json.UserSharedLibrary) error {
	return d.Apply(c.library, u)
}

func (c *Client) refreshTribeSharedGroupFromRequest(ctx context.Context, req *http.Request, d *json.Objects) error {
	return c.refreshSharedGroupFromRequest(ctx, req, d)
}

func (c *Client) refreshUserSharedGroupFromRequest(ctx context.Context, req *http.Request, d *json.Objects) error {
	return c.refreshSharedGroupFromRequest(ctx, req, d)
}

func (c *Client) refreshSharedGroupFromRequest(ctx context.Context, req *http.Request, d *json.Objects) error {
	// grab GroupId from Path
	path := req.URL.Path
	if id, err := types.GroupIdFromPath(path); err == nil {

		// find Group, and Apply data
		g, err := c.library.GetGroup(id)
		if err == nil {
			err = c.refreshSharedGroup(ctx, g, d)
		}
		return err
	}

	return errors.New("Invalid Path %q", path)
}

func (c *Client) refreshSharedGroup(ctx context.Context, g *types.Group, d *json.Objects) error {
	var check errors.ErrorStack

	for _, obj := range d.Items {
		if err := obj.Apply(c.library, nil, g); err != nil {
			check.AppendError(err)
		}
	}

	return check.AsError()
}

func (c *Client) refreshSharedLibrary(ctx context.Context, offset int, users ...json.User) error {
	var check errors.ErrorStack

	for _, w := range users {
		if u, err := w.Apply(c.library); err != nil {
			check.AppendError(err)
		} else if err := c.scheduleUserSharedLibrary(ctx, u); err != nil {
			check.AppendError(err)
		} else if err := c.scheduleUserSharedGroup(ctx, u); err != nil {
			check.AppendError(err)
		}
	}

	return check.AsError()
}

func (c *Client) refreshPurchasesLibrary(ctx context.Context, offset int, objects ...json.Object) error {
	var check errors.ErrorStack

	for _, obj := range objects {
		if err := obj.Apply(c.library, nil, nil); err != nil {
			check.AppendError(err)
		}
	}

	return check.AsError()
}

func (c *Client) refreshPledgesLibrary(ctx context.Context, offset int, objects ...json.Object) error {
	var check errors.ErrorStack

	for _, obj := range objects {
		if err := obj.Apply(c.library, nil, nil); err != nil {
			check.AppendError(err)
		}
	}

	return check.AsError()
}

func (c *Client) refreshTribesLibrary(ctx context.Context, offset int, tribes ...json.Tribe) error {
	var check errors.ErrorStack

	for _, w := range tribes {
		if tribe, err := w.Apply(c.library); err != nil {
			check.AppendError(err)
		} else if err := c.scheduleTribeSharedGroup(ctx, tribe); err != nil {
			check.AppendError(err)
		}
	}

	return check.AsError()
}
