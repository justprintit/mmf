package api

import (
	"github.com/justprintit/mmf/util"
)

// Plain Credentials are used by the scrapper, while Client is used by the old oauth2 API
type Credentials struct {
	Username string
	Password string
}

// for accessing the oauth2 API you first need to register
// a generated ClientId at https://www.myminifactory.com/settings/dev
// ClientSecret then arrives by email "You Have been Authorised to use the MyMiniFactory API"
type Client struct {
	ClientId     string `yaml:"client_key"`
	ClientSecret string `yaml:"client_secret"`

	AccessToken  string `yaml:"access_token,omitempty"`
	RefreshToken string `yaml:"refresh_token,omitempty"`
}

// Empty() checks if client credentials aren't set
func (c *Client) Empty() bool {
	return len(c.ClientId) == 0
}

// Init() generates a new ClientId if none is set
func (c *Client) Init() error {
	if len(c.ClientId) == 0 {
		id, err := util.RandomString(16)
		if err != nil {
			return err
		}
		c.ClientId = id
	}
	return nil
}
