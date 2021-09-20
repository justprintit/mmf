package json

import (
	"log"
	"sort"
	"strings"

	json "github.com/json-iterator/go"

	"github.com/justprintit/mmf/api/library/types"
)

type Users struct {
	Count json.Number `json:"total_count"`
	Items []User      `json:"items"`
}

type User struct {
	Id       string         `json:"id"`
	Username string         `json:"username"`
	Name     string         `json:"name"`
	Avatar   string         `json:"avatar_url"`
	API      map[string]API `json:"apis,omitempty"`
	Groups   Groups         `json:"groups,omitempty"`
}

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
	const merge = true

	if n := len(w.Items); n > 0 {
		if v, err := w.Count.Int64(); err == nil {
			if int64(n) != v {
				log.Printf("Users: expected:%v != actual:%v", v, n)
			}
		}
	}

	for i, v := range w.Items {
		var err error
		u := v.Export()

		log.Printf("User.%v: %s (%s)", i, u.Name, u.Username)

		u, err = d.AddUser(u, merge)
		if err != nil {
			log.Printf("User.%v: Failed to add User: %s", i, err)
			continue
		}

		// Groups
		if n := len(v.Groups.Group); n != v.Groups.Count {
			log.Printf("User.%v: Groups: expected:%v != actual:%v", i, v.Groups.Count, n)
		}

		groups := make([]*Group, 0, len(v.Groups.Group))
		for i := range v.Groups.Group {
			p := &v.Groups.Group[i]
			if _, ok := p.Id.Int(); ok {
				groups = append(groups, p)
			}
		}
		sort.Slice(groups[:], func(i, j int) bool {
			a, _ := groups[i].Id.Int()
			b, _ := groups[j].Id.Int()
			return a < b
		})

		for _, p := range groups {
			if err := p.Apply(d, u, nil); err != nil {
				id, _ := p.Id.Int()
				log.Printf("User.%v: Failed to add Group %q (%v): %s", i, p.Name, id, err)
			}
		}
	}

	return nil
}
