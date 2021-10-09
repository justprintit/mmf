package main

import (
	"go.sancus.dev/config/flags"
	"go.sancus.dev/config/flags/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse runs a MMF browser",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flags.GetMapper(cmd.Flags()).Parse()
		return cfg.Setup()
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)

	cobra.NewMapper(browseCmd.Flags()).
		VarP(&cfg.Server.Port, "port", 'p', "HTTP port to listen")
}
