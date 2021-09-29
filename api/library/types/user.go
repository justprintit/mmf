package types

import (
	"net/url"
	"strings"
	"time"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/util"
)

type User struct {
	entry `json:"-"`

	NextUserSharedLibraryUpdate time.Time `json:"-"`

	Username string
	Name     string
	Avatar   string   `json:",omitempty"`
	Groups   []*Group `json:",omitempty"`
}

func (u *User) GetSharedGroupsURL() string {
	return "/data-library/shared/" + url.PathEscape(u.Username)
}

func (u *User) SanitizedName() string {
	return util.Sanitize(u.Name)
}

func (u *User) Path() string {
	return u.SanitizedName()
}

func (u *User) GroupsAll() []*Group {
	u.entry.Lock()
	defer u.entry.Unlock()

	groups := make([]*Group, 0, len(u.Groups))
	for _, g := range u.Groups {
		all := g.groupsAll()
		groups = append(groups, all...)
	}
	return groups
}

func (u *User) updateName(s string) {
	if len(u.Name) == 0 {
		u.updateString("Name", &u.Name, s)
	}
}

func (u *User) updateAvatar(s string) {
	if s = strings.TrimSpace(s); len(s) > 0 {
		u.updateString("Avatar", &u.Avatar, s)
	}
}

func (u *User) updateString(field string, v *string, s string) {
	before := *v
	after := strings.TrimSpace(s)
	if before != after {
		*v = after
		u.entry.OnUserUpdate(u, field, before, after)
	}
}

func (w *Library) GetUser(user string) (*User, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if u, ok := w.User[user]; ok {
		return u, nil
	} else {
		err := errors.New("%s[%q]: Not Found", "User", user)
		return nil, err
	}
}

func (w *Library) AddUser(u *User, merge bool) (*User, error) {
	var check errors.ErrorStack
	var u0 *User
	var ok bool

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.User == nil {
		w.User = make(map[string]*User, 1)
	}

	if user := u.Username; len(user) == 0 {
		check.MissingField("Username")
	} else if u0, ok = w.User[user]; ok {
		// exists
		if !merge {
			err := errors.New("%s[%q]: Already exists", "User", user)

			check.AppendError(err)
		} else {
			// merge
			u0.updateName(u.Name)
			u0.updateAvatar(u.Avatar)

			// merge groups
			for _, g := range u.Groups {
				if _, err := u0.addGroup(g, true); err != nil {
					check.AppendWrapped(err, "%s[%q]", "User", user)
				}
			}
		}
	} else {
		// new
		name := strings.TrimSpace(u.Name)
		if len(name) == 0 {
			name = u.Username
		}

		u0 = &User{
			Username: u.Username,
			Name:     name,
			Avatar:   u.Avatar,
		}

		// register
		u0.Library = w
		w.User[user] = u0

		w.OnNewUser(u0)

		// add groups
		for _, g := range u.Groups {
			if _, err := u0.addGroup(g, false); err != nil {
				check.AppendWrapped(err, "%s[%q]", "User", user)
			}
		}
	}

	if err := check.AsError(); err != nil {
		w.OnError(err)
		return nil, err
	}

	return u0, nil
}

// UsernameFromURL() attempts to extract the username from a URL
func UsernameFromURL(s string) (string, error) {

	if u, err := url.Parse(s); err == nil {
		if _, p, n := util.NextInPath(u.Path,
			"/users/",
			"/data-library/shared/",
		); n > 0 {
			return p, nil
		}
	}

	return "", ErrInvalidPath(s)
}
