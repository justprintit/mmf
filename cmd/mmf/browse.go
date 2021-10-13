package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/browser"

	"go.sancus.dev/config/flags"
	"go.sancus.dev/config/flags/cobra"
	"go.sancus.dev/core/errors"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse runs a MMF browser",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flags.GetMapper(cmd.Flags()).Parse()
		return cfg.Setup()
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		app, err := NewApp(*cfg, cfgFile)
		if err != nil {
			return err
		}

		// watch signals
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

			for signum := range sig {
				switch signum {
				case syscall.SIGHUP:
					if err := app.Reload(); err != nil {
						log.Println("Reload failed: %s", err)
					}
				case syscall.SIGINT, syscall.SIGTERM:
					// terminate
					log.Println("Terminating...")
					app.Abort()
					return
				}
			}
		}()

		// launch app
		go app.Run()

		// lauch browser
		url := app.URL()
		err = browser.OpenURL(url)
		if err != nil {
			err = errors.Wrap(err, "OpenURL: %q", url)
			app.Abort()
		}

		// and wait until we are done
		app.Wait()
		return err
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)

	cobra.NewMapper(browseCmd.Flags()).
		VarP(&cfg.Server.Port, "port", 'p', "HTTP port to listen")
}
