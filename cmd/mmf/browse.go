package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/juju/persistent-cookiejar"
	"github.com/pkg/browser"
	"golang.org/x/net/publicsuffix"

	"go.sancus.dev/config/flags"
	"go.sancus.dev/config/flags/cobra"
	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/transport"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse runs a MMF browser",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flags.GetMapper(cmd.Flags()).Parse()
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

		// client
		_, err = transport.NewClientWithOptions(
			transport.WithCookieJar(jar),
		)
		if err != nil {
			return err
		}

		// listen
		srv, err := cfg.Server.NewServer()
		if err != nil {
			return err
		}
		done := make(chan error)

		// watch signals
		go func() {
			defer close(done)

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

			for signum := range sig {
				switch signum {
				case syscall.SIGHUP:
					// ignore
				case syscall.SIGINT, syscall.SIGTERM:
					// terminate
					log.Println("Terminating...")
					return
				}
			}
		}()

		// launch server
		go srv.Serve()

		// lauch browser
		url := srv.URL().String()
		err = browser.OpenURL(url)
		if err != nil {
			err = errors.Wrap(err, "OpenURL: %q", url)
		} else {
			// and wait until we are done
			<-done
		}
		srv.Shutdown(context.Background())
		return err
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)

	cobra.NewMapper(browseCmd.Flags()).
		VarP(&cfg.Server.Port, "port", 'p', "HTTP port to listen")
}
