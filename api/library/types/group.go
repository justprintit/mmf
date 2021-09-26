package types

import (
	"fmt"
	"path"
	"strings"
	"time"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/util"
)

type Group struct {
	entry  `json:"-"`
	user   *User  `json:"-"`
	parent *Group `json:"-"`

	NextGroupObjectsUpdate time.Time `json:"-"`

	Name string
	Id   int

	Objects   []*Object `json:",omitempty"`
	Subgroups []*Group  `json:",omitempty"`
}

type Object struct{}

func (g *Group) GetObjectsURL() string {
	return fmt.Sprintf("/data-library/group/%v", g.Id)
}

func (g *Group) SanitizedName() string {
	return util.Sanitize(g.Name)
}

func (g *Group) Path() string {
	var s = []string{g.SanitizedName()}

	for {
		var name string

		if g == nil {
			break
		} else if p := g.parent; p != nil {
			g = p
			name = g.SanitizedName()
		} else {
			u := g.user
			g = nil
			name = u.SanitizedName()
		}

		t := make([]string, 1, len(s)+1)
		t[0] = name
		s = append(t, s...)
	}

	return path.Join(s...)
}

func (g *Group) User() *User {
	g.entry.Lock()
	defer g.entry.Unlock()

	return g.user
}

func (g *Group) GroupsAll() []*Group {
	g.entry.Lock()
	defer g.entry.Unlock()

	return g.groupsAll()
}

func (g *Group) groupsAll() []*Group {
	groups := make([]*Group, 0, len(g.Subgroups)+1)
	groups = append(groups, g)

	for _, sg := range g.Subgroups {
		all := sg.groupsAll()
		groups = append(groups, all...)
	}

	return groups
}

func (g *Group) updateName(s string) {
	if len(g.Name) == 0 {
		g.updateString("Name", &g.Name, s)
	}
}

func (g *Group) updateString(field string, v *string, s string) {
	before := *v
	after := strings.TrimSpace(s)
	if before != after {
		*v = after
		g.entry.OnGroupUpdate(g, field, before, after)
	}
}

func (g *Group) AddSubgroup(sg *Group, merge bool) (*Group, error) {
	g.entry.Lock()
	defer g.entry.Unlock()

	return addGroup(nil, g, sg, merge)
}

func (u *User) AddGroup(g *Group, merge bool) (*Group, error) {
	u.entry.Lock()
	defer u.entry.Unlock()

	return addGroup(u, nil, g, merge)
}

func (u *User) addGroup(g *Group, merge bool) (*Group, error) {
	return addGroup(u, nil, g, merge)
}

func newGroup(u *User, parent *Group, id int, name string) *Group {
	w := u.entry.Library

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
	w.OnNewGroup(g)
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

	w := u.entry.Library

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
			g0.updateName(g.Name)

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
		err := &check
		w.OnUserError(u, err)
		return nil, err
	}

	return g0, nil
}

func (w *Library) registerGroup(g *Group) {
	if w.group == nil {
		w.group = make(map[int]*Group, 1)
	}

	g.Library = w
	w.group[g.Id] = g
}

func (w *Library) GetGroup(id int) (*Group, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if g, ok := w.group[id]; ok {
		return g, nil
	} else {
		err := errors.New("%s[%v]: Not Found", "Group", id)
		return nil, err
	}
}
