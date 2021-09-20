package store

import (
	"github.com/justprintit/mmf/api/library/store/yaml"
	"github.com/justprintit/mmf/api/library/types"
)

type (
	YAMLStore = yaml.Store
)

type NOPStore struct{}

func (nop *NOPStore) Load() (*types.Library, error) {
	return &types.Library{}, nil
}

func (nop *NOPStore) Store(l *types.Library) error {
	return nil
}
