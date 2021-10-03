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

	name     string
	username string

	Avatar string   `json:",omitempty"`
	Groups []*Group `json:",omitempty"`
	Tribes []*Tribe `json:",omitempty"`
}

func (u *User) GetSharedGroupsURL() string {
	return "/data-library/shared/" + url.PathEscape(u.username)
}

func (u *User) Id() string {
	return u.username
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Type() NodeType {
	return UserNode
}

func (u *User) Path() string {
	return u.Name()
}

func (u *User) User() *User {
	return u
}

func (u *User) Parent() Node {
	return nil
}

func (u *User) GroupsAll() []*Group {
	u.entry.Lock()
	defer u.entry.Unlock()

	return u.groupsAll()
}

func (u *User) groupsAll() []*Group {
	n := len(u.Groups)
	for _, tribe := range u.Tribes {
		n += len(tribe.Groups)
	}

	groups := make([]*Group, 0, n)

	for _, g := range u.Groups {
		all := g.groupsAll()
		groups = append(groups, all...)
	}

	for _, tribe := range u.Tribes {
		all := tribe.groupsAll()
		groups = append(groups, all...)
	}

	return groups
}

func (u *User) SetName(name string) {
	name = util.Sanitize(name)
	if name == "" {
		name = u.username
	}
	u.updateString("Name", &u.name, name)
}

func (u *User) updateName(s string) {
	if len(u.name) == 0 {
		u.updateString("Name", &u.name, s)
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
		u.entry.OnNodeUpdate(u, field, before, after)
	}
}

func (u *User) appendGroup(g *Group) {
	u.Groups = append(u.Groups, g)
}

func NewUser(username, name string) *User {
	name = util.Sanitize(name)
	if len(name) == 0 {
		name = username
	}

	return &User{
		username: username,
		name:     name,
	}
}

func (w *Library) GetUser(user string) (*User, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if u := w.getUser(user); u != nil {
		return u, nil
	} else {
		err := errors.New("%s[%q]: Not Found", "User", user)
		return nil, err
	}
}

func (w *Library) AddUser(u *User, merge bool) (*User, error) {
	var check errors.ErrorStack
	var u0 *User

	w.mu.Lock()
	defer w.mu.Unlock()

	if user := u.username; len(user) == 0 {
		check.MissingField("Username")
	} else if u0 = w.getUser(user); u0 != nil {
		// exists
		if !merge {
			err := errors.New("%s[%q]: Already exists", "User", user)

			check.AppendError(err)
		} else {
			// merge
			u0.updateName(u.name)
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
		u0 = NewUser(u.username, u.name)
		u0.Avatar = u.Avatar

		// register
		w.registerUser(u0)

		// add groups
		for _, g := range u.Groups {
			if _, err := u0.addGroup(g, false); err != nil {
				check.AppendWrapped(err, "%s[%q]", "User", user)
			}
		}
	}

	if err := check.AsError(); err != nil {
		w.OnError(u0, err)
		return nil, err
	}

	return u0, nil
}

func (w *Library) getUser(username string) *User {
	if u := w.getNode(UserNode, username); u != nil {
		return u.(*User)
	}
	return nil
}

func (w *Library) registerUser(u *User) {
	u.entry.root = w
	w.addNode(UserNode, u.username, u)
	w.OnNewNode(u)
}

// Users() returns slice of registered usernames
func (w *Library) Users() []string {
	return w.Keys(UserNode)
}

// UsernameFromURL() attempts to extract the username from a URL
func UsernameFromURL(s string) (string, error) {
	if u, err := url.Parse(s); err == nil {
		if p, err := UsernameFromPath(u.Path); err == nil {
			return p, nil
		}
	}
	return "", ErrInvalidPath(s)
}

// UsernameFromPath() attempts to extract the username from a URL.Path
func UsernameFromPath(s string) (string, error) {
	if _, p, _ := util.NextInPathUnescaped(s,
		"/users/",
		"/data-library/shared/",
	); p != "" {
		return p, nil
	}
	return "", ErrInvalidPath(s)
}
