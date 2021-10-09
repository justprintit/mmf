package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"go.sancus.dev/config/flags"
	"go.sancus.dev/config/yaml"
	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api"
	"github.com/justprintit/mmf/web/server"
)

const (
	DefaultConfigFileMode    os.FileMode = 0600 // owner only, because we include the password
	DefaultDirectoryFileMode             = 0755
)

type AuthConfig struct {
	User   api.Credentials `yaml:"user"`
	Client api.Client      `yaml:"api"`
}

type Config struct {
	Auth   AuthConfig
	Server server.ServerConfig `yaml:",omitempty"`

	Data    string `yaml:"data_dir,omitempty"`
	Cookies string `yaml:"cookies" default:"cookies.json"`
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

func (c *Config) Setup() error {
	// data directory
	c.Data = filepath.Clean(c.Data)
	if fi, err := os.Stat(c.Data); fi.IsDir() {
		// ready
	} else if err == nil {
		// exists, but not a directory
		return &os.PathError{
			Op:   fmt.Sprintf("%T.%s", c, "Setup"),
			Path: c.Data,
			Err:  syscall.ENOTDIR,
		}
	} else if err = os.MkdirAll(c.Data, DefaultDirectoryFileMode); err != nil {
		// failed to create
		return err
	}

	// CookieJar
	if len(c.Cookies) == 0 {
		flags.SetDefaults(c.Cookies)
	}
	if strings.IndexRune(c.Cookies, os.PathSeparator) == -1 {
		// no '/', place it inside `data_dir`
		c.Cookies = filepath.Join(c.Data, c.Cookies)
	}
	c.Cookies = filepath.Clean(c.Cookies)

	return nil
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
