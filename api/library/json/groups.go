package json

import (
	"log"
	"sort"
	"strings"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/client/json"
	"github.com/justprintit/mmf/api/library/types"
)

type Groups struct {
	Count int     `json:"total_count"`
	Group []Group `json:"items"`
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
	id int
	s  string
}

func (w *GroupId) Int() (int, bool) {
	if len(w.s) > 0 {
		return 0, false
	} else {
		return w.id, true
	}
}

func (w *GroupId) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		w.id = 0
		return json.Unmarshal(data, &w.s)
	} else {
		w.s = ""
		return json.Unmarshal(data, &w.id)
	}
}

func (w *GroupId) MarshalJSON() ([]byte, error) {
	if len(w.s) > 0 {
		return json.Marshal(w.s)
	} else {
		return json.Marshal(w.id)
	}
}

func (w *Group) Export(recursive bool) *types.Group {
	if id, ok := w.Id.Int(); ok {
		g := &types.Group{
			Id:   id,
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

func (w *Group) Apply(d *types.Library, u *types.User, parent *types.Group) error {
	const merge = true

	if g := w.Export(false); g != nil {
		var err error

		if parent == nil {
			g, err = u.AddGroup(g, merge)
		} else {
			g, err = parent.AddSubgroup(g, merge)
		}

		if err != nil {
			return err
		}

		// subgroups
		if n := len(w.Children); n > 0 {
			// sorted
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
				d.OnUserError(u, err)
				return err
			}
		}
	}

	return nil
}
