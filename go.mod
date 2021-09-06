module github.com/justprintit/mmf

go 1.16

require (
	github.com/frankban/quicktest v1.13.1 // indirect
	github.com/go-resty/resty/v2 v2.6.0
	github.com/json-iterator/go v1.1.11
	github.com/juju/go4 v0.0.0-20160222163258-40d72ab9641a // indirect
	github.com/juju/persistent-cookiejar v0.0.0-20171026135701-d5e5a8405ef9
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/motemen/go-loghttp v0.0.0-20170804080138-974ac5ceac27
	github.com/motemen/go-nuts v0.0.0-20210915132349-615a782f2c69 // indirect
	github.com/sethvargo/go-password v0.2.0
	github.com/spf13/cobra v1.2.1
	go.sancus.dev/config v0.6.1
	go.sancus.dev/core v0.16.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/retry.v1 v1.0.3 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	go.sancus.dev/config => ../../../go.sancus.dev/config
	go.sancus.dev/core => ../../../go.sancus.dev/core
	go.sancus.dev/web => ../../../go.sancus.dev/web
)
