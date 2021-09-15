package types

type Store interface {
	Load() (*Library, error)
	Store(*Library) error
}
