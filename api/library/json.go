package library

import (
	"github.com/justprintit/mmf/api/library/json"
)

func (c *Client) GetSharedLibrary() (*json.Users, error) {
	resp, err := c.GetLibrary("shared")
	if err != nil {
		return nil, err
	}

	out := &json.Users{}
	err = json.NewDecoderBytes(resp.Body()).Decode(out)
	if err != nil {
		return nil, err
	} else {
		return out, nil
	}
}
