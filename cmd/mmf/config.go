package main

import (
	"io"
	"log"

	"go.sancus.dev/config/flags"
	"go.sancus.dev/config/yaml"

	"github.com/justprintit/mmf/api"
)

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
