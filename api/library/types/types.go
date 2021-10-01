package types

import (
	"fmt"
)

type NodeType int

// Each NodeType is an independent ID namespace
const (
	GroupNode NodeType = iota
	ObjectNode
	UserNode

	NodeTypes int = 3
)

func (t NodeType) String() string {
	switch t {
	case UserNode:
		return "User"
	case GroupNode:
		return "Group"
	case ObjectNode:
		return "Object"
	default:
		return fmt.Sprintf("Unknown[%v]", int(t))
	}
}

type Node interface {
	Id() string
	Name() string
	Type() NodeType
	Path() string

	Library() *Library
	User() *User
	Parent() Node
}

type Grouper interface {
	Node

	AddGroup(*Group, bool) (*Group, error)
}

type groupAppender interface {
	appendGroup(*Group)
}
