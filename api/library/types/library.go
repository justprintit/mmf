package types

import (
	"sync"
)

type LibraryEvents struct {
	OnNewNode    func(w *Library, g Node)
	OnNodeUpdate func(w *Library, g Node, field string, before, after interface{})
	OnError      func(w *Library, g Node, err error)
}

type Library struct {
	mu     sync.Mutex
	events LibraryEvents

	index [NodeTypes]map[string]Node
}

func (w *Library) SetEvents(ev LibraryEvents) {
	w.events = ev
}

func (w *Library) getNode(typ NodeType, key string) Node {
	if g, ok := w.index[typ][key]; ok {
		return g
	} else {
		return nil
	}
}

func (w *Library) addNode(typ NodeType, key string, node Node) {
	if w.index[typ] == nil {
		w.index[typ] = make(map[string]Node, 1)
	}
	w.index[typ][key] = node
}

// Keys() returns slice of keys of a given type
func (w *Library) Keys(typ NodeType) []string {
	w.mu.Lock()
	defer w.mu.Unlock()

	s := make([]string, 0, len(w.index[typ]))
	for k := range w.index[typ] {
		s = append(s, k)
	}
	return s
}

func (w *Library) OnNewNode(g Node) {
	if f := w.events.OnNewNode; f != nil {
		f(w, g)
	}
}

func (w *Library) OnNodeUpdate(g Node, field string, before, after interface{}) {
	if f := w.events.OnNodeUpdate; f != nil {
		f(w, g, field, before, after)
	}
}

func (w *Library) OnError(g Node, err error) {
	if f := w.events.OnError; f != nil {
		f(w, g, err)
	}
}

type entry struct {
	root *Library
}

func (e *entry) Lock() bool {
	if e.root == nil {
		// Dummy
		return false
	}
	e.root.mu.Lock()
	return true
}

func (e *entry) Unlock() {
	e.root.mu.Unlock()
}

func (e *entry) Library() *Library {
	return e.root
}

func (e *entry) OnNewNode(g Node) {
	e.root.OnNewNode(g)
}

func (e *entry) OnNodeUpdate(g Node, field string, before, after interface{}) {
	e.root.OnNodeUpdate(g, field, before, after)
}

func (e *entry) OnError(g Node, err error) {
	e.root.OnError(g, err)
}
