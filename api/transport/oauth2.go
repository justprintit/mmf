package transport

import (
	"net/url"
	"path"

	"github.com/justprintit/mmf/api/mmf"
)

func WithOauth2(cred mmf.Client, base string, callback string) ClientOptionFunc {
	// callback = base + path
	if u, err := url.Parse(base); err != nil {
		panic(err)
	} else {
		u.Path = path.Clean(path.Join(u.Path, callback))
		callback = u.String()
	}

	return func(c *Client) error {
		c.callback = callback
		c.client = cred
		return nil
	}
}
