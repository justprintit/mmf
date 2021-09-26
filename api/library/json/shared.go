package json

import (
	"github.com/justprintit/mmf/api/library/types"
)

type UserSharedLibrary struct {
	Objects Objects `json:"objects,omitempty"`
	Groups  Groups  `json:"groups,omitempty"`
}

func (p *UserSharedLibrary) Apply(d *types.Library, u *types.User) error {
	if err := p.Groups.Apply(d, u); err != nil {
		return err
	}
	if err := p.Objects.Apply(d, u, nil); err != nil {
		return err
	}
	return nil
}
