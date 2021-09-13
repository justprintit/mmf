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

func (c *Client) RefreshLibraries(ctx context.Context) error {
	// load persistent data
	if err := c.Reload(); err != nil {
		return err
	}

	// shared with me
	c.Add(func(c *Client, ctx context.Context) error {
		return c.RefreshSharedLibrary(ctx)
	})
	// purchases
	c.Add(func(c *Client, ctx context.Context) error {
		return c.RefreshPurchasesLibrary(ctx)
	})
	// pledges
	c.Add(func(c *Client, ctx context.Context) error {
		return c.RefreshPledgesLibrary(ctx)
	})
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

func (c *Client) RefreshSharedLibrary(ctx context.Context) error {
	var p *Pagination

	page := 1
	ok := true

	for ok {
		// get json items
		d, err := c.GetSharedLibraryPage(ctx, page)
		if err != nil {
			return err
		}

		if p == nil {
			//  first page
			p = c.Pages(len(d.User), d.Count)
		}

		// process in parallel
		offset := p.Size * (page - 1)
		c.Spawn(func(c *Client, ctx context.Context) error {
			return c.refreshSharedLibrary(ctx, offset, d.User...)
		})

		// next page
		page, _, ok = p.Next(page)
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

func (c *Client) RefreshPurchasesLibrary(ctx context.Context) error {
	var p *Pagination

	page := 1
	ok := true

	for ok {
		// get json items
		d, err := c.GetPurchasesLibraryPage(ctx, page)
		if err != nil {
			return err
		}

		if p == nil {
			//  first page
			p = c.PagesN(len(d.Items), d.Count)
		}

		// process in parallel
		offset := p.Size * (page - 1)
		c.Spawn(func(c *Client, ctx context.Context) error {
			return c.refreshPurchasesLibrary(ctx, offset, d.Items...)
		})

		// next page
		page, _, ok = p.Next(page)
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

func (c *Client) RefreshPledgesLibrary(ctx context.Context) error {
	var p *Pagination

	page := 1
	ok := true

	for ok {
		// get json items
		d, err := c.GetPledgesLibraryPage(ctx, page)
		if err != nil {
			return err
		}

		if p == nil {
			//  first page
			p = c.PagesN(len(d.Items), d.Count)
		}

		// process in parallel
		offset := p.Size * (page - 1)
		c.Spawn(func(c *Client, ctx context.Context) error {
			return c.refreshPledgesLibrary(ctx, offset, d.Items...)
		})

		// next page
		page, _, ok = p.Next(page)
	}

	return nil
}
