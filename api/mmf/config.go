package mmf

import (
	"golang.org/x/oauth2"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/util"
)

// Plain User credentials are used by the scrapper, while Client is used by the old oauth2 API
type User struct {
	Username string
	Password string
}

// for accessing the oauth2 API you first need to register
// a generated ClientID at https://www.myminifactory.com/settings/dev
// ClientSecret then arrives by email "You Have been Authorised to use the MyMiniFactory API"
type Client struct {
	ClientID     string `yaml:"client_key"`
	ClientSecret string `yaml:"client_secret"`

	AccessToken  string `yaml:"access_token,omitempty"`
	RefreshToken string `yaml:"refresh_token,omitempty"`
}

// Ok() checks if the Client can be used to compose an oauth2.Config
func (c *Client) Ok() bool {
	if len(c.ClientID) > 0 && len(c.ClientSecret) > 0 {
		return true
	}
	return false
}

// NewOauth2 creates a new oauth2.Config for use MMF's OpenAPI
func (c *Client) NewOauth2(callback string) (*oauth2.Config, error) {
	var (
		id     string
		secret string
		check  errors.ErrorStack
	)

	if callback == "" {
		check.MissingArgument("callback")
	}

	if id = c.ClientID; id == "" {
		check.MissingField("client_key")
	}

	if secret = c.ClientSecret; secret == "" {
		check.MissingField("client_secret")
	}

	if err := check.AsError(); err != nil {
		return nil, err
	}

	cfg := &oauth2.Config{
		RedirectURL:  callback,
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     Endpoint,
	}

	return cfg, nil
}

// Empty() checks if Client credentials aren't set
func (c *Client) Empty() bool {
	return len(c.ClientID) == 0
}

// Init() generates a new ClientID if none is set
func (c *Client) Init() error {
	if len(c.ClientID) == 0 {
		id, err := util.RandomString(16)
		if err != nil {
			return err
		}
		c.ClientID = id
	}
	return nil
}
