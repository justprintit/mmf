package mmf

type Credentials struct {
	Username string `hcl:"username,optional"`
	Password string `hcl:"password,optional"`
}

// register generated ClientId at https://www.myminifactory.com/settings/developer/application
// ClientSecret then arrives by email "You Have been Authorised to use the MyMiniFactory API"
//
type Config struct {
	ClientId     string `hcl:"client_key,optional"`
	ClientSecret string `hcl:"client_secret,optional"`

	Token []TokenConfig `hcl:"token,block"`
}

type TokenConfig struct {
	Username     string `hcl:"username,label"`
	Type         string `hcl:"token_type,label"`
	AccessToken  string `hcl:"access_token,optional"`
	RefreshToken string `hcl:"refresh_token,optional"`
	Expiry       string `hcl:"expiry,optional"`
}
