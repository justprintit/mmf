package transport

import (
	"golang.org/x/oauth2"
)

type ClientEvents struct {
	OnNewCredentials func(user, password string) error
	OnNewClient      func(key, secret string) error
	OnNewToken       func(access, renew string) error
}

func (c *Client) onNewCredentials(user, password string) error {
	if c != nil && c.events.OnNewCredentials != nil {
		return c.events.OnNewCredentials(user, password)
	}
	return nil
}

func (c *Client) onNewClient(key, secret string) error {
	if c != nil && c.events.OnNewClient != nil {
		return c.events.OnNewClient(key, secret)
	}
	return nil
}

func (c *Client) onNewToken(token *oauth2.Token) error {
	if c != nil && c.events.OnNewToken != nil {
		var access string
		var renew string

		if token != nil {
			access = token.AccessToken
			renew = token.RefreshToken
		}

		return c.events.OnNewToken(access, renew)
	}
	return nil
}
