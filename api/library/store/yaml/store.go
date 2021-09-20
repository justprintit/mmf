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
	const merge = false

	check := &errors.ErrorStack{}
	data := &types.Library{}

	// for each user on Basedir
	files, err := filepath.Glob(filepath.Join(store.Basedir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	n := len(files)
	for i, filename := range files {
		log.Printf("File[%v/%v]: %q", i, n, filename)
		if u, err := store.ReadUserFile(filename); err != nil {
			check.AppendWrapped(err, "LoadUserFile")
		} else if _, err := data.AddUser(u, merge); err != nil {
			check.AppendWrapped(err, "AddUser(%q)", u.Username)
		}
	}

	if check.Ok() {
		return data, nil // typed nil != nil
	} else {
		return data, check
	}
}

func (store *Store) Store(data *types.Library) error {
	var check errors.ErrorStack

	// Users
	for _, v := range data.User {
		base := store.Basedir

		if err := store.writeUser(base, v); err != nil {
			check.AppendError(err)
		}
	}

	if !check.Ok() {
		return &check
	}

	return nil
}

func (store *Store) writeUser(base string, user *types.User) error {
	// prepare data
	u, err := store.ExportUser(user)
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
