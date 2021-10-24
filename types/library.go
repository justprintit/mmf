package types

import (
	"sync"
)

type NodeType int

const (
	UndeterminedNode NodeType = iota

	UserNode
	ObjectNode
	CollectionNode
	GroupNode
	TribeNode
	CampaignNode

	NodeTypes int = iota
)

type Node interface {
	Library() *Library
	User() *User
	Type() NodeType

	Id() string
	Name() string
	Parent() Node
}

type entry struct {
	r *Library
	u *User
	t NodeType

	id   string
	name string
}

func (m *entry) Library() *Library {
	return m.r
}

func (m *entry) User() *User {
	return m.u
}

func (m *entry) Type() NodeType {
	return m.t
}

func (m *entry) Id() string {
	return m.id
}

func (m *entry) Name() string {
	return m.name
}

type Library struct {
	sync.Mutex

	Events LibraryEvents
	nodes  [NodeTypes]map[string]Node
}

func (m *Library) GetNode(t NodeType, id string) Node {
	m.Lock()
	defer m.Unlock()

	return m.getNode(t, id)
}

func (m *Library) getNode(t NodeType, id string) Node {
	if m.nodes[t] != nil {
		if v, ok := m.nodes[t][id]; ok {
			return v
		}
	}
	return nil
}

func (m *Library) registerNode(v Node) {
	t := v.Type()
	id := v.Id()

	if m.nodes[t] == nil {
		m.nodes[t] = make(map[string]Node, 1)
	}

	m.nodes[t][id] = v
}
