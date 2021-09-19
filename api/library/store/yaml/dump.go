package yaml

import (
	"io"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/justprintit/mmf/api/library/types"
)

// Export User data for YAML encoding
func (store *Store) ExportUser(w *types.User) (*User, error) {
	name := strings.TrimSpace(w.Name)
	name = strings.ReplaceAll(name, string(os.PathSeparator), "-")

	if len(name) == 0 {
		name = w.Username
	}

	u := &User{
		Username: w.Username,
		Name:     name,
		Avatar:   w.Avatar,
	}

	return u, nil
}

// Export Library for YAML encoding
func (store *Store) ExportSlice(data *types.Library) yaml.MapSlice {
	n := len(data.User)                 // users count
	keys := make([]string, 0, n)        // for sorting, lower case names
	mkeys := make(map[string]*User, n)  // lower case name to user
	slice := make([]yaml.MapItem, 0, n) // return value

	for _, w := range data.User {
		u, _ := store.ExportUser(w)
		key := strings.ToLower(u.Name)

		keys = append(keys, key)
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
	return WriteTo(store.ExportSlice(data), out)
}
