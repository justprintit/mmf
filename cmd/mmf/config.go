package main

import (
	"io"
	"log"

	"go.sancus.dev/config/flags"
	"go.sancus.dev/config/hcl"

	"github.com/justprintit/mmf"
)

type Config struct {
	Auth mmf.Config `hcl:"auth,block"`
}

func NewConfig() *Config {
	c := &Config{}

	if err := flags.SetDefaults(c); err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) ReadInFile(filename string) error {
	return hcl.LoadFile(filename, nil, c)
}

func (c *Config) WriteTo(out io.Writer) (int, error) {
	return hcl.WriteTo(out, c)
}
