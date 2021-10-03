package types

import (
	"path/filepath"
	"strings"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/util"
)

type Tribe struct {
	entry
	user *User

	name string
	id   string

	Avatar string

	Groups []*Group
}

func (g *Tribe) Id() string {
	return g.id
}

func (g *Tribe) Name() string {
	return util.Sanitize(g.name)
}

func (g *Tribe) Type() NodeType {
	return TribeNode
}

func (g *Tribe) Path() string {
	var s = make([]string, 2)
	s[0] = g.user.Name()
	s[1] = g.Name()
	return filepath.Join(s...)
}

func (g *Tribe) User() *User {
	return g.user
}

func (g *Tribe) Parent() Node {
	return g.user
}

func (g *Tribe) GroupsAll() []*Group {
	g.entry.Lock()
	defer g.entry.Unlock()

	return g.groupsAll()
}

func (g *Tribe) groupsAll() []*Group {
	groups := make([]*Group, 0, len(g.Groups))

	for _, sg := range g.Groups {
		all := sg.groupsAll()
		groups = append(groups, all...)
	}

	return groups
}

func (g *Tribe) updateName(s string) {
	if len(g.name) == 0 {
		g.updateString("Name", &g.name, s)
	}
}

func (g *Tribe) updateAvatar(s string) {
	g.updateString("Avatar", &g.Avatar, s)
}

func (g *Tribe) updateString(field string, v *string, s string) {
	before := *v
	after := strings.TrimSpace(s)
	if before != after {
		*v = after
		g.entry.OnNodeUpdate(g, field, before, after)
	}
}

func (g *Tribe) HasNodes() bool {
	return len(g.Groups) > 0
}

func (g *Tribe) Nodes() []Node {
	n := len(g.Groups)
	nodes := make([]Node, 0, n)

	for _, g := range g.Groups {
		nodes = append(nodes, g)
	}

	return nodes
}

func (g *Tribe) AddNode(sg Node, merge bool) (Node, error) {
	if g0, ok := sg.(*Group); ok {
		return g.AddGroup(g0, merge)
	} else {
		err := errors.ErrInvalidValue("%s[%s]: %T (%s)", g.Type(), g.Id(), sg, sg.Type())
		return nil, err
	}
}

func (g *Tribe) appendGroup(sg *Group) {
	g.Groups = append(g.Groups, sg)
}

func NewTribe(id, name string) *Tribe {
	id = strings.TrimSpace(id)
	name = util.Sanitize(name)
	if len(name) == 0 {
		name = id
	}

	return &Tribe{
		id:   id,
		name: name,
	}
}

func newTribe(u *User, id string, name string) *Tribe {
	w := u.entry.Library()

	g := NewTribe(id, name)
	g.user = u
	// register
	g.entry.root = w
	w.addNode(TribeNode, g.id, g)
	// append
	u.Tribes = append(u.Tribes, g)
	// report
	w.OnNewNode(g)
	return g
}

func (w *Library) getTribe(id string) *Tribe {
	if g := w.getNode(TribeNode, id); g != nil {
		return g.(*Tribe)
	}
	return nil
}

func (u *User) AddTribe(g *Tribe, merge bool) (*Tribe, error) {
	var check errors.ErrorStack
	var g0 *Tribe

	u.entry.Lock()
	defer u.entry.Unlock()

	w := u.entry.Library()

	if g0 = w.getTribe(g.id); g0 != nil {
		// exists
		if !merge {
			check.AppendErrorf("%s[%s]: Already exists",
				"Tribe", g.id)
		} else if g0 == g {
			// same same
		} else if g0.user != u {
			check.AppendErrorf("%s[%v]: already assigned to user %q",
				"Tribe", g.id, g0.user.Id())
		} else {
			// merge
			g0.updateName(g.name)
			g0.updateAvatar(g.Avatar)

			// merge groups
			for _, g1 := range g.Groups {
				if _, err := addGroup(u, g0, g1, merge); err != nil {
					check.AppendError(err)
				}
			}
		}
	} else {
		// new
		g0 = newTribe(u, g.id, g.name)
		g0.Avatar = g.Avatar

		for _, g1 := range g.Groups {
			// new tribe, new groups
			if _, err := addGroup(u, g0, g1, false); err != nil {
				check.AppendError(err)
			}
		}
	}

	return g0, check.AsError()
}
