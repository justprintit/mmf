package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/juju/persistent-cookiejar"
	"github.com/motemen/go-loghttp"
	"github.com/spf13/cobra"

	"golang.org/x/net/publicsuffix"

	"github.com/justprintit/mmf/api/library"
)

type Sync struct {
	*library.Client
	Cookies *cookiejar.Jar
	Config  *Config
}

func (m *Sync) Init() error {
	// client
	client, err := library.NewWithOptions(
		library.WithCredentials(m.Config.Auth),
		library.WithTransport(&loghttp.Transport{
			LogRequest:  m.LogRequest,
			LogResponse: m.LogResponse,
		}),
		library.WithCookieJar(m.Cookies),
		library.WithDataDir(m.Config.Data),
	)

	if err != nil {
		return err
	}

	m.Client = client
	m.Client.TraceEnabled = true

	// setup
	m.Add(func(c *library.Client, ctx context.Context) error {
		return c.RefreshLibraries(ctx)
	})

	return nil
}

func (m *Sync) LogRequest(req *http.Request) {
	loghttp.DefaultLogRequest(req)
}

func (m *Sync) LogResponse(resp *http.Response) {
	loghttp.DefaultLogResponse(resp)
}

func (m *Sync) Save() {
	if len(cfg.Cookies) > 0 {
		m.Cookies.Save()
	}
	log.Println("Done.")
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

		// prepare sync client
		sync := &Sync{
			Cookies: jar,
			Config:  cfg,
		}

		if err := sync.Init(); err != nil {
			return err
		}

		// watch signals
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		defer close(sig)

		log.Println("Starting...")
		defer sync.Save() // Save cookies
		sync.Start()

		// and wait
		for {
			select {
			case signum := <-sig:
				// signal received
				switch signum {
				case syscall.SIGHUP:
					if err := sync.Reload(); err != nil {
						log.Println("Reload failed:", err)
					}
				case syscall.SIGINT, syscall.SIGTERM:
					log.Println("Terminating...")
					sync.Cancel()
				}
			case <-sync.Done():
				// client terminated
				return nil
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
