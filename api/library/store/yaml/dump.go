package yaml

import (
	"io"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/library/types"
)

type ExportDepth int

const (
	ExportShallow ExportDepth = iota
	ExportDeep
)

// Export User data for YAML encoding
func (store *Store) ExportUser(w *types.User, depth ExportDepth) (*User, error) {
	name := w.Name()
	if len(name) == 0 {
		name = w.Id()
	}

	u := &User{
		Username: w.Id(),
		Name:     name,
		Avatar:   w.Avatar,
	}

	if depth == ExportDeep {
		var err error

		u.Groups, err = store.ExportGroups(w.Groups, depth)
		if err != nil {
			return u, err
		}
	}

	return u, nil
}

// Export Groups data for YAML encoding
func (store *Store) ExportGroups(w []*types.Group, depth ExportDepth) ([]Group, error) {
	var check errors.ErrorStack

	groups := make([]Group, 0, len(w))

	// groups
	for _, g0 := range w {
		g, err := store.ExportGroup(g0, depth)
		if err != nil {
			check.AppendError(err)
		} else {
			groups = append(groups, g)
		}
	}

	// sorted
	sort.Slice(groups[:], func(i, j int) bool {
		return groups[i].Id.Lt(groups[j].Id)
	})

	return groups, check.AsError()
}

// Export Group data for YAML encoding
func (store *Store) ExportGroup(w *types.Group, depth ExportDepth) (g Group, err error) {
	var id types.Id

	name := w.Name()
	if len(name) == 0 {
		name = w.Id()
	}

	id, err = types.NewId(w.Id())
	if err != nil {
		return
	}

	g = Group{
		Id:   id,
		Name: name,
	}

	if depth == ExportDeep && len(w.Groups) > 0 {
		g.Subgroups, err = store.ExportGroups(w.Groups, depth)
	}

	return
}

// Export Library for YAML encoding
func (store *Store) ExportSlice(data *types.Library, depth ExportDepth) yaml.MapSlice {
	keys := data.Users()                // list of user
	n := len(keys)                      // users count
	mkeys := make(map[string]*User, n)  // lower case name to user
	slice := make([]yaml.MapItem, 0, n) // return value

	// convert unordered list of usernames
	// into a sorted list of lower case names
	for i, username := range keys {
		w, _ := data.GetUser(username)
		u, _ := store.ExportUser(w, depth)
		key := strings.ToLower(u.Name)

		keys[i] = key
		mkeys[key] = u
	}

	sort.Strings(keys)

	for _, k := range keys {
		u := mkeys[k]
		m := yaml.MapItem{
			Key:   u.Name,
			Value: u,
		}

		slice = append(slice, m)
	}

	return slice
}

// Write all data into one file
func (store *Store) WriteTo(data *types.Library, out io.Writer) (int64, error) {
	return WriteTo(store.ExportSlice(data, ExportDeep), out)
}
