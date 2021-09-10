package types

import (
	"go.sancus.dev/core/errors"
)

type Group struct {
	entry `json:"-"`

	Name string
	URL  string
	Id   int

	Objects   []*Object `yaml:",omitempty"`
	Subgroups []*Group  `yaml:",omitempty"`
}

type Object struct{}

func (g *Group) AddSubgroup(sg *Group) error {
	w := g.entry.Library()

	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.checkNewGroups(sg); err != nil {
		return err
	}

	w.registerGroup(sg)
	g.Subgroups = append(g.Subgroups, sg)
	return nil
}

func (u *User) AddGroup(sg *Group) error {
	w := u.entry.Library()

	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.checkNewGroups(sg); err != nil {
		return err
	}

	w.registerGroup(sg)
	u.Groups = append(u.Groups, sg)
	return nil
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
