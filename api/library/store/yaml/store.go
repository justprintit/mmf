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
		OnError: store.logError,
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
			check.AppendWrapped(err, "AddUser(%q)", u.Id())
		} else {
			// Groups and Tribes
			base := filepath.Join(store.Basedir, u.Name())
			if err := store.loadNodes(data, u, nil, base); err != nil {
				check.AppendWrapped(err, "LoadNodes(%q)", u.Id())
			}
		}
	}

	// enable all event loggers
	data.SetEvents(types.LibraryEvents{
		OnNewNode:    store.logNewNode,
		OnNodeUpdate: store.logNodeUpdate,
		OnError:      store.logError,
	})

	if err := check.AsError(); err != nil {
		data.OnError(nil, err)
		return data, err
	}
	return data, nil
}

func (store *Store) Store(data *types.Library) error {
	var check errors.ErrorStack

	// Users
	for _, id := range data.Users() {

		base := store.Basedir

		if v, err := data.GetUser(id); err != nil {
			check.AppendError(err)
		} else if err := store.writeUser(base, v); err != nil {
			check.AppendError(err)
		} else if v.HasNodes() {
			base = filepath.Join(base, v.Name())

			if err := os.MkdirAll(base, StoreDirectoryMode); err != nil {
				check.AppendError(err)
			} else {
				for _, g := range v.Nodes() {
					if err := store.writeNodes(base, g); err != nil {
						check.AppendError(err)
					}
				}
			}
		}
	}

	return check.AsError()
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
	u := &User{}
	adaptor := NewDecoder(f)
	if err = adaptor.Decode(u); err != nil {
		return nil, err
	}

	if u.Name == "" {
		// use filename as user name
		base := strings.TrimSpace(filepath.Base(filename))
		u.Name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	user := types.NewUser(u.Username, u.Name)
	user.Avatar = u.Avatar
	return user, nil
}

func (store *Store) loadNodes(data *types.Library, u *types.User, parent types.Node, base string) error {
	var check errors.ErrorStack

	const unique = false

	// for each group on base
	files, err := filepath.Glob(filepath.Join(base, "*.yaml"))
	if err != nil {
		return err
	}

	for _, filename := range files {
		if g, err := store.ReadNodeFile(filename); err != nil {
			check.AppendWrapped(err, "ReadNodeFile")
		} else if g != nil {

			if parent != nil {
				g, err = parent.AddNode(g, unique)
			} else {
				g, err = u.AddNode(g, unique)
			}

			if err == nil {
				subdir := filepath.Join(base, g.Name())

				if err = store.loadNodes(data, u, g, subdir); err != nil {
					check.AppendWrapped(err, "LoadNodes(%v)", g.Id)
				}
			} else {
				check.AppendWrapped(err, "AddGroup(%v)", g.Id)
			}
		}
	}

	return check.AsError()
}

func (store *Store) ReadNodeFile(filename string) (types.Node, error) {
	// open file for reading
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// decode
	g := &Group{}
	adaptor := NewDecoder(f)
	if err = adaptor.Decode(g); err != nil {
		return nil, err
	}

	if g.Name == "" {
		// use filename as group name
		base := strings.TrimSpace(filepath.Base(filename))
		g.Name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	switch g.Type {
	case "":
		group := types.NewGroup(g.Id.String(), g.Name)
		return group, nil
	case "tribe":
		tribe := types.NewTribe(g.Id.String(), g.Name)
		return tribe, nil
	default:
		err := errors.ErrInvalidValue("%s.%s: %q", "Group", "Type", g.Type)
		return nil, err
	}
}

func (store *Store) writeNodes(base string, group types.Node) error {
	// prepare data
	g, err := store.ExportNode(group, ExportShallow)
	if err != nil || len(g.Name) == 0 {
		return err
	}

	// but the Group name isn't to be encoded in the files
	name := g.Name
	g.Name = ""

	// output file
	filename := filepath.Join(base, name+".yaml")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, StoreFileMode)
	if err != nil {
		return err
	}

	// encode, without Name field
	adapter := NewEncoder(f)
	err = adapter.Encode(g)

	adapter.Close()
	f.Close()

	if err != nil {
		return err
	}

	// subgroups
	if group.HasNodes() {
		var check errors.ErrorStack

		base = filepath.Join(base, name)
		if err := os.MkdirAll(base, StoreDirectoryMode); err != nil {
			check.AppendError(err)
		} else {
			for _, sg := range group.Nodes() {
				if err = store.writeNodes(base, sg); err != nil {
					check.AppendError(err)
				}
			}
		}

		return check.AsError()
	}

	return nil
}
