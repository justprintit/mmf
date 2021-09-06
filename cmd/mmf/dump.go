package main

import (
	"os"

	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "prints current config",
	RunE: func(cmd *cobra.Command, args []string) error {

		_, err := cfg.WriteTo(os.Stdout)
		return err
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
