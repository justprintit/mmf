package transport

type ClientOptionFunc func(*Client) error

func (f ClientOptionFunc) Apply(c *Client) error {
	return f(c)
}

type ClientOption interface {
	Apply(*Client) error
}

func NewClientWithOptions(options ...ClientOption) (*Client, error) {
	c := &Client{}

	for _, opt := range options {
		if err := opt.Apply(c); err != nil {
			return nil, err
		}
	}

	if err := c.SetDefaults(); err != nil {
		return nil, err
	}
	return c, nil
}
