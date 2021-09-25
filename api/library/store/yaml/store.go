package yaml

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/library/types"
)

const (
	StoreDirectoryMode fs.FileMode = 0755
	StoreFileMode      fs.FileMode = 0644
)

type Store struct {
	Basedir string
}

func (store *Store) Load() (*types.Library, error) {
	var check errors.ErrorStack

	const unique = false // !merge

	data := &types.Library{}

	// first only log errors
	data.SetEvents(types.LibraryEvents{
		OnError:      store.logError,
		OnUserError:  store.logUserError,
		OnGroupError: store.logGroupError,
	})

	log.Printf("Loading: %q", store.Basedir)

	// for each user on Basedir
	files, err := filepath.Glob(filepath.Join(store.Basedir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	for _, filename := range files {
		if u, err := store.ReadUserFile(filename); err != nil {
			check.AppendWrapped(err, "LoadUserFile")
		} else if u, err := data.AddUser(u, unique); err != nil {
			check.AppendWrapped(err, "AddUser(%q)", u.Username)
		} else {
			// Groups
			name := filepath.Base(filename)
			name = strings.TrimSuffix(name, filepath.Ext(name))

			base := filepath.Join(store.Basedir, name)
			if err := store.loadGroups(data, u, nil, base); err != nil {
				check.AppendWrapped(err, "LoadGroups(%q)", u.Username)
			}
		}
	}

	// enable all event loggers
	data.SetEvents(types.LibraryEvents{
		OnNewUser:  store.logNewUser,
		OnNewGroup: store.logNewGroup,

		OnUserUpdate:  store.logUserUpdate,
		OnGroupUpdate: store.logGroupUpdate,

		OnError:      store.logError,
		OnUserError:  store.logUserError,
		OnGroupError: store.logGroupError,
	})

	if !check.Ok() {
		err := &check
		data.OnError(err)
		return data, err
	}
	return data, nil
}

func (store *Store) Store(data *types.Library) error {
	var check errors.ErrorStack

	// Users
	for _, v := range data.User {
		base := store.Basedir

		if err := store.writeUser(base, v); err != nil {
			check.AppendError(err)
		}

		// Groups
		if len(v.Groups) > 0 {
			base = filepath.Join(base, v.Name)

			if err := os.MkdirAll(base, StoreDirectoryMode); err != nil {
				check.AppendError(err)
			} else {
				for _, g := range v.Groups {
					if err := store.writeGroups(base, g); err != nil {
						check.AppendError(err)
					}
				}
			}
		}
	}

	if !check.Ok() {
		return &check
	}

	return nil
}

func (store *Store) writeUser(base string, user *types.User) error {
	// prepare data
	u, err := store.ExportUser(user, ExportShallow)
	if err != nil || u == nil {
		return err
	}

	// output file
	filename := filepath.Join(base, u.Name+".yaml")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, StoreFileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	// encode
	adapter := NewEncoder(f)
	defer adapter.Close()

	return adapter.Encode(u)
}

func (store *Store) ReadUserFile(filename string) (*types.User, error) {
	// open file for reading
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// decode
	u := &types.User{}
	adaptor := NewDecoder(f)
	if err = adaptor.Decode(u); err != nil {
		return nil, err
	}

	if len(u.Name) == 0 {
		// use filename as user name
		base := strings.TrimSpace(filepath.Base(filename))
		u.Name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	return u, nil
}

func (store *Store) loadGroups(data *types.Library, u *types.User, parent *types.Group, base string) error {
	var check errors.ErrorStack

	const unique = false

	// for each group on base
	files, err := filepath.Glob(filepath.Join(base, "*.yaml"))
	if err != nil {
		return err
	}

	for _, filename := range files {
		if g, err := store.ReadGroupFile(filename); err != nil {
			check.AppendWrapped(err, "LoadGroupFile")
		} else if g != nil {
			if parent != nil {
				g, err = parent.AddSubgroup(g, unique)
			} else {
				g, err = u.AddGroup(g, unique)
			}

			if err != nil {
				check.AppendWrapped(err, "AddGroup(%v)", g.Id)
			}
		}
	}

	if !check.Ok() {
		return &check
	}
	return nil
}

func (store *Store) ReadGroupFile(filename string) (*types.Group, error) {
	// open file for reading
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// decode
	g := &types.Group{}
	adaptor := NewDecoder(f)
	if err = adaptor.Decode(g); err != nil {
		return nil, err
	}

	if len(g.Name) == 0 {
		// use filename as group name
		base := strings.TrimSpace(filepath.Base(filename))
		g.Name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	return g, nil
}

func (store *Store) writeGroups(base string, group *types.Group) error {
	// prepare data
	g, err := store.ExportGroup(group, ExportShallow)
	if err != nil || len(g.Name) == 0 {
		return err
	}

	// output file
	filename := filepath.Join(base, g.Name+".yaml")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, StoreFileMode)
	if err != nil {
		return err
	}

	// encode
	adapter := NewEncoder(f)
	err = adapter.Encode(g)

	adapter.Close()
	f.Close()

	if err != nil {
		return err
	}

	// subgroups
	if len(group.Subgroups) > 0 {
		var check errors.ErrorStack

		base = filepath.Join(base, g.Name)
		if err := os.MkdirAll(base, StoreDirectoryMode); err != nil {
			check.AppendError(err)
		} else {
			for _, sg := range group.Subgroups {
				if err = store.writeGroups(base, sg); err != nil {
					check.AppendError(err)
				}
			}
		}

		if !check.Ok() {
			return &check
		}
	}

	return nil
}
