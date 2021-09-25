package library

import (
	"context"
	"log"

	"github.com/justprintit/mmf/api/library/json"
)

// Reload persistent Library data
func (c *Client) Reload() error {
	l, err := c.store.Load()
	if err == nil {
		c.library = l
	}
	return err
}

// Stores library data persistently
func (c *Client) Commit() error {
	return c.store.Store(c.library)
}

func (c *Client) refreshSharedLibrary(ctx context.Context, offset int, users ...json.User) error {
	for _, u := range users {
		if err := u.Apply(c.library); err != nil {
			log.Println(err)
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
