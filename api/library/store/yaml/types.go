package yaml

import (
	"github.com/justprintit/mmf/api/library/types"
)

type User struct {
	Username string
	Name     string  `yaml:"-"`
	Avatar   string  `yaml:",omitempty"`
	Groups   []Group `yaml:",omitempty"`
}

type Group struct {
	Id        types.Id
	Name      string  `yaml:",omitempty"`
	Subgroups []Group `yaml:",omitempty"`
}
