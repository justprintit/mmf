package json

import (
	"log"

	"github.com/justprintit/mmf/api/library/types"
)

func (w *Users) Apply(d *types.Library) error {
	n := len(w.User)
	if n != w.Count {
		log.Printf("Users: expected:%v != actual:%v", w.Count, n)
	}

	for i, v := range w.User {
		log.Printf("User[%v/%v]: %s (%s)", i, n, v.Username, v.Name)

		if v.Name == "" {
			v.Name = v.Username
		}

		u := types.User{
			Username: v.Username,
			Name:     v.Name,
			Avatar:   v.Avatar,
		}

		if err := d.AddUser(u); err != nil {
			log.Printf("User[%v/%v]: Failed to add User: %s",
				i, n, v.Avatar, err)
		}
	}

	return nil
}

func (w *Objects) Apply(d *types.Library) error {
	return nil
}
