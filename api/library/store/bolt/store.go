package bolt

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/timshannon/bolthold"

	"github.com/justprintit/mmf/types"
)

const (
	DatabaseFilename = "mmf.db"
	DatabaseFileMode = os.FileMode(0644)
)

type BoltStore struct {
	mu   sync.Mutex
	bh   *bolthold.Store
	data types.Library
}

func New(datadir string, ev types.LibraryEvents) (types.Store, error) {
	filename := filepath.Join(datadir, DatabaseFilename)

	bh, err := bolthold.Open(filename, DatabaseFileMode, nil)
	if err != nil {
		return nil, err
	}

	store := &BoltStore{
		bh: bh,
		data: types.Library{
			Events: ev,
		},
	}

	return store, nil
}

func (m *BoltStore) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.bh.Close()
}

func (m *BoltStore) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
