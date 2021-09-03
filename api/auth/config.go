package auth

// register generated ClientId at https://www.myminifactory.com/settings/developer/application
// ClientSecret then arrives by email "You Have been Authorised to use the MyMiniFactory API"
//
type Config struct {
	ClientId     string `hcl:"client_key,optional"`
	ClientSecret string `hcl:"client_secret,optional"`
}
