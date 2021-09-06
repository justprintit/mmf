package main

import (
	"io"
	"log"
	"os"

	"go.sancus.dev/config/flags"
	"go.sancus.dev/config/hcl"

	"github.com/justprintit/mmf"
)

type Config struct {
	Auth    mmf.Credentials `hcl:"auth,block"`
	App     mmf.Config      `hcl:"app,block"`
	Data    string          `hcl:"data_dir,optional"`
	Cookies string          `hcl:"cookie_jar,optional" default:"cookies.json"`
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

func (c *Config) Save(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		_, err = c.WriteTo(f)
	}
	return err
}
