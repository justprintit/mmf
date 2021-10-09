package main

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"go.sancus.dev/config/flags"
	"go.sancus.dev/config/yaml"
	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api"
)

const DefaultConfigFileMode os.FileMode = 0600 // owner only, because we include the password

type AuthConfig struct {
	User   api.Credentials `yaml:"user"`
	Client api.Client      `yaml:"api"`
}

type Config struct {
	Auth AuthConfig
}

func NewConfig() *Config {
	c := &Config{}

	if err := flags.SetDefaults(c); err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) ReadInFile(filename string) error {
	return yaml.LoadFile(filename, c)
}

func (c *Config) WriteTo(out io.Writer) (int64, error) {
	return yaml.WriteTo(out, c)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: CmdName + " config manages the configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			write bool
			dump  bool
			check errors.ErrorStack
		)

		for _, x := range args {
			switch x {
			case "dump":
				dump = true
			default:
				check.InvalidArgument(x)
			}
		}

		if !cfg.Auth.Client.Empty() {
			// client_id set
		} else if err := cfg.Auth.Client.Init(); err != nil {
			// failed to generate client_id
			check.AppendError(err)
		} else {
			// client_id generated, rewrite config
			write = true
		}

		if cfgReadError == nil {
			// no error
		} else if os.IsNotExist(cfgReadError) {
			// no such file or directory, make one
			write = true
		} else {
			// other read error
			check.AppendError(cfgReadError)
		}

		if dump {
			_, err := cfg.WriteTo(os.Stdout)
			if err != nil {
				check.AppendError(err)
			}
		}

		if write {
			_, err := yaml.WriteFile(cfgFile, cfg, DefaultConfigFileMode)
			if err == nil {
				log.Println(cfgFile, "written")
			} else {
				check.AppendError(err)
			}
		}

		return check.AsError()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
