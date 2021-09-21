package json

import (
	"log"
	"sort"
	"strings"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/client/json"
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

func (w *User) Export(groups bool) *types.User {
	name := strings.TrimSpace(w.Name)
	if len(name) == 0 {
		name = w.Username
	}

	u := &types.User{
		Username: w.Username,
		Name:     name,
		Avatar:   w.Avatar,
	}

	if groups {
		const recursive = true

		u.Groups = w.ExportGroups(recursive)
	}

	return u
}

func (w *User) ExportGroups(recursive bool) []*types.Group {
	n := len(w.Groups.Group)

	if k := w.Groups.Count; k != n {
		log.Printf("User.Groups: expected:%v != actual:%v", k, n)
	}

	// export
	out := make([]*types.Group, 0, n)
	for i := range w.Groups.Group {
		p := &w.Groups.Group[i]
		if _, ok := p.Id.Int(); ok {
			// skip groups with string id
			if g := p.Export(recursive); g != nil {
				out = append(out, g)
			}
		}
	}

	// sort
	sort.Slice(out[:], func(i, j int) bool {
		a := out[i].Id
		b := out[i].Id
		return a < b
	})

	return out
}

func (w *User) Apply(d *types.Library) error {
	const merge = true
	const groups = true

	if u := w.Export(groups); u != nil {
		_, err := d.AddUser(u, merge)
		if err != nil {
			return errors.Wrap(err, "AddUser")
		}
	}

	return nil
}

func (w *Users) Apply(d *types.Library) error {
	var check errors.ErrorStack

	// expected quantity
	if n := len(w.Items); n > 0 {
		if v, err := w.Count.Int64(); err == nil {
			if int64(n) != v {
				log.Printf("Users: expected:%v != actual:%v", v, n)
			}
		}
	}

	// apply them all to the library
	for i, v := range w.Items {
		if err := v.Apply(d); err != nil {
			check.AppendWrapped(err, "User.%v: %q", i, v.Username)
		}
	}

	if !check.Ok() {
		return &check
	}
	return nil
}
