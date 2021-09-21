package types

import (
	"fmt"
	"strings"

	"go.sancus.dev/core/errors"
)

type Group struct {
	entry  `json:"-"`
	user   *User  `json:"-"`
	parent *Group `json:"-"`

	Name string
	Id   int

	Objects   []*Object `json:",omitempty"`
	Subgroups []*Group  `json:",omitempty"`
}

type Object struct{}

func (g *Group) GetObjectsURL() string {
	return fmt.Sprintf("/data-library/group/%v", g.Id)
}

func (g *Group) AddSubgroup(sg *Group, merge bool) (*Group, error) {
	w := g.entry.Library()

	w.mu.Lock()
	defer w.mu.Unlock()

	return addGroup(nil, g, sg, merge)
}

func (u *User) AddGroup(g *Group, merge bool) (*Group, error) {
	w := u.entry.Library()

	w.mu.Lock()
	defer w.mu.Unlock()

	return addGroup(u, nil, g, merge)
}

func (u *User) addGroup(g *Group, merge bool) (*Group, error) {
	return addGroup(u, nil, g, merge)
}

func newGroup(u *User, parent *Group, id int, name string) *Group {
	w := u.entry.Library()

	// sanitize name
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		name = fmt.Sprintf("%v", id)
	}

	g := &Group{
		user:   u,
		parent: parent,

		Id:   id,
		Name: name,
	}

	w.registerGroup(g)
	if parent == nil {
		u.Groups = append(u.Groups, g)
	} else {
		parent.Subgroups = append(parent.Subgroups, g)
	}
	return g
}

func addGroup(u *User, parent *Group, g *Group, merge bool) (*Group, error) {
	var check errors.ErrorStack
	var g0 *Group
	var ok bool

	// validate user
	if u == nil {
		if parent == nil || parent.user == nil {
			err := errors.ErrMissingArgument("User not provided")
			return nil, err
		} else {
			u = parent.user
		}
	}

	w := u.entry.Library()

	if u0, ok := w.User[u.Username]; ok {
		// just in case it's a dummy
		u = u0
	} else {
		err := errors.ErrInvalidArgument("User")
		return nil, err
	}

	if g.Id <= 0 {
		check.InvalidArgument("%s.%s", "Group", "Id")
	} else if g0, ok = w.group[g.Id]; ok {
		var err error

		// exists
		if !merge {
			err = errors.New("%s[%v]: Already exists", "Group", g.Id)
		} else if g0.user != u {
			err = errors.New("%s[%v]: already assigned to user %q",
				"Group", g.Id, g0.user.Username)
		} else {
			if g0.parent == parent {
				// same
			} else if g0.parent != nil && parent != nil && g0.parent.Id == parent.Id {
				// dummy
				parent = g0.parent
			} else if g0.parent == nil {
				err = errors.New("%s[%v]: already attached directly",
					"Group", g.Id)
			} else {
				err = errors.New("%s[%v]: already attached to group %v",
					"Group", g.Id, g0.parent.Id)
			}
		}

		if err != nil {
			check.AppendError(err)
		} else {
			// merge
			if len(g0.Name) == 0 {
				g0.Name = strings.TrimSpace(g.Name)
			}

			// merge subgroups
			for _, sg := range g.Subgroups {
				if _, err := addGroup(u, g0, sg, true); err != nil {
					check.AppendError(err)
				}
			}
		}
	} else {
		// new
		g0 = newGroup(u, parent, g.Id, g.Name)

		for _, sg := range g.Subgroups {
			// new group, new subgroups.
			if _, err := addGroup(u, g0, sg, false); err != nil {
				check.AppendError(err)
			}
		}

	}

	if !check.Ok() {
		return nil, &check
	} else {
		return g0, nil
	}
}

func (w *Library) checkNewGroups(groups ...*Group) error {
	for _, g := range groups {
		if _, ok := w.group[g.Id]; ok {
			return errors.New("%s[%v]: Already exists", "Group", g.Id)
		}
	}
	return nil
}

func (w *Library) registerGroup(g *Group) {
	if w.group == nil {
		w.group = make(map[int]*Group, 1)
	}

	g.root = w
	w.group[g.Id] = g
}
