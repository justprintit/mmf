package json

import (
	"fmt"
	"log"
	"strings"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/library/types"
)

type Tribes struct {
	Count int     `json:"total_count,omitempty"`
	Items []Tribe `json:"items,omitempty"`
}

type Tribe struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	UserAvatar string      `json:"user_avatar"`
	URL        string      `json:"url"`
	Groups     TribeGroups `json:"groups,omitempty"`
}

type TribeGroups struct {
	Count int          `json:"total_count,omitempty"`
	Items []TribeGroup `json:"items,omitempty"`
}

type TribeGroup struct {
	Id           types.Id `json:"id"`
	Name         string   `json:"name"`
	TotalObjects int      `json:"total_count_objects,omitempty"`
	Date         string   `json:"date,omitempty"`
}

func (w *Tribe) Apply(d *types.Library) (*types.Tribe, error) {
	const merge = true

	// Tribe's owner
	username, err := types.UsernameFromURL(w.URL)
	if err != nil {
		return nil, errors.Wrap(err, "TribeUser")
	}

	u := types.NewUser(username, "")
	u.Avatar = w.UserAvatar

	u, err = d.AddUser(u, merge)
	if err != nil {
		return nil, errors.Wrap(err, "AddUser")
	}

	// Tribe
	tribe := types.NewTribe(fmt.Sprintf("%v", w.Id), w.Name)

	tribe, err = u.AddTribe(tribe, merge)
	if err != nil {
		return nil, errors.Wrap(err, "AddTribe")
	}

	// Tribe Groups
	err = w.Groups.Apply(tribe)
	return tribe, err
}

func (w *TribeGroups) Apply(tribe *types.Tribe) error {
	var check errors.ErrorStack

	// remove `all/n` group
	groups := make([]*types.Group, 0, len(w.Items))
	for _, tg := range w.Items {
		id := tg.Id.String()

		if !strings.HasPrefix(id, "all/") {
			g := types.NewGroup(tg.Id.String(), tg.Name)
			groups = append(groups, g)
		}
	}

	if len(groups) != w.Count {
		log.Printf("%s.%s: expected:%v != actual:%v", "Tribe", "Groups", w.Count, len(groups))
	}

	for _, g := range groups {
		if _, err := tribe.AddGroup(g, true); err != nil {
			check.AppendError(err)
		}
	}

	return check.AsError()
}
