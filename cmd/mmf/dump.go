package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/justprintit/mmf/api/library/store/yaml"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "prints current data",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cfg.Setup()
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		store := &yaml.Store{
			Basedir: cfg.Data,
		}

		library, err := store.Load()
		if err == nil {
			_, err = store.WriteTo(library, os.Stdout)
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
