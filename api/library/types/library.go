package types

import (
	"sync"
)

type LibraryEvents struct {
	OnNewUser  func(w *Library, u *User)
	OnNewGroup func(w *Library, g *Group)

	OnUserUpdate  func(w *Library, u *User, field string, before, after interface{})
	OnGroupUpdate func(w *Library, g *Group, field string, before, after interface{})

	OnError      func(w *Library, err error)
	OnUserError  func(w *Library, u *User, err error)
	OnGroupError func(w *Library, g *Group, err error)
}

type Library struct {
	mu     sync.Mutex
	events LibraryEvents

	User map[string]*User `json:",omitempty"`

	// index
	group map[int]*Group
}

func (w *Library) SetEvents(ev LibraryEvents) {
	w.events = ev
}

func (w *Library) OnNewUser(u *User) {
	if f := w.events.OnNewUser; f != nil {
		f(w, u)
	}
}

func (w *Library) OnNewGroup(g *Group) {
	if f := w.events.OnNewGroup; f != nil {
		f(w, g)
	}
}

func (w *Library) OnUserUpdate(u *User, field string, before, after interface{}) {
	if f := w.events.OnUserUpdate; f != nil {
		f(w, u, field, before, after)
	}
}

func (w *Library) OnGroupUpdate(g *Group, field string, before, after interface{}) {
	if f := w.events.OnGroupUpdate; f != nil {
		f(w, g, field, before, after)
	}
}

func (w *Library) OnError(err error) {
	if f := w.events.OnError; f != nil {
		f(w, err)
	}
}

func (w *Library) OnUserError(u *User, err error) {
	if f := w.events.OnUserError; f != nil {
		f(w, u, err)
	}
}

func (w *Library) OnGroupError(g *Group, err error) {
	if f := w.events.OnGroupError; f != nil {
		f(w, g, err)
	}
}

type entry struct {
	*Library
}

func (e *entry) Lock() {
	e.mu.Lock()
}

func (e *entry) Unlock() {
	e.mu.Unlock()
}
