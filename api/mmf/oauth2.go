package mmf

import (
	"golang.org/x/oauth2"
)

// Supported OAuth 2.0 Scopes
type Scope string

const (
	BasicScope    Scope = "basic"
	DownloadScope Scope = "download"
)

// Endpoint is MMF's OAuth 2.0 endpoint
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://auth.myminifactory.com/web/authorize",
	TokenURL:  "https://auth.myminifactory.com/v1/oauth/tokens",
	AuthStyle: oauth2.AuthStyleInHeader,
}
