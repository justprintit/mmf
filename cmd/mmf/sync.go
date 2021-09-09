package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/juju/persistent-cookiejar"
	"github.com/motemen/go-loghttp"
	"github.com/spf13/cobra"

	"golang.org/x/net/publicsuffix"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/library"
	"github.com/justprintit/mmf/api/library/json"
)

type Sync struct {
	sync.Mutex

	Cookies *cookiejar.Jar
	Config  *Config
}

func (m *Sync) Run() error {
	// client
	client := library.NewWithTransport(m.Config.Auth, &loghttp.Transport{
		LogRequest:  m.LogRequest,
		LogResponse: m.LogResponse,
	})
	client.SetCookieJar(m.Cookies)
	client.TraceEnabled = true

	// TODO: handle SIGTERM

	defer m.Save() // Save cookies

	if data, err := client.GetSharedLibrary(); err != nil {
		return errors.Wrap(err, "GetSharedLibrary")
	} else if err = json.Write(data, "  ", os.Stdout); err != nil {
		return err
	}

	if data, err := client.GetPurchasesLibrary(); err != nil {
		return errors.Wrap(err, "GetPurchasesLibrary")
	} else if err = json.Write(data, "  ", os.Stdout); err != nil {
		return err
	}

	return nil
}

func (m *Sync) LogRequest(req *http.Request) {
	loghttp.DefaultLogRequest(req)
}

func (m *Sync) LogResponse(resp *http.Response) {
	loghttp.DefaultLogResponse(resp)
}

func (m *Sync) Save() {
	m.Lock()
	defer m.Unlock()

	if len(cfg.Cookies) > 0 {
		m.Cookies.Save()
	}
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "synchronise library locally",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cfg.Setup()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// cookiejar
		jar, err := cookiejar.New(&cookiejar.Options{
			Filename:         cfg.Cookies,
			PublicSuffixList: publicsuffix.List,
		})
		if err != nil {
			return err
		}

		sync := &Sync{
			Cookies: jar,
			Config:  cfg,
		}

		return sync.Run()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
