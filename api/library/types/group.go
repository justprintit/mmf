package types

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/util"
)

type Group struct {
	entry  `json:"-"`
	user   *User   `json:"-"`
	parent Grouper `json:"-"`

	NextGroupObjectsUpdate time.Time `json:"-"`

	name string
	id   string

	Objects []*Object `json:",omitempty"`
	Groups  []*Group  `json:",omitempty"`
}

type Object struct{}

func (g *Group) GetObjectsURL() string {
	return fmt.Sprintf("/data-library/group/%s", g.id)
}

func (g *Group) Id() string {
	return g.id
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) Type() NodeType {
	return GroupNode
}

func (g *Group) Path() string {
	return filepath.Join(g.Parent().Path(), g.Name())
}

func (g *Group) Parent() Node {
	if g.parent != nil {
		return g.parent
	} else {
		return g.user
	}
}

func (g *Group) User() *User {
	return g.user
}

func (g *Group) GroupsAll() []*Group {
	g.entry.Lock()
	defer g.entry.Unlock()

	return g.groupsAll()
}

func (g *Group) groupsAll() []*Group {
	groups := make([]*Group, 0, len(g.Groups)+1)
	groups = append(groups, g)

	for _, sg := range g.Groups {
		all := sg.groupsAll()
		groups = append(groups, all...)
	}

	return groups
}

func (g *Group) appendGroup(sg *Group) {
	g.Groups = append(g.Groups, sg)
}

func (g *Group) SetName(name string) {
	name = util.Sanitize(name)
	if name == "" {
		name = g.id
	}
	g.updateString("Name", &g.name, name)
}

func (g *Group) updateName(s string) {
	if len(g.name) == 0 {
		g.updateString("Name", &g.name, s)
	}
}

func (g *Group) updateString(field string, v *string, s string) {
	before := *v
	after := strings.TrimSpace(s)
	if before != after {
		*v = after
		g.entry.OnNodeUpdate(g, field, before, after)
	}
}

func (g *Group) AddGroup(sg *Group, merge bool) (*Group, error) {
	if !g.entry.Lock() {
		// Dummy
		g.appendGroup(sg)
		return sg, nil
	}

	defer g.entry.Unlock()
	return addGroup(nil, g, sg, merge)
}

func (u *User) AddGroup(g *Group, merge bool) (*Group, error) {
	if !u.entry.Lock() {
		// Dummy
		u.appendGroup(g)
		return g, nil
	}

	defer u.entry.Unlock()
	return addGroup(u, nil, g, merge)
}

func (u *User) addGroup(g *Group, merge bool) (*Group, error) {
	return addGroup(u, nil, g, merge)
}

// Creates dummy Group, disconnected from the library
func NewGroup(id, name string) *Group {
	return &Group{
		id:   strings.TrimSpace(id),
		name: util.Sanitize(name),
	}
}

func newGroup(u *User, parent Grouper, id string, name string) *Group {
	if len(id) == 0 {
		panic(ErrInvalidValue(id))
	}

	w := u.Library()

	// sanitize name
	name = util.Sanitize(name)
	if len(name) == 0 {
		name = id
	}

	g := &Group{
		user:   u,
		parent: parent,

		id:   id,
		name: name,
	}

	w.registerGroup(g)

	w.OnNewNode(g)
	return g
}

func addGroup(u *User, parent Grouper, g *Group, merge bool) (*Group, error) {
	var check errors.ErrorStack
	var g0 *Group

	// validate user
	if u == nil {
		if parent == nil || parent.User() == nil {
			err := errors.ErrMissingArgument("User not provided")
			return nil, err
		} else {
			u = parent.User()
		}
	}

	w := u.entry.Library()

	if u0 := w.getUser(u.Id()); u0 != nil {
		// just in case it's a dummy
		u = u0
	} else {
		err := errors.ErrInvalidArgument("User")
		return nil, err
	}

	if len(g.id) == 0 {
		check.InvalidArgument("%s.%s", "Group", "Id")
	} else if g0 = w.getGroup(g.id); g0 != nil {
		var err error

		// exists
		if !merge {
			err = errors.New("%s[%s]: Already exists", "Group", g.id)
		} else if g0.user != u {
			err = errors.New("%s[%s]: already assigned to user %q",
				"Group", g.id, g0.user.Id())
		} else {
			if g0.parent == parent {
				// same
			} else if g0.parent != nil && parent != nil && g0.parent.Id() == parent.Id() && g0.parent.Type() == parent.Type() {
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
			g0.updateName(g.name)

			// merge subgroups
			for _, sg := range g.Groups {
				if _, err := addGroup(u, g0, sg, true); err != nil {
					check.AppendError(err)
				}
			}
		}
	} else {
		// new
		g0 = newGroup(u, parent, g.Id(), g.Name())

		for _, sg := range g.Groups {
			// new group, new subgroups.
			if _, err := addGroup(u, g0, sg, false); err != nil {
				check.AppendError(err)
			}
		}

	}

	if err := check.AsError(); err != nil {
		w.OnError(u, err)
		return nil, err
	}

	return g0, nil
}

func (w *Library) getGroup(id string) *Group {
	if g := w.getNode(GroupNode, id); g != nil {
		return g.(*Group)
	}
	return nil
}

func (w *Library) registerGroup(g *Group) {
	var pg groupAppender

	g.entry.root = w
	w.addNode(GroupNode, g.id, g)

	if g.parent == nil {
		pg = g.user
	} else {
		pg = g.parent.(groupAppender)
	}

	pg.appendGroup(g)
}

func (w *Library) GetGroup(id string) (*Group, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if g := w.getGroup(id); g != nil {
		return g, nil
	} else {
		err := errors.New("%s[%v]: Not Found", GroupNode, id)
		return nil, err
	}
}

// Groups() returns slice of registered IDs
func (w *Library) Groups() []string {
	return w.Keys(GroupNode)
}

// GroupIdFromURL() attempts to extract the group Id from a URL
func GroupIdFromURL(s string) (string, error) {
	if u, err := url.Parse(s); err == nil {
		if p, err := GroupIdFromPath(u.Path); err == nil {
			return p, nil
		}
	}
	return "", ErrInvalidPath(s)
}

// GroupIdFromPath() attempts to extract the group Id from a URL.Path
func GroupIdFromPath(s string) (string, error) {
	if _, p, _ := util.NextInPathUnescaped(s,
		"/data-library/group/",
	); p != "" {
		return p, nil
	}
	return "", ErrInvalidPath(s)
}
