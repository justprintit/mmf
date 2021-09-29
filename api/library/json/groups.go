package json

import (
	"log"
	"sort"
	"strings"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/library/types"
)

type Groups struct {
	Count int     `json:"total_count"`
	Items []Group `json:"items"`
}

type Group struct {
	Id           GroupId        `json:"id"`
	Name         string         `json:"name"`
	API          map[string]API `json:"apis,omitempty"`
	TotalObjects int            `json:"total_count_objects,omitempty"`
	Children     []Group        `json:"childrens,omitempty"`
	Objects
}

type GroupId struct {
	types.Id
}

func (w *Group) Export(recursive bool) *types.Group {
	if w.Id.Ok() {
		g := &types.Group{
			Id:   w.Id.Id,
			Name: strings.TrimSpace(w.Name),
		}

		if recursive {
			for _, v := range w.Children {
				cg := v.Export(recursive)
				g.Subgroups = append(g.Subgroups, cg)
			}
		}

		return g
	}
	return nil
}

func (w *Group) Apply(d *types.Library, u *types.User, parent *types.Group) (*types.Group, error) {
	const merge = true

	if g := w.Export(false); g != nil {
		var err error

		if parent == nil {
			g, err = u.AddGroup(g, merge)
		} else {
			g, err = parent.AddSubgroup(g, merge)
		}

		if err != nil {
			return nil, err
		}

		// subgroups
		if n := len(w.Children); n > 0 {
			// sorted
			subgroups := make([]*Group, 0, n)
			for i := range w.Children {
				p := &w.Children[i]
				if !strings.HasPrefix(p.Id.String(), "all/") {
					subgroups = append(subgroups, p)
				}
			}

			sort.Slice(subgroups[:], func(i, j int) bool {
				a, _ := subgroups[i].Id.Int()
				b, _ := subgroups[j].Id.Int()
				return a < b
			})

			for _, p := range subgroups {
				if _, err := p.Apply(d, u, g); err != nil {
					return g, err
				}
			}
		}

		// objects
		if n := len(w.Items); n > 0 {
			var check errors.ErrorStack

			if v, err := w.Count.Int64(); err == nil {
				if int64(n) != v {
					log.Printf("Group[%v].Items: expected:%v != actual:%v", g.Id, v, n)
				}
			}

			for _, obj := range w.Items {
				if err := obj.Apply(d, u, g); err != nil {
					check.AppendError(err)
				}
			}

			if !check.Ok() {
				err := &check
				d.OnGroupError(g, err)
				return g, err
			}
		}
		return g, nil
	}

	return nil, nil
}

func (w *Groups) Apply(d *types.Library, u *types.User) error {
	var check errors.ErrorStack

	// expected quantity
	if n := len(w.Items); n > 0 {
		if n != w.Count {
			log.Printf("Groups: expected:%v != actual:%v", w.Count, n)
		}
	}

	// apply them all to the library
	for i, v := range w.Items {
		if _, err := v.Apply(d, u, nil); err != nil {
			check.AppendWrapped(err, "Group.%v: %q (%v)", i, v.Name, v.Id)
		}
	}

	if !check.Ok() {
		return &check
	}
	return nil
}
