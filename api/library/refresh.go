package library

import (
	"context"
	"log"

	"github.com/justprintit/mmf/api/library/json"
)

// Reload persistent Library data
func (c *Client) Reload() error {
	return nil
}

func (c *Client) refreshSharedLibrary(ctx context.Context, offset int, users ...json.User) error {
	i := offset
	for _, u := range users {
		i++

		if u.Name != u.Id {
			log.Printf("User.%v: %s (%s)", i, u.Name, u.Id)
		} else {
			log.Printf("User.%v: %s", i, u.Id)
		}

	}
	return nil
}

func (c *Client) refreshPurchasesLibrary(ctx context.Context, offset int, objects ...json.Object) error {
	i := offset
	for _, obj := range objects {
		i++
		log.Printf("%s.%v: %q (%v, %q)", "Purchase", i, obj.Name, obj.Id, obj.ObjType)
	}
	return nil
}

func (c *Client) refreshPledgesLibrary(ctx context.Context, offset int, objects ...json.Object) error {
	i := offset
	for _, obj := range objects {
		i++
		log.Printf("%s.%v: %q (%v, %q)", "Pledge", i, obj.Name, obj.Id, obj.ObjType)
	}
	return nil
}

func (c *Client) refreshTribesLibrary(ctx context.Context, offset int, tribes ...json.Tribe) error {
	i := offset
	for _, u := range tribes {
		i++

		log.Printf("Tribe.%v: %s (%v)", i, u.Name, u.Id)
	}
	return nil
}
