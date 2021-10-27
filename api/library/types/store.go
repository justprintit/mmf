package types

type Store interface {
	Load() error
	Close() error
}
