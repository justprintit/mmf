package types

import (
	"net/url"
	"strings"

	"go.sancus.dev/core/errors"
)

type User struct {
	entry `json:"-"`

	Username string
	Name     string
	Avatar   string   `json:",omitempty"`
	Groups   []*Group `json:",omitempty"`
}

func (u *User) GetSharedGroupsURL() string {
	return "/data-library/shared/" + url.PathEscape(u.Username)
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
			if len(u0.Name) == 0 {
				u0.Name = strings.TrimSpace(u.Name)
			}

			if len(u0.Avatar) == 0 {
				u0.Avatar = u.Avatar
			}

			// merge groups
			for _, g := range u.Groups {
				if _, err := u0.AddGroup(g, true); err != nil {
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
		u0.root = w
		w.User[user] = u0

		// add groups
		for _, g := range u.Groups {
			if _, err := u0.AddGroup(g, false); err != nil {
				check.AppendWrapped(err, "%s[%q]", "User", user)
			}
		}
	}

	if !check.Ok() {
		return nil, &check
	} else {
		return u0, nil
	}
}
