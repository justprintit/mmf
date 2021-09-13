package json

import (
	"log"
	"sort"
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
	if n := len(w.Items); n > 0 {
		if v, err := w.Count.Int64(); err == nil {
			if int64(n) != v {
				log.Printf("Users: expected:%v != actual:%v", v, n)
			}
		}
	}

	for i, v := range w.Items {
		u := v.Export()

		log.Printf("User.%v: %s (%s)", i, u.Name, u.Username)

		if err := d.AddUser(u); err != nil {
			log.Printf("User.%v: Failed to add User: %s", i, err)
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

func (w *Group) Export() *types.Group {
	if id, ok := w.Id.Int(); ok {
		return &types.Group{
			Id:   id,
			Name: strings.TrimSpace(w.Name),
		}
	}
	return nil
}

func (w *Group) Apply(d *types.Library, u *types.User, parent *types.Group) error {
	if g := w.Export(); g != nil {
		if parent == nil {
			if err := u.AddGroup(g); err != nil {
				return err
			}
		} else if err := parent.AddSubgroup(g); err != nil {
			return err
		}

		if n := len(w.Children); n > 0 {
			subgroups := make([]*Group, 0, n)
			for i := range w.Children {
				p := &w.Children[i]
				if _, ok := p.Id.Int(); ok {
					subgroups = append(subgroups, p)
				}
			}

			sort.Slice(subgroups[:], func(i, j int) bool {
				a, _ := subgroups[i].Id.Int()
				b, _ := subgroups[j].Id.Int()
				return a < b
			})

			for _, p := range subgroups {
				if err := p.Apply(d, u, g); err != nil {
					return err
				}
			}
		}

		if n := len(w.Items); n > 0 {
			if v, err := w.Count.Int64(); err == nil {
				if int64(n) != v {
					log.Printf("Group[%v].Items: expected:%v != actual:%v", g.Id, v, n)
				}
			}

			return ApplyObjects(d, u, g, w.Items...)
		}
	}
	return nil
}

func (w *Objects) Apply(d *types.Library) error {
	if n := len(w.Items); n > 0 {
		if v, err := w.Count.Int64(); err == nil {
			if int64(n) != v {
				log.Printf("Objects: expected:%v != actual:%v", v, n)
			}
		}

		return ApplyObjects(d, nil, nil, w.Items...)
	}
	return nil
}

func ApplyObjects(d *types.Library, u *types.User, g *types.Group, objects ...Object) error {
	return nil
}
