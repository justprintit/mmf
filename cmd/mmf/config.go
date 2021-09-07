package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

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
	} else if err := os.MkdirAll(c.Data, 0755); err != nil {
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
