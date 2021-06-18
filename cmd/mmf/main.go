package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mmf",
	Short: "mmf gives command line access to MyMiniFactory.com",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
