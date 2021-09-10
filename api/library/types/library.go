package types

import (
	"sync"
)

type Library struct {
	mu sync.Mutex

	User map[string]*User `json:",omitempty"`

	// index
	group map[int]*Group
}

type entry struct {
	root *Library
}

func (e *entry) Library() *Library {
	return e.root
}

func (e *entry) Lock() {
	e.root.mu.Lock()
}

func (e *entry) Unlock() {
	e.root.mu.Unlock()
}
