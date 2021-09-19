package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/juju/persistent-cookiejar"
	"github.com/motemen/go-loghttp"
	"github.com/spf13/cobra"

	"golang.org/x/net/publicsuffix"

	"github.com/justprintit/mmf/api/library"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "downloads something from MMF",
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

		// client
		client, err := library.NewWithOptions(
			library.WithCredentials(cfg.Auth),
			library.WithTransport(&loghttp.Transport{}),
			library.WithCookieJar(jar),
			library.WithDataDir(cfg.Data),
		)
		if err != nil {
			return err
		}

		// cancel
		ctx, cancel := context.WithCancel(context.Background())

		// watch signals
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		defer close(sig)

		done := make(chan error)
		go func() {
			defer close(done)

			for _, s := range args {
				resp, err := client.GetWithContext(ctx, s)
				if err != nil {
					done <- err
					break
				}
				os.Stdout.Write(resp.Body())
			}
		}()

		// and wait
		for {
			select {
			case signum := <-sig:
				// signal received
				switch signum {
				case syscall.SIGHUP:
					// ignore
				case syscall.SIGINT, syscall.SIGTERM:
					log.Println("Terminating...")
					cancel()
				}
			case err, _ := <-done:
				// client terminated
				return err
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
