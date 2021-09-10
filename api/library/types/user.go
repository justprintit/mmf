package types

import (
	"go.sancus.dev/core/errors"
)

type User struct {
	Username string
	Name     string
	Avatar   string `json:",omitempty"`
}

func (w *User) Merge(u User) error {
	if len(u.Username) > 0 {
		UpdateString("Username", &w.Username, u.Username)
	}
	if len(u.Name) > 0 {
		UpdateString("Name", &w.Name, u.Name)
	}
	if len(u.Avatar) > 0 {
		UpdateString("Avatar", &w.Avatar, u.Avatar)
	}
	return nil
}

func (w *Library) AddUser(u User) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.User == nil {
		w.User = make(map[string]User, 1)
	}

	if name := u.Username; len(name) == 0 {
		return errors.ErrMissingField("Name")
	} else if v, ok := w.User[u.Username]; ok {
		// exists
		return v.Merge(u)
	} else {
		// new
		w.User[u.Username] = u
		return nil
	}
}
