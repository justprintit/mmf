package store

import (
	"github.com/justprintit/mmf/types"
)

type DummyStore struct {
	types.Library
}

func NewDummy() (types.Store, error) {
	return &DummyStore{}, nil
}

func (m *DummyStore) Load() error {
	return nil
}

func (m *DummyStore) Close() error {
	return nil
}
