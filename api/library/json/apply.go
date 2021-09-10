package json

import (
	"log"
	"strings"

	"github.com/justprintit/mmf/api/library/types"
)

func (w *User) Export() *types.User {
	name := strings.TrimSpace(w.Name)
	if len(name) == 0 {
		name = w.Username
	}

	return &types.User{
		Username: w.Username,
		Name:     name,
		Avatar:   w.Avatar,
	}
}

func (w *Users) Apply(d *types.Library) error {
	n := len(w.User)
	if n != w.Count {
		log.Printf("Users: expected:%v != actual:%v", w.Count, n)
	}

	for i, v := range w.User {
		u := v.Export()

		log.Printf("User[%v/%v]: %s (%s)", i, n, u.Name, u.Username)

		if err := d.AddUser(u); err != nil {
			log.Printf("User[%v/%v]: Failed to add User: %s",
				i, n, err)
		}
	}

	return nil
}

func (w *Objects) Apply(d *types.Library) error {
	return nil
}
