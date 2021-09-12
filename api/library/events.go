package library

func (c *Client) Run() error {
	c.Start()

	return <-c.Done()
}

func (c *Client) Reload() error {
	return nil
}
