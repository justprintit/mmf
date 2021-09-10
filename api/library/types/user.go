package types

import (
	"go.sancus.dev/core/errors"
)

type User struct {
	Username string
	Name     string
	Avatar   string `json:",omitempty"`
}

func (w *Library) AddUser(u *User) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.User == nil {
		w.User = make(map[string]*User, 1)
	}

	if user := u.Username; len(user) == 0 {
		return errors.ErrMissingField("Username")
	} else if _, ok := w.User[user]; ok {
		// exists
		return errors.New("%s[%q]: Already exists", "User", user)
	} else {
		// new
		w.User[user] = u
		return nil
	}
}
