package auth

import (
	"golang.org/x/oauth2"

	"go.sancus.dev/core/errors"
)

// Endpoint is MMF's OAuth 2.0 endpoint
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://auth.myminifactory.com/web/authorize",
	TokenURL:  "https://auth.myminifactory.com/v1/oauth/tokens",
	AuthStyle: oauth2.AuthStyleInHeader,
}

// Populates a oauth2.Config
func (c Config) NewOauth2(callback string) (*oauth2.Config, error) {
	var (
		id     string
		secret string
		check  errors.ErrorStack
	)

	if len(callback) == 0 {
		check.MissingArgument("callback")
	}

	if id = c.ClientId; len(id) == 0 {
		check.MissingField("client_id")
	}

	if secret = c.ClientSecret; len(secret) == 0 {
		check.MissingField("client_secret")
	}

	if !check.Ok() {
		return nil, &check
	}

	cfg := &oauth2.Config{
		RedirectURL:  callback,
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     Endpoint,
	}

	return cfg, nil
}
