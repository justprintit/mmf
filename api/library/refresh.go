package library

import (
	"context"
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

func (c *Client) RefreshSharedLibrary(ctx context.Context) error {
	_, err := c.GetSharedLibrary(ctx)
	return err
}

func (c *Client) RefreshPurchasesLibrary(ctx context.Context) error {
	_, err := c.GetSharedLibrary(ctx)
	return err
}

func (c *Client) RefreshPledgesLibrary(ctx context.Context) error {
	_, err := c.GetPledgesLibrary(ctx)
	return err
}
