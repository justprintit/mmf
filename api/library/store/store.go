package store

import (
	"github.com/justprintit/mmf/api/library/types"
)

type NOPStore struct{}

func (nop *NOPStore) Load() (*types.Library, error) {
	return &types.Library{}, nil
}

func (nop *NOPStore) Store(l *types.Library) error {
	return nil
}

type YAMLStore struct {
	Basedir string
}

func (store *YAMLStore) Load() (*types.Library, error) {
	return &types.Library{}, nil
}

func (store *YAMLStore) Store(data *types.Library) error {
	return nil
}
